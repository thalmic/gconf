package lib

import (
	"os"
	"strings"
)

// ArgumentLoader defines a loader that loads configuration from command line arguments
type ArgumentLoader struct {
	LowerCase bool
	Prefix    string
	Separator string
}

// NewArgumentLoader creates a new argument loader
func NewArgumentLoader(separator string, prefix string) *ArgumentLoader {
	return &ArgumentLoader{
		Prefix:    prefix,
		Separator: separator,
	}
}

// Load loads environment variables into a configuration map
func (loader *ArgumentLoader) Load() (map[string]interface{}, error) {
	return loader.ParseArguments(os.Args[1:])
}

// ParseArguments parses command line arguments into valid types
func (loader *ArgumentLoader) ParseArguments(args []string) (map[string]interface{}, error) {
	config := map[string]interface{}{}

	for _, arg := range args {

		// If the argument doesn't start with `-`, ignore it
		if !strings.HasPrefix(arg, "-") {
			continue
		}

		// Split the argument up
		parts := strings.Split(arg, "=")

		// If the option doesn't have exactly 2 parts, ignore it
		if len(parts) != 2 {
			continue
		}

		// Read in the key and value
		key := strings.TrimLeft(parts[0], "-")
		value := parts[1]

		// If we have a prefix and the key doesn't match it, ignore this line
		if len(loader.Prefix) > 0 && !strings.HasPrefix(key, loader.Prefix) {
			continue
		}

		// Trim the prefix off the argument name
		trimmedKey := strings.TrimPrefix(key, loader.Prefix)

		// Separate it on the separator if required
		separatedKey := []string{trimmedKey}
		if len(loader.Separator) > 0 {
			separatedKey = strings.Split(trimmedKey, loader.Separator)
		}

		// Parse the value and add it to the final config map
		_, err := Set(config, separatedKey, ParseString(value))
		if err != nil {
			return config, err
		}
	}

	return config, nil
}
