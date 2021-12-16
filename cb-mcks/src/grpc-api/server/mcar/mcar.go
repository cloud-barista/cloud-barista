package mcar

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/core/validation"
	"github.com/cloud-barista/cb-mcks/src/core/model"
	"github.com/cloud-barista/cb-mcks/src/utils/config"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// MCARService - MCKS 서비스 구현
type MCARService struct {
}

// ===== [ Implementations ] =====

func (s *MCARService) Validate(params map[string]string) error {
	valid := validation.Validation{}

	for key, element := range params {
		valid.Required(element, key)
	}

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(fmt.Sprintf("[%s]%s", err.Key, err.Error()))
		}
	}
	return nil
}

func (s *MCARService) ClusterReqDef(clusterReq *model.ClusterReq) error {
	clusterReq.Config.Kubernetes.NetworkCni = lang.NVL(clusterReq.Config.Kubernetes.NetworkCni, config.NETWORKCNI_KILO)
	clusterReq.Config.Kubernetes.PodCidr = lang.NVL(clusterReq.Config.Kubernetes.PodCidr, config.POD_CIDR)
	clusterReq.Config.Kubernetes.ServiceCidr = lang.NVL(clusterReq.Config.Kubernetes.ServiceCidr, config.SERVICE_CIDR)
	clusterReq.Config.Kubernetes.ServiceDnsDomain = lang.NVL(clusterReq.Config.Kubernetes.ServiceDnsDomain, config.SERVICE_DOMAIN)

	return nil
}

func (s *MCARService) ClusterReqValidate(req model.ClusterReq) error {
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

func (s *MCARService) NodeReqValidate(req model.NodeReq) error {
	if len(req.ControlPlane) > 0 {
		return errors.New("control plane node is not supported")
	}
	if len(req.Worker) == 0 {
		return errors.New("worker node must be at least one")
	}

	return nil
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
