package server

import (
	"log"
	"net/http"
)

type MailHandle struct {
	mailQueue <-chan MailTask
}

func (mh *MailHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte("hahha"))
	writer.WriteHeader(200)
}

func Start(addr string, mailQueue <-chan MailTask) {
	http.Handle("/mail/send", &MailHandle{mailQueue: mailQueue})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
