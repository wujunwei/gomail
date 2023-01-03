package smtp

import (
	"fmt"
	"gomail/pkg/config"
	"gomail/pkg/db"
	"log"
	"net"
	"net/http"
)

func Start(smtp config.Smtp, mongo *db.Client) {
	client, err := NewClient(smtp)
	if err == nil {
		http.Handle("/mail/send", &MailHandle{Client: &client, Db: mongo})
		http.Handle("/attachment/upload", &FileHandle{db: mongo})
		fmt.Print("start to listen :8080")
		err = http.ListenAndServe(net.JoinHostPort(smtp.Host, smtp.Port), nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}
