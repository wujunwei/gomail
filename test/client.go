package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	proto2 "gomail/pkg/proto"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	_, _ = conn.Write([]byte(""))
	rec := make([]byte, 100000)
	for {
		n, err := conn.Read(rec)
		if err != nil {
			log.Println(err, n)
			break
		}
		mail := &proto2.Mail{}
		rec = rec[:n]
		err = proto.Unmarshal(rec, mail)
		if err != nil {
			fmt.Println(err, n)
			continue
		}

		fmt.Printf(" get mail: %+v \n", mail)
	}

	_ = conn.Close()
}
