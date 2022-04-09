package zcfg

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func LoadConfig[T any](filePath string) T {
	buff, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var obj T
	if err := yaml.Unmarshal(buff, &obj); err != nil {
		panic(err)
	}
	return obj
}
