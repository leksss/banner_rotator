#!/bin/bash

rm -rf ../../proto/protobuf
mkdir -p ../../proto/protobuf
protoc -I ../../proto/googleapis -I ../../proto/banner \
  --proto_path=../../proto/banner \
  --go_out=../../proto/protobuf \
  --go-grpc_out=../../proto/protobuf \
  --grpc-gateway_out=../../proto/protobuf \
  ../../proto/banner/*.proto
