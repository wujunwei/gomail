package main

import (
	"gomail/pkg/config"
	"gomail/pkg/db"
	"gomail/pkg/imap"
	"gomail/pkg/mailbox"
	"gomail/pkg/proto"
	"gomail/pkg/smtp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mailConfig := config.Load("./config.yml")
	storage, err := db.New(mailConfig.Mongo)
	if err != nil {
		log.Fatal(err)
	}
	s := mailbox.NewGRPCServer()
	smtpClient := smtp.NewClient(mailConfig.Smtp)
	postman := imap.NewPostMan(mailConfig.Imap.MailServers)
	postman.Start()
	mb := mailbox.NewMailBoxService(postman, smtpClient, storage)
	proto.RegisterMailBoxServer(s, mb)
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigs:
		s.GracefulStop()
		postman.Close()
	}
}
