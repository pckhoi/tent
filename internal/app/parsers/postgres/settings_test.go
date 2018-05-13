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
	c.Assert(settings["a"], DeepEquals, []interface{}{"b", "d"})
}

func (s *LocalTestSuite) TestSetSettings(c *C) {
	var b = []interface{}{"b"}
	var d = []interface{}{"d"}
	settings["e"] = b
	c.Assert(len(settings["e"]), Equals, 1)
	c.Assert(settings["e"][0].(string), Equals, "b")
	setSettings("e", d)
	c.Assert(len(settings["e"]), Equals, 1)
	c.Assert(settings["e"][0].(string), Equals, "d")
}
