package router

import (
	"net/http"
	"time"

	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/cloud-barista/cb-mcks/src/core/service"
	"github.com/labstack/echo/v4"

	logger "github.com/sirupsen/logrus"
)

// ListCluster godoc
// @Tags Cluster
// @Summary List all Clusters
// @Description List all Clusters
// @ID ListCluster
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "Namespace ID"
// @Success 200 {object} model.ClusterList
// @Failure 400 {object} app.Status
// @Router /ns/{namespace}/clusters [get]
func ListCluster(c echo.Context) error {
	clusterList, err := service.ListCluster(c.Param("namespace"))
	if err != nil {
		logger.Warnf("(ListCluster) %s'", err.Error())
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	return app.Send(c, http.StatusOK, clusterList)
}

// GetCluster godoc
// @Tags Cluster
// @Summary Get Cluster
// @Description Get Cluster
// @ID GetCluster
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "Namespace ID"
// @Param	cluster	path	string	true  "Cluster Name"
// @Success 200 {object} model.Cluster
// @Failure 400 {object} app.Status
// @Failure 404 {object} app.Status
// @Router /ns/{namespace}/clusters/{cluster} [get]
func GetCluster(c echo.Context) error {
	if err := app.Validate(c, []string{"namespace", "cluster"}); err != nil {
		logger.Warnf("(CreateCluster) %s", err.Error())
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	cluster, err := service.GetCluster(c.Param("namespace"), c.Param("cluster"))
	if err != nil {
		logger.Warnf("(GetCluster) %s", err.Error())
		return app.SendMessage(c, http.StatusNotFound, err.Error())
	}

	return app.Send(c, http.StatusOK, cluster)
}

// CreateCluster godoc
// @Tags Cluster
// @Summary Create Cluster
// @Description Create Cluster
// @ID CreateCluster
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "Namespace ID"
// @Param   minorversion  query    string   true  "string enums"    Enums(1.18, 1.23)
// @Param   patchversion  path	int	true  "Patch version"
// @Param ClusterReq body app.ClusterReq true "Request Body to create cluster"
// @Success 200 {object} model.Cluster
// @Failure 400 {object} app.Status
// @Failure 500 {object} app.Status
// @Router /ns/{namespace}/clusters [post]
func CreateCluster(c echo.Context) error {
	start := time.Now()
	clusterReq := &app.ClusterReq{}
	if err := c.Bind(clusterReq); err != nil {
		logger.Warnf("(CreateCluster) %s", err.Error())
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	app.ClusterReqDef(*clusterReq)

	err := app.ClusterReqValidate(*clusterReq)
	if err != nil {
		logger.Warnf("(CreateCluster) %s", err.Error())
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}
	cluster, err := service.CreateCluster(c.Param("namespace"), c.QueryParam("minorversion"), c.QueryParam("patchversion"), clusterReq)
	if err != nil {
		logger.Warnf("(CreateCluster) %s", err.Error())
		return app.SendMessage(c, http.StatusInternalServerError, err.Error())
	}

	logger.Info("(CreateCluster) Duration = ", time.Since(start))
	return app.Send(c, http.StatusOK, cluster)
}

// DeleteCluster godoc
// @Tags Cluster
// @Summary Delete Cluster
// @Description Delete Cluster
// @ID DeleteCluster
// @Accept json
// @Produce json
// @Param	namespace	path	string	true  "Namespace ID"
// @Param	cluster	path	string	true  "Cluster Name"
// @Success 200 {object} app.Status
// @Failure 400 {object} app.Status
// @Failure 500 {object} app.Status
// @Router /ns/{namespace}/clusters/{cluster} [delete]
func DeleteCluster(c echo.Context) error {
	start := time.Now()

	if err := app.Validate(c, []string{"namespace", "cluster"}); err != nil {
		logger.Warnf("(DeleteCluster) %s", err.Error())
		return app.SendMessage(c, http.StatusBadRequest, err.Error())
	}

	status, err := service.DeleteCluster(c.Param("namespace"), c.Param("cluster"))
	if err != nil {
		logger.Warnf("(DeleteCluster) %s", err.Error())
		return app.SendMessage(c, http.StatusInternalServerError, err.Error())
	}

	logger.Info("(DeleteCluster) Duration = ", time.Since(start))
	return app.Send(c, http.StatusOK, status)
}
