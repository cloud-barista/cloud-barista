package config

type CSP string

const (
	CONTROL_PLANE = "control-plane"
	WORKER        = "worker"

	BOOTSTRAP_FILE = "bootstrap.sh"
	INIT_FILE      = "k8s-init.sh"

	GCP_IMAGE_ID = "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1804-bionic-v20201014"

	CSP_AWS CSP = "aws"
	CSP_GCP CSP = "gcp"
)
