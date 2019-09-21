package imap

import (
	"gomail/config"
	"log"
	"net"
	"sync"
)

func StartAndListen(imap config.Imap) {
	log.Println("start to listen ：" + imap.Host)
	handler := MailHandler{connMap: make(map[MailConn]chan []byte)}
	listener, err := net.Listen(imap.Network, net.JoinHostPort(imap.Host, imap.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		handler.Close()
		_ = listener.Close()

	}()
	for {
		conn, _ := listener.Accept()
		mConn := MailConn{
			Lock:         sync.RWMutex{},
			Conn:         conn,
			Done:         make(chan error),
			readTimeout:  imap.Timeout,
			writeTimeout: imap.Timeout}
		go handler.Serve(mConn)
	}
}
