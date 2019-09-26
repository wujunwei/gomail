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
	RemoteServer  string
	lock          sync.RWMutex
	subscribers   []chan []byte
	User          string
	Password      string
	Done          chan error
	mailBox       *client.Client
}

func (cli *Client) Fetch() chan *imap.Message {
	seqSet := &imap.SeqSet{}
	ch := make(chan *imap.Message, 10)
	go func() {
		cli.Done <- cli.mailBox.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope}, ch)
	}()

	return ch
}

func (cli *Client) addSubscriber(subscriber chan []byte) bool {
	cli.lock.Lock()
	if len(cli.subscribers) >= cli.subscriberMax {
		return false
	}
	cli.subscribers = append(cli.subscribers, subscriber)
	cli.lock.Unlock()
	return true
}

func (cli *Client) Login() (err error) {
	err = cli.mailBox.Login(cli.User, cli.Password)
	return
}

func (cli *Client) Reconnect() (err error) {
	mailClient, err := client.Dial(cli.RemoteServer)
	if err != nil {
		return
	}
	err = mailClient.Login(cli.User, cli.Password)
	cli.mailBox = mailClient
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
		RemoteServer:  imapConfig.RemoteServer,
		User:          imapConfig.Auth.User,
		Password:      imapConfig.Auth.Password,
		Done:          make(chan error, 1),
		subscribers:   make([]chan []byte, 50),
	}
	err = instance.Login()
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
