package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"gomail/pkg/mailbox"
	"gomail/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
)

func main() {
	user := "example_user"
	password := "example_password"
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: mailbox.NewBasicToken(user, password)}),
	}
	conn, err := grpc.Dial("localhost:5000", opts...)
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
