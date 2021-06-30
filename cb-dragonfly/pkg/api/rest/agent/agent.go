package agent

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"

	"github.com/cloud-barista/cb-dragonfly/pkg/core/agent"
)

func InstallTelegraf(c echo.Context) error {
	// form 파라미터 값 가져오기
	nsId := c.FormValue("ns_id")
	mcisId := c.FormValue("mcis_id")
	vmId := c.FormValue("vm_id")
	publicIp := c.FormValue("public_ip")
	userName := c.FormValue("user_name")
	sshKey := c.FormValue("ssh_key")
	cspType := c.FormValue("cspType")
	port := c.FormValue("port")

	// form 파라미터 값 체크
	if nsId == "" || mcisId == "" || vmId == "" || publicIp == "" || userName == "" || sshKey == "" || cspType == "" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to get package. query parameter is missing"))
	}
	if port == "" {
		port = "22"
	}

	errCode, err := agent.InstallTelegraf(nsId, mcisId, vmId, publicIp, userName, sshKey, cspType, port)
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

func UninstallAgent(c echo.Context) error {
	nsId := c.FormValue("ns_id")
	mcisId := c.FormValue("mcis_id")
	vmId := c.FormValue("vm_id")
	publicIp := c.FormValue("public_ip")
	userName := c.FormValue("user_name")
	sshKey := c.FormValue("ssh_key")
	cspType := c.FormValue("cspType")
	port := c.FormValue("port")
	// form 파라미터 값 체크
	if nsId == "" || mcisId == "" || vmId == "" || publicIp == "" || userName == "" || sshKey == "" || cspType == "" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to get package. query parameter is missing"))
	}

	if port == "" {
		port = "22"
	}

	errCode, err := agent.UninstallAgent(nsId, mcisId, vmId, publicIp, userName, sshKey, cspType, port)
	if errCode != http.StatusOK {
		fmt.Println(errCode)
		return c.JSON(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, rest.SetMessage("Agent Uninstallation is finished"))
}
