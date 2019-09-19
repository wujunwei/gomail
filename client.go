package main

import (
	"fmt"
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
	_, _ = conn.Write([]byte("hello"))
	rec := make([]byte, 10)
	for {
		_, err = conn.Read(rec)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(string(rec))
	}

	_ = conn.Close()
}
