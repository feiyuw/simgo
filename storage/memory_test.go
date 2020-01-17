package storage

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	Convey("new memory storage can add and query item", t, func() {
		storage, err := NewMemoryStorage()
		So(err, ShouldBeNil)
		So(storage, ShouldNotBeNil)
		storage.Add("k1", "v1")
		storage.Add("k2", map[string]string{"k21": "v21", "k22": "v22"})
		So(len(storage.M), ShouldEqual, 2)
		storage.Remove("k1")
		So(len(storage.M), ShouldEqual, 1)
		item, err := storage.FindOne("k2")
		So(err, ShouldBeNil)
		So(item, ShouldResemble, map[string]string{"k21": "v21", "k22": "v22"})
		item, err = storage.FindOne("k1")
		So(err, ShouldNotBeNil)
		So(item, ShouldBeNil)
	})
}
