package redis

import (
	"github.com/garyburd/redigo/redis"
	"gomail/config"
	"gomail/server"
)

type RQueue struct {
	Conn redis.Conn
}

func New(config config.Redis) (client RQueue, err error) {
	//todo add option
	conn, err := redis.Dial(config.Network, config.Address)
	client.Conn = conn
	return
}

func (rq *RQueue) push(task string) bool {
	return true
}

func (rq *RQueue) pop() (mt server.MailTask, err error) {
	return
}
