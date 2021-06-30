#!/bin/bash
# kubeadm-config 정의
# - controlPlaneEndpoint 에 LB 지정 (9998 포트)
# - advertise-address 에 Public IP 지정
cat << EOF > kubeadm-config.yaml
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
imageRepository: k8s.gcr.io
controlPlaneEndpoint: $(dig +short myip.opendns.com @resolver1.opendns.com):9998
dns:
  type: CoreDNS
apiServer:
  extraArgs:
    advertise-address: $(dig +short myip.opendns.com @resolver1.opendns.com)
    authorization-mode: Node,RBAC
etcd:
  local:
    dataDir: /var/lib/etcd
networking:
  dnsDomain: $3
  podSubnet: $1
  serviceSubnet: $2
controllerManager: {}
scheduler: {}
EOF

# Control-plane init
sudo kubeadm init --v=5 --upload-certs --config kubeadm-config.yaml

# control-plane leader 의 경우
# - ladybug-bootstrap 데몬이 자동 실행
#systemctl status ladybug-bootstrap
