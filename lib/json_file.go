package lib

import (
	"encoding/json"
	"io/ioutil"
)

// JSONFileLoader defines a loader that loads configurations from a JSON file
type JSONFileLoader struct {
	FilePath       string
	ParseDurations bool
}

// NewJSONFileLoader creates a new JSON file loader
func NewJSONFileLoader(filePath string, parseDurations bool) *JSONFileLoader {
	return &JSONFileLoader{
		FilePath:       filePath,
		ParseDurations: parseDurations,
	}
}

// Load loads a JSON file
func (loader *JSONFileLoader) Load() (map[string]interface{}, error) {
	file, err := ioutil.ReadFile(loader.FilePath)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return loader.ParseJSON(file)
}

// ParseJSON parses json into a configuration map
func (loader *JSONFileLoader) ParseJSON(bytes []byte) (map[string]interface{}, error) {
	config := map[string]interface{}{}
	err := json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	// If we were configured to parse durations, do that
	if loader.ParseDurations {
		return loader.ParseDurationStrings(config), nil
	}

	return config, nil
}

// ParseDurationStrings recursively loops through all the string values in the supplied map and parses them to a duration if possible
func (loader *JSONFileLoader) ParseDurationStrings(m map[string]interface{}) map[string]interface{} {
	for key, value := range m {

		// If the value is a string, apply duration parsing
		stringValue, castStringValue := value.(string)
		if castStringValue {
			m[key] = ParseDurationString(stringValue)
			continue
		}

		// If the value is not a map, keep going
		mapValue, castMapValue := value.(map[string]interface{})
		if !castMapValue {
			continue
		}

		// Value is a map, recurse
		m[key] = loader.ParseDurationStrings(mapValue)
	}

	return m
}
