package imap

import (
	"net"
	"sync"
	"time"
)

type Handler interface {
	Serve(conn net.Conn)
}

type MailHandler struct {
	conns   map[net.Conn]bool
	timeout time.Time
	sync.RWMutex
}

func (mh *MailHandler) Serve(conn net.Conn) {
	if !mh.accept(conn) {
		return
	}
	mh.conns[conn] = true
	_ = conn.SetDeadline(mh.timeout)
	//todo subscribe relation

}

func (mh *MailHandler) accept(conn net.Conn) bool {
	return true
}

func (mh *MailHandler) Close() {
	if mh.conns != nil {
		for conn := range mh.conns {
			_ = conn.Close()
		}
		mh.conns = nil
	}

}
