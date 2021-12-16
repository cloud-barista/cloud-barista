package app

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/core/validation"
	"github.com/cloud-barista/cb-mcks/src/core/model"
	"github.com/cloud-barista/cb-mcks/src/utils/config"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"
	"github.com/labstack/echo/v4"
)

type Status struct {
	Message string `json:"message" example:"Any message"`
}

func SendMessage(c echo.Context, httpCode int, msg string) error {
	return c.JSON(httpCode, Status{Message: msg})
}

func Send(c echo.Context, httpCode int, json interface{}) error {
	return c.JSON(httpCode, json)
}

func Validate(c echo.Context, params []string) error {
	valid := validation.Validation{}

	for _, name := range params {
		valid.Required(c.Param(name), name)
	}

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(fmt.Sprintf("[%s]%s", err.Key, err.Error()))
		}
	}
	return nil
}

func ClusterReqDef(clusterReq model.ClusterReq) {
	clusterReq.Config.Kubernetes.NetworkCni = lang.NVL(clusterReq.Config.Kubernetes.NetworkCni, config.NETWORKCNI_KILO)
	clusterReq.Config.Kubernetes.PodCidr = lang.NVL(clusterReq.Config.Kubernetes.PodCidr, config.POD_CIDR)
	clusterReq.Config.Kubernetes.ServiceCidr = lang.NVL(clusterReq.Config.Kubernetes.ServiceCidr, config.SERVICE_CIDR)
	clusterReq.Config.Kubernetes.ServiceDnsDomain = lang.NVL(clusterReq.Config.Kubernetes.ServiceDnsDomain, config.SERVICE_DOMAIN)
}

func ClusterReqValidate(req model.ClusterReq) error {
	if len(req.ControlPlane) == 0 {
		return errors.New("control plane node must be at least one")
	}
	if len(req.ControlPlane) > 1 {
		return errors.New("only one control plane node is supported")
	}
	if len(req.Worker) == 0 {
		return errors.New("worker node must be at least one")
	}
	if !(req.Config.Kubernetes.NetworkCni == config.NETWORKCNI_CANAL || req.Config.Kubernetes.NetworkCni == config.NETWORKCNI_KILO) {
		return errors.New("network cni allows only canal or kilo")
	}

	if len(req.Name) == 0 {
		return errors.New("cluster name is empty")
	} else {
		err := lang.CheckName(req.Name)
		if err != nil {
			return err
		}
	}

	if len(req.Config.Kubernetes.PodCidr) > 0 {
		err := lang.CheckIpCidr("podCidr", req.Config.Kubernetes.PodCidr)
		if err != nil {
			return err
		}
	}
	if len(req.Config.Kubernetes.ServiceCidr) > 0 {
		err := lang.CheckIpCidr("serviceCidr", req.Config.Kubernetes.ServiceCidr)
		if err != nil {
			return err
		}
	}

	return nil
}

func NodeReqValidate(req model.NodeReq) error {
	if len(req.ControlPlane) > 0 {
		return errors.New("control plane node is not supported")
	}
	if len(req.Worker) == 0 {
		return errors.New("worker node must be at least one")
	}

	return nil
}
