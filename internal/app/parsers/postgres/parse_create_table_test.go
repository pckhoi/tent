package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestCreateTableCustomType(c *C) {
	val, err := tryParse(`
        CREATE TYPE pet AS ENUM ('dog', 'cat', 'camel');
        create table my_table (
            my_pet pet
        );
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
				ID:        "pet",
				Content: map[string]string{
					"type":   "enum",
					"labels": "'dog','cat','camel'",
				},
			},
			[]storage.DataRow{
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_pet",
					Content: map[string]string{
						"type": "pet",
					},
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestCreateTableUnknownType(c *C) {
	val, err := tryParse(`
        create table my_table (
            my_office office
        );
    `)
	c.Assert(val, DeepEquals, []interface{}{[]storage.DataRow(nil)})
	c.Assert(err, ErrorMatches, `.+Type office is not defined`)
}

func (s *LocalTestSuite) TestCreateTableStmt(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            id integer NOT NULL,
            an_int int,
            smallnum smallint,
            bignum bigint,
            decimalnum decimal,
            number numeric,
            realty real,
            small_serial smallserial,
            the_serial serial,
            big_serial bigserial,
            important boolean NOT NULL
        );
    `)
	if err != nil {
		c.Error(err)
	}
	c.Assert(
		val,
		DeepEquals,
		[]interface{}{
			[]storage.DataRow{
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "id",
					Content: map[string]string{
						"type":     "integer",
						"not_null": "true",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "an_int",
					Content: map[string]string{
						"type": "integer",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "smallnum",
					Content: map[string]string{
						"type": "smallint",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "bignum",
					Content: map[string]string{
						"type": "bigint",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "decimalnum",
					Content: map[string]string{
						"type": "decimal",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "number",
					Content: map[string]string{
						"type": "numeric",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "realty",
					Content: map[string]string{
						"type": "real",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "small_serial",
					Content: map[string]string{
						"type": "smallserial",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "the_serial",
					Content: map[string]string{
						"type": "serial",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "big_serial",
					Content: map[string]string{
						"type": "bigserial",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "important",
					Content: map[string]string{
						"type":     "boolean",
						"not_null": "true",
					},
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestDateTimeTypes(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            datetime_with_tz timestamp with time zone,
            datetime_with_prec timestamp 6,
            datetime_with_prec_and_tz timestamp 5 with time zone,
            datetime_plain timestamp,
            datetime_without_tz timestamp without time zone,
            time_with_tz time with time zone,
            time_with_prec time 6,
            time_with_prec_and_tz time 5 with time zone,
            time_plain time,
            time_without_tz time without time zone,
            date_plain date
        );
    `)
	if err != nil {
		c.Error(err)
	}
	c.Assert(
		val,
		DeepEquals,
		[]interface{}{
			[]storage.DataRow{
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "datetime_with_tz",
					Content: map[string]string{
						"type": "datetimetz",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "datetime_with_prec",
					Content: map[string]string{
						"type":          "datetime",
						"sec_precision": "6",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "datetime_with_prec_and_tz",
					Content: map[string]string{
						"type":          "datetimetz",
						"sec_precision": "5",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "datetime_plain",
					Content: map[string]string{
						"type": "datetime",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "datetime_without_tz",
					Content: map[string]string{
						"type": "datetime",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "time_with_tz",
					Content: map[string]string{
						"type": "timetz",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "time_with_prec",
					Content: map[string]string{
						"type":          "time",
						"sec_precision": "6",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "time_with_prec_and_tz",
					Content: map[string]string{
						"type":          "timetz",
						"sec_precision": "5",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "time_plain",
					Content: map[string]string{
						"type": "time",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "time_without_tz",
					Content: map[string]string{
						"type": "time",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "date_plain",
					Content: map[string]string{
						"type": "date",
					},
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestCharacterTypes(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            firstname character varying(10),
            lastname character(12),
            middlename varchar(5),
            title char(6),
            streetname character varying,
            letter character,
            summary text
        );
    `)
	if err != nil {
		c.Error(err)
	}
	c.Assert(
		val,
		DeepEquals,
		[]interface{}{
			[]storage.DataRow{
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "firstname",
					Content: map[string]string{
						"type":   "varchar",
						"length": "10",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "lastname",
					Content: map[string]string{
						"type":   "char",
						"length": "12",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "middlename",
					Content: map[string]string{
						"type":   "varchar",
						"length": "5",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "title",
					Content: map[string]string{
						"type":   "char",
						"length": "6",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "streetname",
					Content: map[string]string{
						"type": "varchar",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "letter",
					Content: map[string]string{
						"type":   "char",
						"length": "1",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "summary",
					Content: map[string]string{
						"type": "text",
					},
				},
			},
		},
	)
}
