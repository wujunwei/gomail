package mailbox

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DefaultMailBoxServer struct {
}

func (s DefaultMailBoxServer) Send(context.Context, *Mail) (*SendMailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (s DefaultMailBoxServer) ListServer(context.Context, *empty.Empty) (*ServerList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListServer not implemented")
}
func (s DefaultMailBoxServer) Upload(MailBox_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (s DefaultMailBoxServer) Watch(*Server, MailBox_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (s DefaultMailBoxServer) mustEmbedUnimplementedMailBoxServer() {

}
