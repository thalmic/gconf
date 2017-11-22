package test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/thalmic/gconf/lib"
	"testing"
)

func TestUse(t *testing.T) {

	Convey("Adds a loaded config to the config map", t, func() {
		config := lib.NewConfig()
		config.Use(lib.NewMapLoader(map[string]interface{}{"one": 1}))
		So(config.Map, ShouldResemble, map[string]interface{}{"one": 1})
	})

	Convey("Merges the new config with previous configs", t, func() {
		config := lib.NewConfig()
		config.Use(lib.NewMapLoader(map[string]interface{}{"one": 1}))
		config.Use(lib.NewMapLoader(map[string]interface{}{"two": 2}))
		So(config.Map, ShouldResemble, map[string]interface{}{"one": 1, "two": 2})
	})

	Convey("Panics if the config failed to load", t, func() {
		config := lib.NewConfig()
		So(func() { config.Use(lib.NewJSONFileLoader("")) }, ShouldPanic)
	})
}

func TestToStructure(t *testing.T) {

	Convey("Converts the config map to a structure", t, func() {
		config := lib.NewConfig()
		config.Use(lib.NewMapLoader(map[string]interface{}{"One": 1}))

		structure := struct{ One int }{}
		err := config.ToStructure(&structure)

		So(structure.One, ShouldEqual, 1)
		So(err, ShouldBeNil)
	})
}

func TestGetters(t *testing.T) {
	config := lib.NewConfig()
	config.Use(lib.NewMapLoader(map[string]interface{}{
		"One": 1,
		"Map": map[string]interface{}{
			"Two": "Hi",
		},
	}))

	Convey("Get", t, func() {

		Convey("Gets a top level value", func() {
			value, err := config.Get("One")
			So(value, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})

		Convey("Gets a nested value from the underlying map", func() {
			value, err := config.Get("Map:Two")
			So(value, ShouldEqual, "Hi")
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the key doesn't exist", func() {
			value, err := config.Get("Non-existent")
			So(value, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("GetSubConfig", t, func() {

		Convey("Constructs a new config object containing the sub-map", func() {
			value, err := config.GetSubConfig("Map")
			So(value.Map, ShouldResemble, map[string]interface{}{"Two": "Hi"})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the sub-config can't be cast to a map", func() {
			result, err := config.GetSubConfig("One")
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}
