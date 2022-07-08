package agent

import (
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent"

	agentcommon "github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent/common"
)

// InstallTelegraf 에이전트 설치
// @Summary Install Agent
// @Description 모니터링 에이전트 설치
// @Tags [Agent] Monitoring Agent
// @Accept  json
// @Produce  json
// @Param agentInfo body rest.AgentType true "Details for an Agent Install object"
// @Success 200 {object} rest.SimpleMsg
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /agent [post]
func InstallTelegraf(c echo.Context) error {
	params := &rest.AgentType{}
	if err := c.Bind(params); err != nil {
		return err
	}

	if !checkEmptyFormParam(params.ServiceType) {
		return c.JSON(http.StatusBadRequest, rest.SetMessage("empty agent type parameter"))
	}

	if util.CheckMCK8SType(params.ServiceType) {
		// 토큰 값이 비어있을 경우
		if !checkEmptyFormParam(params.ClientToken) {
			// 키 기반 연동일 때 데이터 확인
			if !checkEmptyFormParam(params.NsId, params.Mck8sId, params.APIServerURL, params.ServerCA, params.ClientCA) {
				return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter for mck8s agent installation by key"))
			} else {
				// 토큰 기반 연동일 때 데이터 확인
				if !checkEmptyFormParam(params.NsId, params.Mck8sId, params.APIServerURL) {
					return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter for mck8s agent installation by token"))
				}
			}
		}
	} else if util.CheckMCISType(params.ServiceType) {
		// MCIS 에이전트 form 파라미터 값 체크
		if !checkEmptyFormParam(params.NsId, params.McisId, params.VmId, params.PublicIp, params.UserName, params.SshKey, params.CspType) {
			return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter for mcis agent installation"))
		}
		if params.Port == "" {
			params.Port = "22"
		}
	} else {
		return c.JSON(http.StatusBadRequest, rest.SetMessage(fmt.Sprintf("unsupported agentType: %s", params.ServiceType)))
	}

	requestInfo := &agentcommon.AgentInstallInfo{
		NsId:         params.NsId,
		McisId:       params.McisId,
		VmId:         params.VmId,
		PublicIp:     params.PublicIp,
		UserName:     params.UserName,
		SshKey:       params.SshKey,
		CspType:      params.CspType,
		Port:         params.Port,
		ServiceType:  params.ServiceType,
		Mck8sId:      params.Mck8sId,
		APIServerURL: params.APIServerURL,
		ServerCA:     params.ServerCA,
		ClientCA:     params.ClientCA,
		ClientKey:    params.ClientKey,
		ClientToken:  params.ClientToken,
	}

	errCode, err := agent.InstallAgent(*requestInfo)
	if errCode != http.StatusOK {
		return c.JSON(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, rest.SetMessage("agent installation is finished"))
}

// TODO: WINDOW Version
func GetWindowInstaller(c echo.Context) error {
	rootPath := os.Getenv("CBMON_ROOT")
	filePath := rootPath + "/file/pkg/windows/installer/cbinstaller_windows_amd64.zip"
	return c.File(filePath)
}

// Telegraf config 파일 다운로드
func GetTelegrafConfFile(c echo.Context) error {
	// Query 파라미터 가져오기
	nsId := c.QueryParam("ns_id")
	mcisId := c.QueryParam("mcis_id")
	vmId := c.QueryParam("vm_id")
	cspType := c.QueryParam("csp_type")

	// Query 파라미터 값 체크
	if nsId == "" || mcisId == "" || vmId == "" || cspType == "" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("query parameter is missing"))
	}

	rootPath := os.Getenv("CBMON_ROOT")
	filePath := rootPath + "/file/conf/telegraf.conf"

	read, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// 파일 내의 변수 값 설정 (hostId, collectorServer)
	strConf := string(read)
	strConf = strings.ReplaceAll(strConf, "{{ns_id}}", nsId)
	strConf = strings.ReplaceAll(strConf, "{{mcis_id}}", mcisId)
	strConf = strings.ReplaceAll(strConf, "{{vm_id}}", vmId)
	strConf = strings.ReplaceAll(strConf, "{{csp_type}}", cspType)

	return c.Blob(http.StatusOK, "text/plain", []byte(strConf))
}

// Telegraf package 파일 다운로드
func GetTelegrafPkgFile(c echo.Context) error {
	// Query 파라미터 가져오기
	osType := c.QueryParam("osType")
	arch := c.QueryParam("arch")

	// Query 파라미터 값 체크
	if osType == "" || arch == "" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to get package. query parameter is missing"))
	}

	// osType, architecture 지원 여부 체크
	osType = strings.ToLower(osType)
	if osType != "ubuntu" && osType != "centos" && osType != "windows" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to get package. not supported OS type"))
	}
	if !strings.Contains(arch, "32") && !strings.Contains(arch, "64") {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to get package. not supported architecture"))
	}

	if strings.Contains(arch, "64") {
		arch = "x64"
	} else {
		arch = "x32"
	}

	// 제공 설치 파일 탐색
	rootPath := os.Getenv("CBMON_ROOT")
	filepath := rootPath + fmt.Sprintf("/file/pkg/%s/%s/", osType, arch)
	filename, err := agentcommon.GetPackageName(filepath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get package. osType %s not supported", osType)))
	}
	file := filepath + filename
	return c.File(file)
}

// UninstallAgent 에이전트 삭제
// @Summary Uninstall Agent
// @Description 모니터링 에이전트 제거
// @Tags [Agent] Monitoring Agent
// @Accept  json
// @Produce  json
// @Param agentInfo body rest.AgentType true "Details for an Agent Remove object"
// @Success 200 {object} rest.SimpleMsg
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /agent [delete]
func UninstallAgent(c echo.Context) error {
	params := &rest.AgentType{}
	if err := c.Bind(params); err != nil {
		return err
	}

	if !checkEmptyFormParam(params.ServiceType) {
		return c.JSON(http.StatusBadRequest, rest.SetMessage("empty agent type parameter"))
	}

	if util.CheckMCK8SType(params.ServiceType) {
		// 토큰 값이 비어있을 경우
		if !checkEmptyFormParam(params.ClientToken) {
			// 키 기반 연동일 때 데이터 확인
			if !checkEmptyFormParam(params.NsId, params.Mck8sId, params.APIServerURL, params.ServerCA, params.ClientCA, params.ClientKey) {
				return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter for mck8s agent uninstallation by key"))
			} else {
				// 토큰 기반 연동일 때 데이터 확인
				if !checkEmptyFormParam(params.NsId, params.Mck8sId, params.APIServerURL) {
					return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter for mck8s agent uninstallation by token"))
				}
			}
		}
	}

	// MCIS 에이전트 form 파라미터 값 체크
	if util.CheckMCISType(params.ServiceType) {
		if !checkEmptyFormParam(params.NsId, params.McisId, params.VmId, params.PublicIp, params.UserName, params.SshKey, params.CspType) {
			return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter for mcis agent uninstallation"))
		}
		if params.Port == "" {
			params.Port = "22"
		}
	}

	requestInfo := agentcommon.AgentInstallInfo{
		NsId:         params.NsId,
		McisId:       params.McisId,
		VmId:         params.VmId,
		PublicIp:     params.PublicIp,
		UserName:     params.UserName,
		SshKey:       params.SshKey,
		CspType:      params.CspType,
		Port:         params.Port,
		ServiceType:  params.ServiceType,
		Mck8sId:      params.Mck8sId,
		APIServerURL: params.APIServerURL,
		ServerCA:     params.ServerCA,
		ClientCA:     params.ClientCA,
		ClientKey:    params.ClientKey,
		ClientToken:  params.ClientToken,
	}
	errCode, err := agent.UninstallAgent(requestInfo)
	if errCode != http.StatusOK {
		return c.JSON(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, rest.SetMessage("Agent Uninstallation is finished"))
}

func checkEmptyFormParam(datas ...string) bool {
	for _, data := range datas {
		if len(strings.TrimSpace(data)) == 0 {
			return false
		}
	}
	return true
}
