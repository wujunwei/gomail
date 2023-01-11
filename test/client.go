package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"gomail/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	cli := proto.NewMailBoxClient(conn)
	server, err := cli.ListServer(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(server.Items)
}
