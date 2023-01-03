package main

import (
	"gomail/pkg/config"
	"gomail/pkg/db"
	"gomail/pkg/imap"
	"gomail/pkg/smtp"
	"log"
)

func main() {
	mailConfig := config.Load("./config.yml")
	mongo, err := db.New(mailConfig.Mongo)
	if err != nil {
		log.Fatal(err)
	}
	go smtp.Start(mailConfig.Smtp, mongo)
	imap.StartAndListen(mailConfig.Imap)
}
