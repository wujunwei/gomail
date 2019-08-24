package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

func Load(path string) (conf Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &conf)
	return 
}


