#!/bin/bash
export GO_PATH=~/go
export PATH=$PATH:/$GO_PATH/bin
protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=.