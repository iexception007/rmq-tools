#!/usr/bin/env bash

# go mod init rmq-tools
# go mod tidy
# go mod vendor

go build -v -o rmq-tools main.go
if [ $? -ne 0 ]; then
    exit 1
fi
#set ROCKETMQ_CLIENT_GO_LOG_LEVEL=off
export ROCKETMQ_GO_LOG_LEVEL=error
./rmq-tools --role=receiver