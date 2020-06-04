package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
)

// ===== [ Constants and Variables ] =====

const (
	// HeaderCompleteResponseValue - Response가 정상적으로 처리된 경우에 Header에 사용할 값
	HeaderCompleteResponseValue = "true"
	// HeaderIncompleteResponseValue - Response가 비 정상적으로 처리된 경우에 Header에 사용할 값
	HeaderIncompleteResponseValue = "false"
)

var (
	// MessageResponseHeaderName - Response 오류인 경우 클라이언트에 표시할 Header 명
	MessageResponseHeaderName = "X-" + core.AppName + "-Messages"
	// CompleteResponseHeaderName - Response 정상 여부를 클라이언트에 표시할 Header 명
	CompleteResponseHeaderName = "X-" + core.AppName + "-Completed"
	// HeadersToSend - Router에 도착한 Request에서 Proxy로 전달할 필수 Header 정보들
	HeadersToSend = []string{"Content-Type"}
	// HeadersToNotSend - Router에 도착한 Request에서 Proxy로 전달하지 않을 Header 정보들
	HeadersToNotSend = []string{"Accept-Encoding"}
	// UserAgentHeaderValue - Proxy로 전달할 User-Agent Header 값
	UserAgentHeaderValue = []string{core.AppUserAgent}

	// ErrInternalError - 처리 중에 발생한 오류
	ErrInternalError = errors.New("internal server error")
	// ErrPrivateKey - TLS 설정에서 사용할 Private Key 미 설정 오류
	ErrPrivateKey = errors.New("private key not defined")
	// ErrPublicKey - TLS 설정에서 사용할 Public Key 미 설정 오류
	ErrPublicKey = errors.New("public key not defined")

	onceTransportConfig sync.Once

	defaultCurves = []tls.CurveID{
		tls.CurveP521,
		tls.CurveP384,
		tls.CurveP256,
	}
	defaultCipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	}
	versions = map[string]uint16{
		"SSL3.0": tls.VersionSSL30,
		"TLS10":  tls.VersionTLS10,
		"TLS11":  tls.VersionTLS11,
		"TLS12":  tls.VersionTLS12,
	}
)

// ===== [ Types ] =====

// ToHTTPError - 처리 중에 발생한 오류를 StatusCode로 처리하는 함수 정의
type ToHTTPError func(error) int

// ===== [ Private Functions ] =====

func parseTLSVersion(key string) uint16 {
	if v, ok := versions[key]; ok {
		return v
	}
	return tls.VersionTLS12
}

func parseCurveIDs(conf *config.TLSConfig) []tls.CurveID {
	l := len(conf.CurvePreferences)
	if l == 0 {
		return defaultCurves
	}

	curves := make([]tls.CurveID, len(conf.CurvePreferences))
	for i := range curves {
		curves[i] = tls.CurveID(conf.CurvePreferences[i])
	}
	return curves
}

func parseCipherSuites(conf *config.TLSConfig) []uint16 {
	l := len(conf.CipherSuites)
	if l == 0 {
		return defaultCipherSuites
	}

	cs := make([]uint16, l)
	for i := range cs {
		cs[i] = uint16(conf.CipherSuites[i])
	}
	return cs
}

// ===== [ Public Functions ] =====

// DefaultToHTTPError - 발생한 오류들을 전부 InternalServerError로 처리하는 기본 오류 처리기
func DefaultToHTTPError(_ error) int {
	return http.StatusInternalServerError
}

// InitHTTPDefaultTransport - 설정 기준으로 단 한번 설정되는 HTTP 설정 초기화
func InitHTTPDefaultTransport(sConf config.ServiceConfig) {
	onceTransportConfig.Do(func() {
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:       sConf.DialerTimeout,
				KeepAlive:     sConf.DialerKeepAlive,
				FallbackDelay: sConf.DialerFallbackDelay,
				DualStack:     true,
			}).DialContext,
			DisableCompression:    sConf.DisableCompression,
			DisableKeepAlives:     sConf.DisableKeepAlives,
			MaxIdleConns:          sConf.MaxIdleConnections,
			MaxIdleConnsPerHost:   sConf.MaxIdleConnectionsPerHost,
			IdleConnTimeout:       sConf.IdleConnectionTimeout,
			ResponseHeaderTimeout: sConf.ResponseHeaderTimeout,
			ExpectContinueTimeout: sConf.ExpectContinueTimeout,
			TLSHandshakeTimeout:   10 * time.Second,
		}
	})
}

// RunServer - 지정된 Context와 설정 및 Handler 기반으로 동작하는 HTT Server 구동
func RunServer(ctx context.Context, sConf config.ServiceConfig, handler http.Handler) error {
	done := make(chan error)
	s := NewServer(sConf, handler)

	if s.TLSConfig == nil {
		go func() {
			done <- s.ListenAndServe()
		}()
	} else {
		if sConf.TLS.PublicKey == "" {
			return ErrPublicKey
		}
		if sConf.TLS.PrivateKey == "" {
			return ErrPrivateKey
		}
		go func() {
			done <- s.ListenAndServeTLS(sConf.TLS.PublicKey, sConf.TLS.PrivateKey)
		}()
	}

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	}
}

// NewServer - 지정한 설정과 http handler 기준으로 동작하는 http server 반환
func NewServer(sConf config.ServiceConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", sConf.Port),
		Handler:           handler,
		ReadTimeout:       sConf.ReadTimeout,
		WriteTimeout:      sConf.WriteTimeout,
		ReadHeaderTimeout: sConf.ReadHeaderTimeout,
		IdleTimeout:       sConf.IdleTimeout,
		TLSConfig:         ParseTLSConfig(sConf.TLS),
	}
}

// ParseTLSConfig - 서비스 설정에 지정된 TLS 설정을 기준으로 tls 모듈에 대한 설정 반환
func ParseTLSConfig(tlsConf *config.TLSConfig) *tls.Config {
	if tlsConf == nil {
		return nil
	}
	if tlsConf.IsDisabled {
		return nil
	}

	return &tls.Config{
		MinVersion:               parseTLSVersion(tlsConf.MinVersion),
		MaxVersion:               parseTLSVersion(tlsConf.MaxVersion),
		CurvePreferences:         parseCurveIDs(tlsConf),
		PreferServerCipherSuites: tlsConf.PreferServerCipherSuites,
		CipherSuites:             parseCipherSuites(tlsConf),
	}
}
