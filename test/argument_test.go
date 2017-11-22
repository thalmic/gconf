package test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/thalmic/gconf/lib"
	"testing"
)

func TestParseArguments(t *testing.T) {

	Convey("Returns an empty map if there are no arguments", t, func() {
		loader := lib.NewArgumentLoader("", "")
		result, err := loader.ParseArguments([]string{})
		So(result, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})

	Convey("Parses arguments without a separator or prefix", t, func() {
		loader := lib.NewArgumentLoader("", "")

		Convey("Ignores arguments not starting with '-' or '--'", func() {
			result, err := loader.ParseArguments([]string{"string=testing"})
			So(result, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("Ignores arguments without an '='", func() {
			result, err := loader.ParseArguments([]string{"testing"})
			So(result, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("Parses arguments starting with '-'", func() {
			result, err := loader.ParseArguments([]string{"-string=testing"})
			So(result, ShouldResemble, map[string]interface{}{"string": "testing"})
			So(err, ShouldBeNil)
		})

		Convey("Parses arguments starting with '--'", func() {
			result, err := loader.ParseArguments([]string{"--string=testing"})
			So(result, ShouldResemble, map[string]interface{}{"string": "testing"})
			So(err, ShouldBeNil)
		})
	})

	Convey("Parses arguments with a separator", t, func() {
		loader := lib.NewArgumentLoader("__", "")

		Convey("Nests the configuration using the separator", func() {
			result, err := loader.ParseArguments([]string{"--map__string=testing"})
			So(result, ShouldResemble, map[string]interface{}{"map": map[string]interface{}{"string": "testing"}})
			So(err, ShouldBeNil)
		})

		Convey("Fails when a nested option would override another option", func() {
			result, err := loader.ParseArguments([]string{"--test=testing", "--test__stuff=things"})
			So(result, ShouldResemble, map[string]interface{}{"test": "testing"})
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Parses arguments with a prefix", t, func() {
		loader := lib.NewArgumentLoader("", "TEST")

		Convey("Strips the prefix from the argument", func() {
			result, err := loader.ParseArguments([]string{"--TESTing=testing"})
			So(result, ShouldResemble, map[string]interface{}{"ing": "testing"})
			So(err, ShouldBeNil)
		})

		Convey("Ignores arguments without the prefix", func() {
			result, err := loader.ParseArguments([]string{"--notTesting=testing"})
			So(result, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})
	})
}
