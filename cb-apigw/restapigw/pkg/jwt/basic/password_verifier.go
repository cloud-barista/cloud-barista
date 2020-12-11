// Package basic -
package basic

import (
	"net/http"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
)

// ===== [ Constants and Variables ] =====
const (
	contentTypeJSON string = "application/json"
)

var ()

// ===== [ Types ] =====
type (
	// user - 로그인 사용자 정보 구조
	user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// PasswordVerifier - 패스워드 검증을 위한 사용자들 구조
	PasswordVerifier struct {
		users []*user
	}
)

// ===== [ Implementations ] =====

// Equals - 지정한 사용자와 동일한지 여부 검증
func (u *user) Equals(user *user) bool {
	return user.Username == u.Username && user.Password == u.Password
}

// getUserFromRequest - 지정한 Request에서 사용자 정보 추출
func (pv *PasswordVerifier) getUserFromRequest(req *http.Request) (*user, error) {
	var u *user

	// Request에서 Basic Auth 정보 추출
	username, password, ok := req.BasicAuth()
	u = &user{
		Username: username,
		Password: password,
	}

	// Basic Auth 정보가 아니라면 Request에서 포맷정보 추출해서 검증
	if !ok {
		contentType := filterFlags(req.Header.Get("Content-Type"))
		switch contentType {
		case contentTypeJSON:
			err := core.JSONDecode(req.Body, &u)
			if nil != err {
				return u, errors.Wrap(err, "could not parse the json body")
			}
		default:
			req.ParseForm()

			u = &user{
				Username: req.Form.Get("username"),
				Password: req.Form.Get("password"),
			}
		}
	}

	return u, nil
}

// Verify - 지정한 Request 정보를 기준으로 사용자 정보를 검증
func (pv *PasswordVerifier) Verify(req *http.Request, hc *http.Client) (bool, error) {
	currUser, err := pv.getUserFromRequest(req)
	if nil != err {
		return false, errors.Wrap(err, "could not get user from request")
	}

	for _, user := range pv.users {
		if user.Equals((currUser)) {
			return true, nil
		}
	}

	return false, errors.New("incorret username or password")
}

// ===== [ Private Functions ] =====

// 지정한 Content에서 필요 정보 추출
func filterFlags(content string) string {
	for i, char := range content {
		if ' ' == char || ';' == char {
			return content[:i]
		}
	}
	return content
}

// ===== [ Public Functions ] =====

// NewPasswordVerifier - 사용자의 비밀번호 검증을 위한 인스턴스 생성
func NewPasswordVerifier(users []*user) *PasswordVerifier {
	return &PasswordVerifier{users}
}
