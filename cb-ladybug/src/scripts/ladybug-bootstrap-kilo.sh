#!/bin/bash
# 스크립트 생성
echo -e '#!/bin/sh
ifconfig $(ip route get 8.8.8.8 | awk \047{ print $5; exit }\047):1 $(dig +short myip.opendns.com @resolver1.opendns.com) netmask 255.255.255.255 broadcast 0.0.0.0 up
echo "KUBELET_EXTRA_ARGS=-\"-node-ip=$(dig +short myip.opendns.com @resolver1.opendns.com)\"" > /etc/default/kubelet
if [ -f "/etc/kubernetes/kubelet.conf" ]; then
  systemctl restart kubelet
  kubectl --kubeconfig=/etc/kubernetes/kubelet.conf annotate node $(hostname) kilo.squat.ai/location=$(hostname) --overwrite
  kubectl --kubeconfig=/etc/kubernetes/kubelet.conf annotate node $(hostname) kilo.squat.ai/force-endpoint=$(dig +short myip.opendns.com @resolver1.opendns.com):51820 --overwrite
  kubectl --kubeconfig=/etc/kubernetes/kubelet.conf annotate node $(hostname) kilo.squat.ai/persistent-keepalive=25 --overwrite
fi
exit 0
fi' | sudo tee /lib/systemd/system/ladybug-bootstrap > /dev/null

# 실행권한
sudo chmod +x /lib/systemd/system/ladybug-bootstrap
