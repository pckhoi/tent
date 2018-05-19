package postgres

import (
	"fmt"
	"github.com/pckhoi/tent/internal/app/storage"
	. "gopkg.in/check.v1"
	"strings"
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

func (s *LocalTestSuite) TestCreateTableWithColumnConstraints(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            id integer not null,
            book varchar null,
            cats smallint constraint cats_amount check ((cats > 5)),
            dogs smallint check ((dogs < 10)) no inherit
        );
    `, Debug(true))
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
					ID:        "book",
					Content: map[string]string{
						"type":     "varchar",
						"not_null": "false",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "cats",
					Content: map[string]string{
						"type":            "smallint",
						"constraint_name": "cats_amount",
						"check_def":       "((cats > 5))",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "dogs",
					Content: map[string]string{
						"type":             "smallint",
						"check_def":        "((dogs < 10))",
						"check_no_inherit": "true",
					},
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestCreateTableStmt(c *C) {
	types := map[string]string{
		"integer":       "integer",
		"int":           "integer",
		"smallint":      "smallint",
		"bigint":        "bigint",
		"decimal":       "decimal",
		"numeric":       "numeric",
		"real":          "real",
		"smallserial":   "smallserial",
		"serial":        "serial",
		"bigserial":     "bigserial",
		"boolean":       "boolean",
		"money":         "money",
		"bytea":         "bytea",
		"point":         "point",
		"line":          "line",
		"lseg":          "lseg",
		"box":           "box",
		"path":          "path",
		"polygon":       "polygon",
		"circle":        "circle",
		"cidr":          "cidr",
		"inet":          "inet",
		"macaddr":       "macaddr",
		"uuid":          "uuid",
		"xml":           "xml",
		"json":          "json",
		"jsonb":         "jsonb",
		"oid":           "oid",
		"regproc":       "regproc",
		"regprocedure":  "regprocedure",
		"regoper":       "regoper",
		"regoperator":   "regoperator",
		"regclass":      "regclass",
		"regtype":       "regtype",
		"regrole":       "regrole",
		"regnamespace":  "regnamespace",
		"regconfig":     "regconfig",
		"regdictionary": "regdictionary",
	}

	fields := []string{}
	rows := []storage.DataRow{}
	for k, v := range types {
		fields = append(fields, fmt.Sprintf("my_%s %s", k, k))
		rows = append(rows, storage.DataRow{
			TableName: "schema/my_table",
			ID:        fmt.Sprintf("my_%s", k),
			Content: map[string]string{
				"type": v,
			},
		})
	}
	statement := fmt.Sprintf("CREATE TABLE my_table (%s);", strings.Join(fields, ", "))

	val, err := tryParse(statement)
	if err != nil {
		c.Error(err)
	}
	c.Assert(
		val,
		DeepEquals,
		[]interface{}{
			rows,
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

func (s *LocalTestSuite) TestBitTypes(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            my_bit_var bit varying(10),
            my_bit bit(12),
            my_unlimited_bit bit varying,
            single_bit bit
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
					ID:        "my_bit_var",
					Content: map[string]string{
						"type":   "bitvar",
						"length": "10",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_bit",
					Content: map[string]string{
						"type":   "bit",
						"length": "12",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_unlimited_bit",
					Content: map[string]string{
						"type": "bitvar",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "single_bit",
					Content: map[string]string{
						"type":   "bit",
						"length": "1",
					},
				},
			},
		},
	)
}
