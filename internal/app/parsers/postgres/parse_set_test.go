package postgres

import (
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestSetStmt(c *C) {
	settings["my_key"] = nil
	testParse(c, "SET my_key = 0;")
	c.Assert(len(settings["my_key"]), Equals, 1)
	c.Assert(settings["my_key"][0].(int64), Equals, int64(0))
}

func (s *LocalTestSuite) TestSetMultipleValues(c *C) {
	settings["search_path"] = nil
	testParse(c, "SET search_path = public, pg_catalog;")
	c.Assert(len(settings["search_path"]), Equals, 2)
	c.Assert(interfaceToString(settings["search_path"][0]), Equals, "public")
	c.Assert(interfaceToString(settings["search_path"][1]), Equals, "pg_catalog")
}

func (s *LocalTestSuite) TestSetStringConstants(c *C) {
	settings["client_encoding"] = nil
	testParse(c, "SET client_encoding TO 'UTF8';")
	c.Assert(len(settings["client_encoding"]), Equals, 1)
	c.Assert(interfaceToString(settings["client_encoding"][0]), Equals, "UTF8")
}

func (s *LocalTestSuite) TestSetLowerCase(c *C) {
	settings["default_with_oids"] = nil
	testParse(c, "set default_with_oids to false;")
	c.Assert(len(settings["default_with_oids"]), Equals, 1)
	c.Assert(settings["default_with_oids"][0].(bool), Equals, false)
}
