package test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/thalmic/gconf/lib"
)

func TestJSONFileLoad(t *testing.T) {

	Convey("Returns an error when the file can't be found", t, func() {
		result, err := lib.NewJSONFileLoader("", false).Load()
		So(result, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})

	Convey("Loads a JSON file", t, func() {
		result, err := lib.NewJSONFileLoader("test.json", false).Load()
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
	loader := lib.NewJSONFileLoader("", false)

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

		Convey("Returns the original map when duration parsing is disabled", func() {
			result, err := loader.ParseJSON([]byte(`{"a": "3s"}`))
			So(result, ShouldResemble, map[string]interface{}{"a": "3s"})
			So(err, ShouldBeNil)
		})

		Convey("Returns a modified map when duration parsing is enabled", func() {
			l := lib.NewJSONFileLoader("", true)
			result, err := l.ParseJSON([]byte(`{"a": "3s"}`))
			So(result, ShouldResemble, map[string]interface{}{"a": 3 * time.Second})
			So(err, ShouldBeNil)
		})
	})
}

func TestParseDurationStrings(t *testing.T) {
	loader := lib.NewJSONFileLoader("", true)

	Convey("Parses string durations", t, func() {
		m := loader.ParseDurationStrings(map[string]interface{}{"a": "3s"})
		So(m, ShouldResemble, map[string]interface{}{"a": 3 * time.Second})
	})

	Convey("Recurses into submaps", t, func() {
		m := loader.ParseDurationStrings(map[string]interface{}{"a": map[string]interface{}{"b": "3s"}})
		So(m, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": 3 * time.Second}})
	})
}
