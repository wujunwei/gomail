package smtp

import (
	"fmt"
	"gomail/pkg/config"
	"gomail/pkg/db"
	"net"
	"net/http"
)

func Start(smtp config.Smtp, mongo *db.Client) {
	client := NewClient(smtp)
	http.Handle("/mail/send", &MailHandle{Client: client, Db: mongo})
	http.Handle("/attachment/upload", &FileHandle{db: mongo})
	fmt.Print("start to listen :8080")
	panic(http.ListenAndServe(net.JoinHostPort(smtp.Host, smtp.Port), nil))
}
