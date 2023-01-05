package db

import (
	"errors"
	"gomail/pkg/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
)

type Client struct {
	DB         *mgo.Database
	gridPrefix string
}

func (client *Client) Upload(filename string, contentType string, stream io.ReadCloser) (string, error) {
	defer func() { _ = stream.Close() }()
	gridFS := client.DB.GridFS(client.gridPrefix)
	file, err := gridFS.Create(filename)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	file.SetContentType(contentType)
	_, err = io.Copy(file, stream)
	if err != nil {
		return "", err
	}
	id := file.Id().(bson.ObjectId).Hex()
	return id, nil
}

func (client *Client) Download(id string) (File, error) {
	mongoId := bson.ObjectIdHex(id)
	if !mongoId.Valid() {
		return nil, errors.New("invalid file id")
	}
	gridFS := client.DB.GridFS(client.gridPrefix)
	file, err := gridFS.OpenId(mongoId)
	return file, err
}

func (client *Client) Close() {
	client.DB.Session.Close()
}

func New(mongoConfig config.Mongo) (Storage, error) {
	session, err := mgo.Dial(mongoConfig.Url)
	if err != nil {
		return nil, err
	}
	client := &Client{DB: session.DB(mongoConfig.Db), gridPrefix: mongoConfig.GridPrefix}
	return client, nil
}
