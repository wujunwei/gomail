package smtp

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	. "gomail/pkg/config"
	"gomail/pkg/util"
	"io"
	"log"
	"net/smtp"
	"strings"
	"time"
)

const splitLine = "\r\n"

type MailTask struct {
	MessageId   string
	From        string     `json:"from"`
	To          []string   `json:"to"`
	Cc          []string   `json:"cc"`
	Bcc         []string   `json:"bcc"`
	Subject     string     `json:"subject"`
	ReplyId     string     `json:"reply_id"`
	Body        string     `json:"body"`
	ContentType string     `json:"content_type"`
	Attachment  Attachment `json:"attachment"`
}

type Attachment struct {
	Id          string `json:"id"`
	Name        string
	Reader      io.Reader
	ContentType string
	WithFile    bool `json:"with_file"`
}

type Client interface {
	Send(task MailTask) (string, error)
	BuildStruct(task MailTask) *bytes.Buffer
	WriteHeader(io.StringWriter, map[string]string) error
}

type MailClient struct {
	HostName string
	Auth     smtp.Auth
	Addr     string
}

func (c MailClient) generatorMessageId() string {
	randomByte, _ := util.Alpha(uint64(32))
	hash := sha256.New()
	hash.Write(randomByte)
	randomStr := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	randomStr = strings.ReplaceAll(randomStr, "=", "")
	randomStr = strings.ReplaceAll(randomStr, "/", "")
	randomStr = strings.ReplaceAll(randomStr, "+", "")
	return fmt.Sprintf("<%s@%s>", randomStr, c.HostName)
}

func (c MailClient) WriteHeader(buffer io.StringWriter, Header map[string]string) error {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + splitLine
	}
	header += splitLine
	_, err := buffer.WriteString(header)
	return err
}
func (c MailClient) writeFile(buffer *bytes.Buffer, fileName io.Reader) {
	file, err := io.ReadAll(fileName)
	if err != nil {
		panic(err.Error())
	}
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString(splitLine)
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString(splitLine)
		}
	}
}
func (c MailClient) BuildStruct(task MailTask) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = task.From
	Header["To"] = strings.Join(task.To, ";")
	Header["Cc"] = strings.Join(task.Cc, ";")
	Header["Bcc"] = strings.Join(task.Bcc, ";")
	Header["Subject"] = task.Subject
	if task.MessageId == "" {
		task.MessageId = c.generatorMessageId()
	}
	Header["Message-Id"] = task.MessageId
	Header["In-Reply-To"] = task.ReplyId
	Header["References"] = task.ReplyId
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	_ = c.WriteHeader(buffer, Header)
	body := splitLine + "--" + boundary + splitLine
	body += "Content-Type:" + task.ContentType + splitLine
	body += splitLine + task.Body + splitLine
	buffer.WriteString(body)

	if task.Attachment.WithFile {
		attachment := splitLine + "--" + boundary + splitLine
		attachment += "Content-Transfer-Encoding:base64" + splitLine
		attachment += "Content-Disposition:attachment" + splitLine
		attachment += "Content-Type:" + task.Attachment.ContentType + ";name=\"" + task.Attachment.Name + "\"" + splitLine
		buffer.WriteString(attachment)
		defer func() {
			if err := recover(); err != nil {
				log.Fatalln(err)
			}
		}()
		c.writeFile(buffer, task.Attachment.Reader)
	}

	buffer.WriteString(splitLine + "--" + boundary + "--")
	return buffer
}

func (c MailClient) Send(task MailTask) (messageId string, err error) {
	if task.From == "" {
		err = errors.New("unknown json string")
		return
	}
	messageId = c.generatorMessageId()
	task.MessageId = messageId
	buffer := c.BuildStruct(task)
	err = smtp.SendMail(c.Addr, c.Auth, task.From, task.To, buffer.Bytes())
	return
}

func NewClient(smtpConfig Smtp) Client {
	//auth
	MailSender := MailClient{HostName: smtpConfig.Host, Addr: smtpConfig.RemoteServer}
	MailSender.Auth = smtp.PlainAuth("", smtpConfig.User, smtpConfig.Password, strings.Split(smtpConfig.RemoteServer, ":")[0])
	return MailSender
}
