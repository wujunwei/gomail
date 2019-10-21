package imap

import (
	"gomail/config"
	"log"
	"net"
	"time"
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
			Conn:         conn,
			msgChan:      make(chan []byte, 50),
			Done:         make(chan error, 1),
			readTimeout:  imap.Timeout * time.Second,
			writeTimeout: imap.Timeout * time.Second,
		}
		go handler.Serve(mConn)
	}
}
