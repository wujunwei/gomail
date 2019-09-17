package imap

import (
	"gomail/config"
	"log"
	"net"
)

func StartAndListen(imap config.Imap) {
	listener, err := net.Listen(imap.Network, net.JoinHostPort(imap.Host, imap.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	for {
		conn, _ := listener.Accept()
		conn.Close()
	}
}
