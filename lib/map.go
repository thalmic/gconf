package lib

// MapLoader defines a loader that loads configurations from a map
type MapLoader struct {
	Map map[string]interface{}
}

// NewMapLoader creates a new map loader
func NewMapLoader(stringMap map[string]interface{}) *MapLoader {
	return &MapLoader{
		Map: stringMap,
	}
}

// Load returns the underlying map
func (loader *MapLoader) Load() (map[string]interface{}, error) {
	return loader.Map, nil
}
