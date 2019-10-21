package imap

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"gomail/config"
	"gomail/response"
	"log"
	"sync"
	"time"
)

type Mail struct {
	MessageId   string
	From        string   `json:"from"`
	To          []string `json:"to"`
	Cc          []string `json:"cc"`
	Bcc         []string `json:"bcc"`
	Subject     string   `json:"subject"`
	ReplyId     string   `json:"reply_id"`
	Body        string   `json:"body"`
	ContentType string   `json:"content_type"`
	Attachment  string   `json:"id"`
}

//alive check， subscribe restart client
type Postman struct {
	Lock     sync.RWMutex
	mailPool map[string]*Client
}

func (postman *Postman) Subscribe(user, password string, conn *MailConn) (err error) {
	chooseBox, ok := postman.mailPool[user]
	log.Println(user + " 开始订阅")
	if !ok || password != chooseBox.Password {
		err = errors.New("user is not exist ot password invalid")
		return
	}
	if !chooseBox.addSubscriber(conn) {
		err = errors.New("up to the max subscribe client")
	}
	log.Println(user + " 订阅成功")
	return
}

func (postman *Postman) UnSubscribe(user string, conn *MailConn) {
	chooseBox, ok := postman.mailPool[user]
	if !ok {
		return
	}
	chooseBox.unSubscribe(conn)
	return
}

func (postman *Postman) addClients(accounts []config.Account) {
	for _, account := range accounts {
		client, err := New(account)
		if err != nil {
			log.Println(err)
			continue
		}
		postman.mailPool[account.Auth.User] = client
	}
}

func (postman *Postman) StartToFetch() {
	for _, client := range postman.mailPool {
		go func() {
			ticker := time.Tick(client.flushTime * time.Second)
			for {
				select {
				case <-ticker:
					mailChan := client.Fetch()
					for msg := range mailChan {
						for _, listener := range client.subscribers {
							listener <- postman.openMessage(msg)
						}
					}
				case err := <-client.Done:
					if err != nil {
						log.Println(err)
						err = client.Reconnect()
						if err != nil {
							log.Println("retry :" + err.Error())
						}
					}
				}
			}
		}()
	}
}

func (postman *Postman) openMessage(msg *imap.Message) (res []byte) {
	var section imap.BodySectionName
	mr, _ := mail.CreateReader(msg.GetBody(&section))
	res, err := response.ConstructMsg(mr)
	if err != nil {
		log.Println(err)
	}
	return
}

func (postman *Postman) Close() {
	for _, cli := range postman.mailPool {
		cli.Close()
	}
}

func NewPostMan(accounts []config.Account) (postman *Postman) {
	postman = &Postman{
		Lock:     sync.RWMutex{},
		mailPool: make(map[string]*Client, len(accounts)),
	}
	postman.addClients(accounts)
	return
}
