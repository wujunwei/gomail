package imap

import (
	"github.com/axgle/mahonia"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"gomail/pkg/config"
	"io"
	"log"
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

type Client struct {
	flushTime       time.Duration
	subscriberLimit int
	RemoteServer    string
	lock            sync.Mutex
	subscribers     map[*MailConn]chan []byte
	User            string
	Password        string
	Done            chan error
	mailBox         *client.Client
}

func (cli *Client) Fetch() (chan *imap.Message, *imap.SeqSet) {
	if nil != cli.mailBox.Noop() {
		return nil, nil
	}
	seqSet := &imap.SeqSet{}
	ch := make(chan *imap.Message, 100)

	seqids, err := cli.SearchUnseen()
	if err != nil {
		log.Println(cli.User, " fetch unsee error: ", err)
		go func() { cli.Done <- err }()
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

func (cli *Client) addSubscriber(conn *MailConn) bool {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	if len(cli.subscribers) >= cli.subscriberLimit {
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
	if err != nil {
		_, _ = cli.mailBox.Select("INBOX", false)
	}

	return
}

func (cli *Client) Reconnect() (err error) {
	cli.mailBox, err = client.DialTLS(cli.RemoteServer, nil)
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
	for _, sub := range cli.subscribers {
		close(sub)
	}
	cli.lock.Unlock()
}

func New(imapConfig config.MailServer) (instance *Client, err error) {
	imapClient, err := client.DialTLS(imapConfig.RemoteServer, nil)
	if err != nil {
		return
	}
	imapClient.Timeout = imapConfig.Timeout * time.Second
	instance = &Client{
		flushTime:       imapConfig.FlushTime,
		subscriberLimit: 50,
		mailBox:         imapClient,
		RemoteServer:    imapConfig.RemoteServer,
		User:            imapConfig.Auth.User,
		Password:        imapConfig.Auth.Password,
		Done:            make(chan error, 1),
		subscribers:     make(map[*MailConn]chan []byte, 50),
	}
	err = instance.Login()
	_, _ = instance.mailBox.Select("INBOX", false)
	return
}
