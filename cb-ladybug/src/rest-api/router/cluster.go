package router

import (
	"net/http"
	"time"

	"github.com/cloud-barista/cb-ladybug/src/core/model"
	"github.com/cloud-barista/cb-ladybug/src/core/service"
	"github.com/cloud-barista/cb-ladybug/src/utils/app"
	"github.com/labstack/echo/v4"

	logger "github.com/sirupsen/logrus"
)

// ListCluster
// @Tags Cluster
// @Summary List Cluster
// @Description List Cluster
// @ID ListCluster
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "namespace"
// @Success 200 {object} model.ClusterList
// @Router /ns/{namespace}/clusters [get]
func ListCluster(c echo.Context) error {
	clusterList, err := service.ListCluster(c.Param("namespace"))
	if err != nil {
		logger.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	return app.Send(c, http.StatusOK, clusterList)
}

// GetCluster
// @Tags Cluster
// @Summary Get Cluster
// @Description Get Cluster
// @ID GetCluster
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "namespace"
// @Param	cluster	path	string	true  "cluster"
// @Success 200 {object} model.Cluster
// @Router /ns/{namespace}/clusters/{cluster} [get]
func GetCluster(c echo.Context) error {
	if err := app.Validate(c, []string{"namespace", "cluster"}); err != nil {
		logger.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	cluster, err := service.GetCluster(c.Param("namespace"), c.Param("cluster"))
	if err != nil {
		logger.Infof("not found a cluster (namespace=%s, cluster=%s, cause=%s)", c.Param("namespace"), c.Param("cluster"), err)
		return app.SendMessage(c, http.StatusNotFound, err.Error())
	}

	return app.Send(c, http.StatusOK, cluster)
}

// CreateCluster
// @Tags Cluster
// @Summary Create Cluster
// @Description Create Cluster
// @ID CreateCluster
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "namespace"
// @Param json body model.ClusterReq true "Reuest json"
// @Success 200 {object} model.Cluster
// @Router /ns/{namespace}/clusters [post]
func CreateCluster(c echo.Context) error {
	start := time.Now()
	clusterReq := &model.ClusterReq{}
	if err := c.Bind(clusterReq); err != nil {
		logger.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	app.ClusterReqDef(*clusterReq)

	err := app.ClusterReqValidate(*clusterReq)
	if err != nil {
		logger.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	cluster, err := service.CreateCluster(c.Param("namespace"), clusterReq)
	if err != nil {
		logger.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	duration := time.Since(start)
	logger.Info("duration := ", duration)
	return app.Send(c, http.StatusOK, cluster)
}

// DeleteCluster
// @Tags Cluster
// @Summary Delete a cluster
// @Description Delete a cluster
// @ID DeleteCluster
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "namespace"
// @Param	cluster	path	string	true  "cluster"
// @Success 200 {object} model.Status
// @Router /ns/{namespace}/clusters/{cluster} [delete]
func DeleteCluster(c echo.Context) error {
	if err := app.Validate(c, []string{"namespace", "cluster"}); err != nil {
		logger.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	status, err := service.DeleteCluster(c.Param("namespace"), c.Param("cluster"))
	if err != nil {
		logger.Error(err)
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	return app.Send(c, http.StatusOK, status)
}
