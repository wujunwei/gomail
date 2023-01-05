package mailbox

import (
	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	return s
}
