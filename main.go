package main

import (
	"fmt"
	"gomail/config"
)

func main() {
	conf ,err := config.Load("./config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v", conf)
	//client := smtp.Client{}
	//err := client.Auth(smtp.PlainAuth("", "1262193323@qq.com", "kwjklcboqznsbabc", "smtp.qq.com"))
	//if err != nil {
	//	fmt.Print(err)
	//	os.Exit(0)
	//}
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
