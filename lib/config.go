package lib

import "github.com/mitchellh/mapstructure"

// Loader defines a generic loader interface
type Loader interface {
	Load() (map[string]interface{}, error)
}

// Config defines the overall configuration structure
type Config struct {
	Map map[string]interface{}
}

// NewConfig creates a new configuration structure
func NewConfig() *Config {
	return &Config{
		Map: map[string]interface{}{},
	}
}

// Use adds a loader to the configuration loading chain
func (config *Config) Use(loader Loader) {

	// Load in the config map from this loader
	loadedMap, err := loader.Load()
	if err != nil {
		panic(err)
	}

	// Merge it with our existing values
	Merge(config.Map, loadedMap)
}

// ToStructure maps the loaded configuration to a structure
func (config *Config) ToStructure(structure interface{}) error {
	return mapstructure.Decode(config.Map, structure)
}

// Get gets a key from the loaded configuration
func (config *Config) Get(key string) (interface{}, error) {
	return Get(config.Map, SplitKey(key))
}

// GetSubConfig gets a loaded submap as a configuration structure
func (config *Config) GetSubConfig(key string) (*Config, error) {
	value, err := config.GetMap(key)
	if err != nil {
		return nil, err
	}
	return &Config{
		Map: value,
	}, nil
}

// GetMap gets a map from the loaded configuration
func (config *Config) GetMap(key string) (map[string]interface{}, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastMap(value)
}

// GetSlice gets a slice from the loaded configuration
func (config *Config) GetSlice(key string) ([]interface{}, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastSlice(value)
}

// GetStringSlice gets a string slice from the loaded configuration
func (config *Config) GetStringSlice(key string) ([]string, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastStringSlice(value)
}

// GetString gets a string from the loaded configuration
func (config *Config) GetString(key string) (string, error) {
	value, err := config.Get(key)
	if err != nil {
		return "", err
	}
	return CastString(value)
}

// GetIntegerSlice gets a integer slice from the loaded configuration
func (config *Config) GetIntegerSlice(key string) ([]int, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastIntegerSlice(value)
}

// GetInteger gets a integer from the loaded configuration
func (config *Config) GetInteger(key string) (int, error) {
	value, err := config.Get(key)
	if err != nil {
		return 0, err
	}
	return CastInteger(value)
}

// GetBooleanSlice gets a boolean slice from the loaded configuration
func (config *Config) GetBooleanSlice(key string) ([]bool, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastBooleanSlice(value)
}

// GetBoolean gets a boolean from the loaded configuration
func (config *Config) GetBoolean(key string) (bool, error) {
	value, err := config.Get(key)
	if err != nil {
		return false, err
	}
	return CastBoolean(value)
}

// GetFloatSlice gets a float slice from the loaded configuration
func (config *Config) GetFloatSlice(key string) ([]float64, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastFloatSlice(value)
}

// GetFloat gets a float from the loaded configuration
func (config *Config) GetFloat(key string) (float64, error) {
	value, err := config.Get(key)
	if err != nil {
		return 0, err
	}
	return CastFloat(value)
}

// Set sets a value in the loaded configuration
func (config *Config) Set(key string, value interface{}) error {
	_, err := Set(config.Map, SplitKey(key), value)
	return err
}
