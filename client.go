package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	_, _ = conn.Write([]byte("helo"))
	rec := make([]byte, 10)
	conn.Read(rec)
	print(rec)
}
