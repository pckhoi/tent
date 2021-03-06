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
            "order" int not null,
            cats smallint constraint cats_amount check ((cats > 5)),
            dogs smallint check ((dogs < 10)) no inherit,
            constraint my_table_pets_check check ((cats + dogs > 2)),
            check ((cats + dogs < 15)) no inherit
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
					ID:        "book",
					Content: map[string]string{
						"type":     "varchar",
						"not_null": "false",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "order",
					Content: map[string]string{
						"type":     "integer",
						"not_null": "true",
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
				storage.DataRow{
					TableName: "constraint/my_table",
					ID:        "my_table_pets_check",
					Content: map[string]string{
						"check_def":       "((cats + dogs > 2))",
						"constraint_name": "my_table_pets_check",
					},
				},
				storage.DataRow{
					TableName: "constraint/my_table",
					ID:        "0",
					Content: map[string]string{
						"check_def":        "((cats + dogs < 15))",
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

func (s *LocalTestSuite) TestArrayTypes(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            my_integer_list integer[],
            my_char_list character varying(20)[][]
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
					ID:        "my_integer_list",
					Content: map[string]string{
						"type":             "integer",
						"array_dimensions": "1",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_char_list",
					Content: map[string]string{
						"type":             "varchar",
						"length":           "20",
						"array_dimensions": "2",
					},
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestNumericTypes(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            var1 numeric,
            var2 numeric(6),
            var3 numeric(6, 3)
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
					ID:        "var1",
					Content: map[string]string{
						"type": "numeric",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "var2",
					Content: map[string]string{
						"type":      "numeric",
						"precision": "6",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "var3",
					Content: map[string]string{
						"type":      "numeric",
						"precision": "6",
						"scale":     "3",
					},
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestGeographyTypes(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            geog geography(POINT),
            my_point geography(POINT,4269),
            other_point geography(POINT,4326),
            my_line geography(LINESTRING),
            my_poly geography(POLYGON,4267),
            many_point geography(MULTIPOINT),
            many_line geography(MULTILINESTRING),
            many_poly geography(MULTIPOLYGON),
            many_geo geography(GEOMETRYCOLLECTION)
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
					ID:        "geog",
					Content: map[string]string{
						"type":    "geography",
						"subtype": "point",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_point",
					Content: map[string]string{
						"type":    "geography",
						"subtype": "point",
						"srid":    "4269",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "other_point",
					Content: map[string]string{
						"type":    "geography",
						"subtype": "point",
						"srid":    "4326",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_line",
					Content: map[string]string{
						"type":    "geography",
						"subtype": "linestring",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_poly",
					Content: map[string]string{
						"type":    "geography",
						"subtype": "polygon",
						"srid":    "4267",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "many_point",
					Content: map[string]string{
						"type":    "geography",
						"subtype": "multipoint",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "many_line",
					Content: map[string]string{
						"type":    "geography",
						"subtype": "multilinestring",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "many_poly",
					Content: map[string]string{
						"type":    "geography",
						"subtype": "multipolygon",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "many_geo",
					Content: map[string]string{
						"type":    "geography",
						"subtype": "geometrycollection",
					},
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestGeometryTypes(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            geom geometry(POINT),
            my_point geometry(POINT,4269),
            other_point geometry(POINT,4326),
            my_line geometry(LINESTRING),
            my_poly geometry(POLYGON,4267),
            many_point geometry(MULTIPOINT),
            many_line geometry(MULTILINESTRING),
            many_poly geometry(MULTIPOLYGON),
            many_geo geometry(GEOMETRYCOLLECTION)
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
					ID:        "geom",
					Content: map[string]string{
						"type":    "geometry",
						"subtype": "point",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_point",
					Content: map[string]string{
						"type":    "geometry",
						"subtype": "point",
						"srid":    "4269",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "other_point",
					Content: map[string]string{
						"type":    "geometry",
						"subtype": "point",
						"srid":    "4326",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_line",
					Content: map[string]string{
						"type":    "geometry",
						"subtype": "linestring",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "my_poly",
					Content: map[string]string{
						"type":    "geometry",
						"subtype": "polygon",
						"srid":    "4267",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "many_point",
					Content: map[string]string{
						"type":    "geometry",
						"subtype": "multipoint",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "many_line",
					Content: map[string]string{
						"type":    "geometry",
						"subtype": "multilinestring",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "many_poly",
					Content: map[string]string{
						"type":    "geometry",
						"subtype": "multipolygon",
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "many_geo",
					Content: map[string]string{
						"type":    "geometry",
						"subtype": "geometrycollection",
					},
				},
			},
		},
	)
}

func (s *LocalTestSuite) TestCollation(c *C) {
	val, err := tryParse(`
        CREATE TABLE my_table (
            var1 varchar collate "en_US",
            var2 varchar collate pg_catalog."C"
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
					ID:        "var1",
					Content: map[string]string{
						"type":      "varchar",
						"collation": `"en_US"`,
					},
				},
				storage.DataRow{
					TableName: "schema/my_table",
					ID:        "var2",
					Content: map[string]string{
						"type":      "varchar",
						"collation": `pg_catalog."C"`,
					},
				},
			},
		},
	)
}
