// Cloud Info Manager's Rest Runtime of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2020.06.

package adminweb

import (
	"fmt"
	cr "github.com/cloud-barista/cb-spider/api-runtime/common-runtime"

	"strconv"
	"net/http"
	"strings"
	"github.com/labstack/echo/v4"
	"encoding/json"
)

// number, Spider's NameId, CSP's SystemId, checkbox
func makeMgmtTRList_html(bgcolor string, height string, fontSize string, infoList cr.AllResourceList) string {
        if bgcolor == "" { bgcolor = "#FFFFFF" }
        if height == "" { height = "30" }
        if fontSize == "" { fontSize = "2" }

        // make base TR frame for info list
        strTR := fmt.Sprintf(`
                <tr bgcolor="%s" align="center" height="%s">
                    <td>
                            <font size=%s>$$NUM$$</font>
                    </td>
                    <td $$NAMEIDSTYLE$$>
                            <font size=%s>$$NAMEID$$</font>
                    </td>
                    <td $$SYTEMIDSTYLE$$>
                            <font size=%s>$$SYTEMID$$</font>
                    </td>
                    <td>
                        <input type="checkbox" name="check_box" value=$$ID$$>
                    </td>
                </tr>
                `, bgcolor, height, fontSize, fontSize, fontSize)

        strData := ""
        // set data and make TR list
        for i, one := range infoList.AllList.MappedList{
                str := strings.ReplaceAll(strTR, "$$NUM$$", strconv.Itoa(i+1))
                str = strings.ReplaceAll(str, "$$NAMEID$$", one.NameId)
                str = strings.ReplaceAll(str, "$$SYTEMID$$", one.SystemId)
                str = strings.ReplaceAll(str, "$$ID$$", "::NAMEID::" + one.NameId) // MappedList: contain "::NAMEID::"
                str = strings.ReplaceAll(str, "$$NAMEIDSTYLE$$", `style="background-color:#F0F3FF;"`)
                str = strings.ReplaceAll(str, "$$SYTEMIDSTYLE$$", `style="background-color:#F0F3FF;"`)
                strData += str
        }
        for i, one := range infoList.AllList.OnlySpiderList{
                str := strings.ReplaceAll(strTR, "$$NUM$$", strconv.Itoa(i+1))
                str = strings.ReplaceAll(str, "$$NAMEID$$", one.NameId)
                str = strings.ReplaceAll(str, "$$SYTEMID$$", "( " + one.SystemId + " )")
                str = strings.ReplaceAll(str, "$$ID$$", "::NAMEID::" + one.NameId) // OnlySpiderList: contain "::NAMEID::"

                str = strings.ReplaceAll(str, "$$NAMEIDSTYLE$$", `style="background-color:#F0F3FF;"`)
                str = strings.ReplaceAll(str, "$$SYTEMIDSTYLE$$", ``)
                strData += str
        }
        for i, one := range infoList.AllList.OnlyCSPList{
                str := strings.ReplaceAll(strTR, "$$NUM$$", strconv.Itoa(i+1))
                str = strings.ReplaceAll(str, "$$NAMEID$$", "( " + one.NameId + " )")
                str = strings.ReplaceAll(str, "$$SYTEMID$$", one.SystemId)
                str = strings.ReplaceAll(str, "$$ID$$", one.SystemId) // OnlyCSPList: not contain "::NAMEID::"

                str = strings.ReplaceAll(str, "$$NAMEIDSTYLE$$", ``)
                str = strings.ReplaceAll(str, "$$SYTEMIDSTYLE$$", `style="background-color:#F0F3FF;"`)
                strData += str
        }        

        return strData
}

//====================================== VPC

// make the string of javascript function
func makeDeleteVPCMgmtFunc_js() string {
// delete for MappedList & OnlySpiderList
// curl -sX DELETE http://localhost:1024/spider/vpc/vpc-01?force=true -H 'Content-Type: application/json' -d '{ "ConnectionName": "aws-ohio-config"}
// delete for OnlyCSPList
// curl -sX DELETE http://localhost:1024/spider/cspvpc/vpc-0b0d0d30794eab379 -H 'Content-Type: application/json' -d '{ "ConnectionName": "aws-ohio-config"}' |json_pp

        strFunc := `
                function deleteVPCMgmt() {
                        var connConfig = parent.frames["top_frame"].document.getElementById("connConfig").innerHTML;
                        var checkboxes = document.getElementsByName('check_box');
                        for (var i = 0; i < checkboxes.length; i++) { // @todo make parallel executions
                                if (checkboxes[i].checked) {
                                        var xhr = new XMLHttpRequest();
                                        if(checkboxes[i].value.includes("::NAMEID::")) { // MappedList & OnlySpiderList
                                            xhr.open("DELETE", "$$SPIDER_SERVER$$/spider/vpc/" + checkboxes[i].value.replace("::NAMEID::", "") + "?force=true", false);
                                        }else { // OnlyCSPList
                                            xhr.open("DELETE", "$$SPIDER_SERVER$$/spider/cspvpc/" + checkboxes[i].value, false);
                                        }

                                        xhr.setRequestHeader('Content-Type', 'application/json');
					sendJson = '{ "ConnectionName": "' + connConfig + '"}'
                                        xhr.send(sendJson);
                                }
                        }
			location.reload();
                }
        `
        strFunc = strings.ReplaceAll(strFunc, "$$SPIDER_SERVER$$", "http://" + cr.HostIPorName + cr.ServicePort) // cr.ServicePort = ":1024"
        return strFunc
}

func VPCMgmt(c echo.Context) error {
        cblog.Info("call VPCMgmt()")

	connConfig := c.Param("ConnectConfig")
	if connConfig == "region not set" {
		htmlStr :=  `
			<html>
			<head>
			    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
			    <script type="text/javascript">
				alert(connConfig)
			    </script>
			</head>
			<body>
				<br>
				<br>
				<label style="font-size:24px;color:#606262;">&nbsp;&nbsp;&nbsp;Please select a Connection Configuration! (MENU: 2.CONNECTION)</label>	
			</body>
		`

		return c.HTML(http.StatusOK, htmlStr)
	}
	
        // make page header
        htmlStr :=  `
                <html>
                <head>
                    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                    <script type="text/javascript">
                `
        // (1) make Javascript Function
                htmlStr += makeCheckBoxToggleFunc_js()                
                htmlStr += makeDeleteVPCMgmtFunc_js()


        htmlStr += `
                    </script>
                </head>

                <body>
                    <table border="0" bordercolordark="#F8F8FF" cellpadding="0" cellspacing="1" bgcolor="#FFFFFF"  style="font-size:small;">
                `

        // (2) make Table Action TR
                // colspan, f5_href, delete_href, fontSize
                //htmlStr += makeActionTR_html("4", "vpc", "deleteVPCMgmt()", "2")
                htmlStr += makeActionTR_html("4", "", "deleteVPCMgmt()", "2")


        // (3) make Table Header TR
                nameWidthList := []NameWidth {
                    {"Spider's NameId", "300"},
                    {"CSP's SystemId", "300"},                    
                }
                htmlStr +=  makeTitleTRList_html("#DDDDDD", "2", nameWidthList, true)


        // (4) make TR list with info list
        // (4-1) get info list 
                resBody, err := getAllResourceList_with_Connection_JsonByte(connConfig, "vpc")
                if err != nil {
                        cblog.Error(err)
                        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
                }
                
                var info cr.AllResourceList
                
                json.Unmarshal(resBody, &info)

        // (4-2) make TR list with info list
                htmlStr += makeMgmtTRList_html("", "", "", info)

        // make page tail
        htmlStr += `
                    </table>
		    <hr>
                </body>
                </html>
        `

//fmt.Println(htmlStr)
        return c.HTML(http.StatusOK, htmlStr)
}

//====================================== Security Group

// make the string of javascript function
func makeDeleteSecurityGroupMgmtFunc_js() string {
// delete for MappedList & OnlySpiderList
// curl -sX DELETE http://localhost:1024/spider/securitygroup/sg-01?force=true -H 'Content-Type: application/json' -d '{ "ConnectionName": "aws-ohio-config"}
// delete for OnlyCSPList
// curl -sX DELETE http://localhost:1024/spider/cspsecuritygroup/sg-0b0d0d30794eab379 -H 'Content-Type: application/json' -d '{ "ConnectionName": "aws-ohio-config"}' |json_pp

        strFunc := `
                function deleteSecurityGroupMgmt() {
                        var connConfig = parent.frames["top_frame"].document.getElementById("connConfig").innerHTML;
                        var checkboxes = document.getElementsByName('check_box');
                        for (var i = 0; i < checkboxes.length; i++) { // @todo make parallel executions
                                if (checkboxes[i].checked) {
                                        var xhr = new XMLHttpRequest();
                                        if(checkboxes[i].value.includes("::NAMEID::")) { // if MappedList & OnlySpiderList
                                            xhr.open("DELETE", "$$SPIDER_SERVER$$/spider/securitygroup/" + checkboxes[i].value.replace("::NAMEID::", "") + "?force=true", false);
                                        }else { // OnlyCSPList
                                            xhr.open("DELETE", "$$SPIDER_SERVER$$/spider/cspsecuritygroup/" + checkboxes[i].value, false);
                                        }

                                        xhr.setRequestHeader('Content-Type', 'application/json');
                                        sendJson = '{ "ConnectionName": "' + connConfig + '"}'
                                        xhr.send(sendJson);
                                }
                        }
                        location.reload();
                }
        `
        strFunc = strings.ReplaceAll(strFunc, "$$SPIDER_SERVER$$", "http://" + cr.HostIPorName + cr.ServicePort) // cr.ServicePort = ":1024"
        return strFunc
}

func SecurityGroupMgmt(c echo.Context) error {
        cblog.Info("call SecurityGroupMgmt()")

        connConfig := c.Param("ConnectConfig")
        if connConfig == "region not set" {
                htmlStr :=  `
                        <html>
                        <head>
                            <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                            <script type="text/javascript">
                                alert(connConfig)
                            </script>
                        </head>
                        <body>
                                <br>
                                <br>
                                <label style="font-size:24px;color:#606262;">&nbsp;&nbsp;&nbsp;Please select a Connection Configuration! (MENU: 2.CONNECTION)</label>
                        </body>
                `

                return c.HTML(http.StatusOK, htmlStr)
        }

        // make page header
        htmlStr :=  `
                <html>
                <head>
                    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                    <script type="text/javascript">
                `
        // (1) make Javascript Function
                htmlStr += makeCheckBoxToggleFunc_js()
                htmlStr += makeDeleteSecurityGroupMgmtFunc_js()


        htmlStr += `
                    </script>
                </head>

                <body>
                    <table border="0" bordercolordark="#F8F8FF" cellpadding="0" cellspacing="1" bgcolor="#FFFFFF"  style="font-size:small;">
                `

        // (2) make Table Action TR
                // colspan, f5_href, delete_href, fontSize
                //htmlStr += makeActionTR_html("4", "vpc", "deleteSecurityGroupMgmt()", "2")
                htmlStr += makeActionTR_html("4", "", "deleteSecurityGroupMgmt()", "2")


        // (3) make Table Header TR
                nameWidthList := []NameWidth {
                    {"Spider's NameId", "300"},
                    {"CSP's SystemId", "300"},
                }
                htmlStr +=  makeTitleTRList_html("#DDDDDD", "2", nameWidthList, true)


        // (4) make TR list with info list
        // (4-1) get info list
                resBody, err := getAllResourceList_with_Connection_JsonByte(connConfig, "securitygroup")
                if err != nil {
                        cblog.Error(err)
                        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
                }

                var info cr.AllResourceList

                json.Unmarshal(resBody, &info)

        // (4-2) make TR list with info list
                htmlStr += makeMgmtTRList_html("", "", "", info)

        // make page tail
        htmlStr += `
                    </table>
                    <hr>
                </body>
                </html>
        `

//fmt.Println(htmlStr)
        return c.HTML(http.StatusOK, htmlStr)
}

//====================================== KeyPair

// make the string of javascript function
func makeDeleteKeyPairMgmtFunc_js() string {
// delete for MappedList & OnlySpiderList
// curl -sX DELETE http://localhost:1024/spider/keypair/keypair-01?force=true -H 'Content-Type: application/json' -d '{ "ConnectionName": "aws-ohio-config"}
// delete for OnlyCSPList
// curl -sX DELETE http://localhost:1024/spider/cspkeypair/0b0d0d30794eab379 -H 'Content-Type: application/json' -d '{ "ConnectionName": "aws-ohio-config"}' |json_pp

        strFunc := `
                function deleteKeyPairMgmt() {
                        var connConfig = parent.frames["top_frame"].document.getElementById("connConfig").innerHTML;
                        var checkboxes = document.getElementsByName('check_box');
                        for (var i = 0; i < checkboxes.length; i++) { // @todo make parallel executions
                                if (checkboxes[i].checked) {
                                        var xhr = new XMLHttpRequest();
                                        if(checkboxes[i].value.includes("::NAMEID::")) { // MappedList & OnlySpiderList
                                            xhr.open("DELETE", "$$SPIDER_SERVER$$/spider/keypair/" + checkboxes[i].value.replace("::NAMEID::", "") + "?force=true", false);
                                        }else { // OnlyCSPList
                                            xhr.open("DELETE", "$$SPIDER_SERVER$$/spider/cspkeypair/" + checkboxes[i].value, false);
                                        }

                                        xhr.setRequestHeader('Content-Type', 'application/json');
                    sendJson = '{ "ConnectionName": "' + connConfig + '"}'
                                        xhr.send(sendJson);
                                }
                        }
            location.reload();
                }
        `
        strFunc = strings.ReplaceAll(strFunc, "$$SPIDER_SERVER$$", "http://" + cr.HostIPorName + cr.ServicePort) // cr.ServicePort = ":1024"
        return strFunc
}

func KeyPairMgmt(c echo.Context) error {
        cblog.Info("call KeyPairMgmt()")

    connConfig := c.Param("ConnectConfig")
    if connConfig == "region not set" {
        htmlStr :=  `
            <html>
            <head>
                <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                <script type="text/javascript">
                alert(connConfig)
                </script>
            </head>
            <body>
                <br>
                <br>
                <label style="font-size:24px;color:#606262;">&nbsp;&nbsp;&nbsp;Please select a Connection Configuration! (MENU: 2.CONNECTION)</label>   
            </body>
        `

        return c.HTML(http.StatusOK, htmlStr)
    }
    
        // make page header
        htmlStr :=  `
                <html>
                <head>
                    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                    <script type="text/javascript">
                `
        // (1) make Javascript Function
                htmlStr += makeCheckBoxToggleFunc_js()                
                htmlStr += makeDeleteKeyPairMgmtFunc_js()


        htmlStr += `
                    </script>
                </head>

                <body>
                    <table border="0" bordercolordark="#F8F8FF" cellpadding="0" cellspacing="1" bgcolor="#FFFFFF"  style="font-size:small;">
                `

        // (2) make Table Action TR
                // colspan, f5_href, delete_href, fontSize
                //htmlStr += makeActionTR_html("4", "keypair", "deleteKeyPairMgmt()", "2")
                htmlStr += makeActionTR_html("4", "", "deleteKeyPairMgmt()", "2")


        // (3) make Table Header TR
                nameWidthList := []NameWidth {
                    {"Spider's NameId", "300"},
                    {"CSP's SystemId", "300"},                    
                }
                htmlStr +=  makeTitleTRList_html("#DDDDDD", "2", nameWidthList, true)


        // (4) make TR list with info list
        // (4-1) get info list 
                resBody, err := getAllResourceList_with_Connection_JsonByte(connConfig, "keypair")
                if err != nil {
                        cblog.Error(err)
                        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
                }
                
                var info cr.AllResourceList
                
                json.Unmarshal(resBody, &info)

        // (4-2) make TR list with info list
                htmlStr += makeMgmtTRList_html("", "", "", info)

        // make page tail
        htmlStr += `
                    </table>
            <hr>
                </body>
                </html>
        `

//fmt.Println(htmlStr)
        return c.HTML(http.StatusOK, htmlStr)
}

//====================================== VM

// make the string of javascript function
func makeDeleteVMMgmtFunc_js() string {
// delete for MappedList & OnlySpiderList
// curl -sX DELETE http://localhost:1024/spider/vm/vm-01?force=true -H 'Content-Type: application/json' -d '{ "ConnectionName": "aws-ohio-config"}
// delete for OnlyCSPList
// curl -sX DELETE http://localhost:1024/spider/cspvm/0b0d0d30794eab379 -H 'Content-Type: application/json' -d '{ "ConnectionName": "aws-ohio-config"}' |json_pp

        strFunc := `
                function deleteVMMgmt() {
                        var connConfig = parent.frames["top_frame"].document.getElementById("connConfig").innerHTML;
                        var checkboxes = document.getElementsByName('check_box');
                        for (var i = 0; i < checkboxes.length; i++) { // @todo make parallel executions
                                if (checkboxes[i].checked) {
                                        var xhr = new XMLHttpRequest();
                                        if(checkboxes[i].value.includes("::NAMEID::")) { // MappedList & OnlySpiderList
                                            xhr.open("DELETE", "$$SPIDER_SERVER$$/spider/vm/" + checkboxes[i].value.replace("::NAMEID::", "") + "?force=true", false);
                                        }else { // OnlyCSPList
                                            xhr.open("DELETE", "$$SPIDER_SERVER$$/spider/cspvm/" + checkboxes[i].value, false);
                                        }

                                        xhr.setRequestHeader('Content-Type', 'application/json');
                    sendJson = '{ "ConnectionName": "' + connConfig + '"}'
                                        xhr.send(sendJson);
                                }
                        }
            location.reload();
                }
        `
        strFunc = strings.ReplaceAll(strFunc, "$$SPIDER_SERVER$$", "http://" + cr.HostIPorName + cr.ServicePort) // cr.ServicePort = ":1024"
        return strFunc
}

func VMMgmt(c echo.Context) error {
        cblog.Info("call VMMgmt()")

    connConfig := c.Param("ConnectConfig")
    if connConfig == "region not set" {
        htmlStr :=  `
            <html>
            <head>
                <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                <script type="text/javascript">
                alert(connConfig)
                </script>
            </head>
            <body>
                <br>
                <br>
                <label style="font-size:24px;color:#606262;">&nbsp;&nbsp;&nbsp;Please select a Connection Configuration! (MENU: 2.CONNECTION)</label>
            </body>
        `

        return c.HTML(http.StatusOK, htmlStr)
    }

        // make page header
        htmlStr :=  `
                <html>
                <head>
                    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                    <script type="text/javascript">
                `
        // (1) make Javascript Function
                htmlStr += makeCheckBoxToggleFunc_js()
                htmlStr += makeDeleteVMMgmtFunc_js()


        htmlStr += `
                    </script>
                </head>

                <body>
                    <table border="0" bordercolordark="#F8F8FF" cellpadding="0" cellspacing="1" bgcolor="#FFFFFF"  style="font-size:small;">
                `

        // (2) make Table Action TR
                // colspan, f5_href, delete_href, fontSize
                //htmlStr += makeActionTR_html("4", "vm", "deleteVMMgmt()", "2")
                htmlStr += makeActionTR_html("4", "", "deleteVMMgmt()", "2")


        // (3) make Table Header TR
                nameWidthList := []NameWidth {
                    {"Spider's NameId", "300"},
                    {"CSP's SystemId", "300"},
                }
                htmlStr +=  makeTitleTRList_html("#DDDDDD", "2", nameWidthList, true)


        // (4) make TR list with info list
        // (4-1) get info list
                resBody, err := getAllResourceList_with_Connection_JsonByte(connConfig, "vm")
                if err != nil {
                        cblog.Error(err)
                        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
                }

                var info cr.AllResourceList

                json.Unmarshal(resBody, &info)

        // (4-2) make TR list with info list
                htmlStr += makeMgmtTRList_html("", "", "", info)

        // make page tail
        htmlStr += `
                    </table>
            <hr>
                </body>
                </html>
        `

//fmt.Println(htmlStr)
        return c.HTML(http.StatusOK, htmlStr)
}

