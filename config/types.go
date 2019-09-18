package config

import "time"

type Mongo struct {
	Url        string `yaml:"url"`
	Db         string `yaml:"db"`
	GridPrefix string `yaml:"grid_prefix"`
}

type Smtp struct {
	RemoteServer string `yaml:"remote_server"`
	User         string `yaml:"user"`
	Password     string `yaml:"pwd"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
}

type Imap struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	Accounts []Account     `yaml:"account"`
	Network  string        `yaml:"network"`
	Timeout  time.Duration `yaml:"timeout"`
}

type Account struct {
	RemoteServer string        `yaml:"remote_server"`
	User         string        `yaml:"user"`
	Password     string        `yaml:"pwd"`
	FlushTime    time.Duration `yaml:"flush_time"`
}

type Config struct {
	Smtp  Smtp   `yaml:"smtp"`
	Imap  Imap   `yaml:"imap"`
	Name  string `yaml:"name"`
	Mongo Mongo  `yaml:"mongo"`
}
