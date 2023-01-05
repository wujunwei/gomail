package db

import (
	"io"
	"time"
)

type Storage interface {
	Upload(filename string, contentType string, stream io.ReadCloser) (id string, err error)
	Download(id string) (File, error)
}

type File interface {
	io.ReadSeekCloser
	ContentType() string
	Name() string
	MD5() (md5 string)
	UploadDate() time.Time
}
