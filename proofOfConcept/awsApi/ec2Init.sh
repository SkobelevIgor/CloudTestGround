#!/bin/bash

wget https://raw.githubusercontent.com/SkobelevIgor/CloudTestGround/dev/proofOfConcept/awsApi/server/server.go

sudo yum install -y go

echo "export GOPATH=~" >> ~/.bash_profile
echo "export GOBIN=~/bin" >> ~/.bash_profile
export GOPATH=~
export GOBIN=~/bin

go install server.go

cat << EOF >> /lib/systemd/system/server-agent.service
[Unit]
Description=server-agent

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=$GOBIN/server

[Install]
WantedBy=multi-user.target
EOF

service server-agent start
service server-agent enable



