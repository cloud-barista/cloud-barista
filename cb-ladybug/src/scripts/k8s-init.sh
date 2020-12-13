# kubeadm init
sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --token-ttl=0
sudo kubectl apply -f https://docs.projectcalico.org/manifests/calico.yaml --kubeconfig=/etc/kubernetes/admin.conf
# direct run & get variables for k8s join
END_POINT="$(hostname -i):6443"; echo "${END_POINT}"
TOKEN="$(sudo kubeadm token list | tail -n 1 | awk '{print $1}')"; echo "${TOKEN}"
HASH=$(openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'); echo "${HASH}"