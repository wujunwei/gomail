package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

var MailConfig = Config{}

func init() {
	data, err := ioutil.ReadFile("../config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &MailConfig)
	return
}
