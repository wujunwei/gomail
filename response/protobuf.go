package response

import "github.com/emersion/go-message/mail"

//todo protobuf

func ConstructMsg(mr *mail.Reader) *Mail {
	header := mr.Header
	subject, _ := header.Subject()
	toAddress, _ := header.AddressList("To")
	fromAddress, _ := header.AddressList("From")
	msgStruct := &Mail{
		MessageId:  0,
		Subject:    subject,
		To:         changeAddress2str(toAddress),
		From:       changeAddress2str(fromAddress),
		Text:       nil,
		Attachment: nil,
	}

	return msgStruct
}

func changeAddress2str(addresses []*mail.Address) (to []string) {
	to = make([]string, len(addresses))
	for key, address := range addresses {
		to[key] = address.String()
	}
	return
}
