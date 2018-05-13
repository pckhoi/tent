package postgres

import (
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestSeStmt(c *C) {
	settings["search_path"] = nil
	settings["client_encoding"] = nil
	settings["default_with_oids"] = nil

	val, err := tryParse(`
		SET search_path = public;
		SET client_encoding TO 'UTF8';
		set default_with_oids to false;
    `)
	if err != nil {
		c.Error(err)
	}
	c.Assert(val, DeepEquals, []interface{}{nil, nil, nil})

	c.Assert(len(settings["search_path"]), Equals, 1)
	c.Assert(interfaceToString(settings["search_path"][0]), Equals, "public")

	c.Assert(len(settings["client_encoding"]), Equals, 1)
	c.Assert(interfaceToString(settings["client_encoding"][0]), Equals, "UTF8")

	c.Assert(len(settings["default_with_oids"]), Equals, 1)
	c.Assert(settings["default_with_oids"][0].(bool), Equals, false)
}
