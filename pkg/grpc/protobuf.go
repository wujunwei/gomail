package grpc

import (
	"github.com/emersion/go-message/mail"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
)

func ConstructMsg(mr *mail.Reader) ([]byte, error) {
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
		MessageId:  header.Get("Message-Id"),
		Subject:    subject,
		To:         changeAddress2str(toAddress),
		From:       changeAddress2str(fromAddress),
		Text:       text,
		Attachment: attachBody,
	}
	return proto.Marshal(msgStruct)
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
