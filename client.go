package main

import (
	"fmt"
	"mime"
)

func main() {
	fds := mime.WordDecoder{}
	header, _ := fds.Decode("=?utf-8?q?=E5=90=B4_=E4=BF=8A=E4=BC=9F?=")
	fmt.Print(header)
	//conn, err := net.Dial("tcp", "localhost:8080")
	//if err != nil {
	//	fmt.Print(err)
	//	os.Exit(0)
	//}
	//_, _ = conn.Write([]byte("1262193323@qq.com:kwjklcboqznsbabc"))
	//rec := make([]byte, 1000)
	//for {
	//	_, err = conn.Read(rec)
	//	if err != nil {
	//		log.Println(err)
	//		break
	//	}
	//	fmt.Println(string(rec))
	//}
	//
	//_ = conn.Close()
}
