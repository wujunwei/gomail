package imap

import (
	"github.com/axgle/mahonia"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"gomail/pkg/config"
	"gomail/pkg/proto"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func init() {
	message.CharsetReader = func(charset string, input io.Reader) (reader io.Reader, e error) {
		if strings.ToLower(charset) == "gb2312" {
			charset = "GB18030"
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

type Watcher interface {
	Subscribe(serverName, id string, ch chan *proto.Mail) error
	UnSubscribe(serverName, id string)
	Start()
	Close()
	ListServer() []string
}

type Client struct {
	flushTime       time.Duration
	subscriberLimit int
	Host            string
	Port            string
	lock            sync.Mutex
	subscribers     map[string]chan *proto.Mail
	User            string
	Password        string
	Done            chan error
	mailBox         *client.Client
}

func (cli *Client) Fetch() (chan *imap.Message, *imap.SeqSet) {
	if err := cli.mailBox.Noop(); err != nil {
		cli.Done <- err
		return nil, nil
	}
	seqSet := &imap.SeqSet{}
	ch := make(chan *imap.Message, 100)

	seqids, err := cli.SearchUnseen()
	if err != nil {
		log.Println(cli.User, " fetch unsee error: ", err)
		cli.Done <- err
		close(ch)
		return ch, nil
	}
	if len(seqids) == 0 {
		log.Println(cli.User, " 没有邮件")
		close(ch)
		return ch, nil
	}
	seqSet.AddNum(seqids...)

	go func() {
		err := cli.mailBox.Fetch(seqSet, []imap.FetchItem{imap.FetchBody + "[]", imap.FetchFlags, imap.FetchUid}, ch)
		if err != nil {
			cli.Done <- err
		}
	}()

	return ch, seqSet
}

func (cli *Client) SearchUnseen() (ids []uint32, err error) {
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err = cli.mailBox.Search(criteria)
	return
}

func (cli *Client) See(seqSet *imap.SeqSet) {
	cli.Done <- cli.mailBox.Store(seqSet, imap.AddFlags, []interface{}{imap.SeenFlag}, nil)
}

func (cli *Client) addSubscriber(id string, ch chan *proto.Mail) bool {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	if len(cli.subscribers) >= cli.subscriberLimit {
		return false
	}
	cli.subscribers[id] = ch
	return true
}

func (cli *Client) unSubscribe(id string) {
	cli.lock.Lock()
	delete(cli.subscribers, id)
	cli.lock.Unlock()
}

func (cli *Client) Login() (err error) {
	err = cli.mailBox.Login(cli.User, cli.Password)
	if err != nil {
		_, _ = cli.mailBox.Select("INBOX", false)
	}

	return
}

func (cli *Client) Reconnect() (err error) {
	cli.mailBox, err = client.DialTLS(net.JoinHostPort(cli.Host, cli.Port), nil)
	if err != nil {
		return
	}
	err = cli.mailBox.Login(cli.User, cli.Password)
	_, _ = cli.mailBox.Select("INBOX", false)
	return
}

func (cli *Client) Close() {
	cli.lock.Lock()
	_ = cli.mailBox.Close()
	cli.lock.Unlock()
}

func New(imapConfig config.MailServer) (instance *Client, err error) {
	remote := net.JoinHostPort(imapConfig.Host, imapConfig.Port)
	imapClient, err := client.DialTLS(remote, nil)
	if err != nil {
		return
	}
	imapClient.Timeout = imapConfig.Timeout * time.Second
	instance = &Client{
		flushTime:       imapConfig.FlushTime,
		subscriberLimit: 50,
		mailBox:         imapClient,
		Host:            imapConfig.Host,
		Port:            imapConfig.Port,
		User:            imapConfig.Auth.User,
		Password:        imapConfig.Auth.Password,
		Done:            make(chan error, 10),
		subscribers:     make(map[string]chan *proto.Mail, 50),
	}
	err = instance.Login()
	_, _ = instance.mailBox.Select("INBOX", false)
	return
}
