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
			pos:  position{line: 11, col: 1, offset: 169},
			expr: &actionExpr{
				pos: position{line: 11, col: 8, offset: 176},
				run: (*parser).callonSQL1,
				expr: &labeledExpr{
					pos:   position{line: 11, col: 8, offset: 176},
					label: "stmts",
					expr: &oneOrMoreExpr{
						pos: position{line: 11, col: 14, offset: 182},
						expr: &ruleRefExpr{
							pos:  position{line: 11, col: 14, offset: 182},
							name: "Stmt",
						},
					},
				},
			},
		},
		{
			name: "Stmt",
			pos:  position{line: 15, col: 1, offset: 215},
			expr: &actionExpr{
				pos: position{line: 15, col: 9, offset: 223},
				run: (*parser).callonStmt1,
				expr: &seqExpr{
					pos: position{line: 15, col: 9, offset: 223},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 15, col: 9, offset: 223},
							expr: &ruleRefExpr{
								pos:  position{line: 15, col: 9, offset: 223},
								name: "Comment",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 15, col: 18, offset: 232},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 15, col: 20, offset: 234},
							label: "stmt",
							expr: &choiceExpr{
								pos: position{line: 15, col: 27, offset: 241},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 15, col: 27, offset: 241},
										name: "SetStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 15, col: 37, offset: 251},
										name: "CreateTableStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 15, col: 55, offset: 269},
										name: "CreateSeqStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 15, col: 71, offset: 285},
										name: "CreateExtensionStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 15, col: 93, offset: 307},
										name: "CreateTypeStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 15, col: 110, offset: 324},
										name: "AlterTableStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 15, col: 127, offset: 341},
										name: "AlterSeqStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 15, col: 142, offset: 356},
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
			pos:  position{line: 29, col: 1, offset: 1933},
			expr: &actionExpr{
				pos: position{line: 29, col: 20, offset: 1952},
				run: (*parser).callonCreateTableStmt1,
				expr: &seqExpr{
					pos: position{line: 29, col: 20, offset: 1952},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 29, col: 20, offset: 1952},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 29, col: 30, offset: 1962},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 29, col: 33, offset: 1965},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 29, col: 42, offset: 1974},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 29, col: 45, offset: 1977},
							label: "tablename",
							expr: &ruleRefExpr{
								pos:  position{line: 29, col: 55, offset: 1987},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 29, col: 61, offset: 1993},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 29, col: 63, offset: 1995},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 29, col: 67, offset: 1999},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 29, col: 69, offset: 2001},
							label: "defs",
							expr: &seqExpr{
								pos: position{line: 29, col: 76, offset: 2008},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 29, col: 76, offset: 2008},
										name: "TableDef",
									},
									&zeroOrMoreExpr{
										pos: position{line: 29, col: 85, offset: 2017},
										expr: &seqExpr{
											pos: position{line: 29, col: 87, offset: 2019},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 29, col: 87, offset: 2019},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 29, col: 89, offset: 2021},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 29, col: 93, offset: 2025},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 29, col: 95, offset: 2027},
													name: "TableDef",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 29, col: 109, offset: 2041},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 29, col: 111, offset: 2043},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 29, col: 115, offset: 2047},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 29, col: 117, offset: 2049},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 29, col: 121, offset: 2053},
							expr: &ruleRefExpr{
								pos:  position{line: 29, col: 121, offset: 2053},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "TableDef",
			pos:  position{line: 49, col: 1, offset: 2650},
			expr: &choiceExpr{
				pos: position{line: 49, col: 13, offset: 2662},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 49, col: 13, offset: 2662},
						name: "TableConstr",
					},
					&ruleRefExpr{
						pos:  position{line: 49, col: 27, offset: 2676},
						name: "ColumnDef",
					},
				},
			},
		},
		{
			name: "ColumnDef",
			pos:  position{line: 51, col: 1, offset: 2687},
			expr: &actionExpr{
				pos: position{line: 51, col: 14, offset: 2700},
				run: (*parser).callonColumnDef1,
				expr: &seqExpr{
					pos: position{line: 51, col: 14, offset: 2700},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 51, col: 14, offset: 2700},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 51, col: 20, offset: 2706},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 51, col: 20, offset: 2706},
										name: "DblQuotedString",
									},
									&ruleRefExpr{
										pos:  position{line: 51, col: 38, offset: 2724},
										name: "StringConst",
									},
									&ruleRefExpr{
										pos:  position{line: 51, col: 52, offset: 2738},
										name: "Ident",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 51, col: 59, offset: 2745},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 51, col: 62, offset: 2748},
							label: "dataType",
							expr: &ruleRefExpr{
								pos:  position{line: 51, col: 71, offset: 2757},
								name: "DataType",
							},
						},
						&labeledExpr{
							pos:   position{line: 51, col: 80, offset: 2766},
							label: "constraint",
							expr: &zeroOrOneExpr{
								pos: position{line: 51, col: 91, offset: 2777},
								expr: &ruleRefExpr{
									pos:  position{line: 51, col: 91, offset: 2777},
									name: "ColumnConstraint",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ColumnConstraint",
			pos:  position{line: 68, col: 1, offset: 3266},
			expr: &actionExpr{
				pos: position{line: 68, col: 21, offset: 3286},
				run: (*parser).callonColumnConstraint1,
				expr: &seqExpr{
					pos: position{line: 68, col: 21, offset: 3286},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 68, col: 21, offset: 3286},
							label: "nameOpt",
							expr: &zeroOrOneExpr{
								pos: position{line: 68, col: 29, offset: 3294},
								expr: &seqExpr{
									pos: position{line: 68, col: 31, offset: 3296},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 68, col: 31, offset: 3296},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 68, col: 34, offset: 3299},
											val:        "constraint",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 68, col: 48, offset: 3313},
											name: "_1",
										},
										&choiceExpr{
											pos: position{line: 68, col: 52, offset: 3317},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 68, col: 52, offset: 3317},
													name: "StringConst",
												},
												&ruleRefExpr{
													pos:  position{line: 68, col: 66, offset: 3331},
													name: "Ident",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 68, col: 76, offset: 3341},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 68, col: 78, offset: 3343},
							label: "constraint",
							expr: &choiceExpr{
								pos: position{line: 68, col: 91, offset: 3356},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 68, col: 91, offset: 3356},
										name: "NotNullCls",
									},
									&ruleRefExpr{
										pos:  position{line: 68, col: 104, offset: 3369},
										name: "NullCls",
									},
									&ruleRefExpr{
										pos:  position{line: 68, col: 114, offset: 3379},
										name: "CheckCls",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TableConstr",
			pos:  position{line: 83, col: 1, offset: 3745},
			expr: &actionExpr{
				pos: position{line: 83, col: 16, offset: 3760},
				run: (*parser).callonTableConstr1,
				expr: &seqExpr{
					pos: position{line: 83, col: 16, offset: 3760},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 83, col: 16, offset: 3760},
							label: "nameOpt",
							expr: &zeroOrOneExpr{
								pos: position{line: 83, col: 24, offset: 3768},
								expr: &seqExpr{
									pos: position{line: 83, col: 26, offset: 3770},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 83, col: 26, offset: 3770},
											val:        "constraint",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 83, col: 40, offset: 3784},
											name: "_1",
										},
										&choiceExpr{
											pos: position{line: 83, col: 44, offset: 3788},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 83, col: 44, offset: 3788},
													name: "StringConst",
												},
												&ruleRefExpr{
													pos:  position{line: 83, col: 58, offset: 3802},
													name: "Ident",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 83, col: 68, offset: 3812},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 83, col: 70, offset: 3814},
							label: "constraint",
							expr: &ruleRefExpr{
								pos:  position{line: 83, col: 81, offset: 3825},
								name: "CheckCls",
							},
						},
					},
				},
			},
		},
		{
			name: "NotNullCls",
			pos:  position{line: 100, col: 1, offset: 4226},
			expr: &actionExpr{
				pos: position{line: 100, col: 15, offset: 4240},
				run: (*parser).callonNotNullCls1,
				expr: &seqExpr{
					pos: position{line: 100, col: 15, offset: 4240},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 100, col: 15, offset: 4240},
							val:        "not",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 22, offset: 4247},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 100, col: 25, offset: 4250},
							val:        "null",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "NullCls",
			pos:  position{line: 106, col: 1, offset: 4332},
			expr: &actionExpr{
				pos: position{line: 106, col: 12, offset: 4343},
				run: (*parser).callonNullCls1,
				expr: &litMatcher{
					pos:        position{line: 106, col: 12, offset: 4343},
					val:        "null",
					ignoreCase: true,
				},
			},
		},
		{
			name: "CheckCls",
			pos:  position{line: 112, col: 1, offset: 4426},
			expr: &actionExpr{
				pos: position{line: 112, col: 13, offset: 4438},
				run: (*parser).callonCheckCls1,
				expr: &seqExpr{
					pos: position{line: 112, col: 13, offset: 4438},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 112, col: 13, offset: 4438},
							val:        "check",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 112, col: 22, offset: 4447},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 112, col: 25, offset: 4450},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 112, col: 30, offset: 4455},
								name: "WrappedExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 112, col: 42, offset: 4467},
							label: "noInherit",
							expr: &zeroOrOneExpr{
								pos: position{line: 112, col: 52, offset: 4477},
								expr: &seqExpr{
									pos: position{line: 112, col: 54, offset: 4479},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 112, col: 54, offset: 4479},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 112, col: 57, offset: 4482},
											val:        "no",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 112, col: 63, offset: 4488},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 112, col: 66, offset: 4491},
											val:        "inherit",
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
			name: "WrappedExpr",
			pos:  position{line: 122, col: 1, offset: 4684},
			expr: &actionExpr{
				pos: position{line: 122, col: 16, offset: 4699},
				run: (*parser).callonWrappedExpr1,
				expr: &seqExpr{
					pos: position{line: 122, col: 16, offset: 4699},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 122, col: 16, offset: 4699},
							val:        "(",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 122, col: 20, offset: 4703},
							expr: &ruleRefExpr{
								pos:  position{line: 122, col: 20, offset: 4703},
								name: "Expr",
							},
						},
						&litMatcher{
							pos:        position{line: 122, col: 26, offset: 4709},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 126, col: 1, offset: 4749},
			expr: &choiceExpr{
				pos: position{line: 126, col: 9, offset: 4757},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 126, col: 9, offset: 4757},
						name: "WrappedExpr",
					},
					&oneOrMoreExpr{
						pos: position{line: 126, col: 23, offset: 4771},
						expr: &charClassMatcher{
							pos:        position{line: 126, col: 23, offset: 4771},
							val:        "[^()]",
							chars:      []rune{'(', ')'},
							ignoreCase: false,
							inverted:   true,
						},
					},
				},
			},
		},
		{
			name: "DataType",
			pos:  position{line: 140, col: 1, offset: 5793},
			expr: &actionExpr{
				pos: position{line: 140, col: 13, offset: 5805},
				run: (*parser).callonDataType1,
				expr: &seqExpr{
					pos: position{line: 140, col: 13, offset: 5805},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 140, col: 13, offset: 5805},
							label: "t",
							expr: &choiceExpr{
								pos: position{line: 140, col: 17, offset: 5809},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 140, col: 17, offset: 5809},
										name: "TimestampT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 30, offset: 5822},
										name: "TimeT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 38, offset: 5830},
										name: "NumericT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 49, offset: 5841},
										name: "VarcharT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 60, offset: 5852},
										name: "CharT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 68, offset: 5860},
										name: "BitVarT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 78, offset: 5870},
										name: "BitT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 85, offset: 5877},
										name: "IntT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 92, offset: 5884},
										name: "PgOidT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 101, offset: 5893},
										name: "PostgisT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 112, offset: 5904},
										name: "OtherT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 121, offset: 5913},
										name: "CustomT",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 140, col: 131, offset: 5923},
							label: "brackets",
							expr: &zeroOrMoreExpr{
								pos: position{line: 140, col: 140, offset: 5932},
								expr: &litMatcher{
									pos:        position{line: 140, col: 142, offset: 5934},
									val:        "[]",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TimestampT",
			pos:  position{line: 153, col: 1, offset: 6245},
			expr: &actionExpr{
				pos: position{line: 153, col: 15, offset: 6259},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 153, col: 15, offset: 6259},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 153, col: 15, offset: 6259},
							val:        "timestamp",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 153, col: 28, offset: 6272},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 153, col: 33, offset: 6277},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 153, col: 46, offset: 6290},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 153, col: 59, offset: 6303},
								expr: &choiceExpr{
									pos: position{line: 153, col: 61, offset: 6305},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 153, col: 61, offset: 6305},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 153, col: 70, offset: 6314},
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
			pos:  position{line: 166, col: 1, offset: 6593},
			expr: &actionExpr{
				pos: position{line: 166, col: 10, offset: 6602},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 166, col: 10, offset: 6602},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 166, col: 10, offset: 6602},
							val:        "time",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 166, col: 18, offset: 6610},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 166, col: 23, offset: 6615},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 166, col: 36, offset: 6628},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 166, col: 49, offset: 6641},
								expr: &choiceExpr{
									pos: position{line: 166, col: 51, offset: 6643},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 166, col: 51, offset: 6643},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 166, col: 60, offset: 6652},
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
			pos:  position{line: 179, col: 1, offset: 6923},
			expr: &actionExpr{
				pos: position{line: 179, col: 17, offset: 6939},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 179, col: 17, offset: 6939},
					expr: &seqExpr{
						pos: position{line: 179, col: 19, offset: 6941},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 179, col: 19, offset: 6941},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 179, col: 22, offset: 6944},
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
			pos:  position{line: 186, col: 1, offset: 7072},
			expr: &actionExpr{
				pos: position{line: 186, col: 11, offset: 7082},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 186, col: 11, offset: 7082},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 186, col: 11, offset: 7082},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 186, col: 14, offset: 7085},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 22, offset: 7093},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 186, col: 25, offset: 7096},
							val:        "time",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 33, offset: 7104},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 186, col: 36, offset: 7107},
							val:        "zone",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 190, col: 1, offset: 7141},
			expr: &actionExpr{
				pos: position{line: 190, col: 14, offset: 7154},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 190, col: 14, offset: 7154},
					expr: &seqExpr{
						pos: position{line: 190, col: 16, offset: 7156},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 190, col: 16, offset: 7156},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 190, col: 19, offset: 7159},
								val:        "without",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 190, col: 30, offset: 7170},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 190, col: 33, offset: 7173},
								val:        "time",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 190, col: 41, offset: 7181},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 190, col: 44, offset: 7184},
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
			pos:  position{line: 194, col: 1, offset: 7222},
			expr: &actionExpr{
				pos: position{line: 194, col: 10, offset: 7231},
				run: (*parser).callonCharT1,
				expr: &seqExpr{
					pos: position{line: 194, col: 10, offset: 7231},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 194, col: 12, offset: 7233},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 194, col: 12, offset: 7233},
									val:        "character",
									ignoreCase: true,
								},
								&litMatcher{
									pos:        position{line: 194, col: 27, offset: 7248},
									val:        "char",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 194, col: 37, offset: 7258},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 194, col: 44, offset: 7265},
								expr: &seqExpr{
									pos: position{line: 194, col: 46, offset: 7267},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 194, col: 46, offset: 7267},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 50, offset: 7271},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 194, col: 61, offset: 7282},
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
			pos:  position{line: 206, col: 1, offset: 7537},
			expr: &actionExpr{
				pos: position{line: 206, col: 13, offset: 7549},
				run: (*parser).callonVarcharT1,
				expr: &seqExpr{
					pos: position{line: 206, col: 13, offset: 7549},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 206, col: 15, offset: 7551},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 206, col: 17, offset: 7553},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 206, col: 17, offset: 7553},
											val:        "character",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 206, col: 30, offset: 7566},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 206, col: 33, offset: 7569},
											val:        "varying",
											ignoreCase: true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 206, col: 48, offset: 7584},
									val:        "varchar",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 206, col: 61, offset: 7597},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 206, col: 68, offset: 7604},
								expr: &seqExpr{
									pos: position{line: 206, col: 70, offset: 7606},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 206, col: 70, offset: 7606},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 206, col: 74, offset: 7610},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 206, col: 85, offset: 7621},
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
			pos:  position{line: 217, col: 1, offset: 7856},
			expr: &actionExpr{
				pos: position{line: 217, col: 9, offset: 7864},
				run: (*parser).callonBitT1,
				expr: &seqExpr{
					pos: position{line: 217, col: 9, offset: 7864},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 217, col: 9, offset: 7864},
							val:        "bit",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 217, col: 16, offset: 7871},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 217, col: 23, offset: 7878},
								expr: &seqExpr{
									pos: position{line: 217, col: 25, offset: 7880},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 217, col: 25, offset: 7880},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 217, col: 29, offset: 7884},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 217, col: 40, offset: 7895},
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
			pos:  position{line: 229, col: 1, offset: 8149},
			expr: &actionExpr{
				pos: position{line: 229, col: 12, offset: 8160},
				run: (*parser).callonBitVarT1,
				expr: &seqExpr{
					pos: position{line: 229, col: 12, offset: 8160},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 229, col: 12, offset: 8160},
							val:        "bit",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 229, col: 19, offset: 8167},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 229, col: 22, offset: 8170},
							val:        "varying",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 229, col: 33, offset: 8181},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 229, col: 40, offset: 8188},
								expr: &seqExpr{
									pos: position{line: 229, col: 42, offset: 8190},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 229, col: 42, offset: 8190},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 229, col: 46, offset: 8194},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 229, col: 57, offset: 8205},
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
			pos:  position{line: 240, col: 1, offset: 8439},
			expr: &actionExpr{
				pos: position{line: 240, col: 9, offset: 8447},
				run: (*parser).callonIntT1,
				expr: &choiceExpr{
					pos: position{line: 240, col: 11, offset: 8449},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 240, col: 11, offset: 8449},
							val:        "integer",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 24, offset: 8462},
							val:        "int",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "NumericT",
			pos:  position{line: 246, col: 1, offset: 8544},
			expr: &actionExpr{
				pos: position{line: 246, col: 13, offset: 8556},
				run: (*parser).callonNumericT1,
				expr: &seqExpr{
					pos: position{line: 246, col: 13, offset: 8556},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 246, col: 13, offset: 8556},
							val:        "numeric",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 246, col: 24, offset: 8567},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 246, col: 29, offset: 8572},
								expr: &seqExpr{
									pos: position{line: 246, col: 31, offset: 8574},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 246, col: 31, offset: 8574},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 246, col: 35, offset: 8578},
											name: "NonZNumber",
										},
										&zeroOrOneExpr{
											pos: position{line: 246, col: 46, offset: 8589},
											expr: &seqExpr{
												pos: position{line: 246, col: 48, offset: 8591},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 246, col: 48, offset: 8591},
														val:        ",",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 246, col: 52, offset: 8595},
														name: "_",
													},
													&ruleRefExpr{
														pos:  position{line: 246, col: 54, offset: 8597},
														name: "NonZNumber",
													},
												},
											},
										},
										&litMatcher{
											pos:        position{line: 246, col: 68, offset: 8611},
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
			name: "PostgisT",
			pos:  position{line: 261, col: 1, offset: 9014},
			expr: &actionExpr{
				pos: position{line: 261, col: 13, offset: 9026},
				run: (*parser).callonPostgisT1,
				expr: &seqExpr{
					pos: position{line: 261, col: 13, offset: 9026},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 261, col: 13, offset: 9026},
							label: "t",
							expr: &choiceExpr{
								pos: position{line: 261, col: 17, offset: 9030},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 261, col: 17, offset: 9030},
										val:        "geography",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 261, col: 32, offset: 9045},
										val:        "geometry",
										ignoreCase: true,
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 261, col: 46, offset: 9059},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 261, col: 50, offset: 9063},
							label: "subtype",
							expr: &choiceExpr{
								pos: position{line: 261, col: 60, offset: 9073},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 261, col: 60, offset: 9073},
										val:        "point",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 261, col: 71, offset: 9084},
										val:        "linestring",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 261, col: 87, offset: 9100},
										val:        "polygon",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 261, col: 100, offset: 9113},
										val:        "multipoint",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 261, col: 116, offset: 9129},
										val:        "multilinestring",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 261, col: 137, offset: 9150},
										val:        "multipolygon",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 261, col: 155, offset: 9168},
										val:        "geometrycollection",
										ignoreCase: true,
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 261, col: 179, offset: 9192},
							label: "srid",
							expr: &zeroOrOneExpr{
								pos: position{line: 261, col: 184, offset: 9197},
								expr: &seqExpr{
									pos: position{line: 261, col: 185, offset: 9198},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 261, col: 185, offset: 9198},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 261, col: 189, offset: 9202},
											name: "NonZNumber",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 261, col: 202, offset: 9215},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PgOidT",
			pos:  position{line: 273, col: 1, offset: 9537},
			expr: &actionExpr{
				pos: position{line: 273, col: 11, offset: 9547},
				run: (*parser).callonPgOidT1,
				expr: &choiceExpr{
					pos: position{line: 273, col: 13, offset: 9549},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 273, col: 13, offset: 9549},
							val:        "oid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 22, offset: 9558},
							val:        "regprocedure",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 40, offset: 9576},
							val:        "regproc",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 53, offset: 9589},
							val:        "regoperator",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 70, offset: 9606},
							val:        "regoper",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 83, offset: 9619},
							val:        "regclass",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 97, offset: 9633},
							val:        "regtype",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 110, offset: 9646},
							val:        "regrole",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 123, offset: 9659},
							val:        "regnamespace",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 141, offset: 9677},
							val:        "regconfig",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 273, col: 156, offset: 9692},
							val:        "regdictionary",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 279, col: 1, offset: 9806},
			expr: &actionExpr{
				pos: position{line: 279, col: 11, offset: 9816},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 279, col: 13, offset: 9818},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 279, col: 13, offset: 9818},
							val:        "date",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 23, offset: 9828},
							val:        "smallint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 37, offset: 9842},
							val:        "bigint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 49, offset: 9854},
							val:        "decimal",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 62, offset: 9867},
							val:        "real",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 72, offset: 9877},
							val:        "smallserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 89, offset: 9894},
							val:        "serial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 101, offset: 9906},
							val:        "bigserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 116, offset: 9921},
							val:        "boolean",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 129, offset: 9934},
							val:        "text",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 139, offset: 9944},
							val:        "money",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 150, offset: 9955},
							val:        "bytea",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 161, offset: 9966},
							val:        "point",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 172, offset: 9977},
							val:        "line",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 182, offset: 9987},
							val:        "lseg",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 192, offset: 9997},
							val:        "box",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 201, offset: 10006},
							val:        "path",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 211, offset: 10016},
							val:        "polygon",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 224, offset: 10029},
							val:        "circle",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 236, offset: 10041},
							val:        "cidr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 246, offset: 10051},
							val:        "inet",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 256, offset: 10061},
							val:        "macaddr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 269, offset: 10074},
							val:        "uuid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 279, offset: 10084},
							val:        "xml",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 288, offset: 10093},
							val:        "jsonb",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 279, col: 299, offset: 10104},
							val:        "json",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CustomT",
			pos:  position{line: 285, col: 1, offset: 10209},
			expr: &actionExpr{
				pos: position{line: 285, col: 13, offset: 10221},
				run: (*parser).callonCustomT1,
				expr: &ruleRefExpr{
					pos:  position{line: 285, col: 13, offset: 10221},
					name: "Ident",
				},
			},
		},
		{
			name: "CreateSeqStmt",
			pos:  position{line: 306, col: 1, offset: 11686},
			expr: &actionExpr{
				pos: position{line: 306, col: 18, offset: 11703},
				run: (*parser).callonCreateSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 306, col: 18, offset: 11703},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 306, col: 18, offset: 11703},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 306, col: 28, offset: 11713},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 306, col: 31, offset: 11716},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 306, col: 43, offset: 11728},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 306, col: 46, offset: 11731},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 306, col: 51, offset: 11736},
								name: "Ident",
							},
						},
						&labeledExpr{
							pos:   position{line: 306, col: 57, offset: 11742},
							label: "verses",
							expr: &zeroOrMoreExpr{
								pos: position{line: 306, col: 64, offset: 11749},
								expr: &ruleRefExpr{
									pos:  position{line: 306, col: 64, offset: 11749},
									name: "CreateSeqVerse",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 306, col: 80, offset: 11765},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 306, col: 82, offset: 11767},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 306, col: 86, offset: 11771},
							expr: &ruleRefExpr{
								pos:  position{line: 306, col: 86, offset: 11771},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateSeqVerse",
			pos:  position{line: 320, col: 1, offset: 12165},
			expr: &actionExpr{
				pos: position{line: 320, col: 19, offset: 12183},
				run: (*parser).callonCreateSeqVerse1,
				expr: &labeledExpr{
					pos:   position{line: 320, col: 19, offset: 12183},
					label: "verse",
					expr: &choiceExpr{
						pos: position{line: 320, col: 27, offset: 12191},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 320, col: 27, offset: 12191},
								name: "IncrementBy",
							},
							&ruleRefExpr{
								pos:  position{line: 320, col: 41, offset: 12205},
								name: "MinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 320, col: 52, offset: 12216},
								name: "NoMinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 320, col: 65, offset: 12229},
								name: "MaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 320, col: 76, offset: 12240},
								name: "NoMaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 320, col: 89, offset: 12253},
								name: "Start",
							},
							&ruleRefExpr{
								pos:  position{line: 320, col: 97, offset: 12261},
								name: "Cache",
							},
							&ruleRefExpr{
								pos:  position{line: 320, col: 105, offset: 12269},
								name: "Cycle",
							},
							&ruleRefExpr{
								pos:  position{line: 320, col: 113, offset: 12277},
								name: "OwnedBy",
							},
						},
					},
				},
			},
		},
		{
			name: "IncrementBy",
			pos:  position{line: 324, col: 1, offset: 12314},
			expr: &actionExpr{
				pos: position{line: 324, col: 16, offset: 12329},
				run: (*parser).callonIncrementBy1,
				expr: &seqExpr{
					pos: position{line: 324, col: 16, offset: 12329},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 324, col: 16, offset: 12329},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 324, col: 19, offset: 12332},
							val:        "increment",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 324, col: 32, offset: 12345},
							expr: &seqExpr{
								pos: position{line: 324, col: 33, offset: 12346},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 324, col: 33, offset: 12346},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 324, col: 36, offset: 12349},
										val:        "by",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 324, col: 44, offset: 12357},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 324, col: 47, offset: 12360},
							label: "num",
							expr: &ruleRefExpr{
								pos:  position{line: 324, col: 51, offset: 12364},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "MinValue",
			pos:  position{line: 330, col: 1, offset: 12478},
			expr: &actionExpr{
				pos: position{line: 330, col: 13, offset: 12490},
				run: (*parser).callonMinValue1,
				expr: &seqExpr{
					pos: position{line: 330, col: 13, offset: 12490},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 330, col: 13, offset: 12490},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 330, col: 16, offset: 12493},
							val:        "minvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 330, col: 28, offset: 12505},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 330, col: 31, offset: 12508},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 330, col: 35, offset: 12512},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMinValue",
			pos:  position{line: 336, col: 1, offset: 12625},
			expr: &actionExpr{
				pos: position{line: 336, col: 15, offset: 12639},
				run: (*parser).callonNoMinValue1,
				expr: &seqExpr{
					pos: position{line: 336, col: 15, offset: 12639},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 336, col: 15, offset: 12639},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 336, col: 18, offset: 12642},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 336, col: 24, offset: 12648},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 336, col: 27, offset: 12651},
							val:        "minvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "MaxValue",
			pos:  position{line: 340, col: 1, offset: 12688},
			expr: &actionExpr{
				pos: position{line: 340, col: 13, offset: 12700},
				run: (*parser).callonMaxValue1,
				expr: &seqExpr{
					pos: position{line: 340, col: 13, offset: 12700},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 340, col: 13, offset: 12700},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 340, col: 16, offset: 12703},
							val:        "maxvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 340, col: 28, offset: 12715},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 340, col: 31, offset: 12718},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 340, col: 35, offset: 12722},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMaxValue",
			pos:  position{line: 346, col: 1, offset: 12835},
			expr: &actionExpr{
				pos: position{line: 346, col: 15, offset: 12849},
				run: (*parser).callonNoMaxValue1,
				expr: &seqExpr{
					pos: position{line: 346, col: 15, offset: 12849},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 346, col: 15, offset: 12849},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 346, col: 18, offset: 12852},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 346, col: 24, offset: 12858},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 346, col: 27, offset: 12861},
							val:        "maxvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "Start",
			pos:  position{line: 350, col: 1, offset: 12898},
			expr: &actionExpr{
				pos: position{line: 350, col: 10, offset: 12907},
				run: (*parser).callonStart1,
				expr: &seqExpr{
					pos: position{line: 350, col: 10, offset: 12907},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 350, col: 10, offset: 12907},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 350, col: 13, offset: 12910},
							val:        "start",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 350, col: 22, offset: 12919},
							expr: &seqExpr{
								pos: position{line: 350, col: 23, offset: 12920},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 350, col: 23, offset: 12920},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 350, col: 26, offset: 12923},
										val:        "with",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 350, col: 36, offset: 12933},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 350, col: 39, offset: 12936},
							label: "start",
							expr: &ruleRefExpr{
								pos:  position{line: 350, col: 45, offset: 12942},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cache",
			pos:  position{line: 356, col: 1, offset: 13054},
			expr: &actionExpr{
				pos: position{line: 356, col: 10, offset: 13063},
				run: (*parser).callonCache1,
				expr: &seqExpr{
					pos: position{line: 356, col: 10, offset: 13063},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 356, col: 10, offset: 13063},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 356, col: 13, offset: 13066},
							val:        "cache",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 356, col: 22, offset: 13075},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 356, col: 25, offset: 13078},
							label: "cache",
							expr: &ruleRefExpr{
								pos:  position{line: 356, col: 31, offset: 13084},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cycle",
			pos:  position{line: 362, col: 1, offset: 13196},
			expr: &actionExpr{
				pos: position{line: 362, col: 10, offset: 13205},
				run: (*parser).callonCycle1,
				expr: &seqExpr{
					pos: position{line: 362, col: 10, offset: 13205},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 362, col: 10, offset: 13205},
							label: "no",
							expr: &zeroOrOneExpr{
								pos: position{line: 362, col: 13, offset: 13208},
								expr: &seqExpr{
									pos: position{line: 362, col: 14, offset: 13209},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 362, col: 14, offset: 13209},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 362, col: 17, offset: 13212},
											val:        "no",
											ignoreCase: true,
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 362, col: 25, offset: 13220},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 362, col: 28, offset: 13223},
							val:        "cycle",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OwnedBy",
			pos:  position{line: 373, col: 1, offset: 13407},
			expr: &actionExpr{
				pos: position{line: 373, col: 12, offset: 13418},
				run: (*parser).callonOwnedBy1,
				expr: &seqExpr{
					pos: position{line: 373, col: 12, offset: 13418},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 373, col: 12, offset: 13418},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 373, col: 15, offset: 13421},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 373, col: 24, offset: 13430},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 373, col: 27, offset: 13433},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 373, col: 33, offset: 13439},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 373, col: 36, offset: 13442},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 373, col: 43, offset: 13449},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 373, col: 43, offset: 13449},
										val:        "none",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 373, col: 53, offset: 13459},
										name: "TableDotCol",
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
			pos:  position{line: 391, col: 1, offset: 14918},
			expr: &actionExpr{
				pos: position{line: 391, col: 19, offset: 14936},
				run: (*parser).callonCreateTypeStmt1,
				expr: &seqExpr{
					pos: position{line: 391, col: 19, offset: 14936},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 391, col: 19, offset: 14936},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 391, col: 29, offset: 14946},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 391, col: 32, offset: 14949},
							val:        "type",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 391, col: 40, offset: 14957},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 391, col: 43, offset: 14960},
							label: "typename",
							expr: &ruleRefExpr{
								pos:  position{line: 391, col: 52, offset: 14969},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 391, col: 58, offset: 14975},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 391, col: 61, offset: 14978},
							val:        "as",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 391, col: 67, offset: 14984},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 391, col: 70, offset: 14987},
							label: "typedef",
							expr: &ruleRefExpr{
								pos:  position{line: 391, col: 78, offset: 14995},
								name: "EnumDef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 391, col: 86, offset: 15003},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 391, col: 88, offset: 15005},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 391, col: 92, offset: 15009},
							expr: &ruleRefExpr{
								pos:  position{line: 391, col: 92, offset: 15009},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "EnumDef",
			pos:  position{line: 397, col: 1, offset: 15125},
			expr: &actionExpr{
				pos: position{line: 397, col: 12, offset: 15136},
				run: (*parser).callonEnumDef1,
				expr: &seqExpr{
					pos: position{line: 397, col: 12, offset: 15136},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 397, col: 12, offset: 15136},
							val:        "ENUM",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 397, col: 19, offset: 15143},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 397, col: 21, offset: 15145},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 397, col: 25, offset: 15149},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 397, col: 27, offset: 15151},
							label: "vals",
							expr: &seqExpr{
								pos: position{line: 397, col: 34, offset: 15158},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 397, col: 34, offset: 15158},
										name: "StringConst",
									},
									&zeroOrMoreExpr{
										pos: position{line: 397, col: 46, offset: 15170},
										expr: &seqExpr{
											pos: position{line: 397, col: 48, offset: 15172},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 397, col: 48, offset: 15172},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 397, col: 50, offset: 15174},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 397, col: 54, offset: 15178},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 397, col: 56, offset: 15180},
													name: "StringConst",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 397, col: 74, offset: 15198},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 397, col: 76, offset: 15200},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AlterTableStmt",
			pos:  position{line: 422, col: 1, offset: 16830},
			expr: &actionExpr{
				pos: position{line: 422, col: 19, offset: 16848},
				run: (*parser).callonAlterTableStmt1,
				expr: &seqExpr{
					pos: position{line: 422, col: 19, offset: 16848},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 422, col: 19, offset: 16848},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 28, offset: 16857},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 422, col: 31, offset: 16860},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 40, offset: 16869},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 422, col: 43, offset: 16872},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 422, col: 48, offset: 16877},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 54, offset: 16883},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 422, col: 57, offset: 16886},
							val:        "owner",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 66, offset: 16895},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 422, col: 69, offset: 16898},
							val:        "to",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 75, offset: 16904},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 422, col: 78, offset: 16907},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 422, col: 84, offset: 16913},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 90, offset: 16919},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 422, col: 92, offset: 16921},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 422, col: 96, offset: 16925},
							expr: &ruleRefExpr{
								pos:  position{line: 422, col: 96, offset: 16925},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "AlterSeqStmt",
			pos:  position{line: 436, col: 1, offset: 18071},
			expr: &actionExpr{
				pos: position{line: 436, col: 17, offset: 18087},
				run: (*parser).callonAlterSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 436, col: 17, offset: 18087},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 436, col: 17, offset: 18087},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 436, col: 26, offset: 18096},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 436, col: 29, offset: 18099},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 436, col: 41, offset: 18111},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 436, col: 44, offset: 18114},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 436, col: 49, offset: 18119},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 436, col: 55, offset: 18125},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 436, col: 58, offset: 18128},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 436, col: 67, offset: 18137},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 436, col: 70, offset: 18140},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 436, col: 76, offset: 18146},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 436, col: 79, offset: 18149},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 436, col: 85, offset: 18155},
								name: "TableDotCol",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 436, col: 97, offset: 18167},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 436, col: 99, offset: 18169},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 436, col: 103, offset: 18173},
							expr: &ruleRefExpr{
								pos:  position{line: 436, col: 103, offset: 18173},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "TableDotCol",
			pos:  position{line: 440, col: 1, offset: 18252},
			expr: &actionExpr{
				pos: position{line: 440, col: 16, offset: 18267},
				run: (*parser).callonTableDotCol1,
				expr: &seqExpr{
					pos: position{line: 440, col: 16, offset: 18267},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 440, col: 16, offset: 18267},
							label: "table",
							expr: &ruleRefExpr{
								pos:  position{line: 440, col: 22, offset: 18273},
								name: "Ident",
							},
						},
						&litMatcher{
							pos:        position{line: 440, col: 28, offset: 18279},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 440, col: 32, offset: 18283},
							label: "column",
							expr: &ruleRefExpr{
								pos:  position{line: 440, col: 39, offset: 18290},
								name: "Ident",
							},
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 454, col: 1, offset: 19618},
			expr: &actionExpr{
				pos: position{line: 454, col: 25, offset: 19642},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 454, col: 25, offset: 19642},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 454, col: 25, offset: 19642},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 454, col: 36, offset: 19653},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 454, col: 39, offset: 19656},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 454, col: 45, offset: 19662},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 454, col: 48, offset: 19665},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 454, col: 61, offset: 19678},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 454, col: 63, offset: 19680},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 454, col: 73, offset: 19690},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 454, col: 79, offset: 19696},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 454, col: 81, offset: 19698},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 454, col: 87, offset: 19704},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 454, col: 89, offset: 19706},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 454, col: 97, offset: 19714},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 454, col: 109, offset: 19726},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 454, col: 111, offset: 19728},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 454, col: 115, offset: 19732},
							expr: &ruleRefExpr{
								pos:  position{line: 454, col: 115, offset: 19732},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 458, col: 1, offset: 19821},
			expr: &actionExpr{
				pos: position{line: 458, col: 24, offset: 19844},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 458, col: 24, offset: 19844},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 458, col: 24, offset: 19844},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 458, col: 34, offset: 19854},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 458, col: 37, offset: 19857},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 458, col: 50, offset: 19870},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 458, col: 53, offset: 19873},
							expr: &seqExpr{
								pos: position{line: 458, col: 55, offset: 19875},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 458, col: 55, offset: 19875},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 458, col: 61, offset: 19881},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 458, col: 64, offset: 19884},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 458, col: 71, offset: 19891},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 458, col: 74, offset: 19894},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 458, col: 84, offset: 19904},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 458, col: 90, offset: 19910},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 458, col: 100, offset: 19920},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 458, col: 106, offset: 19926},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 458, col: 109, offset: 19929},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 458, col: 117, offset: 19937},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 458, col: 120, offset: 19940},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 458, col: 130, offset: 19950},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 458, col: 133, offset: 19953},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 458, col: 140, offset: 19960},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 458, col: 146, offset: 19966},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 458, col: 148, offset: 19968},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 458, col: 152, offset: 19972},
							expr: &ruleRefExpr{
								pos:  position{line: 458, col: 152, offset: 19972},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 462, col: 1, offset: 20063},
			expr: &actionExpr{
				pos: position{line: 462, col: 12, offset: 20074},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 462, col: 12, offset: 20074},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 462, col: 12, offset: 20074},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 462, col: 19, offset: 20081},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 462, col: 21, offset: 20083},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 462, col: 25, offset: 20087},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 462, col: 29, offset: 20091},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 462, col: 33, offset: 20095},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 462, col: 33, offset: 20095},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 462, col: 39, offset: 20101},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 462, col: 47, offset: 20109},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 462, col: 49, offset: 20111},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 462, col: 56, offset: 20118},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 462, col: 77, offset: 20139},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 462, col: 79, offset: 20141},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 462, col: 83, offset: 20145},
							expr: &ruleRefExpr{
								pos:  position{line: 462, col: 83, offset: 20145},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 467, col: 1, offset: 20227},
			expr: &actionExpr{
				pos: position{line: 467, col: 8, offset: 20234},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 467, col: 8, offset: 20234},
					expr: &charClassMatcher{
						pos:        position{line: 467, col: 8, offset: 20234},
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
			pos:  position{line: 482, col: 1, offset: 21074},
			expr: &actionExpr{
				pos: position{line: 482, col: 25, offset: 21098},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 482, col: 25, offset: 21098},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 482, col: 32, offset: 21105},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 482, col: 32, offset: 21105},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 482, col: 38, offset: 21111},
								expr: &seqExpr{
									pos: position{line: 482, col: 40, offset: 21113},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 482, col: 40, offset: 21113},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 482, col: 42, offset: 21115},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 482, col: 46, offset: 21119},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 482, col: 48, offset: 21121},
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
			pos:  position{line: 494, col: 1, offset: 21411},
			expr: &choiceExpr{
				pos: position{line: 494, col: 12, offset: 21422},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 494, col: 12, offset: 21422},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 494, col: 21, offset: 21431},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 494, col: 31, offset: 21441},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 494, col: 45, offset: 21455},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 496, col: 1, offset: 21464},
			expr: &actionExpr{
				pos: position{line: 496, col: 16, offset: 21479},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 496, col: 16, offset: 21479},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 496, col: 16, offset: 21479},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 496, col: 20, offset: 21483},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 496, col: 26, offset: 21489},
								expr: &choiceExpr{
									pos: position{line: 496, col: 27, offset: 21490},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 496, col: 27, offset: 21490},
											val:        "[^'\\n]",
											chars:      []rune{'\'', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 496, col: 36, offset: 21499},
											val:        "''",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 496, col: 43, offset: 21506},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "DblQuotedString",
			pos:  position{line: 500, col: 1, offset: 21558},
			expr: &actionExpr{
				pos: position{line: 500, col: 20, offset: 21577},
				run: (*parser).callonDblQuotedString1,
				expr: &seqExpr{
					pos: position{line: 500, col: 20, offset: 21577},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 500, col: 20, offset: 21577},
							val:        "\"",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 500, col: 24, offset: 21581},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 500, col: 30, offset: 21587},
								expr: &choiceExpr{
									pos: position{line: 500, col: 31, offset: 21588},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 500, col: 31, offset: 21588},
											val:        "[^\"\\n]",
											chars:      []rune{'"', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 500, col: 40, offset: 21597},
											val:        "\"\"",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 500, col: 49, offset: 21606},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 504, col: 1, offset: 21670},
			expr: &actionExpr{
				pos: position{line: 504, col: 10, offset: 21679},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 504, col: 10, offset: 21679},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 504, col: 10, offset: 21679},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 504, col: 18, offset: 21687},
							expr: &charClassMatcher{
								pos:        position{line: 504, col: 18, offset: 21687},
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
			pos:  position{line: 508, col: 1, offset: 21740},
			expr: &actionExpr{
				pos: position{line: 508, col: 11, offset: 21750},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 508, col: 13, offset: 21752},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 508, col: 13, offset: 21752},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 508, col: 19, offset: 21758},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 508, col: 19, offset: 21758},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 508, col: 24, offset: 21763},
									expr: &charClassMatcher{
										pos:        position{line: 508, col: 24, offset: 21763},
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
			pos:  position{line: 513, col: 1, offset: 21858},
			expr: &actionExpr{
				pos: position{line: 513, col: 15, offset: 21872},
				run: (*parser).callonNonZNumber1,
				expr: &seqExpr{
					pos: position{line: 513, col: 15, offset: 21872},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 513, col: 15, offset: 21872},
							val:        "[1-9]",
							ranges:     []rune{'1', '9'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 513, col: 20, offset: 21877},
							expr: &charClassMatcher{
								pos:        position{line: 513, col: 20, offset: 21877},
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
			pos:  position{line: 518, col: 1, offset: 21970},
			expr: &actionExpr{
				pos: position{line: 518, col: 12, offset: 21981},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 518, col: 12, offset: 21981},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 518, col: 20, offset: 21989},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 518, col: 20, offset: 21989},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 518, col: 34, offset: 22003},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 522, col: 1, offset: 22045},
			expr: &actionExpr{
				pos: position{line: 522, col: 16, offset: 22060},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 522, col: 18, offset: 22062},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 522, col: 18, offset: 22062},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 522, col: 27, offset: 22071},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 522, col: 27, offset: 22071},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 522, col: 31, offset: 22075},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 522, col: 49, offset: 22093},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 522, col: 55, offset: 22099},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 526, col: 1, offset: 22145},
			expr: &choiceExpr{
				pos: position{line: 526, col: 24, offset: 22168},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 526, col: 24, offset: 22168},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 526, col: 33, offset: 22177},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 526, col: 41, offset: 22185},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 526, col: 48, offset: 22192},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 526, col: 54, offset: 22198},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 528, col: 1, offset: 22205},
			expr: &actionExpr{
				pos: position{line: 528, col: 17, offset: 22221},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 528, col: 19, offset: 22223},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 528, col: 19, offset: 22223},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 528, col: 29, offset: 22233},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 528, col: 29, offset: 22233},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 528, col: 33, offset: 22237},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 528, col: 52, offset: 22256},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 528, col: 58, offset: 22262},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 532, col: 1, offset: 22310},
			expr: &choiceExpr{
				pos: position{line: 532, col: 25, offset: 22334},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 532, col: 25, offset: 22334},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 532, col: 35, offset: 22344},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 532, col: 42, offset: 22351},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 532, col: 50, offset: 22359},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 532, col: 56, offset: 22365},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 544, col: 1, offset: 22880},
			expr: &actionExpr{
				pos: position{line: 544, col: 12, offset: 22891},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 544, col: 14, offset: 22893},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 544, col: 14, offset: 22893},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 544, col: 34, offset: 22913},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 548, col: 1, offset: 22957},
			expr: &seqExpr{
				pos: position{line: 548, col: 21, offset: 22977},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 548, col: 21, offset: 22977},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 548, col: 26, offset: 22982},
						expr: &anyMatcher{
							line: 548, col: 26, offset: 22982,
						},
					},
					&litMatcher{
						pos:        position{line: 548, col: 29, offset: 22985},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 548, col: 34, offset: 22990},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 550, col: 1, offset: 22995},
			expr: &seqExpr{
				pos: position{line: 550, col: 22, offset: 23016},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 550, col: 22, offset: 23016},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 550, col: 27, offset: 23021},
						expr: &charClassMatcher{
							pos:        position{line: 550, col: 27, offset: 23021},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 550, col: 36, offset: 23030},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 552, col: 1, offset: 23035},
			expr: &seqExpr{
				pos: position{line: 552, col: 9, offset: 23043},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 552, col: 9, offset: 23043},
						expr: &charClassMatcher{
							pos:        position{line: 552, col: 9, offset: 23043},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 552, col: 17, offset: 23051},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 552, col: 17, offset: 23051},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 552, col: 26, offset: 23060},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 552, col: 35, offset: 23069},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 552, col: 42, offset: 23076},
								val:        "\n",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 554, col: 1, offset: 23083},
			expr: &zeroOrMoreExpr{
				pos: position{line: 554, col: 19, offset: 23101},
				expr: &charClassMatcher{
					pos:        position{line: 554, col: 19, offset: 23101},
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
			pos:         position{line: 556, col: 1, offset: 23113},
			expr: &oneOrMoreExpr{
				pos: position{line: 556, col: 31, offset: 23143},
				expr: &charClassMatcher{
					pos:        position{line: 556, col: 31, offset: 23143},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 558, col: 1, offset: 23155},
			expr: &notExpr{
				pos: position{line: 558, col: 8, offset: 23162},
				expr: &anyMatcher{
					line: 558, col: 9, offset: 23163,
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

func (c *current) onCreateTableStmt1(tablename, defs interface{}) (interface{}, error) {
	defsSlice := []map[string]string{}
	valsSlice := toIfaceSlice(defs)
	if valsSlice[0] == nil {
		defsSlice = append(defsSlice, nil)
	} else {
		defsSlice = append(defsSlice, valsSlice[0].(map[string]string))
	}
	restSlice := toIfaceSlice(valsSlice[1])
	for _, v := range restSlice {
		vSlice := toIfaceSlice(v)
		if vSlice[3] == nil {
			defsSlice = append(defsSlice, nil)
		} else {
			defsSlice = append(defsSlice, vSlice[3].(map[string]string))
		}
	}
	return parseCreateTableStmt(tablename, defsSlice)
}

func (p *parser) callonCreateTableStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCreateTableStmt1(stack["tablename"], stack["defs"])
}

func (c *current) onColumnDef1(name, dataType, constraint interface{}) (interface{}, error) {
	if dataType == nil {
		return nil, nil
	}
	result := make(map[string]string)
	if err := mergo.Merge(&result, dataType.(map[string]string), mergo.WithOverride); err != nil {
		return nil, err
	}
	if constraint != nil {
		if err := mergo.Merge(&result, constraint.(map[string]string), mergo.WithOverride); err != nil {
			return nil, err
		}
	}
	result["name"] = interfaceToString(name)
	return result, nil
}

func (p *parser) callonColumnDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onColumnDef1(stack["name"], stack["dataType"], stack["constraint"])
}

func (c *current) onColumnConstraint1(nameOpt, constraint interface{}) (interface{}, error) {
	properties := make(map[string]string)

	if err := mergo.Merge(&properties, constraint.(map[string]string), mergo.WithOverride); err != nil {
		return nil, err
	}

	nameSlice := toIfaceSlice(nameOpt)
	if nameSlice != nil {
		properties["constraint_name"] = interfaceToString(nameSlice[3])
	}

	return properties, nil
}

func (p *parser) callonColumnConstraint1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onColumnConstraint1(stack["nameOpt"], stack["constraint"])
}

func (c *current) onTableConstr1(nameOpt, constraint interface{}) (interface{}, error) {
	properties := map[string]string{
		"table_constraint": "true",
	}

	if err := mergo.Merge(&properties, constraint.(map[string]string), mergo.WithOverride); err != nil {
		return nil, err
	}

	nameSlice := toIfaceSlice(nameOpt)
	if nameSlice != nil {
		properties["constraint_name"] = interfaceToString(nameSlice[2])
	}

	return properties, nil
}

func (p *parser) callonTableConstr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTableConstr1(stack["nameOpt"], stack["constraint"])
}

func (c *current) onNotNullCls1() (interface{}, error) {
	return map[string]string{
		"not_null": "true",
	}, nil
}

func (p *parser) callonNotNullCls1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNotNullCls1()
}

func (c *current) onNullCls1() (interface{}, error) {
	return map[string]string{
		"not_null": "false",
	}, nil
}

func (p *parser) callonNullCls1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNullCls1()
}

func (c *current) onCheckCls1(expr, noInherit interface{}) (interface{}, error) {
	result := map[string]string{
		"check_def": expr.(string),
	}
	if noInherit != nil {
		result["check_no_inherit"] = "true"
	}
	return result, nil
}

func (p *parser) callonCheckCls1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCheckCls1(stack["expr"], stack["noInherit"])
}

func (c *current) onWrappedExpr1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonWrappedExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onWrappedExpr1()
}

func (c *current) onDataType1(t, brackets interface{}) (interface{}, error) {
	bracketsSlice := toIfaceSlice(brackets)
	var result map[string]string
	if t != nil {
		result = t.(map[string]string)
		if l := len(bracketsSlice); l > 0 {
			result["array_dimensions"] = strconv.Itoa(l)
		}
		return result, nil
	}
	return nil, nil
}

func (p *parser) callonDataType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDataType1(stack["t"], stack["brackets"])
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

func (c *current) onNumericT1(args interface{}) (interface{}, error) {
	result := map[string]string{
		"type": "numeric",
	}
	if args != nil {
		argsSlice := toIfaceSlice(args)
		result["precision"] = strconv.FormatInt(argsSlice[1].(int64), 10)
		if argsSlice[2] != nil {
			slice := toIfaceSlice(argsSlice[2])
			result["scale"] = strconv.FormatInt(slice[2].(int64), 10)
		}
	}
	return result, nil
}

func (p *parser) callonNumericT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumericT1(stack["args"])
}

func (c *current) onPostgisT1(t, subtype, srid interface{}) (interface{}, error) {
	result := map[string]string{
		"type":    strings.ToLower(string(t.([]byte))),
		"subtype": strings.ToLower(string(subtype.([]byte))),
	}
	if srid != nil {
		sridSlice := toIfaceSlice(srid)
		result["srid"] = strconv.FormatInt(sridSlice[1].(int64), 10)
	}
	return result, nil
}

func (p *parser) callonPostgisT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPostgisT1(stack["t"], stack["subtype"], stack["srid"])
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
	return map[string]string{
		"owned_by": name.(string),
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

func (c *current) onAlterSeqStmt1(name, owner interface{}) (interface{}, error) {
	return parseAlterSequenceStmt(name.(Identifier), owner.(string))
}

func (p *parser) callonAlterSeqStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAlterSeqStmt1(stack["name"], stack["owner"])
}

func (c *current) onTableDotCol1(table, column interface{}) (interface{}, error) {
	return parseTableDotColumn(table.(Identifier), column.(Identifier)), nil
}

func (p *parser) callonTableDotCol1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTableDotCol1(stack["table"], stack["column"])
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

func (c *current) onDblQuotedString1(value interface{}) (interface{}, error) {
	return DoubleQuotedString(toByteSlice(value)), nil
}

func (p *parser) callonDblQuotedString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDblQuotedString1(stack["value"])
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
