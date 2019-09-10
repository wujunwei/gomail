package config

import "time"

type Mail struct {
	User     string `yaml:"user"`
	Password string `yaml:"pwd"`
	Smtp     string `yaml:"smtp"`
	Imap     string `yaml:"imap"`
}

type Mongo struct {
	Url        string `yaml:"url"`
	Db         string `yaml:"db"`
	GridPrefix string `yaml:"gridPrefix"`
}

type Config struct {
	Host       string        `yaml:"host"`
	Port       string        `yaml:"port"`
	Name       string        `yaml:"name"`
	Mail       Mail          `yaml:"mail"`
	WorkNumber int           `yaml:"workNumber"`
	Mongo      Mongo         `yaml:"mongo"`
	Timeout    time.Duration `yaml:"timeout"`
}
