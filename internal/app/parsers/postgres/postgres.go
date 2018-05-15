// Package postgres parses SQL dump files produced pgdump command
package postgres

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/imdario/mergo"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "SQL",
			pos:  position{line: 10, col: 1, offset: 151},
			expr: &actionExpr{
				pos: position{line: 10, col: 8, offset: 158},
				run: (*parser).callonSQL1,
				expr: &labeledExpr{
					pos:   position{line: 10, col: 8, offset: 158},
					label: "stmts",
					expr: &oneOrMoreExpr{
						pos: position{line: 10, col: 14, offset: 164},
						expr: &ruleRefExpr{
							pos:  position{line: 10, col: 14, offset: 164},
							name: "Stmt",
						},
					},
				},
			},
		},
		{
			name: "Stmt",
			pos:  position{line: 14, col: 1, offset: 197},
			expr: &actionExpr{
				pos: position{line: 14, col: 9, offset: 205},
				run: (*parser).callonStmt1,
				expr: &seqExpr{
					pos: position{line: 14, col: 9, offset: 205},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 14, col: 9, offset: 205},
							expr: &ruleRefExpr{
								pos:  position{line: 14, col: 9, offset: 205},
								name: "Comment",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 14, col: 18, offset: 214},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 14, col: 20, offset: 216},
							label: "stmt",
							expr: &choiceExpr{
								pos: position{line: 14, col: 27, offset: 223},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 14, col: 27, offset: 223},
										name: "SetStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 14, col: 37, offset: 233},
										name: "CreateTableStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 14, col: 55, offset: 251},
										name: "CreateSeqStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 14, col: 71, offset: 267},
										name: "CreateExtensionStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 14, col: 93, offset: 289},
										name: "CreateTypeStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 14, col: 110, offset: 306},
										name: "AlterTableStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 14, col: 127, offset: 323},
										name: "CommentExtensionStmt",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CreateTableStmt",
			pos:  position{line: 28, col: 1, offset: 1900},
			expr: &actionExpr{
				pos: position{line: 28, col: 20, offset: 1919},
				run: (*parser).callonCreateTableStmt1,
				expr: &seqExpr{
					pos: position{line: 28, col: 20, offset: 1919},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 28, col: 20, offset: 1919},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 30, offset: 1929},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 28, col: 33, offset: 1932},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 42, offset: 1941},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 45, offset: 1944},
							label: "tablename",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 55, offset: 1954},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 61, offset: 1960},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 28, col: 63, offset: 1962},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 67, offset: 1966},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 69, offset: 1968},
							label: "fields",
							expr: &seqExpr{
								pos: position{line: 28, col: 78, offset: 1977},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 28, col: 78, offset: 1977},
										name: "FieldDef",
									},
									&zeroOrMoreExpr{
										pos: position{line: 28, col: 87, offset: 1986},
										expr: &seqExpr{
											pos: position{line: 28, col: 89, offset: 1988},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 28, col: 89, offset: 1988},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 28, col: 91, offset: 1990},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 28, col: 95, offset: 1994},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 28, col: 97, offset: 1996},
													name: "FieldDef",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 111, offset: 2010},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 28, col: 113, offset: 2012},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 117, offset: 2016},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 28, col: 119, offset: 2018},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 123, offset: 2022},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "FieldDef",
			pos:  position{line: 48, col: 1, offset: 2640},
			expr: &actionExpr{
				pos: position{line: 48, col: 13, offset: 2652},
				run: (*parser).callonFieldDef1,
				expr: &seqExpr{
					pos: position{line: 48, col: 13, offset: 2652},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 48, col: 13, offset: 2652},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 48, col: 18, offset: 2657},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 48, col: 24, offset: 2663},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 48, col: 27, offset: 2666},
							label: "dataType",
							expr: &ruleRefExpr{
								pos:  position{line: 48, col: 36, offset: 2675},
								name: "DataType",
							},
						},
						&labeledExpr{
							pos:   position{line: 48, col: 45, offset: 2684},
							label: "notnull",
							expr: &zeroOrOneExpr{
								pos: position{line: 48, col: 53, offset: 2692},
								expr: &seqExpr{
									pos: position{line: 48, col: 55, offset: 2694},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 48, col: 55, offset: 2694},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 48, col: 58, offset: 2697},
											val:        "not",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 48, col: 65, offset: 2704},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 48, col: 68, offset: 2707},
											val:        "null",
											ignoreCase: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DataType",
			pos:  position{line: 60, col: 1, offset: 2955},
			expr: &actionExpr{
				pos: position{line: 60, col: 13, offset: 2967},
				run: (*parser).callonDataType1,
				expr: &labeledExpr{
					pos:   position{line: 60, col: 13, offset: 2967},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 60, col: 17, offset: 2971},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 60, col: 17, offset: 2971},
								name: "TimestampT",
							},
							&ruleRefExpr{
								pos:  position{line: 60, col: 30, offset: 2984},
								name: "TimeT",
							},
							&ruleRefExpr{
								pos:  position{line: 60, col: 38, offset: 2992},
								name: "VarcharT",
							},
							&ruleRefExpr{
								pos:  position{line: 60, col: 49, offset: 3003},
								name: "CharT",
							},
							&ruleRefExpr{
								pos:  position{line: 60, col: 57, offset: 3011},
								name: "BitVarT",
							},
							&ruleRefExpr{
								pos:  position{line: 60, col: 67, offset: 3021},
								name: "BitT",
							},
							&ruleRefExpr{
								pos:  position{line: 60, col: 74, offset: 3028},
								name: "IntT",
							},
							&ruleRefExpr{
								pos:  position{line: 60, col: 81, offset: 3035},
								name: "PgOidT",
							},
							&ruleRefExpr{
								pos:  position{line: 60, col: 90, offset: 3044},
								name: "OtherT",
							},
							&ruleRefExpr{
								pos:  position{line: 60, col: 99, offset: 3053},
								name: "CustomT",
							},
						},
					},
				},
			},
		},
		{
			name: "TimestampT",
			pos:  position{line: 64, col: 1, offset: 3086},
			expr: &actionExpr{
				pos: position{line: 64, col: 15, offset: 3100},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 64, col: 15, offset: 3100},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 64, col: 15, offset: 3100},
							val:        "timestamp",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 64, col: 28, offset: 3113},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 64, col: 33, offset: 3118},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 64, col: 46, offset: 3131},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 64, col: 59, offset: 3144},
								expr: &choiceExpr{
									pos: position{line: 64, col: 61, offset: 3146},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 64, col: 61, offset: 3146},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 64, col: 70, offset: 3155},
											name: "WithoutTZ",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TimeT",
			pos:  position{line: 77, col: 1, offset: 3434},
			expr: &actionExpr{
				pos: position{line: 77, col: 10, offset: 3443},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 77, col: 10, offset: 3443},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 77, col: 10, offset: 3443},
							val:        "time",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 77, col: 18, offset: 3451},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 77, col: 23, offset: 3456},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 77, col: 36, offset: 3469},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 77, col: 49, offset: 3482},
								expr: &choiceExpr{
									pos: position{line: 77, col: 51, offset: 3484},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 77, col: 51, offset: 3484},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 77, col: 60, offset: 3493},
											name: "WithoutTZ",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SecPrecision",
			pos:  position{line: 90, col: 1, offset: 3764},
			expr: &actionExpr{
				pos: position{line: 90, col: 17, offset: 3780},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 90, col: 17, offset: 3780},
					expr: &seqExpr{
						pos: position{line: 90, col: 19, offset: 3782},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 90, col: 19, offset: 3782},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 90, col: 22, offset: 3785},
								val:        "[0-6]",
								ranges:     []rune{'0', '6'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "WithTZ",
			pos:  position{line: 97, col: 1, offset: 3913},
			expr: &actionExpr{
				pos: position{line: 97, col: 11, offset: 3923},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 97, col: 11, offset: 3923},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 97, col: 11, offset: 3923},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 97, col: 14, offset: 3926},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 97, col: 22, offset: 3934},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 97, col: 25, offset: 3937},
							val:        "time",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 97, col: 33, offset: 3945},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 97, col: 36, offset: 3948},
							val:        "zone",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 101, col: 1, offset: 3982},
			expr: &actionExpr{
				pos: position{line: 101, col: 14, offset: 3995},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 101, col: 14, offset: 3995},
					expr: &seqExpr{
						pos: position{line: 101, col: 16, offset: 3997},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 101, col: 16, offset: 3997},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 101, col: 19, offset: 4000},
								val:        "without",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 101, col: 30, offset: 4011},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 101, col: 33, offset: 4014},
								val:        "time",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 101, col: 41, offset: 4022},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 101, col: 44, offset: 4025},
								val:        "zone",
								ignoreCase: true,
							},
						},
					},
				},
			},
		},
		{
			name: "CharT",
			pos:  position{line: 105, col: 1, offset: 4063},
			expr: &actionExpr{
				pos: position{line: 105, col: 10, offset: 4072},
				run: (*parser).callonCharT1,
				expr: &seqExpr{
					pos: position{line: 105, col: 10, offset: 4072},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 105, col: 12, offset: 4074},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 105, col: 12, offset: 4074},
									val:        "character",
									ignoreCase: true,
								},
								&litMatcher{
									pos:        position{line: 105, col: 27, offset: 4089},
									val:        "char",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 105, col: 37, offset: 4099},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 105, col: 44, offset: 4106},
								expr: &seqExpr{
									pos: position{line: 105, col: 46, offset: 4108},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 105, col: 46, offset: 4108},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 105, col: 50, offset: 4112},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 105, col: 61, offset: 4123},
											val:        ")",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "VarcharT",
			pos:  position{line: 117, col: 1, offset: 4378},
			expr: &actionExpr{
				pos: position{line: 117, col: 13, offset: 4390},
				run: (*parser).callonVarcharT1,
				expr: &seqExpr{
					pos: position{line: 117, col: 13, offset: 4390},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 117, col: 15, offset: 4392},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 117, col: 17, offset: 4394},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 117, col: 17, offset: 4394},
											val:        "character",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 117, col: 30, offset: 4407},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 117, col: 33, offset: 4410},
											val:        "varying",
											ignoreCase: true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 117, col: 48, offset: 4425},
									val:        "varchar",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 117, col: 61, offset: 4438},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 117, col: 68, offset: 4445},
								expr: &seqExpr{
									pos: position{line: 117, col: 70, offset: 4447},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 117, col: 70, offset: 4447},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 117, col: 74, offset: 4451},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 117, col: 85, offset: 4462},
											val:        ")",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "BitT",
			pos:  position{line: 128, col: 1, offset: 4697},
			expr: &actionExpr{
				pos: position{line: 128, col: 9, offset: 4705},
				run: (*parser).callonBitT1,
				expr: &seqExpr{
					pos: position{line: 128, col: 9, offset: 4705},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 128, col: 9, offset: 4705},
							val:        "bit",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 128, col: 16, offset: 4712},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 128, col: 23, offset: 4719},
								expr: &seqExpr{
									pos: position{line: 128, col: 25, offset: 4721},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 128, col: 25, offset: 4721},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 128, col: 29, offset: 4725},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 128, col: 40, offset: 4736},
											val:        ")",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "BitVarT",
			pos:  position{line: 140, col: 1, offset: 4990},
			expr: &actionExpr{
				pos: position{line: 140, col: 12, offset: 5001},
				run: (*parser).callonBitVarT1,
				expr: &seqExpr{
					pos: position{line: 140, col: 12, offset: 5001},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 140, col: 12, offset: 5001},
							val:        "bit",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 140, col: 19, offset: 5008},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 140, col: 22, offset: 5011},
							val:        "varying",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 140, col: 33, offset: 5022},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 140, col: 40, offset: 5029},
								expr: &seqExpr{
									pos: position{line: 140, col: 42, offset: 5031},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 140, col: 42, offset: 5031},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 140, col: 46, offset: 5035},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 140, col: 57, offset: 5046},
											val:        ")",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "IntT",
			pos:  position{line: 151, col: 1, offset: 5280},
			expr: &actionExpr{
				pos: position{line: 151, col: 9, offset: 5288},
				run: (*parser).callonIntT1,
				expr: &choiceExpr{
					pos: position{line: 151, col: 11, offset: 5290},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 151, col: 11, offset: 5290},
							val:        "integer",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 151, col: 24, offset: 5303},
							val:        "int",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "PgOidT",
			pos:  position{line: 157, col: 1, offset: 5385},
			expr: &actionExpr{
				pos: position{line: 157, col: 11, offset: 5395},
				run: (*parser).callonPgOidT1,
				expr: &choiceExpr{
					pos: position{line: 157, col: 13, offset: 5397},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 157, col: 13, offset: 5397},
							val:        "oid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 22, offset: 5406},
							val:        "regprocedure",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 40, offset: 5424},
							val:        "regproc",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 53, offset: 5437},
							val:        "regoperator",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 70, offset: 5454},
							val:        "regoper",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 83, offset: 5467},
							val:        "regclass",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 97, offset: 5481},
							val:        "regtype",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 110, offset: 5494},
							val:        "regrole",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 123, offset: 5507},
							val:        "regnamespace",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 141, offset: 5525},
							val:        "regconfig",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 157, col: 156, offset: 5540},
							val:        "regdictionary",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 163, col: 1, offset: 5654},
			expr: &actionExpr{
				pos: position{line: 163, col: 11, offset: 5664},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 163, col: 13, offset: 5666},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 163, col: 13, offset: 5666},
							val:        "date",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 23, offset: 5676},
							val:        "smallint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 37, offset: 5690},
							val:        "bigint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 49, offset: 5702},
							val:        "decimal",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 62, offset: 5715},
							val:        "numeric",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 75, offset: 5728},
							val:        "real",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 85, offset: 5738},
							val:        "smallserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 102, offset: 5755},
							val:        "serial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 114, offset: 5767},
							val:        "bigserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 129, offset: 5782},
							val:        "boolean",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 142, offset: 5795},
							val:        "text",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 152, offset: 5805},
							val:        "money",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 163, offset: 5816},
							val:        "bytea",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 174, offset: 5827},
							val:        "point",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 185, offset: 5838},
							val:        "line",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 195, offset: 5848},
							val:        "lseg",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 205, offset: 5858},
							val:        "box",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 214, offset: 5867},
							val:        "path",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 224, offset: 5877},
							val:        "polygon",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 237, offset: 5890},
							val:        "circle",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 249, offset: 5902},
							val:        "cidr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 259, offset: 5912},
							val:        "inet",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 269, offset: 5922},
							val:        "macaddr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 282, offset: 5935},
							val:        "uuid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 292, offset: 5945},
							val:        "xml",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 301, offset: 5954},
							val:        "jsonb",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 163, col: 312, offset: 5965},
							val:        "json",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CustomT",
			pos:  position{line: 169, col: 1, offset: 6070},
			expr: &actionExpr{
				pos: position{line: 169, col: 13, offset: 6082},
				run: (*parser).callonCustomT1,
				expr: &ruleRefExpr{
					pos:  position{line: 169, col: 13, offset: 6082},
					name: "Ident",
				},
			},
		},
		{
			name: "CreateSeqStmt",
			pos:  position{line: 190, col: 1, offset: 7547},
			expr: &actionExpr{
				pos: position{line: 190, col: 18, offset: 7564},
				run: (*parser).callonCreateSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 190, col: 18, offset: 7564},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 190, col: 18, offset: 7564},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 190, col: 28, offset: 7574},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 190, col: 31, offset: 7577},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 190, col: 43, offset: 7589},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 190, col: 46, offset: 7592},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 190, col: 51, offset: 7597},
								name: "Ident",
							},
						},
						&labeledExpr{
							pos:   position{line: 190, col: 57, offset: 7603},
							label: "verses",
							expr: &zeroOrMoreExpr{
								pos: position{line: 190, col: 64, offset: 7610},
								expr: &ruleRefExpr{
									pos:  position{line: 190, col: 64, offset: 7610},
									name: "CreateSeqVerse",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 190, col: 80, offset: 7626},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 190, col: 82, offset: 7628},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 190, col: 86, offset: 7632},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "CreateSeqVerse",
			pos:  position{line: 204, col: 1, offset: 8025},
			expr: &actionExpr{
				pos: position{line: 204, col: 19, offset: 8043},
				run: (*parser).callonCreateSeqVerse1,
				expr: &labeledExpr{
					pos:   position{line: 204, col: 19, offset: 8043},
					label: "verse",
					expr: &choiceExpr{
						pos: position{line: 204, col: 27, offset: 8051},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 204, col: 27, offset: 8051},
								name: "IncrementBy",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 41, offset: 8065},
								name: "MinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 52, offset: 8076},
								name: "NoMinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 65, offset: 8089},
								name: "MaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 76, offset: 8100},
								name: "NoMaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 89, offset: 8113},
								name: "Start",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 97, offset: 8121},
								name: "Cache",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 105, offset: 8129},
								name: "Cycle",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 113, offset: 8137},
								name: "OwnedBy",
							},
						},
					},
				},
			},
		},
		{
			name: "IncrementBy",
			pos:  position{line: 208, col: 1, offset: 8174},
			expr: &actionExpr{
				pos: position{line: 208, col: 16, offset: 8189},
				run: (*parser).callonIncrementBy1,
				expr: &seqExpr{
					pos: position{line: 208, col: 16, offset: 8189},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 208, col: 16, offset: 8189},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 208, col: 19, offset: 8192},
							val:        "increment",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 208, col: 32, offset: 8205},
							expr: &seqExpr{
								pos: position{line: 208, col: 33, offset: 8206},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 208, col: 33, offset: 8206},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 208, col: 36, offset: 8209},
										val:        "by",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 208, col: 44, offset: 8217},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 208, col: 47, offset: 8220},
							label: "num",
							expr: &ruleRefExpr{
								pos:  position{line: 208, col: 51, offset: 8224},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "MinValue",
			pos:  position{line: 214, col: 1, offset: 8338},
			expr: &actionExpr{
				pos: position{line: 214, col: 13, offset: 8350},
				run: (*parser).callonMinValue1,
				expr: &seqExpr{
					pos: position{line: 214, col: 13, offset: 8350},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 214, col: 13, offset: 8350},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 214, col: 16, offset: 8353},
							val:        "minvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 214, col: 28, offset: 8365},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 214, col: 31, offset: 8368},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 214, col: 35, offset: 8372},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMinValue",
			pos:  position{line: 220, col: 1, offset: 8485},
			expr: &actionExpr{
				pos: position{line: 220, col: 15, offset: 8499},
				run: (*parser).callonNoMinValue1,
				expr: &seqExpr{
					pos: position{line: 220, col: 15, offset: 8499},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 220, col: 15, offset: 8499},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 220, col: 18, offset: 8502},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 220, col: 24, offset: 8508},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 220, col: 27, offset: 8511},
							val:        "minvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "MaxValue",
			pos:  position{line: 224, col: 1, offset: 8548},
			expr: &actionExpr{
				pos: position{line: 224, col: 13, offset: 8560},
				run: (*parser).callonMaxValue1,
				expr: &seqExpr{
					pos: position{line: 224, col: 13, offset: 8560},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 224, col: 13, offset: 8560},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 224, col: 16, offset: 8563},
							val:        "maxvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 224, col: 28, offset: 8575},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 224, col: 31, offset: 8578},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 224, col: 35, offset: 8582},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMaxValue",
			pos:  position{line: 230, col: 1, offset: 8695},
			expr: &actionExpr{
				pos: position{line: 230, col: 15, offset: 8709},
				run: (*parser).callonNoMaxValue1,
				expr: &seqExpr{
					pos: position{line: 230, col: 15, offset: 8709},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 230, col: 15, offset: 8709},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 230, col: 18, offset: 8712},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 230, col: 24, offset: 8718},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 230, col: 27, offset: 8721},
							val:        "maxvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "Start",
			pos:  position{line: 234, col: 1, offset: 8758},
			expr: &actionExpr{
				pos: position{line: 234, col: 10, offset: 8767},
				run: (*parser).callonStart1,
				expr: &seqExpr{
					pos: position{line: 234, col: 10, offset: 8767},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 234, col: 10, offset: 8767},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 234, col: 13, offset: 8770},
							val:        "start",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 234, col: 22, offset: 8779},
							expr: &seqExpr{
								pos: position{line: 234, col: 23, offset: 8780},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 234, col: 23, offset: 8780},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 234, col: 26, offset: 8783},
										val:        "with",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 36, offset: 8793},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 234, col: 39, offset: 8796},
							label: "start",
							expr: &ruleRefExpr{
								pos:  position{line: 234, col: 45, offset: 8802},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cache",
			pos:  position{line: 240, col: 1, offset: 8914},
			expr: &actionExpr{
				pos: position{line: 240, col: 10, offset: 8923},
				run: (*parser).callonCache1,
				expr: &seqExpr{
					pos: position{line: 240, col: 10, offset: 8923},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 240, col: 10, offset: 8923},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 240, col: 13, offset: 8926},
							val:        "cache",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 22, offset: 8935},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 240, col: 25, offset: 8938},
							label: "cache",
							expr: &ruleRefExpr{
								pos:  position{line: 240, col: 31, offset: 8944},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cycle",
			pos:  position{line: 246, col: 1, offset: 9056},
			expr: &actionExpr{
				pos: position{line: 246, col: 10, offset: 9065},
				run: (*parser).callonCycle1,
				expr: &seqExpr{
					pos: position{line: 246, col: 10, offset: 9065},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 246, col: 10, offset: 9065},
							label: "no",
							expr: &zeroOrOneExpr{
								pos: position{line: 246, col: 13, offset: 9068},
								expr: &seqExpr{
									pos: position{line: 246, col: 14, offset: 9069},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 246, col: 14, offset: 9069},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 246, col: 17, offset: 9072},
											val:        "no",
											ignoreCase: true,
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 246, col: 25, offset: 9080},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 246, col: 28, offset: 9083},
							val:        "cycle",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OwnedBy",
			pos:  position{line: 257, col: 1, offset: 9267},
			expr: &actionExpr{
				pos: position{line: 257, col: 12, offset: 9278},
				run: (*parser).callonOwnedBy1,
				expr: &seqExpr{
					pos: position{line: 257, col: 12, offset: 9278},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 257, col: 12, offset: 9278},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 257, col: 15, offset: 9281},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 257, col: 24, offset: 9290},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 257, col: 27, offset: 9293},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 257, col: 33, offset: 9299},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 257, col: 36, offset: 9302},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 257, col: 43, offset: 9309},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 257, col: 43, offset: 9309},
										val:        "none",
										ignoreCase: true,
									},
									&seqExpr{
										pos: position{line: 257, col: 55, offset: 9321},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 257, col: 55, offset: 9321},
												name: "Ident",
											},
											&litMatcher{
												pos:        position{line: 257, col: 61, offset: 9327},
												val:        ".",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 257, col: 65, offset: 9331},
												name: "Ident",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CreateTypeStmt",
			pos:  position{line: 278, col: 2, offset: 10977},
			expr: &actionExpr{
				pos: position{line: 278, col: 20, offset: 10995},
				run: (*parser).callonCreateTypeStmt1,
				expr: &seqExpr{
					pos: position{line: 278, col: 20, offset: 10995},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 278, col: 20, offset: 10995},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 278, col: 30, offset: 11005},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 278, col: 33, offset: 11008},
							val:        "type",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 278, col: 41, offset: 11016},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 278, col: 44, offset: 11019},
							label: "typename",
							expr: &ruleRefExpr{
								pos:  position{line: 278, col: 53, offset: 11028},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 278, col: 59, offset: 11034},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 278, col: 62, offset: 11037},
							val:        "as",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 278, col: 68, offset: 11043},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 278, col: 71, offset: 11046},
							label: "typedef",
							expr: &ruleRefExpr{
								pos:  position{line: 278, col: 79, offset: 11054},
								name: "EnumDef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 278, col: 87, offset: 11062},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 278, col: 89, offset: 11064},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 278, col: 93, offset: 11068},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "EnumDef",
			pos:  position{line: 284, col: 2, offset: 11185},
			expr: &actionExpr{
				pos: position{line: 284, col: 13, offset: 11196},
				run: (*parser).callonEnumDef1,
				expr: &seqExpr{
					pos: position{line: 284, col: 13, offset: 11196},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 284, col: 13, offset: 11196},
							val:        "ENUM",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 284, col: 20, offset: 11203},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 284, col: 22, offset: 11205},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 284, col: 26, offset: 11209},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 284, col: 28, offset: 11211},
							label: "vals",
							expr: &seqExpr{
								pos: position{line: 284, col: 35, offset: 11218},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 284, col: 35, offset: 11218},
										name: "StringConst",
									},
									&zeroOrMoreExpr{
										pos: position{line: 284, col: 47, offset: 11230},
										expr: &seqExpr{
											pos: position{line: 284, col: 49, offset: 11232},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 284, col: 49, offset: 11232},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 284, col: 51, offset: 11234},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 284, col: 55, offset: 11238},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 284, col: 57, offset: 11240},
													name: "StringConst",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 284, col: 75, offset: 11258},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 284, col: 77, offset: 11260},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AlterTableStmt",
			pos:  position{line: 309, col: 1, offset: 12891},
			expr: &actionExpr{
				pos: position{line: 309, col: 19, offset: 12909},
				run: (*parser).callonAlterTableStmt1,
				expr: &seqExpr{
					pos: position{line: 309, col: 19, offset: 12909},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 309, col: 19, offset: 12909},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 309, col: 28, offset: 12918},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 309, col: 31, offset: 12921},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 309, col: 40, offset: 12930},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 309, col: 43, offset: 12933},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 309, col: 48, offset: 12938},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 309, col: 54, offset: 12944},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 309, col: 57, offset: 12947},
							val:        "owner",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 309, col: 66, offset: 12956},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 309, col: 69, offset: 12959},
							val:        "to",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 309, col: 75, offset: 12965},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 309, col: 78, offset: 12968},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 309, col: 84, offset: 12974},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 309, col: 90, offset: 12980},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 309, col: 92, offset: 12982},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 309, col: 96, offset: 12986},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 323, col: 1, offset: 14305},
			expr: &actionExpr{
				pos: position{line: 323, col: 25, offset: 14329},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 323, col: 25, offset: 14329},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 323, col: 25, offset: 14329},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 323, col: 36, offset: 14340},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 323, col: 39, offset: 14343},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 323, col: 45, offset: 14349},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 323, col: 48, offset: 14352},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 323, col: 61, offset: 14365},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 323, col: 63, offset: 14367},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 323, col: 73, offset: 14377},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 323, col: 79, offset: 14383},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 323, col: 81, offset: 14385},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 323, col: 87, offset: 14391},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 323, col: 89, offset: 14393},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 323, col: 97, offset: 14401},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 323, col: 109, offset: 14413},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 323, col: 111, offset: 14415},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 323, col: 115, offset: 14419},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 327, col: 1, offset: 14507},
			expr: &actionExpr{
				pos: position{line: 327, col: 24, offset: 14530},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 327, col: 24, offset: 14530},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 327, col: 24, offset: 14530},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 327, col: 34, offset: 14540},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 327, col: 37, offset: 14543},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 327, col: 50, offset: 14556},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 327, col: 53, offset: 14559},
							expr: &seqExpr{
								pos: position{line: 327, col: 55, offset: 14561},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 327, col: 55, offset: 14561},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 327, col: 61, offset: 14567},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 327, col: 64, offset: 14570},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 327, col: 71, offset: 14577},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 327, col: 74, offset: 14580},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 327, col: 84, offset: 14590},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 327, col: 90, offset: 14596},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 327, col: 100, offset: 14606},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 327, col: 106, offset: 14612},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 327, col: 109, offset: 14615},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 327, col: 117, offset: 14623},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 327, col: 120, offset: 14626},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 327, col: 130, offset: 14636},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 327, col: 133, offset: 14639},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 327, col: 140, offset: 14646},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 327, col: 146, offset: 14652},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 327, col: 148, offset: 14654},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 327, col: 152, offset: 14658},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 331, col: 1, offset: 14748},
			expr: &actionExpr{
				pos: position{line: 331, col: 12, offset: 14759},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 331, col: 12, offset: 14759},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 331, col: 12, offset: 14759},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 19, offset: 14766},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 331, col: 21, offset: 14768},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 331, col: 25, offset: 14772},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 29, offset: 14776},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 331, col: 33, offset: 14780},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 331, col: 33, offset: 14780},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 331, col: 39, offset: 14786},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 47, offset: 14794},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 331, col: 49, offset: 14796},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 331, col: 56, offset: 14803},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 77, offset: 14824},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 331, col: 79, offset: 14826},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 83, offset: 14830},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 336, col: 1, offset: 14911},
			expr: &actionExpr{
				pos: position{line: 336, col: 8, offset: 14918},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 336, col: 8, offset: 14918},
					expr: &charClassMatcher{
						pos:        position{line: 336, col: 8, offset: 14918},
						val:        "[a-z_]i",
						chars:      []rune{'_'},
						ranges:     []rune{'a', 'z'},
						ignoreCase: true,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "CommaSeparatedValues",
			pos:  position{line: 351, col: 1, offset: 15758},
			expr: &actionExpr{
				pos: position{line: 351, col: 25, offset: 15782},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 351, col: 25, offset: 15782},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 351, col: 32, offset: 15789},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 351, col: 32, offset: 15789},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 351, col: 38, offset: 15795},
								expr: &seqExpr{
									pos: position{line: 351, col: 40, offset: 15797},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 351, col: 40, offset: 15797},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 351, col: 42, offset: 15799},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 351, col: 46, offset: 15803},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 351, col: 48, offset: 15805},
											name: "Value",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 363, col: 1, offset: 16095},
			expr: &choiceExpr{
				pos: position{line: 363, col: 12, offset: 16106},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 363, col: 12, offset: 16106},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 363, col: 21, offset: 16115},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 363, col: 31, offset: 16125},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 363, col: 45, offset: 16139},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 365, col: 1, offset: 16148},
			expr: &actionExpr{
				pos: position{line: 365, col: 16, offset: 16163},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 365, col: 16, offset: 16163},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 365, col: 16, offset: 16163},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 365, col: 20, offset: 16167},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 365, col: 26, offset: 16173},
								expr: &choiceExpr{
									pos: position{line: 365, col: 27, offset: 16174},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 365, col: 27, offset: 16174},
											val:        "[^'\\n]",
											chars:      []rune{'\'', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 365, col: 36, offset: 16183},
											val:        "''",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 365, col: 43, offset: 16190},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 369, col: 1, offset: 16242},
			expr: &actionExpr{
				pos: position{line: 369, col: 10, offset: 16251},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 369, col: 10, offset: 16251},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 369, col: 10, offset: 16251},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 369, col: 18, offset: 16259},
							expr: &charClassMatcher{
								pos:        position{line: 369, col: 18, offset: 16259},
								val:        "[a-z_0-9$]i",
								chars:      []rune{'_', '$'},
								ranges:     []rune{'a', 'z', '0', '9'},
								ignoreCase: true,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Number",
			pos:  position{line: 373, col: 1, offset: 16312},
			expr: &actionExpr{
				pos: position{line: 373, col: 11, offset: 16322},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 373, col: 13, offset: 16324},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 373, col: 13, offset: 16324},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 373, col: 19, offset: 16330},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 373, col: 19, offset: 16330},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 373, col: 24, offset: 16335},
									expr: &charClassMatcher{
										pos:        position{line: 373, col: 24, offset: 16335},
										val:        "[0-9]",
										ranges:     []rune{'0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NonZNumber",
			pos:  position{line: 378, col: 1, offset: 16430},
			expr: &actionExpr{
				pos: position{line: 378, col: 15, offset: 16444},
				run: (*parser).callonNonZNumber1,
				expr: &seqExpr{
					pos: position{line: 378, col: 15, offset: 16444},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 378, col: 15, offset: 16444},
							val:        "[1-9]",
							ranges:     []rune{'1', '9'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 378, col: 20, offset: 16449},
							expr: &charClassMatcher{
								pos:        position{line: 378, col: 20, offset: 16449},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Boolean",
			pos:  position{line: 383, col: 1, offset: 16542},
			expr: &actionExpr{
				pos: position{line: 383, col: 12, offset: 16553},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 383, col: 12, offset: 16553},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 383, col: 20, offset: 16561},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 383, col: 20, offset: 16561},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 383, col: 34, offset: 16575},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 387, col: 1, offset: 16617},
			expr: &actionExpr{
				pos: position{line: 387, col: 16, offset: 16632},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 387, col: 18, offset: 16634},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 387, col: 18, offset: 16634},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 387, col: 27, offset: 16643},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 387, col: 27, offset: 16643},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 387, col: 31, offset: 16647},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 387, col: 49, offset: 16665},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 387, col: 55, offset: 16671},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 391, col: 1, offset: 16717},
			expr: &choiceExpr{
				pos: position{line: 391, col: 24, offset: 16740},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 391, col: 24, offset: 16740},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 391, col: 33, offset: 16749},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 391, col: 41, offset: 16757},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 391, col: 48, offset: 16764},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 391, col: 54, offset: 16770},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 393, col: 1, offset: 16777},
			expr: &actionExpr{
				pos: position{line: 393, col: 17, offset: 16793},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 393, col: 19, offset: 16795},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 393, col: 19, offset: 16795},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 393, col: 29, offset: 16805},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 393, col: 29, offset: 16805},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 393, col: 33, offset: 16809},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 393, col: 52, offset: 16828},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 393, col: 58, offset: 16834},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 397, col: 1, offset: 16882},
			expr: &choiceExpr{
				pos: position{line: 397, col: 25, offset: 16906},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 397, col: 25, offset: 16906},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 397, col: 35, offset: 16916},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 397, col: 42, offset: 16923},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 397, col: 50, offset: 16931},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 397, col: 56, offset: 16937},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 409, col: 1, offset: 17452},
			expr: &actionExpr{
				pos: position{line: 409, col: 12, offset: 17463},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 409, col: 14, offset: 17465},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 409, col: 14, offset: 17465},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 409, col: 34, offset: 17485},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 413, col: 1, offset: 17529},
			expr: &seqExpr{
				pos: position{line: 413, col: 21, offset: 17549},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 413, col: 21, offset: 17549},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 413, col: 26, offset: 17554},
						expr: &seqExpr{
							pos: position{line: 413, col: 28, offset: 17556},
							exprs: []interface{}{
								&zeroOrMoreExpr{
									pos: position{line: 413, col: 28, offset: 17556},
									expr: &anyMatcher{
										line: 413, col: 28, offset: 17556,
									},
								},
								&zeroOrOneExpr{
									pos: position{line: 413, col: 31, offset: 17559},
									expr: &ruleRefExpr{
										pos:  position{line: 413, col: 31, offset: 17559},
										name: "MultilineComment",
									},
								},
								&zeroOrMoreExpr{
									pos: position{line: 413, col: 49, offset: 17577},
									expr: &anyMatcher{
										line: 413, col: 49, offset: 17577,
									},
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 413, col: 55, offset: 17583},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 413, col: 60, offset: 17588},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 415, col: 1, offset: 17593},
			expr: &seqExpr{
				pos: position{line: 415, col: 22, offset: 17614},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 415, col: 22, offset: 17614},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 415, col: 27, offset: 17619},
						expr: &charClassMatcher{
							pos:        position{line: 415, col: 27, offset: 17619},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 415, col: 36, offset: 17628},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 417, col: 1, offset: 17633},
			expr: &seqExpr{
				pos: position{line: 417, col: 9, offset: 17641},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 417, col: 9, offset: 17641},
						expr: &charClassMatcher{
							pos:        position{line: 417, col: 9, offset: 17641},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 417, col: 17, offset: 17649},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 417, col: 17, offset: 17649},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 417, col: 26, offset: 17658},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 417, col: 35, offset: 17667},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 417, col: 42, offset: 17674},
								val:        "\n",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 417, col: 49, offset: 17681},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 419, col: 1, offset: 17687},
			expr: &zeroOrMoreExpr{
				pos: position{line: 419, col: 19, offset: 17705},
				expr: &charClassMatcher{
					pos:        position{line: 419, col: 19, offset: 17705},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name:        "_1",
			displayName: "\"at least 1 whitespace\"",
			pos:         position{line: 421, col: 1, offset: 17717},
			expr: &oneOrMoreExpr{
				pos: position{line: 421, col: 31, offset: 17747},
				expr: &charClassMatcher{
					pos:        position{line: 421, col: 31, offset: 17747},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 423, col: 1, offset: 17759},
			expr: &notExpr{
				pos: position{line: 423, col: 8, offset: 17766},
				expr: &anyMatcher{
					line: 423, col: 9, offset: 17767,
				},
			},
		},
	},
}

func (c *current) onSQL1(stmts interface{}) (interface{}, error) {
	return stmts, nil
}

func (p *parser) callonSQL1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSQL1(stack["stmts"])
}

func (c *current) onStmt1(stmt interface{}) (interface{}, error) {
	return stmt, nil
}

func (p *parser) callonStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStmt1(stack["stmt"])
}

func (c *current) onCreateTableStmt1(tablename, fields interface{}) (interface{}, error) {
	fieldsSlice := []map[string]string{}
	valsSlice := toIfaceSlice(fields)
	if valsSlice[0] == nil {
		fieldsSlice = append(fieldsSlice, nil)
	} else {
		fieldsSlice = append(fieldsSlice, valsSlice[0].(map[string]string))
	}
	restSlice := toIfaceSlice(valsSlice[1])
	for _, v := range restSlice {
		vSlice := toIfaceSlice(v)
		if vSlice[3] == nil {
			fieldsSlice = append(fieldsSlice, nil)
		} else {
			fieldsSlice = append(fieldsSlice, vSlice[3].(map[string]string))
		}
	}
	return parseCreateTableStmt(tablename, fieldsSlice)
}

func (p *parser) callonCreateTableStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCreateTableStmt1(stack["tablename"], stack["fields"])
}

func (c *current) onFieldDef1(name, dataType, notnull interface{}) (interface{}, error) {
	if dataType == nil {
		return nil, nil
	}
	result := dataType.(map[string]string)
	result["name"] = interfaceToString(name)
	if notnull != nil {
		result["not_null"] = "true"
	}
	return result, nil
}

func (p *parser) callonFieldDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldDef1(stack["name"], stack["dataType"], stack["notnull"])
}

func (c *current) onDataType1(t interface{}) (interface{}, error) {
	return t, nil
}

func (p *parser) callonDataType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDataType1(stack["t"])
}

func (c *current) onTimestampT1(prec, withTimeZone interface{}) (interface{}, error) {
	var result = make(map[string]string)
	if withTimeZone.(bool) {
		result["type"] = "datetimetz"
	} else {
		result["type"] = "datetime"
	}
	if prec != nil {
		result["sec_precision"] = prec.(string)
	}
	return result, nil
}

func (p *parser) callonTimestampT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimestampT1(stack["prec"], stack["withTimeZone"])
}

func (c *current) onTimeT1(prec, withTimeZone interface{}) (interface{}, error) {
	var result = make(map[string]string)
	if withTimeZone.(bool) {
		result["type"] = "timetz"
	} else {
		result["type"] = "time"
	}
	if prec != nil {
		result["sec_precision"] = prec.(string)
	}
	return result, nil
}

func (p *parser) callonTimeT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimeT1(stack["prec"], stack["withTimeZone"])
}

func (c *current) onSecPrecision1() (interface{}, error) {
	if len(c.text) > 0 {
		return strings.TrimLeft(string(c.text), " \r\t\n"), nil
	}
	return nil, nil
}

func (p *parser) callonSecPrecision1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSecPrecision1()
}

func (c *current) onWithTZ1() (interface{}, error) {
	return true, nil
}

func (p *parser) callonWithTZ1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onWithTZ1()
}

func (c *current) onWithoutTZ1() (interface{}, error) {
	return false, nil
}

func (p *parser) callonWithoutTZ1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onWithoutTZ1()
}

func (c *current) onCharT1(length interface{}) (interface{}, error) {
	result := map[string]string{
		"type":   "char",
		"length": "1",
	}
	if length != nil {
		slice := toIfaceSlice(length)
		result["length"] = strconv.FormatInt(slice[1].(int64), 10)
	}
	return result, nil
}

func (p *parser) callonCharT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharT1(stack["length"])
}

func (c *current) onVarcharT1(length interface{}) (interface{}, error) {
	result := map[string]string{
		"type": "varchar",
	}
	if length != nil {
		slice := toIfaceSlice(length)
		result["length"] = strconv.FormatInt(slice[1].(int64), 10)
	}
	return result, nil
}

func (p *parser) callonVarcharT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVarcharT1(stack["length"])
}

func (c *current) onBitT1(length interface{}) (interface{}, error) {
	result := map[string]string{
		"type":   "bit",
		"length": "1",
	}
	if length != nil {
		slice := toIfaceSlice(length)
		result["length"] = strconv.FormatInt(slice[1].(int64), 10)
	}
	return result, nil
}

func (p *parser) callonBitT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBitT1(stack["length"])
}

func (c *current) onBitVarT1(length interface{}) (interface{}, error) {
	result := map[string]string{
		"type": "bitvar",
	}
	if length != nil {
		slice := toIfaceSlice(length)
		result["length"] = strconv.FormatInt(slice[1].(int64), 10)
	}
	return result, nil
}

func (p *parser) callonBitVarT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBitVarT1(stack["length"])
}

func (c *current) onIntT1() (interface{}, error) {
	return map[string]string{
		"type": "integer",
	}, nil
}

func (p *parser) callonIntT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntT1()
}

func (c *current) onPgOidT1() (interface{}, error) {
	return map[string]string{
		"type": strings.ToLower(string(c.text)),
	}, nil
}

func (p *parser) callonPgOidT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPgOidT1()
}

func (c *current) onOtherT1() (interface{}, error) {
	return map[string]string{
		"type": strings.ToLower(string(c.text)),
	}, nil
}

func (p *parser) callonOtherT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOtherT1()
}

func (c *current) onCustomT1() (interface{}, error) {
	typeName := strings.ToLower(string(c.text))
	err := typeExists(typeName)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"type": strings.ToLower(string(c.text)),
	}, nil
}

func (p *parser) callonCustomT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCustomT1()
}

func (c *current) onCreateSeqStmt1(name, verses interface{}) (interface{}, error) {
	versesSlice := toIfaceSlice(verses)
	properties := make(map[string]string)
	for _, verse := range versesSlice {
		if verse == nil {
			continue
		}
		if err := mergo.Merge(&properties, verse.(map[string]string), mergo.WithOverride); err != nil {
			return nil, err
		}
	}
	return parseCreateSeq(name.(Identifier), properties)
}

func (p *parser) callonCreateSeqStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCreateSeqStmt1(stack["name"], stack["verses"])
}

func (c *current) onCreateSeqVerse1(verse interface{}) (interface{}, error) {
	return verse, nil
}

func (p *parser) callonCreateSeqVerse1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCreateSeqVerse1(stack["verse"])
}

func (c *current) onIncrementBy1(num interface{}) (interface{}, error) {
	return map[string]string{
		"increment": strconv.FormatInt(num.(int64), 10),
	}, nil
}

func (p *parser) callonIncrementBy1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIncrementBy1(stack["num"])
}

func (c *current) onMinValue1(val interface{}) (interface{}, error) {
	return map[string]string{
		"minvalue": strconv.FormatInt(val.(int64), 10),
	}, nil
}

func (p *parser) callonMinValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMinValue1(stack["val"])
}

func (c *current) onNoMinValue1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonNoMinValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNoMinValue1()
}

func (c *current) onMaxValue1(val interface{}) (interface{}, error) {
	return map[string]string{
		"maxvalue": strconv.FormatInt(val.(int64), 10),
	}, nil
}

func (p *parser) callonMaxValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMaxValue1(stack["val"])
}

func (c *current) onNoMaxValue1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonNoMaxValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNoMaxValue1()
}

func (c *current) onStart1(start interface{}) (interface{}, error) {
	return map[string]string{
		"start": strconv.FormatInt(start.(int64), 10),
	}, nil
}

func (p *parser) callonStart1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStart1(stack["start"])
}

func (c *current) onCache1(cache interface{}) (interface{}, error) {
	return map[string]string{
		"cache": strconv.FormatInt(cache.(int64), 10),
	}, nil
}

func (p *parser) callonCache1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCache1(stack["cache"])
}

func (c *current) onCycle1(no interface{}) (interface{}, error) {
	if no != nil {
		return map[string]string{
			"cycle": "false",
		}, nil
	}
	return map[string]string{
		"cycle": "true",
	}, nil
}

func (p *parser) callonCycle1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCycle1(stack["no"])
}

func (c *current) onOwnedBy1(name interface{}) (interface{}, error) {
	if _, ok := name.([]byte); ok {
		return nil, nil
	}
	nameSlice := toIfaceSlice(name)
	table := strings.ToLower(interfaceToString(nameSlice[0]))
	column := strings.ToLower(interfaceToString(nameSlice[2]))
	return map[string]string{
		"owned_by": strings.Join([]string{table, column}, "/"),
	}, nil
}

func (p *parser) callonOwnedBy1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOwnedBy1(stack["name"])
}

func (c *current) onCreateTypeStmt1(typename, typedef interface{}) (interface{}, error) {
	enum := typedef.(Enum)
	enum.Name = typename.(Identifier)
	return parseCreateTypeEnumStmt(enum)

}

func (p *parser) callonCreateTypeStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCreateTypeStmt1(stack["typename"], stack["typedef"])
}

func (c *current) onEnumDef1(vals interface{}) (interface{}, error) {
	labels := []String{}
	valsSlice := toIfaceSlice(vals)
	labels = append(labels, valsSlice[0].(String))
	restSlice := toIfaceSlice(valsSlice[1])
	for _, v := range restSlice {
		vSlice := toIfaceSlice(v)
		labels = append(labels, vSlice[3].(String))
	}
	return Enum{
		Name:   "",
		Labels: labels,
	}, nil

}

func (p *parser) callonEnumDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumDef1(stack["vals"])
}

func (c *current) onAlterTableStmt1(name, owner interface{}) (interface{}, error) {
	return parseAlterTableStmt(name.(Identifier), owner.(Identifier))
}

func (p *parser) callonAlterTableStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAlterTableStmt1(stack["name"], stack["owner"])
}

func (c *current) onCommentExtensionStmt1(extension, comment interface{}) (interface{}, error) {
	return parseCommentExtensionStmt(extension.(Identifier), comment.(String))
}

func (p *parser) callonCommentExtensionStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCommentExtensionStmt1(stack["extension"], stack["comment"])
}

func (c *current) onCreateExtensionStmt1(extension, schema interface{}) (interface{}, error) {
	return parseCreateExtensionStmt(extension.(Identifier), schema.(Identifier))
}

func (p *parser) callonCreateExtensionStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCreateExtensionStmt1(stack["extension"], stack["schema"])
}

func (c *current) onSetStmt1(key, values interface{}) (interface{}, error) {
	setSettings(key.(string), toIfaceSlice(values))
	return nil, nil
}

func (p *parser) callonSetStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSetStmt1(stack["key"], stack["values"])
}

func (c *current) onKey1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonKey1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onKey1()
}

func (c *current) onCommaSeparatedValues1(vals interface{}) (interface{}, error) {
	res := []interface{}{}
	valsSlice := toIfaceSlice(vals)
	res = append(res, valsSlice[0])
	restSlice := toIfaceSlice(valsSlice[1])
	for _, v := range restSlice {
		vSlice := toIfaceSlice(v)
		res = append(res, vSlice[3])
	}
	return res, nil
}

func (p *parser) callonCommaSeparatedValues1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCommaSeparatedValues1(stack["vals"])
}

func (c *current) onStringConst1(value interface{}) (interface{}, error) {
	return String(toByteSlice(value)), nil
}

func (p *parser) callonStringConst1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringConst1(stack["value"])
}

func (c *current) onIdent1() (interface{}, error) {
	return Identifier(c.text), nil
}

func (p *parser) callonIdent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdent1()
}

func (c *current) onNumber1() (interface{}, error) {
	number, _ := strconv.ParseInt(string(c.text), 10, 64)
	return number, nil
}

func (p *parser) callonNumber1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumber1()
}

func (c *current) onNonZNumber1() (interface{}, error) {
	number, _ := strconv.ParseInt(string(c.text), 10, 64)
	return number, nil
}

func (p *parser) callonNonZNumber1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonZNumber1()
}

func (c *current) onBoolean1(value interface{}) (interface{}, error) {
	return value, nil
}

func (p *parser) callonBoolean1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolean1(stack["value"])
}

func (c *current) onBooleanTrue1() (interface{}, error) {
	return true, nil
}

func (p *parser) callonBooleanTrue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBooleanTrue1()
}

func (c *current) onBooleanFalse1() (interface{}, error) {
	return false, nil
}

func (p *parser) callonBooleanFalse1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBooleanFalse1()
}

func (c *current) onComment1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonComment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onComment1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEntrypoint is returned when the specified entrypoint rule
	// does not exit.
	errInvalidEntrypoint = errors.New("invalid entrypoint")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Entrypoint creates an Option to set the rule name to use as entrypoint.
// The rule name must have been specified in the -alternate-entrypoints
// if generating the parser with the -optimize-grammar flag, otherwise
// it may have been optimized out. Passing an empty string sets the
// entrypoint to the first rule in the grammar.
//
// The default is to start parsing at the first rule in the grammar.
func Entrypoint(ruleName string) Option {
	return func(p *parser) Option {
		oldEntrypoint := p.entrypoint
		p.entrypoint = ruleName
		if ruleName == "" {
			p.entrypoint = g.rules[0].name
		}
		return Entrypoint(oldEntrypoint)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// AllowInvalidUTF8 creates an Option to allow invalid UTF-8 bytes.
// Every invalid UTF-8 byte is treated as a utf8.RuneError (U+FFFD)
// by character class matchers and is matched by the any matcher.
// The returned matched value, c.text and c.offset are NOT affected.
//
// The default is false.
func AllowInvalidUTF8(b bool) Option {
	return func(p *parser) Option {
		old := p.allowInvalidUTF8
		p.allowInvalidUTF8 = b
		return AllowInvalidUTF8(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// InitState creates an Option to set a key to a certain value in
// the global "state" store.
func InitState(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.state[key]
		p.cur.state[key] = value
		return InitState(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// state is a store for arbitrary key,value pairs that the user wants to be
	// tied to the backtracking of the parser.
	// This is always rolled back if a parsing rule fails.
	state storeDict

	// globalStore is a general store for the user to store arbitrary key-value
	// pairs that they need to manage and that they do not want tied to the
	// backtracking of the parser. This is only modified by the user and never
	// rolled back by the parser. It is always up to the user to keep this in a
	// consistent state.
	globalStore storeDict
}

type storeDict map[string]interface{}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type stateCodeExpr struct {
	pos position
	run func(*parser) error
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			state:       make(storeDict),
			globalStore: make(storeDict),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
		// start rule is rule [0] unless an alternate entrypoint is specified
		entrypoint: g.rules[0].name,
		emptyState: make(storeDict),
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64
	// entrypoint for the parser
	entrypoint string

	allowInvalidUTF8 bool

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}

	// emptyState contains an empty storeDict, which is used to optimize cloneState if global "state" store is not used.
	emptyState storeDict
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError && n == 1 { // see utf8.DecodeRune
		if !p.allowInvalidUTF8 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// Cloner is implemented by any value that has a Clone method, which returns a
// copy of the value. This is mainly used for types which are not passed by
// value (e.g map, slice, chan) or structs that contain such types.
//
// This is used in conjunction with the global state feature to create proper
// copies of the state to allow the parser to properly restore the state in
// the case of backtracking.
type Cloner interface {
	Clone() interface{}
}

// clone and return parser current state.
func (p *parser) cloneState() storeDict {
	if p.debug {
		defer p.out(p.in("cloneState"))
	}

	if len(p.cur.state) == 0 {
		if len(p.emptyState) > 0 {
			p.emptyState = make(storeDict)
		}
		return p.emptyState
	}

	state := make(storeDict, len(p.cur.state))
	for k, v := range p.cur.state {
		if c, ok := v.(Cloner); ok {
			state[k] = c.Clone()
		} else {
			state[k] = v
		}
	}
	return state
}

// restore parser current state to the state storeDict.
// every restoreState should applied only one time for every cloned state
func (p *parser) restoreState(state storeDict) {
	if p.debug {
		defer p.out(p.in("restoreState"))
	}
	p.cur.state = state
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	startRule, ok := p.rules[p.entrypoint]
	if !ok {
		p.addErr(errInvalidEntrypoint)
		return nil, p.errs.err()
	}

	p.read() // advance to first rune
	val, ok = p.parseRule(startRule)
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}

		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *stateCodeExpr:
		val, ok = p.parseStateCodeExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		state := p.cloneState()
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		p.restoreState(state)

		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	state := p.cloneState()

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		// EOF - see utf8.DecodeRune
		p.failAt(false, p.pt.position, ".")
		return nil, false
	}
	start := p.pt
	p.read()
	p.failAt(true, start.position, ".")
	return p.sliceFrom(start), true
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError && p.pt.w == 0 { // see utf8.DecodeRune
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		state := p.cloneState()

		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
		p.restoreState(state)
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	state := p.cloneState()

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	state := p.cloneState()
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restoreState(state)
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseStateCodeExpr(state *stateCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseStateCodeExpr"))
	}

	err := state.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, true
}

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
