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
	"net"
	"net/smtp"
	"strings"
	"time"
)

const (
	SplitLine       = "\r\n"
	Boundary        = "GoBoundary"
	BoundarySign    = "--"
	DefaultEncoding = "base64"
)

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
	buf  *bytes.Buffer
	Host string
	Auth smtp.Auth
	Port string
}

func (c *MailTool) generatorMessageId() string {
	randomByte, _ := random.Alpha(uint64(32))
	hash := sha256.New()
	hash.Write(randomByte)
	randomStr := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	randomStr = strings.ReplaceAll(randomStr, "=", "")
	randomStr = strings.ReplaceAll(randomStr, "/", "")
	randomStr = strings.ReplaceAll(randomStr, "+", "")
	return fmt.Sprintf("<%s@%s>", randomStr, c.Host)
}

func (c *MailTool) writeHeader(Header map[string]string) {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + SplitLine
	}
	c.buf.WriteString(header)
	c.WriteSplitLine()
}

func (c *MailTool) writeFile(reader io.Reader) {
	file, err := io.ReadAll(reader)
	if err != nil {
		panic(err.Error())
	}
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	for index, line := 0, len(payload); index < line; index++ {
		c.buf.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			c.buf.WriteString(SplitLine)
		}
	}
}

func (c *MailTool) WriteSplitLine() {
	c.buf.WriteString(SplitLine)
}

func (c *MailTool) WriteBody(body []byte) {
	c.buf.WriteString(SplitLine)
	c.buf.Write(body)
	c.buf.WriteString(SplitLine)
}
func (c *MailTool) buildHeader(task MailTask) map[string]string {
	Header := make(map[string]string)
	Header["From"] = task.From
	Header["To"] = strings.Join(task.To, ";")
	Header["Cc"] = strings.Join(task.Cc, ";")
	Header["Bcc"] = strings.Join(task.Bcc, ";")
	Header["Subject"] = task.Subject
	Header["Message-Id"] = task.MessageId
	Header["In-Reply-To"] = task.ReplyId
	Header["References"] = task.ReplyId
	Header["Content-Type"] = "multipart/mixed;boundary=" + Boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	return Header
}

func (c *MailTool) writeContentType(contentType string) {
	c.buf.WriteString("Content-Type:" + contentType)
}

func (c *MailTool) writeEncoding(encode string) {
	c.buf.WriteString("Content-Transfer-Encoding:" + encode)
}
func (c *MailTool) writeContentDisposition() {
	c.buf.WriteString("Content-Disposition:attachment")
}

func (c *MailTool) writeContentTypeAndName(ty, name string) {
	c.buf.WriteString(fmt.Sprintf("Content-Type:%s;name=\"%s\"", ty, name))

}

func (c *MailTool) writeAttachment(att Attachment) {
	if att.WithFile {
		return
	}
	c.WriteSplitLine()
	c.writeBoundary(false)
	c.WriteSplitLine()
	c.writeEncoding(DefaultEncoding)
	c.WriteSplitLine()
	c.writeContentDisposition()
	c.WriteSplitLine()
	c.writeContentTypeAndName(att.ContentType(), att.Name())
	c.WriteSplitLine()
	c.writeFile(att.File)
	_ = att.Close()
}
func (c *MailTool) writeBoundary(end bool) {
	if end {
		c.buf.WriteString(BoundarySign + Boundary + BoundarySign)
	} else {
		c.buf.WriteString(BoundarySign + Boundary)
	}
}

func (c *MailTool) build(task MailTask) *bytes.Buffer {
	c.writeHeader(c.buildHeader(task))
	c.WriteSplitLine()
	c.writeBoundary(false)
	c.WriteSplitLine()
	c.writeContentType(task.ContentType)
	c.WriteSplitLine()
	c.WriteBody(task.Body)
	c.WriteSplitLine()
	c.writeAttachment(task.Attachment)
	c.WriteSplitLine()
	c.writeBoundary(true)
	return c.buf
}

func (c *MailTool) Send(task MailTask) (messageId string, err error) {
	if task.From == "" {
		err = errors.New("unknown json string")
		return
	}
	messageId = c.generatorMessageId()
	task.MessageId = messageId
	buffer := c.build(task)
	c.reset()
	err = smtp.SendMail(net.JoinHostPort(c.Host, c.Port), c.Auth, task.From, task.To, buffer.Bytes())
	return
}

func (c *MailTool) reset() {
	c.buf.Reset()
}

func NewClient(smtpConfig Smtp) Tool {
	//auth
	MailSender := &MailTool{
		Port: smtpConfig.Port,
		Host: smtpConfig.Host,
		buf:  bytes.NewBuffer(nil),
		Auth: smtp.PlainAuth("", smtpConfig.User, smtpConfig.Password, smtpConfig.Host),
	}
	return MailSender
}
