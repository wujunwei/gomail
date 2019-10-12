package imap

import (
	"fmt"
	"gomail/config"
	"log"
	"net"
	"sync"
	"time"
)

type Handler interface {
	Serve(conn net.Conn)
	Close()
}

type MailConn struct {
	auth         config.Auth
	Conn         net.Conn
	Lock         sync.RWMutex
	Done         chan error
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (mc *MailConn) accept() bool {
	return mc.Conn != nil
}
func (mc *MailConn) Close() {
	_ = mc.Conn.Close()
	close(mc.Done)
}

func (mc *MailConn) Write(bytes []byte) (err error) {
	err = mc.Conn.SetWriteDeadline(time.Now().Add(mc.writeTimeout))
	_, err = mc.Conn.Write(bytes)
	return
}

type MailHandler struct {
	connMap map[MailConn]chan []byte
}

func (mh *MailHandler) Close() {
	if mh.connMap != nil {
		for conn := range mh.connMap {
			conn.Close()
		}
		mh.connMap = nil
	}

}

func (mh *MailHandler) Serve(conn MailConn) {
	log.Println("one connection comes")
	if !conn.accept() {
		return
	}
	log.Println("accept !")
	msgChan := make(chan []byte, 50)
	mh.connMap[conn] = msgChan
	defer func() {
		fmt.Println("close!")
	}()
out:
	for {
		select {
		case msg := <-msgChan:
			{
				conn.Done <- conn.Write(msg)
			}
		case err := <-conn.Done:
			{
				if err != nil {
					log.Println(err)
					break out
				}
			}
		}
	}

}
