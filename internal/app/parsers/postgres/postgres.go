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
										name: "AlterSeqStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 14, col: 142, offset: 338},
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
			pos:  position{line: 28, col: 1, offset: 1915},
			expr: &actionExpr{
				pos: position{line: 28, col: 20, offset: 1934},
				run: (*parser).callonCreateTableStmt1,
				expr: &seqExpr{
					pos: position{line: 28, col: 20, offset: 1934},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 28, col: 20, offset: 1934},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 30, offset: 1944},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 28, col: 33, offset: 1947},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 42, offset: 1956},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 45, offset: 1959},
							label: "tablename",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 55, offset: 1969},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 61, offset: 1975},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 28, col: 63, offset: 1977},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 67, offset: 1981},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 69, offset: 1983},
							label: "defs",
							expr: &seqExpr{
								pos: position{line: 28, col: 76, offset: 1990},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 28, col: 76, offset: 1990},
										name: "TableDef",
									},
									&zeroOrMoreExpr{
										pos: position{line: 28, col: 85, offset: 1999},
										expr: &seqExpr{
											pos: position{line: 28, col: 87, offset: 2001},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 28, col: 87, offset: 2001},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 28, col: 89, offset: 2003},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 28, col: 93, offset: 2007},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 28, col: 95, offset: 2009},
													name: "TableDef",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 109, offset: 2023},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 28, col: 111, offset: 2025},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 115, offset: 2029},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 28, col: 117, offset: 2031},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 28, col: 121, offset: 2035},
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 121, offset: 2035},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "TableDef",
			pos:  position{line: 48, col: 1, offset: 2632},
			expr: &choiceExpr{
				pos: position{line: 48, col: 13, offset: 2644},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 48, col: 13, offset: 2644},
						name: "TableConstr",
					},
					&ruleRefExpr{
						pos:  position{line: 48, col: 27, offset: 2658},
						name: "ColumnDef",
					},
				},
			},
		},
		{
			name: "ColumnDef",
			pos:  position{line: 50, col: 1, offset: 2669},
			expr: &actionExpr{
				pos: position{line: 50, col: 14, offset: 2682},
				run: (*parser).callonColumnDef1,
				expr: &seqExpr{
					pos: position{line: 50, col: 14, offset: 2682},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 50, col: 14, offset: 2682},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 50, col: 19, offset: 2687},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 50, col: 25, offset: 2693},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 50, col: 28, offset: 2696},
							label: "dataType",
							expr: &ruleRefExpr{
								pos:  position{line: 50, col: 37, offset: 2705},
								name: "DataType",
							},
						},
						&labeledExpr{
							pos:   position{line: 50, col: 46, offset: 2714},
							label: "constraint",
							expr: &zeroOrOneExpr{
								pos: position{line: 50, col: 57, offset: 2725},
								expr: &ruleRefExpr{
									pos:  position{line: 50, col: 57, offset: 2725},
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
			pos:  position{line: 67, col: 1, offset: 3214},
			expr: &actionExpr{
				pos: position{line: 67, col: 21, offset: 3234},
				run: (*parser).callonColumnConstraint1,
				expr: &seqExpr{
					pos: position{line: 67, col: 21, offset: 3234},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 67, col: 21, offset: 3234},
							label: "nameOpt",
							expr: &zeroOrOneExpr{
								pos: position{line: 67, col: 29, offset: 3242},
								expr: &seqExpr{
									pos: position{line: 67, col: 31, offset: 3244},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 67, col: 31, offset: 3244},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 67, col: 34, offset: 3247},
											val:        "constraint",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 67, col: 48, offset: 3261},
											name: "_1",
										},
										&choiceExpr{
											pos: position{line: 67, col: 52, offset: 3265},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 67, col: 52, offset: 3265},
													name: "StringConst",
												},
												&ruleRefExpr{
													pos:  position{line: 67, col: 66, offset: 3279},
													name: "Ident",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 67, col: 76, offset: 3289},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 67, col: 78, offset: 3291},
							label: "constraint",
							expr: &choiceExpr{
								pos: position{line: 67, col: 91, offset: 3304},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 67, col: 91, offset: 3304},
										name: "NotNullCls",
									},
									&ruleRefExpr{
										pos:  position{line: 67, col: 104, offset: 3317},
										name: "NullCls",
									},
									&ruleRefExpr{
										pos:  position{line: 67, col: 114, offset: 3327},
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
			pos:  position{line: 82, col: 1, offset: 3693},
			expr: &actionExpr{
				pos: position{line: 82, col: 16, offset: 3708},
				run: (*parser).callonTableConstr1,
				expr: &seqExpr{
					pos: position{line: 82, col: 16, offset: 3708},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 82, col: 16, offset: 3708},
							label: "nameOpt",
							expr: &zeroOrOneExpr{
								pos: position{line: 82, col: 24, offset: 3716},
								expr: &seqExpr{
									pos: position{line: 82, col: 26, offset: 3718},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 82, col: 26, offset: 3718},
											val:        "constraint",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 82, col: 40, offset: 3732},
											name: "_1",
										},
										&choiceExpr{
											pos: position{line: 82, col: 44, offset: 3736},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 82, col: 44, offset: 3736},
													name: "StringConst",
												},
												&ruleRefExpr{
													pos:  position{line: 82, col: 58, offset: 3750},
													name: "Ident",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 82, col: 68, offset: 3760},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 82, col: 70, offset: 3762},
							label: "constraint",
							expr: &ruleRefExpr{
								pos:  position{line: 82, col: 81, offset: 3773},
								name: "CheckCls",
							},
						},
					},
				},
			},
		},
		{
			name: "NotNullCls",
			pos:  position{line: 99, col: 1, offset: 4174},
			expr: &actionExpr{
				pos: position{line: 99, col: 15, offset: 4188},
				run: (*parser).callonNotNullCls1,
				expr: &seqExpr{
					pos: position{line: 99, col: 15, offset: 4188},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 99, col: 15, offset: 4188},
							val:        "not",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 22, offset: 4195},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 99, col: 25, offset: 4198},
							val:        "null",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "NullCls",
			pos:  position{line: 105, col: 1, offset: 4280},
			expr: &actionExpr{
				pos: position{line: 105, col: 12, offset: 4291},
				run: (*parser).callonNullCls1,
				expr: &litMatcher{
					pos:        position{line: 105, col: 12, offset: 4291},
					val:        "null",
					ignoreCase: true,
				},
			},
		},
		{
			name: "CheckCls",
			pos:  position{line: 111, col: 1, offset: 4374},
			expr: &actionExpr{
				pos: position{line: 111, col: 13, offset: 4386},
				run: (*parser).callonCheckCls1,
				expr: &seqExpr{
					pos: position{line: 111, col: 13, offset: 4386},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 111, col: 13, offset: 4386},
							val:        "check",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 111, col: 22, offset: 4395},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 111, col: 25, offset: 4398},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 111, col: 30, offset: 4403},
								name: "WrappedExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 111, col: 42, offset: 4415},
							label: "noInherit",
							expr: &zeroOrOneExpr{
								pos: position{line: 111, col: 52, offset: 4425},
								expr: &seqExpr{
									pos: position{line: 111, col: 54, offset: 4427},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 111, col: 54, offset: 4427},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 111, col: 57, offset: 4430},
											val:        "no",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 111, col: 63, offset: 4436},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 111, col: 66, offset: 4439},
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
			pos:  position{line: 121, col: 1, offset: 4632},
			expr: &actionExpr{
				pos: position{line: 121, col: 16, offset: 4647},
				run: (*parser).callonWrappedExpr1,
				expr: &seqExpr{
					pos: position{line: 121, col: 16, offset: 4647},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 121, col: 16, offset: 4647},
							val:        "(",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 121, col: 20, offset: 4651},
							expr: &ruleRefExpr{
								pos:  position{line: 121, col: 20, offset: 4651},
								name: "Expr",
							},
						},
						&litMatcher{
							pos:        position{line: 121, col: 26, offset: 4657},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 125, col: 1, offset: 4697},
			expr: &choiceExpr{
				pos: position{line: 125, col: 9, offset: 4705},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 125, col: 9, offset: 4705},
						name: "WrappedExpr",
					},
					&oneOrMoreExpr{
						pos: position{line: 125, col: 23, offset: 4719},
						expr: &charClassMatcher{
							pos:        position{line: 125, col: 23, offset: 4719},
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
			pos:  position{line: 137, col: 1, offset: 5739},
			expr: &actionExpr{
				pos: position{line: 137, col: 13, offset: 5751},
				run: (*parser).callonDataType1,
				expr: &labeledExpr{
					pos:   position{line: 137, col: 13, offset: 5751},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 137, col: 17, offset: 5755},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 137, col: 17, offset: 5755},
								name: "TimestampT",
							},
							&ruleRefExpr{
								pos:  position{line: 137, col: 30, offset: 5768},
								name: "TimeT",
							},
							&ruleRefExpr{
								pos:  position{line: 137, col: 38, offset: 5776},
								name: "VarcharT",
							},
							&ruleRefExpr{
								pos:  position{line: 137, col: 49, offset: 5787},
								name: "CharT",
							},
							&ruleRefExpr{
								pos:  position{line: 137, col: 57, offset: 5795},
								name: "BitVarT",
							},
							&ruleRefExpr{
								pos:  position{line: 137, col: 67, offset: 5805},
								name: "BitT",
							},
							&ruleRefExpr{
								pos:  position{line: 137, col: 74, offset: 5812},
								name: "IntT",
							},
							&ruleRefExpr{
								pos:  position{line: 137, col: 81, offset: 5819},
								name: "PgOidT",
							},
							&ruleRefExpr{
								pos:  position{line: 137, col: 90, offset: 5828},
								name: "OtherT",
							},
							&ruleRefExpr{
								pos:  position{line: 137, col: 99, offset: 5837},
								name: "CustomT",
							},
						},
					},
				},
			},
		},
		{
			name: "TimestampT",
			pos:  position{line: 141, col: 1, offset: 5870},
			expr: &actionExpr{
				pos: position{line: 141, col: 15, offset: 5884},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 141, col: 15, offset: 5884},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 141, col: 15, offset: 5884},
							val:        "timestamp",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 141, col: 28, offset: 5897},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 141, col: 33, offset: 5902},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 141, col: 46, offset: 5915},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 141, col: 59, offset: 5928},
								expr: &choiceExpr{
									pos: position{line: 141, col: 61, offset: 5930},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 141, col: 61, offset: 5930},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 141, col: 70, offset: 5939},
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
			pos:  position{line: 154, col: 1, offset: 6218},
			expr: &actionExpr{
				pos: position{line: 154, col: 10, offset: 6227},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 154, col: 10, offset: 6227},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 154, col: 10, offset: 6227},
							val:        "time",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 154, col: 18, offset: 6235},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 154, col: 23, offset: 6240},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 154, col: 36, offset: 6253},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 154, col: 49, offset: 6266},
								expr: &choiceExpr{
									pos: position{line: 154, col: 51, offset: 6268},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 154, col: 51, offset: 6268},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 154, col: 60, offset: 6277},
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
			pos:  position{line: 167, col: 1, offset: 6548},
			expr: &actionExpr{
				pos: position{line: 167, col: 17, offset: 6564},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 167, col: 17, offset: 6564},
					expr: &seqExpr{
						pos: position{line: 167, col: 19, offset: 6566},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 167, col: 19, offset: 6566},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 167, col: 22, offset: 6569},
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
			pos:  position{line: 174, col: 1, offset: 6697},
			expr: &actionExpr{
				pos: position{line: 174, col: 11, offset: 6707},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 174, col: 11, offset: 6707},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 174, col: 11, offset: 6707},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 174, col: 14, offset: 6710},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 174, col: 22, offset: 6718},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 174, col: 25, offset: 6721},
							val:        "time",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 174, col: 33, offset: 6729},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 174, col: 36, offset: 6732},
							val:        "zone",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 178, col: 1, offset: 6766},
			expr: &actionExpr{
				pos: position{line: 178, col: 14, offset: 6779},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 178, col: 14, offset: 6779},
					expr: &seqExpr{
						pos: position{line: 178, col: 16, offset: 6781},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 178, col: 16, offset: 6781},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 178, col: 19, offset: 6784},
								val:        "without",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 178, col: 30, offset: 6795},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 178, col: 33, offset: 6798},
								val:        "time",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 178, col: 41, offset: 6806},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 178, col: 44, offset: 6809},
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
			pos:  position{line: 182, col: 1, offset: 6847},
			expr: &actionExpr{
				pos: position{line: 182, col: 10, offset: 6856},
				run: (*parser).callonCharT1,
				expr: &seqExpr{
					pos: position{line: 182, col: 10, offset: 6856},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 182, col: 12, offset: 6858},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 182, col: 12, offset: 6858},
									val:        "character",
									ignoreCase: true,
								},
								&litMatcher{
									pos:        position{line: 182, col: 27, offset: 6873},
									val:        "char",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 182, col: 37, offset: 6883},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 182, col: 44, offset: 6890},
								expr: &seqExpr{
									pos: position{line: 182, col: 46, offset: 6892},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 182, col: 46, offset: 6892},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 182, col: 50, offset: 6896},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 182, col: 61, offset: 6907},
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
			pos:  position{line: 194, col: 1, offset: 7162},
			expr: &actionExpr{
				pos: position{line: 194, col: 13, offset: 7174},
				run: (*parser).callonVarcharT1,
				expr: &seqExpr{
					pos: position{line: 194, col: 13, offset: 7174},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 194, col: 15, offset: 7176},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 194, col: 17, offset: 7178},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 194, col: 17, offset: 7178},
											val:        "character",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 30, offset: 7191},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 194, col: 33, offset: 7194},
											val:        "varying",
											ignoreCase: true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 194, col: 48, offset: 7209},
									val:        "varchar",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 194, col: 61, offset: 7222},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 194, col: 68, offset: 7229},
								expr: &seqExpr{
									pos: position{line: 194, col: 70, offset: 7231},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 194, col: 70, offset: 7231},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 74, offset: 7235},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 194, col: 85, offset: 7246},
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
			pos:  position{line: 205, col: 1, offset: 7481},
			expr: &actionExpr{
				pos: position{line: 205, col: 9, offset: 7489},
				run: (*parser).callonBitT1,
				expr: &seqExpr{
					pos: position{line: 205, col: 9, offset: 7489},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 205, col: 9, offset: 7489},
							val:        "bit",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 205, col: 16, offset: 7496},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 205, col: 23, offset: 7503},
								expr: &seqExpr{
									pos: position{line: 205, col: 25, offset: 7505},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 205, col: 25, offset: 7505},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 205, col: 29, offset: 7509},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 205, col: 40, offset: 7520},
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
			pos:  position{line: 217, col: 1, offset: 7774},
			expr: &actionExpr{
				pos: position{line: 217, col: 12, offset: 7785},
				run: (*parser).callonBitVarT1,
				expr: &seqExpr{
					pos: position{line: 217, col: 12, offset: 7785},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 217, col: 12, offset: 7785},
							val:        "bit",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 217, col: 19, offset: 7792},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 217, col: 22, offset: 7795},
							val:        "varying",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 217, col: 33, offset: 7806},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 217, col: 40, offset: 7813},
								expr: &seqExpr{
									pos: position{line: 217, col: 42, offset: 7815},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 217, col: 42, offset: 7815},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 217, col: 46, offset: 7819},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 217, col: 57, offset: 7830},
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
			pos:  position{line: 228, col: 1, offset: 8064},
			expr: &actionExpr{
				pos: position{line: 228, col: 9, offset: 8072},
				run: (*parser).callonIntT1,
				expr: &choiceExpr{
					pos: position{line: 228, col: 11, offset: 8074},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 228, col: 11, offset: 8074},
							val:        "integer",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 228, col: 24, offset: 8087},
							val:        "int",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "PgOidT",
			pos:  position{line: 234, col: 1, offset: 8169},
			expr: &actionExpr{
				pos: position{line: 234, col: 11, offset: 8179},
				run: (*parser).callonPgOidT1,
				expr: &choiceExpr{
					pos: position{line: 234, col: 13, offset: 8181},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 234, col: 13, offset: 8181},
							val:        "oid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 22, offset: 8190},
							val:        "regprocedure",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 40, offset: 8208},
							val:        "regproc",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 53, offset: 8221},
							val:        "regoperator",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 70, offset: 8238},
							val:        "regoper",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 83, offset: 8251},
							val:        "regclass",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 97, offset: 8265},
							val:        "regtype",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 110, offset: 8278},
							val:        "regrole",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 123, offset: 8291},
							val:        "regnamespace",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 141, offset: 8309},
							val:        "regconfig",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 234, col: 156, offset: 8324},
							val:        "regdictionary",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 240, col: 1, offset: 8438},
			expr: &actionExpr{
				pos: position{line: 240, col: 11, offset: 8448},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 240, col: 13, offset: 8450},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 240, col: 13, offset: 8450},
							val:        "date",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 23, offset: 8460},
							val:        "smallint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 37, offset: 8474},
							val:        "bigint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 49, offset: 8486},
							val:        "decimal",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 62, offset: 8499},
							val:        "numeric",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 75, offset: 8512},
							val:        "real",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 85, offset: 8522},
							val:        "smallserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 102, offset: 8539},
							val:        "serial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 114, offset: 8551},
							val:        "bigserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 129, offset: 8566},
							val:        "boolean",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 142, offset: 8579},
							val:        "text",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 152, offset: 8589},
							val:        "money",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 163, offset: 8600},
							val:        "bytea",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 174, offset: 8611},
							val:        "point",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 185, offset: 8622},
							val:        "line",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 195, offset: 8632},
							val:        "lseg",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 205, offset: 8642},
							val:        "box",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 214, offset: 8651},
							val:        "path",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 224, offset: 8661},
							val:        "polygon",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 237, offset: 8674},
							val:        "circle",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 249, offset: 8686},
							val:        "cidr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 259, offset: 8696},
							val:        "inet",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 269, offset: 8706},
							val:        "macaddr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 282, offset: 8719},
							val:        "uuid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 292, offset: 8729},
							val:        "xml",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 301, offset: 8738},
							val:        "jsonb",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 312, offset: 8749},
							val:        "json",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CustomT",
			pos:  position{line: 246, col: 1, offset: 8854},
			expr: &actionExpr{
				pos: position{line: 246, col: 13, offset: 8866},
				run: (*parser).callonCustomT1,
				expr: &ruleRefExpr{
					pos:  position{line: 246, col: 13, offset: 8866},
					name: "Ident",
				},
			},
		},
		{
			name: "CreateSeqStmt",
			pos:  position{line: 267, col: 1, offset: 10331},
			expr: &actionExpr{
				pos: position{line: 267, col: 18, offset: 10348},
				run: (*parser).callonCreateSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 267, col: 18, offset: 10348},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 267, col: 18, offset: 10348},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 267, col: 28, offset: 10358},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 267, col: 31, offset: 10361},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 267, col: 43, offset: 10373},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 267, col: 46, offset: 10376},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 267, col: 51, offset: 10381},
								name: "Ident",
							},
						},
						&labeledExpr{
							pos:   position{line: 267, col: 57, offset: 10387},
							label: "verses",
							expr: &zeroOrMoreExpr{
								pos: position{line: 267, col: 64, offset: 10394},
								expr: &ruleRefExpr{
									pos:  position{line: 267, col: 64, offset: 10394},
									name: "CreateSeqVerse",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 267, col: 80, offset: 10410},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 267, col: 82, offset: 10412},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 267, col: 86, offset: 10416},
							expr: &ruleRefExpr{
								pos:  position{line: 267, col: 86, offset: 10416},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateSeqVerse",
			pos:  position{line: 281, col: 1, offset: 10810},
			expr: &actionExpr{
				pos: position{line: 281, col: 19, offset: 10828},
				run: (*parser).callonCreateSeqVerse1,
				expr: &labeledExpr{
					pos:   position{line: 281, col: 19, offset: 10828},
					label: "verse",
					expr: &choiceExpr{
						pos: position{line: 281, col: 27, offset: 10836},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 281, col: 27, offset: 10836},
								name: "IncrementBy",
							},
							&ruleRefExpr{
								pos:  position{line: 281, col: 41, offset: 10850},
								name: "MinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 281, col: 52, offset: 10861},
								name: "NoMinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 281, col: 65, offset: 10874},
								name: "MaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 281, col: 76, offset: 10885},
								name: "NoMaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 281, col: 89, offset: 10898},
								name: "Start",
							},
							&ruleRefExpr{
								pos:  position{line: 281, col: 97, offset: 10906},
								name: "Cache",
							},
							&ruleRefExpr{
								pos:  position{line: 281, col: 105, offset: 10914},
								name: "Cycle",
							},
							&ruleRefExpr{
								pos:  position{line: 281, col: 113, offset: 10922},
								name: "OwnedBy",
							},
						},
					},
				},
			},
		},
		{
			name: "IncrementBy",
			pos:  position{line: 285, col: 1, offset: 10959},
			expr: &actionExpr{
				pos: position{line: 285, col: 16, offset: 10974},
				run: (*parser).callonIncrementBy1,
				expr: &seqExpr{
					pos: position{line: 285, col: 16, offset: 10974},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 285, col: 16, offset: 10974},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 285, col: 19, offset: 10977},
							val:        "increment",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 285, col: 32, offset: 10990},
							expr: &seqExpr{
								pos: position{line: 285, col: 33, offset: 10991},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 285, col: 33, offset: 10991},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 285, col: 36, offset: 10994},
										val:        "by",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 285, col: 44, offset: 11002},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 285, col: 47, offset: 11005},
							label: "num",
							expr: &ruleRefExpr{
								pos:  position{line: 285, col: 51, offset: 11009},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "MinValue",
			pos:  position{line: 291, col: 1, offset: 11123},
			expr: &actionExpr{
				pos: position{line: 291, col: 13, offset: 11135},
				run: (*parser).callonMinValue1,
				expr: &seqExpr{
					pos: position{line: 291, col: 13, offset: 11135},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 291, col: 13, offset: 11135},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 291, col: 16, offset: 11138},
							val:        "minvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 291, col: 28, offset: 11150},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 291, col: 31, offset: 11153},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 291, col: 35, offset: 11157},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMinValue",
			pos:  position{line: 297, col: 1, offset: 11270},
			expr: &actionExpr{
				pos: position{line: 297, col: 15, offset: 11284},
				run: (*parser).callonNoMinValue1,
				expr: &seqExpr{
					pos: position{line: 297, col: 15, offset: 11284},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 297, col: 15, offset: 11284},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 297, col: 18, offset: 11287},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 297, col: 24, offset: 11293},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 297, col: 27, offset: 11296},
							val:        "minvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "MaxValue",
			pos:  position{line: 301, col: 1, offset: 11333},
			expr: &actionExpr{
				pos: position{line: 301, col: 13, offset: 11345},
				run: (*parser).callonMaxValue1,
				expr: &seqExpr{
					pos: position{line: 301, col: 13, offset: 11345},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 301, col: 13, offset: 11345},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 301, col: 16, offset: 11348},
							val:        "maxvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 301, col: 28, offset: 11360},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 301, col: 31, offset: 11363},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 301, col: 35, offset: 11367},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMaxValue",
			pos:  position{line: 307, col: 1, offset: 11480},
			expr: &actionExpr{
				pos: position{line: 307, col: 15, offset: 11494},
				run: (*parser).callonNoMaxValue1,
				expr: &seqExpr{
					pos: position{line: 307, col: 15, offset: 11494},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 307, col: 15, offset: 11494},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 307, col: 18, offset: 11497},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 307, col: 24, offset: 11503},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 307, col: 27, offset: 11506},
							val:        "maxvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "Start",
			pos:  position{line: 311, col: 1, offset: 11543},
			expr: &actionExpr{
				pos: position{line: 311, col: 10, offset: 11552},
				run: (*parser).callonStart1,
				expr: &seqExpr{
					pos: position{line: 311, col: 10, offset: 11552},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 311, col: 10, offset: 11552},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 311, col: 13, offset: 11555},
							val:        "start",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 311, col: 22, offset: 11564},
							expr: &seqExpr{
								pos: position{line: 311, col: 23, offset: 11565},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 311, col: 23, offset: 11565},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 311, col: 26, offset: 11568},
										val:        "with",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 311, col: 36, offset: 11578},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 311, col: 39, offset: 11581},
							label: "start",
							expr: &ruleRefExpr{
								pos:  position{line: 311, col: 45, offset: 11587},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cache",
			pos:  position{line: 317, col: 1, offset: 11699},
			expr: &actionExpr{
				pos: position{line: 317, col: 10, offset: 11708},
				run: (*parser).callonCache1,
				expr: &seqExpr{
					pos: position{line: 317, col: 10, offset: 11708},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 317, col: 10, offset: 11708},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 317, col: 13, offset: 11711},
							val:        "cache",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 22, offset: 11720},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 317, col: 25, offset: 11723},
							label: "cache",
							expr: &ruleRefExpr{
								pos:  position{line: 317, col: 31, offset: 11729},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cycle",
			pos:  position{line: 323, col: 1, offset: 11841},
			expr: &actionExpr{
				pos: position{line: 323, col: 10, offset: 11850},
				run: (*parser).callonCycle1,
				expr: &seqExpr{
					pos: position{line: 323, col: 10, offset: 11850},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 323, col: 10, offset: 11850},
							label: "no",
							expr: &zeroOrOneExpr{
								pos: position{line: 323, col: 13, offset: 11853},
								expr: &seqExpr{
									pos: position{line: 323, col: 14, offset: 11854},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 323, col: 14, offset: 11854},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 323, col: 17, offset: 11857},
											val:        "no",
											ignoreCase: true,
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 323, col: 25, offset: 11865},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 323, col: 28, offset: 11868},
							val:        "cycle",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OwnedBy",
			pos:  position{line: 334, col: 1, offset: 12052},
			expr: &actionExpr{
				pos: position{line: 334, col: 12, offset: 12063},
				run: (*parser).callonOwnedBy1,
				expr: &seqExpr{
					pos: position{line: 334, col: 12, offset: 12063},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 334, col: 12, offset: 12063},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 334, col: 15, offset: 12066},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 334, col: 24, offset: 12075},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 334, col: 27, offset: 12078},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 334, col: 33, offset: 12084},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 334, col: 36, offset: 12087},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 334, col: 43, offset: 12094},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 334, col: 43, offset: 12094},
										val:        "none",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 334, col: 53, offset: 12104},
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
			pos:  position{line: 352, col: 1, offset: 13563},
			expr: &actionExpr{
				pos: position{line: 352, col: 19, offset: 13581},
				run: (*parser).callonCreateTypeStmt1,
				expr: &seqExpr{
					pos: position{line: 352, col: 19, offset: 13581},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 352, col: 19, offset: 13581},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 352, col: 29, offset: 13591},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 352, col: 32, offset: 13594},
							val:        "type",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 352, col: 40, offset: 13602},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 352, col: 43, offset: 13605},
							label: "typename",
							expr: &ruleRefExpr{
								pos:  position{line: 352, col: 52, offset: 13614},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 352, col: 58, offset: 13620},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 352, col: 61, offset: 13623},
							val:        "as",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 352, col: 67, offset: 13629},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 352, col: 70, offset: 13632},
							label: "typedef",
							expr: &ruleRefExpr{
								pos:  position{line: 352, col: 78, offset: 13640},
								name: "EnumDef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 352, col: 86, offset: 13648},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 352, col: 88, offset: 13650},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 352, col: 92, offset: 13654},
							expr: &ruleRefExpr{
								pos:  position{line: 352, col: 92, offset: 13654},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "EnumDef",
			pos:  position{line: 358, col: 1, offset: 13770},
			expr: &actionExpr{
				pos: position{line: 358, col: 12, offset: 13781},
				run: (*parser).callonEnumDef1,
				expr: &seqExpr{
					pos: position{line: 358, col: 12, offset: 13781},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 358, col: 12, offset: 13781},
							val:        "ENUM",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 358, col: 19, offset: 13788},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 358, col: 21, offset: 13790},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 358, col: 25, offset: 13794},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 358, col: 27, offset: 13796},
							label: "vals",
							expr: &seqExpr{
								pos: position{line: 358, col: 34, offset: 13803},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 358, col: 34, offset: 13803},
										name: "StringConst",
									},
									&zeroOrMoreExpr{
										pos: position{line: 358, col: 46, offset: 13815},
										expr: &seqExpr{
											pos: position{line: 358, col: 48, offset: 13817},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 358, col: 48, offset: 13817},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 358, col: 50, offset: 13819},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 358, col: 54, offset: 13823},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 358, col: 56, offset: 13825},
													name: "StringConst",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 358, col: 74, offset: 13843},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 358, col: 76, offset: 13845},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AlterTableStmt",
			pos:  position{line: 383, col: 1, offset: 15475},
			expr: &actionExpr{
				pos: position{line: 383, col: 19, offset: 15493},
				run: (*parser).callonAlterTableStmt1,
				expr: &seqExpr{
					pos: position{line: 383, col: 19, offset: 15493},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 383, col: 19, offset: 15493},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 383, col: 28, offset: 15502},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 383, col: 31, offset: 15505},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 383, col: 40, offset: 15514},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 383, col: 43, offset: 15517},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 383, col: 48, offset: 15522},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 383, col: 54, offset: 15528},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 383, col: 57, offset: 15531},
							val:        "owner",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 383, col: 66, offset: 15540},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 383, col: 69, offset: 15543},
							val:        "to",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 383, col: 75, offset: 15549},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 383, col: 78, offset: 15552},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 383, col: 84, offset: 15558},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 383, col: 90, offset: 15564},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 383, col: 92, offset: 15566},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 383, col: 96, offset: 15570},
							expr: &ruleRefExpr{
								pos:  position{line: 383, col: 96, offset: 15570},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "AlterSeqStmt",
			pos:  position{line: 397, col: 1, offset: 16716},
			expr: &actionExpr{
				pos: position{line: 397, col: 17, offset: 16732},
				run: (*parser).callonAlterSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 397, col: 17, offset: 16732},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 397, col: 17, offset: 16732},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 397, col: 26, offset: 16741},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 397, col: 29, offset: 16744},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 397, col: 41, offset: 16756},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 397, col: 44, offset: 16759},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 397, col: 49, offset: 16764},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 397, col: 55, offset: 16770},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 397, col: 58, offset: 16773},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 397, col: 67, offset: 16782},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 397, col: 70, offset: 16785},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 397, col: 76, offset: 16791},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 397, col: 79, offset: 16794},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 397, col: 85, offset: 16800},
								name: "TableDotCol",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 397, col: 97, offset: 16812},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 397, col: 99, offset: 16814},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 397, col: 103, offset: 16818},
							expr: &ruleRefExpr{
								pos:  position{line: 397, col: 103, offset: 16818},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "TableDotCol",
			pos:  position{line: 401, col: 1, offset: 16897},
			expr: &actionExpr{
				pos: position{line: 401, col: 16, offset: 16912},
				run: (*parser).callonTableDotCol1,
				expr: &seqExpr{
					pos: position{line: 401, col: 16, offset: 16912},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 401, col: 16, offset: 16912},
							label: "table",
							expr: &ruleRefExpr{
								pos:  position{line: 401, col: 22, offset: 16918},
								name: "Ident",
							},
						},
						&litMatcher{
							pos:        position{line: 401, col: 28, offset: 16924},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 401, col: 32, offset: 16928},
							label: "column",
							expr: &ruleRefExpr{
								pos:  position{line: 401, col: 39, offset: 16935},
								name: "Ident",
							},
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 415, col: 1, offset: 18263},
			expr: &actionExpr{
				pos: position{line: 415, col: 25, offset: 18287},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 415, col: 25, offset: 18287},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 415, col: 25, offset: 18287},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 415, col: 36, offset: 18298},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 415, col: 39, offset: 18301},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 415, col: 45, offset: 18307},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 415, col: 48, offset: 18310},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 415, col: 61, offset: 18323},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 415, col: 63, offset: 18325},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 415, col: 73, offset: 18335},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 415, col: 79, offset: 18341},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 415, col: 81, offset: 18343},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 415, col: 87, offset: 18349},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 415, col: 89, offset: 18351},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 415, col: 97, offset: 18359},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 415, col: 109, offset: 18371},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 415, col: 111, offset: 18373},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 415, col: 115, offset: 18377},
							expr: &ruleRefExpr{
								pos:  position{line: 415, col: 115, offset: 18377},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 419, col: 1, offset: 18466},
			expr: &actionExpr{
				pos: position{line: 419, col: 24, offset: 18489},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 419, col: 24, offset: 18489},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 419, col: 24, offset: 18489},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 419, col: 34, offset: 18499},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 419, col: 37, offset: 18502},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 419, col: 50, offset: 18515},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 419, col: 53, offset: 18518},
							expr: &seqExpr{
								pos: position{line: 419, col: 55, offset: 18520},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 419, col: 55, offset: 18520},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 419, col: 61, offset: 18526},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 419, col: 64, offset: 18529},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 419, col: 71, offset: 18536},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 419, col: 74, offset: 18539},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 419, col: 84, offset: 18549},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 419, col: 90, offset: 18555},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 419, col: 100, offset: 18565},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 419, col: 106, offset: 18571},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 419, col: 109, offset: 18574},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 419, col: 117, offset: 18582},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 419, col: 120, offset: 18585},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 419, col: 130, offset: 18595},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 419, col: 133, offset: 18598},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 419, col: 140, offset: 18605},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 419, col: 146, offset: 18611},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 419, col: 148, offset: 18613},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 419, col: 152, offset: 18617},
							expr: &ruleRefExpr{
								pos:  position{line: 419, col: 152, offset: 18617},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 423, col: 1, offset: 18708},
			expr: &actionExpr{
				pos: position{line: 423, col: 12, offset: 18719},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 423, col: 12, offset: 18719},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 423, col: 12, offset: 18719},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 423, col: 19, offset: 18726},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 423, col: 21, offset: 18728},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 423, col: 25, offset: 18732},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 423, col: 29, offset: 18736},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 423, col: 33, offset: 18740},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 423, col: 33, offset: 18740},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 423, col: 39, offset: 18746},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 423, col: 47, offset: 18754},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 423, col: 49, offset: 18756},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 423, col: 56, offset: 18763},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 423, col: 77, offset: 18784},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 423, col: 79, offset: 18786},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 423, col: 83, offset: 18790},
							expr: &ruleRefExpr{
								pos:  position{line: 423, col: 83, offset: 18790},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 428, col: 1, offset: 18872},
			expr: &actionExpr{
				pos: position{line: 428, col: 8, offset: 18879},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 428, col: 8, offset: 18879},
					expr: &charClassMatcher{
						pos:        position{line: 428, col: 8, offset: 18879},
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
			pos:  position{line: 443, col: 1, offset: 19719},
			expr: &actionExpr{
				pos: position{line: 443, col: 25, offset: 19743},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 443, col: 25, offset: 19743},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 443, col: 32, offset: 19750},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 443, col: 32, offset: 19750},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 443, col: 38, offset: 19756},
								expr: &seqExpr{
									pos: position{line: 443, col: 40, offset: 19758},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 443, col: 40, offset: 19758},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 443, col: 42, offset: 19760},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 443, col: 46, offset: 19764},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 443, col: 48, offset: 19766},
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
			pos:  position{line: 455, col: 1, offset: 20056},
			expr: &choiceExpr{
				pos: position{line: 455, col: 12, offset: 20067},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 455, col: 12, offset: 20067},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 455, col: 21, offset: 20076},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 455, col: 31, offset: 20086},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 455, col: 45, offset: 20100},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 457, col: 1, offset: 20109},
			expr: &actionExpr{
				pos: position{line: 457, col: 16, offset: 20124},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 457, col: 16, offset: 20124},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 457, col: 16, offset: 20124},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 457, col: 20, offset: 20128},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 457, col: 26, offset: 20134},
								expr: &choiceExpr{
									pos: position{line: 457, col: 27, offset: 20135},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 457, col: 27, offset: 20135},
											val:        "[^'\\n]",
											chars:      []rune{'\'', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 457, col: 36, offset: 20144},
											val:        "''",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 457, col: 43, offset: 20151},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 461, col: 1, offset: 20203},
			expr: &actionExpr{
				pos: position{line: 461, col: 10, offset: 20212},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 461, col: 10, offset: 20212},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 461, col: 10, offset: 20212},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 461, col: 18, offset: 20220},
							expr: &charClassMatcher{
								pos:        position{line: 461, col: 18, offset: 20220},
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
			pos:  position{line: 465, col: 1, offset: 20273},
			expr: &actionExpr{
				pos: position{line: 465, col: 11, offset: 20283},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 465, col: 13, offset: 20285},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 465, col: 13, offset: 20285},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 465, col: 19, offset: 20291},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 465, col: 19, offset: 20291},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 465, col: 24, offset: 20296},
									expr: &charClassMatcher{
										pos:        position{line: 465, col: 24, offset: 20296},
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
			pos:  position{line: 470, col: 1, offset: 20391},
			expr: &actionExpr{
				pos: position{line: 470, col: 15, offset: 20405},
				run: (*parser).callonNonZNumber1,
				expr: &seqExpr{
					pos: position{line: 470, col: 15, offset: 20405},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 470, col: 15, offset: 20405},
							val:        "[1-9]",
							ranges:     []rune{'1', '9'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 470, col: 20, offset: 20410},
							expr: &charClassMatcher{
								pos:        position{line: 470, col: 20, offset: 20410},
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
			pos:  position{line: 475, col: 1, offset: 20503},
			expr: &actionExpr{
				pos: position{line: 475, col: 12, offset: 20514},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 475, col: 12, offset: 20514},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 475, col: 20, offset: 20522},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 475, col: 20, offset: 20522},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 475, col: 34, offset: 20536},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 479, col: 1, offset: 20578},
			expr: &actionExpr{
				pos: position{line: 479, col: 16, offset: 20593},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 479, col: 18, offset: 20595},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 479, col: 18, offset: 20595},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 479, col: 27, offset: 20604},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 479, col: 27, offset: 20604},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 479, col: 31, offset: 20608},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 479, col: 49, offset: 20626},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 479, col: 55, offset: 20632},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 483, col: 1, offset: 20678},
			expr: &choiceExpr{
				pos: position{line: 483, col: 24, offset: 20701},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 483, col: 24, offset: 20701},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 483, col: 33, offset: 20710},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 483, col: 41, offset: 20718},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 483, col: 48, offset: 20725},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 483, col: 54, offset: 20731},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 485, col: 1, offset: 20738},
			expr: &actionExpr{
				pos: position{line: 485, col: 17, offset: 20754},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 485, col: 19, offset: 20756},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 485, col: 19, offset: 20756},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 485, col: 29, offset: 20766},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 485, col: 29, offset: 20766},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 485, col: 33, offset: 20770},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 485, col: 52, offset: 20789},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 485, col: 58, offset: 20795},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 489, col: 1, offset: 20843},
			expr: &choiceExpr{
				pos: position{line: 489, col: 25, offset: 20867},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 489, col: 25, offset: 20867},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 489, col: 35, offset: 20877},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 489, col: 42, offset: 20884},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 489, col: 50, offset: 20892},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 489, col: 56, offset: 20898},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 501, col: 1, offset: 21413},
			expr: &actionExpr{
				pos: position{line: 501, col: 12, offset: 21424},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 501, col: 14, offset: 21426},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 501, col: 14, offset: 21426},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 501, col: 34, offset: 21446},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 505, col: 1, offset: 21490},
			expr: &seqExpr{
				pos: position{line: 505, col: 21, offset: 21510},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 505, col: 21, offset: 21510},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 505, col: 26, offset: 21515},
						expr: &anyMatcher{
							line: 505, col: 26, offset: 21515,
						},
					},
					&litMatcher{
						pos:        position{line: 505, col: 29, offset: 21518},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 505, col: 34, offset: 21523},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 507, col: 1, offset: 21528},
			expr: &seqExpr{
				pos: position{line: 507, col: 22, offset: 21549},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 507, col: 22, offset: 21549},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 507, col: 27, offset: 21554},
						expr: &charClassMatcher{
							pos:        position{line: 507, col: 27, offset: 21554},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 507, col: 36, offset: 21563},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 509, col: 1, offset: 21568},
			expr: &seqExpr{
				pos: position{line: 509, col: 9, offset: 21576},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 509, col: 9, offset: 21576},
						expr: &charClassMatcher{
							pos:        position{line: 509, col: 9, offset: 21576},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 509, col: 17, offset: 21584},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 509, col: 17, offset: 21584},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 509, col: 26, offset: 21593},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 509, col: 35, offset: 21602},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 509, col: 42, offset: 21609},
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
			pos:         position{line: 511, col: 1, offset: 21616},
			expr: &zeroOrMoreExpr{
				pos: position{line: 511, col: 19, offset: 21634},
				expr: &charClassMatcher{
					pos:        position{line: 511, col: 19, offset: 21634},
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
			pos:         position{line: 513, col: 1, offset: 21646},
			expr: &oneOrMoreExpr{
				pos: position{line: 513, col: 31, offset: 21676},
				expr: &charClassMatcher{
					pos:        position{line: 513, col: 31, offset: 21676},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 515, col: 1, offset: 21688},
			expr: &notExpr{
				pos: position{line: 515, col: 8, offset: 21695},
				expr: &anyMatcher{
					line: 515, col: 9, offset: 21696,
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
