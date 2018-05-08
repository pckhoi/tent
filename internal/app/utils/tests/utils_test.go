package test

import (
	"github.com/pckhoi/tent/internal/app/utils"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestMergeStringSlices(c *C) {
	c.Assert(
		utils.MergeStringSlices(
			[]string{"a", "b"},
			[]string{"b", "c"},
		),
		DeepEquals,
		[]string{"a", "b", "c"},
	)
	c.Assert(
		utils.MergeStringSlices(
			[]string{"a", "b"},
			[]string{"d", "a"},
		),
		DeepEquals,
		[]string{"a", "b", "d"},
	)
	c.Assert(
		utils.MergeStringSlices(
			nil,
			[]string{"d", "a"},
		),
		DeepEquals,
		[]string{"d", "a"},
	)
	c.Assert(
		utils.MergeStringSlices(
			[]string{"d", "a"},
			nil,
		),
		DeepEquals,
		[]string{"d", "a"},
	)
	c.Assert(
		utils.MergeStringSlices(
			[]string{"b", "a"},
			[]string{"a"},
		),
		DeepEquals,
		[]string{"b", "a"},
	)
}

func (s *MySuite) TestStringSliceEqual(c *C) {
	c.Assert(
		utils.StringSliceEqual(
			[]string{"a", "b"},
			[]string{"a", "b"}),
		Equals,
		true,
	)
	c.Assert(
		utils.StringSliceEqual(
			nil,
			nil),
		Equals,
		true,
	)
	c.Assert(
		utils.StringSliceEqual(
			[]string{"a"},
			nil),
		Equals,
		false,
	)
	c.Assert(
		utils.StringSliceEqual(
			[]string{"a"},
			[]string{"a", "b"}),
		Equals,
		false,
	)
}
