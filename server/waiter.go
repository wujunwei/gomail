package server

import (
	"log"
	"net/http"
)

func (mh *MailHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	client, err := NewClient()
	writer.WriteHeader(200)
}

func Start(addr string, mailQueue <-chan MailTask) {
	http.Handle("/mail/send", &MailHandle{mailQueue: mailQueue})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
