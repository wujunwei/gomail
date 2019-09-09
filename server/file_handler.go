package server

import (
	"gomail/server/response"
	"net/http"
)

type FileHandle struct {
	//todo
}

func (fh *FileHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	_, _ = writer.Write(response.Success("ok"))
}
