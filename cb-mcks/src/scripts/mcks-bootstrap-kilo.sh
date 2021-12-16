#!/bin/bash
# 스크립트 생성
echo -e '#!/bin/sh
ifconfig $(ip route get 8.8.8.8 | awk \047{ print $5; exit }\047):1 {{PUBLICIP}} netmask 255.255.255.255 broadcast 0.0.0.0 up
echo "KUBELET_EXTRA_ARGS=-\"-node-ip={{PUBLICIP}}\"" > /etc/default/kubelet
if [ -f "/etc/kubernetes/kubelet.conf" ]; then
  systemctl restart kubelet
  kubectl --kubeconfig=/etc/kubernetes/kubelet.conf annotate node $(hostname | awk \047{print tolower($0)}\047) kilo.squat.ai/location=$(hostname | awk \047{print tolower($0)}\047) --overwrite
  kubectl --kubeconfig=/etc/kubernetes/kubelet.conf annotate node $(hostname | awk \047{print tolower($0)}\047) kilo.squat.ai/force-endpoint={{PUBLICIP}}:51820 --overwrite
  kubectl --kubeconfig=/etc/kubernetes/kubelet.conf annotate node $(hostname | awk \047{print tolower($0)}\047) kilo.squat.ai/persistent-keepalive=25 --overwrite
fi
exit 0
fi' | sed "s/{{PUBLICIP}}/$1/g" | sudo tee /lib/systemd/system/mcks-bootstrap > /dev/null

# 실행권한
sudo chmod +x /lib/systemd/system/mcks-bootstrap
