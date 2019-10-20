package imap

import (
	"github.com/axgle/mahonia"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"gomail/config"
	"io"
	"sync"
	"time"
)

func init() {
	message.CharsetReader = func(charset string, input io.Reader) (reader io.Reader, e error) {
		decoder := mahonia.NewDecoder(charset)
		if decoder != nil {
			reader = decoder.NewReader(input)
		} else {
			reader = input
		}
		return
	}
}

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
	ch := make(chan *imap.Message, 100)
	status, _ := cli.mailBox.Select("INBOX", true)
	if status.UnseenSeqNum == 0 {
		close(ch)
		return ch
	}
	seqSet.AddRange(status.UnseenSeqNum, status.UnseenSeqNum+status.Unseen-1)
	go func() {
		cli.Done <- cli.mailBox.Fetch(seqSet, []imap.FetchItem{imap.FetchBody + "[]"}, ch)
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
	cli.mailBox, err = client.Dial(cli.RemoteServer)
	if err != nil {
		return
	}
	err = cli.mailBox.Login(cli.User, cli.Password)
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
