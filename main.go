package main

import (
	"gomail/config"
	"gomail/server/db"
	"gomail/server/imap"
	"gomail/server/smtp"
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
