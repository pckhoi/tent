package postgres

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type LocalTestSuite struct{}

var _ = Suite(&LocalTestSuite{})

func testParse(c *C, s string) interface{} {
	res, err := Parse("test", []byte(s))
	if err != nil {
		c.Error(err)
	}
	return res
}
