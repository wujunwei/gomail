package server

import (
	"encoding/json"
	"gomail/server/response"
	"io/ioutil"
	"log"
	"net/http"
)

type MailHandle struct {
	Client MailClient
}

func (mh *MailHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var jsonData []byte
	var task = MailTask{}

	jsonData, _ = ioutil.ReadAll(request.Body)
	err := json.Unmarshal(jsonData, &task)
	if err != nil {
		log.Print(err)
	}
	err = mh.Client.Send(task)
	//fmt.Println("end!")
	if err != nil {
		log.Print(err)
		_, _ = writer.Write(response.Fail(1, "error"))
	} else {
		_, _ = writer.Write(response.Success(nil))
	}
}
