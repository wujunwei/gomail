package handler

import (
	"gomail/config"
	"gomail/server"
	"log"
	"net/http"
)

var (
	mailWorker chan server.MailClient
)

func init() {
	mailWorker = make(chan server.MailClient, config.MailConfig.WorkNumber)
	for i := 0; i < 10; i++ {
		client, err := server.NewClient()
		if err != nil {
			log.Fatal(err)
		}
		mailWorker <- client
	}

}

type MailHandle struct {
}

func (mh *MailHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	writer.WriteHeader(200)
}
