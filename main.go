package main

import (
	"flag"
	"fmt"
	"gomail/pkg/config"
	"gomail/pkg/db"
	"gomail/pkg/imap"
	"gomail/pkg/mailbox"
	"gomail/pkg/mailbox/auth"
	"gomail/pkg/proto"
	"gomail/pkg/smtp"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "path for config file")
}
func main() {
	flag.Parse()
	mailConfig := config.Load(configFile)
	mongo, err := db.New(mailConfig.Mongo)
	if err != nil {
		log.Fatal(err)
	}
	interceptor := auth.NewAuthInterceptor(mongo, mongo)
	s := mailbox.NewGRPCServer(grpc.StreamInterceptor(interceptor.StreamAuth),
		grpc.UnaryInterceptor(interceptor.UnaryAuth))
	smtpClient := smtp.NewClient(mailConfig.Smtp)
	postman := imap.NewPostMan(mailConfig.Imap.MailServers)
	postman.Start()
	mb := mailbox.NewMailBoxService(postman, smtpClient, mongo, mongo)
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
