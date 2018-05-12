package test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/thalmic/gconf/lib"
)

func TestHas(t *testing.T) {

	Convey("Returns false if the map doesn't have the key", t, func() {
		result := lib.Has(map[string]interface{}{}, "something")
		So(result, ShouldBeFalse)
	})

	Convey("Returns true if the map has the key", t, func() {
		result := lib.Has(map[string]interface{}{"something": "stuff"}, "something")
		So(result, ShouldBeTrue)
	})
}

func TestSet(t *testing.T) {

	Convey("Returns the input map when no keys are specified", t, func() {
		result, err := lib.Set(map[string]interface{}{}, []string{}, nil)
		So(result, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})

	Convey("Sets a non-nested key to the specified value", t, func() {
		result, err := lib.Set(map[string]interface{}{}, []string{"a"}, "testing")
		So(result, ShouldResemble, map[string]interface{}{"a": "testing"})
		So(err, ShouldBeNil)
	})

	Convey("Sets a nested key to the specified value", t, func() {
		result, err := lib.Set(map[string]interface{}{}, []string{"a", "b", "c"}, "testing")
		So(result, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": "testing"}}})
		So(err, ShouldBeNil)
	})

	Convey("Returns an error when a non-nested key is already present", t, func() {
		result, err := lib.Set(map[string]interface{}{"a": true}, []string{"a"}, nil)
		So(result, ShouldResemble, map[string]interface{}{"a": true})
		So(err, ShouldNotBeNil)
	})

	Convey("Returns an error when a nested key is already present", t, func() {
		result, err := lib.Set(map[string]interface{}{"a": map[string]interface{}{"b": true}}, []string{"a", "b"}, nil)
		So(result, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": true}})
		So(err, ShouldNotBeNil)
	})

	Convey("Returns an error when a nested key is already present as a nested key", t, func() {
		result, err := lib.Set(map[string]interface{}{"a": map[string]interface{}{"b": true}}, []string{"a", "b", "c"}, nil)
		So(result, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": true}})
		So(err, ShouldNotBeNil)
	})
}

func TestGet(t *testing.T) {

	Convey("Gets a non-nested key", t, func() {
		result, err := lib.Get(map[string]interface{}{"string": "woohoo"}, []string{"string"})
		So(result, ShouldEqual, "woohoo")
		So(err, ShouldBeNil)
	})

	Convey("Gets a nested key", t, func() {
		result, err := lib.Get(map[string]interface{}{"map": map[string]interface{}{"string": "woohoo"}}, []string{"map", "string"})
		So(result, ShouldEqual, "woohoo")
		So(err, ShouldBeNil)
	})

	Convey("Returns an error when a non-existent non-nested key is requested", t, func() {
		result, err := lib.Get(map[string]interface{}{}, []string{"non-existent"})
		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("Returns an error when a non-existent nested key is requested", t, func() {
		result, err := lib.Get(map[string]interface{}{"string": "woohoo"}, []string{"string", "subString"})
		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestMerge(t *testing.T) {

	Convey("Merges non-nested keys", t, func() {
		result := lib.Merge(map[string]interface{}{"one": 1}, map[string]interface{}{"two": 2})
		So(result, ShouldResemble, map[string]interface{}{"one": 1, "two": 2})
	})

	Convey("Doesn't override existing keys", t, func() {
		result := lib.Merge(map[string]interface{}{"one": 1}, map[string]interface{}{"one": 2})
		So(result, ShouldResemble, map[string]interface{}{"one": 1})
	})

	Convey("Merges nested keys", t, func() {
		result := lib.Merge(map[string]interface{}{"one": map[string]interface{}{"one": 1}}, map[string]interface{}{"one": map[string]interface{}{"two": 2}})
		So(result, ShouldResemble, map[string]interface{}{"one": map[string]interface{}{"one": 1, "two": 2}})
	})

	Convey("Doesn't override existing nested keys", t, func() {
		result := lib.Merge(map[string]interface{}{"one": map[string]interface{}{"one": 1}}, map[string]interface{}{"one": map[string]interface{}{"one": 2}})
		So(result, ShouldResemble, map[string]interface{}{"one": map[string]interface{}{"one": 1}})
	})
}

func TestParseString(t *testing.T) {

	Convey("Parses booleans", t, func() {

		Convey("Parses 'true' into true", func() {
			result := lib.ParseString("true")
			So(result, ShouldEqual, true)
		})

		Convey("Parses 'false' into false", func() {
			result := lib.ParseString("false")
			So(result, ShouldEqual, false)
		})
	})

	Convey("Parses numbers", t, func() {

		Convey("Parses integers", func() {
			result := lib.ParseString("10")
			So(result, ShouldEqual, 10)
		})

		Convey("Parses floats", func() {
			result := lib.ParseString("10.5")
			So(result, ShouldEqual, 10.5)
		})

		Convey("Parses floats that end with .0", func() {
			result := lib.ParseString("10.0")
			So(result, ShouldEqual, 10.0)
		})
	})

	Convey("Parses JSON arrays", t, func() {
		result := lib.ParseString("[1,2,3]")
		So(result, ShouldResemble, []interface{}{float64(1), float64(2), float64(3)})
	})

	Convey("Parses JSON objects", t, func() {
		result := lib.ParseString(`{"a": 1, "b":2}`)
		So(result, ShouldResemble, map[string]interface{}{"a": float64(1), "b": float64(2)})
	})

	Convey("Parses durations", t, func() {
		result := lib.ParseString("3s")
		So(result, ShouldEqual, 3*time.Second)
	})

	Convey("Leaves every other format untouched", t, func() {
		result := lib.ParseString("Hello")
		So(result, ShouldEqual, "Hello")
	})
}

func TestParseDuration(t *testing.T) {

	Convey("Returns a duration when the string matches duration format", t, func() {
		result := lib.ParseDurationString("3s")
		So(result, ShouldEqual, 3*time.Second)
	})

	Convey("Returns the original string when the string doesn't match duration format", t, func() {
		result := lib.ParseDurationString("definitely not 3s")
		So(result, ShouldEqual, "definitely not 3s")
	})
}

func TestCasts(t *testing.T) {

	Convey("CastMap", t, func() {

		Convey("Casts an empty interface into a map", func() {
			result, err := lib.CastMap(map[string]interface{}{})
			So(result, ShouldResemble, map[string]interface{}{})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastMap(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("CastSlice", t, func() {

		Convey("Casts an empty interface into a slice", func() {
			result, err := lib.CastSlice([]interface{}{1, "two"})
			So(result, ShouldResemble, []interface{}{1, "two"})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("CastStringSlice", t, func() {

		Convey("Casts an empty interface into a string slice", func() {
			result, err := lib.CastStringSlice([]string{"one", "two"})
			So(result, ShouldResemble, []string{"one", "two"})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastStringSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("CastString", t, func() {

		Convey("Casts an empty interface into a string", func() {
			result, err := lib.CastString("one")
			So(result, ShouldEqual, "one")
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastString(5)
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("CastIntegerSlice", t, func() {

		Convey("Casts an empty interface into a integer slice", func() {
			result, err := lib.CastIntegerSlice([]int{1, 2})
			So(result, ShouldResemble, []int{1, 2})
			So(err, ShouldBeNil)
		})

		Convey("Attempts to cast a float slice to an int", func() {
			result, err := lib.CastIntegerSlice([]float64{1, 2})
			So(result, ShouldResemble, []int{1, 2})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error when casting a float slice to an int slice would truncate a value", func() {
			result, err := lib.CastIntegerSlice([]float64{1.5, 2})
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastIntegerSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("CastInteger", t, func() {

		Convey("Casts an empty interface into a integer", func() {
			result, err := lib.CastInteger(1)
			So(result, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})

		Convey("Attempts to cast a float to an int", func() {
			result, err := lib.CastInteger(float64(1))
			So(result, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})

		Convey("Returns an error when a converting to a int would truncate the value", func() {
			result, err := lib.CastInteger(1.5)
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastInteger("Hello")
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("CastBooleanSlice", t, func() {

		Convey("Casts an empty interface into a boolean slice", func() {
			result, err := lib.CastBooleanSlice([]bool{true, false})
			So(result, ShouldResemble, []bool{true, false})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastBooleanSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("CastBoolean", t, func() {

		Convey("Casts an empty interface into a boolean", func() {
			result, err := lib.CastBoolean(true)
			So(result, ShouldEqual, true)
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastBoolean(5)
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("CastFloatSlice", t, func() {

		Convey("Casts an empty interface into a float slice", func() {
			result, err := lib.CastFloatSlice([]float64{1, 2})
			So(result, ShouldResemble, []float64{1, 2})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastFloatSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("CastFloat", t, func() {

		Convey("Casts an empty interface into a float", func() {
			result, err := lib.CastFloat(3.3)
			So(result, ShouldEqual, 3.3)
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := lib.CastFloat("Hello")
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})
	})
}
