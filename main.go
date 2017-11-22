package gconf

import (
	"github.com/thalmic/gconf/lib"
	"sync"
)

var configSingleton *lib.Config
var once sync.Once

func New() *lib.Config {
	return lib.NewConfig()
}

func Instance() *lib.Config {
	once.Do(func() {
		configSingleton = lib.NewConfig()
	})
	return configSingleton
}

func Arguments(separator string, prefix string) *lib.ArgumentLoader {
	return lib.NewArgumentLoader(separator, prefix)
}

func Environment(lowerCase bool, separator string, prefix string) *lib.EnvironmentLoader {
	return lib.NewEnvironmentLoader(lowerCase, separator, prefix)
}

func JSONFile(filePath string) *lib.JSONFileLoader {
	return lib.NewJSONFileLoader(filePath)
}

func Map(stringMap map[string]interface{}) *lib.MapLoader {
	return lib.NewMapLoader(stringMap)
}
