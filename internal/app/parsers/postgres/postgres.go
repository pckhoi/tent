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
							label: "fields",
							expr: &seqExpr{
								pos: position{line: 28, col: 78, offset: 1992},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 28, col: 78, offset: 1992},
										name: "FieldDef",
									},
									&zeroOrMoreExpr{
										pos: position{line: 28, col: 87, offset: 2001},
										expr: &seqExpr{
											pos: position{line: 28, col: 89, offset: 2003},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 28, col: 89, offset: 2003},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 28, col: 91, offset: 2005},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 28, col: 95, offset: 2009},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 28, col: 97, offset: 2011},
													name: "FieldDef",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 111, offset: 2025},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 28, col: 113, offset: 2027},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 117, offset: 2031},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 28, col: 119, offset: 2033},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 28, col: 123, offset: 2037},
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 123, offset: 2037},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "FieldDef",
			pos:  position{line: 48, col: 1, offset: 2656},
			expr: &actionExpr{
				pos: position{line: 48, col: 13, offset: 2668},
				run: (*parser).callonFieldDef1,
				expr: &seqExpr{
					pos: position{line: 48, col: 13, offset: 2668},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 48, col: 13, offset: 2668},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 48, col: 18, offset: 2673},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 48, col: 24, offset: 2679},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 48, col: 27, offset: 2682},
							label: "dataType",
							expr: &ruleRefExpr{
								pos:  position{line: 48, col: 36, offset: 2691},
								name: "DataType",
							},
						},
						&labeledExpr{
							pos:   position{line: 48, col: 45, offset: 2700},
							label: "constraint",
							expr: &zeroOrOneExpr{
								pos: position{line: 48, col: 56, offset: 2711},
								expr: &ruleRefExpr{
									pos:  position{line: 48, col: 56, offset: 2711},
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
			pos:  position{line: 62, col: 1, offset: 3076},
			expr: &actionExpr{
				pos: position{line: 62, col: 21, offset: 3096},
				run: (*parser).callonColumnConstraint1,
				expr: &seqExpr{
					pos: position{line: 62, col: 21, offset: 3096},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 62, col: 21, offset: 3096},
							label: "nameOpt",
							expr: &zeroOrOneExpr{
								pos: position{line: 62, col: 29, offset: 3104},
								expr: &seqExpr{
									pos: position{line: 62, col: 31, offset: 3106},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 62, col: 31, offset: 3106},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 62, col: 34, offset: 3109},
											val:        "constraint",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 62, col: 48, offset: 3123},
											name: "_1",
										},
										&choiceExpr{
											pos: position{line: 62, col: 52, offset: 3127},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 62, col: 52, offset: 3127},
													name: "StringConst",
												},
												&ruleRefExpr{
													pos:  position{line: 62, col: 66, offset: 3141},
													name: "Ident",
												},
											},
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 62, col: 76, offset: 3151},
							label: "constraintClauses",
							expr: &oneOrMoreExpr{
								pos: position{line: 62, col: 94, offset: 3169},
								expr: &choiceExpr{
									pos: position{line: 62, col: 96, offset: 3171},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 62, col: 96, offset: 3171},
											name: "NotNullCls",
										},
										&ruleRefExpr{
											pos:  position{line: 62, col: 109, offset: 3184},
											name: "NullCls",
										},
										&ruleRefExpr{
											pos:  position{line: 62, col: 119, offset: 3194},
											name: "CheckCls",
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
			name: "NotNullCls",
			pos:  position{line: 82, col: 1, offset: 3726},
			expr: &actionExpr{
				pos: position{line: 82, col: 15, offset: 3740},
				run: (*parser).callonNotNullCls1,
				expr: &seqExpr{
					pos: position{line: 82, col: 15, offset: 3740},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 82, col: 15, offset: 3740},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 82, col: 18, offset: 3743},
							val:        "not",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 82, col: 25, offset: 3750},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 82, col: 28, offset: 3753},
							val:        "null",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "NullCls",
			pos:  position{line: 88, col: 1, offset: 3835},
			expr: &actionExpr{
				pos: position{line: 88, col: 12, offset: 3846},
				run: (*parser).callonNullCls1,
				expr: &seqExpr{
					pos: position{line: 88, col: 12, offset: 3846},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 88, col: 12, offset: 3846},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 88, col: 15, offset: 3849},
							val:        "null",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CheckCls",
			pos:  position{line: 94, col: 1, offset: 3932},
			expr: &actionExpr{
				pos: position{line: 94, col: 13, offset: 3944},
				run: (*parser).callonCheckCls1,
				expr: &seqExpr{
					pos: position{line: 94, col: 13, offset: 3944},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 94, col: 13, offset: 3944},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 94, col: 16, offset: 3947},
							val:        "check",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 94, col: 25, offset: 3956},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 94, col: 28, offset: 3959},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 94, col: 33, offset: 3964},
								name: "WrappedExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 94, col: 45, offset: 3976},
							label: "noInherit",
							expr: &zeroOrOneExpr{
								pos: position{line: 94, col: 55, offset: 3986},
								expr: &seqExpr{
									pos: position{line: 94, col: 57, offset: 3988},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 94, col: 57, offset: 3988},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 94, col: 60, offset: 3991},
											val:        "no",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 94, col: 66, offset: 3997},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 94, col: 69, offset: 4000},
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
			pos:  position{line: 104, col: 1, offset: 4193},
			expr: &actionExpr{
				pos: position{line: 104, col: 16, offset: 4208},
				run: (*parser).callonWrappedExpr1,
				expr: &seqExpr{
					pos: position{line: 104, col: 16, offset: 4208},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 104, col: 16, offset: 4208},
							val:        "(",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 104, col: 20, offset: 4212},
							expr: &ruleRefExpr{
								pos:  position{line: 104, col: 20, offset: 4212},
								name: "Expr",
							},
						},
						&litMatcher{
							pos:        position{line: 104, col: 26, offset: 4218},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 108, col: 1, offset: 4258},
			expr: &choiceExpr{
				pos: position{line: 108, col: 9, offset: 4266},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 108, col: 9, offset: 4266},
						name: "WrappedExpr",
					},
					&oneOrMoreExpr{
						pos: position{line: 108, col: 23, offset: 4280},
						expr: &charClassMatcher{
							pos:        position{line: 108, col: 23, offset: 4280},
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
			pos:  position{line: 120, col: 1, offset: 5300},
			expr: &actionExpr{
				pos: position{line: 120, col: 13, offset: 5312},
				run: (*parser).callonDataType1,
				expr: &labeledExpr{
					pos:   position{line: 120, col: 13, offset: 5312},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 120, col: 17, offset: 5316},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 120, col: 17, offset: 5316},
								name: "TimestampT",
							},
							&ruleRefExpr{
								pos:  position{line: 120, col: 30, offset: 5329},
								name: "TimeT",
							},
							&ruleRefExpr{
								pos:  position{line: 120, col: 38, offset: 5337},
								name: "VarcharT",
							},
							&ruleRefExpr{
								pos:  position{line: 120, col: 49, offset: 5348},
								name: "CharT",
							},
							&ruleRefExpr{
								pos:  position{line: 120, col: 57, offset: 5356},
								name: "BitVarT",
							},
							&ruleRefExpr{
								pos:  position{line: 120, col: 67, offset: 5366},
								name: "BitT",
							},
							&ruleRefExpr{
								pos:  position{line: 120, col: 74, offset: 5373},
								name: "IntT",
							},
							&ruleRefExpr{
								pos:  position{line: 120, col: 81, offset: 5380},
								name: "PgOidT",
							},
							&ruleRefExpr{
								pos:  position{line: 120, col: 90, offset: 5389},
								name: "OtherT",
							},
							&ruleRefExpr{
								pos:  position{line: 120, col: 99, offset: 5398},
								name: "CustomT",
							},
						},
					},
				},
			},
		},
		{
			name: "TimestampT",
			pos:  position{line: 124, col: 1, offset: 5431},
			expr: &actionExpr{
				pos: position{line: 124, col: 15, offset: 5445},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 124, col: 15, offset: 5445},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 124, col: 15, offset: 5445},
							val:        "timestamp",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 124, col: 28, offset: 5458},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 124, col: 33, offset: 5463},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 124, col: 46, offset: 5476},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 124, col: 59, offset: 5489},
								expr: &choiceExpr{
									pos: position{line: 124, col: 61, offset: 5491},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 124, col: 61, offset: 5491},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 124, col: 70, offset: 5500},
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
			pos:  position{line: 137, col: 1, offset: 5779},
			expr: &actionExpr{
				pos: position{line: 137, col: 10, offset: 5788},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 137, col: 10, offset: 5788},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 137, col: 10, offset: 5788},
							val:        "time",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 137, col: 18, offset: 5796},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 137, col: 23, offset: 5801},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 137, col: 36, offset: 5814},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 137, col: 49, offset: 5827},
								expr: &choiceExpr{
									pos: position{line: 137, col: 51, offset: 5829},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 137, col: 51, offset: 5829},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 137, col: 60, offset: 5838},
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
			pos:  position{line: 150, col: 1, offset: 6109},
			expr: &actionExpr{
				pos: position{line: 150, col: 17, offset: 6125},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 150, col: 17, offset: 6125},
					expr: &seqExpr{
						pos: position{line: 150, col: 19, offset: 6127},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 150, col: 19, offset: 6127},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 150, col: 22, offset: 6130},
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
			pos:  position{line: 157, col: 1, offset: 6258},
			expr: &actionExpr{
				pos: position{line: 157, col: 11, offset: 6268},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 157, col: 11, offset: 6268},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 157, col: 11, offset: 6268},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 157, col: 14, offset: 6271},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 157, col: 22, offset: 6279},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 157, col: 25, offset: 6282},
							val:        "time",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 157, col: 33, offset: 6290},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 157, col: 36, offset: 6293},
							val:        "zone",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 161, col: 1, offset: 6327},
			expr: &actionExpr{
				pos: position{line: 161, col: 14, offset: 6340},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 161, col: 14, offset: 6340},
					expr: &seqExpr{
						pos: position{line: 161, col: 16, offset: 6342},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 161, col: 16, offset: 6342},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 161, col: 19, offset: 6345},
								val:        "without",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 161, col: 30, offset: 6356},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 161, col: 33, offset: 6359},
								val:        "time",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 161, col: 41, offset: 6367},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 161, col: 44, offset: 6370},
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
			pos:  position{line: 165, col: 1, offset: 6408},
			expr: &actionExpr{
				pos: position{line: 165, col: 10, offset: 6417},
				run: (*parser).callonCharT1,
				expr: &seqExpr{
					pos: position{line: 165, col: 10, offset: 6417},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 165, col: 12, offset: 6419},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 165, col: 12, offset: 6419},
									val:        "character",
									ignoreCase: true,
								},
								&litMatcher{
									pos:        position{line: 165, col: 27, offset: 6434},
									val:        "char",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 165, col: 37, offset: 6444},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 165, col: 44, offset: 6451},
								expr: &seqExpr{
									pos: position{line: 165, col: 46, offset: 6453},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 165, col: 46, offset: 6453},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 165, col: 50, offset: 6457},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 165, col: 61, offset: 6468},
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
			pos:  position{line: 177, col: 1, offset: 6723},
			expr: &actionExpr{
				pos: position{line: 177, col: 13, offset: 6735},
				run: (*parser).callonVarcharT1,
				expr: &seqExpr{
					pos: position{line: 177, col: 13, offset: 6735},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 177, col: 15, offset: 6737},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 177, col: 17, offset: 6739},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 177, col: 17, offset: 6739},
											val:        "character",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 177, col: 30, offset: 6752},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 177, col: 33, offset: 6755},
											val:        "varying",
											ignoreCase: true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 177, col: 48, offset: 6770},
									val:        "varchar",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 177, col: 61, offset: 6783},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 177, col: 68, offset: 6790},
								expr: &seqExpr{
									pos: position{line: 177, col: 70, offset: 6792},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 177, col: 70, offset: 6792},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 177, col: 74, offset: 6796},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 177, col: 85, offset: 6807},
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
			pos:  position{line: 188, col: 1, offset: 7042},
			expr: &actionExpr{
				pos: position{line: 188, col: 9, offset: 7050},
				run: (*parser).callonBitT1,
				expr: &seqExpr{
					pos: position{line: 188, col: 9, offset: 7050},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 188, col: 9, offset: 7050},
							val:        "bit",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 188, col: 16, offset: 7057},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 188, col: 23, offset: 7064},
								expr: &seqExpr{
									pos: position{line: 188, col: 25, offset: 7066},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 188, col: 25, offset: 7066},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 29, offset: 7070},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 188, col: 40, offset: 7081},
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
			pos:  position{line: 200, col: 1, offset: 7335},
			expr: &actionExpr{
				pos: position{line: 200, col: 12, offset: 7346},
				run: (*parser).callonBitVarT1,
				expr: &seqExpr{
					pos: position{line: 200, col: 12, offset: 7346},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 200, col: 12, offset: 7346},
							val:        "bit",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 200, col: 19, offset: 7353},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 200, col: 22, offset: 7356},
							val:        "varying",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 200, col: 33, offset: 7367},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 200, col: 40, offset: 7374},
								expr: &seqExpr{
									pos: position{line: 200, col: 42, offset: 7376},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 200, col: 42, offset: 7376},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 200, col: 46, offset: 7380},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 200, col: 57, offset: 7391},
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
			pos:  position{line: 211, col: 1, offset: 7625},
			expr: &actionExpr{
				pos: position{line: 211, col: 9, offset: 7633},
				run: (*parser).callonIntT1,
				expr: &choiceExpr{
					pos: position{line: 211, col: 11, offset: 7635},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 211, col: 11, offset: 7635},
							val:        "integer",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 211, col: 24, offset: 7648},
							val:        "int",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "PgOidT",
			pos:  position{line: 217, col: 1, offset: 7730},
			expr: &actionExpr{
				pos: position{line: 217, col: 11, offset: 7740},
				run: (*parser).callonPgOidT1,
				expr: &choiceExpr{
					pos: position{line: 217, col: 13, offset: 7742},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 217, col: 13, offset: 7742},
							val:        "oid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 22, offset: 7751},
							val:        "regprocedure",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 40, offset: 7769},
							val:        "regproc",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 53, offset: 7782},
							val:        "regoperator",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 70, offset: 7799},
							val:        "regoper",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 83, offset: 7812},
							val:        "regclass",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 97, offset: 7826},
							val:        "regtype",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 110, offset: 7839},
							val:        "regrole",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 123, offset: 7852},
							val:        "regnamespace",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 141, offset: 7870},
							val:        "regconfig",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 217, col: 156, offset: 7885},
							val:        "regdictionary",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 223, col: 1, offset: 7999},
			expr: &actionExpr{
				pos: position{line: 223, col: 11, offset: 8009},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 223, col: 13, offset: 8011},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 223, col: 13, offset: 8011},
							val:        "date",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 23, offset: 8021},
							val:        "smallint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 37, offset: 8035},
							val:        "bigint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 49, offset: 8047},
							val:        "decimal",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 62, offset: 8060},
							val:        "numeric",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 75, offset: 8073},
							val:        "real",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 85, offset: 8083},
							val:        "smallserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 102, offset: 8100},
							val:        "serial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 114, offset: 8112},
							val:        "bigserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 129, offset: 8127},
							val:        "boolean",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 142, offset: 8140},
							val:        "text",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 152, offset: 8150},
							val:        "money",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 163, offset: 8161},
							val:        "bytea",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 174, offset: 8172},
							val:        "point",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 185, offset: 8183},
							val:        "line",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 195, offset: 8193},
							val:        "lseg",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 205, offset: 8203},
							val:        "box",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 214, offset: 8212},
							val:        "path",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 224, offset: 8222},
							val:        "polygon",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 237, offset: 8235},
							val:        "circle",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 249, offset: 8247},
							val:        "cidr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 259, offset: 8257},
							val:        "inet",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 269, offset: 8267},
							val:        "macaddr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 282, offset: 8280},
							val:        "uuid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 292, offset: 8290},
							val:        "xml",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 301, offset: 8299},
							val:        "jsonb",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 223, col: 312, offset: 8310},
							val:        "json",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CustomT",
			pos:  position{line: 229, col: 1, offset: 8415},
			expr: &actionExpr{
				pos: position{line: 229, col: 13, offset: 8427},
				run: (*parser).callonCustomT1,
				expr: &ruleRefExpr{
					pos:  position{line: 229, col: 13, offset: 8427},
					name: "Ident",
				},
			},
		},
		{
			name: "CreateSeqStmt",
			pos:  position{line: 250, col: 1, offset: 9892},
			expr: &actionExpr{
				pos: position{line: 250, col: 18, offset: 9909},
				run: (*parser).callonCreateSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 250, col: 18, offset: 9909},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 250, col: 18, offset: 9909},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 250, col: 28, offset: 9919},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 250, col: 31, offset: 9922},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 250, col: 43, offset: 9934},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 250, col: 46, offset: 9937},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 250, col: 51, offset: 9942},
								name: "Ident",
							},
						},
						&labeledExpr{
							pos:   position{line: 250, col: 57, offset: 9948},
							label: "verses",
							expr: &zeroOrMoreExpr{
								pos: position{line: 250, col: 64, offset: 9955},
								expr: &ruleRefExpr{
									pos:  position{line: 250, col: 64, offset: 9955},
									name: "CreateSeqVerse",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 250, col: 80, offset: 9971},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 250, col: 82, offset: 9973},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 250, col: 86, offset: 9977},
							expr: &ruleRefExpr{
								pos:  position{line: 250, col: 86, offset: 9977},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateSeqVerse",
			pos:  position{line: 264, col: 1, offset: 10371},
			expr: &actionExpr{
				pos: position{line: 264, col: 19, offset: 10389},
				run: (*parser).callonCreateSeqVerse1,
				expr: &labeledExpr{
					pos:   position{line: 264, col: 19, offset: 10389},
					label: "verse",
					expr: &choiceExpr{
						pos: position{line: 264, col: 27, offset: 10397},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 264, col: 27, offset: 10397},
								name: "IncrementBy",
							},
							&ruleRefExpr{
								pos:  position{line: 264, col: 41, offset: 10411},
								name: "MinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 264, col: 52, offset: 10422},
								name: "NoMinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 264, col: 65, offset: 10435},
								name: "MaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 264, col: 76, offset: 10446},
								name: "NoMaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 264, col: 89, offset: 10459},
								name: "Start",
							},
							&ruleRefExpr{
								pos:  position{line: 264, col: 97, offset: 10467},
								name: "Cache",
							},
							&ruleRefExpr{
								pos:  position{line: 264, col: 105, offset: 10475},
								name: "Cycle",
							},
							&ruleRefExpr{
								pos:  position{line: 264, col: 113, offset: 10483},
								name: "OwnedBy",
							},
						},
					},
				},
			},
		},
		{
			name: "IncrementBy",
			pos:  position{line: 268, col: 1, offset: 10520},
			expr: &actionExpr{
				pos: position{line: 268, col: 16, offset: 10535},
				run: (*parser).callonIncrementBy1,
				expr: &seqExpr{
					pos: position{line: 268, col: 16, offset: 10535},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 268, col: 16, offset: 10535},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 268, col: 19, offset: 10538},
							val:        "increment",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 268, col: 32, offset: 10551},
							expr: &seqExpr{
								pos: position{line: 268, col: 33, offset: 10552},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 268, col: 33, offset: 10552},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 268, col: 36, offset: 10555},
										val:        "by",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 268, col: 44, offset: 10563},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 268, col: 47, offset: 10566},
							label: "num",
							expr: &ruleRefExpr{
								pos:  position{line: 268, col: 51, offset: 10570},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "MinValue",
			pos:  position{line: 274, col: 1, offset: 10684},
			expr: &actionExpr{
				pos: position{line: 274, col: 13, offset: 10696},
				run: (*parser).callonMinValue1,
				expr: &seqExpr{
					pos: position{line: 274, col: 13, offset: 10696},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 274, col: 13, offset: 10696},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 274, col: 16, offset: 10699},
							val:        "minvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 274, col: 28, offset: 10711},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 274, col: 31, offset: 10714},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 274, col: 35, offset: 10718},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMinValue",
			pos:  position{line: 280, col: 1, offset: 10831},
			expr: &actionExpr{
				pos: position{line: 280, col: 15, offset: 10845},
				run: (*parser).callonNoMinValue1,
				expr: &seqExpr{
					pos: position{line: 280, col: 15, offset: 10845},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 280, col: 15, offset: 10845},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 280, col: 18, offset: 10848},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 280, col: 24, offset: 10854},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 280, col: 27, offset: 10857},
							val:        "minvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "MaxValue",
			pos:  position{line: 284, col: 1, offset: 10894},
			expr: &actionExpr{
				pos: position{line: 284, col: 13, offset: 10906},
				run: (*parser).callonMaxValue1,
				expr: &seqExpr{
					pos: position{line: 284, col: 13, offset: 10906},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 284, col: 13, offset: 10906},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 284, col: 16, offset: 10909},
							val:        "maxvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 284, col: 28, offset: 10921},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 284, col: 31, offset: 10924},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 284, col: 35, offset: 10928},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMaxValue",
			pos:  position{line: 290, col: 1, offset: 11041},
			expr: &actionExpr{
				pos: position{line: 290, col: 15, offset: 11055},
				run: (*parser).callonNoMaxValue1,
				expr: &seqExpr{
					pos: position{line: 290, col: 15, offset: 11055},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 290, col: 15, offset: 11055},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 290, col: 18, offset: 11058},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 290, col: 24, offset: 11064},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 290, col: 27, offset: 11067},
							val:        "maxvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "Start",
			pos:  position{line: 294, col: 1, offset: 11104},
			expr: &actionExpr{
				pos: position{line: 294, col: 10, offset: 11113},
				run: (*parser).callonStart1,
				expr: &seqExpr{
					pos: position{line: 294, col: 10, offset: 11113},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 294, col: 10, offset: 11113},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 294, col: 13, offset: 11116},
							val:        "start",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 294, col: 22, offset: 11125},
							expr: &seqExpr{
								pos: position{line: 294, col: 23, offset: 11126},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 294, col: 23, offset: 11126},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 294, col: 26, offset: 11129},
										val:        "with",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 294, col: 36, offset: 11139},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 294, col: 39, offset: 11142},
							label: "start",
							expr: &ruleRefExpr{
								pos:  position{line: 294, col: 45, offset: 11148},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cache",
			pos:  position{line: 300, col: 1, offset: 11260},
			expr: &actionExpr{
				pos: position{line: 300, col: 10, offset: 11269},
				run: (*parser).callonCache1,
				expr: &seqExpr{
					pos: position{line: 300, col: 10, offset: 11269},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 300, col: 10, offset: 11269},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 300, col: 13, offset: 11272},
							val:        "cache",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 300, col: 22, offset: 11281},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 300, col: 25, offset: 11284},
							label: "cache",
							expr: &ruleRefExpr{
								pos:  position{line: 300, col: 31, offset: 11290},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cycle",
			pos:  position{line: 306, col: 1, offset: 11402},
			expr: &actionExpr{
				pos: position{line: 306, col: 10, offset: 11411},
				run: (*parser).callonCycle1,
				expr: &seqExpr{
					pos: position{line: 306, col: 10, offset: 11411},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 306, col: 10, offset: 11411},
							label: "no",
							expr: &zeroOrOneExpr{
								pos: position{line: 306, col: 13, offset: 11414},
								expr: &seqExpr{
									pos: position{line: 306, col: 14, offset: 11415},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 306, col: 14, offset: 11415},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 306, col: 17, offset: 11418},
											val:        "no",
											ignoreCase: true,
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 306, col: 25, offset: 11426},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 306, col: 28, offset: 11429},
							val:        "cycle",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OwnedBy",
			pos:  position{line: 317, col: 1, offset: 11613},
			expr: &actionExpr{
				pos: position{line: 317, col: 12, offset: 11624},
				run: (*parser).callonOwnedBy1,
				expr: &seqExpr{
					pos: position{line: 317, col: 12, offset: 11624},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 317, col: 12, offset: 11624},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 317, col: 15, offset: 11627},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 24, offset: 11636},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 317, col: 27, offset: 11639},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 33, offset: 11645},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 317, col: 36, offset: 11648},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 317, col: 43, offset: 11655},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 317, col: 43, offset: 11655},
										val:        "none",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 317, col: 53, offset: 11665},
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
			pos:  position{line: 335, col: 1, offset: 13124},
			expr: &actionExpr{
				pos: position{line: 335, col: 19, offset: 13142},
				run: (*parser).callonCreateTypeStmt1,
				expr: &seqExpr{
					pos: position{line: 335, col: 19, offset: 13142},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 335, col: 19, offset: 13142},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 335, col: 29, offset: 13152},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 335, col: 32, offset: 13155},
							val:        "type",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 335, col: 40, offset: 13163},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 335, col: 43, offset: 13166},
							label: "typename",
							expr: &ruleRefExpr{
								pos:  position{line: 335, col: 52, offset: 13175},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 335, col: 58, offset: 13181},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 335, col: 61, offset: 13184},
							val:        "as",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 335, col: 67, offset: 13190},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 335, col: 70, offset: 13193},
							label: "typedef",
							expr: &ruleRefExpr{
								pos:  position{line: 335, col: 78, offset: 13201},
								name: "EnumDef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 335, col: 86, offset: 13209},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 335, col: 88, offset: 13211},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 335, col: 92, offset: 13215},
							expr: &ruleRefExpr{
								pos:  position{line: 335, col: 92, offset: 13215},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "EnumDef",
			pos:  position{line: 341, col: 1, offset: 13331},
			expr: &actionExpr{
				pos: position{line: 341, col: 12, offset: 13342},
				run: (*parser).callonEnumDef1,
				expr: &seqExpr{
					pos: position{line: 341, col: 12, offset: 13342},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 341, col: 12, offset: 13342},
							val:        "ENUM",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 341, col: 19, offset: 13349},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 341, col: 21, offset: 13351},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 341, col: 25, offset: 13355},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 341, col: 27, offset: 13357},
							label: "vals",
							expr: &seqExpr{
								pos: position{line: 341, col: 34, offset: 13364},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 341, col: 34, offset: 13364},
										name: "StringConst",
									},
									&zeroOrMoreExpr{
										pos: position{line: 341, col: 46, offset: 13376},
										expr: &seqExpr{
											pos: position{line: 341, col: 48, offset: 13378},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 341, col: 48, offset: 13378},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 341, col: 50, offset: 13380},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 341, col: 54, offset: 13384},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 341, col: 56, offset: 13386},
													name: "StringConst",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 341, col: 74, offset: 13404},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 341, col: 76, offset: 13406},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AlterTableStmt",
			pos:  position{line: 366, col: 1, offset: 15036},
			expr: &actionExpr{
				pos: position{line: 366, col: 19, offset: 15054},
				run: (*parser).callonAlterTableStmt1,
				expr: &seqExpr{
					pos: position{line: 366, col: 19, offset: 15054},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 366, col: 19, offset: 15054},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 366, col: 28, offset: 15063},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 366, col: 31, offset: 15066},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 366, col: 40, offset: 15075},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 366, col: 43, offset: 15078},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 366, col: 48, offset: 15083},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 366, col: 54, offset: 15089},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 366, col: 57, offset: 15092},
							val:        "owner",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 366, col: 66, offset: 15101},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 366, col: 69, offset: 15104},
							val:        "to",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 366, col: 75, offset: 15110},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 366, col: 78, offset: 15113},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 366, col: 84, offset: 15119},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 366, col: 90, offset: 15125},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 366, col: 92, offset: 15127},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 366, col: 96, offset: 15131},
							expr: &ruleRefExpr{
								pos:  position{line: 366, col: 96, offset: 15131},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "AlterSeqStmt",
			pos:  position{line: 380, col: 1, offset: 16277},
			expr: &actionExpr{
				pos: position{line: 380, col: 17, offset: 16293},
				run: (*parser).callonAlterSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 380, col: 17, offset: 16293},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 380, col: 17, offset: 16293},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 380, col: 26, offset: 16302},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 380, col: 29, offset: 16305},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 380, col: 41, offset: 16317},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 380, col: 44, offset: 16320},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 380, col: 49, offset: 16325},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 380, col: 55, offset: 16331},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 380, col: 58, offset: 16334},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 380, col: 67, offset: 16343},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 380, col: 70, offset: 16346},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 380, col: 76, offset: 16352},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 380, col: 79, offset: 16355},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 380, col: 85, offset: 16361},
								name: "TableDotCol",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 380, col: 97, offset: 16373},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 380, col: 99, offset: 16375},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 380, col: 103, offset: 16379},
							expr: &ruleRefExpr{
								pos:  position{line: 380, col: 103, offset: 16379},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "TableDotCol",
			pos:  position{line: 384, col: 1, offset: 16458},
			expr: &actionExpr{
				pos: position{line: 384, col: 16, offset: 16473},
				run: (*parser).callonTableDotCol1,
				expr: &seqExpr{
					pos: position{line: 384, col: 16, offset: 16473},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 384, col: 16, offset: 16473},
							label: "table",
							expr: &ruleRefExpr{
								pos:  position{line: 384, col: 22, offset: 16479},
								name: "Ident",
							},
						},
						&litMatcher{
							pos:        position{line: 384, col: 28, offset: 16485},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 384, col: 32, offset: 16489},
							label: "column",
							expr: &ruleRefExpr{
								pos:  position{line: 384, col: 39, offset: 16496},
								name: "Ident",
							},
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 398, col: 1, offset: 17824},
			expr: &actionExpr{
				pos: position{line: 398, col: 25, offset: 17848},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 398, col: 25, offset: 17848},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 398, col: 25, offset: 17848},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 398, col: 36, offset: 17859},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 398, col: 39, offset: 17862},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 398, col: 45, offset: 17868},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 398, col: 48, offset: 17871},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 398, col: 61, offset: 17884},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 398, col: 63, offset: 17886},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 398, col: 73, offset: 17896},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 398, col: 79, offset: 17902},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 398, col: 81, offset: 17904},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 398, col: 87, offset: 17910},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 398, col: 89, offset: 17912},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 398, col: 97, offset: 17920},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 398, col: 109, offset: 17932},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 398, col: 111, offset: 17934},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 398, col: 115, offset: 17938},
							expr: &ruleRefExpr{
								pos:  position{line: 398, col: 115, offset: 17938},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 402, col: 1, offset: 18027},
			expr: &actionExpr{
				pos: position{line: 402, col: 24, offset: 18050},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 402, col: 24, offset: 18050},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 402, col: 24, offset: 18050},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 34, offset: 18060},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 402, col: 37, offset: 18063},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 50, offset: 18076},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 402, col: 53, offset: 18079},
							expr: &seqExpr{
								pos: position{line: 402, col: 55, offset: 18081},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 402, col: 55, offset: 18081},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 402, col: 61, offset: 18087},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 402, col: 64, offset: 18090},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 402, col: 71, offset: 18097},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 402, col: 74, offset: 18100},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 402, col: 84, offset: 18110},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 402, col: 90, offset: 18116},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 402, col: 100, offset: 18126},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 106, offset: 18132},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 402, col: 109, offset: 18135},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 117, offset: 18143},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 402, col: 120, offset: 18146},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 130, offset: 18156},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 402, col: 133, offset: 18159},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 402, col: 140, offset: 18166},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 146, offset: 18172},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 402, col: 148, offset: 18174},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 402, col: 152, offset: 18178},
							expr: &ruleRefExpr{
								pos:  position{line: 402, col: 152, offset: 18178},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 406, col: 1, offset: 18269},
			expr: &actionExpr{
				pos: position{line: 406, col: 12, offset: 18280},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 406, col: 12, offset: 18280},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 406, col: 12, offset: 18280},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 406, col: 19, offset: 18287},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 406, col: 21, offset: 18289},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 406, col: 25, offset: 18293},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 406, col: 29, offset: 18297},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 406, col: 33, offset: 18301},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 406, col: 33, offset: 18301},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 406, col: 39, offset: 18307},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 406, col: 47, offset: 18315},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 406, col: 49, offset: 18317},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 406, col: 56, offset: 18324},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 406, col: 77, offset: 18345},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 406, col: 79, offset: 18347},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 406, col: 83, offset: 18351},
							expr: &ruleRefExpr{
								pos:  position{line: 406, col: 83, offset: 18351},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 411, col: 1, offset: 18433},
			expr: &actionExpr{
				pos: position{line: 411, col: 8, offset: 18440},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 411, col: 8, offset: 18440},
					expr: &charClassMatcher{
						pos:        position{line: 411, col: 8, offset: 18440},
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
			pos:  position{line: 426, col: 1, offset: 19280},
			expr: &actionExpr{
				pos: position{line: 426, col: 25, offset: 19304},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 426, col: 25, offset: 19304},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 426, col: 32, offset: 19311},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 426, col: 32, offset: 19311},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 426, col: 38, offset: 19317},
								expr: &seqExpr{
									pos: position{line: 426, col: 40, offset: 19319},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 426, col: 40, offset: 19319},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 426, col: 42, offset: 19321},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 426, col: 46, offset: 19325},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 426, col: 48, offset: 19327},
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
			pos:  position{line: 438, col: 1, offset: 19617},
			expr: &choiceExpr{
				pos: position{line: 438, col: 12, offset: 19628},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 438, col: 12, offset: 19628},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 438, col: 21, offset: 19637},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 438, col: 31, offset: 19647},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 438, col: 45, offset: 19661},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 440, col: 1, offset: 19670},
			expr: &actionExpr{
				pos: position{line: 440, col: 16, offset: 19685},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 440, col: 16, offset: 19685},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 440, col: 16, offset: 19685},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 440, col: 20, offset: 19689},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 440, col: 26, offset: 19695},
								expr: &choiceExpr{
									pos: position{line: 440, col: 27, offset: 19696},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 440, col: 27, offset: 19696},
											val:        "[^'\\n]",
											chars:      []rune{'\'', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 440, col: 36, offset: 19705},
											val:        "''",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 440, col: 43, offset: 19712},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 444, col: 1, offset: 19764},
			expr: &actionExpr{
				pos: position{line: 444, col: 10, offset: 19773},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 444, col: 10, offset: 19773},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 444, col: 10, offset: 19773},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 444, col: 18, offset: 19781},
							expr: &charClassMatcher{
								pos:        position{line: 444, col: 18, offset: 19781},
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
			pos:  position{line: 448, col: 1, offset: 19834},
			expr: &actionExpr{
				pos: position{line: 448, col: 11, offset: 19844},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 448, col: 13, offset: 19846},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 448, col: 13, offset: 19846},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 448, col: 19, offset: 19852},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 448, col: 19, offset: 19852},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 448, col: 24, offset: 19857},
									expr: &charClassMatcher{
										pos:        position{line: 448, col: 24, offset: 19857},
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
			pos:  position{line: 453, col: 1, offset: 19952},
			expr: &actionExpr{
				pos: position{line: 453, col: 15, offset: 19966},
				run: (*parser).callonNonZNumber1,
				expr: &seqExpr{
					pos: position{line: 453, col: 15, offset: 19966},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 453, col: 15, offset: 19966},
							val:        "[1-9]",
							ranges:     []rune{'1', '9'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 453, col: 20, offset: 19971},
							expr: &charClassMatcher{
								pos:        position{line: 453, col: 20, offset: 19971},
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
			pos:  position{line: 458, col: 1, offset: 20064},
			expr: &actionExpr{
				pos: position{line: 458, col: 12, offset: 20075},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 458, col: 12, offset: 20075},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 458, col: 20, offset: 20083},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 458, col: 20, offset: 20083},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 458, col: 34, offset: 20097},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 462, col: 1, offset: 20139},
			expr: &actionExpr{
				pos: position{line: 462, col: 16, offset: 20154},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 462, col: 18, offset: 20156},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 462, col: 18, offset: 20156},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 462, col: 27, offset: 20165},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 462, col: 27, offset: 20165},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 462, col: 31, offset: 20169},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 462, col: 49, offset: 20187},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 462, col: 55, offset: 20193},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 466, col: 1, offset: 20239},
			expr: &choiceExpr{
				pos: position{line: 466, col: 24, offset: 20262},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 466, col: 24, offset: 20262},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 466, col: 33, offset: 20271},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 466, col: 41, offset: 20279},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 466, col: 48, offset: 20286},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 466, col: 54, offset: 20292},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 468, col: 1, offset: 20299},
			expr: &actionExpr{
				pos: position{line: 468, col: 17, offset: 20315},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 468, col: 19, offset: 20317},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 468, col: 19, offset: 20317},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 468, col: 29, offset: 20327},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 468, col: 29, offset: 20327},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 468, col: 33, offset: 20331},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 468, col: 52, offset: 20350},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 468, col: 58, offset: 20356},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 472, col: 1, offset: 20404},
			expr: &choiceExpr{
				pos: position{line: 472, col: 25, offset: 20428},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 472, col: 25, offset: 20428},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 472, col: 35, offset: 20438},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 472, col: 42, offset: 20445},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 472, col: 50, offset: 20453},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 472, col: 56, offset: 20459},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 484, col: 1, offset: 20974},
			expr: &actionExpr{
				pos: position{line: 484, col: 12, offset: 20985},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 484, col: 14, offset: 20987},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 484, col: 14, offset: 20987},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 484, col: 34, offset: 21007},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 488, col: 1, offset: 21051},
			expr: &seqExpr{
				pos: position{line: 488, col: 21, offset: 21071},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 488, col: 21, offset: 21071},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 488, col: 26, offset: 21076},
						expr: &anyMatcher{
							line: 488, col: 26, offset: 21076,
						},
					},
					&litMatcher{
						pos:        position{line: 488, col: 29, offset: 21079},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 488, col: 34, offset: 21084},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 490, col: 1, offset: 21089},
			expr: &seqExpr{
				pos: position{line: 490, col: 22, offset: 21110},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 490, col: 22, offset: 21110},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 490, col: 27, offset: 21115},
						expr: &charClassMatcher{
							pos:        position{line: 490, col: 27, offset: 21115},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 490, col: 36, offset: 21124},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 492, col: 1, offset: 21129},
			expr: &seqExpr{
				pos: position{line: 492, col: 9, offset: 21137},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 492, col: 9, offset: 21137},
						expr: &charClassMatcher{
							pos:        position{line: 492, col: 9, offset: 21137},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 492, col: 17, offset: 21145},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 492, col: 17, offset: 21145},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 492, col: 26, offset: 21154},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 492, col: 35, offset: 21163},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 492, col: 42, offset: 21170},
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
			pos:         position{line: 494, col: 1, offset: 21177},
			expr: &zeroOrMoreExpr{
				pos: position{line: 494, col: 19, offset: 21195},
				expr: &charClassMatcher{
					pos:        position{line: 494, col: 19, offset: 21195},
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
			pos:         position{line: 496, col: 1, offset: 21207},
			expr: &oneOrMoreExpr{
				pos: position{line: 496, col: 31, offset: 21237},
				expr: &charClassMatcher{
					pos:        position{line: 496, col: 31, offset: 21237},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 498, col: 1, offset: 21249},
			expr: &notExpr{
				pos: position{line: 498, col: 8, offset: 21256},
				expr: &anyMatcher{
					line: 498, col: 9, offset: 21257,
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

func (c *current) onFieldDef1(name, dataType, constraint interface{}) (interface{}, error) {
	if dataType == nil {
		return nil, nil
	}
	result := dataType.(map[string]string)
	if constraint != nil {
		if err := mergo.Merge(&result, constraint.(map[string]string), mergo.WithOverride); err != nil {
			return nil, err
		}
	}
	result["name"] = interfaceToString(name)
	return result, nil
}

func (p *parser) callonFieldDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldDef1(stack["name"], stack["dataType"], stack["constraint"])
}

func (c *current) onColumnConstraint1(nameOpt, constraintClauses interface{}) (interface{}, error) {
	clausesSlice := toIfaceSlice(constraintClauses)
	properties := make(map[string]string)
	for _, clause := range clausesSlice {
		if clause == nil {
			continue
		}
		if err := mergo.Merge(&properties, clause.(map[string]string), mergo.WithOverride); err != nil {
			return nil, err
		}
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
	return p.cur.onColumnConstraint1(stack["nameOpt"], stack["constraintClauses"])
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
