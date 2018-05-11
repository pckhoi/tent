package postgres

import (
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestInt64Value(c *C) {
	settings["my_key"] = nil
	testParse(c, "SET my_key = 0;")
	c.Assert(settings["my_key"][0].(int64), Equals, int64(0))
}

func (s *LocalTestSuite) TestIdentifierValue(c *C) {
	settings["my_key"] = nil
	testParse(c, "SET my_key = my_value;")
	c.Assert(interfaceToString(settings["my_key"][0]), Equals, "my_value")
}

func (s *LocalTestSuite) TestStringConstantValue(c *C) {
	settings["my_key"] = nil
	testParse(c, "SET my_key TO 'UTF8';")
	c.Assert(interfaceToString(settings["my_key"][0]), Equals, "UTF8")
}

func (s *LocalTestSuite) TestCommaSeparatedValues(c *C) {
	settings["my_key"] = nil
	testParse(c, "SET my_key = public, pg_catalog;")
	c.Assert(interfaceToString(settings["my_key"][0]), Equals, "public")
	c.Assert(interfaceToString(settings["my_key"][1]), Equals, "pg_catalog")
}

func (s *LocalTestSuite) TestFalseBooleanValues(c *C) {
	settings["my_key"] = nil
	testParse(c, "set my_key to FALSE, false, no, off, f, n, 'false', 'no', 'off', 'f', 'n';")
	c.Assert(len(settings["my_key"]), Equals, 11)
	for _, v := range settings["my_key"] {
		c.Assert(v.(bool), Equals, false)
	}
}

func (s *LocalTestSuite) TestTrueBooleanValues(c *C) {
	settings["my_key"] = nil
	testParse(c, "set my_key to TRUE, true, yes, on, t, y, 'true', 'yes', 'on', 't', 'y';")
	c.Assert(len(settings["my_key"]), Equals, 11)
	for _, v := range settings["my_key"] {
		c.Assert(v.(bool), Equals, true)
	}
}
