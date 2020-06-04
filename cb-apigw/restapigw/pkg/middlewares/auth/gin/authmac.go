package gin

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// token - HMAC 처리를 위한 정보 구조
type token struct {
	secretKey string
	accessKey string
	timestamp string
	duration  string
}

// ===== [ Implementations ] =====

// generateToken - HMAC Token을 생성하고 운영을 위한 추가 정보를 포함한 Auth Token 생성
func (t token) generateToken() string {
	return hex.EncodeToString(append(t.makeToken(), []byte("||"+t.timestamp+"|"+t.duration+"|"+t.accessKey)...))
}

// makeToken - HMAC Token 생성
func (t token) makeToken() []byte {
	data := t.duration + "^" + t.timestamp + "^" + t.accessKey

	h := hmac.New(sha256.New, []byte(t.secretKey))
	h.Write([]byte(data))

	return h.Sum(nil)
}

// resolveToken - 전달된 원본 Auth Token을 HMAC Token과 추가 정보로 분리
func (t token) resolveToken(rawToken string, ids []string) error {
	tokenBytes, err := hex.DecodeString(rawToken)
	if err != nil {
		return err
	}

	data := bytes.Split(tokenBytes, []byte("||"))
	if len(data) != 2 || len(data[0]) == 0 || len(data[1]) == 0 {
		return errors.New("invalid auth token")
	}

	tokenInfo := strings.Split(string(data[1]), "|")
	if len(tokenInfo) != 3 {
		return errors.New("invalid auth token")
	}

	t.timestamp = tokenInfo[0]
	t.duration = tokenInfo[1]
	t.accessKey = tokenInfo[2]

	newToken := t.makeToken()
	if bytes.Compare(data[0], newToken) != 0 {
		return errors.New("invalid auth token")
	} else if !t.checkDuration() {
		return errors.New("invalid time limit or expired")
	}

	for _, v := range ids {
		if v == t.accessKey {
			return nil
		}
	}

	return errors.New("invalid access id")
}

// checkDuration - 현재 시간과 Token에 지정된 Duration 검증
func (t token) checkDuration() bool {
	td, err := parseDuration(t.duration)
	if err != nil {
		return false
	}

	ts, err := time.Parse(time.UnixDate, t.timestamp)
	if err != nil {
		return false
	}

	ts = ts.Add(td)

	return ts.Sub(getTime()) >= 0
}

// setTimestamp - Token 구성을 위한 기준 시간 (현재 시간) 설정
func (t token) setTimestamp() {
	t.timestamp = getTimestamp()
}

// ===== [ Private Functions ] =====

// getTime - 현재 시간 반환 (로컬)
func getTime() time.Time {
	return time.Now().UTC()
}

// getTimestamp - 현재 시간의 Timestamp 문자열 반환
func getTimestamp() string {
	return getTime().Format(time.UnixDate)
}

// parseDuration - 유효 액세스 기간을 검증하기 위한 time.Duration 정보 검증
func parseDuration(duration string) (time.Duration, error) {
	td, err := time.ParseDuration(duration)
	if err != nil {
		return 0, err
	}

	return td, nil
}

// ===== [ Public Functions ] =====

// ValidateToken - description
func ValidateToken(secretKey string, rawToken string, ids []string) error {
	if secretKey == "" {
		return errors.New("secretKey is required")
	}
	if rawToken == "" {
		return errors.New("rawToken that used to validate is required")
	}

	hmacToken := token{
		secretKey: secretKey,
	}
	if err := hmacToken.resolveToken(rawToken, ids); err != nil {
		return err
	}

	return nil
}

// CreateToken - HMAC 토큰 생성
func CreateToken(secretKey string, accessID string, duration string) (string, error) {
	if secretKey == "" {
		return "", errors.New("secretKey is required")
	}
	if accessID == "" {
		return "", errors.New("accessID is required")
	}
	if _, err := parseDuration(duration); err != nil {
		return "", err
	}

	hmacToken := token{
		secretKey: secretKey,
		accessKey: accessID,
		duration:  duration,
	}
	hmacToken.setTimestamp()

	return hmacToken.generateToken(), nil
}
