// Package provider -
package provider

import (
	"net/http"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type (
	// Verifier - Provider에 대한 검증 인터페이스
	Verifier interface {
		Verify(req *http.Request, hc *http.Client) (bool, error)
	}

	// VerifierBasket - Verifier들 관리 구조
	VerifierBasket struct {
		verifiers []Verifier
	}
)

// ===== [ Implementations ] =====

// Verify - 지정한 Request들의 정보를 기준으로 Verify 수행
func (vb *VerifierBasket) Verify(req *http.Request, hc *http.Client) (bool, error) {
	var wrappedErrors error
	for _, verifier := range vb.verifiers {
		verified, err := verifier.Verify(req, hc)
		if nil != err {
			wrappedErrors = errors.Wrap(err, "verification failed")
			continue
		}
		if verified {
			return true, nil
		}
	}
	return false, wrappedErrors
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewVerifierBasket - 지정한 Verifier 들을 관리하는 Basket 생성
func NewVerifierBasket(verifiers ...Verifier) *VerifierBasket {
	return &VerifierBasket{verifiers: verifiers}
}
