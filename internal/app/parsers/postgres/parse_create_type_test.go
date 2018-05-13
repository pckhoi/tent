package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestTypeExists(c *C) {
	err := typeExists("animal")
	c.Assert(
		err.Error(),
		Equals,
		"Type animal is not defined",
	)
}

func (s *LocalTestSuite) TestCreateTypeStmt(c *C) {
	val, err := tryParse(`
        CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy', 'so''so');
    `)
	if err != nil {
		c.Error(err)
	}
	c.Assert(
		val,
		DeepEquals,
		[]interface{}{
			storage.DataRow{
				TableName: "custom/type",
				ID:        "mood",
				Content: map[string]string{
					"type":   "enum",
					"labels": "'sad','ok','happy','so''so'",
				},
			},
		},
	)
	c.Assert(
		typeExists("mood"),
		Equals,
		nil,
	)
}
