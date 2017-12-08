package lib

import (
	"os"
	"strings"
)

type EnvironmentLoader struct {
	LowerCase bool
	Prefix    string
	Separator string
}

func NewEnvironmentLoader(lowerCase bool, separator string, prefix string) *EnvironmentLoader {
	return &EnvironmentLoader{
		LowerCase: lowerCase,
		Prefix:    prefix,
		Separator: separator,
	}
}

func (loader *EnvironmentLoader) Load() (map[string]interface{}, error) {
	return loader.ParseEnvironment(os.Environ())
}

func (loader *EnvironmentLoader) ParseEnvironment(environmentData []string) (map[string]interface{}, error) {
	config := map[string]interface{}{}

	for _, environmentLine := range environmentData {

		// Split the env entry on =
		keyValue := strings.Split(environmentLine, "=")

		// If there was no equals, ignore this line
		if keyValue == nil || len(keyValue) < 2 {
			continue
		}

		// If we have a configured prefix and the key doesn't match it, ignore this line
		if len(loader.Prefix) > 0 && !strings.HasPrefix(keyValue[0], loader.Prefix) {
			continue
		}

		// Trim the prefix off the key and trim the separator if it's there as a prefix (that would result in an empty key)
		trimmedKey := strings.TrimPrefix(keyValue[0], loader.Prefix)
		trimmedKey = strings.TrimPrefix(trimmedKey, loader.Separator)

		// Ignore keys that are empty after trimming
		if len(trimmedKey) == 0 {
			continue
		}

		// Lowercase the key if that option is enabled
		if loader.LowerCase {
			trimmedKey = strings.ToLower(trimmedKey)
		}

		// Separate it on the separator if required
		separatedKeys := []string{trimmedKey}
		if len(loader.Separator) > 0 {
			separatedKeys = strings.Split(trimmedKey, loader.Separator)
		}

		// Set the nested value in the map
		equalIndex := strings.Index(environmentLine, "=")
		value := environmentLine[(equalIndex + 1):]

		_, err := Set(config, separatedKeys, value)
		if err != nil {
			return config, err
		}
	}

	return config, nil
}
