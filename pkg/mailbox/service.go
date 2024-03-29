package mailbox

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"gomail/pkg/db"
	"gomail/pkg/imap"
	"gomail/pkg/mailbox/auth"
	"gomail/pkg/proto"
	"gomail/pkg/smtp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"sync"
)

type DefaultMailBoxService struct {
	proto.UnimplementedMailBoxServer
	Watcher  imap.Watcher
	Registry db.Storage
	Session  db.Session
	Tool     smtp.Sender
	lock     sync.Mutex
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
	md, ok := metadata.FromIncomingContext(ws.Context())
	if !ok {
		return status.Error(codes.Unknown, "header not found")
	}
	temp := md.Get("UserID")
	if len(temp) == 0 {
		return status.Error(codes.Unknown, "user not found")
	}
	id := temp[0]
	u := &auth.User{}
	err := s.Session.Get(map[string]interface{}{"_id": id}, u)
	if err != nil {
		return err
	}
	done := make(chan error)
	msgChan := make(chan *proto.Mail, 50)
	sub, err := s.Watcher.Subscribe(ser.GetName(), id, u.Weight, msgChan)
	if err != nil {
		return err
	}
	defer func() {
		s.Watcher.UnSubscribe(sub)
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

func (s *DefaultMailBoxService) Register(_ context.Context, u *proto.User) (*proto.UserResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.Session.Exist(map[string]interface{}{"name": u.Name, "password": u.Password}) {
		return nil, status.Error(codes.AlreadyExists, "user existed")
	}
	id, err := s.Session.Set(&auth.User{Password: u.Password, Name: u.Name, Weight: u.Weight})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error when saving user %v", err)
	}
	return &proto.UserResponse{
		ID:   id,
		Name: u.Name,
	}, nil
}

func (s *DefaultMailBoxService) Login(_ context.Context, u *proto.User) (*proto.UserResponse, error) {
	var user = &auth.User{}
	if err := s.Session.Get(map[string]interface{}{"name": u.Name, "password": u.Password}, u); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return &proto.UserResponse{
		ID:   user.ID,
		Name: u.Name,
	}, nil
}

func NewMailBoxService(watcher imap.Watcher, client smtp.Sender, storage db.Storage, session db.Session) *DefaultMailBoxService {
	return &DefaultMailBoxService{
		Watcher:  watcher,
		Tool:     client,
		Registry: storage,
		Session:  session,
		lock:     sync.Mutex{},
	}
}
