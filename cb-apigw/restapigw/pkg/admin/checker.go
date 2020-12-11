// Package admin -
package admin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/admin/health"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/api"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	ginAdapter "github.com/cloud-barista/cb-apigw/restapigw/pkg/core/adapters/gin"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====
// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====

// doStatusRequest - Health Checking용 요청 처리
func doStatusRequest(def *config.EndpointConfig, closeBody bool, logger logging.Logger) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, def.HealthCheck.URL, nil)
	if nil != err {
		logger.WithError(err).Error("[ADMIN Server] Creating the request for the health check failed")
		return nil, err
	}

	// 요청 후에 연결 종료 여부 설정
	req.Header.Set("Connection", "close")

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		logger.WithError(err).Error("[ADMIN Server] Making the request for the health check failed")
		return resp, err
	}

	if closeBody {
		defer resp.Body.Close()
	}

	return resp, err
}

// check - 지정한 API Definition에 대한 검증
func check(def *config.EndpointConfig, logger logging.Logger) func() error {
	return func() error {
		resp, err := doStatusRequest(def, true, logger)
		if nil != err {
			return fmt.Errorf("%s health check endpoint %s is unreachable", def.Name, def.HealthCheck.URL)
		}

		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("%s is not available at the moment", def.Name)
		}

		if resp.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("%s is partially available at the moment", def.Name)
		}

		return nil
	}
}

// findValidAPIHealthChecks - 지정된 API Definition들 중에 HealthCheck 정보가 있는 것들만 추출
func findValidAPIHealthChecks(maps []*api.DefinitionMap) []*config.EndpointConfig {
	validDefs := make([]*config.EndpointConfig, 0)

	for _, dm := range maps {
		for _, def := range dm.Definitions {
			if def.Active && "" != def.HealthCheck.URL {
				validDefs = append(validDefs, def)
			}
		}
	}

	return validDefs
}

// ===== [ Public Functions ] =====

// NewOverviewHandler - 모든 검증 대상에 대한 상태 검증용 핸들러 구성
func NewOverviewHandler(conf *api.Configuration, logger logging.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		defs := findValidAPIHealthChecks(conf.DefinitionMaps)

		logger.WithField("len", len(defs)).Debug("[ADMIN Server] Loading health check endpoints")
		health.Reset()

		for _, def := range defs {
			logger.WithField("name", def.Name).Debug("[ADMIN Server] Registering health check")
			health.Register(health.Config{
				Name:      def.Name,
				Timeout:   time.Second * time.Duration(def.HealthCheck.Timeout),
				SkipOnErr: true,
				Check:     check(def, logger),
			})
		}

		health.HandlerFunc(rw, req)
	}
}

// NewStatusHandler - 단일 검증 대상에 대한 상태 검증용 핸들러 구성
func NewStatusHandler(conf *api.Configuration, logger logging.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		defs := findValidAPIHealthChecks(conf.DefinitionMaps)

		name := ginAdapter.URLParam(req, "name")
		for _, def := range defs {
			if name == def.Name {
				resp, err := doStatusRequest(def, false, logger)
				if nil != err {
					logger.WithField("name", name).WithError(err).Error("[ADMIN Server] Error requesting service health status")
					rw.WriteHeader(http.StatusInternalServerError)
					rw.Write([]byte(err.Error()))
					return
				}

				body, err := ioutil.ReadAll(resp.Body)
				if closeErr := resp.Body.Close(); nil != closeErr {
					logger.WithField("name", name).WithError(closeErr).Error("[ADMIN Server] Error closing health status body")
				}

				if nil != err {
					logger.WithField("name", name).WithError(err).Error("[ADMIN Server] Error reading health status body")
					rw.WriteHeader(http.StatusInternalServerError)
					rw.Write([]byte(err.Error()))
					return
				}

				rw.WriteHeader(resp.StatusCode)
				rw.Write(body)
				return
			}
		}

		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Definition name is not found"))
	}
}
