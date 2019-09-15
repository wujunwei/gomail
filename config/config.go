package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

func Load(path string) (config Config) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	return
}
