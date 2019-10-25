package main

import (
	"gomail/config"
	"gomail/db"
	"gomail/imap"
	"gomail/smtp"
	"log"
)

func main() {
	mailConfig := config.Load("./config.yml")
	mongo, err := db.New(mailConfig.Mongo)
	if err != nil {
		log.Fatal(err)
	}
	go smtp.Start(mailConfig.Smtp, mongo)
	imap.StartAndListen(mailConfig.Imap)
}

//func init() {
//	message.CharsetReader = func(charset string, input io.Reader) (reader io.Reader, e error) {
//		if strings.ToLower(charset) == "gb2312" {
//			charset = "GB18030"
//		}
//		decoder := mahonia.NewDecoder(charset)
//		if decoder != nil {
//			reader = decoder.NewReader(input)
//		} else {
//			reader = input
//		}
//		return
//	}
//}
//func main() {
//	log.Println("Connecting to server...")
//
//	// Connect to server
//	c, err := client.DialTLS("imap.qq.com:993", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Println("Connected")
//
//	// Don't forget to logout
//	defer func() { _ = c.Logout() }()
//
//	// Login
//	if err := c.Login("1262193323@qq.com", "kwjklcboqznsbabc"); err != nil {
//		log.Fatal(err)
//	}
//	log.Println("Logged in")
//	// Select INBOX
//	mbox, err := c.Select("INBOX", false)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("Flags for INBOX: %+v", mbox)
//
//	criteria := imap.NewSearchCriteria()
//	criteria.WithoutFlags = []string{imap.AnsweredFlag}
//	ids, err := c.Search(criteria)
//	var section imap.BodySectionName
//	seqset := new(imap.SeqSet)
//	seqset.AddNum(ids...)
//	messages := make(chan *imap.Message, 10)
//	done := make(chan error, 1)
//	go func() {
//		done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
//	}()
//
//	log.Println("Last unseen messages:")
//	for msg := range messages {
//		mr, _ := mail.CreateReader(msg.GetBody(&section))
//		header := mr.Header
//
//		if date, err := header.Date(); err == nil {
//			log.Println("Date:", date)
//		}
//		if from, err := header.AddressList("From"); err == nil {
//			log.Println("From:", from[0].Name, from[0])
//		}
//		if to, err := header.AddressList("To"); err == nil {
//			log.Println("To:", to)
//		}
//		if subject, err := header.Subject(); err == nil {
//			log.Println("Subject:", subject)
//		}
//
//		// Process each message's part
//		//for {
//		//	p, err := mr.NextPart()
//		//	if err == io.EOF {
//		//		break
//		//	} else if err != nil {
//		//		log.Fatal(err)
//		//	}
//		//
//		//	switch h := p.Header.(type) {
//		//	case *mail.InlineHeader:
//		//		// This is the message's text (can be plain-text or HTML)
//		//		b, _ := ioutil.ReadAll(p.Body)
//		//		t, _, _ := h.ContentType()
//		//		log.Println("Got text: ", string(b), t)
//		//	case *mail.AttachmentHeader:
//		//		// This is an attachment
//		//		filename, _ := h.Filename()
//		//		log.Println("Got attachment: ", filename)
//		//	}
//		//}
//	}
//
//	if err := <-done; err != nil {
//		fmt.Println(err)
//	}
//
//	log.Println("Done!")
//}
