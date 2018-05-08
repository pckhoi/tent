package test

import (
	"github.com/pckhoi/tent/internal/app/parsers/postgres"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})
var filename = "test"

func tryParse(c *C, s string) interface{} {
	res, err := postgres.Parse(filename, []byte(s))
	if err != nil {
		c.Error(err)
	}
	return res
}

func (s *MySuite) TestSetEqualNumber(c *C) {
	c.Assert(
		tryParse(c, "SET my_key = 0;"),
		DeepEquals,
		[]interface{}{
			postgres.Update{
				TableName: "postgres_settings",
				Row: map[string]interface{}{
					"name":    "my_key",
					"setting": "0",
					"type":    "int64",
				},
			},
		},
	)
}

func (s *MySuite) TestSetEqualStringConstant(c *C) {
	c.Assert(
		tryParse(c, "SET client_encoding = 'UTF8';"),
		DeepEquals,
		[]interface{}{
			postgres.Update{
				TableName: "postgres_settings",
				Row: map[string]interface{}{
					"name":    "client_encoding",
					"setting": "UTF8",
					"type":    "postgres.String",
				},
			},
		},
	)
}
