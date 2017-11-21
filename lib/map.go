package lib

type MapLoader struct {
	Map map[string]interface{}
}

func NewMapLoader(stringMap map[string]interface{}) *MapLoader {
	return &MapLoader{
		Map: stringMap,
	}
}

func (loader *MapLoader) Load() (map[string]interface{}, error) {
	return loader.Map, nil
}
