package smtp

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	. "gomail/pkg/config"
	"gomail/pkg/db"
	"gomail/pkg/util/random"
	"io"
	"log"
	"net"
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
	Body        []byte     `json:"body"`
	ContentType string     `json:"content_type"`
	Attachment  Attachment `json:"attachment"`
}

type Attachment struct {
	db.File
	WithFile bool `json:"with_file"`
}

type Tool interface {
	Send(task MailTask) (string, error)
}

type MailTool struct {
	Host string
	Auth smtp.Auth
	Port string
}

func (c MailTool) generatorMessageId() string {
	randomByte, _ := random.Alpha(uint64(32))
	hash := sha256.New()
	hash.Write(randomByte)
	randomStr := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	randomStr = strings.ReplaceAll(randomStr, "=", "")
	randomStr = strings.ReplaceAll(randomStr, "/", "")
	randomStr = strings.ReplaceAll(randomStr, "+", "")
	return fmt.Sprintf("<%s@%s>", randomStr, c.Host)
}

func (c MailTool) writeHeader(buffer io.StringWriter, Header map[string]string) error {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + splitLine
	}
	header += splitLine
	_, err := buffer.WriteString(header)
	return err
}
func (c MailTool) writeFile(buffer *bytes.Buffer, reader io.Reader) {
	file, err := io.ReadAll(reader)
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
func (c MailTool) build(task MailTask) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = task.From
	Header["To"] = strings.Join(task.To, ";")
	Header["Cc"] = strings.Join(task.Cc, ";")
	Header["Bcc"] = strings.Join(task.Bcc, ";")
	Header["Subject"] = task.Subject
	Header["Message-Id"] = task.MessageId
	Header["In-Reply-To"] = task.ReplyId
	Header["References"] = task.ReplyId
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	_ = c.writeHeader(buffer, Header)
	buffer.WriteString(splitLine + "--" + boundary + splitLine)
	buffer.WriteString("Content-Type:" + task.ContentType + splitLine)
	buffer.WriteString(splitLine)
	buffer.Write(task.Body)
	buffer.WriteString(splitLine)

	if task.Attachment.WithFile {
		attachment := splitLine + "--" + boundary + splitLine
		attachment += "Content-Transfer-Encoding:base64" + splitLine
		attachment += "Content-Disposition:attachment" + splitLine
		attachment += "Content-Type:" + task.Attachment.ContentType() + ";name=\"" + task.Attachment.Name() + "\"" + splitLine
		buffer.WriteString(attachment)
		defer func() {
			if err := recover(); err != nil {
				log.Fatalln(err)
			}
		}()
		c.writeFile(buffer, task.Attachment)
	}

	buffer.WriteString(splitLine + "--" + boundary + "--")
	return buffer
}

func (c MailTool) Send(task MailTask) (messageId string, err error) {
	if task.From == "" {
		err = errors.New("unknown json string")
		return
	}
	messageId = c.generatorMessageId()
	task.MessageId = messageId
	buffer := c.build(task)
	err = smtp.SendMail(net.JoinHostPort(c.Host, c.Port), c.Auth, task.From, task.To, buffer.Bytes())
	return
}

func NewClient(smtpConfig Smtp) Tool {
	//auth
	MailSender := MailTool{Port: smtpConfig.Port, Host: smtpConfig.Host}
	MailSender.Auth = smtp.PlainAuth("", smtpConfig.User, smtpConfig.Password, smtpConfig.Host)
	return MailSender
}
