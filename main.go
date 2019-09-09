package main

import (
	"gomail/config"
	"gomail/server"
	"net"
)

func main() {
	server.Start(net.JoinHostPort(config.MailConfig.Host, config.MailConfig.Port))
}
