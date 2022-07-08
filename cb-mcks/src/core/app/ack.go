package app

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/core/validation"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"
	"github.com/labstack/echo/v4"
)

func SendMessage(c echo.Context, httpCode int, msg string) error {
	status := Status{Kind: KIND_STATUS, Code: httpCode, Message: msg}
	return c.JSON(httpCode, status)
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
			return errors.New(fmt.Sprintf("[%s] %s", err.Key, err.Error()))
		}
	}
	return nil
}

func NewStatus(code int, message string) *Status {
	return &Status{
		Kind:    KIND_STATUS,
		Code:    code,
		Message: message,
	}
}

func ClusterReqDef(clusterReq ClusterReq) {
	clusterReq.Config.Kubernetes.PodCidr = lang.NVL(clusterReq.Config.Kubernetes.PodCidr, POD_CIDR)
	clusterReq.Config.Kubernetes.ServiceCidr = lang.NVL(clusterReq.Config.Kubernetes.ServiceCidr, SERVICE_CIDR)
	clusterReq.Config.Kubernetes.ServiceDnsDomain = lang.NVL(clusterReq.Config.Kubernetes.ServiceDnsDomain, SERVICE_DOMAIN)
}

func ClusterReqValidate(req ClusterReq) error {
	if len(req.ControlPlane) == 0 {
		return errors.New("Control plane node must be at least one")
	}
	if len(req.ControlPlane) > 1 {
		return errors.New("Only one control plane node is supported")
	}
	if len(req.Worker) == 0 {
		return errors.New("Worker node must be at least one")
	}
	if !(req.Config.Kubernetes.NetworkCni == NETWORKCNI_CANAL || req.Config.Kubernetes.NetworkCni == NETWORKCNI_KILO) {
		return errors.New("Network-cni allows only canal or kilo")
	}

	if len(req.Name) == 0 {
		return errors.New("Cluster name is empty")
	} else {
		err := lang.VerifyClusterName(req.Name)
		if err != nil {
			return err
		}
	}

	if len(req.Config.Kubernetes.PodCidr) > 0 {
		err := lang.VerifyCIDR("podCidr", req.Config.Kubernetes.PodCidr)
		if err != nil {
			return err
		}
	}
	if len(req.Config.Kubernetes.ServiceCidr) > 0 {
		err := lang.VerifyCIDR("serviceCidr", req.Config.Kubernetes.ServiceCidr)
		if err != nil {
			return err
		}
	}

	return nil
}

func NodeReqValidate(req NodeReq) error {
	if len(req.ControlPlane) > 0 {
		return errors.New("Control plane node is not supported")
	}
	if len(req.Worker) == 0 {
		return errors.New("Worker node must be at least one")
	}

	return nil
}
