package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

var MailConfig = Config{}

func init() {
	data, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &MailConfig)
	if err != nil {
		log.Fatal(err)
	}
	return
}
