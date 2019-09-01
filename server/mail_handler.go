package server

import (
	"log"
	"net/http"
)

type MailHandle struct {
	Pool Pool
}

func (mh *MailHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//todo deal with json
	err := mh.Pool.Get().Client.Send(MailTask{from: "1262193323@qq.com", to: []string{"wjw3323@live.com"}, subject: "test", body: "哈哈哈，我收到了"})
	//fmt.Println("end!")
	if err != nil {
		log.Print(err)
	} else {
		_, _ = writer.Write([]byte("hahaahha"))
		writer.WriteHeader(200)
	}
}
