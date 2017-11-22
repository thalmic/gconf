package lib

import "github.com/mitchellh/mapstructure"

type Loader interface {
	Load() (map[string]interface{}, error)
}

type Config struct {
	Map map[string]interface{}
}

func NewConfig() *Config {
	return &Config{
		Map: map[string]interface{}{},
	}
}

func (config *Config) Use(loader Loader) {

	// Load in the config map from this loader
	loadedMap, err := loader.Load()
	if err != nil {
		panic(err)
	}

	// Merge it with our existing values
	Merge(config.Map, loadedMap)
}

func (config *Config) ToStructure(structure interface{}) error {
	return mapstructure.Decode(config.Map, structure)
}

func (config *Config) Get(key string) (interface{}, error) {
	return Get(config.Map, SplitKey(key))
}

func (config *Config) GetSubConfig(key string) (*Config, error) {
	value, err := config.GetMap(key)
	if err != nil {
		return nil, err
	}
	return &Config{
		Map: value,
	}, nil
}

func (config *Config) GetMap(key string) (map[string]interface{}, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastMap(value)
}

func (config *Config) GetSlice(key string) ([]interface{}, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastSlice(value)
}

func (config *Config) GetStringSlice(key string) ([]string, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastStringSlice(value)
}

func (config *Config) GetString(key string) (string, error) {
	value, err := config.Get(key)
	if err != nil {
		return "", err
	}
	return CastString(value)
}

func (config *Config) GetIntegerSlice(key string) ([]int, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastIntegerSlice(value)
}

func (config *Config) GetInteger(key string) (int, error) {
	value, err := config.Get(key)
	if err != nil {
		return 0, err
	}
	return CastInteger(value)
}

func (config *Config) GetBooleanSlice(key string) ([]bool, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastBooleanSlice(value)
}

func (config *Config) GetBoolean(key string) (bool, error) {
	value, err := config.Get(key)
	if err != nil {
		return false, err
	}
	return CastBoolean(value)
}

func (config *Config) GetFloatSlice(key string) ([]float64, error) {
	value, err := config.Get(key)
	if err != nil {
		return nil, err
	}
	return CastFloatSlice(value)
}

func (config *Config) GetFloat(key string) (float64, error) {
	value, err := config.Get(key)
	if err != nil {
		return 0, err
	}
	return CastFloat(value)
}

func (config *Config) Set(key string, value interface{}) error {
	_, err := Set(config.Map, SplitKey(key), value)
	return err
}
