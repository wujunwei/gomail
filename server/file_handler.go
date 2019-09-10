package server

import (
	"gomail/server/db"
	"gomail/server/response"
	"net/http"
)

type FileHandle struct {
	db *db.Client
}

func (fh *FileHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	_ = request.ParseMultipartForm(10000)
	files, ok := request.MultipartForm.File["upload"]

	if ok {
		result := make([]string, len(files))
		for i, file := range files {
			reader, _ := file.Open()
			objId, _ := fh.db.Upload(file.Filename, file.Header.Get("Content-Type"), reader)
			result[i] = objId
		}
		_, _ = writer.Write(response.Success(result))
	} else {
		_, _ = writer.Write(response.Fail(-1, "can not find 'upload'"))
	}
}
