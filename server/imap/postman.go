package imap

import (
	"gomail/config"
	"log"
	"net"
)

func StartAndListen(imap config.Imap) {
	log.Println("start to listen" + imap.Host)
	handler := MailHandler{connMap: make(map[net.Conn]chan []byte), readTimeout: imap.Timeout, writeTimeout: imap.Timeout}
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
		go handler.Serve(conn)
	}
}
