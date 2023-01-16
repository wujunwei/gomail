package imap

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"gomail/pkg/config"
	"gomail/pkg/proto"
	"io"
	"log"
	"sync"
	"time"
)

// Postman alive check， subscribe restart client
type Postman struct {
	mailPool map[string]*Client
	lock     *sync.Mutex
}

func (postman *Postman) Subscribe(serverName, id string, weight int32, ch chan *proto.Mail) (*Subscriber, error) {
	chooseBox, ok := postman.mailPool[serverName]
	if !ok {
		return nil, errors.New("server is invalid")
	}
	sub := &Subscriber{
		Weight:     weight,
		ID:         id,
		Channel:    ch,
		serverName: serverName,
	}
	if !chooseBox.addSubscriber(sub) {
		return nil, errors.New("up to the max subscribe client")
	}
	log.Println(serverName + " subscribe successfully")
	return sub, nil
}

func (postman *Postman) UnSubscribe(sub *Subscriber) {
	chooseBox, ok := postman.mailPool[sub.serverName]
	if !ok {
		return
	}
	chooseBox.unSubscribe(sub)
	return
}

func (postman *Postman) addClients(accounts []config.MailServer) {
	postman.lock.Lock()
	defer postman.lock.Unlock()
	for _, account := range accounts {
		_, ok := postman.mailPool[account.Name]
		if ok {
			continue
		}
		client, err := New(account)
		if err != nil {
			log.Println(err)
			continue
		}
		postman.mailPool[account.Name] = client
	}
}

func (postman *Postman) Start() {
	for _, cli := range postman.mailPool {
		go func(client *Client) {
			ticker := time.NewTicker(client.flushTime * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					mailChan, seqSet := client.Fetch()
					for msg := range mailChan {
						message, err := postman.openMessage(msg)
						if err != nil {
							log.Printf("open message: %s", err)
							continue
						}
						log.Println("start to push msg , subscribers :", client.subscribers.Size())
						client.subscribers.Each(func(index int, a *Subscriber) {
							log.Println("pushing message !!")
							a.Channel <- message
						})
					}
					if seqSet != nil {
						log.Println("start to see")
						go client.See(seqSet)
						log.Println("saw !")
					}

				case err := <-client.Done: //处理异常需开启协程
					if err != nil {
						log.Println("error happen:", err)
						err = client.Reconnect()
						if err != nil {
							log.Println("retry :" + err.Error())
							return
						} else {
							log.Println("retry success !")
						}
					}
				}
			}
		}(cli)
	}
}

func (postman *Postman) ListServer() []string {
	server := make([]string, len(postman.mailPool))
	i := 0
	for s := range postman.mailPool {
		server[i] = s
		i++
	}
	return server
}

func (postman *Postman) openMessage(msg *imap.Message) (*proto.Mail, error) {
	var section imap.BodySectionName
	mr, err := mail.CreateReader(msg.GetBody(&section))
	if err != nil {
		log.Println("construct message error:", err)
		return nil, err
	}
	email := postman.parseMsg(mr)
	return email, nil
}
func (postman *Postman) parseMsg(mr *mail.Reader) *proto.Mail {
	header := mr.Header
	subject, _ := header.Subject()
	log.Println(subject)
	toAddress, _ := header.AddressList("To")
	fromAddress, _ := header.AddressList("From")
	var attachBody *proto.Body
	var text []*proto.Body
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
			text = append(text, &proto.Body{MainBody: b, ContentType: t})
		case *mail.AttachmentHeader:
			// This is an attachment
			contentType, _, _ := h.ContentType()
			b, _ := io.ReadAll(p.Body)
			attachBody = &proto.Body{ContentType: contentType, MainBody: b}
		}
	}
	msgStruct := &proto.Mail{
		MessageID: header.Get("Message-Id"),
		Subject:   subject,
		To:        changeAddress2str(toAddress),
		//From:       changeAddress2str(fromAddress),
		Text:       text,
		Attachment: attachBody,
	}
	if len(fromAddress) > 0 {
		msgStruct.From = &proto.Address{Name: fromAddress[0].Name, Address: fromAddress[0].Address}
	}
	return msgStruct
}

func changeAddress2str(addresses []*mail.Address) (to []*proto.Address) {
	to = make([]*proto.Address, len(addresses))
	for key, address := range addresses {
		to[key] = &proto.Address{
			Name:    address.Name,
			Address: address.Address,
		}
	}
	return
}

func (postman *Postman) Close() {
	for _, cli := range postman.mailPool {
		cli.Close()
	}
}

func NewPostMan(accounts []config.MailServer) Watcher {
	postman := &Postman{
		mailPool: make(map[string]*Client, len(accounts)),
		lock:     &sync.Mutex{},
	}
	postman.addClients(accounts)
	return postman
}
