package logger

import (
	"bytes"
	"errors"
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)
	log.SetFlags(0)

	Convey("set to different log level", t, func() {
		SetLogLevel("debug")
		So(logLevel, ShouldEqual, DEBUG)
		SetLogLevel("info")
		So(logLevel, ShouldEqual, INFO)
		SetLogLevel("warn")
		So(logLevel, ShouldEqual, WARN)
		SetLogLevel("error")
		So(logLevel, ShouldEqual, ERROR)
		SetLogLevel("fatal")
		So(logLevel, ShouldEqual, FATAL)
	})

	Convey("message show in different level", t, func() {
		SetLogLevel("debug")
		defer buf.Reset()

		Debug("test", "hello", "world")
		So(buf.String(), ShouldEqual, "D\ttest\thello world\n")
		buf.Reset()

		Info("test", "hello", "world")
		So(buf.String(), ShouldEqual, "I\ttest\thello world\n")
		buf.Reset()

		Warnf("demo", "hello %s world", "china")
		So(buf.String(), ShouldEqual, "W\tdemo\thello china world\n")
		buf.Reset()

		Errorf("err", "hello %s world", errors.New("china"))
		So(buf.String(), ShouldEqual, "E\terr\thello china world\n")
		buf.Reset()

		SetLogLevel("error")
		Debug("xxx", "xxxx")
		Info("xxx", "don")
		Warn("xxx", "don")
		So(buf.String(), ShouldEqual, "")
		Error("xxx", "don")
		So(buf.String(), ShouldEqual, "E\txxx\tdon\n")
		buf.Reset()
	})
}
