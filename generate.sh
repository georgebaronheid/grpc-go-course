#!/usr/bin/env bash

#Generate pb for greet.proto
#protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=.

protoc sum/sumpb/sum.proto --go_out=. --go-grpc_out=.