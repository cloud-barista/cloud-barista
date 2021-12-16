package request

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	gc "github.com/cloud-barista/cb-mcks/src/grpc-api/common"
	"github.com/cloud-barista/cb-mcks/src/grpc-api/config"
	"github.com/cloud-barista/cb-mcks/src/grpc-api/logger"
	pb "github.com/cloud-barista/cb-mcks/src/grpc-api/protobuf/cbmcks"
	"github.com/cloud-barista/cb-mcks/src/grpc-api/request/mcar"

	"google.golang.org/grpc"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// MCARApi - MCKS API 구조 정의
type MCARApi struct {
	gConf        *config.GrpcConfig
	conn         *grpc.ClientConn
	jaegerCloser io.Closer
	clientMCAR   pb.MCARClient
	requestMCAR  *mcar.MCARRequest
	inType       string
	outType      string
}

// ClusterCreateRequest - CLUSTER 생성 요청 구조 Wrapper 정의
type ClusterCreateRequest struct {
	Namespace string     `yaml:"namespace" json:"namespace"`
	Item      ClusterReq `yaml:"ReqInfo" json:"ReqInfo"`
}

// ClusterReq - CLUSTER 생성 요청 구조 정의
type ClusterReq struct {
	Name         string       `yaml:"name" json:"name"`
	ControlPlane []NodeConfig `yaml:"controlPlane" json:"controlPlane"`
	Worker       []NodeConfig `yaml:"worker" json:"worker"`
	Config       Config       `yaml:"config" json:"config"`
}

// NodeConfig - Node 환경설정 구조 정의
type NodeConfig struct {
	Connection string `yaml:"connection" json:"connection"`
	Count      int    `yaml:"count" json:"count"`
	Spec       string `yaml:"spec" json:"spec"`
}

// Config - 클러스터 환경설정 구조 정의
type Config struct {
	Kubernetes Kubernetes `yaml:"kubernetes" json:"kubernetes"`
}

// Kubernetes - 쿠버네티스 환경설정 구조 정의
type Kubernetes struct {
	NetworkCni       string `yaml:"networkCni" json:"networkCni"`
	PodCidr          string `yaml:"podCidr" json:"podCidr"`
	ServiceCidr      string `yaml:"serviceCidr" json:"serviceCidr"`
	ServiceDnsDomain string `yaml:"serviceDnsDomain" json:"serviceDnsDomain"`
}

// NodeCreateRequest - NODE 생성 요청 구조 Wrapper 정의
type NodeCreateRequest struct {
	Namespace string  `yaml:"namespace" json:"namespace"`
	Cluster   string  `yaml:"cluster" json:"cluster"`
	Item      NodeReq `yaml:"ReqInfo" json:"ReqInfo"`
}

// NodeReq - NODE 생성 요청 구조 정의
type NodeReq struct {
	ControlPlane []NodeConfig `yaml:"controlPlane" json:"controlPlane"`
	Worker       []NodeConfig `yaml:"worker" json:"worker"`
}

// ===== [ Implementations ] =====

// SetServerAddr - MCKS 서버 주소 설정
func (m *MCARApi) SetServerAddr(addr string) error {
	if addr == "" {
		return errors.New("parameter is empty")
	}

	m.gConf.GSL.MCKSCli.ServerAddr = addr
	return nil
}

// GetServerAddr - MCKS 서버 주소 값 조회
func (m *MCARApi) GetServerAddr() (string, error) {
	return m.gConf.GSL.MCKSCli.ServerAddr, nil
}

// SetTLSCA - TLS CA 설정
func (m *MCARApi) SetTLSCA(tlsCAFile string) error {
	if tlsCAFile == "" {
		return errors.New("parameter is empty")
	}

	if m.gConf.GSL.MCKSCli.TLS == nil {
		m.gConf.GSL.MCKSCli.TLS = &config.TLSConfig{}
	}

	m.gConf.GSL.MCKSCli.TLS.TLSCA = tlsCAFile
	return nil
}

// GetTLSCA - TLS CA 값 조회
func (m *MCARApi) GetTLSCA() (string, error) {
	if m.gConf.GSL.MCKSCli.TLS == nil {
		return "", nil
	}

	return m.gConf.GSL.MCKSCli.TLS.TLSCA, nil
}

// SetTimeout - Timeout 설정
func (m *MCARApi) SetTimeout(timeout time.Duration) error {
	m.gConf.GSL.MCKSCli.Timeout = timeout
	return nil
}

// GetTimeout - Timeout 값 조회
func (m *MCARApi) GetTimeout() (time.Duration, error) {
	return m.gConf.GSL.MCKSCli.Timeout, nil
}

// SetJWTToken - JWT 인증 토큰 설정
func (m *MCARApi) SetJWTToken(token string) error {
	if token == "" {
		return errors.New("parameter is empty")
	}

	if m.gConf.GSL.MCKSCli.Interceptors == nil {
		m.gConf.GSL.MCKSCli.Interceptors = &config.InterceptorsConfig{}
		m.gConf.GSL.MCKSCli.Interceptors.AuthJWT = &config.AuthJWTConfig{}
	}
	if m.gConf.GSL.MCKSCli.Interceptors.AuthJWT == nil {
		m.gConf.GSL.MCKSCli.Interceptors.AuthJWT = &config.AuthJWTConfig{}
	}

	m.gConf.GSL.MCKSCli.Interceptors.AuthJWT.JWTToken = token
	return nil
}

// GetJWTToken - JWT 인증 토큰 값 조회
func (m *MCARApi) GetJWTToken() (string, error) {
	if m.gConf.GSL.MCKSCli.Interceptors == nil {
		return "", nil
	}
	if m.gConf.GSL.MCKSCli.Interceptors.AuthJWT == nil {
		return "", nil
	}

	return m.gConf.GSL.MCKSCli.Interceptors.AuthJWT.JWTToken, nil
}

// SetConfigPath - 환경설정 파일 설정
func (m *MCARApi) SetConfigPath(configFile string) error {
	logger := logger.NewLogger()

	// Viper 를 사용하는 설정 파서 생성
	parser := config.MakeParser()

	var (
		gConf config.GrpcConfig
		err   error
	)

	if configFile == "" {
		logger.Error("Please, provide the path to your configuration file")
		return errors.New("configuration file are not specified")
	}

	logger.Debug("Parsing configuration file: ", configFile)
	if gConf, err = parser.GrpcParse(configFile); err != nil {
		logger.Error("ERROR - Parsing the configuration file.\n", err.Error())
		return err
	}

	// MCKS CLIENT 필수 입력 항목 체크
	mckscli := gConf.GSL.MCKSCli

	if mckscli == nil {
		return errors.New("mckscli field are not specified")
	}

	if mckscli.ServerAddr == "" {
		return errors.New("mckscli.server_addr field are not specified")
	}

	if mckscli.Timeout == 0 {
		mckscli.Timeout = 90 * time.Second
	}

	if mckscli.TLS != nil {
		if mckscli.TLS.TLSCA == "" {
			return errors.New("mckscli.tls.tls_ca field are not specified")
		}
	}

	if mckscli.Interceptors != nil {
		if mckscli.Interceptors.AuthJWT != nil {
			if mckscli.Interceptors.AuthJWT.JWTToken == "" {
				return errors.New("mckscli.interceptors.auth_jwt.jwt_token field are not specified")
			}
		}
		if mckscli.Interceptors.Opentracing != nil {
			if mckscli.Interceptors.Opentracing.Jaeger != nil {
				if mckscli.Interceptors.Opentracing.Jaeger.Endpoint == "" {
					return errors.New("mckscli.interceptors.opentracing.jaeger.endpoint field are not specified")
				}
			}
		}
	}

	m.gConf = &gConf
	return nil
}

// Open - 연결 설정
func (m *MCARApi) Open() error {

	mckscli := m.gConf.GSL.MCKSCli

	// grpc 커넥션 생성
	cbconn, closer, err := gc.NewCBConnection(mckscli)
	if err != nil {
		return err
	}

	if closer != nil {
		m.jaegerCloser = closer
	}

	m.conn = cbconn.Conn

	// grpc 클라이언트 생성
	m.clientMCAR = pb.NewMCARClient(m.conn)

	// grpc 호출 Wrapper
	m.requestMCAR = &mcar.MCARRequest{Client: m.clientMCAR, Timeout: mckscli.Timeout, InType: m.inType, OutType: m.outType}

	return nil
}

// Close - 연결 종료
func (m *MCARApi) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
	if m.jaegerCloser != nil {
		m.jaegerCloser.Close()
	}

	m.jaegerCloser = nil
	m.conn = nil
	m.clientMCAR = nil
	m.requestMCAR = nil
}

// SetInType - 입력 문서 타입 설정 (json/yaml)
func (m *MCARApi) SetInType(in string) error {
	if in == "json" {
		m.inType = in
	} else if in == "yaml" {
		m.inType = in
	} else {
		return errors.New("input type is not supported")
	}

	if m.requestMCAR != nil {
		m.requestMCAR.InType = m.inType
	}

	return nil
}

// GetInType - 입력 문서 타입 값 조회
func (m *MCARApi) GetInType() (string, error) {
	return m.inType, nil
}

// SetOutType - 출력 문서 타입 설정 (json/yaml)
func (m *MCARApi) SetOutType(out string) error {
	if out == "json" {
		m.outType = out
	} else if out == "yaml" {
		m.outType = out
	} else {
		return errors.New("output type is not supported")
	}

	if m.requestMCAR != nil {
		m.requestMCAR.OutType = m.outType
	}

	return nil
}

// GetOutType - 출력 문서 타입 값 조회
func (m *MCARApi) GetOutType() (string, error) {
	return m.outType, nil
}

// Healthy - 상태확인
func (m *MCARApi) Healthy() (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	return m.requestMCAR.Healthy()
}

// CreateCluster - Cluster 생성
func (m *MCARApi) CreateCluster(doc string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	m.requestMCAR.InData = doc
	return m.requestMCAR.CreateCluster()
}

// CreateClusterByParam - Cluster 생성
func (m *MCARApi) CreateClusterByParam(req *ClusterCreateRequest) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	holdType, _ := m.GetInType()
	m.SetInType("json")
	j, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	m.requestMCAR.InData = string(j)
	result, err := m.requestMCAR.CreateCluster()
	m.SetInType(holdType)

	return result, err
}

// ListCluster - Cluster 목록
func (m *MCARApi) ListCluster(doc string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	m.requestMCAR.InData = doc
	return m.requestMCAR.ListCluster()
}

// ListClusterByParam - Cluster 목록
func (m *MCARApi) ListClusterByParam(namespace string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	holdType, _ := m.GetInType()
	m.SetInType("json")
	m.requestMCAR.InData = `{"namespace":"` + namespace + `"}`
	result, err := m.requestMCAR.ListCluster()
	m.SetInType(holdType)

	return result, err
}

// GetCluster - Cluster 조회
func (m *MCARApi) GetCluster(doc string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	m.requestMCAR.InData = doc
	return m.requestMCAR.GetCluster()
}

// GetClusterByParam - Cluster 조회
func (m *MCARApi) GetClusterByParam(namespace string, cluster string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	holdType, _ := m.GetInType()
	m.SetInType("json")
	m.requestMCAR.InData = `{"namespace":"` + namespace + `", "cluster":"` + cluster + `"}`
	result, err := m.requestMCAR.GetCluster()
	m.SetInType(holdType)

	return result, err
}

// DeleteCluster - Cluster 삭제
func (m *MCARApi) DeleteCluster(doc string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	m.requestMCAR.InData = doc
	return m.requestMCAR.DeleteCluster()
}

// DeleteClusterByParam - Cluster 삭제
func (m *MCARApi) DeleteClusterByParam(namespace string, cluster string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	holdType, _ := m.GetInType()
	m.SetInType("json")
	m.requestMCAR.InData = `{"namespace":"` + namespace + `", "cluster":"` + cluster + `"}`
	result, err := m.requestMCAR.DeleteCluster()
	m.SetInType(holdType)

	return result, err
}

// AddNode - Node 추가
func (m *MCARApi) AddNode(doc string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	m.requestMCAR.InData = doc
	return m.requestMCAR.AddNode()
}

// AddNodeByParam - Node 추가
func (m *MCARApi) AddNodeByParam(req *NodeCreateRequest) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	holdType, _ := m.GetInType()
	m.SetInType("json")
	j, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	m.requestMCAR.InData = string(j)
	result, err := m.requestMCAR.AddNode()
	m.SetInType(holdType)

	return result, err
}

// ListNode - Node 목록
func (m *MCARApi) ListNode(doc string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	m.requestMCAR.InData = doc
	return m.requestMCAR.ListNode()
}

// ListNodeByParam - Node 목록
func (m *MCARApi) ListNodeByParam(namespace string, cluster string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	holdType, _ := m.GetInType()
	m.SetInType("json")
	m.requestMCAR.InData = `{"namespace":"` + namespace + `", "cluster":"` + cluster + `"}`
	result, err := m.requestMCAR.ListNode()
	m.SetInType(holdType)

	return result, err
}

// GetNode - Node 조회
func (m *MCARApi) GetNode(doc string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	m.requestMCAR.InData = doc
	return m.requestMCAR.GetNode()
}

// GetNodeByParam - Node 조회
func (m *MCARApi) GetNodeByParam(namespace string, cluster string, node string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	holdType, _ := m.GetInType()
	m.SetInType("json")
	m.requestMCAR.InData = `{"namespace":"` + namespace + `", "cluster":"` + cluster + `", "node":"` + node + `"}`
	result, err := m.requestMCAR.GetNode()
	m.SetInType(holdType)

	return result, err
}

// RemoveNode - Node 삭제
func (m *MCARApi) RemoveNode(doc string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	m.requestMCAR.InData = doc
	return m.requestMCAR.RemoveNode()
}

// RemoveNodeByParam - Node 삭제
func (m *MCARApi) RemoveNodeByParam(namespace string, cluster string, node string) (string, error) {
	if m.requestMCAR == nil {
		return "", errors.New("The Open() function must be called")
	}

	holdType, _ := m.GetInType()
	m.SetInType("json")
	m.requestMCAR.InData = `{"namespace":"` + namespace + `", "cluster":"` + cluster + `", "node":"` + node + `"}`
	result, err := m.requestMCAR.RemoveNode()
	m.SetInType(holdType)

	return result, err
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewMCARManager - MCAR API 객체 생성
func NewMCARManager() (m *MCARApi) {

	m = &MCARApi{}
	m.gConf = &config.GrpcConfig{}
	m.gConf.GSL.MCKSCli = &config.GrpcClientConfig{}

	m.jaegerCloser = nil
	m.conn = nil
	m.clientMCAR = nil
	m.requestMCAR = nil

	m.inType = "json"
	m.outType = "json"

	return
}
