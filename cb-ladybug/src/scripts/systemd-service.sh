#!/bin/bash
# sysatemd for persistent ip address
sudo bash -c 'cat > /lib/systemd/system/ladybug-bootstrap.service <<EOF
[Unit]
Description=Ladybug bootstrap script
After=multi-user.target
StartLimitIntervalSec=60
StartLimitBurst=3
[Service]
ExecStart=/lib/systemd/system/ladybug-bootstrap
Restart=on-failure
RestartSec=10s
[Install]
WantedBy=kubelet.service
EOF'

# reload 및 reboot 시 실행하도록 지정
sudo systemctl daemon-reload
sudo systemctl enable ladybug-bootstrap
