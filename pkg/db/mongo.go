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
	collection string
}
type WrapObject struct {
	Id  interface{} "_id"
	Obj interface{} "Obj"
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

// Set obj should include _id field
func (client *Client) Set(obj interface{}) (string, error) {
	err := client.DB.C(client.collection).Insert(obj)
	return "", err
}

func (client *Client) Get(conditions map[string]interface{}, result interface{}) error {
	return client.DB.C(client.collection).Find(bson.M(conditions)).One(result)
}

func (client *Client) Exist(condition map[string]interface{}) bool {
	n, err := client.DB.C(client.collection).Find(bson.M(condition)).Count()
	if err != nil {
		return false
	}
	return n > 0
}

func (client *Client) Close() {
	client.DB.Session.Close()
}

func New(mongoConfig config.Mongo) (Storage, error) {
	session, err := mgo.Dial(mongoConfig.Url)
	if err != nil {
		return nil, err
	}
	db := session.DB(mongoConfig.Db)
	if mongoConfig.User != "" {
		err = db.Login(mongoConfig.User, mongoConfig.Password)
		if err != nil {
			return nil, err
		}
	}
	client := &Client{DB: db, gridPrefix: mongoConfig.GridPrefix, collection: mongoConfig.Collection}
	return client, nil
}
