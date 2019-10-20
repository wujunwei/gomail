//package main
//
//import (
//	"gomail/config"
//	"gomail/imap"
//)
//
//func main() {
//	mailConfig := config.Load("./config.yml")
//	//mongo, err := db.New(mailConfig.Mongo)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//go smtp.Start(mailConfig.Smtp, mongo)
//	imap.StartAndListen(mailConfig.Imap)
//}

package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
	"io"
	"io/ioutil"
	"log"
)

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap-mail.outlook.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer func() { _ = c.Logout() }()

	// Login
	if err := c.Login("wjw3323@live.com", "126219"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Flags for INBOX: %+v", mbox)

	// Get the last 4 messages
	from := mbox.Messages - 8
	to := mbox.Messages - 8
	//if mbox.Messages > 3 {
	//	// We're using unsigned integers here, only substract if the result is > 0
	//	from = mbox.Messages
	//}
	var section imap.BodySectionName
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
	}()

	log.Println("Last unseen messages:")
	for msg := range messages {
		mr, _ := mail.CreateReader(msg.GetBody(&section))
		header := mr.Header
		message.CharsetReader = func(charset string, input io.Reader) (reader io.Reader, e error) {
			decoder := mahonia.NewDecoder(charset)
			if decoder != nil {
				reader = decoder.NewReader(input)
			} else {
				reader = input
			}
			return
		}
		if date, err := header.Date(); err == nil {
			log.Println("Date:", date)
		}
		if from, err := header.AddressList("From"); err == nil {
			log.Println("From:", from)
		}
		if to, err := header.AddressList("To"); err == nil {
			log.Println("To:", to)
		}
		if subject, err := header.Subject(); err == nil {
			log.Println("Subject:", subject)
		}

		// Process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				t, _, _ := h.ContentType()
				log.Println("Got text: ", string(b), t)
			case *mail.AttachmentHeader:
				// This is an attachment
				filename, _ := h.Filename()
				log.Println("Got attachment: ", filename)
			}
		}
	}

	if err := <-done; err != nil {
		fmt.Println(err)
	}

	log.Println("Done!")
}
