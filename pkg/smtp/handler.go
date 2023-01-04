package smtp

import (
	"encoding/json"
	"gomail/pkg/db"
	"gomail/pkg/grpc"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
)

type MailHandle struct {
	Client Client
	Db     *db.Client
}

func (mh *MailHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var jsonData []byte
	var task = MailTask{}
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			err := r.(error)
			_, _ = writer.Write(grpc.Fail(1, err.Error()))
		}
	}()
	writer.Header().Add("Content-Type", "application/json")
	jsonData, _ = io.ReadAll(request.Body)
	err := json.Unmarshal(jsonData, &task)
	if err != nil {
		panic(err)
	}
	if task.Attachment.WithFile {
		file, err := mh.Db.Download(bson.ObjectIdHex(task.Attachment.Id))
		if err != nil {
			panic(err)
		}
		task.Attachment.ContentType = file.ContentType()
		task.Attachment.Name = file.Name()
		task.Attachment.Reader = file
	}

	MessageId, err := mh.Client.Send(task)
	if err != nil {
		panic(err)
	}
	_, _ = writer.Write(grpc.Success(MessageId))
}

type FileHandle struct {
	db *db.Client
}

func (fh *FileHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	_ = request.ParseMultipartForm(10240)
	files, ok := request.MultipartForm.File["upload"]

	if ok {
		result := make([]string, len(files))
		for i, file := range files {
			reader, _ := file.Open()
			objId, _ := fh.db.Upload(file.Filename, file.Header.Get("Content-Type"), reader)
			result[i] = objId
		}
		_, _ = writer.Write(grpc.Success(result))
	} else {
		_, _ = writer.Write(grpc.Fail(-1, "can not find 'upload'"))
	}
}
