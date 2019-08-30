package server

import (
	"bytes"
	"encoding/base64"
	. "gomail/config"
	"io"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
	"time"
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
	writeHeader(headers []string) *MailClient
}

type MailClient struct {
	smtp.Client
	io.Writer
}

func (mClient *MailClient) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}
func (mClient *MailClient) writeFile(buffer *bytes.Buffer, fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}
func (mClient *MailClient) BuildStruct(task MailTask) *MailClient {
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = task.from
	Header["To"] = strings.Join(task.to, ";")
	Header["Cc"] = strings.Join(task.cc, ";")
	Header["Bcc"] = strings.Join(task.bcc, ";")
	Header["Subject"] = task.subject
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	mClient.writeHeader(buffer, Header)
	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type:" + task.contentType + "\r\n"
	body += "\r\n" + task.body + "\r\n"
	buffer.WriteString(body)

	if task.attachment.withFile {
		attachment := "\r\n--" + boundary + "\r\n"
		attachment += "Content-Transfer-Encoding:base64\r\n"
		attachment += "Content-Disposition:attachment\r\n"
		attachment += "Content-Type:" + task.attachment.contentType + ";name=\"" + task.attachment.name + "\"\r\n"
		buffer.WriteString(attachment)
		defer func() {
			if err := recover(); err != nil {
				log.Fatalln(err)
			}
		}()
		mClient.writeFile(buffer, task.attachment.name)
	}

	buffer.WriteString("\r\n--" + boundary + "--")
	return mClient
}

func (mClient *MailClient) Send(task MailTask) (err error) {
	err = mClient.BuildStruct(task).Mail(task.from)
	return
}

func NewClient() (MailSender MailClient, err error) {
	MailSender.Client = smtp.Client{}
	//auth
	err = MailSender.Auth(smtp.PlainAuth("", MailConfig.Mail.User, MailConfig.Mail.Password, MailConfig.Mail.Smtp))
	if err != nil {
		return
	}
	MailSender.Writer, err = MailSender.Data()
	return
}
