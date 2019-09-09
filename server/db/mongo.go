package db

import "gopkg.in/mgo.v2"

type Client struct {
	Session *mgo.Session
}

func (client *Client) Upload() {

}

func New(url string) (client *Client, err error) {
	session, err := mgo.Dial(url)
	client = &Client{Session: session}
	return
}
