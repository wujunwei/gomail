package grpc

import (
	"github.com/emersion/go-message/mail"
	"github.com/golang/protobuf/proto"
	"gomail/pkg/mailbox"
	"io"
	"log"
)

func ConstructMsg(mr *mail.Reader) ([]byte, error) {
	header := mr.Header
	subject, _ := header.Subject()
	log.Println(subject)
	toAddress, _ := header.AddressList("To")
	fromAddress, _ := header.AddressList("From")
	var attachBody *mailbox.Body
	var text []*mailbox.Body
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
			text = append(text, &mailbox.Body{MainBody: b, ContentType: t})
		case *mail.AttachmentHeader:
			// This is an attachment
			contentType, _, _ := h.ContentType()
			b, _ := io.ReadAll(p.Body)
			attachBody = &mailbox.Body{ContentType: contentType, MainBody: b}
		}
	}
	msgStruct := &mailbox.Mail{
		MessageID:  header.Get("Message-Id"),
		Subject:    subject,
		To:         changeAddress2str(toAddress),
		From:       changeAddress2str(fromAddress),
		Text:       text,
		Attachment: attachBody,
	}
	return proto.Marshal(msgStruct)
}

func changeAddress2str(addresses []*mail.Address) (to []*mailbox.Address) {
	to = make([]*mailbox.Address, len(addresses))
	for key, address := range addresses {
		to[key] = &mailbox.Address{
			Name:    address.Name,
			Address: address.Address,
		}
	}
	return
}
