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

func (mc *MailConn) Read(b []byte) (n int, err error) {
	err = mc.Conn.SetReadDeadline(time.Now().Add(mc.readTimeout))
	n, err = mc.Conn.Read(b)
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
	msg := make([]byte, 1)
	_, _ = conn.Read(msg)
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
					delete(mh.connMap, conn) //清退连接池
					break out
				}
			}
		}
	}

}
