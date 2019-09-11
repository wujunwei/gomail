package main

import (
	"gomail/config"
	"gomail/server/smtp"
	"net"
)

func main() {
	smtp.Start(net.JoinHostPort(config.MailConfig.Host, config.MailConfig.Port))
}
