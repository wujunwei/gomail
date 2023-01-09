package mailbox

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"gomail/pkg/db"
	"gomail/pkg/imap"
	"gomail/pkg/proto"
	"gomail/pkg/smtp"
	"gomail/pkg/util/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

type DefaultMailBoxService struct {
	proto.UnimplementedMailBoxServer
	Watcher  imap.Watcher
	Registry db.Storage
	Tool     smtp.Tool
}

func (s *DefaultMailBoxService) Send(_ context.Context, t *proto.MailTask) (*proto.SendMailResponse, error) {
	task := smtp.MailTask{
		From:        AddressString(t.From),
		To:          AddressStrings(t.To),
		Cc:          AddressStrings(t.Cc),
		Bcc:         AddressStrings(t.Bcc),
		Subject:     t.Subject,
		ReplyId:     t.ReplyId,
		Body:        t.Text.MainBody,
		ContentType: t.Text.ContentType,
	}
	if t.Attachment != nil && t.Attachment.WithAttachment {
		file, err := s.Registry.Download(t.Attachment.AttachmentID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error happen %v", err)
		}
		task.Attachment = smtp.Attachment{
			File:     file,
			WithFile: true,
		}
	}
	msgID, err := s.Tool.Send(task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error happen %v", err)
	}
	return &proto.SendMailResponse{MsgID: msgID}, nil
}
func (s *DefaultMailBoxService) ListServer(context.Context, *empty.Empty) (*proto.ServerList, error) {
	resp := &proto.ServerList{}
	for _, name := range s.Watcher.ListServer() {
		resp.Items = append(resp.Items, &proto.Server{Name: name})
	}
	return resp, nil
}
func (s *DefaultMailBoxService) Upload(us proto.MailBox_UploadServer) error {
	uf, err := us.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "error happen %v", err)
	}
	errChan := make(chan error, 1)
	defer close(errChan)
	pr, pw := io.Pipe()
	go func() {
		defer func() { _ = pw.Close() }()
		for {
			uf, err := us.Recv()
			if err != nil {
				errChan <- err
				return
			}
			_, err = pw.Write(uf.GetContent())
			if err != nil {
				errChan <- err
				return
			}
		}
	}()
	id, err := s.Registry.Upload(uf.GetName(), uf.GetContentType(), pr)
	if err != nil {
		return err
	}
	err = <-errChan
	if err != nil {
		return err
	}
	return us.SendAndClose(&proto.UploadResponse{FileID: id})
}
func (s *DefaultMailBoxService) Watch(ser *proto.Server, ws proto.MailBox_WatchServer) error {
	done := make(chan error)
	msgChan := make(chan *proto.Mail, 50)
	id, _ := random.Alpha(16)
	err := s.Watcher.Subscribe(ser.GetName(), string(id), msgChan)
	if err != nil {
		return err
	}
	defer func() {
		s.Watcher.UnSubscribe(ser.GetName(), string(id))
		close(msgChan)
	}()
	for {
		select {
		case msg := <-msgChan:
			{
				err := ws.Send(msg)
				if err != nil {
					return err
				}
			}
		case err := <-done:
			{
				if err != nil {
					log.Println(err, "client clean up !")
					return err
				}
			}
		}
	}
}

func NewMailBoxService(watcher imap.Watcher, client smtp.Tool, storage db.Storage) *DefaultMailBoxService {
	return &DefaultMailBoxService{Watcher: watcher, Tool: client, Registry: storage}
}
