package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStringSliceExcept(t *testing.T) {
	Convey("StringSliceExcept", t, func() {
		Convey("return new slice without unwanted items", func() {
			result := StringSliceExcept([]string{
				"1",
				"2",
				"3",
			}, []string{
				"1",
				"3",
			})
			So(len(result), ShouldEqual, 1)
			So(result[0], ShouldEqual, "2")
		})

		Convey("should return all items if no items is filtered", func() {
			result := StringSliceExcept([]string{
				"1",
				"2",
				"3",
			}, []string{
				"4",
			})
			So(len(result), ShouldEqual, 3)
		})

		Convey("works with duplicated items to filter", func() {
			result := StringSliceExcept([]string{
				"1",
				"2",
				"3",
				"4",
				"5",
				"6",
				"7",
				"8",
				"9",
			}, []string{
				"4",
				"4",
				"1",
				"2",
				"3",
				"1",
				"2",
				"3",
				"7",
				"8",
				"9",
			})
			So(len(result), ShouldEqual, 2)
		})
	})
}

func TestStringSliceContainAll(t *testing.T) {
	Convey("StringSliceContainAll", t, func() {
		Convey("return true on container have all elements", func() {
			result := StringSliceContainAll([]string{
				"god",
				"man",
			}, []string{
				"god",
			})
			So(result, ShouldEqual, true)
		})
		Convey("return true on target slice is empty", func() {
			result := StringSliceContainAll([]string{
				"god",
				"man",
			}, []string{})
			So(result, ShouldEqual, true)
		})
		Convey("return false on container don't have all elements", func() {
			result := StringSliceContainAll([]string{
				"god",
				"man",
			}, []string{
				"devil",
			})
			So(result, ShouldEqual, false)
		})
	})
}
