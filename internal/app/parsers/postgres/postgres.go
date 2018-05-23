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
										name: "VarcharT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 49, offset: 5841},
										name: "CharT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 57, offset: 5849},
										name: "BitVarT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 67, offset: 5859},
										name: "BitT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 74, offset: 5866},
										name: "IntT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 81, offset: 5873},
										name: "PgOidT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 90, offset: 5882},
										name: "GeographyT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 103, offset: 5895},
										name: "OtherT",
									},
									&ruleRefExpr{
										pos:  position{line: 140, col: 112, offset: 5904},
										name: "CustomT",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 140, col: 122, offset: 5914},
							label: "brackets",
							expr: &zeroOrMoreExpr{
								pos: position{line: 140, col: 131, offset: 5923},
								expr: &litMatcher{
									pos:        position{line: 140, col: 133, offset: 5925},
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
			pos:  position{line: 153, col: 1, offset: 6236},
			expr: &actionExpr{
				pos: position{line: 153, col: 15, offset: 6250},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 153, col: 15, offset: 6250},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 153, col: 15, offset: 6250},
							val:        "timestamp",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 153, col: 28, offset: 6263},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 153, col: 33, offset: 6268},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 153, col: 46, offset: 6281},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 153, col: 59, offset: 6294},
								expr: &choiceExpr{
									pos: position{line: 153, col: 61, offset: 6296},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 153, col: 61, offset: 6296},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 153, col: 70, offset: 6305},
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
			pos:  position{line: 166, col: 1, offset: 6584},
			expr: &actionExpr{
				pos: position{line: 166, col: 10, offset: 6593},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 166, col: 10, offset: 6593},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 166, col: 10, offset: 6593},
							val:        "time",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 166, col: 18, offset: 6601},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 166, col: 23, offset: 6606},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 166, col: 36, offset: 6619},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 166, col: 49, offset: 6632},
								expr: &choiceExpr{
									pos: position{line: 166, col: 51, offset: 6634},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 166, col: 51, offset: 6634},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 166, col: 60, offset: 6643},
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
			pos:  position{line: 179, col: 1, offset: 6914},
			expr: &actionExpr{
				pos: position{line: 179, col: 17, offset: 6930},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 179, col: 17, offset: 6930},
					expr: &seqExpr{
						pos: position{line: 179, col: 19, offset: 6932},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 179, col: 19, offset: 6932},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 179, col: 22, offset: 6935},
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
			pos:  position{line: 186, col: 1, offset: 7063},
			expr: &actionExpr{
				pos: position{line: 186, col: 11, offset: 7073},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 186, col: 11, offset: 7073},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 186, col: 11, offset: 7073},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 186, col: 14, offset: 7076},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 22, offset: 7084},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 186, col: 25, offset: 7087},
							val:        "time",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 33, offset: 7095},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 186, col: 36, offset: 7098},
							val:        "zone",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 190, col: 1, offset: 7132},
			expr: &actionExpr{
				pos: position{line: 190, col: 14, offset: 7145},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 190, col: 14, offset: 7145},
					expr: &seqExpr{
						pos: position{line: 190, col: 16, offset: 7147},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 190, col: 16, offset: 7147},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 190, col: 19, offset: 7150},
								val:        "without",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 190, col: 30, offset: 7161},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 190, col: 33, offset: 7164},
								val:        "time",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 190, col: 41, offset: 7172},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 190, col: 44, offset: 7175},
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
			pos:  position{line: 194, col: 1, offset: 7213},
			expr: &actionExpr{
				pos: position{line: 194, col: 10, offset: 7222},
				run: (*parser).callonCharT1,
				expr: &seqExpr{
					pos: position{line: 194, col: 10, offset: 7222},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 194, col: 12, offset: 7224},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 194, col: 12, offset: 7224},
									val:        "character",
									ignoreCase: true,
								},
								&litMatcher{
									pos:        position{line: 194, col: 27, offset: 7239},
									val:        "char",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 194, col: 37, offset: 7249},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 194, col: 44, offset: 7256},
								expr: &seqExpr{
									pos: position{line: 194, col: 46, offset: 7258},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 194, col: 46, offset: 7258},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 50, offset: 7262},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 194, col: 61, offset: 7273},
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
			pos:  position{line: 206, col: 1, offset: 7528},
			expr: &actionExpr{
				pos: position{line: 206, col: 13, offset: 7540},
				run: (*parser).callonVarcharT1,
				expr: &seqExpr{
					pos: position{line: 206, col: 13, offset: 7540},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 206, col: 15, offset: 7542},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 206, col: 17, offset: 7544},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 206, col: 17, offset: 7544},
											val:        "character",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 206, col: 30, offset: 7557},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 206, col: 33, offset: 7560},
											val:        "varying",
											ignoreCase: true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 206, col: 48, offset: 7575},
									val:        "varchar",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 206, col: 61, offset: 7588},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 206, col: 68, offset: 7595},
								expr: &seqExpr{
									pos: position{line: 206, col: 70, offset: 7597},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 206, col: 70, offset: 7597},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 206, col: 74, offset: 7601},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 206, col: 85, offset: 7612},
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
			pos:  position{line: 217, col: 1, offset: 7847},
			expr: &actionExpr{
				pos: position{line: 217, col: 9, offset: 7855},
				run: (*parser).callonBitT1,
				expr: &seqExpr{
					pos: position{line: 217, col: 9, offset: 7855},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 217, col: 9, offset: 7855},
							val:        "bit",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 217, col: 16, offset: 7862},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 217, col: 23, offset: 7869},
								expr: &seqExpr{
									pos: position{line: 217, col: 25, offset: 7871},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 217, col: 25, offset: 7871},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 217, col: 29, offset: 7875},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 217, col: 40, offset: 7886},
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
			pos:  position{line: 229, col: 1, offset: 8140},
			expr: &actionExpr{
				pos: position{line: 229, col: 12, offset: 8151},
				run: (*parser).callonBitVarT1,
				expr: &seqExpr{
					pos: position{line: 229, col: 12, offset: 8151},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 229, col: 12, offset: 8151},
							val:        "bit",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 229, col: 19, offset: 8158},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 229, col: 22, offset: 8161},
							val:        "varying",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 229, col: 33, offset: 8172},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 229, col: 40, offset: 8179},
								expr: &seqExpr{
									pos: position{line: 229, col: 42, offset: 8181},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 229, col: 42, offset: 8181},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 229, col: 46, offset: 8185},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 229, col: 57, offset: 8196},
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
			pos:  position{line: 240, col: 1, offset: 8430},
			expr: &actionExpr{
				pos: position{line: 240, col: 9, offset: 8438},
				run: (*parser).callonIntT1,
				expr: &choiceExpr{
					pos: position{line: 240, col: 11, offset: 8440},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 240, col: 11, offset: 8440},
							val:        "integer",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 240, col: 24, offset: 8453},
							val:        "int",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "GeographyT",
			pos:  position{line: 246, col: 1, offset: 8535},
			expr: &actionExpr{
				pos: position{line: 246, col: 15, offset: 8549},
				run: (*parser).callonGeographyT1,
				expr: &seqExpr{
					pos: position{line: 246, col: 15, offset: 8549},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 246, col: 15, offset: 8549},
							val:        "geography",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 246, col: 28, offset: 8562},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 246, col: 32, offset: 8566},
							label: "subtype",
							expr: &choiceExpr{
								pos: position{line: 246, col: 42, offset: 8576},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 246, col: 42, offset: 8576},
										val:        "point",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 246, col: 53, offset: 8587},
										val:        "linestring",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 246, col: 69, offset: 8603},
										val:        "polygon",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 246, col: 82, offset: 8616},
										val:        "multipoint",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 246, col: 98, offset: 8632},
										val:        "multilinestring",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 246, col: 119, offset: 8653},
										val:        "multipolygon",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 246, col: 137, offset: 8671},
										val:        "geometrycollection",
										ignoreCase: true,
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 246, col: 161, offset: 8695},
							label: "srid",
							expr: &zeroOrOneExpr{
								pos: position{line: 246, col: 166, offset: 8700},
								expr: &seqExpr{
									pos: position{line: 246, col: 167, offset: 8701},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 246, col: 167, offset: 8701},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 246, col: 171, offset: 8705},
											name: "NonZNumber",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 246, col: 184, offset: 8718},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PgOidT",
			pos:  position{line: 259, col: 1, offset: 9062},
			expr: &actionExpr{
				pos: position{line: 259, col: 11, offset: 9072},
				run: (*parser).callonPgOidT1,
				expr: &choiceExpr{
					pos: position{line: 259, col: 13, offset: 9074},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 259, col: 13, offset: 9074},
							val:        "oid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 22, offset: 9083},
							val:        "regprocedure",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 40, offset: 9101},
							val:        "regproc",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 53, offset: 9114},
							val:        "regoperator",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 70, offset: 9131},
							val:        "regoper",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 83, offset: 9144},
							val:        "regclass",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 97, offset: 9158},
							val:        "regtype",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 110, offset: 9171},
							val:        "regrole",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 123, offset: 9184},
							val:        "regnamespace",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 141, offset: 9202},
							val:        "regconfig",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 259, col: 156, offset: 9217},
							val:        "regdictionary",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 265, col: 1, offset: 9331},
			expr: &actionExpr{
				pos: position{line: 265, col: 11, offset: 9341},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 265, col: 13, offset: 9343},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 265, col: 13, offset: 9343},
							val:        "date",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 23, offset: 9353},
							val:        "smallint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 37, offset: 9367},
							val:        "bigint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 49, offset: 9379},
							val:        "decimal",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 62, offset: 9392},
							val:        "numeric",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 75, offset: 9405},
							val:        "real",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 85, offset: 9415},
							val:        "smallserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 102, offset: 9432},
							val:        "serial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 114, offset: 9444},
							val:        "bigserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 129, offset: 9459},
							val:        "boolean",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 142, offset: 9472},
							val:        "text",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 152, offset: 9482},
							val:        "money",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 163, offset: 9493},
							val:        "bytea",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 174, offset: 9504},
							val:        "point",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 185, offset: 9515},
							val:        "line",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 195, offset: 9525},
							val:        "lseg",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 205, offset: 9535},
							val:        "box",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 214, offset: 9544},
							val:        "path",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 224, offset: 9554},
							val:        "polygon",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 237, offset: 9567},
							val:        "circle",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 249, offset: 9579},
							val:        "cidr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 259, offset: 9589},
							val:        "inet",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 269, offset: 9599},
							val:        "macaddr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 282, offset: 9612},
							val:        "uuid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 292, offset: 9622},
							val:        "xml",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 301, offset: 9631},
							val:        "jsonb",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 265, col: 312, offset: 9642},
							val:        "json",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CustomT",
			pos:  position{line: 271, col: 1, offset: 9747},
			expr: &actionExpr{
				pos: position{line: 271, col: 13, offset: 9759},
				run: (*parser).callonCustomT1,
				expr: &ruleRefExpr{
					pos:  position{line: 271, col: 13, offset: 9759},
					name: "Ident",
				},
			},
		},
		{
			name: "CreateSeqStmt",
			pos:  position{line: 292, col: 1, offset: 11224},
			expr: &actionExpr{
				pos: position{line: 292, col: 18, offset: 11241},
				run: (*parser).callonCreateSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 292, col: 18, offset: 11241},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 292, col: 18, offset: 11241},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 292, col: 28, offset: 11251},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 292, col: 31, offset: 11254},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 292, col: 43, offset: 11266},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 292, col: 46, offset: 11269},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 292, col: 51, offset: 11274},
								name: "Ident",
							},
						},
						&labeledExpr{
							pos:   position{line: 292, col: 57, offset: 11280},
							label: "verses",
							expr: &zeroOrMoreExpr{
								pos: position{line: 292, col: 64, offset: 11287},
								expr: &ruleRefExpr{
									pos:  position{line: 292, col: 64, offset: 11287},
									name: "CreateSeqVerse",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 292, col: 80, offset: 11303},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 292, col: 82, offset: 11305},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 292, col: 86, offset: 11309},
							expr: &ruleRefExpr{
								pos:  position{line: 292, col: 86, offset: 11309},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateSeqVerse",
			pos:  position{line: 306, col: 1, offset: 11703},
			expr: &actionExpr{
				pos: position{line: 306, col: 19, offset: 11721},
				run: (*parser).callonCreateSeqVerse1,
				expr: &labeledExpr{
					pos:   position{line: 306, col: 19, offset: 11721},
					label: "verse",
					expr: &choiceExpr{
						pos: position{line: 306, col: 27, offset: 11729},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 306, col: 27, offset: 11729},
								name: "IncrementBy",
							},
							&ruleRefExpr{
								pos:  position{line: 306, col: 41, offset: 11743},
								name: "MinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 306, col: 52, offset: 11754},
								name: "NoMinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 306, col: 65, offset: 11767},
								name: "MaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 306, col: 76, offset: 11778},
								name: "NoMaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 306, col: 89, offset: 11791},
								name: "Start",
							},
							&ruleRefExpr{
								pos:  position{line: 306, col: 97, offset: 11799},
								name: "Cache",
							},
							&ruleRefExpr{
								pos:  position{line: 306, col: 105, offset: 11807},
								name: "Cycle",
							},
							&ruleRefExpr{
								pos:  position{line: 306, col: 113, offset: 11815},
								name: "OwnedBy",
							},
						},
					},
				},
			},
		},
		{
			name: "IncrementBy",
			pos:  position{line: 310, col: 1, offset: 11852},
			expr: &actionExpr{
				pos: position{line: 310, col: 16, offset: 11867},
				run: (*parser).callonIncrementBy1,
				expr: &seqExpr{
					pos: position{line: 310, col: 16, offset: 11867},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 310, col: 16, offset: 11867},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 310, col: 19, offset: 11870},
							val:        "increment",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 310, col: 32, offset: 11883},
							expr: &seqExpr{
								pos: position{line: 310, col: 33, offset: 11884},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 310, col: 33, offset: 11884},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 310, col: 36, offset: 11887},
										val:        "by",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 310, col: 44, offset: 11895},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 310, col: 47, offset: 11898},
							label: "num",
							expr: &ruleRefExpr{
								pos:  position{line: 310, col: 51, offset: 11902},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "MinValue",
			pos:  position{line: 316, col: 1, offset: 12016},
			expr: &actionExpr{
				pos: position{line: 316, col: 13, offset: 12028},
				run: (*parser).callonMinValue1,
				expr: &seqExpr{
					pos: position{line: 316, col: 13, offset: 12028},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 316, col: 13, offset: 12028},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 316, col: 16, offset: 12031},
							val:        "minvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 316, col: 28, offset: 12043},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 316, col: 31, offset: 12046},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 316, col: 35, offset: 12050},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMinValue",
			pos:  position{line: 322, col: 1, offset: 12163},
			expr: &actionExpr{
				pos: position{line: 322, col: 15, offset: 12177},
				run: (*parser).callonNoMinValue1,
				expr: &seqExpr{
					pos: position{line: 322, col: 15, offset: 12177},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 322, col: 15, offset: 12177},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 322, col: 18, offset: 12180},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 322, col: 24, offset: 12186},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 322, col: 27, offset: 12189},
							val:        "minvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "MaxValue",
			pos:  position{line: 326, col: 1, offset: 12226},
			expr: &actionExpr{
				pos: position{line: 326, col: 13, offset: 12238},
				run: (*parser).callonMaxValue1,
				expr: &seqExpr{
					pos: position{line: 326, col: 13, offset: 12238},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 326, col: 13, offset: 12238},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 326, col: 16, offset: 12241},
							val:        "maxvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 326, col: 28, offset: 12253},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 326, col: 31, offset: 12256},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 326, col: 35, offset: 12260},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMaxValue",
			pos:  position{line: 332, col: 1, offset: 12373},
			expr: &actionExpr{
				pos: position{line: 332, col: 15, offset: 12387},
				run: (*parser).callonNoMaxValue1,
				expr: &seqExpr{
					pos: position{line: 332, col: 15, offset: 12387},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 332, col: 15, offset: 12387},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 332, col: 18, offset: 12390},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 332, col: 24, offset: 12396},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 332, col: 27, offset: 12399},
							val:        "maxvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "Start",
			pos:  position{line: 336, col: 1, offset: 12436},
			expr: &actionExpr{
				pos: position{line: 336, col: 10, offset: 12445},
				run: (*parser).callonStart1,
				expr: &seqExpr{
					pos: position{line: 336, col: 10, offset: 12445},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 336, col: 10, offset: 12445},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 336, col: 13, offset: 12448},
							val:        "start",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 336, col: 22, offset: 12457},
							expr: &seqExpr{
								pos: position{line: 336, col: 23, offset: 12458},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 336, col: 23, offset: 12458},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 336, col: 26, offset: 12461},
										val:        "with",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 336, col: 36, offset: 12471},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 336, col: 39, offset: 12474},
							label: "start",
							expr: &ruleRefExpr{
								pos:  position{line: 336, col: 45, offset: 12480},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cache",
			pos:  position{line: 342, col: 1, offset: 12592},
			expr: &actionExpr{
				pos: position{line: 342, col: 10, offset: 12601},
				run: (*parser).callonCache1,
				expr: &seqExpr{
					pos: position{line: 342, col: 10, offset: 12601},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 342, col: 10, offset: 12601},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 342, col: 13, offset: 12604},
							val:        "cache",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 342, col: 22, offset: 12613},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 342, col: 25, offset: 12616},
							label: "cache",
							expr: &ruleRefExpr{
								pos:  position{line: 342, col: 31, offset: 12622},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cycle",
			pos:  position{line: 348, col: 1, offset: 12734},
			expr: &actionExpr{
				pos: position{line: 348, col: 10, offset: 12743},
				run: (*parser).callonCycle1,
				expr: &seqExpr{
					pos: position{line: 348, col: 10, offset: 12743},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 348, col: 10, offset: 12743},
							label: "no",
							expr: &zeroOrOneExpr{
								pos: position{line: 348, col: 13, offset: 12746},
								expr: &seqExpr{
									pos: position{line: 348, col: 14, offset: 12747},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 348, col: 14, offset: 12747},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 348, col: 17, offset: 12750},
											val:        "no",
											ignoreCase: true,
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 348, col: 25, offset: 12758},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 348, col: 28, offset: 12761},
							val:        "cycle",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OwnedBy",
			pos:  position{line: 359, col: 1, offset: 12945},
			expr: &actionExpr{
				pos: position{line: 359, col: 12, offset: 12956},
				run: (*parser).callonOwnedBy1,
				expr: &seqExpr{
					pos: position{line: 359, col: 12, offset: 12956},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 359, col: 12, offset: 12956},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 359, col: 15, offset: 12959},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 359, col: 24, offset: 12968},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 359, col: 27, offset: 12971},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 359, col: 33, offset: 12977},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 359, col: 36, offset: 12980},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 359, col: 43, offset: 12987},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 359, col: 43, offset: 12987},
										val:        "none",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 359, col: 53, offset: 12997},
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
			pos:  position{line: 377, col: 1, offset: 14456},
			expr: &actionExpr{
				pos: position{line: 377, col: 19, offset: 14474},
				run: (*parser).callonCreateTypeStmt1,
				expr: &seqExpr{
					pos: position{line: 377, col: 19, offset: 14474},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 377, col: 19, offset: 14474},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 377, col: 29, offset: 14484},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 377, col: 32, offset: 14487},
							val:        "type",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 377, col: 40, offset: 14495},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 377, col: 43, offset: 14498},
							label: "typename",
							expr: &ruleRefExpr{
								pos:  position{line: 377, col: 52, offset: 14507},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 377, col: 58, offset: 14513},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 377, col: 61, offset: 14516},
							val:        "as",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 377, col: 67, offset: 14522},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 377, col: 70, offset: 14525},
							label: "typedef",
							expr: &ruleRefExpr{
								pos:  position{line: 377, col: 78, offset: 14533},
								name: "EnumDef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 377, col: 86, offset: 14541},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 377, col: 88, offset: 14543},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 377, col: 92, offset: 14547},
							expr: &ruleRefExpr{
								pos:  position{line: 377, col: 92, offset: 14547},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "EnumDef",
			pos:  position{line: 383, col: 1, offset: 14663},
			expr: &actionExpr{
				pos: position{line: 383, col: 12, offset: 14674},
				run: (*parser).callonEnumDef1,
				expr: &seqExpr{
					pos: position{line: 383, col: 12, offset: 14674},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 383, col: 12, offset: 14674},
							val:        "ENUM",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 383, col: 19, offset: 14681},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 383, col: 21, offset: 14683},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 383, col: 25, offset: 14687},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 383, col: 27, offset: 14689},
							label: "vals",
							expr: &seqExpr{
								pos: position{line: 383, col: 34, offset: 14696},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 383, col: 34, offset: 14696},
										name: "StringConst",
									},
									&zeroOrMoreExpr{
										pos: position{line: 383, col: 46, offset: 14708},
										expr: &seqExpr{
											pos: position{line: 383, col: 48, offset: 14710},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 383, col: 48, offset: 14710},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 383, col: 50, offset: 14712},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 383, col: 54, offset: 14716},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 383, col: 56, offset: 14718},
													name: "StringConst",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 383, col: 74, offset: 14736},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 383, col: 76, offset: 14738},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AlterTableStmt",
			pos:  position{line: 408, col: 1, offset: 16368},
			expr: &actionExpr{
				pos: position{line: 408, col: 19, offset: 16386},
				run: (*parser).callonAlterTableStmt1,
				expr: &seqExpr{
					pos: position{line: 408, col: 19, offset: 16386},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 408, col: 19, offset: 16386},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 408, col: 28, offset: 16395},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 408, col: 31, offset: 16398},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 408, col: 40, offset: 16407},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 408, col: 43, offset: 16410},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 408, col: 48, offset: 16415},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 408, col: 54, offset: 16421},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 408, col: 57, offset: 16424},
							val:        "owner",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 408, col: 66, offset: 16433},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 408, col: 69, offset: 16436},
							val:        "to",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 408, col: 75, offset: 16442},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 408, col: 78, offset: 16445},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 408, col: 84, offset: 16451},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 408, col: 90, offset: 16457},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 408, col: 92, offset: 16459},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 408, col: 96, offset: 16463},
							expr: &ruleRefExpr{
								pos:  position{line: 408, col: 96, offset: 16463},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "AlterSeqStmt",
			pos:  position{line: 422, col: 1, offset: 17609},
			expr: &actionExpr{
				pos: position{line: 422, col: 17, offset: 17625},
				run: (*parser).callonAlterSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 422, col: 17, offset: 17625},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 422, col: 17, offset: 17625},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 26, offset: 17634},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 422, col: 29, offset: 17637},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 41, offset: 17649},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 422, col: 44, offset: 17652},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 422, col: 49, offset: 17657},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 55, offset: 17663},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 422, col: 58, offset: 17666},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 67, offset: 17675},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 422, col: 70, offset: 17678},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 76, offset: 17684},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 422, col: 79, offset: 17687},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 422, col: 85, offset: 17693},
								name: "TableDotCol",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 422, col: 97, offset: 17705},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 422, col: 99, offset: 17707},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 422, col: 103, offset: 17711},
							expr: &ruleRefExpr{
								pos:  position{line: 422, col: 103, offset: 17711},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "TableDotCol",
			pos:  position{line: 426, col: 1, offset: 17790},
			expr: &actionExpr{
				pos: position{line: 426, col: 16, offset: 17805},
				run: (*parser).callonTableDotCol1,
				expr: &seqExpr{
					pos: position{line: 426, col: 16, offset: 17805},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 426, col: 16, offset: 17805},
							label: "table",
							expr: &ruleRefExpr{
								pos:  position{line: 426, col: 22, offset: 17811},
								name: "Ident",
							},
						},
						&litMatcher{
							pos:        position{line: 426, col: 28, offset: 17817},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 426, col: 32, offset: 17821},
							label: "column",
							expr: &ruleRefExpr{
								pos:  position{line: 426, col: 39, offset: 17828},
								name: "Ident",
							},
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 440, col: 1, offset: 19156},
			expr: &actionExpr{
				pos: position{line: 440, col: 25, offset: 19180},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 440, col: 25, offset: 19180},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 440, col: 25, offset: 19180},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 440, col: 36, offset: 19191},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 440, col: 39, offset: 19194},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 440, col: 45, offset: 19200},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 440, col: 48, offset: 19203},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 440, col: 61, offset: 19216},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 440, col: 63, offset: 19218},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 440, col: 73, offset: 19228},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 440, col: 79, offset: 19234},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 440, col: 81, offset: 19236},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 440, col: 87, offset: 19242},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 440, col: 89, offset: 19244},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 440, col: 97, offset: 19252},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 440, col: 109, offset: 19264},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 440, col: 111, offset: 19266},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 440, col: 115, offset: 19270},
							expr: &ruleRefExpr{
								pos:  position{line: 440, col: 115, offset: 19270},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 444, col: 1, offset: 19359},
			expr: &actionExpr{
				pos: position{line: 444, col: 24, offset: 19382},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 444, col: 24, offset: 19382},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 444, col: 24, offset: 19382},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 444, col: 34, offset: 19392},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 444, col: 37, offset: 19395},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 444, col: 50, offset: 19408},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 444, col: 53, offset: 19411},
							expr: &seqExpr{
								pos: position{line: 444, col: 55, offset: 19413},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 444, col: 55, offset: 19413},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 444, col: 61, offset: 19419},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 444, col: 64, offset: 19422},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 444, col: 71, offset: 19429},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 444, col: 74, offset: 19432},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 444, col: 84, offset: 19442},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 444, col: 90, offset: 19448},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 444, col: 100, offset: 19458},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 444, col: 106, offset: 19464},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 444, col: 109, offset: 19467},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 444, col: 117, offset: 19475},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 444, col: 120, offset: 19478},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 444, col: 130, offset: 19488},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 444, col: 133, offset: 19491},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 444, col: 140, offset: 19498},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 444, col: 146, offset: 19504},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 444, col: 148, offset: 19506},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 444, col: 152, offset: 19510},
							expr: &ruleRefExpr{
								pos:  position{line: 444, col: 152, offset: 19510},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 448, col: 1, offset: 19601},
			expr: &actionExpr{
				pos: position{line: 448, col: 12, offset: 19612},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 448, col: 12, offset: 19612},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 448, col: 12, offset: 19612},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 448, col: 19, offset: 19619},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 448, col: 21, offset: 19621},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 448, col: 25, offset: 19625},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 448, col: 29, offset: 19629},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 448, col: 33, offset: 19633},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 448, col: 33, offset: 19633},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 448, col: 39, offset: 19639},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 448, col: 47, offset: 19647},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 448, col: 49, offset: 19649},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 448, col: 56, offset: 19656},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 448, col: 77, offset: 19677},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 448, col: 79, offset: 19679},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 448, col: 83, offset: 19683},
							expr: &ruleRefExpr{
								pos:  position{line: 448, col: 83, offset: 19683},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 453, col: 1, offset: 19765},
			expr: &actionExpr{
				pos: position{line: 453, col: 8, offset: 19772},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 453, col: 8, offset: 19772},
					expr: &charClassMatcher{
						pos:        position{line: 453, col: 8, offset: 19772},
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
			pos:  position{line: 468, col: 1, offset: 20612},
			expr: &actionExpr{
				pos: position{line: 468, col: 25, offset: 20636},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 468, col: 25, offset: 20636},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 468, col: 32, offset: 20643},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 468, col: 32, offset: 20643},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 468, col: 38, offset: 20649},
								expr: &seqExpr{
									pos: position{line: 468, col: 40, offset: 20651},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 468, col: 40, offset: 20651},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 468, col: 42, offset: 20653},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 468, col: 46, offset: 20657},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 468, col: 48, offset: 20659},
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
			pos:  position{line: 480, col: 1, offset: 20949},
			expr: &choiceExpr{
				pos: position{line: 480, col: 12, offset: 20960},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 480, col: 12, offset: 20960},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 480, col: 21, offset: 20969},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 480, col: 31, offset: 20979},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 480, col: 45, offset: 20993},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 482, col: 1, offset: 21002},
			expr: &actionExpr{
				pos: position{line: 482, col: 16, offset: 21017},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 482, col: 16, offset: 21017},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 482, col: 16, offset: 21017},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 482, col: 20, offset: 21021},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 482, col: 26, offset: 21027},
								expr: &choiceExpr{
									pos: position{line: 482, col: 27, offset: 21028},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 482, col: 27, offset: 21028},
											val:        "[^'\\n]",
											chars:      []rune{'\'', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 482, col: 36, offset: 21037},
											val:        "''",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 482, col: 43, offset: 21044},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "DblQuotedString",
			pos:  position{line: 486, col: 1, offset: 21096},
			expr: &actionExpr{
				pos: position{line: 486, col: 20, offset: 21115},
				run: (*parser).callonDblQuotedString1,
				expr: &seqExpr{
					pos: position{line: 486, col: 20, offset: 21115},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 486, col: 20, offset: 21115},
							val:        "\"",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 486, col: 24, offset: 21119},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 486, col: 30, offset: 21125},
								expr: &choiceExpr{
									pos: position{line: 486, col: 31, offset: 21126},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 486, col: 31, offset: 21126},
											val:        "[^\"\\n]",
											chars:      []rune{'"', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 486, col: 40, offset: 21135},
											val:        "\"\"",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 486, col: 49, offset: 21144},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 490, col: 1, offset: 21208},
			expr: &actionExpr{
				pos: position{line: 490, col: 10, offset: 21217},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 490, col: 10, offset: 21217},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 490, col: 10, offset: 21217},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 490, col: 18, offset: 21225},
							expr: &charClassMatcher{
								pos:        position{line: 490, col: 18, offset: 21225},
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
			pos:  position{line: 494, col: 1, offset: 21278},
			expr: &actionExpr{
				pos: position{line: 494, col: 11, offset: 21288},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 494, col: 13, offset: 21290},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 494, col: 13, offset: 21290},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 494, col: 19, offset: 21296},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 494, col: 19, offset: 21296},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 494, col: 24, offset: 21301},
									expr: &charClassMatcher{
										pos:        position{line: 494, col: 24, offset: 21301},
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
			pos:  position{line: 499, col: 1, offset: 21396},
			expr: &actionExpr{
				pos: position{line: 499, col: 15, offset: 21410},
				run: (*parser).callonNonZNumber1,
				expr: &seqExpr{
					pos: position{line: 499, col: 15, offset: 21410},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 499, col: 15, offset: 21410},
							val:        "[1-9]",
							ranges:     []rune{'1', '9'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 499, col: 20, offset: 21415},
							expr: &charClassMatcher{
								pos:        position{line: 499, col: 20, offset: 21415},
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
			pos:  position{line: 504, col: 1, offset: 21508},
			expr: &actionExpr{
				pos: position{line: 504, col: 12, offset: 21519},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 504, col: 12, offset: 21519},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 504, col: 20, offset: 21527},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 504, col: 20, offset: 21527},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 504, col: 34, offset: 21541},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 508, col: 1, offset: 21583},
			expr: &actionExpr{
				pos: position{line: 508, col: 16, offset: 21598},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 508, col: 18, offset: 21600},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 508, col: 18, offset: 21600},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 508, col: 27, offset: 21609},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 508, col: 27, offset: 21609},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 508, col: 31, offset: 21613},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 508, col: 49, offset: 21631},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 508, col: 55, offset: 21637},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 512, col: 1, offset: 21683},
			expr: &choiceExpr{
				pos: position{line: 512, col: 24, offset: 21706},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 512, col: 24, offset: 21706},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 512, col: 33, offset: 21715},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 512, col: 41, offset: 21723},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 512, col: 48, offset: 21730},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 512, col: 54, offset: 21736},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 514, col: 1, offset: 21743},
			expr: &actionExpr{
				pos: position{line: 514, col: 17, offset: 21759},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 514, col: 19, offset: 21761},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 514, col: 19, offset: 21761},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 514, col: 29, offset: 21771},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 514, col: 29, offset: 21771},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 514, col: 33, offset: 21775},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 514, col: 52, offset: 21794},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 514, col: 58, offset: 21800},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 518, col: 1, offset: 21848},
			expr: &choiceExpr{
				pos: position{line: 518, col: 25, offset: 21872},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 518, col: 25, offset: 21872},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 518, col: 35, offset: 21882},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 518, col: 42, offset: 21889},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 518, col: 50, offset: 21897},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 518, col: 56, offset: 21903},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 530, col: 1, offset: 22418},
			expr: &actionExpr{
				pos: position{line: 530, col: 12, offset: 22429},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 530, col: 14, offset: 22431},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 530, col: 14, offset: 22431},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 530, col: 34, offset: 22451},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 534, col: 1, offset: 22495},
			expr: &seqExpr{
				pos: position{line: 534, col: 21, offset: 22515},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 534, col: 21, offset: 22515},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 534, col: 26, offset: 22520},
						expr: &anyMatcher{
							line: 534, col: 26, offset: 22520,
						},
					},
					&litMatcher{
						pos:        position{line: 534, col: 29, offset: 22523},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 534, col: 34, offset: 22528},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 536, col: 1, offset: 22533},
			expr: &seqExpr{
				pos: position{line: 536, col: 22, offset: 22554},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 536, col: 22, offset: 22554},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 536, col: 27, offset: 22559},
						expr: &charClassMatcher{
							pos:        position{line: 536, col: 27, offset: 22559},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 536, col: 36, offset: 22568},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 538, col: 1, offset: 22573},
			expr: &seqExpr{
				pos: position{line: 538, col: 9, offset: 22581},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 538, col: 9, offset: 22581},
						expr: &charClassMatcher{
							pos:        position{line: 538, col: 9, offset: 22581},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 538, col: 17, offset: 22589},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 538, col: 17, offset: 22589},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 538, col: 26, offset: 22598},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 538, col: 35, offset: 22607},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 538, col: 42, offset: 22614},
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
			pos:         position{line: 540, col: 1, offset: 22621},
			expr: &zeroOrMoreExpr{
				pos: position{line: 540, col: 19, offset: 22639},
				expr: &charClassMatcher{
					pos:        position{line: 540, col: 19, offset: 22639},
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
			pos:         position{line: 542, col: 1, offset: 22651},
			expr: &oneOrMoreExpr{
				pos: position{line: 542, col: 31, offset: 22681},
				expr: &charClassMatcher{
					pos:        position{line: 542, col: 31, offset: 22681},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 544, col: 1, offset: 22693},
			expr: &notExpr{
				pos: position{line: 544, col: 8, offset: 22700},
				expr: &anyMatcher{
					line: 544, col: 9, offset: 22701,
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

func (c *current) onGeographyT1(subtype, srid interface{}) (interface{}, error) {
	result := map[string]string{
		"type": "geography",
	}
	subtypeString := string(subtype.([]byte))
	result["geography_type"] = strings.ToLower(subtypeString)
	if srid != nil {
		sridSlice := toIfaceSlice(srid)
		result["srid"] = strconv.FormatInt(sridSlice[1].(int64), 10)
	}
	return result, nil
}

func (p *parser) callonGeographyT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGeographyT1(stack["subtype"], stack["srid"])
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
