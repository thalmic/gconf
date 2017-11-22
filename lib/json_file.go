package lib

import (
	"encoding/json"
	"io/ioutil"
)

type JSONFileLoader struct {
	FilePath string
}

func NewJSONFileLoader(filePath string) *JSONFileLoader {
	return &JSONFileLoader{
		FilePath: filePath,
	}
}

func (loader *JSONFileLoader) Load() (map[string]interface{}, error) {
	file, err := ioutil.ReadFile(loader.FilePath)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return loader.ParseJSON(file)
}

func (loader *JSONFileLoader) ParseJSON(bytes []byte) (map[string]interface{}, error) {
	config := map[string]interface{}{}
	err := json.Unmarshal(bytes, &config)
	return config, err
}
