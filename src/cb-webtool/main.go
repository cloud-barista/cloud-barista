package main

import (
	"html/template"
	"io"
	"net/http"

	controller "github.com/cloud-barista/cb-webtool/src/controller"
	echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func init() {

}

type user struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

type connectionInfo struct {
	RegionName     string `json:"regionname"`
	ConfigName     string
	ProviderName   string `json:"ProviderName"`
	CredentialName string `json:"CredentialName"`
	DriverName     string `json:"DriverName"`
}

type TemplateRender struct {
	templates *template.Template
}

func (t *TemplateRender) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func requestApi(method string, restUrl string, body io.Reader) {

}

func main() {
	e := echo.New()

	e.Use(echosession.New())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://210.207.104.150"},
	// 	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	// }))

	e.Static("/assets", "./src/static/assets")

	// paresGlob 를 사용하여 모든 경로에 있는 파일을 가져올 경우 사용하면 되겠다.
	// 사용한다음에 해당 파일을 불러오면 되네.
	// 서브디렉토리에 있는 걸 확인하기가 힘드네....
	renderer := &TemplateRender{
		templates: template.Must(template.ParseGlob(`./src/views/*.html`)),
	}

	e.Renderer = renderer

	e.GET("/", controller.IndexController)
	e.GET("/dashboard", controller.DashBoard)

	//login 관련
	e.GET("/login", controller.LoginForm)
	e.POST("/login/proc", controller.LoginController)
	e.POST("/regUser", controller.RegUserConrtoller)
	e.GET("/logout", controller.LogoutForm)
	e.GET("/logout/proc", controller.LoginController)

	// Monitoring Control
	e.GET("/monitoring", controller.MornitoringListForm)

	// MCIS
	e.GET("/MCIS/reg", controller.McisRegForm)
	e.GET("/MCIS/reg/:mcis_id/:mcis_name", controller.VMAddForm)
	e.POST("/MCIS/reg/proc", controller.McisRegController)
	e.GET("/MCIS/list", controller.McisListForm)

	// Resource
	e.GET("/Resource/board", controller.ResourceBoard)

	// 웹툴에서 사용할 rest
	e.GET("/SET/NS/:nsid", controller.SetNameSpace)

	// 웹툴에서 처리할 NameSpace
	e.GET("/NS/list", controller.NsListForm)
	e.GET("/NS/reg", controller.NsRegForm)
	e.POST("/NS/reg/proc", controller.NsRegController)
	e.GET("/GET/ns", controller.GetNameSpace)

	// 웹툴에서 처리할 Connection
	e.GET("/Connection/list", controller.ConnectionListForm)
	e.GET("/Connection/reg", controller.ConnectionRegForm)
	e.POST("/Connection/reg/proc", controller.NsRegController)

	// 웹툴에서 처리할 Region
	e.GET("/Region/list", controller.RegionListForm)
	e.GET("/Region/reg", controller.RegionRegForm)
	e.POST("/Region/reg/proc", controller.NsRegController)

	// 웹툴에서 처리할 Credential
	e.GET("/Credential/list", controller.CredertialListForm)
	e.GET("/Credential/reg", controller.CredertialRegForm)

	// 웹툴에서 처리할 Driver
	e.GET("/Driver/list", controller.DriverListForm)
	e.GET("/Driver/reg", controller.DriverRegForm)

	// 웹툴에서 Select Pop
	e.GET("/Pop/spec", controller.PopSpec)

	e.Logger.Fatal(e.Start(":1234"))

}
