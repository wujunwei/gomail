package server

import (
	"encoding/json"
	"gomail/server/response"
	"io/ioutil"
	"log"
	"net/http"
)

type MailHandle struct {
	Client *MailClient
}

func (mh *MailHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var jsonData []byte
	var task = MailTask{}
	defer func() {
		r := recover()
		log.Println(r)
		err := r.(error)
		_, _ = writer.Write(response.Fail(1, err.Error()))
	}()
	writer.Header().Add("Content-Type", "application/json")
	jsonData, _ = ioutil.ReadAll(request.Body)
	err := json.Unmarshal(jsonData, &task)
	if err != nil {
		panic(err)
	}
	MessageId, err := mh.Client.Send(task)
	if err != nil {
		panic(err)
	}
	_, _ = writer.Write(response.Success(MessageId))
}
