#!/bin/bash

wget https://raw.githubusercontent.com/SkobelevIgor/CloudTestGround/dev/proofOfConcept/awsApi/server/server.go

sudo yum install -y go

echo "export GOPATH=~" >> ~/.bash_profile
echo "export GOBIN=~/bin" >> ~/.bash_profile
export GOPATH=~
export GOBIN=~/bin

go install server.go

chmod u+x $GOBIN/server

cat << EOF >> /lib/systemd/system/server-agent.service
[Unit]
Description=server-agent

[Service]
Type=simple
Restart=on-failure
RestartSec=5s
ExecStart=$GOBIN/server

[Install]
WantedBy=multi-user.target
EOF

systemctl start server-agent
systemctl enable server-agent


