package imap

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Handler interface {
	Serve(conn net.Conn)
	Close()
}

type MailConn struct {
	Conn         net.Conn
	Done         chan error
	msgChan      chan []byte
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (mc *MailConn) accept() bool {
	return mc.Conn != nil
}
func (mc *MailConn) Close() {
	_ = mc.Conn.Close()
	close(mc.msgChan)
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
	postman *Postman
}

func (mh *MailHandler) PostmanStart() {
	mh.postman.StartToFetch()
}

func (mh *MailHandler) CleanUpClient(user string, conn *MailConn) {
	mh.postman.UnSubscribe(user, conn)
	conn.Close()
}

func (mh *MailHandler) Close() {
	mh.postman.Close()
}

func (mh *MailHandler) Serve(conn *MailConn) {
	log.Println("one connection comes")
	defer func() {
		fmt.Println("close!")
	}()
	if !conn.accept() {
		return
	}
	log.Println("accept !")
	msg := make([]byte, 1024)
	_, err := conn.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	account := strings.Split(string(msg), ":")
	if len(account) < 2 {
		conn.Close()
		return
	}
	conn.Done <- mh.postman.Subscribe(account[0], account[1], conn)
out:
	for {
		select {
		case msg := <-conn.msgChan:
			{
				go func() { conn.Done <- conn.Write(msg) }()
			}
		case err := <-conn.Done:
			{
				if err != nil {
					log.Println(err, "client clean up !")
					mh.CleanUpClient(account[0], conn) //清退连接
					break out
				}
			}
		}
	}

}
