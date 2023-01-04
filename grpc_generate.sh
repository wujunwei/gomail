#!/usr/bin/env sh

protoc --go_out=pkg/mailbox/ --go_opt=paths=source_relative  \
    --go-grpc_out=pkg/mailbox/ --go-grpc_opt=paths=source_relative \
    --proto_path=proto/ \
     proto/mail.proto