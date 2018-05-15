package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestCreateSequence(c *C) {
	val, err := tryParse(`
        CREATE SEQUENCE my_seq;
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
				ID:        "my_seq",
				Content:   map[string]string{},
			},
		},
	)
}

func (s *LocalTestSuite) TestCreateSequenceWithIncrement(c *C) {
	val, err := tryParse(`
        CREATE SEQUENCE my_seq
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1
            CYCLE
            OWNED BY table.column;
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
				ID:        "my_seq",
				Content: map[string]string{
					"start":     "1",
					"increment": "1",
					"cache":     "1",
					"cycle":     "true",
					"owned_by":  "table/column",
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestCreateSequenceWithMinMax(c *C) {
	val, err := tryParse(`
        CREATE SEQUENCE my_seq
            INCREMENT 1
            START 2
            MAXVALUE 10000
            MINVALUE 2
            OWNED BY NONE
            NO CYCLE;
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
				ID:        "my_seq",
				Content: map[string]string{
					"start":     "2",
					"increment": "1",
					"maxvalue":  "10000",
					"minvalue":  "2",
					"cycle":     "false",
				},
			},
		},
	)
}
