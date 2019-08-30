package server

import "time"

type Pool interface {
	Get() *Node
}
type Node struct {
	client   MailClient
	LastTime time.Time
}
type WorkerPool struct {
}
