package server

import (
	"fmt"
	"gomail/server/db"
	"log"
	"net/http"
)

func Start(addr string) {
	client, err := NewClient()
	mongo, err := db.New()
	if err == nil {
		http.Handle("/mail/send", &MailHandle{Client: &client, Db: mongo})
		http.Handle("/attachment/upload", &FileHandle{db: mongo})
		fmt.Print("start to listen :8080")
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}
