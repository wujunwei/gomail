package main

import (
	"fmt"
	"gomail/pkg/config"
	"gomail/pkg/db"
	"gomail/pkg/imap"
	"gomail/pkg/mailbox"
	"gomail/pkg/proto"
	"gomail/pkg/smtp"
	"log"
	"net"
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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", mailConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		err := s.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	log.Println("server start !")
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigs:
		s.GracefulStop()
		postman.Close()
	}
}
