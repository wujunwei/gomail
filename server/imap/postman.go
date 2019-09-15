package imap

import (
	"gomail/config"
	"net"
)

func StartAndListen(imap config.Imap) (listener net.Listener, err error) {
	listener, err = net.Listen(imap.Network, net.JoinHostPort(imap.Host, imap.Port))
	if err != nil {

	}
	defer func() { _ = listener.Close() }()
	for {
		conn, _ := listener.Accept()
		go func() {
			conn.Close() //todo construct a tcp handle
		}()
	}

	return
}
