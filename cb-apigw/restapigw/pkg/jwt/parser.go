// Package jwt -
package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
	"github.com/dgrijalva/jwt-go"
)

// ===== [ Constants and Variables ] =====

var (
	// ErrSigningMethodMismatch - 다른 서명 메서드가 적용된 경우 오류
	ErrSigningMethodMismatch = errors.New("signing method mismatch")
	// ErrFailedToParseToken - 비밀키와 만기등을 검증하기 위한 파싱 처리에 실패한 경우 오류
	ErrFailedToParseToken = errors.New("failed to parse token")
	// ErrUnsupportedSigningMethod - 지원되지 않는 서명 메서드인 경우 오류
	ErrUnsupportedSigningMethod = errors.New("unsupported signing method")
	// ErrInvalidPEMBlock - PEM-encoded 가 아닌 키 구성인 경우 오류
	ErrInvalidPEMBlock = errors.New("invalid RSA: not PEM-encoded")
	// ErrNotRSAPublicKey - RSA 공개 키 형식이 아닌 경우 오류
	ErrNotRSAPublicKey = errors.New("invalid RSA: expected PUBLIC KEY block type")
	// ErrBadPublicKey - RSA 공개 키가 잘못된 경우 오류
	ErrBadPublicKey = errors.New("invalid RSA: failed to assert public key")
)

// ===== [ Types ] =====

type (
	// SigningMethod - 서명 처리에 필요한 알고리즘, 키 구조
	SigningMethod struct {
		// 알고리즘 (HS256, HS384, HS512, RS256, RS384, RS512)
		Alg string `yaml:"alg" json:"alg"`
		Key string `yaml:"key" json:"key"`
	}

	// ParserConfig - JWT 파서 및 토큰 검증 구조
	ParserConfig struct {
		SigningMethods []SigningMethod
		// Request 정보에서 토큰 추출을 위한 정보 "<source>:<name>" 형식 (ex. "header:Authorization", "query:name", "cookie:name")
		TokenLookup string
		// 만기 시간용
		Leeway int64
	}

	// Parser - Parser 설정 관리 구조
	Parser struct {
		Config ParserConfig
	}
)

// ===== [ Implementations ] =====

func (p *Parser) jwtFromHeader(req *http.Request, key string) (string, error) {
	authHeader := req.Header.Get(key)

	if "" == authHeader {
		return "", errors.New("auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(2 == len(parts) && "Bearer" == parts[0]) {
		return "", errors.New("invalid auth header")
	}

	return parts[1], nil
}

func (p *Parser) jwtFromQuery(req *http.Request, key string) (string, error) {
	token := req.URL.Query().Get(key)

	if "" == token {
		return "", errors.New("query token empty")
	}

	return token, nil
}

func (p *Parser) jwtFromCookie(req *http.Request, key string) (string, error) {
	cookie, _ := req.Cookie(key)

	if nil == cookie {
		return "", errors.New("cookie token empty")
	}

	return cookie.Value, nil
}

// Parse a JWT token and validates it
func (p *Parser) Parse(tokenString string) (*jwt.Token, error) {
	for _, method := range p.Config.SigningMethods {
		token, err := jwt.ParseWithClaims(tokenString, NewAppClaims(p.Config.Leeway), func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != method.Alg {
				return nil, ErrSigningMethodMismatch
			}

			switch token.Method.(type) {
			case *jwt.SigningMethodHMAC:
				return []byte(method.Key), nil
			case *jwt.SigningMethodRSA:
				block, _ := pem.Decode([]byte(method.Key))
				if nil == block {
					return nil, ErrInvalidPEMBlock
				}
				if got, want := block.Type, "PUBLIC KEY"; got != want {
					return nil, ErrNotRSAPublicKey
				}
				pub, err := x509.ParsePKIXPublicKey(block.Bytes)
				if nil != err {
					return nil, err
				}

				if _, ok := pub.(*rsa.PublicKey); !ok {
					return nil, ErrBadPublicKey
				}

				return pub, nil
			default:
				return nil, ErrUnsupportedSigningMethod
			}
		})

		if nil != err {
			if err == ErrSigningMethodMismatch {
				continue
			}

			if validationErr, ok := err.(*jwt.ValidationError); ok && (0 < validationErr.Errors&jwt.ValidationErrorUnverifiable || 0 < validationErr.Errors&jwt.ValidationErrorSignatureInvalid) {
				continue
			}
		}

		return token, err
	}

	return nil, ErrFailedToParseToken
}

// ParseFromRequest - Request에서 Token 추출 및 검증
func (p *Parser) ParseFromRequest(req *http.Request) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split(p.Config.TokenLookup, ":")
	switch parts[0] {
	case "header":
		token, err = p.jwtFromHeader(req, parts[1])
	case "query":
		token, err = p.jwtFromQuery(req, parts[1])
	case "cookie":
		token, err = p.jwtFromCookie(req, parts[1])
	}

	if nil != err {
		return nil, err
	}

	return p.Parse(token)
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====
