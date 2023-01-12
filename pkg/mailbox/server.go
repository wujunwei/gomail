package mailbox

import (
	"google.golang.org/grpc"
)

func NewGRPCServer(opts ...grpc.ServerOption) *grpc.Server {
	s := grpc.NewServer(opts...)

	return s
}
