package imap

import "sync"

//alive check， subscribe restart client
type Postman struct {
	Lock        sync.RWMutex
	subscribers map[string]chan []byte
	mailPool    map[string]Client
}

func (postman *Postman) Subscribe(user string, msgChan chan []byte) {
	postman.subscribers[user] = msgChan
}
