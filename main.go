package main

import (
	"gomail/config"
	"gomail/server"
	"net"
)

func main() {
	//fmt.Printf("%+v", config.MailConfig)
	server.Start(net.JoinHostPort(config.MailConfig.Host, config.MailConfig.Port))
	//m := make(chan server.MailTask, 10)
	//server.Start("127.0.0.1:80", m)
	//conf, err := config.Load("./config.yml")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("%v", conf)
	//client := smtp.Client{}
	//err := client.Auth(smtp.PlainAuth("", "1262193323@qq.com", "kwjklcboqznsbabc", "smtp.qq.com"))
	//client.Mail("wjw3323@live.com")
	//auth := smtp.PlainAuth("", "1262193323@qq.com", "kwjklcboqznsbabc", "smtp.qq.com")
	//to := []string{"wjw3323@live.com"}
	//nickname := "adam"
	//user := "1262193323@qq.com"
	//subject := "test mail"
	//contentType := "Content-Type: text/plain; charset=UTF-8"
	//body := "This is the email body."
	//msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
	//	"<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	//err := smtp.SendMail("smtp.qq.com:587", auth, user, to, msg)
	//if err != nil {
	//	fmt.Printf("send mail error: %v", err)
	//}
}
