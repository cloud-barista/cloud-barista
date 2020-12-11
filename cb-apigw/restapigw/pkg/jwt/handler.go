// Package jwt -
package jwt

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/admin/response"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/jwt/provider"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

// ===== [ Constants and Variables ] =====

const (
	bearer = "bearer"
)

// ===== [ Types ] =====
type (
	// Handler - JWT 처리와 연계되는 구조
	Handler struct {
		Guard Guard
	}
)

// ===== [ Implementations ] =====

// Login - JWT Token 기반으로 인증 처리
func (h *Handler) Login(cc *config.CredentialsConfig, logger logging.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		accessToken, err := extractAccessToken(req)

		if nil != err {
			logger.WithError(err).Debug("[ADMIN API] failed to extract access token")
		}

		httpClient := getClient(accessToken)
		factory := provider.Factory{}
		p := factory.Build(req.URL.Query().Get("provider"), cc)

		verified, err := p.Verify(req, httpClient)

		if nil != err || !verified {
			response.Errorf(rw, req, http.StatusUnauthorized, err)
			return
		}

		if 0 == h.Guard.Timeout {
			h.Guard.Timeout = time.Hour
		}

		claims, err := p.GetClaims(httpClient)
		if nil != err {
			response.Errorf(rw, req, http.StatusBadRequest, err)
			return
		}

		token, err := IssueAdminToken(h.Guard.SigningMethod, claims, h.Guard.Timeout)
		if nil != err {
			response.Errorf(rw, req, http.StatusUnauthorized, err)
			return
		}

		response.Write(rw, req, token)
	}
}

// Logout - JWT Token 기반으로 인증 해제 처리
func (h *Handler) Logout(cc *config.CredentialsConfig, logger logging.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		accessToken, err := extractAccessToken(req)

		if nil != err || "" == accessToken {
			logger.WithError(err).Debug("[ADMIN API] failed to extract access token")
			response.Errorf(rw, req, http.StatusBadRequest, err)
			return
		}

		response.Write(rw, req, nil)
	}
}

// Refresh - 현재 JWT에 대한 Refresh 처리
func (h *Handler) Refresh() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		parser := Parser{h.Guard.ParserConfig}
		token, _ := parser.ParseFromRequest(req)
		claims := token.Claims.(*AppClaims).MapClaims

		origIat := int64(claims["iat"].(float64))

		if origIat < time.Now().Add(-h.Guard.MaxRefresh).Unix() {
			response.Errorf(rw, req, http.StatusUnauthorized, errors.New("token is expired"))
			return
		}

		// Refresh 용 토큰 생성
		newToken := jwt.New(jwt.GetSigningMethod(h.Guard.SigningMethod.Alg))
		newClaims := newToken.Claims.(jwt.MapClaims)

		for key := range claims {
			newClaims[key] = claims[key]
		}

		expire := time.Now().Add(h.Guard.Timeout)
		newClaims["sub"] = claims["sub"]
		newClaims["exp"] = expire.Unix()
		newClaims["iat"] = origIat

		// 현재는 HSxxx 형식 알고리즘만 지원
		tokenString, err := newToken.SignedString([]byte(h.Guard.SigningMethod.Key))
		if nil != err {
			response.Errorf(rw, req, http.StatusUnauthorized, errors.New("create JWT Token failed"))
			return
		}

		response.Write(rw, req, &AccessToken{
			Type:  "Bearer",
			Token: tokenString,
			//Expires: expire.Format(time.RFC3339),
			Expires: expire.Unix(),
		})
	}
}

// ===== [ Private Functions ] =====

// extractAccessToken - Request 정보에서 인증 정보 추출
func extractAccessToken(req *http.Request) (string, error) {
	// Access Key들을 검증하는 OAuth 활용
	authHeaderValue := req.Header.Get("Authorization")
	parts := strings.Split(authHeaderValue, " ")
	if 2 > len(parts) {
		return "", errors.New("attempted access with malformed header, no auth header found")
	}

	if strings.ToLower(parts[0]) != bearer {
		return "", errors.New("bearer token malformed")
	}

	return parts[1], nil
}

// getClient - 지정한 토큰을 처리하기 위한 HTTP Client 구성
func getClient(token string) *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return oauth2.NewClient(ctx, ts)
}

// ===== [ Public Functions ] =====
