package server

import (
	"log"
	"net/http"
)

func Start(addr string) {
	pool, err := NewPool()
	if err == nil {
		defer pool.Close()
		http.Handle("/mail/send", &MailHandle{Pool: &pool})

		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}
