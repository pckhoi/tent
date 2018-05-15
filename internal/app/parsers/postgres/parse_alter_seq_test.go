package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestAlterSequenceOwner(c *C) {
	val, err := tryParse(`
        ALTER SEQUENCE my_sequence OWNED BY my_table.my_column;
    `)
	if err != nil {
		c.Error(err)
	}
	c.Assert(
		val,
		DeepEquals,
		[]interface{}{
			storage.DataRow{
				TableName: "custom/sequence",
				ID:        "my_sequence",
				Content: map[string]string{
					"owned_by": "my_table/my_column",
				},
			},
		},
	)
}
