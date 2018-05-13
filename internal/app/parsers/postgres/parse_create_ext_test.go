package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestCreateExtensionStmt(c *C) {
	val, err := tryParse(`
        CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
        CREATE EXTENSION postgis WITH SCHEMA public;
        create extension if not exists postgis with schema athur;
    `)
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
			storage.DataRow{
				TableName: "postgres_extensions",
				ID:        "postgis",
				Content: map[string]string{
					"name":   "postgis",
					"schema": "public",
				},
			},
			storage.DataRow{
				TableName: "postgres_extensions",
				ID:        "postgis",
				Content: map[string]string{
					"name":   "postgis",
					"schema": "athur",
				},
			},
		},
	)
}
