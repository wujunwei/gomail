package imap

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"gomail/config"
	"sync"
	"time"
)

type Client struct {
	flushTime     time.Duration
	subscriberMax int
	lock          sync.RWMutex
	subscribers   []chan []byte
	User          string
	Password      string
	Done          chan error
	mailBox       *client.Client
}

func (client *Client) Fetch() chan *imap.Message {
	seqSet := &imap.SeqSet{}
	ch := make(chan *imap.Message, 10)
	go func() {
		client.Done <- client.mailBox.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope}, ch)
	}()

	return ch
}

func (client *Client) addSubscriber(subscriber chan []byte) bool {
	client.lock.Lock()
	if len(client.subscribers) >= client.subscriberMax {
		return false
	}
	client.subscribers = append(client.subscribers, subscriber)
	client.lock.Unlock()
	return true
}

func (client *Client) Login() (err error) {
	err = client.mailBox.Login(client.User, client.Password)
	return
}

func New(imapConfig config.Account) (instance *Client, err error) {
	imapClient, err := client.Dial(imapConfig.RemoteServer)

	if err != nil {
		return
	}
	imapClient.Timeout = imapConfig.Timeout * time.Second
	instance = &Client{
		flushTime:     imapConfig.FlushTime,
		subscriberMax: 50,
		lock:          sync.RWMutex{},
		mailBox:       imapClient,
		User:          imapConfig.Auth.User,
		Password:      imapConfig.Auth.Password,
		Done:          make(chan error, 1),
		subscribers:   make([]chan []byte, 50),
	}
	instance.Login()
	return
}

//func main() {
//
//	var err error
//	log.Println("Connecting to server...")
//	c, err = client.DialTLS(server, nil)
//	//连接失败报错
//	if err != nil {
//		log.Fatal(err)
//	}
//	//登陆
//	if err := c.Login(username, password); err != nil {
//		log.Fatal(err)
//	}
//	log.Println("Logged in")
//
//	seqset := &imap.SeqSet{}
//
//	messages := make(chan *imap.Message, 10)
//	done := make(chan error, 1)
//	go func() {
//		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
//	}()
//
//	log.Println("Last 4 messages:")
//	for msg := range messages {
//		fmt.Printf("%+v\n", msg.Envelope.MessageId)
//	}
//
//	if err := <-done; err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println("Done!")
//}
