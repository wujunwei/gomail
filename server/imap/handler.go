package imap

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Handler interface {
	Serve(conn net.Conn)
}

type MailHandler struct {
	connMap      map[net.Conn]chan []byte
	readTimeout  time.Duration
	writeTimeout time.Duration
	Lock         sync.RWMutex
}

func (mh *MailHandler) Serve(conn net.Conn) {
	log.Println("one connection comes")
	if !mh.accept(conn) {
		return
	}
	log.Println("accept !")
	msgChan := make(chan []byte, 50)
	mh.connMap[conn] = msgChan
	heartbeatTimer := time.NewTimer(time.Duration(30 * time.Second))
	defer func() {
		heartbeatTimer.Stop()
		fmt.Println("close!")
	}()
out:
	for {
		select {
		case msg := <-msgChan:
			{
				_ = conn.SetWriteDeadline(time.Now().Add(mh.writeTimeout))
				_, err := conn.Write(msg)
				if err != nil {
					log.Println(err)
					break out
				}
			}
		case <-heartbeatTimer.C:
			log.Println("tick")
			_ = conn.SetReadDeadline(time.Now().Add(mh.readTimeout))
			data := make([]byte, 1024)
			_, err := conn.Read(data)
			log.Println(string(data))
			_, err = conn.Write([]byte("pong"))
			if err != nil {
				log.Println(err)
				break out
			}
			heartbeatTimer.Reset(30 * time.Second)
		}
	}

}

func (mh *MailHandler) accept(conn net.Conn) bool {

	return conn != nil
}

func (mh *MailHandler) Close() {
	if mh.connMap != nil {
		for conn := range mh.connMap {
			_ = conn.Close()
		}
		mh.connMap = nil
	}

}
