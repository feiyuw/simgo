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

func TestAtoUint64(t *testing.T) {
	Convey("string to uint64", t, func() {
		v, err := AtoUint64("1")
		So(err, ShouldBeNil)
		So(v, ShouldEqual, uint64(1))
		v, err = AtoUint64("123")
		So(err, ShouldBeNil)
		So(v, ShouldEqual, uint64(123))
		_, err = AtoUint64("-123")
		So(err, ShouldNotBeNil)
		_, err = AtoUint64("")
		So(err, ShouldNotBeNil)
		_, err = AtoUint64("12a")
		So(err, ShouldNotBeNil)
		_, err = AtoUint64("12.34")
		So(err, ShouldNotBeNil)
	})
}
