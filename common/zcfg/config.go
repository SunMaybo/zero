package zcfg

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func LoadConfig(filePath string, dst interface{}) {
	buff, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(buff, dst); err != nil {
		panic(err)
	}
	return
}
