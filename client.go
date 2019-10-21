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
	_, _ = conn.Write([]byte("1262193323@qq.com:kwjklcboqznsbabc"))
	rec := make([]byte, 1000)
	for {
		_, err = conn.Read(rec)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(string(rec))
	}

	_ = conn.Close()
}
