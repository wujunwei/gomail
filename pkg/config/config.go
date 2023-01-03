package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func Load(path string) (config Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	return
}
