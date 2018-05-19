package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
	"strings"
)

func (s *LocalTestSuite) TestSingleLineComment(c *C) {
	val, err := tryParse(
		strings.Join(
			[]string{
				"--",
				"-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:",
				"--",
				"",
				"CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;",
			},
			"\n",
		),
	)
	if err != nil {
		c.Error(err)
	}
	c.Assert(
		val,
		DeepEquals,
		[]interface{}{
			storage.DataRow{
				TableName: "postgres_extensions",
				ID:        "plpgsql",
				Content: map[string]string{
					"name":   "plpgsql",
					"schema": "pg_catalog",
				},
			},
		},
	)
}
