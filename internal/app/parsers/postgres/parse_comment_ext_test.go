package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestCommentExtensionStmt(c *C) {
	c.Assert(
		testParse(c, `
            COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';
            comment on extension postgis is 'PostGIS geometry, geography';
        `),
		DeepEquals,
		[]interface{}{
			storage.DataRow{
				TableName: "postgres_extensions",
				ID:        "plpgsql",
				Content: map[string]string{
					"name":    "plpgsql",
					"comment": "PL/pgSQL procedural language",
				},
			},
			storage.DataRow{
				TableName: "postgres_extensions",
				ID:        "postgis",
				Content: map[string]string{
					"name":    "postgis",
					"comment": "PostGIS geometry, geography",
				},
			},
		},
	)
}
