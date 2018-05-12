package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
)

func (s *LocalTestSuite) TestCreateTableStmt(c *C) {
	c.Assert(
		testParse(c, `
            CREATE TABLE my_table (
                id integer NOT NULL,
                smallnum smallint,
                bignum bigint,
                decimalnum decimal,
                number numeric,
                realty real,
                small_serial smallserial,
                the_serial serial,
                big_serial bigserial,
                important boolean NOT NULL,
                last_activity timestamp with time zone
            );
        `),
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
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "last_activity",
					Content: map[string]string{
						"type": "datetimetz",
					},
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestDateTimeTypes(c *C) {
	c.Assert(
		testParse(c, `
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
        `),
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
