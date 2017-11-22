# gconf
A simple hierarchical configuration reader, heavily inspired by [nconf](https://github.com/indexzero/nconf).

## Installation
Installation is possible via dep (recommended):
```
dep ensure -add github.com/thalmic/gconf
```

Or via `go get`:
```
go get github.com/thalmic/gconf
```

## Basic Usage
```go
import "github.com/thalmic/gconf"

// Construct
config := gconf.NewConfig() // Create a brand new set of configs
config := gconf.Instance()  // Or use a built in singleton

// Load some configs. In case of collisions, the first loader wins
config.Use(gconf.Arguments("separator", "prefix"))                      // From command line arguments
config.Use(gconf.Environment(false, "separator", "prefix"))             // From environment variables
config.Use(gconf.JSONFile("some_file.json"))                            // From a JSON file
config.Use(gconf.Map(map[string]interface{}{ "SomeKey": "SomeValue" })) // From an arbitrary map

// Convert to a structure or grab the final underlying map
err := config.ToStructure(&MyAwesomeConfigStructure)
configMap := config.Map

// Get an arbitrary value or a map
val, err := config.Get("something")          // interface{}
val, err := config.GetMap("something")       // string[]interface{}
val, err := config.GetSubConfig("something") // A whole new config object containing the sub-map

// Get standard types
val, err := config.GetString("something")  // string
val, err := config.GetInteger("something") // int
val, err := config.GetBoolean("something") // bool
val, err := config.GetFloat("something")   // float64

// Get slices
val, err := config.GetStringSlice("something")  // []string
val, err := config.GetIntegerSlice("something") // []int
val, err := config.GetBooleanSlice("something") // []bool
val, err := config.GetFloatSlice("something")   // []float64

// Set an arbitrary key in memory to an arbitrary value (useful for testing)
config.Set("key", "value")
```

## Loaders
Four config loaders come with this library. More information about these can be found below.

### Arguments
The arguments loader (`gconf.Arguments()`) has 2 parameters:
* separator: The separator to use (more info on this below).
* prefix: The prefix to use. When specified, this loader will ignore arguments that don't start with the specified prefix.
Only arguments with values are read. Arguments should be supplied using the standard `-` or `--` prefix:
```
go run main.go --test1=1 -test2=2     // Reads in test1 and test2
go run main.go test1=1 -test2 --test3 // Ignores all these flags
go run main.go --PREFIXtest=5         // Reads in test=5 if the prefix is configured to "PREFIX"
```

### Environment
The environment loader (`gconf.Environment()`) has 3 parameters:
* lowerCase: A bool defining if env vars should be lower-cased before reading them in.
* separator: The separator to use (more info on this below).
* prefix: The prefix to use. When specified, this loader will ignore environment variables that don't start with the specified prefix.
```
TEST=1 go run main.go       // Reads in TEST=1 if lowerCase is false, and test=1 if lowerCase is true
PREFIXTEST=1 go run main.go // Same as above as long as prefix is set to "PREFIX", reads in nothing otherwise
```

### JSONFile
The JSON file loader (`gconf.JSONFile`) only has 1 parameter:
* filePath: The file path of the JSON file to use.

### Map
The map loader (`gconf.Map`) only has 1 parameter:
* stringMap: The `map[string]interface{}` to add to the config.
This loader should be used for defaulting values not found in any other loaders.

### Extensions
Adding a new loader is very simple, simply create a structure that extends the following interface:
```go
type Loader interface {
	Load() (map[string]interface{}, error)
}
```
If you find yourself using a loader often, please consider opening a PR for it.

## Nested Configuration
A big advantage of using gconf is support for nested configuration values. For example, let's say you load the following
JSON file:
```json
{
  "object": {
    "value": 1
  }
}
```

The result will be as follows:
```go
map[string]interface{}{
	"object": map[string]interface{}{
		"value": 1
	}
}
```

However, if you take this approach, it can be hard to override nested values from command line arguments or environment
variables. That's where the `separator` parameter comes in. For example, let's say we use a environment loader:
```go
config.Use(gconf.Environment(true, "__", ""))
```

We then run the program as follows:
```
OBJECT__VALUE=1 go run main.go
```

This will result in the following map:
```go
map[string]interface{}{
	"object": map[string]interface{}{
		"value": 1
	}
}
```

There are several options for reading in these nested values:
```go
val := config.Map["object"].(map[string]interface{})["value"].(int) // The standard way to get from a nested map :(
val, err := config.getMap("object")["value"].(int)                  // A little bit simpler, but still not ideal
val, err := config.getSubConfig("object").GetInteger("value")       // No more casts :)
val, err := config.GetInteger("object:value")                       // Simple and intuitive :D
```

## Command Line and Environment Parsing
gconf will parse environment and command line parameters into various primitive types. For example, if you are using both
command line and environment loaders and run your program as follows:
```
ENV_INT=1 ENV_FLOAT=3.3 ENV_BOOL=true go run main.go -argInt=1 -argFloat=3.3 -argBool=true
```

gconf will parse the configuration into a map as follows:
```go
map[string]interface{}{
	"ENV_INT": int(1),
	"ENV_FLOAT": float64(3.3),
	"ENV_BOOL": bool(true),
	"argInt": 1,
	"argFloat": float64(3.3),
	"argBool": bool(true),
}
```

gconf also supports intelligent slice and object parsing in argument and environment variables. For example, if you are
using both the environment and argument loaders and run your program as follows:
```
ENV_SLICE="[1, 2, 3]" ENV_OBJECT="{\"key\":\"value\"}" go run main.go -argSlice="[1, 2, 3]" -argObject="{\"key\":\"value\"}"
```

This will be parsed into a map as follows:
```go
map[string]interface{}{
	"ENV_SLICE": []interface{}{1, 2, 3},
	"ENV_OBJECT": map[string]interface{}{ "key": "value" },
	"argSlice": []interface{}{1, 2, 3},
	"argObject": map[string]interface{}{ "key": "value" },
}
```

## Structure Copying
gconf uses the awesome [mapstructure](https://github.com/mitchellh/mapstructure) library under the hood for copying a 
map to a structure. That means that it supports mapstructure's structure tagging out of the box. You can take a look at 
the mapstructure [godoc](https://godoc.org/github.com/mitchellh/mapstructure#Decode) for more information.
