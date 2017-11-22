package test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/thalmic/gconf/lib"
	"testing"
)

func TestJSONFileLoad(t *testing.T) {

	Convey("Returns an error when the file can't be found", t, func() {
		result, err := lib.NewJSONFileLoader("").Load()
		So(result, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})

	Convey("Loads a JSON file", t, func() {
		result, err := lib.NewJSONFileLoader("test.json").Load()
		So(err, ShouldBeNil)
		So(result, ShouldResemble, map[string]interface{}{
			"string":  "woohoo",
			"boolean": true,
			"integer": 10.0,
			"float":   3.5,
			"array":   []interface{}{"woohoo", true, float64(10), 3.5},
			"object": map[string]interface{}{
				"string":  "woohoo",
				"boolean": true,
				"integer": float64(10),
				"float":   3.5,
			},
		})
	})
}

func TestParseJSON(t *testing.T) {
	loader := lib.NewJSONFileLoader("")

	Convey("Returns an error when parsing invalid JSON", t, func() {
		result, err := loader.ParseJSON([]byte(""))
		So(result, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})

	Convey("Parses valid JSON into a map", t, func() {

		Convey("Parses JSON without sub objects", func() {
			result, err := loader.ParseJSON([]byte(`{"a": "b"}`))
			So(result, ShouldResemble, map[string]interface{}{"a": "b"})
			So(err, ShouldBeNil)
		})

		Convey("Parses JSON with sub objects", func() {
			result, err := loader.ParseJSON([]byte(`{"a": { "b": "c" } }`))
			So(result, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": "c"}})
			So(err, ShouldBeNil)
		})
	})
}
