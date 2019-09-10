package db

import (
	"gomail/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"mime/multipart"
)

type Client struct {
	DB         *mgo.Database
	gridPrefix string
}

func (client *Client) Upload(filename string, contentType string, stream multipart.File) (id string, err error) {
	defer func() { _ = stream.Close() }()
	gridFS := client.DB.GridFS(client.gridPrefix)
	file, err := gridFS.Create(filename)
	defer func() { _ = file.Close() }()
	if err != nil {
		log.Println(err)
		return
	}
	by, err := ioutil.ReadAll(stream)
	file.SetContentType(contentType)
	_, err = file.Write(by)
	id = file.Id().(bson.ObjectId).Hex()
	return
}

func (client *Client) Download(id bson.ObjectId) (file *mgo.GridFile, err error) {
	if id.Valid() {

	}
	gridFS := client.DB.GridFS(client.gridPrefix)
	file, err = gridFS.OpenId(id)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (client *Client) Close() {
	client.DB.Session.Close()
}

func New() (client *Client, err error) {
	mongoConfig := config.MailConfig.Mongo
	session, err := mgo.Dial(mongoConfig.Url)
	if err != nil {
		return
	}
	client = &Client{DB: session.DB(mongoConfig.Db), gridPrefix: mongoConfig.GridPrefix}
	return
}
