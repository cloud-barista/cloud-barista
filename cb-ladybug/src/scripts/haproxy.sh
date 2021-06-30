#!/bin/bash
sudo add-apt-repository -y ppa:vbernat/haproxy-1.7
sudo apt update
sudo apt install -y haproxy

sudo bash -c "cat << EOF > /etc/haproxy/haproxy.cfg
global
  log 127.0.0.1 local0
  maxconn 2000
  uid 0
  gid 0
  daemon
defaults
  log global
  mode tcp
  option dontlognull
  timeout connect 5000ms
  timeout client 50000ms
  timeout server 50000ms
frontend apiserver
  bind :9998
  default_backend apiserver
backend apiserver
  balance roundrobin
# {{SERVERS}} will be replaced to provisioned kubernetes api-servers' ip addresses as below:
#   server  api1  111.222.333.444:6443  check
# DO NOT EDIT {{SERVERS}}
{{SERVERS}}
EOF"

# haproxy 재시작
sudo systemctl restart haproxy
