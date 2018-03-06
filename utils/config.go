package utils

import (
	"encoding/json"
	"io/ioutil"
)

func LoadConfig(path string, v interface{}) (err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, v)
	return
}

func ParseConfigs(configs map[string]interface{}) (err error) {
	for path, v := range configs {
		if err = LoadConfig(path, v); err != nil {
			return
		}
	}
	return
}
