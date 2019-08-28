package server

import (
	. "gomail/config"
	"io"
	"log"
	"net/smtp"
)

type MailTask struct {
	from          string
	to            []string
	cc            []string
	bcc           []string
	subject       string
	LastMessageId string
	body          string
	contentType   string
	attachment    Attachment
}

type Attachment struct {
	name        string
	contentType string
	withFile    bool
}

type Client interface {
	Send(task MailTask) (ok bool, err error)
	BuildStruct(task MailTask) *MailClient
	WriteHeaders(headers []string) *MailClient
}

type MailClient struct {
	smtp.Client
	io.Writer
}

func (mClient *MailClient) BuildStruct(task MailTask) *MailClient {
	//todo 拼装邮件结构
	return mClient
}

func (mClient *MailClient) Send(task MailTask) (ok bool, err error) {
	err = mClient.BuildStruct(task).Mail(task.from)
	if err != nil {
		return
	}
	return
}

var client = MailClient{}

func init() {
	client.Client = smtp.Client{}
	//auth
	err := client.Auth(smtp.PlainAuth("", MailConfig.Mail.User, MailConfig.Mail.Password, MailConfig.Mail.Smtp))
	if err != nil {
		log.Fatal(err)
	}
	client.Writer, _ = client.Data()

}
