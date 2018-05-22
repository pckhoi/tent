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
							expr: &ruleRefExpr{
								pos:  position{line: 51, col: 19, offset: 2705},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 51, col: 25, offset: 2711},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 51, col: 28, offset: 2714},
							label: "dataType",
							expr: &ruleRefExpr{
								pos:  position{line: 51, col: 37, offset: 2723},
								name: "DataType",
							},
						},
						&labeledExpr{
							pos:   position{line: 51, col: 46, offset: 2732},
							label: "constraint",
							expr: &zeroOrOneExpr{
								pos: position{line: 51, col: 57, offset: 2743},
								expr: &ruleRefExpr{
									pos:  position{line: 51, col: 57, offset: 2743},
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
			pos:  position{line: 68, col: 1, offset: 3232},
			expr: &actionExpr{
				pos: position{line: 68, col: 21, offset: 3252},
				run: (*parser).callonColumnConstraint1,
				expr: &seqExpr{
					pos: position{line: 68, col: 21, offset: 3252},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 68, col: 21, offset: 3252},
							label: "nameOpt",
							expr: &zeroOrOneExpr{
								pos: position{line: 68, col: 29, offset: 3260},
								expr: &seqExpr{
									pos: position{line: 68, col: 31, offset: 3262},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 68, col: 31, offset: 3262},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 68, col: 34, offset: 3265},
											val:        "constraint",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 68, col: 48, offset: 3279},
											name: "_1",
										},
										&choiceExpr{
											pos: position{line: 68, col: 52, offset: 3283},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 68, col: 52, offset: 3283},
													name: "StringConst",
												},
												&ruleRefExpr{
													pos:  position{line: 68, col: 66, offset: 3297},
													name: "Ident",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 68, col: 76, offset: 3307},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 68, col: 78, offset: 3309},
							label: "constraint",
							expr: &choiceExpr{
								pos: position{line: 68, col: 91, offset: 3322},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 68, col: 91, offset: 3322},
										name: "NotNullCls",
									},
									&ruleRefExpr{
										pos:  position{line: 68, col: 104, offset: 3335},
										name: "NullCls",
									},
									&ruleRefExpr{
										pos:  position{line: 68, col: 114, offset: 3345},
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
			pos:  position{line: 83, col: 1, offset: 3711},
			expr: &actionExpr{
				pos: position{line: 83, col: 16, offset: 3726},
				run: (*parser).callonTableConstr1,
				expr: &seqExpr{
					pos: position{line: 83, col: 16, offset: 3726},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 83, col: 16, offset: 3726},
							label: "nameOpt",
							expr: &zeroOrOneExpr{
								pos: position{line: 83, col: 24, offset: 3734},
								expr: &seqExpr{
									pos: position{line: 83, col: 26, offset: 3736},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 83, col: 26, offset: 3736},
											val:        "constraint",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 83, col: 40, offset: 3750},
											name: "_1",
										},
										&choiceExpr{
											pos: position{line: 83, col: 44, offset: 3754},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 83, col: 44, offset: 3754},
													name: "StringConst",
												},
												&ruleRefExpr{
													pos:  position{line: 83, col: 58, offset: 3768},
													name: "Ident",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 83, col: 68, offset: 3778},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 83, col: 70, offset: 3780},
							label: "constraint",
							expr: &ruleRefExpr{
								pos:  position{line: 83, col: 81, offset: 3791},
								name: "CheckCls",
							},
						},
					},
				},
			},
		},
		{
			name: "NotNullCls",
			pos:  position{line: 100, col: 1, offset: 4192},
			expr: &actionExpr{
				pos: position{line: 100, col: 15, offset: 4206},
				run: (*parser).callonNotNullCls1,
				expr: &seqExpr{
					pos: position{line: 100, col: 15, offset: 4206},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 100, col: 15, offset: 4206},
							val:        "not",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 22, offset: 4213},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 100, col: 25, offset: 4216},
							val:        "null",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "NullCls",
			pos:  position{line: 106, col: 1, offset: 4298},
			expr: &actionExpr{
				pos: position{line: 106, col: 12, offset: 4309},
				run: (*parser).callonNullCls1,
				expr: &litMatcher{
					pos:        position{line: 106, col: 12, offset: 4309},
					val:        "null",
					ignoreCase: true,
				},
			},
		},
		{
			name: "CheckCls",
			pos:  position{line: 112, col: 1, offset: 4392},
			expr: &actionExpr{
				pos: position{line: 112, col: 13, offset: 4404},
				run: (*parser).callonCheckCls1,
				expr: &seqExpr{
					pos: position{line: 112, col: 13, offset: 4404},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 112, col: 13, offset: 4404},
							val:        "check",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 112, col: 22, offset: 4413},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 112, col: 25, offset: 4416},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 112, col: 30, offset: 4421},
								name: "WrappedExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 112, col: 42, offset: 4433},
							label: "noInherit",
							expr: &zeroOrOneExpr{
								pos: position{line: 112, col: 52, offset: 4443},
								expr: &seqExpr{
									pos: position{line: 112, col: 54, offset: 4445},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 112, col: 54, offset: 4445},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 112, col: 57, offset: 4448},
											val:        "no",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 112, col: 63, offset: 4454},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 112, col: 66, offset: 4457},
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
			pos:  position{line: 122, col: 1, offset: 4650},
			expr: &actionExpr{
				pos: position{line: 122, col: 16, offset: 4665},
				run: (*parser).callonWrappedExpr1,
				expr: &seqExpr{
					pos: position{line: 122, col: 16, offset: 4665},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 122, col: 16, offset: 4665},
							val:        "(",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 122, col: 20, offset: 4669},
							expr: &ruleRefExpr{
								pos:  position{line: 122, col: 20, offset: 4669},
								name: "Expr",
							},
						},
						&litMatcher{
							pos:        position{line: 122, col: 26, offset: 4675},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 126, col: 1, offset: 4715},
			expr: &choiceExpr{
				pos: position{line: 126, col: 9, offset: 4723},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 126, col: 9, offset: 4723},
						name: "WrappedExpr",
					},
					&oneOrMoreExpr{
						pos: position{line: 126, col: 23, offset: 4737},
						expr: &charClassMatcher{
							pos:        position{line: 126, col: 23, offset: 4737},
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
			pos:  position{line: 138, col: 1, offset: 5757},
			expr: &actionExpr{
				pos: position{line: 138, col: 13, offset: 5769},
				run: (*parser).callonDataType1,
				expr: &seqExpr{
					pos: position{line: 138, col: 13, offset: 5769},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 138, col: 13, offset: 5769},
							label: "t",
							expr: &choiceExpr{
								pos: position{line: 138, col: 17, offset: 5773},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 138, col: 17, offset: 5773},
										name: "TimestampT",
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 30, offset: 5786},
										name: "TimeT",
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 38, offset: 5794},
										name: "VarcharT",
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 49, offset: 5805},
										name: "CharT",
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 57, offset: 5813},
										name: "BitVarT",
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 67, offset: 5823},
										name: "BitT",
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 74, offset: 5830},
										name: "IntT",
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 81, offset: 5837},
										name: "PgOidT",
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 90, offset: 5846},
										name: "OtherT",
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 99, offset: 5855},
										name: "CustomT",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 138, col: 109, offset: 5865},
							label: "brackets",
							expr: &zeroOrMoreExpr{
								pos: position{line: 138, col: 118, offset: 5874},
								expr: &litMatcher{
									pos:        position{line: 138, col: 120, offset: 5876},
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
			pos:  position{line: 151, col: 1, offset: 6187},
			expr: &actionExpr{
				pos: position{line: 151, col: 15, offset: 6201},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 151, col: 15, offset: 6201},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 151, col: 15, offset: 6201},
							val:        "timestamp",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 151, col: 28, offset: 6214},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 151, col: 33, offset: 6219},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 151, col: 46, offset: 6232},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 151, col: 59, offset: 6245},
								expr: &choiceExpr{
									pos: position{line: 151, col: 61, offset: 6247},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 151, col: 61, offset: 6247},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 151, col: 70, offset: 6256},
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
			pos:  position{line: 164, col: 1, offset: 6535},
			expr: &actionExpr{
				pos: position{line: 164, col: 10, offset: 6544},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 164, col: 10, offset: 6544},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 164, col: 10, offset: 6544},
							val:        "time",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 164, col: 18, offset: 6552},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 164, col: 23, offset: 6557},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 164, col: 36, offset: 6570},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 164, col: 49, offset: 6583},
								expr: &choiceExpr{
									pos: position{line: 164, col: 51, offset: 6585},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 164, col: 51, offset: 6585},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 164, col: 60, offset: 6594},
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
			pos:  position{line: 177, col: 1, offset: 6865},
			expr: &actionExpr{
				pos: position{line: 177, col: 17, offset: 6881},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 177, col: 17, offset: 6881},
					expr: &seqExpr{
						pos: position{line: 177, col: 19, offset: 6883},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 177, col: 19, offset: 6883},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 177, col: 22, offset: 6886},
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
			pos:  position{line: 184, col: 1, offset: 7014},
			expr: &actionExpr{
				pos: position{line: 184, col: 11, offset: 7024},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 184, col: 11, offset: 7024},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 184, col: 11, offset: 7024},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 184, col: 14, offset: 7027},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 184, col: 22, offset: 7035},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 184, col: 25, offset: 7038},
							val:        "time",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 184, col: 33, offset: 7046},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 184, col: 36, offset: 7049},
							val:        "zone",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 188, col: 1, offset: 7083},
			expr: &actionExpr{
				pos: position{line: 188, col: 14, offset: 7096},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 188, col: 14, offset: 7096},
					expr: &seqExpr{
						pos: position{line: 188, col: 16, offset: 7098},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 188, col: 16, offset: 7098},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 188, col: 19, offset: 7101},
								val:        "without",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 188, col: 30, offset: 7112},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 188, col: 33, offset: 7115},
								val:        "time",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 188, col: 41, offset: 7123},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 188, col: 44, offset: 7126},
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
			pos:  position{line: 192, col: 1, offset: 7164},
			expr: &actionExpr{
				pos: position{line: 192, col: 10, offset: 7173},
				run: (*parser).callonCharT1,
				expr: &seqExpr{
					pos: position{line: 192, col: 10, offset: 7173},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 192, col: 12, offset: 7175},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 192, col: 12, offset: 7175},
									val:        "character",
									ignoreCase: true,
								},
								&litMatcher{
									pos:        position{line: 192, col: 27, offset: 7190},
									val:        "char",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 192, col: 37, offset: 7200},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 192, col: 44, offset: 7207},
								expr: &seqExpr{
									pos: position{line: 192, col: 46, offset: 7209},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 192, col: 46, offset: 7209},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 192, col: 50, offset: 7213},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 192, col: 61, offset: 7224},
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
			pos:  position{line: 204, col: 1, offset: 7479},
			expr: &actionExpr{
				pos: position{line: 204, col: 13, offset: 7491},
				run: (*parser).callonVarcharT1,
				expr: &seqExpr{
					pos: position{line: 204, col: 13, offset: 7491},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 204, col: 15, offset: 7493},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 204, col: 17, offset: 7495},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 204, col: 17, offset: 7495},
											val:        "character",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 204, col: 30, offset: 7508},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 204, col: 33, offset: 7511},
											val:        "varying",
											ignoreCase: true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 204, col: 48, offset: 7526},
									val:        "varchar",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 204, col: 61, offset: 7539},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 204, col: 68, offset: 7546},
								expr: &seqExpr{
									pos: position{line: 204, col: 70, offset: 7548},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 204, col: 70, offset: 7548},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 204, col: 74, offset: 7552},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 204, col: 85, offset: 7563},
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
			pos:  position{line: 215, col: 1, offset: 7798},
			expr: &actionExpr{
				pos: position{line: 215, col: 9, offset: 7806},
				run: (*parser).callonBitT1,
				expr: &seqExpr{
					pos: position{line: 215, col: 9, offset: 7806},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 215, col: 9, offset: 7806},
							val:        "bit",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 215, col: 16, offset: 7813},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 215, col: 23, offset: 7820},
								expr: &seqExpr{
									pos: position{line: 215, col: 25, offset: 7822},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 215, col: 25, offset: 7822},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 215, col: 29, offset: 7826},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 215, col: 40, offset: 7837},
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
			pos:  position{line: 227, col: 1, offset: 8091},
			expr: &actionExpr{
				pos: position{line: 227, col: 12, offset: 8102},
				run: (*parser).callonBitVarT1,
				expr: &seqExpr{
					pos: position{line: 227, col: 12, offset: 8102},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 227, col: 12, offset: 8102},
							val:        "bit",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 227, col: 19, offset: 8109},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 227, col: 22, offset: 8112},
							val:        "varying",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 227, col: 33, offset: 8123},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 227, col: 40, offset: 8130},
								expr: &seqExpr{
									pos: position{line: 227, col: 42, offset: 8132},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 227, col: 42, offset: 8132},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 227, col: 46, offset: 8136},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 227, col: 57, offset: 8147},
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
			pos:  position{line: 238, col: 1, offset: 8381},
			expr: &actionExpr{
				pos: position{line: 238, col: 9, offset: 8389},
				run: (*parser).callonIntT1,
				expr: &choiceExpr{
					pos: position{line: 238, col: 11, offset: 8391},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 238, col: 11, offset: 8391},
							val:        "integer",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 238, col: 24, offset: 8404},
							val:        "int",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "PgOidT",
			pos:  position{line: 244, col: 1, offset: 8486},
			expr: &actionExpr{
				pos: position{line: 244, col: 11, offset: 8496},
				run: (*parser).callonPgOidT1,
				expr: &choiceExpr{
					pos: position{line: 244, col: 13, offset: 8498},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 244, col: 13, offset: 8498},
							val:        "oid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 22, offset: 8507},
							val:        "regprocedure",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 40, offset: 8525},
							val:        "regproc",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 53, offset: 8538},
							val:        "regoperator",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 70, offset: 8555},
							val:        "regoper",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 83, offset: 8568},
							val:        "regclass",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 97, offset: 8582},
							val:        "regtype",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 110, offset: 8595},
							val:        "regrole",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 123, offset: 8608},
							val:        "regnamespace",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 141, offset: 8626},
							val:        "regconfig",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 244, col: 156, offset: 8641},
							val:        "regdictionary",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 250, col: 1, offset: 8755},
			expr: &actionExpr{
				pos: position{line: 250, col: 11, offset: 8765},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 250, col: 13, offset: 8767},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 250, col: 13, offset: 8767},
							val:        "date",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 23, offset: 8777},
							val:        "smallint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 37, offset: 8791},
							val:        "bigint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 49, offset: 8803},
							val:        "decimal",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 62, offset: 8816},
							val:        "numeric",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 75, offset: 8829},
							val:        "real",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 85, offset: 8839},
							val:        "smallserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 102, offset: 8856},
							val:        "serial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 114, offset: 8868},
							val:        "bigserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 129, offset: 8883},
							val:        "boolean",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 142, offset: 8896},
							val:        "text",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 152, offset: 8906},
							val:        "money",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 163, offset: 8917},
							val:        "bytea",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 174, offset: 8928},
							val:        "point",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 185, offset: 8939},
							val:        "line",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 195, offset: 8949},
							val:        "lseg",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 205, offset: 8959},
							val:        "box",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 214, offset: 8968},
							val:        "path",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 224, offset: 8978},
							val:        "polygon",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 237, offset: 8991},
							val:        "circle",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 249, offset: 9003},
							val:        "cidr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 259, offset: 9013},
							val:        "inet",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 269, offset: 9023},
							val:        "macaddr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 282, offset: 9036},
							val:        "uuid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 292, offset: 9046},
							val:        "xml",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 301, offset: 9055},
							val:        "jsonb",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 250, col: 312, offset: 9066},
							val:        "json",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CustomT",
			pos:  position{line: 256, col: 1, offset: 9171},
			expr: &actionExpr{
				pos: position{line: 256, col: 13, offset: 9183},
				run: (*parser).callonCustomT1,
				expr: &ruleRefExpr{
					pos:  position{line: 256, col: 13, offset: 9183},
					name: "Ident",
				},
			},
		},
		{
			name: "CreateSeqStmt",
			pos:  position{line: 277, col: 1, offset: 10648},
			expr: &actionExpr{
				pos: position{line: 277, col: 18, offset: 10665},
				run: (*parser).callonCreateSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 277, col: 18, offset: 10665},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 277, col: 18, offset: 10665},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 277, col: 28, offset: 10675},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 277, col: 31, offset: 10678},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 277, col: 43, offset: 10690},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 277, col: 46, offset: 10693},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 277, col: 51, offset: 10698},
								name: "Ident",
							},
						},
						&labeledExpr{
							pos:   position{line: 277, col: 57, offset: 10704},
							label: "verses",
							expr: &zeroOrMoreExpr{
								pos: position{line: 277, col: 64, offset: 10711},
								expr: &ruleRefExpr{
									pos:  position{line: 277, col: 64, offset: 10711},
									name: "CreateSeqVerse",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 277, col: 80, offset: 10727},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 277, col: 82, offset: 10729},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 277, col: 86, offset: 10733},
							expr: &ruleRefExpr{
								pos:  position{line: 277, col: 86, offset: 10733},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateSeqVerse",
			pos:  position{line: 291, col: 1, offset: 11127},
			expr: &actionExpr{
				pos: position{line: 291, col: 19, offset: 11145},
				run: (*parser).callonCreateSeqVerse1,
				expr: &labeledExpr{
					pos:   position{line: 291, col: 19, offset: 11145},
					label: "verse",
					expr: &choiceExpr{
						pos: position{line: 291, col: 27, offset: 11153},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 291, col: 27, offset: 11153},
								name: "IncrementBy",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 41, offset: 11167},
								name: "MinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 52, offset: 11178},
								name: "NoMinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 65, offset: 11191},
								name: "MaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 76, offset: 11202},
								name: "NoMaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 89, offset: 11215},
								name: "Start",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 97, offset: 11223},
								name: "Cache",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 105, offset: 11231},
								name: "Cycle",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 113, offset: 11239},
								name: "OwnedBy",
							},
						},
					},
				},
			},
		},
		{
			name: "IncrementBy",
			pos:  position{line: 295, col: 1, offset: 11276},
			expr: &actionExpr{
				pos: position{line: 295, col: 16, offset: 11291},
				run: (*parser).callonIncrementBy1,
				expr: &seqExpr{
					pos: position{line: 295, col: 16, offset: 11291},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 295, col: 16, offset: 11291},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 295, col: 19, offset: 11294},
							val:        "increment",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 295, col: 32, offset: 11307},
							expr: &seqExpr{
								pos: position{line: 295, col: 33, offset: 11308},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 295, col: 33, offset: 11308},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 295, col: 36, offset: 11311},
										val:        "by",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 295, col: 44, offset: 11319},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 295, col: 47, offset: 11322},
							label: "num",
							expr: &ruleRefExpr{
								pos:  position{line: 295, col: 51, offset: 11326},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "MinValue",
			pos:  position{line: 301, col: 1, offset: 11440},
			expr: &actionExpr{
				pos: position{line: 301, col: 13, offset: 11452},
				run: (*parser).callonMinValue1,
				expr: &seqExpr{
					pos: position{line: 301, col: 13, offset: 11452},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 301, col: 13, offset: 11452},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 301, col: 16, offset: 11455},
							val:        "minvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 301, col: 28, offset: 11467},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 301, col: 31, offset: 11470},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 301, col: 35, offset: 11474},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMinValue",
			pos:  position{line: 307, col: 1, offset: 11587},
			expr: &actionExpr{
				pos: position{line: 307, col: 15, offset: 11601},
				run: (*parser).callonNoMinValue1,
				expr: &seqExpr{
					pos: position{line: 307, col: 15, offset: 11601},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 307, col: 15, offset: 11601},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 307, col: 18, offset: 11604},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 307, col: 24, offset: 11610},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 307, col: 27, offset: 11613},
							val:        "minvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "MaxValue",
			pos:  position{line: 311, col: 1, offset: 11650},
			expr: &actionExpr{
				pos: position{line: 311, col: 13, offset: 11662},
				run: (*parser).callonMaxValue1,
				expr: &seqExpr{
					pos: position{line: 311, col: 13, offset: 11662},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 311, col: 13, offset: 11662},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 311, col: 16, offset: 11665},
							val:        "maxvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 311, col: 28, offset: 11677},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 311, col: 31, offset: 11680},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 311, col: 35, offset: 11684},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMaxValue",
			pos:  position{line: 317, col: 1, offset: 11797},
			expr: &actionExpr{
				pos: position{line: 317, col: 15, offset: 11811},
				run: (*parser).callonNoMaxValue1,
				expr: &seqExpr{
					pos: position{line: 317, col: 15, offset: 11811},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 317, col: 15, offset: 11811},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 317, col: 18, offset: 11814},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 24, offset: 11820},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 317, col: 27, offset: 11823},
							val:        "maxvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "Start",
			pos:  position{line: 321, col: 1, offset: 11860},
			expr: &actionExpr{
				pos: position{line: 321, col: 10, offset: 11869},
				run: (*parser).callonStart1,
				expr: &seqExpr{
					pos: position{line: 321, col: 10, offset: 11869},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 321, col: 10, offset: 11869},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 321, col: 13, offset: 11872},
							val:        "start",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 321, col: 22, offset: 11881},
							expr: &seqExpr{
								pos: position{line: 321, col: 23, offset: 11882},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 321, col: 23, offset: 11882},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 321, col: 26, offset: 11885},
										val:        "with",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 321, col: 36, offset: 11895},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 321, col: 39, offset: 11898},
							label: "start",
							expr: &ruleRefExpr{
								pos:  position{line: 321, col: 45, offset: 11904},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cache",
			pos:  position{line: 327, col: 1, offset: 12016},
			expr: &actionExpr{
				pos: position{line: 327, col: 10, offset: 12025},
				run: (*parser).callonCache1,
				expr: &seqExpr{
					pos: position{line: 327, col: 10, offset: 12025},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 327, col: 10, offset: 12025},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 327, col: 13, offset: 12028},
							val:        "cache",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 327, col: 22, offset: 12037},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 327, col: 25, offset: 12040},
							label: "cache",
							expr: &ruleRefExpr{
								pos:  position{line: 327, col: 31, offset: 12046},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cycle",
			pos:  position{line: 333, col: 1, offset: 12158},
			expr: &actionExpr{
				pos: position{line: 333, col: 10, offset: 12167},
				run: (*parser).callonCycle1,
				expr: &seqExpr{
					pos: position{line: 333, col: 10, offset: 12167},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 333, col: 10, offset: 12167},
							label: "no",
							expr: &zeroOrOneExpr{
								pos: position{line: 333, col: 13, offset: 12170},
								expr: &seqExpr{
									pos: position{line: 333, col: 14, offset: 12171},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 333, col: 14, offset: 12171},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 333, col: 17, offset: 12174},
											val:        "no",
											ignoreCase: true,
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 333, col: 25, offset: 12182},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 333, col: 28, offset: 12185},
							val:        "cycle",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OwnedBy",
			pos:  position{line: 344, col: 1, offset: 12369},
			expr: &actionExpr{
				pos: position{line: 344, col: 12, offset: 12380},
				run: (*parser).callonOwnedBy1,
				expr: &seqExpr{
					pos: position{line: 344, col: 12, offset: 12380},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 344, col: 12, offset: 12380},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 344, col: 15, offset: 12383},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 344, col: 24, offset: 12392},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 344, col: 27, offset: 12395},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 344, col: 33, offset: 12401},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 344, col: 36, offset: 12404},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 344, col: 43, offset: 12411},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 344, col: 43, offset: 12411},
										val:        "none",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 344, col: 53, offset: 12421},
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
			pos:  position{line: 362, col: 1, offset: 13880},
			expr: &actionExpr{
				pos: position{line: 362, col: 19, offset: 13898},
				run: (*parser).callonCreateTypeStmt1,
				expr: &seqExpr{
					pos: position{line: 362, col: 19, offset: 13898},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 362, col: 19, offset: 13898},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 362, col: 29, offset: 13908},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 362, col: 32, offset: 13911},
							val:        "type",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 362, col: 40, offset: 13919},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 362, col: 43, offset: 13922},
							label: "typename",
							expr: &ruleRefExpr{
								pos:  position{line: 362, col: 52, offset: 13931},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 362, col: 58, offset: 13937},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 362, col: 61, offset: 13940},
							val:        "as",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 362, col: 67, offset: 13946},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 362, col: 70, offset: 13949},
							label: "typedef",
							expr: &ruleRefExpr{
								pos:  position{line: 362, col: 78, offset: 13957},
								name: "EnumDef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 362, col: 86, offset: 13965},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 362, col: 88, offset: 13967},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 362, col: 92, offset: 13971},
							expr: &ruleRefExpr{
								pos:  position{line: 362, col: 92, offset: 13971},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "EnumDef",
			pos:  position{line: 368, col: 1, offset: 14087},
			expr: &actionExpr{
				pos: position{line: 368, col: 12, offset: 14098},
				run: (*parser).callonEnumDef1,
				expr: &seqExpr{
					pos: position{line: 368, col: 12, offset: 14098},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 368, col: 12, offset: 14098},
							val:        "ENUM",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 368, col: 19, offset: 14105},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 368, col: 21, offset: 14107},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 368, col: 25, offset: 14111},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 368, col: 27, offset: 14113},
							label: "vals",
							expr: &seqExpr{
								pos: position{line: 368, col: 34, offset: 14120},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 368, col: 34, offset: 14120},
										name: "StringConst",
									},
									&zeroOrMoreExpr{
										pos: position{line: 368, col: 46, offset: 14132},
										expr: &seqExpr{
											pos: position{line: 368, col: 48, offset: 14134},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 368, col: 48, offset: 14134},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 368, col: 50, offset: 14136},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 368, col: 54, offset: 14140},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 368, col: 56, offset: 14142},
													name: "StringConst",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 368, col: 74, offset: 14160},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 368, col: 76, offset: 14162},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AlterTableStmt",
			pos:  position{line: 393, col: 1, offset: 15792},
			expr: &actionExpr{
				pos: position{line: 393, col: 19, offset: 15810},
				run: (*parser).callonAlterTableStmt1,
				expr: &seqExpr{
					pos: position{line: 393, col: 19, offset: 15810},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 393, col: 19, offset: 15810},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 393, col: 28, offset: 15819},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 393, col: 31, offset: 15822},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 393, col: 40, offset: 15831},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 393, col: 43, offset: 15834},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 393, col: 48, offset: 15839},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 393, col: 54, offset: 15845},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 393, col: 57, offset: 15848},
							val:        "owner",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 393, col: 66, offset: 15857},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 393, col: 69, offset: 15860},
							val:        "to",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 393, col: 75, offset: 15866},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 393, col: 78, offset: 15869},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 393, col: 84, offset: 15875},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 393, col: 90, offset: 15881},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 393, col: 92, offset: 15883},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 393, col: 96, offset: 15887},
							expr: &ruleRefExpr{
								pos:  position{line: 393, col: 96, offset: 15887},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "AlterSeqStmt",
			pos:  position{line: 407, col: 1, offset: 17033},
			expr: &actionExpr{
				pos: position{line: 407, col: 17, offset: 17049},
				run: (*parser).callonAlterSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 407, col: 17, offset: 17049},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 407, col: 17, offset: 17049},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 407, col: 26, offset: 17058},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 407, col: 29, offset: 17061},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 407, col: 41, offset: 17073},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 407, col: 44, offset: 17076},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 407, col: 49, offset: 17081},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 407, col: 55, offset: 17087},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 407, col: 58, offset: 17090},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 407, col: 67, offset: 17099},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 407, col: 70, offset: 17102},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 407, col: 76, offset: 17108},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 407, col: 79, offset: 17111},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 407, col: 85, offset: 17117},
								name: "TableDotCol",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 407, col: 97, offset: 17129},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 407, col: 99, offset: 17131},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 407, col: 103, offset: 17135},
							expr: &ruleRefExpr{
								pos:  position{line: 407, col: 103, offset: 17135},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "TableDotCol",
			pos:  position{line: 411, col: 1, offset: 17214},
			expr: &actionExpr{
				pos: position{line: 411, col: 16, offset: 17229},
				run: (*parser).callonTableDotCol1,
				expr: &seqExpr{
					pos: position{line: 411, col: 16, offset: 17229},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 411, col: 16, offset: 17229},
							label: "table",
							expr: &ruleRefExpr{
								pos:  position{line: 411, col: 22, offset: 17235},
								name: "Ident",
							},
						},
						&litMatcher{
							pos:        position{line: 411, col: 28, offset: 17241},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 411, col: 32, offset: 17245},
							label: "column",
							expr: &ruleRefExpr{
								pos:  position{line: 411, col: 39, offset: 17252},
								name: "Ident",
							},
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 425, col: 1, offset: 18580},
			expr: &actionExpr{
				pos: position{line: 425, col: 25, offset: 18604},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 425, col: 25, offset: 18604},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 425, col: 25, offset: 18604},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 425, col: 36, offset: 18615},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 425, col: 39, offset: 18618},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 425, col: 45, offset: 18624},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 425, col: 48, offset: 18627},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 425, col: 61, offset: 18640},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 425, col: 63, offset: 18642},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 425, col: 73, offset: 18652},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 425, col: 79, offset: 18658},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 425, col: 81, offset: 18660},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 425, col: 87, offset: 18666},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 425, col: 89, offset: 18668},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 425, col: 97, offset: 18676},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 425, col: 109, offset: 18688},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 425, col: 111, offset: 18690},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 425, col: 115, offset: 18694},
							expr: &ruleRefExpr{
								pos:  position{line: 425, col: 115, offset: 18694},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 429, col: 1, offset: 18783},
			expr: &actionExpr{
				pos: position{line: 429, col: 24, offset: 18806},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 429, col: 24, offset: 18806},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 429, col: 24, offset: 18806},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 429, col: 34, offset: 18816},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 429, col: 37, offset: 18819},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 429, col: 50, offset: 18832},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 429, col: 53, offset: 18835},
							expr: &seqExpr{
								pos: position{line: 429, col: 55, offset: 18837},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 429, col: 55, offset: 18837},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 429, col: 61, offset: 18843},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 429, col: 64, offset: 18846},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 429, col: 71, offset: 18853},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 429, col: 74, offset: 18856},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 429, col: 84, offset: 18866},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 429, col: 90, offset: 18872},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 429, col: 100, offset: 18882},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 429, col: 106, offset: 18888},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 429, col: 109, offset: 18891},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 429, col: 117, offset: 18899},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 429, col: 120, offset: 18902},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 429, col: 130, offset: 18912},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 429, col: 133, offset: 18915},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 429, col: 140, offset: 18922},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 429, col: 146, offset: 18928},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 429, col: 148, offset: 18930},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 429, col: 152, offset: 18934},
							expr: &ruleRefExpr{
								pos:  position{line: 429, col: 152, offset: 18934},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 433, col: 1, offset: 19025},
			expr: &actionExpr{
				pos: position{line: 433, col: 12, offset: 19036},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 433, col: 12, offset: 19036},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 433, col: 12, offset: 19036},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 19, offset: 19043},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 433, col: 21, offset: 19045},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 433, col: 25, offset: 19049},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 29, offset: 19053},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 433, col: 33, offset: 19057},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 433, col: 33, offset: 19057},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 433, col: 39, offset: 19063},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 47, offset: 19071},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 433, col: 49, offset: 19073},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 433, col: 56, offset: 19080},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 77, offset: 19101},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 433, col: 79, offset: 19103},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 433, col: 83, offset: 19107},
							expr: &ruleRefExpr{
								pos:  position{line: 433, col: 83, offset: 19107},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 438, col: 1, offset: 19189},
			expr: &actionExpr{
				pos: position{line: 438, col: 8, offset: 19196},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 438, col: 8, offset: 19196},
					expr: &charClassMatcher{
						pos:        position{line: 438, col: 8, offset: 19196},
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
			pos:  position{line: 453, col: 1, offset: 20036},
			expr: &actionExpr{
				pos: position{line: 453, col: 25, offset: 20060},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 453, col: 25, offset: 20060},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 453, col: 32, offset: 20067},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 453, col: 32, offset: 20067},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 453, col: 38, offset: 20073},
								expr: &seqExpr{
									pos: position{line: 453, col: 40, offset: 20075},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 453, col: 40, offset: 20075},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 453, col: 42, offset: 20077},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 453, col: 46, offset: 20081},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 453, col: 48, offset: 20083},
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
			pos:  position{line: 465, col: 1, offset: 20373},
			expr: &choiceExpr{
				pos: position{line: 465, col: 12, offset: 20384},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 465, col: 12, offset: 20384},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 465, col: 21, offset: 20393},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 465, col: 31, offset: 20403},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 465, col: 45, offset: 20417},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 467, col: 1, offset: 20426},
			expr: &actionExpr{
				pos: position{line: 467, col: 16, offset: 20441},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 467, col: 16, offset: 20441},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 467, col: 16, offset: 20441},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 467, col: 20, offset: 20445},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 467, col: 26, offset: 20451},
								expr: &choiceExpr{
									pos: position{line: 467, col: 27, offset: 20452},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 467, col: 27, offset: 20452},
											val:        "[^'\\n]",
											chars:      []rune{'\'', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 467, col: 36, offset: 20461},
											val:        "''",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 467, col: 43, offset: 20468},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 471, col: 1, offset: 20520},
			expr: &actionExpr{
				pos: position{line: 471, col: 10, offset: 20529},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 471, col: 10, offset: 20529},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 471, col: 10, offset: 20529},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 471, col: 18, offset: 20537},
							expr: &charClassMatcher{
								pos:        position{line: 471, col: 18, offset: 20537},
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
			pos:  position{line: 475, col: 1, offset: 20590},
			expr: &actionExpr{
				pos: position{line: 475, col: 11, offset: 20600},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 475, col: 13, offset: 20602},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 475, col: 13, offset: 20602},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 475, col: 19, offset: 20608},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 475, col: 19, offset: 20608},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 475, col: 24, offset: 20613},
									expr: &charClassMatcher{
										pos:        position{line: 475, col: 24, offset: 20613},
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
			pos:  position{line: 480, col: 1, offset: 20708},
			expr: &actionExpr{
				pos: position{line: 480, col: 15, offset: 20722},
				run: (*parser).callonNonZNumber1,
				expr: &seqExpr{
					pos: position{line: 480, col: 15, offset: 20722},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 480, col: 15, offset: 20722},
							val:        "[1-9]",
							ranges:     []rune{'1', '9'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 480, col: 20, offset: 20727},
							expr: &charClassMatcher{
								pos:        position{line: 480, col: 20, offset: 20727},
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
			pos:  position{line: 485, col: 1, offset: 20820},
			expr: &actionExpr{
				pos: position{line: 485, col: 12, offset: 20831},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 485, col: 12, offset: 20831},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 485, col: 20, offset: 20839},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 485, col: 20, offset: 20839},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 485, col: 34, offset: 20853},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 489, col: 1, offset: 20895},
			expr: &actionExpr{
				pos: position{line: 489, col: 16, offset: 20910},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 489, col: 18, offset: 20912},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 489, col: 18, offset: 20912},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 489, col: 27, offset: 20921},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 489, col: 27, offset: 20921},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 489, col: 31, offset: 20925},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 489, col: 49, offset: 20943},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 489, col: 55, offset: 20949},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 493, col: 1, offset: 20995},
			expr: &choiceExpr{
				pos: position{line: 493, col: 24, offset: 21018},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 493, col: 24, offset: 21018},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 493, col: 33, offset: 21027},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 493, col: 41, offset: 21035},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 493, col: 48, offset: 21042},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 493, col: 54, offset: 21048},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 495, col: 1, offset: 21055},
			expr: &actionExpr{
				pos: position{line: 495, col: 17, offset: 21071},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 495, col: 19, offset: 21073},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 495, col: 19, offset: 21073},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 495, col: 29, offset: 21083},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 495, col: 29, offset: 21083},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 495, col: 33, offset: 21087},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 495, col: 52, offset: 21106},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 495, col: 58, offset: 21112},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 499, col: 1, offset: 21160},
			expr: &choiceExpr{
				pos: position{line: 499, col: 25, offset: 21184},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 499, col: 25, offset: 21184},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 499, col: 35, offset: 21194},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 499, col: 42, offset: 21201},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 499, col: 50, offset: 21209},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 499, col: 56, offset: 21215},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 511, col: 1, offset: 21730},
			expr: &actionExpr{
				pos: position{line: 511, col: 12, offset: 21741},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 511, col: 14, offset: 21743},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 511, col: 14, offset: 21743},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 511, col: 34, offset: 21763},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 515, col: 1, offset: 21807},
			expr: &seqExpr{
				pos: position{line: 515, col: 21, offset: 21827},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 515, col: 21, offset: 21827},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 515, col: 26, offset: 21832},
						expr: &anyMatcher{
							line: 515, col: 26, offset: 21832,
						},
					},
					&litMatcher{
						pos:        position{line: 515, col: 29, offset: 21835},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 515, col: 34, offset: 21840},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 517, col: 1, offset: 21845},
			expr: &seqExpr{
				pos: position{line: 517, col: 22, offset: 21866},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 517, col: 22, offset: 21866},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 517, col: 27, offset: 21871},
						expr: &charClassMatcher{
							pos:        position{line: 517, col: 27, offset: 21871},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 517, col: 36, offset: 21880},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 519, col: 1, offset: 21885},
			expr: &seqExpr{
				pos: position{line: 519, col: 9, offset: 21893},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 519, col: 9, offset: 21893},
						expr: &charClassMatcher{
							pos:        position{line: 519, col: 9, offset: 21893},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 519, col: 17, offset: 21901},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 519, col: 17, offset: 21901},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 519, col: 26, offset: 21910},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 519, col: 35, offset: 21919},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 519, col: 42, offset: 21926},
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
			pos:         position{line: 521, col: 1, offset: 21933},
			expr: &zeroOrMoreExpr{
				pos: position{line: 521, col: 19, offset: 21951},
				expr: &charClassMatcher{
					pos:        position{line: 521, col: 19, offset: 21951},
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
			pos:         position{line: 523, col: 1, offset: 21963},
			expr: &oneOrMoreExpr{
				pos: position{line: 523, col: 31, offset: 21993},
				expr: &charClassMatcher{
					pos:        position{line: 523, col: 31, offset: 21993},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 525, col: 1, offset: 22005},
			expr: &notExpr{
				pos: position{line: 525, col: 8, offset: 22012},
				expr: &anyMatcher{
					line: 525, col: 9, offset: 22013,
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
