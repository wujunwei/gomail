package main

import (
	"gomail/config"
	"gomail/server/db"
	"gomail/server/smtp"
	"log"
)

func main() {
	mailConfig := config.Load("./config.yml")
	mongo, err := db.New(mailConfig.Mongo)
	if err != nil {
		log.Fatal(err)
	}
	smtp.Start(mailConfig.Smtp, mongo)
}
