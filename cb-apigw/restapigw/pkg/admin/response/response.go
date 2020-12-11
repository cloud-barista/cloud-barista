// Package response - Admin Response 처리 기능 제공 패키지
package response

import (
	"net/http"
	"strconv"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/render"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====

type (
	// ReturnData - Response 데이터 구조
	ReturnData struct {
		Error   bool        `json:"error"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====

// getFields - Returns logging fields from request
func getFields(req *http.Request) logging.Fields {
	return logging.Fields{
		"host":       req.Host,
		"address":    req.RemoteAddr,
		"method":     req.Method,
		"requestURI": req.RequestURI,
		"proto":      req.Proto,
		"userAgent":  req.UserAgent(),
	}
}

// ===== [ Public Functions ] =====

// Errorf - 오류 발생 시 반환 처리
func Errorf(rw http.ResponseWriter, req *http.Request, code int, err error) {
	msg := err.Error()

	log := logging.GetLogger()
	log.SetFields(getFields(req)).WithError(err).Debug("[API SERVER] Processed Code: " + strconv.Itoa(code) + ", Message: " + msg)

	returnData := ReturnData{
		Error:   true,
		Code:    code,
		Message: msg,
		Data:    nil,
	}

	render.JSON(rw, code, &returnData)
}

// Write - 정상 처리 시 반환 처리
func Write(rw http.ResponseWriter, req *http.Request, data interface{}) {
	returnData := ReturnData{
		Error:   false,
		Code:    0,
		Message: "",
		Data:    data,
	}

	render.JSON(rw, http.StatusOK, &returnData)
}
