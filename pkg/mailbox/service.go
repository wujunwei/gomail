package mailbox

import (
	"context"
	"github.com/emersion/go-message/mail"
	"github.com/golang/protobuf/ptypes/empty"
	"gomail/pkg/db"
	"gomail/pkg/imap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

type DefaultMailBoxService struct {
	Postman  *imap.Postman
	Registry db.Storage
}

func (s *DefaultMailBoxService) Send(context.Context, *Mail) (*SendMailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (s *DefaultMailBoxService) ListServer(context.Context, *empty.Empty) (*ServerList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListServer not implemented")
}
func (s *DefaultMailBoxService) Upload(us MailBox_UploadServer) error {
	uf, err := us.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "error happen %v", err)
	}
	IDChan := make(chan string)
	errChan := make(chan error)
	pr, pw := io.Pipe()
	go func() {
		id, err := s.Registry.Upload(uf.GetName(), uf.GetContentType(), pr)
		if err != nil {
			errChan <- err
			return
		}
		IDChan <- id
	}()
	//todo use io。pipe to upload async
	for {
		uf, err := us.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return status.Errorf(codes.InvalidArgument, "error happen %v", err)
		}
		_, err = pw.Write(uf.GetContent())
		if err != nil {
			return status.Errorf(codes.Internal, "error happen %v", err)
		}
	}
	select {
	case err := <-errChan:
		return err
	case id := <-IDChan:
		return us.SendAndClose(&UploadResponse{FileID: id})
	}
}
func (s *DefaultMailBoxService) Watch(*Server, MailBox_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (s *DefaultMailBoxService) mustEmbedUnimplementedMailBoxServer() {

}

func (s *DefaultMailBoxService) ParseMsg(mr *mail.Reader) *Mail {
	header := mr.Header
	subject, _ := header.Subject()
	log.Println(subject)
	toAddress, _ := header.AddressList("To")
	fromAddress, _ := header.AddressList("From")
	var attachBody *Body
	var text []*Body
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := io.ReadAll(p.Body)
			t, _, _ := h.ContentType()
			text = append(text, &Body{MainBody: b, ContentType: t})
		case *mail.AttachmentHeader:
			// This is an attachment
			contentType, _, _ := h.ContentType()
			b, _ := io.ReadAll(p.Body)
			attachBody = &Body{ContentType: contentType, MainBody: b}
		}
	}
	msgStruct := &Mail{
		MessageID:  header.Get("Message-Id"),
		Subject:    subject,
		To:         changeAddress2str(toAddress),
		From:       changeAddress2str(fromAddress),
		Text:       text,
		Attachment: attachBody,
	}
	return msgStruct
}

func changeAddress2str(addresses []*mail.Address) (to []*Address) {
	to = make([]*Address, len(addresses))
	for key, address := range addresses {
		to[key] = &Address{
			Name:    address.Name,
			Address: address.Address,
		}
	}
	return
}

func NewMailBoxService(postman *imap.Postman, storage db.Storage) *DefaultMailBoxService {
	return &DefaultMailBoxService{Postman: postman, Registry: storage}
}
