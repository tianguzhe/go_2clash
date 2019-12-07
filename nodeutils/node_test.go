package nodeutils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDecodeInfoByByte(t *testing.T) {

	Convey("base64补全", t, func() {
		Convey("1位补全", func() {
			So(DecodeInfoByByte("a"), ShouldEqual, "a===")
		})
		Convey("2位补全", func() {
			So(DecodeInfoByByte("aa"), ShouldEqual, "aa==")
		})
		Convey("3位补全", func() {
			So(DecodeInfoByByte("aaa"), ShouldEqual, "aaa=")
		})
		Convey("4位补全", func() {
			So(DecodeInfoByByte("aaaa"), ShouldEqual, "aaaa")
		})
	})
}
