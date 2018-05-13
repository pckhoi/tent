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
)

var g = &grammar{
	rules: []*rule{
		{
			name: "SQL",
			pos:  position{line: 6, col: 1, offset: 96},
			expr: &actionExpr{
				pos: position{line: 6, col: 8, offset: 103},
				run: (*parser).callonSQL1,
				expr: &labeledExpr{
					pos:   position{line: 6, col: 8, offset: 103},
					label: "stmts",
					expr: &oneOrMoreExpr{
						pos: position{line: 6, col: 14, offset: 109},
						expr: &ruleRefExpr{
							pos:  position{line: 6, col: 14, offset: 109},
							name: "Stmt",
						},
					},
				},
			},
		},
		{
			name: "Stmt",
			pos:  position{line: 10, col: 1, offset: 142},
			expr: &actionExpr{
				pos: position{line: 10, col: 9, offset: 150},
				run: (*parser).callonStmt1,
				expr: &seqExpr{
					pos: position{line: 10, col: 9, offset: 150},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 10, col: 9, offset: 150},
							expr: &ruleRefExpr{
								pos:  position{line: 10, col: 9, offset: 150},
								name: "Comment",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 10, col: 18, offset: 159},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 10, col: 20, offset: 161},
							label: "stmt",
							expr: &choiceExpr{
								pos: position{line: 10, col: 27, offset: 168},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 10, col: 27, offset: 168},
										name: "SetStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 10, col: 37, offset: 178},
										name: "CreateTableStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 10, col: 55, offset: 196},
										name: "CreateExtensionStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 10, col: 77, offset: 218},
										name: "CreateTypeStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 10, col: 94, offset: 235},
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
			pos:  position{line: 24, col: 1, offset: 1812},
			expr: &actionExpr{
				pos: position{line: 24, col: 20, offset: 1831},
				run: (*parser).callonCreateTableStmt1,
				expr: &seqExpr{
					pos: position{line: 24, col: 20, offset: 1831},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 24, col: 20, offset: 1831},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 30, offset: 1841},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 24, col: 33, offset: 1844},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 42, offset: 1853},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 24, col: 45, offset: 1856},
							label: "tablename",
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 55, offset: 1866},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 61, offset: 1872},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 24, col: 63, offset: 1874},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 67, offset: 1878},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 24, col: 69, offset: 1880},
							label: "fields",
							expr: &seqExpr{
								pos: position{line: 24, col: 78, offset: 1889},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 24, col: 78, offset: 1889},
										name: "FieldDef",
									},
									&zeroOrMoreExpr{
										pos: position{line: 24, col: 87, offset: 1898},
										expr: &seqExpr{
											pos: position{line: 24, col: 89, offset: 1900},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 24, col: 89, offset: 1900},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 24, col: 91, offset: 1902},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 24, col: 95, offset: 1906},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 24, col: 97, offset: 1908},
													name: "FieldDef",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 111, offset: 1922},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 24, col: 113, offset: 1924},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 117, offset: 1928},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 24, col: 119, offset: 1930},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 123, offset: 1934},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "FieldDef",
			pos:  position{line: 44, col: 1, offset: 2552},
			expr: &actionExpr{
				pos: position{line: 44, col: 13, offset: 2564},
				run: (*parser).callonFieldDef1,
				expr: &seqExpr{
					pos: position{line: 44, col: 13, offset: 2564},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 44, col: 13, offset: 2564},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 44, col: 18, offset: 2569},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 44, col: 24, offset: 2575},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 44, col: 27, offset: 2578},
							label: "dataType",
							expr: &ruleRefExpr{
								pos:  position{line: 44, col: 36, offset: 2587},
								name: "DataType",
							},
						},
						&labeledExpr{
							pos:   position{line: 44, col: 45, offset: 2596},
							label: "notnull",
							expr: &zeroOrOneExpr{
								pos: position{line: 44, col: 53, offset: 2604},
								expr: &seqExpr{
									pos: position{line: 44, col: 55, offset: 2606},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 44, col: 55, offset: 2606},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 44, col: 58, offset: 2609},
											val:        "not",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 44, col: 65, offset: 2616},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 44, col: 68, offset: 2619},
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
			pos:  position{line: 56, col: 1, offset: 2867},
			expr: &actionExpr{
				pos: position{line: 56, col: 13, offset: 2879},
				run: (*parser).callonDataType1,
				expr: &labeledExpr{
					pos:   position{line: 56, col: 13, offset: 2879},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 56, col: 17, offset: 2883},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 56, col: 17, offset: 2883},
								name: "TimestampT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 30, offset: 2896},
								name: "TimeT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 38, offset: 2904},
								name: "VarcharT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 49, offset: 2915},
								name: "CharT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 57, offset: 2923},
								name: "IntT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 64, offset: 2930},
								name: "OtherT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 73, offset: 2939},
								name: "CustomT",
							},
						},
					},
				},
			},
		},
		{
			name: "TimestampT",
			pos:  position{line: 60, col: 1, offset: 2972},
			expr: &actionExpr{
				pos: position{line: 60, col: 15, offset: 2986},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 60, col: 15, offset: 2986},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 60, col: 15, offset: 2986},
							val:        "timestamp",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 60, col: 28, offset: 2999},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 60, col: 33, offset: 3004},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 60, col: 46, offset: 3017},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 60, col: 59, offset: 3030},
								expr: &choiceExpr{
									pos: position{line: 60, col: 61, offset: 3032},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 60, col: 61, offset: 3032},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 60, col: 70, offset: 3041},
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
			pos:  position{line: 73, col: 1, offset: 3320},
			expr: &actionExpr{
				pos: position{line: 73, col: 10, offset: 3329},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 73, col: 10, offset: 3329},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 73, col: 10, offset: 3329},
							val:        "time",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 73, col: 18, offset: 3337},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 73, col: 23, offset: 3342},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 73, col: 36, offset: 3355},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 73, col: 49, offset: 3368},
								expr: &choiceExpr{
									pos: position{line: 73, col: 51, offset: 3370},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 73, col: 51, offset: 3370},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 73, col: 60, offset: 3379},
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
			pos:  position{line: 86, col: 1, offset: 3650},
			expr: &actionExpr{
				pos: position{line: 86, col: 17, offset: 3666},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 86, col: 17, offset: 3666},
					expr: &seqExpr{
						pos: position{line: 86, col: 19, offset: 3668},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 86, col: 19, offset: 3668},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 86, col: 22, offset: 3671},
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
			pos:  position{line: 93, col: 1, offset: 3799},
			expr: &actionExpr{
				pos: position{line: 93, col: 11, offset: 3809},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 93, col: 11, offset: 3809},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 93, col: 11, offset: 3809},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 93, col: 14, offset: 3812},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 93, col: 22, offset: 3820},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 93, col: 25, offset: 3823},
							val:        "time",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 93, col: 33, offset: 3831},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 93, col: 36, offset: 3834},
							val:        "zone",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 97, col: 1, offset: 3868},
			expr: &actionExpr{
				pos: position{line: 97, col: 14, offset: 3881},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 97, col: 14, offset: 3881},
					expr: &seqExpr{
						pos: position{line: 97, col: 16, offset: 3883},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 97, col: 16, offset: 3883},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 97, col: 19, offset: 3886},
								val:        "without",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 97, col: 30, offset: 3897},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 97, col: 33, offset: 3900},
								val:        "time",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 97, col: 41, offset: 3908},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 97, col: 44, offset: 3911},
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
			pos:  position{line: 101, col: 1, offset: 3949},
			expr: &actionExpr{
				pos: position{line: 101, col: 10, offset: 3958},
				run: (*parser).callonCharT1,
				expr: &seqExpr{
					pos: position{line: 101, col: 10, offset: 3958},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 101, col: 12, offset: 3960},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 101, col: 12, offset: 3960},
									val:        "character",
									ignoreCase: true,
								},
								&litMatcher{
									pos:        position{line: 101, col: 27, offset: 3975},
									val:        "char",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 101, col: 37, offset: 3985},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 101, col: 44, offset: 3992},
								expr: &seqExpr{
									pos: position{line: 101, col: 46, offset: 3994},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 101, col: 46, offset: 3994},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 101, col: 50, offset: 3998},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 101, col: 61, offset: 4009},
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
			pos:  position{line: 113, col: 1, offset: 4264},
			expr: &actionExpr{
				pos: position{line: 113, col: 13, offset: 4276},
				run: (*parser).callonVarcharT1,
				expr: &seqExpr{
					pos: position{line: 113, col: 13, offset: 4276},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 113, col: 15, offset: 4278},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 113, col: 17, offset: 4280},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 113, col: 17, offset: 4280},
											val:        "character",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 113, col: 30, offset: 4293},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 113, col: 33, offset: 4296},
											val:        "varying",
											ignoreCase: true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 113, col: 48, offset: 4311},
									val:        "varchar",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 113, col: 61, offset: 4324},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 113, col: 68, offset: 4331},
								expr: &seqExpr{
									pos: position{line: 113, col: 70, offset: 4333},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 113, col: 70, offset: 4333},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 113, col: 74, offset: 4337},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 113, col: 85, offset: 4348},
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
			pos:  position{line: 124, col: 1, offset: 4583},
			expr: &actionExpr{
				pos: position{line: 124, col: 9, offset: 4591},
				run: (*parser).callonIntT1,
				expr: &choiceExpr{
					pos: position{line: 124, col: 11, offset: 4593},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 124, col: 11, offset: 4593},
							val:        "integer",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 124, col: 24, offset: 4606},
							val:        "int",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 130, col: 1, offset: 4688},
			expr: &actionExpr{
				pos: position{line: 130, col: 11, offset: 4698},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 130, col: 13, offset: 4700},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 130, col: 13, offset: 4700},
							val:        "date",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 23, offset: 4710},
							val:        "smallint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 37, offset: 4724},
							val:        "bigint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 49, offset: 4736},
							val:        "decimal",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 62, offset: 4749},
							val:        "numeric",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 75, offset: 4762},
							val:        "real",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 85, offset: 4772},
							val:        "smallserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 102, offset: 4789},
							val:        "serial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 114, offset: 4801},
							val:        "bigserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 129, offset: 4816},
							val:        "boolean",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 142, offset: 4829},
							val:        "text",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 152, offset: 4839},
							val:        "money",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 130, col: 163, offset: 4850},
							val:        "bytea",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CustomT",
			pos:  position{line: 136, col: 1, offset: 4956},
			expr: &actionExpr{
				pos: position{line: 136, col: 13, offset: 4968},
				run: (*parser).callonCustomT1,
				expr: &ruleRefExpr{
					pos:  position{line: 136, col: 13, offset: 4968},
					name: "Ident",
				},
			},
		},
		{
			name: "CreateTypeStmt",
			pos:  position{line: 157, col: 2, offset: 6498},
			expr: &actionExpr{
				pos: position{line: 157, col: 20, offset: 6516},
				run: (*parser).callonCreateTypeStmt1,
				expr: &seqExpr{
					pos: position{line: 157, col: 20, offset: 6516},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 157, col: 20, offset: 6516},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 157, col: 30, offset: 6526},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 157, col: 33, offset: 6529},
							val:        "type",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 157, col: 41, offset: 6537},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 157, col: 44, offset: 6540},
							label: "typename",
							expr: &ruleRefExpr{
								pos:  position{line: 157, col: 53, offset: 6549},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 157, col: 59, offset: 6555},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 157, col: 62, offset: 6558},
							val:        "as",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 157, col: 68, offset: 6564},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 157, col: 71, offset: 6567},
							label: "typedef",
							expr: &ruleRefExpr{
								pos:  position{line: 157, col: 79, offset: 6575},
								name: "EnumDef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 157, col: 87, offset: 6583},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 157, col: 89, offset: 6585},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 157, col: 93, offset: 6589},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "EnumDef",
			pos:  position{line: 163, col: 2, offset: 6706},
			expr: &actionExpr{
				pos: position{line: 163, col: 13, offset: 6717},
				run: (*parser).callonEnumDef1,
				expr: &seqExpr{
					pos: position{line: 163, col: 13, offset: 6717},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 163, col: 13, offset: 6717},
							val:        "ENUM",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 163, col: 20, offset: 6724},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 163, col: 22, offset: 6726},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 163, col: 26, offset: 6730},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 163, col: 28, offset: 6732},
							label: "vals",
							expr: &seqExpr{
								pos: position{line: 163, col: 35, offset: 6739},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 163, col: 35, offset: 6739},
										name: "StringConst",
									},
									&zeroOrMoreExpr{
										pos: position{line: 163, col: 47, offset: 6751},
										expr: &seqExpr{
											pos: position{line: 163, col: 49, offset: 6753},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 163, col: 49, offset: 6753},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 163, col: 51, offset: 6755},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 163, col: 55, offset: 6759},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 163, col: 57, offset: 6761},
													name: "StringConst",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 163, col: 75, offset: 6779},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 163, col: 77, offset: 6781},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 188, col: 1, offset: 8385},
			expr: &actionExpr{
				pos: position{line: 188, col: 25, offset: 8409},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 188, col: 25, offset: 8409},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 188, col: 25, offset: 8409},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 36, offset: 8420},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 188, col: 39, offset: 8423},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 45, offset: 8429},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 188, col: 48, offset: 8432},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 61, offset: 8445},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 63, offset: 8447},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 188, col: 73, offset: 8457},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 79, offset: 8463},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 188, col: 81, offset: 8465},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 87, offset: 8471},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 89, offset: 8473},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 188, col: 97, offset: 8481},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 109, offset: 8493},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 188, col: 111, offset: 8495},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 115, offset: 8499},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 192, col: 1, offset: 8587},
			expr: &actionExpr{
				pos: position{line: 192, col: 24, offset: 8610},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 192, col: 24, offset: 8610},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 192, col: 24, offset: 8610},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 192, col: 34, offset: 8620},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 192, col: 37, offset: 8623},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 192, col: 50, offset: 8636},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 192, col: 53, offset: 8639},
							expr: &seqExpr{
								pos: position{line: 192, col: 55, offset: 8641},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 192, col: 55, offset: 8641},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 192, col: 61, offset: 8647},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 192, col: 64, offset: 8650},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 192, col: 71, offset: 8657},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 192, col: 74, offset: 8660},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 192, col: 84, offset: 8670},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 192, col: 90, offset: 8676},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 192, col: 100, offset: 8686},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 192, col: 106, offset: 8692},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 192, col: 109, offset: 8695},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 192, col: 117, offset: 8703},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 192, col: 120, offset: 8706},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 192, col: 130, offset: 8716},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 192, col: 133, offset: 8719},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 192, col: 140, offset: 8726},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 192, col: 146, offset: 8732},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 192, col: 148, offset: 8734},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 192, col: 152, offset: 8738},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 196, col: 1, offset: 8828},
			expr: &actionExpr{
				pos: position{line: 196, col: 12, offset: 8839},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 196, col: 12, offset: 8839},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 196, col: 12, offset: 8839},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 196, col: 19, offset: 8846},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 196, col: 21, offset: 8848},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 196, col: 25, offset: 8852},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 196, col: 29, offset: 8856},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 196, col: 33, offset: 8860},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 196, col: 33, offset: 8860},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 196, col: 39, offset: 8866},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 196, col: 47, offset: 8874},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 196, col: 49, offset: 8876},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 196, col: 56, offset: 8883},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 196, col: 77, offset: 8904},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 196, col: 79, offset: 8906},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 196, col: 83, offset: 8910},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 201, col: 1, offset: 8991},
			expr: &actionExpr{
				pos: position{line: 201, col: 8, offset: 8998},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 201, col: 8, offset: 8998},
					expr: &charClassMatcher{
						pos:        position{line: 201, col: 8, offset: 8998},
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
			pos:  position{line: 216, col: 1, offset: 9838},
			expr: &actionExpr{
				pos: position{line: 216, col: 25, offset: 9862},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 216, col: 25, offset: 9862},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 216, col: 32, offset: 9869},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 216, col: 32, offset: 9869},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 216, col: 38, offset: 9875},
								expr: &seqExpr{
									pos: position{line: 216, col: 40, offset: 9877},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 216, col: 40, offset: 9877},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 216, col: 42, offset: 9879},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 216, col: 46, offset: 9883},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 216, col: 48, offset: 9885},
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
			pos:  position{line: 228, col: 1, offset: 10175},
			expr: &choiceExpr{
				pos: position{line: 228, col: 12, offset: 10186},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 228, col: 12, offset: 10186},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 228, col: 21, offset: 10195},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 228, col: 31, offset: 10205},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 228, col: 45, offset: 10219},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 230, col: 1, offset: 10228},
			expr: &actionExpr{
				pos: position{line: 230, col: 16, offset: 10243},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 230, col: 16, offset: 10243},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 230, col: 16, offset: 10243},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 230, col: 20, offset: 10247},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 230, col: 26, offset: 10253},
								expr: &choiceExpr{
									pos: position{line: 230, col: 27, offset: 10254},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 230, col: 27, offset: 10254},
											val:        "[^'\\n]",
											chars:      []rune{'\'', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 230, col: 36, offset: 10263},
											val:        "''",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 230, col: 43, offset: 10270},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 234, col: 1, offset: 10322},
			expr: &actionExpr{
				pos: position{line: 234, col: 10, offset: 10331},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 234, col: 10, offset: 10331},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 234, col: 10, offset: 10331},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 234, col: 18, offset: 10339},
							expr: &charClassMatcher{
								pos:        position{line: 234, col: 18, offset: 10339},
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
			pos:  position{line: 238, col: 1, offset: 10392},
			expr: &actionExpr{
				pos: position{line: 238, col: 11, offset: 10402},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 238, col: 13, offset: 10404},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 238, col: 13, offset: 10404},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 238, col: 19, offset: 10410},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 238, col: 19, offset: 10410},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 238, col: 24, offset: 10415},
									expr: &charClassMatcher{
										pos:        position{line: 238, col: 24, offset: 10415},
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
			pos:  position{line: 243, col: 1, offset: 10510},
			expr: &actionExpr{
				pos: position{line: 243, col: 15, offset: 10524},
				run: (*parser).callonNonZNumber1,
				expr: &seqExpr{
					pos: position{line: 243, col: 15, offset: 10524},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 243, col: 15, offset: 10524},
							val:        "[1-9]",
							ranges:     []rune{'1', '9'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 243, col: 20, offset: 10529},
							expr: &charClassMatcher{
								pos:        position{line: 243, col: 20, offset: 10529},
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
			pos:  position{line: 248, col: 1, offset: 10622},
			expr: &actionExpr{
				pos: position{line: 248, col: 12, offset: 10633},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 248, col: 12, offset: 10633},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 248, col: 20, offset: 10641},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 248, col: 20, offset: 10641},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 248, col: 34, offset: 10655},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 252, col: 1, offset: 10697},
			expr: &actionExpr{
				pos: position{line: 252, col: 16, offset: 10712},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 252, col: 18, offset: 10714},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 252, col: 18, offset: 10714},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 252, col: 27, offset: 10723},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 252, col: 27, offset: 10723},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 252, col: 31, offset: 10727},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 252, col: 49, offset: 10745},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 252, col: 55, offset: 10751},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 256, col: 1, offset: 10797},
			expr: &choiceExpr{
				pos: position{line: 256, col: 24, offset: 10820},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 256, col: 24, offset: 10820},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 256, col: 33, offset: 10829},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 256, col: 41, offset: 10837},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 256, col: 48, offset: 10844},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 256, col: 54, offset: 10850},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 258, col: 1, offset: 10857},
			expr: &actionExpr{
				pos: position{line: 258, col: 17, offset: 10873},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 258, col: 19, offset: 10875},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 258, col: 19, offset: 10875},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 258, col: 29, offset: 10885},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 258, col: 29, offset: 10885},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 258, col: 33, offset: 10889},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 258, col: 52, offset: 10908},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 258, col: 58, offset: 10914},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 262, col: 1, offset: 10962},
			expr: &choiceExpr{
				pos: position{line: 262, col: 25, offset: 10986},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 262, col: 25, offset: 10986},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 262, col: 35, offset: 10996},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 262, col: 42, offset: 11003},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 262, col: 50, offset: 11011},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 262, col: 56, offset: 11017},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 274, col: 1, offset: 11532},
			expr: &actionExpr{
				pos: position{line: 274, col: 12, offset: 11543},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 274, col: 14, offset: 11545},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 274, col: 14, offset: 11545},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 274, col: 34, offset: 11565},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 278, col: 1, offset: 11609},
			expr: &seqExpr{
				pos: position{line: 278, col: 21, offset: 11629},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 278, col: 21, offset: 11629},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 278, col: 26, offset: 11634},
						expr: &seqExpr{
							pos: position{line: 278, col: 28, offset: 11636},
							exprs: []interface{}{
								&zeroOrMoreExpr{
									pos: position{line: 278, col: 28, offset: 11636},
									expr: &anyMatcher{
										line: 278, col: 28, offset: 11636,
									},
								},
								&zeroOrOneExpr{
									pos: position{line: 278, col: 31, offset: 11639},
									expr: &ruleRefExpr{
										pos:  position{line: 278, col: 31, offset: 11639},
										name: "MultilineComment",
									},
								},
								&zeroOrMoreExpr{
									pos: position{line: 278, col: 49, offset: 11657},
									expr: &anyMatcher{
										line: 278, col: 49, offset: 11657,
									},
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 278, col: 55, offset: 11663},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 278, col: 60, offset: 11668},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 280, col: 1, offset: 11673},
			expr: &seqExpr{
				pos: position{line: 280, col: 22, offset: 11694},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 280, col: 22, offset: 11694},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 280, col: 27, offset: 11699},
						expr: &charClassMatcher{
							pos:        position{line: 280, col: 27, offset: 11699},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 280, col: 36, offset: 11708},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 282, col: 1, offset: 11713},
			expr: &seqExpr{
				pos: position{line: 282, col: 9, offset: 11721},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 282, col: 9, offset: 11721},
						expr: &charClassMatcher{
							pos:        position{line: 282, col: 9, offset: 11721},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 282, col: 17, offset: 11729},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 282, col: 17, offset: 11729},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 282, col: 26, offset: 11738},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 282, col: 35, offset: 11747},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 282, col: 42, offset: 11754},
								val:        "\n",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 282, col: 49, offset: 11761},
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
			pos:         position{line: 284, col: 1, offset: 11767},
			expr: &zeroOrMoreExpr{
				pos: position{line: 284, col: 19, offset: 11785},
				expr: &charClassMatcher{
					pos:        position{line: 284, col: 19, offset: 11785},
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
			pos:         position{line: 286, col: 1, offset: 11797},
			expr: &oneOrMoreExpr{
				pos: position{line: 286, col: 31, offset: 11827},
				expr: &charClassMatcher{
					pos:        position{line: 286, col: 31, offset: 11827},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 288, col: 1, offset: 11839},
			expr: &notExpr{
				pos: position{line: 288, col: 8, offset: 11846},
				expr: &anyMatcher{
					line: 288, col: 9, offset: 11847,
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
