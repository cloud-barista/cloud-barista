package agent

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent"
)

// InstallTelegraf 에이전트 설치
// @Summary Install agent to vm
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
	inputServiceType := strings.ToLower(params.ServiceType)
	// form 파라미터 값 체크
	if params.NsId == "" || params.McisId == "" || params.VmId == "" || params.PublicIp == "" || params.UserName == "" || params.SshKey == "" || params.CspType == "" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to get package. query parameter is missing"))
	}
	if inputServiceType == "" || inputServiceType == "default" {
		inputServiceType = "mcis"
	} else {
		inputServiceType = "mcks"
	}
	if params.Port == "" {
		params.Port = "22"
	}

	errCode, err := agent.InstallAgent(params.NsId, params.McisId, params.VmId, params.PublicIp, params.UserName, params.SshKey, params.CspType, params.Port, inputServiceType)
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
	filename, err := agent.GetPackageName(filepath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get package. osType %s not supported", osType)))
	}
	file := filepath + filename
	return c.File(file)
}

// UninstallAgent 에이전트 삭제
// @Summary Uninstall agent to vm
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

	// form 파라미터 값 체크
	if params.NsId == "" || params.McisId == "" || params.VmId == "" || params.PublicIp == "" || params.UserName == "" || params.SshKey == "" || params.CspType == "" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to get package. query parameter is missing"))
	}

	if params.Port == "" {
		params.Port = "22"
	}

	errCode, err := agent.UninstallAgent(params.NsId, params.McisId, params.VmId, params.PublicIp, params.UserName, params.SshKey, params.CspType, params.Port)
	if errCode != http.StatusOK {
		fmt.Println(errCode)
		return c.JSON(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, rest.SetMessage("Agent Uninstallation is finished"))
}
