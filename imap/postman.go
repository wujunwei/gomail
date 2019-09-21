package imap

import (
	"errors"
	"gomail/config"
	"log"
	"sync"
	"time"
)

//alive check， subscribe restart client
type Postman struct {
	Lock     sync.RWMutex
	mailPool map[string]*Client
}

func (postman *Postman) Subscribe(user, password string, msgChan chan []byte) (err error) {
	chooseBox, ok := postman.mailPool[user]
	if ok && password != chooseBox.Password {
		err = errors.New("user is not exist ot password invalid")
		return
	}
	if !chooseBox.addSubscriber(msgChan) {
		err = errors.New("up to the max subscribe client")
	}
	return
}

func (postman *Postman) addClient(accounts []config.Account) {
	for _, account := range accounts {
		client, err := New(account)
		if err != nil {
			log.Println(err)
			continue
		}
		postman.mailPool[account.Auth.User] = client
	}
}

func (postman *Postman) StartToFetch() {
	for _, client := range postman.mailPool {
		go func() {
			ticker := time.Tick(client.flushTime * time.Second)
			for {
				select {
				case <-ticker:
					mailChan := client.Fetch()
					for mail := range mailChan {
						//todo deal mail
						for _, listener := range client.subscribers {
							listener <- []byte(mail.Envelope.Subject)
						}
					}
				case err := <-client.Done:
					if err != nil {
						log.Println(err)
						err = client.Login()
						if err != nil {
							log.Println("retry :" + err.Error())
						}
					}
				}
			}
		}()
	}
}

func NewPostMan(accounts []config.Account) (postman *Postman) {
	postman = &Postman{
		Lock:     sync.RWMutex{},
		mailPool: make(map[string]*Client, len(accounts)),
	}
	postman.addClient(accounts)
	return
}
