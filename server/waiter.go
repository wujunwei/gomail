package server

import (
	"fmt"
	"log"
	"net/http"
)

func Start(addr string) {
	client, err := NewClient()
	if err == nil {
		http.Handle("/mail/send", &MailHandle{Client: &client})
		http.Handle("/attachment/upload", &FileHandle{})
		fmt.Print("start to listen :8080")
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}
