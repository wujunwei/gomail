package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	. "gomail/config"
	"gomail/server/util"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
	"time"
)

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
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	WithFile    bool   `json:"with_file"`
}

type Client interface {
	Send(task MailTask) (ok bool, err error)
	BuildStruct(task MailTask) *MailClient
	writeHeader(headers []string) *MailClient
}

type MailClient struct {
	HostName string
	Auth     smtp.Auth
	Addr     string
}

func (mClient MailClient) generatorMessageId() string {
	randomByte, _ := util.Alpha(uint64(32))
	hash := sha256.New()
	hash.Write(randomByte)
	randomStr := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	randomStr = strings.ReplaceAll(randomStr, "=", "")
	randomStr = strings.ReplaceAll(randomStr, "/", "")
	randomStr = strings.ReplaceAll(randomStr, "+", "")
	return fmt.Sprintf("<%s@%s>", randomStr, mClient.HostName)
}

func (mClient MailClient) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}
func (mClient MailClient) writeFile(buffer *bytes.Buffer, fileName string) {
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
func (mClient MailClient) BuildStruct(task MailTask) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = task.From
	Header["To"] = strings.Join(task.To, ";")
	Header["Cc"] = strings.Join(task.Cc, ";")
	Header["Bcc"] = strings.Join(task.Bcc, ";")
	Header["Subject"] = task.Subject
	if task.MessageId == "" {
		task.MessageId = mClient.generatorMessageId()
	}
	Header["Message-Id"] = task.MessageId
	Header["In-Reply-To"] = task.ReplyId
	Header["References"] = task.ReplyId
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	mClient.writeHeader(buffer, Header)
	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type:" + task.ContentType + "\r\n"
	body += "\r\n" + task.Body + "\r\n"
	buffer.WriteString(body)

	if task.Attachment.WithFile {
		attachment := "\r\n--" + boundary + "\r\n"
		attachment += "Content-Transfer-Encoding:base64\r\n"
		attachment += "Content-Disposition:attachment\r\n"
		attachment += "Content-Type:" + task.Attachment.ContentType + ";name=\"" + task.Attachment.Name + "\"\r\n"
		buffer.WriteString(attachment)
		defer func() {
			if err := recover(); err != nil {
				log.Fatalln(err)
			}
		}()
		mClient.writeFile(buffer, task.Attachment.Name)
	}

	buffer.WriteString("\r\n--" + boundary + "--")
	return buffer
}

func (mClient MailClient) Send(task MailTask) (messageId string, err error) {
	messageId = mClient.generatorMessageId()
	task.MessageId = messageId
	buffer := mClient.BuildStruct(task)
	err = smtp.SendMail(mClient.Addr, mClient.Auth, task.From, task.To, buffer.Bytes())
	return
}

func NewClient() (MailSender MailClient, err error) {
	//auth
	MailSender.HostName = MailConfig.Host
	MailSender.Addr = MailConfig.Mail.Smtp
	MailSender.Auth = smtp.PlainAuth("", MailConfig.Mail.User, MailConfig.Mail.Password, strings.Split(MailConfig.Mail.Smtp, ":")[0])
	return
}
