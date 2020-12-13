package router

import (
	"net/http"
	"time"

	"github.com/cloud-barista/cb-ladybug/src/core/common"
	"github.com/cloud-barista/cb-ladybug/src/core/model"
	"github.com/cloud-barista/cb-ladybug/src/rest-api/service"
	"github.com/cloud-barista/cb-ladybug/src/utils/app"

	"github.com/labstack/echo/v4"
)

// ListNode
// @Tags Node
// @Summary List Node
// @Description List Node
// @ID ListNode
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "namespace"
// @Param	cluster	path	string	true  "cluster"
// @Success 200 {object} model.NodeList
// @Router /ns/{namespace}/clusters/{cluster}/nodes [get]
func ListNode(c echo.Context) error {
	if err := app.Validate(c, []string{"cluster"}); err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	nodeList, err := service.ListNode(c.Param("namespace"), c.Param("cluster"))
	if err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	return app.Send(c, http.StatusOK, nodeList)
}

// GetNode
// @Tags Node
// @Summary Get Node
// @Description Get Node
// @ID GetNode
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "namespace"
// @Param	cluster	path	string	true  "cluster"
// @Param	node	path	string	true  "node"
// @Success 200 {object} model.Node
// @Router /ns/{namespace}/clusters/{cluster}/nodes/{node} [get]
func GetNode(c echo.Context) error {
	if err := app.Validate(c, []string{"cluster", "node"}); err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	node, err := service.GetNode(c.Param("namespace"), c.Param("cluster"), c.Param("node"))
	if err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	return app.Send(c, http.StatusOK, node)
}

// AddNode
// @Tags Node
// @Summary Add Node
// @Description Add Node
// @ID AddNode
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "namespace"
// @Param	cluster	path	string	true  "cluster"
// @Param json body model.NodeReq true "Reuest json"
// @Success 200 {object} model.Node
// @Router /ns/{namespace}/clusters/{cluster}/nodes [post]
func AddNode(c echo.Context) error {
	start := time.Now()
	if err := app.Validate(c, []string{"cluster"}); err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	nodeReq := &model.NodeReq{}
	if err := c.Bind(nodeReq); err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	err := app.NodeReqValidate(c, *nodeReq)
	if err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	node, err := service.AddNode(c.Param("namespace"), c.Param("cluster"), nodeReq)
	if err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	duration := time.Since(start)
	common.CBLog.Info(" duration := ", duration)
	return app.Send(c, http.StatusOK, node)
}

// RemoveNode
// @Tags Node
// @Summary Remove Node
// @Description Remove Node
// @ID RemoveNode
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "namespace"
// @Param	cluster	path	string	true  "cluster"
// @Param	node	path	string	true  "node"
// @Success 200 {object} model.Status
// @Router /ns/{namespace}/clusters/{cluster}/nodes/{node} [delete]
func RemoveNode(c echo.Context) error {
	if err := app.Validate(c, []string{"cluster", "node"}); err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	status, err := service.RemoveNode(c.Param("namespace"), c.Param("cluster"), c.Param("node"))
	if err != nil {
		common.CBLog.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	return app.Send(c, http.StatusOK, status)
}
