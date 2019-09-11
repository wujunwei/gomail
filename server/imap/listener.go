package imap

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
	"os"
)

const (
	server   = "imap.qq.com:993"
	username = "1262193323@qq.com"
	password = "kwjklcboqznsbabc"
)

var c *client.Client

func init() {

}

func main() {

	var err error
	log.Println("Connecting to server...")
	c, err = client.DialTLS(server, nil)
	//连接失败报错
	if err != nil {
		log.Fatal(err)
	}
	//登陆
	if err := c.Login(username, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")
	//mailboxes := make(chan *imap.MailboxInfo, 20)
	//go func() {
	//	_ = c.List("", "*", mailboxes)
	//}()
	////列取邮件夹
	//for m := range mailboxes {
	//	mbox, err := c.Select(m.Name, false)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	to := mbox.Items
	//	for key, mail := range to{
	//		fmt.Printf("%s   : %+v \n", key,mail)
	//	}
	//	log.Printf("%s : %d", m.Name, mbox.Messages)
	//}
	//os.Exit(0)

	seqset := &imap.SeqSet{}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Println("Last 4 messages:")
	for msg := range messages {
		fmt.Printf("%+v\n", msg.Envelope.From[0].PersonalName)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}
