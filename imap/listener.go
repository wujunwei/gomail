package imap

import (
	"gomail/config"
	"log"
	"net"
	"sync"
)

func StartAndListen(imap config.Imap) {

	handler := MailHandler{postman: NewPostMan(imap.Accounts)}
	handler.PostmanStart() // 开启协程定时获取邮件数据
	listener, err := net.Listen(imap.Network, net.JoinHostPort(imap.Host, imap.Port))
	log.Println("start to listen ：" + imap.Host)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		handler.Close()
		_ = listener.Close()

	}()
	for {
		conn, _ := listener.Accept()
		mConn := &MailConn{
			Lock:         sync.RWMutex{},
			Conn:         conn,
			msgChan:      make(chan []byte, 50),
			Done:         make(chan error),
			readTimeout:  imap.Timeout,
			writeTimeout: imap.Timeout,
		}
		go handler.Serve(mConn)
	}
}
