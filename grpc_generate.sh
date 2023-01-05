#!/usr/bin/env sh

protoc --go_out=pkg/proto/ --go_opt=paths=source_relative  \
    --go-grpc_out=pkg/proto/ --go-grpc_opt=paths=source_relative \
    --proto_path=pkg/proto/ \
     pkg/proto/mail.proto