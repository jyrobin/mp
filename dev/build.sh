#!/usr/bin/env bash                                                                       

MP_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && cd .. && pwd )"
[ -z $GOBIN ] && GOBIN=~/go/bin

echo go build -o $GOBIN/mt $MP_DIR/cmd/mt/main.go 
go build -o $GOBIN/mt $MP_DIR/cmd/mt/main.go 

echo go build -o $GOBIN/mj $MP_DIR/cmd/mj/main.go 
go build -o $GOBIN/mj $MP_DIR/cmd/mj/main.go 
