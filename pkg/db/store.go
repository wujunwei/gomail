package db

import (
	"io"
	"time"
)

type Storage interface {
	Upload(filename string, contentType string, stream io.ReadCloser) (id string, err error)
	Download(id string) (File, error)
	Set(obj interface{}) (string, error)
	Get(condition map[string]interface{}, result interface{}) error
	Exist(condition map[string]interface{}) bool
}

type File interface {
	io.ReadSeekCloser
	ContentType() string
	Name() string
	MD5() (md5 string)
	UploadDate() time.Time
}
