#!/bin/bash
# sysatemd for persistent ip address
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

# reload 및 reboot 시 실행하도록 지정
sudo systemctl daemon-reload
sudo systemctl enable mcks-bootstrap
