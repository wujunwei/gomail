package imap

import (
	"github.com/axgle/mahonia"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"gomail/config"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

func init() {
	message.CharsetReader = func(charset string, input io.Reader) (reader io.Reader, e error) {
		if strings.ToLower(charset) == "gb2312" {
			charset = "gb2312"
		}
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
	subscribers   map[*MailConn]chan []byte
	User          string
	Password      string
	Done          chan error
	mailBox       *client.Client
}

func (cli *Client) Fetch() chan *imap.Message {
	seqSet := &imap.SeqSet{}
	ch := make(chan *imap.Message, 100)
	status, err := cli.mailBox.Select("INBOX", false)
	if err != nil || status.UnseenSeqNum == 0 {
		log.Println(status.UnseenSeqNum)
		close(ch)
		return ch
	}
	seqSet.AddRange(status.UnseenSeqNum, status.UnseenSeqNum+status.Unseen-1)
	go func() {
		cli.Done <- cli.mailBox.Fetch(seqSet, []imap.FetchItem{imap.FetchBody + "[]"}, ch)
	}()

	return ch
}

func (cli *Client) addSubscriber(conn *MailConn) bool {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	if len(cli.subscribers) >= cli.subscriberMax {
		return false
	}
	cli.subscribers[conn] = conn.msgChan
	return true
}

func (cli *Client) unSubscribe(conn *MailConn) {
	cli.lock.Lock()
	delete(cli.subscribers, conn)
	cli.lock.Unlock()
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

func (cli *Client) Close() {
	cli.lock.Lock()
	_ = cli.mailBox.Close()
	for _, sub := range cli.subscribers {
		close(sub)
	}
	cli.lock.Unlock()
}

func New(imapConfig config.Account) (instance *Client, err error) {
	imapClient, err := client.DialTLS(imapConfig.RemoteServer, nil)
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
		subscribers:   make(map[*MailConn]chan []byte, 50),
	}
	err = instance.Login()
	return
}
