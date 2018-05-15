package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestAlterTableOwner(c *C) {
	val, err := tryParse(`
        ALTER TABLE my_table OWNER TO john;
    `)
	if err != nil {
		c.Error(err)
	}
	c.Assert(
		val,
		DeepEquals,
		[]interface{}{
			storage.DataRow{
				TableName: "custom/table",
				ID:        "my_table",
				Content: map[string]string{
					"owner": "john",
				},
			},
		},
	)
}
