package main

import (
	"gomail/pkg/config"
	"gomail/pkg/db"
	"gomail/pkg/imap"
	"gomail/pkg/mailbox"
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
	postman := imap.NewPostMan(mailConfig.Imap.MailServers)
	postman.StartToFetch()
	mb := mailbox.NewMailBoxService(postman, storage)
	mailbox.RegisterMailBoxServer(s, mb)
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigs:
		s.GracefulStop()
		postman.Close()
	}
}
