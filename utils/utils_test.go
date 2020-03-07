package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMin(t *testing.T) {
	Convey("one number", t, func() {
		So(Min(1), ShouldEqual, 1)
	})

	Convey("Two equal numbers", t, func() {
		So(Min(1, 1), ShouldEqual, 1)
	})

	Convey("Two different numbers", t, func() {
		So(Min(1, 2), ShouldEqual, 1)
	})

	Convey("multiple different numbers", t, func() {
		So(Min(1, 2, 1, 4, 3, 0), ShouldEqual, 0)
	})

	Convey("less than 0 numbers", t, func() {
		So(Min(-1, -3, 0, 1, 4), ShouldEqual, -3)
	})
}
