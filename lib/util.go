package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Has checks if the supplied map contains the supplied key
func Has(m map[string]interface{}, key string) bool {
	_, keyExists := m[key]
	return keyExists
}

// Set sets the value of a nested key in the supplied map
func Set(m map[string]interface{}, keys []string, value interface{}) (map[string]interface{}, error) {

	// If we're not adding interface{} more keys, return this map
	if len(keys) == 0 {
		return m, nil
	}

	key := keys[0]
	keyExists := Has(m, key)

	// Last key, just write it in and return
	if len(keys) == 1 {

		// Key that we're trying to set already exists
		if keyExists {
			return m, fmt.Errorf("configuration option '%s' already present", key)
		}

		m[key] = value
		return m, nil
	}

	// Initialize a new map. We'll put this in the parent map if there isn't already a key there
	castValue := map[string]interface{}{}

	// The key already exists but if it's a map we can still go into it
	if keyExists {
		var castSuccessfully bool
		castValue, castSuccessfully = m[key].(map[string]interface{})

		if !castSuccessfully {
			return m, fmt.Errorf("configuration option '%s' already present and not a map", key)
		}
	}

	// Recurse and add the next nested value
	submap, err := Set(castValue, keys[1:], value)
	if err != nil {
		return m, err
	}

	m[key] = submap
	return m, nil
}

// Get gets the value of a nested key in the supplied map
func Get(m map[string]interface{}, keys []string) (interface{}, error) {

	key := keys[0]
	keyExists := Has(m, key)

	if !keyExists {
		return nil, fmt.Errorf("key '%s' was not found", key)
	}

	// If this is the last key, return the value
	if len(keys) == 1 {
		return m[key], nil
	}

	// Not the last key in the chain, make sure the next key is a map
	mapValue, castMapValue := m[key].(map[string]interface{})
	if !castMapValue {
		return nil, fmt.Errorf("key '%s' is not a map that can contain sub keys", key)
	}

	return Get(mapValue, keys[1:])
}

// Merge merges two maps recursively
func Merge(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {

	for key, value := range map2 {

		// If we don't have the key in map 1, just take the whole thing
		if !Has(map1, key) {
			map1[key] = value
			continue
		}

		// We have the key in map 1 and map 2, let's see if it's a map in both so we can merge those
		map1Value, castMap1Value := map1[key].(map[string]interface{})
		map2Value, castMap2Value := map2[key].(map[string]interface{})

		// If we failed to cast one of these to a map then we can't merge them. Just ignore the key
		if !castMap1Value || !castMap2Value {
			continue
		}

		// Both of them are maps, keep merging
		map1[key] = Merge(map1Value, map2Value)
	}

	return map1
}

// ParseString parses a string into a variety of types
func ParseString(value string) interface{} {

	// Check if it's an int
	intValue, err := strconv.ParseInt(value, 10, 0)
	if err == nil {
		return int(intValue) // Cast into an integer (parse int returns a int64)
	}

	// Check if it's a float
	floatValue, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return floatValue
	}

	// Check if it's a bool
	boolValue, err := strconv.ParseBool(value)
	if err == nil {
		return boolValue
	}

	// Convert to bytes so we can do JSON checks
	bytes := []byte(value)

	// Check if it's a JSON object
	var jsonObject map[string]interface{}
	jsonObjectError := json.Unmarshal(bytes, &jsonObject)
	if jsonObjectError == nil {
		return jsonObject
	}

	// Check if it's a JSON array
	var jsonArray []interface{}
	jsonArrayError := json.Unmarshal(bytes, &jsonArray)
	if jsonArrayError == nil {
		return jsonArray
	}

	// Finally, try it as a duration
	return ParseDurationString(value)
}

// ParseDurationString attempts to parse a string into a duration, returning the original value if parsing failed
func ParseDurationString(value string) interface{} {
	durationValue, err := time.ParseDuration(value)
	if err == nil {
		return durationValue
	}
	return value
}

// SplitKey splits the supplied key into an array
func SplitKey(key string) []string {
	return strings.Split(key, ":")
}

// CastMap casts a value to a map
func CastMap(obj interface{}) (map[string]interface{}, error) {
	value, cast := obj.(map[string]interface{})
	if !cast {
		return nil, errors.New("failed to cast value to map")
	}
	return value, nil
}

// CastSlice casts a value to a slice
func CastSlice(obj interface{}) ([]interface{}, error) {
	value, cast := obj.([]interface{})
	if !cast {
		return nil, errors.New("failed to cast value to slice")
	}
	return value, nil
}

// CastStringSlice casts a value to a string slice
func CastStringSlice(obj interface{}) ([]string, error) {
	value, cast := obj.([]string)
	if !cast {
		return nil, errors.New("failed to cast value to string slice")
	}
	return value, nil
}

// CastString casts a value to a string
func CastString(obj interface{}) (string, error) {
	value, cast := obj.(string)
	if !cast {
		return "", errors.New("failed to cast value to string")
	}
	return value, nil
}

// CastIntegerSlice casts a value to an integer slice
func CastIntegerSlice(obj interface{}) ([]int, error) {
	value, cast := obj.([]int)
	if !cast {

		// Try to cast it into a float slice (the default when reading JSON)
		floatSlice, cast := obj.([]float64)
		if !cast {
			return nil, errors.New("failed to cast value to integer slice")
		}

		// Managed to convert to a float slice, convert to a integer slice
		intSlice := make([]int, len(floatSlice))
		for i, v := range floatSlice {
			intSlice[i] = int(v)

			if float64(intSlice[i]) != v {
				return nil, errors.New("failed to cast value to integer slice")
			}
		}

		return intSlice, nil
	}
	return value, nil
}

// CastInteger casts a value to an integer
func CastInteger(obj interface{}) (int, error) {
	value, cast := obj.(int) // Try a regular integer cast
	if !cast {

		// Try to cast it into a float (the default when reading JSON) and then convert to int
		floatValue, cast := obj.(float64)
		if !cast {
			return 0, errors.New("failed to cast value to integer")
		}

		if float64(int(floatValue)) != floatValue {
			return 0, errors.New("failed to cast value to integer")
		}

		return int(floatValue), nil
	}
	return value, nil
}

// CastBooleanSlice casts a value to a boolean slice
func CastBooleanSlice(obj interface{}) ([]bool, error) {
	value, cast := obj.([]bool)
	if !cast {
		return nil, errors.New("failed to cast value to boolean slice")
	}
	return value, nil
}

// CastBoolean casts a value to a boolean
func CastBoolean(obj interface{}) (bool, error) {
	value, cast := obj.(bool)
	if !cast {
		return false, errors.New("failed to cast value to boolean")
	}
	return value, nil
}

// CastFloatSlice casts a value to a float slice
func CastFloatSlice(obj interface{}) ([]float64, error) {
	value, cast := obj.([]float64)
	if !cast {
		return nil, errors.New("failed to cast value to float slice")
	}
	return value, nil
}

// CastFloat casts a value to a float
func CastFloat(obj interface{}) (float64, error) {
	value, cast := obj.(float64)
	if !cast {
		return 0, errors.New("failed to cast value to float")
	}
	return value, nil
}
