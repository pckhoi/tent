package postgres

import (
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestUpdateSettings(c *C) {
	var b = []interface{}{"b"}
	var d = []interface{}{"d"}
	settings["a"] = b
	c.Assert(len(settings["a"]), Equals, 1)
	c.Assert(settings["a"][0].(string), Equals, "b")
	updateSettings("a", d)
	c.Assert(len(settings["a"]), Equals, 1)
	c.Assert(settings["a"][0].(string), Equals, "d")
}
