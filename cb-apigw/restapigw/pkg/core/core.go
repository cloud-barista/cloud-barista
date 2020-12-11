// Package core - Defines variables/constants and provides utilty functions
package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	jsoniter "github.com/json-iterator/go"
)

// ===== [ Constants and Variables ] =====

const (
	// AppName - 어플리케이션 명
	AppName = "cb-restapigw"
	// AppVersion - 어플리케이션 버전
	AppVersion = "0.1.0"
	// AppHeaderName - 어플리케이션 식별을 위한 Header 관리 명
	AppHeaderName = "X-CB-RESTAPIGW"
	// AppUserAgent - Backend 전달에 사용할 User Agent Header 값
	AppUserAgent = AppName + " version " + AppVersion
	// CollectionTag - Backend의 Array를 Json 객체의 데이터로 반환 처리를 위한 Tag Name
	CollectionTag = "collection"
	// WrappingTag - Backend의 Array 직접 반환 처리를 위한 Tag Name
	WrappingTag = "!!wrapping!!"
	// Bypass - Endpoint/Backend Bypass 처리용 식별자
	Bypass = "*bypass"

	// RequestIDKey - Request ID 추적을 위한 Key
	RequestIDKey requestIDKeyType = iota
)

// ===== [ Types ] =====

type (
	// requestIDType - RequestID 식별 형식
	requestIDKeyType int

	// WrappedError - 원본 오류를 관리하는 오류 형식
	WrappedError struct {
		code          int
		message       string
		originalError error
	}
)

// ===== [ Implementations ] =====

// Code - Wrapping된 오류 코드 반환
func (we WrappedError) Code() int {
	return we.code
}

// Error - 오류 메시지 반환
func (we WrappedError) Error() string {
	return fmt.Sprintf("%d, %s", we.code, we.message)
}

// GetError - 원본 오류 반환
func (we WrappedError) GetError() error {
	return we.originalError
}

// ===== [ Private Functions ] =====

// getClientIPByRequestRemoteAddr - Request의 Remote Addr를 통한 IP 검증
func getClientIPByRequestRemoteAddr(req *http.Request) (string, error) {
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if nil != err {
		log.Printf("debug: Getting req.RemoteAddr: %v\n", err)
		return "", err
	}
	log.Printf("debug: With req.RemoteAddr found IP: %v, Port: %v\n", ip, port)

	userIP := net.ParseIP(ip)
	if nil == userIP {
		message := fmt.Sprintf("debug: Parsing IP from Request.RemoteAddr got nothing.")
		log.Println(message)
		return "", fmt.Errorf(message)
	}

	log.Printf("debug: Found IP: %v\n", userIP)
	return userIP.String(), nil
}

// getClientIPByHeaders - Request Header를 통한 IP 검증
func getClientIPByHeaders(req *http.Request) (string, error) {
	ipSlice := []string{}
	ipSlice = append(ipSlice, req.Header.Get("X-Forwarded-For"))
	ipSlice = append(ipSlice, req.Header.Get("x-forwarded-for"))
	ipSlice = append(ipSlice, req.Header.Get("X-FORWARDED-FOR"))

	for _, v := range ipSlice {
		log.Printf("debug: client request header check gives ip: %v\n", v)
		if "" != v {
			return v, nil
		}
	}

	err := errors.New("error: Could not find clients IP address from the Request Headers")
	return "", err
}

// getMyInterfaceAddr - Private network IP를 통한 IP 검증
func getMyInterfaceAddr() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if nil != err {
		return nil, err
	}
	addresses := []net.IP{}
	for _, iface := range ifaces {
		if 0 == iface.Flags&net.FlagUp {
			continue // interface down
		}
		if 0 != iface.Flags&net.FlagLoopback {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if nil != err {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if nil == ip || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if nil == ip {
				continue // not an ipv4 address
			}
			addresses = append(addresses, ip)
		}
	}

	if 0 == len(addresses) {
		return nil, fmt.Errorf("no address found, net.InterfaceAddrs: %v", addresses)
	}

	// only need first
	return addresses[0], nil
}

// ===== [ Public Functions ] =====

// ContextWithSignal - OS Interrupt signal 연계 처리를 위한 Context 구성
func ContextWithSignal(ctx context.Context) context.Context {
	newCtx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signals:
			cancel()
			close(signals)
		}
	}()

	return newCtx
}

// NewWrappedError - 원본 오류를 관리하는 오류 생성
func NewWrappedError(code int, message string, originalError error) error {
	return WrappedError{
		code:          code,
		message:       message,
		originalError: originalError,
	}
}

// GetStrings - 지정된 맵 데이터에서 지정된 이름에 해당하는 데이터를 []string 으로 반환
func GetStrings(data map[string]interface{}, name string) []string {
	result := []string{}
	if datas, ok := data[name]; ok {
		if data, ok := datas.([]interface{}); ok {
			for _, val := range data {
				if strVal, ok := val.(string); ok {
					result = append(result, strVal)
				}
			}
		}
	}
	return result
}

// GetString - 지정된 맵 데이터에서 지정한 키에 해당하는 데이터를 string 으로 반환
func GetString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

// GetBool - 지정된 맵 데이터에서 지정한 키에 해당하는 데이터를 bool 으로 반환
func GetBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// GetInt64 - 지정된 맵 데이터에서 지정한 키에 해당하는 데이터를 int64 로 반환
func GetInt64(data map[string]interface{}, key string) int64 {
	if val, ok := data[key]; ok {
		switch i := val.(type) {
		case int64:
			return i
		case int:
			return int64(i)
		case float64:
			return int64(i)
		}
	}
	return 0
}

// ContainsString returns true if a string is present in a iteratee.
func ContainsString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// GetResponseString - http.Response Body를 문자열로 반환
func GetResponseString(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return "", err
	}

	defer func() {
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}()

	return string(body), nil
}

// GetClientIPHelper - Request 기반의 Client IP를 검증
func GetClientIPHelper(req *http.Request) (string, error) {
	// Try parse "Origin" from header
	url, err := url.Parse(req.Header.Get("Origin"))
	if nil == err {
		host := url.Host
		ip, _, err := net.SplitHostPort(host)
		if nil == err {
			log.Printf("debug: Found IP using Header (Origin) sniffing, ip: %v\n", ip)
			return ip, nil
		}
	}

	// Try parse request
	ip, err := getClientIPByRequestRemoteAddr(req)
	if nil == err {
		log.Printf("debug: Found IP using Request, ip: %v\n", ip)
		return ip, nil
	}

	// Try parse "X-Forwarder" from header
	ip, err = getClientIPByHeaders(req)
	if nil == err {
		log.Printf("debug: Found IP using Request Headers (X-Forwarder) sniffing, ip: %v\n", ip)
		return ip, nil
	}

	err = errors.New("error: Could not find clients IP address")
	return "", err
}

// GetLastPart - 지정한 문자열을 지정한 문자로 분리하고 마지막 부분 반환
func GetLastPart(source, seperater string) string {
	if "" == source {
		return source
	}

	srcs := strings.Split(source, seperater)
	if 1 == len(srcs) {
		return srcs[0]
	}
	return srcs[len(srcs)-1]
}

// ToJSON - 지정 정보를 JSON 문자열로 변환
func ToJSON(data interface{}) string {
	bytes, err := JSONMarshal(data)
	if nil != err {
		log.Println("error on convert to json")
		return ""
	}
	return string(bytes)
}

// FromJSON - 지정한 JSON 문자열을 지정한 struct로 변환
func FromJSON(data string, target interface{}) error {
	return JSONUnmarshal([]byte(data), target)
}

// JSONDecode - 지정한 Source의 JSON 정보를 지정한 Target으로 설정
func JSONDecode(source io.Reader, target interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	dec := json.NewDecoder(source)
	dec.UseNumber()
	return dec.Decode(target)
}

// JSONMarshal - 지정한 정보를 JSON 으로 Marshal 처리
func JSONMarshal(data interface{}) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(data)
}

// JSONUnmarshal - 지정한 JSON 문자열 Byte를 지정한 Target으로 설정
func JSONUnmarshal(source []byte, target interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(source, target)
}

// GCD - 지정한 두개의 숫자에 대한 최대 공약수 계산
func GCD(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// RemoveSlice - 지정한 배열 구조에서 지정한 인덱스의 값을 삭제하고 반환
func RemoveSlice(arr []interface{}, idx int) []interface{} {
	arr[idx] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}
