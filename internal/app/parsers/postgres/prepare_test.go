package postgres

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type LocalTestSuite struct{}

var _ = Suite(&LocalTestSuite{})

func tryParse(s string, opts ...Option) (interface{}, error) {
	return Parse("test", []byte(s), opts...)
}
