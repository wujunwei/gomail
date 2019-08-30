package server

import (
	. "gomail/server/handler"
	"log"
	"net/http"
)

func Start(addr string) {
	http.Handle("/mail/send", &MailHandle{})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
