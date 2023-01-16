package config

import "time"

type Mongo struct {
	Url        string `yaml:"url"`
	Db         string `yaml:"db"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	GridPrefix string `yaml:"grid_prefix"`
	Collection string `yaml:"collection"`
}

type Smtp struct {
	User     string `yaml:"user"`
	Password string `yaml:"pwd"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type Imap struct {
	MailServers []MailServer  `yaml:"mailServers"`
	Network     string        `yaml:"network"`
	Timeout     time.Duration `yaml:"timeout"`
}

type MailServer struct {
	Host      string        `yaml:"host"`
	Port      string        `yaml:"port"`
	Auth      Auth          `yaml:"auth"`
	Name      string        `yaml:"name"`
	Timeout   time.Duration `yaml:"timeout"`
	FlushTime time.Duration `yaml:"flush_time"`
}

type Auth struct {
	User     string `yaml:"user"`
	Password string `yaml:"pwd"`
}

type Config struct {
	Smtp  Smtp   `yaml:"smtp"`
	Imap  Imap   `yaml:"imap"`
	Name  string `yaml:"name"`
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Mongo Mongo  `yaml:"mongo"`
}
