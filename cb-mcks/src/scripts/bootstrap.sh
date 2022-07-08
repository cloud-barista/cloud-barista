#!/bin/bash
K8S_VERSION="$1"	# curl https://packages.cloud.google.com/apt/dists/kubernetes-xenial/main/binary-amd64/Packages
CSP="$2"
HOSTNAME="$3"
PUBLIC_IP="$4"			# openstack
NETWORK_CNI="$5"

# hostname
sudo hostnamectl set-hostname ${HOSTNAME}

if [[ "${K8S_VERSION}" == "1.23"* ]]; then 

sudo swapoff -a && sed -i '/swap/s/^/#/' /etc/fstab
# br_netfilter
sudo cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

sudo cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sudo sysctl --system

#  default packages
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg2

# apg-get add repoisities
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add
sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
sudo apt-get update

# container runtime
sudo apt-get install -y containerd.io=1.2.13-2
cat <<EOF | sudo tee /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter

cat <<EOF | sudo tee /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-iptables  = 1
net.ipv4.ip_forward                 = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF
sudo sysctl --system

sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
sudo systemctl restart containerd
fi

if [[ "${K8S_VERSION}" == "1.18"* ]]; then 
# packages
`echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections`
sudo killall apt apt-get > /dev/null 2>&1
sudo rm -vf /var/lib/apt/lists/lock
sudo rm -vf /var/cache/apt/archives/lock
sudo rm -vf /var/lib/dpkg/lock*
sudo dpkg --configure -a

sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg2

sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y containerd.io=1.2.13-2 docker-ce=5:19.03.11~3-0~ubuntu-$(lsb_release -cs) docker-ce-cli=5:19.03.11~3-0~ubuntu-$(lsb_release -cs)


sudo bash -c 'cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF'

sudo swapoff -a && sed -i '/swap/s/^/#/' /etc/fstab

sudo mkdir -p /etc/systemd/system/docker.service.d
sudo systemctl daemon-reload
sudo systemctl restart docker

sudo curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add
sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
sudo apt-get update
fi

# kubeadm , kubelet, kubectl
sudo apt-get install -y kubeadm=${K8S_VERSION} kubelet=${K8S_VERSION} kubectl=${K8S_VERSION}
sudo apt-mark hold kubeadm kubelet kubectl

if [ "${CSP}" != "openstack" ]; then 
	PUBLIC_IP='$(dig +short myip.opendns.com @resolver1.opendns.com)'
fi

if [ "${NETWORK_CNI}" == "kilo" ]; then 
# install wireguard
sudo add-apt-repository -y ppa:wireguard/wireguard
sudo apt-get update
sudo apt-get install -y wireguard
# mcks-bootstrap
echo -e '#!/bin/sh
IFACE="$(ip route get 8.8.8.8 | awk \047{ print $5; exit }\047)"
PUBLIC_IP="{{PUBLIC_IP}}"
ifconfig ${IFACE}:1 ${PUBLIC_IP} netmask 255.255.255.255  broadcast 0.0.0.0 up
echo "KUBELET_EXTRA_ARGS=-\"-node-ip=${PUBLIC_IP}\"" > /etc/default/kubelet
if [ -f "/etc/kubernetes/kubelet.conf" ]; then
  systemctl restart kubelet
  kubectl --kubeconfig=/etc/kubernetes/kubelet.conf annotate node {{HOSTNAME}} kilo.squat.ai/force-endpoint=${PUBLIC_IP}:51820 --overwrite
fi
exit 0
fi' | sed "s/{{HOSTNAME}}/${HOSTNAME}/g" | sed "s/{{PUBLIC_IP}}/${PUBLIC_IP}/g" | sudo tee /lib/systemd/system/mcks-bootstrap > /dev/null
sudo chmod +x /lib/systemd/system/mcks-bootstrap
fi

if [ "${NETWORK_CNI}" == "canal" ]; then 
# mcks-bootstrap
echo -e '#!/bin/sh
IFACE="$(ip route get 8.8.8.8 | awk \047{ print $5; exit }\047)"
PUBLIC_IP="{{PUBLIC_IP}}"
ifconfig ${IFACE}:1 ${PUBLIC_IP} netmask 255.255.255.255  broadcast 0.0.0.0 up
echo "KUBELET_EXTRA_ARGS=-\"-node-ip=${PUBLIC_IP}\"" > /etc/default/kubelet
if [ -f "/etc/kubernetes/kubelet.conf" ]; then
  systemctl restart kubelet
  kubectl --kubeconfig=/etc/kubernetes/kubelet.conf annotate node {{HOSTNAME}} projectcalico.org/IPv4Address=${PUBLIC_IP} --overwrite
  R="$(kubectl --kubeconfig=/etc/kubernetes/kubelet.conf annotate node {{HOSTNAME}} flannel.alpha.coreos.com/public-ip-overwrite=${PUBLIC_IP} --overwrite)"
  if echo "$R" | grep "annotated"; then
    R=$(kubectl --kubeconfig=/etc/kubernetes/kubelet.conf get nodes --no-headers | awk \047END { print NR }\047)
    echo "nodes count = ${R}"
    if [ "$R" != "1" ]; then
      systemctl restart docker
      echo "docker daemon restarted"
    fi
    exit 0
  else
    exit 1
  fi
fi
exit 0
fi' | sed "s/{{HOSTNAME}}/${HOSTNAME}/g" | sed "s/{{PUBLIC_IP}}/${PUBLIC_IP}/g" | sudo tee /lib/systemd/system/mcks-bootstrap > /dev/null
sudo chmod +x /lib/systemd/system/mcks-bootstrap
fi

# setup bootstrap service deamon
sudo bash -c 'cat > /lib/systemd/system/mcks-bootstrap.service <<EOF
[Unit]
Description=MCKS bootstrap script
After=multi-user.target
StartLimitIntervalSec=60
StartLimitBurst=3
[Service]
ExecStart=/lib/systemd/system/mcks-bootstrap
Restart=on-failure
RestartSec=10s
[Install]
WantedBy=kubelet.service
EOF'
sudo systemctl daemon-reload
sudo systemctl enable mcks-bootstrap