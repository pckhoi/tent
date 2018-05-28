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
							label: "collation",
							expr: &zeroOrOneExpr{
								pos: position{line: 51, col: 90, offset: 2776},
								expr: &ruleRefExpr{
									pos:  position{line: 51, col: 90, offset: 2776},
									name: "Collate",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 51, col: 99, offset: 2785},
							label: "constraint",
							expr: &zeroOrOneExpr{
								pos: position{line: 51, col: 110, offset: 2796},
								expr: &ruleRefExpr{
									pos:  position{line: 51, col: 110, offset: 2796},
									name: "ColumnConstraint",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Collate",
			pos:  position{line: 71, col: 1, offset: 3366},
			expr: &actionExpr{
				pos: position{line: 71, col: 12, offset: 3377},
				run: (*parser).callonCollate1,
				expr: &seqExpr{
					pos: position{line: 71, col: 12, offset: 3377},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 71, col: 12, offset: 3377},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 71, col: 15, offset: 3380},
							val:        "collate",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 71, col: 26, offset: 3391},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 71, col: 29, offset: 3394},
							label: "collation",
							expr: &ruleRefExpr{
								pos:  position{line: 71, col: 39, offset: 3404},
								name: "Collation",
							},
						},
					},
				},
			},
		},
		{
			name: "Collation",
			pos:  position{line: 75, col: 1, offset: 3445},
			expr: &actionExpr{
				pos: position{line: 75, col: 14, offset: 3458},
				run: (*parser).callonCollation1,
				expr: &seqExpr{
					pos: position{line: 75, col: 14, offset: 3458},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 75, col: 14, offset: 3458},
							expr: &seqExpr{
								pos: position{line: 75, col: 16, offset: 3460},
								exprs: []interface{}{
									&choiceExpr{
										pos: position{line: 75, col: 18, offset: 3462},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 75, col: 18, offset: 3462},
												name: "DblQuotedString",
											},
											&ruleRefExpr{
												pos:  position{line: 75, col: 36, offset: 3480},
												name: "Ident",
											},
										},
									},
									&litMatcher{
										pos:        position{line: 75, col: 44, offset: 3488},
										val:        ".",
										ignoreCase: false,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 51, offset: 3495},
							name: "DblQuotedString",
						},
					},
				},
			},
		},
		{
			name: "ColumnConstraint",
			pos:  position{line: 79, col: 1, offset: 3547},
			expr: &actionExpr{
				pos: position{line: 79, col: 21, offset: 3567},
				run: (*parser).callonColumnConstraint1,
				expr: &seqExpr{
					pos: position{line: 79, col: 21, offset: 3567},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 79, col: 21, offset: 3567},
							label: "nameOpt",
							expr: &zeroOrOneExpr{
								pos: position{line: 79, col: 29, offset: 3575},
								expr: &seqExpr{
									pos: position{line: 79, col: 31, offset: 3577},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 79, col: 31, offset: 3577},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 79, col: 34, offset: 3580},
											val:        "constraint",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 79, col: 48, offset: 3594},
											name: "_1",
										},
										&choiceExpr{
											pos: position{line: 79, col: 52, offset: 3598},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 79, col: 52, offset: 3598},
													name: "StringConst",
												},
												&ruleRefExpr{
													pos:  position{line: 79, col: 66, offset: 3612},
													name: "Ident",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 79, col: 76, offset: 3622},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 79, col: 78, offset: 3624},
							label: "constraint",
							expr: &choiceExpr{
								pos: position{line: 79, col: 91, offset: 3637},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 79, col: 91, offset: 3637},
										name: "NotNullCls",
									},
									&ruleRefExpr{
										pos:  position{line: 79, col: 104, offset: 3650},
										name: "NullCls",
									},
									&ruleRefExpr{
										pos:  position{line: 79, col: 114, offset: 3660},
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
			pos:  position{line: 94, col: 1, offset: 4026},
			expr: &actionExpr{
				pos: position{line: 94, col: 16, offset: 4041},
				run: (*parser).callonTableConstr1,
				expr: &seqExpr{
					pos: position{line: 94, col: 16, offset: 4041},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 94, col: 16, offset: 4041},
							label: "nameOpt",
							expr: &zeroOrOneExpr{
								pos: position{line: 94, col: 24, offset: 4049},
								expr: &seqExpr{
									pos: position{line: 94, col: 26, offset: 4051},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 94, col: 26, offset: 4051},
											val:        "constraint",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 94, col: 40, offset: 4065},
											name: "_1",
										},
										&choiceExpr{
											pos: position{line: 94, col: 44, offset: 4069},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 94, col: 44, offset: 4069},
													name: "StringConst",
												},
												&ruleRefExpr{
													pos:  position{line: 94, col: 58, offset: 4083},
													name: "Ident",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 94, col: 68, offset: 4093},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 94, col: 70, offset: 4095},
							label: "constraint",
							expr: &ruleRefExpr{
								pos:  position{line: 94, col: 81, offset: 4106},
								name: "CheckCls",
							},
						},
					},
				},
			},
		},
		{
			name: "NotNullCls",
			pos:  position{line: 111, col: 1, offset: 4507},
			expr: &actionExpr{
				pos: position{line: 111, col: 15, offset: 4521},
				run: (*parser).callonNotNullCls1,
				expr: &seqExpr{
					pos: position{line: 111, col: 15, offset: 4521},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 111, col: 15, offset: 4521},
							val:        "not",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 111, col: 22, offset: 4528},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 111, col: 25, offset: 4531},
							val:        "null",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "NullCls",
			pos:  position{line: 117, col: 1, offset: 4613},
			expr: &actionExpr{
				pos: position{line: 117, col: 12, offset: 4624},
				run: (*parser).callonNullCls1,
				expr: &litMatcher{
					pos:        position{line: 117, col: 12, offset: 4624},
					val:        "null",
					ignoreCase: true,
				},
			},
		},
		{
			name: "CheckCls",
			pos:  position{line: 123, col: 1, offset: 4707},
			expr: &actionExpr{
				pos: position{line: 123, col: 13, offset: 4719},
				run: (*parser).callonCheckCls1,
				expr: &seqExpr{
					pos: position{line: 123, col: 13, offset: 4719},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 123, col: 13, offset: 4719},
							val:        "check",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 123, col: 22, offset: 4728},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 123, col: 25, offset: 4731},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 123, col: 30, offset: 4736},
								name: "WrappedExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 123, col: 42, offset: 4748},
							label: "noInherit",
							expr: &zeroOrOneExpr{
								pos: position{line: 123, col: 52, offset: 4758},
								expr: &seqExpr{
									pos: position{line: 123, col: 54, offset: 4760},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 123, col: 54, offset: 4760},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 123, col: 57, offset: 4763},
											val:        "no",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 123, col: 63, offset: 4769},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 123, col: 66, offset: 4772},
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
			pos:  position{line: 133, col: 1, offset: 4965},
			expr: &actionExpr{
				pos: position{line: 133, col: 16, offset: 4980},
				run: (*parser).callonWrappedExpr1,
				expr: &seqExpr{
					pos: position{line: 133, col: 16, offset: 4980},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 133, col: 16, offset: 4980},
							val:        "(",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 133, col: 20, offset: 4984},
							expr: &ruleRefExpr{
								pos:  position{line: 133, col: 20, offset: 4984},
								name: "Expr",
							},
						},
						&litMatcher{
							pos:        position{line: 133, col: 26, offset: 4990},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 137, col: 1, offset: 5030},
			expr: &choiceExpr{
				pos: position{line: 137, col: 9, offset: 5038},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 137, col: 9, offset: 5038},
						name: "WrappedExpr",
					},
					&oneOrMoreExpr{
						pos: position{line: 137, col: 23, offset: 5052},
						expr: &charClassMatcher{
							pos:        position{line: 137, col: 23, offset: 5052},
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
			pos:  position{line: 151, col: 1, offset: 6074},
			expr: &actionExpr{
				pos: position{line: 151, col: 13, offset: 6086},
				run: (*parser).callonDataType1,
				expr: &seqExpr{
					pos: position{line: 151, col: 13, offset: 6086},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 151, col: 13, offset: 6086},
							label: "t",
							expr: &choiceExpr{
								pos: position{line: 151, col: 17, offset: 6090},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 151, col: 17, offset: 6090},
										name: "TimestampT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 30, offset: 6103},
										name: "TimeT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 38, offset: 6111},
										name: "NumericT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 49, offset: 6122},
										name: "VarcharT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 60, offset: 6133},
										name: "CharT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 68, offset: 6141},
										name: "BitVarT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 78, offset: 6151},
										name: "BitT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 85, offset: 6158},
										name: "IntT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 92, offset: 6165},
										name: "PgOidT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 101, offset: 6174},
										name: "PostgisT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 112, offset: 6185},
										name: "OtherT",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 121, offset: 6194},
										name: "CustomT",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 151, col: 131, offset: 6204},
							label: "brackets",
							expr: &zeroOrMoreExpr{
								pos: position{line: 151, col: 140, offset: 6213},
								expr: &litMatcher{
									pos:        position{line: 151, col: 142, offset: 6215},
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
			pos:  position{line: 164, col: 1, offset: 6526},
			expr: &actionExpr{
				pos: position{line: 164, col: 15, offset: 6540},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 164, col: 15, offset: 6540},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 164, col: 15, offset: 6540},
							val:        "timestamp",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 164, col: 28, offset: 6553},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 164, col: 33, offset: 6558},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 164, col: 46, offset: 6571},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 164, col: 59, offset: 6584},
								expr: &choiceExpr{
									pos: position{line: 164, col: 61, offset: 6586},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 164, col: 61, offset: 6586},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 164, col: 70, offset: 6595},
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
			pos:  position{line: 177, col: 1, offset: 6874},
			expr: &actionExpr{
				pos: position{line: 177, col: 10, offset: 6883},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 177, col: 10, offset: 6883},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 177, col: 10, offset: 6883},
							val:        "time",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 177, col: 18, offset: 6891},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 177, col: 23, offset: 6896},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 177, col: 36, offset: 6909},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 177, col: 49, offset: 6922},
								expr: &choiceExpr{
									pos: position{line: 177, col: 51, offset: 6924},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 177, col: 51, offset: 6924},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 177, col: 60, offset: 6933},
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
			pos:  position{line: 190, col: 1, offset: 7204},
			expr: &actionExpr{
				pos: position{line: 190, col: 17, offset: 7220},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 190, col: 17, offset: 7220},
					expr: &seqExpr{
						pos: position{line: 190, col: 19, offset: 7222},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 190, col: 19, offset: 7222},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 190, col: 22, offset: 7225},
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
			pos:  position{line: 197, col: 1, offset: 7353},
			expr: &actionExpr{
				pos: position{line: 197, col: 11, offset: 7363},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 197, col: 11, offset: 7363},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 197, col: 11, offset: 7363},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 197, col: 14, offset: 7366},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 197, col: 22, offset: 7374},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 197, col: 25, offset: 7377},
							val:        "time",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 197, col: 33, offset: 7385},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 197, col: 36, offset: 7388},
							val:        "zone",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 201, col: 1, offset: 7422},
			expr: &actionExpr{
				pos: position{line: 201, col: 14, offset: 7435},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 201, col: 14, offset: 7435},
					expr: &seqExpr{
						pos: position{line: 201, col: 16, offset: 7437},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 201, col: 16, offset: 7437},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 201, col: 19, offset: 7440},
								val:        "without",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 201, col: 30, offset: 7451},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 201, col: 33, offset: 7454},
								val:        "time",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 201, col: 41, offset: 7462},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 201, col: 44, offset: 7465},
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
			pos:  position{line: 205, col: 1, offset: 7503},
			expr: &actionExpr{
				pos: position{line: 205, col: 10, offset: 7512},
				run: (*parser).callonCharT1,
				expr: &seqExpr{
					pos: position{line: 205, col: 10, offset: 7512},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 205, col: 12, offset: 7514},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 205, col: 12, offset: 7514},
									val:        "character",
									ignoreCase: true,
								},
								&litMatcher{
									pos:        position{line: 205, col: 27, offset: 7529},
									val:        "char",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 205, col: 37, offset: 7539},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 205, col: 44, offset: 7546},
								expr: &seqExpr{
									pos: position{line: 205, col: 46, offset: 7548},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 205, col: 46, offset: 7548},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 205, col: 50, offset: 7552},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 205, col: 61, offset: 7563},
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
			pos:  position{line: 217, col: 1, offset: 7818},
			expr: &actionExpr{
				pos: position{line: 217, col: 13, offset: 7830},
				run: (*parser).callonVarcharT1,
				expr: &seqExpr{
					pos: position{line: 217, col: 13, offset: 7830},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 217, col: 15, offset: 7832},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 217, col: 17, offset: 7834},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 217, col: 17, offset: 7834},
											val:        "character",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 217, col: 30, offset: 7847},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 217, col: 33, offset: 7850},
											val:        "varying",
											ignoreCase: true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 217, col: 48, offset: 7865},
									val:        "varchar",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 217, col: 61, offset: 7878},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 217, col: 68, offset: 7885},
								expr: &seqExpr{
									pos: position{line: 217, col: 70, offset: 7887},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 217, col: 70, offset: 7887},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 217, col: 74, offset: 7891},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 217, col: 85, offset: 7902},
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
			pos:  position{line: 228, col: 1, offset: 8137},
			expr: &actionExpr{
				pos: position{line: 228, col: 9, offset: 8145},
				run: (*parser).callonBitT1,
				expr: &seqExpr{
					pos: position{line: 228, col: 9, offset: 8145},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 228, col: 9, offset: 8145},
							val:        "bit",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 228, col: 16, offset: 8152},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 228, col: 23, offset: 8159},
								expr: &seqExpr{
									pos: position{line: 228, col: 25, offset: 8161},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 228, col: 25, offset: 8161},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 228, col: 29, offset: 8165},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 228, col: 40, offset: 8176},
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
			pos:  position{line: 240, col: 1, offset: 8430},
			expr: &actionExpr{
				pos: position{line: 240, col: 12, offset: 8441},
				run: (*parser).callonBitVarT1,
				expr: &seqExpr{
					pos: position{line: 240, col: 12, offset: 8441},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 240, col: 12, offset: 8441},
							val:        "bit",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 19, offset: 8448},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 240, col: 22, offset: 8451},
							val:        "varying",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 240, col: 33, offset: 8462},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 240, col: 40, offset: 8469},
								expr: &seqExpr{
									pos: position{line: 240, col: 42, offset: 8471},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 240, col: 42, offset: 8471},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 240, col: 46, offset: 8475},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 240, col: 57, offset: 8486},
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
			pos:  position{line: 251, col: 1, offset: 8720},
			expr: &actionExpr{
				pos: position{line: 251, col: 9, offset: 8728},
				run: (*parser).callonIntT1,
				expr: &choiceExpr{
					pos: position{line: 251, col: 11, offset: 8730},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 251, col: 11, offset: 8730},
							val:        "integer",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 251, col: 24, offset: 8743},
							val:        "int",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "NumericT",
			pos:  position{line: 257, col: 1, offset: 8825},
			expr: &actionExpr{
				pos: position{line: 257, col: 13, offset: 8837},
				run: (*parser).callonNumericT1,
				expr: &seqExpr{
					pos: position{line: 257, col: 13, offset: 8837},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 257, col: 13, offset: 8837},
							val:        "numeric",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 257, col: 24, offset: 8848},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 257, col: 29, offset: 8853},
								expr: &seqExpr{
									pos: position{line: 257, col: 31, offset: 8855},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 257, col: 31, offset: 8855},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 257, col: 35, offset: 8859},
											name: "NonZNumber",
										},
										&zeroOrOneExpr{
											pos: position{line: 257, col: 46, offset: 8870},
											expr: &seqExpr{
												pos: position{line: 257, col: 48, offset: 8872},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 257, col: 48, offset: 8872},
														val:        ",",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 257, col: 52, offset: 8876},
														name: "_",
													},
													&ruleRefExpr{
														pos:  position{line: 257, col: 54, offset: 8878},
														name: "NonZNumber",
													},
												},
											},
										},
										&litMatcher{
											pos:        position{line: 257, col: 68, offset: 8892},
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
			pos:  position{line: 272, col: 1, offset: 9295},
			expr: &actionExpr{
				pos: position{line: 272, col: 13, offset: 9307},
				run: (*parser).callonPostgisT1,
				expr: &seqExpr{
					pos: position{line: 272, col: 13, offset: 9307},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 272, col: 13, offset: 9307},
							label: "t",
							expr: &choiceExpr{
								pos: position{line: 272, col: 17, offset: 9311},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 272, col: 17, offset: 9311},
										val:        "geography",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 272, col: 32, offset: 9326},
										val:        "geometry",
										ignoreCase: true,
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 272, col: 46, offset: 9340},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 272, col: 50, offset: 9344},
							label: "subtype",
							expr: &choiceExpr{
								pos: position{line: 272, col: 60, offset: 9354},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 272, col: 60, offset: 9354},
										val:        "point",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 272, col: 71, offset: 9365},
										val:        "linestring",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 272, col: 87, offset: 9381},
										val:        "polygon",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 272, col: 100, offset: 9394},
										val:        "multipoint",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 272, col: 116, offset: 9410},
										val:        "multilinestring",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 272, col: 137, offset: 9431},
										val:        "multipolygon",
										ignoreCase: true,
									},
									&litMatcher{
										pos:        position{line: 272, col: 155, offset: 9449},
										val:        "geometrycollection",
										ignoreCase: true,
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 272, col: 179, offset: 9473},
							label: "srid",
							expr: &zeroOrOneExpr{
								pos: position{line: 272, col: 184, offset: 9478},
								expr: &seqExpr{
									pos: position{line: 272, col: 185, offset: 9479},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 272, col: 185, offset: 9479},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 272, col: 189, offset: 9483},
											name: "NonZNumber",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 272, col: 202, offset: 9496},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PgOidT",
			pos:  position{line: 284, col: 1, offset: 9818},
			expr: &actionExpr{
				pos: position{line: 284, col: 11, offset: 9828},
				run: (*parser).callonPgOidT1,
				expr: &choiceExpr{
					pos: position{line: 284, col: 13, offset: 9830},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 284, col: 13, offset: 9830},
							val:        "oid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 22, offset: 9839},
							val:        "regprocedure",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 40, offset: 9857},
							val:        "regproc",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 53, offset: 9870},
							val:        "regoperator",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 70, offset: 9887},
							val:        "regoper",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 83, offset: 9900},
							val:        "regclass",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 97, offset: 9914},
							val:        "regtype",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 110, offset: 9927},
							val:        "regrole",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 123, offset: 9940},
							val:        "regnamespace",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 141, offset: 9958},
							val:        "regconfig",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 284, col: 156, offset: 9973},
							val:        "regdictionary",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 290, col: 1, offset: 10087},
			expr: &actionExpr{
				pos: position{line: 290, col: 11, offset: 10097},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 290, col: 13, offset: 10099},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 290, col: 13, offset: 10099},
							val:        "date",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 23, offset: 10109},
							val:        "smallint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 37, offset: 10123},
							val:        "bigint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 49, offset: 10135},
							val:        "decimal",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 62, offset: 10148},
							val:        "real",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 72, offset: 10158},
							val:        "smallserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 89, offset: 10175},
							val:        "serial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 101, offset: 10187},
							val:        "bigserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 116, offset: 10202},
							val:        "boolean",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 129, offset: 10215},
							val:        "text",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 139, offset: 10225},
							val:        "money",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 150, offset: 10236},
							val:        "bytea",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 161, offset: 10247},
							val:        "point",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 172, offset: 10258},
							val:        "line",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 182, offset: 10268},
							val:        "lseg",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 192, offset: 10278},
							val:        "box",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 201, offset: 10287},
							val:        "path",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 211, offset: 10297},
							val:        "polygon",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 224, offset: 10310},
							val:        "circle",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 236, offset: 10322},
							val:        "cidr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 246, offset: 10332},
							val:        "inet",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 256, offset: 10342},
							val:        "macaddr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 269, offset: 10355},
							val:        "uuid",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 279, offset: 10365},
							val:        "xml",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 288, offset: 10374},
							val:        "jsonb",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 290, col: 299, offset: 10385},
							val:        "json",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CustomT",
			pos:  position{line: 296, col: 1, offset: 10490},
			expr: &actionExpr{
				pos: position{line: 296, col: 13, offset: 10502},
				run: (*parser).callonCustomT1,
				expr: &ruleRefExpr{
					pos:  position{line: 296, col: 13, offset: 10502},
					name: "Ident",
				},
			},
		},
		{
			name: "CreateSeqStmt",
			pos:  position{line: 317, col: 1, offset: 11967},
			expr: &actionExpr{
				pos: position{line: 317, col: 18, offset: 11984},
				run: (*parser).callonCreateSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 317, col: 18, offset: 11984},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 317, col: 18, offset: 11984},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 28, offset: 11994},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 317, col: 31, offset: 11997},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 43, offset: 12009},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 317, col: 46, offset: 12012},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 317, col: 51, offset: 12017},
								name: "Ident",
							},
						},
						&labeledExpr{
							pos:   position{line: 317, col: 57, offset: 12023},
							label: "verses",
							expr: &zeroOrMoreExpr{
								pos: position{line: 317, col: 64, offset: 12030},
								expr: &ruleRefExpr{
									pos:  position{line: 317, col: 64, offset: 12030},
									name: "CreateSeqVerse",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 80, offset: 12046},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 317, col: 82, offset: 12048},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 317, col: 86, offset: 12052},
							expr: &ruleRefExpr{
								pos:  position{line: 317, col: 86, offset: 12052},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateSeqVerse",
			pos:  position{line: 331, col: 1, offset: 12446},
			expr: &actionExpr{
				pos: position{line: 331, col: 19, offset: 12464},
				run: (*parser).callonCreateSeqVerse1,
				expr: &labeledExpr{
					pos:   position{line: 331, col: 19, offset: 12464},
					label: "verse",
					expr: &choiceExpr{
						pos: position{line: 331, col: 27, offset: 12472},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 331, col: 27, offset: 12472},
								name: "IncrementBy",
							},
							&ruleRefExpr{
								pos:  position{line: 331, col: 41, offset: 12486},
								name: "MinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 331, col: 52, offset: 12497},
								name: "NoMinValue",
							},
							&ruleRefExpr{
								pos:  position{line: 331, col: 65, offset: 12510},
								name: "MaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 331, col: 76, offset: 12521},
								name: "NoMaxValue",
							},
							&ruleRefExpr{
								pos:  position{line: 331, col: 89, offset: 12534},
								name: "Start",
							},
							&ruleRefExpr{
								pos:  position{line: 331, col: 97, offset: 12542},
								name: "Cache",
							},
							&ruleRefExpr{
								pos:  position{line: 331, col: 105, offset: 12550},
								name: "Cycle",
							},
							&ruleRefExpr{
								pos:  position{line: 331, col: 113, offset: 12558},
								name: "OwnedBy",
							},
						},
					},
				},
			},
		},
		{
			name: "IncrementBy",
			pos:  position{line: 335, col: 1, offset: 12595},
			expr: &actionExpr{
				pos: position{line: 335, col: 16, offset: 12610},
				run: (*parser).callonIncrementBy1,
				expr: &seqExpr{
					pos: position{line: 335, col: 16, offset: 12610},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 335, col: 16, offset: 12610},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 335, col: 19, offset: 12613},
							val:        "increment",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 335, col: 32, offset: 12626},
							expr: &seqExpr{
								pos: position{line: 335, col: 33, offset: 12627},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 335, col: 33, offset: 12627},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 335, col: 36, offset: 12630},
										val:        "by",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 335, col: 44, offset: 12638},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 335, col: 47, offset: 12641},
							label: "num",
							expr: &ruleRefExpr{
								pos:  position{line: 335, col: 51, offset: 12645},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "MinValue",
			pos:  position{line: 341, col: 1, offset: 12759},
			expr: &actionExpr{
				pos: position{line: 341, col: 13, offset: 12771},
				run: (*parser).callonMinValue1,
				expr: &seqExpr{
					pos: position{line: 341, col: 13, offset: 12771},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 341, col: 13, offset: 12771},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 341, col: 16, offset: 12774},
							val:        "minvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 341, col: 28, offset: 12786},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 341, col: 31, offset: 12789},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 341, col: 35, offset: 12793},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMinValue",
			pos:  position{line: 347, col: 1, offset: 12906},
			expr: &actionExpr{
				pos: position{line: 347, col: 15, offset: 12920},
				run: (*parser).callonNoMinValue1,
				expr: &seqExpr{
					pos: position{line: 347, col: 15, offset: 12920},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 347, col: 15, offset: 12920},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 347, col: 18, offset: 12923},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 347, col: 24, offset: 12929},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 347, col: 27, offset: 12932},
							val:        "minvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "MaxValue",
			pos:  position{line: 351, col: 1, offset: 12969},
			expr: &actionExpr{
				pos: position{line: 351, col: 13, offset: 12981},
				run: (*parser).callonMaxValue1,
				expr: &seqExpr{
					pos: position{line: 351, col: 13, offset: 12981},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 351, col: 13, offset: 12981},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 351, col: 16, offset: 12984},
							val:        "maxvalue",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 351, col: 28, offset: 12996},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 351, col: 31, offset: 12999},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 351, col: 35, offset: 13003},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "NoMaxValue",
			pos:  position{line: 357, col: 1, offset: 13116},
			expr: &actionExpr{
				pos: position{line: 357, col: 15, offset: 13130},
				run: (*parser).callonNoMaxValue1,
				expr: &seqExpr{
					pos: position{line: 357, col: 15, offset: 13130},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 357, col: 15, offset: 13130},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 357, col: 18, offset: 13133},
							val:        "no",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 357, col: 24, offset: 13139},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 357, col: 27, offset: 13142},
							val:        "maxvalue",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "Start",
			pos:  position{line: 361, col: 1, offset: 13179},
			expr: &actionExpr{
				pos: position{line: 361, col: 10, offset: 13188},
				run: (*parser).callonStart1,
				expr: &seqExpr{
					pos: position{line: 361, col: 10, offset: 13188},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 361, col: 10, offset: 13188},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 361, col: 13, offset: 13191},
							val:        "start",
							ignoreCase: true,
						},
						&zeroOrOneExpr{
							pos: position{line: 361, col: 22, offset: 13200},
							expr: &seqExpr{
								pos: position{line: 361, col: 23, offset: 13201},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 361, col: 23, offset: 13201},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 361, col: 26, offset: 13204},
										val:        "with",
										ignoreCase: true,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 361, col: 36, offset: 13214},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 361, col: 39, offset: 13217},
							label: "start",
							expr: &ruleRefExpr{
								pos:  position{line: 361, col: 45, offset: 13223},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cache",
			pos:  position{line: 367, col: 1, offset: 13335},
			expr: &actionExpr{
				pos: position{line: 367, col: 10, offset: 13344},
				run: (*parser).callonCache1,
				expr: &seqExpr{
					pos: position{line: 367, col: 10, offset: 13344},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 367, col: 10, offset: 13344},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 367, col: 13, offset: 13347},
							val:        "cache",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 367, col: 22, offset: 13356},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 367, col: 25, offset: 13359},
							label: "cache",
							expr: &ruleRefExpr{
								pos:  position{line: 367, col: 31, offset: 13365},
								name: "NonZNumber",
							},
						},
					},
				},
			},
		},
		{
			name: "Cycle",
			pos:  position{line: 373, col: 1, offset: 13477},
			expr: &actionExpr{
				pos: position{line: 373, col: 10, offset: 13486},
				run: (*parser).callonCycle1,
				expr: &seqExpr{
					pos: position{line: 373, col: 10, offset: 13486},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 373, col: 10, offset: 13486},
							label: "no",
							expr: &zeroOrOneExpr{
								pos: position{line: 373, col: 13, offset: 13489},
								expr: &seqExpr{
									pos: position{line: 373, col: 14, offset: 13490},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 373, col: 14, offset: 13490},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 373, col: 17, offset: 13493},
											val:        "no",
											ignoreCase: true,
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 373, col: 25, offset: 13501},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 373, col: 28, offset: 13504},
							val:        "cycle",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OwnedBy",
			pos:  position{line: 384, col: 1, offset: 13688},
			expr: &actionExpr{
				pos: position{line: 384, col: 12, offset: 13699},
				run: (*parser).callonOwnedBy1,
				expr: &seqExpr{
					pos: position{line: 384, col: 12, offset: 13699},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 384, col: 12, offset: 13699},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 384, col: 15, offset: 13702},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 384, col: 24, offset: 13711},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 384, col: 27, offset: 13714},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 384, col: 33, offset: 13720},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 384, col: 36, offset: 13723},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 384, col: 43, offset: 13730},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 384, col: 43, offset: 13730},
										val:        "none",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 384, col: 53, offset: 13740},
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
			pos:  position{line: 402, col: 1, offset: 15199},
			expr: &actionExpr{
				pos: position{line: 402, col: 19, offset: 15217},
				run: (*parser).callonCreateTypeStmt1,
				expr: &seqExpr{
					pos: position{line: 402, col: 19, offset: 15217},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 402, col: 19, offset: 15217},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 29, offset: 15227},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 402, col: 32, offset: 15230},
							val:        "type",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 40, offset: 15238},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 402, col: 43, offset: 15241},
							label: "typename",
							expr: &ruleRefExpr{
								pos:  position{line: 402, col: 52, offset: 15250},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 58, offset: 15256},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 402, col: 61, offset: 15259},
							val:        "as",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 67, offset: 15265},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 402, col: 70, offset: 15268},
							label: "typedef",
							expr: &ruleRefExpr{
								pos:  position{line: 402, col: 78, offset: 15276},
								name: "EnumDef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 402, col: 86, offset: 15284},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 402, col: 88, offset: 15286},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 402, col: 92, offset: 15290},
							expr: &ruleRefExpr{
								pos:  position{line: 402, col: 92, offset: 15290},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "EnumDef",
			pos:  position{line: 408, col: 1, offset: 15406},
			expr: &actionExpr{
				pos: position{line: 408, col: 12, offset: 15417},
				run: (*parser).callonEnumDef1,
				expr: &seqExpr{
					pos: position{line: 408, col: 12, offset: 15417},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 408, col: 12, offset: 15417},
							val:        "ENUM",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 408, col: 19, offset: 15424},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 408, col: 21, offset: 15426},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 408, col: 25, offset: 15430},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 408, col: 27, offset: 15432},
							label: "vals",
							expr: &seqExpr{
								pos: position{line: 408, col: 34, offset: 15439},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 408, col: 34, offset: 15439},
										name: "StringConst",
									},
									&zeroOrMoreExpr{
										pos: position{line: 408, col: 46, offset: 15451},
										expr: &seqExpr{
											pos: position{line: 408, col: 48, offset: 15453},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 408, col: 48, offset: 15453},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 408, col: 50, offset: 15455},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 408, col: 54, offset: 15459},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 408, col: 56, offset: 15461},
													name: "StringConst",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 408, col: 74, offset: 15479},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 408, col: 76, offset: 15481},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AlterTableStmt",
			pos:  position{line: 433, col: 1, offset: 17111},
			expr: &actionExpr{
				pos: position{line: 433, col: 19, offset: 17129},
				run: (*parser).callonAlterTableStmt1,
				expr: &seqExpr{
					pos: position{line: 433, col: 19, offset: 17129},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 433, col: 19, offset: 17129},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 28, offset: 17138},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 433, col: 31, offset: 17141},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 40, offset: 17150},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 433, col: 43, offset: 17153},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 433, col: 48, offset: 17158},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 54, offset: 17164},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 433, col: 57, offset: 17167},
							val:        "owner",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 66, offset: 17176},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 433, col: 69, offset: 17179},
							val:        "to",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 75, offset: 17185},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 433, col: 78, offset: 17188},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 433, col: 84, offset: 17194},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 433, col: 90, offset: 17200},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 433, col: 92, offset: 17202},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 433, col: 96, offset: 17206},
							expr: &ruleRefExpr{
								pos:  position{line: 433, col: 96, offset: 17206},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "AlterSeqStmt",
			pos:  position{line: 447, col: 1, offset: 18352},
			expr: &actionExpr{
				pos: position{line: 447, col: 17, offset: 18368},
				run: (*parser).callonAlterSeqStmt1,
				expr: &seqExpr{
					pos: position{line: 447, col: 17, offset: 18368},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 447, col: 17, offset: 18368},
							val:        "alter",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 447, col: 26, offset: 18377},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 447, col: 29, offset: 18380},
							val:        "sequence",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 447, col: 41, offset: 18392},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 447, col: 44, offset: 18395},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 447, col: 49, offset: 18400},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 447, col: 55, offset: 18406},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 447, col: 58, offset: 18409},
							val:        "owned",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 447, col: 67, offset: 18418},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 447, col: 70, offset: 18421},
							val:        "by",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 447, col: 76, offset: 18427},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 447, col: 79, offset: 18430},
							label: "owner",
							expr: &ruleRefExpr{
								pos:  position{line: 447, col: 85, offset: 18436},
								name: "TableDotCol",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 447, col: 97, offset: 18448},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 447, col: 99, offset: 18450},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 447, col: 103, offset: 18454},
							expr: &ruleRefExpr{
								pos:  position{line: 447, col: 103, offset: 18454},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "TableDotCol",
			pos:  position{line: 451, col: 1, offset: 18533},
			expr: &actionExpr{
				pos: position{line: 451, col: 16, offset: 18548},
				run: (*parser).callonTableDotCol1,
				expr: &seqExpr{
					pos: position{line: 451, col: 16, offset: 18548},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 451, col: 16, offset: 18548},
							label: "table",
							expr: &ruleRefExpr{
								pos:  position{line: 451, col: 22, offset: 18554},
								name: "Ident",
							},
						},
						&litMatcher{
							pos:        position{line: 451, col: 28, offset: 18560},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 451, col: 32, offset: 18564},
							label: "column",
							expr: &ruleRefExpr{
								pos:  position{line: 451, col: 39, offset: 18571},
								name: "Ident",
							},
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 465, col: 1, offset: 19899},
			expr: &actionExpr{
				pos: position{line: 465, col: 25, offset: 19923},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 465, col: 25, offset: 19923},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 465, col: 25, offset: 19923},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 465, col: 36, offset: 19934},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 465, col: 39, offset: 19937},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 465, col: 45, offset: 19943},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 465, col: 48, offset: 19946},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 465, col: 61, offset: 19959},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 465, col: 63, offset: 19961},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 465, col: 73, offset: 19971},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 465, col: 79, offset: 19977},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 465, col: 81, offset: 19979},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 465, col: 87, offset: 19985},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 465, col: 89, offset: 19987},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 465, col: 97, offset: 19995},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 465, col: 109, offset: 20007},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 465, col: 111, offset: 20009},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 465, col: 115, offset: 20013},
							expr: &ruleRefExpr{
								pos:  position{line: 465, col: 115, offset: 20013},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 469, col: 1, offset: 20102},
			expr: &actionExpr{
				pos: position{line: 469, col: 24, offset: 20125},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 469, col: 24, offset: 20125},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 469, col: 24, offset: 20125},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 469, col: 34, offset: 20135},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 469, col: 37, offset: 20138},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 469, col: 50, offset: 20151},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 469, col: 53, offset: 20154},
							expr: &seqExpr{
								pos: position{line: 469, col: 55, offset: 20156},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 469, col: 55, offset: 20156},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 469, col: 61, offset: 20162},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 469, col: 64, offset: 20165},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 469, col: 71, offset: 20172},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 469, col: 74, offset: 20175},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 469, col: 84, offset: 20185},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 469, col: 90, offset: 20191},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 469, col: 100, offset: 20201},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 469, col: 106, offset: 20207},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 469, col: 109, offset: 20210},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 469, col: 117, offset: 20218},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 469, col: 120, offset: 20221},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 469, col: 130, offset: 20231},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 469, col: 133, offset: 20234},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 469, col: 140, offset: 20241},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 469, col: 146, offset: 20247},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 469, col: 148, offset: 20249},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 469, col: 152, offset: 20253},
							expr: &ruleRefExpr{
								pos:  position{line: 469, col: 152, offset: 20253},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 473, col: 1, offset: 20344},
			expr: &actionExpr{
				pos: position{line: 473, col: 12, offset: 20355},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 473, col: 12, offset: 20355},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 473, col: 12, offset: 20355},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 473, col: 19, offset: 20362},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 473, col: 21, offset: 20364},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 473, col: 25, offset: 20368},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 473, col: 29, offset: 20372},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 473, col: 33, offset: 20376},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 473, col: 33, offset: 20376},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 473, col: 39, offset: 20382},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 473, col: 47, offset: 20390},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 473, col: 49, offset: 20392},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 473, col: 56, offset: 20399},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 473, col: 77, offset: 20420},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 473, col: 79, offset: 20422},
							val:        ";",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 473, col: 83, offset: 20426},
							expr: &ruleRefExpr{
								pos:  position{line: 473, col: 83, offset: 20426},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 478, col: 1, offset: 20508},
			expr: &actionExpr{
				pos: position{line: 478, col: 8, offset: 20515},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 478, col: 8, offset: 20515},
					expr: &charClassMatcher{
						pos:        position{line: 478, col: 8, offset: 20515},
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
			pos:  position{line: 493, col: 1, offset: 21355},
			expr: &actionExpr{
				pos: position{line: 493, col: 25, offset: 21379},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 493, col: 25, offset: 21379},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 493, col: 32, offset: 21386},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 493, col: 32, offset: 21386},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 493, col: 38, offset: 21392},
								expr: &seqExpr{
									pos: position{line: 493, col: 40, offset: 21394},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 493, col: 40, offset: 21394},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 493, col: 42, offset: 21396},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 493, col: 46, offset: 21400},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 493, col: 48, offset: 21402},
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
			pos:  position{line: 505, col: 1, offset: 21692},
			expr: &choiceExpr{
				pos: position{line: 505, col: 12, offset: 21703},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 505, col: 12, offset: 21703},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 505, col: 21, offset: 21712},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 505, col: 31, offset: 21722},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 505, col: 45, offset: 21736},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 507, col: 1, offset: 21745},
			expr: &actionExpr{
				pos: position{line: 507, col: 16, offset: 21760},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 507, col: 16, offset: 21760},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 507, col: 16, offset: 21760},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 507, col: 20, offset: 21764},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 507, col: 26, offset: 21770},
								expr: &choiceExpr{
									pos: position{line: 507, col: 27, offset: 21771},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 507, col: 27, offset: 21771},
											val:        "[^'\\n]",
											chars:      []rune{'\'', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 507, col: 36, offset: 21780},
											val:        "''",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 507, col: 43, offset: 21787},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "DblQuotedString",
			pos:  position{line: 511, col: 1, offset: 21839},
			expr: &actionExpr{
				pos: position{line: 511, col: 20, offset: 21858},
				run: (*parser).callonDblQuotedString1,
				expr: &seqExpr{
					pos: position{line: 511, col: 20, offset: 21858},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 511, col: 20, offset: 21858},
							val:        "\"",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 511, col: 24, offset: 21862},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 511, col: 30, offset: 21868},
								expr: &choiceExpr{
									pos: position{line: 511, col: 31, offset: 21869},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 511, col: 31, offset: 21869},
											val:        "[^\"\\n]",
											chars:      []rune{'"', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 511, col: 40, offset: 21878},
											val:        "\"\"",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 511, col: 49, offset: 21887},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 515, col: 1, offset: 21951},
			expr: &actionExpr{
				pos: position{line: 515, col: 10, offset: 21960},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 515, col: 10, offset: 21960},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 515, col: 10, offset: 21960},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 515, col: 18, offset: 21968},
							expr: &charClassMatcher{
								pos:        position{line: 515, col: 18, offset: 21968},
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
			pos:  position{line: 519, col: 1, offset: 22021},
			expr: &actionExpr{
				pos: position{line: 519, col: 11, offset: 22031},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 519, col: 13, offset: 22033},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 519, col: 13, offset: 22033},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 519, col: 19, offset: 22039},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 519, col: 19, offset: 22039},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 519, col: 24, offset: 22044},
									expr: &charClassMatcher{
										pos:        position{line: 519, col: 24, offset: 22044},
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
			pos:  position{line: 524, col: 1, offset: 22139},
			expr: &actionExpr{
				pos: position{line: 524, col: 15, offset: 22153},
				run: (*parser).callonNonZNumber1,
				expr: &seqExpr{
					pos: position{line: 524, col: 15, offset: 22153},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 524, col: 15, offset: 22153},
							val:        "[1-9]",
							ranges:     []rune{'1', '9'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 524, col: 20, offset: 22158},
							expr: &charClassMatcher{
								pos:        position{line: 524, col: 20, offset: 22158},
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
			pos:  position{line: 529, col: 1, offset: 22251},
			expr: &actionExpr{
				pos: position{line: 529, col: 12, offset: 22262},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 529, col: 12, offset: 22262},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 529, col: 20, offset: 22270},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 529, col: 20, offset: 22270},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 529, col: 34, offset: 22284},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 533, col: 1, offset: 22326},
			expr: &actionExpr{
				pos: position{line: 533, col: 16, offset: 22341},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 533, col: 18, offset: 22343},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 533, col: 18, offset: 22343},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 533, col: 27, offset: 22352},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 533, col: 27, offset: 22352},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 533, col: 31, offset: 22356},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 533, col: 49, offset: 22374},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 533, col: 55, offset: 22380},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 537, col: 1, offset: 22426},
			expr: &choiceExpr{
				pos: position{line: 537, col: 24, offset: 22449},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 537, col: 24, offset: 22449},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 537, col: 33, offset: 22458},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 537, col: 41, offset: 22466},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 537, col: 48, offset: 22473},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 537, col: 54, offset: 22479},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 539, col: 1, offset: 22486},
			expr: &actionExpr{
				pos: position{line: 539, col: 17, offset: 22502},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 539, col: 19, offset: 22504},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 539, col: 19, offset: 22504},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 539, col: 29, offset: 22514},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 539, col: 29, offset: 22514},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 539, col: 33, offset: 22518},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 539, col: 52, offset: 22537},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 539, col: 58, offset: 22543},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 543, col: 1, offset: 22591},
			expr: &choiceExpr{
				pos: position{line: 543, col: 25, offset: 22615},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 543, col: 25, offset: 22615},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 543, col: 35, offset: 22625},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 543, col: 42, offset: 22632},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 543, col: 50, offset: 22640},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 543, col: 56, offset: 22646},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 555, col: 1, offset: 23161},
			expr: &actionExpr{
				pos: position{line: 555, col: 12, offset: 23172},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 555, col: 14, offset: 23174},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 555, col: 14, offset: 23174},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 555, col: 34, offset: 23194},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 559, col: 1, offset: 23238},
			expr: &seqExpr{
				pos: position{line: 559, col: 21, offset: 23258},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 559, col: 21, offset: 23258},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 559, col: 26, offset: 23263},
						expr: &anyMatcher{
							line: 559, col: 26, offset: 23263,
						},
					},
					&litMatcher{
						pos:        position{line: 559, col: 29, offset: 23266},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 559, col: 34, offset: 23271},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 561, col: 1, offset: 23276},
			expr: &seqExpr{
				pos: position{line: 561, col: 22, offset: 23297},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 561, col: 22, offset: 23297},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 561, col: 27, offset: 23302},
						expr: &charClassMatcher{
							pos:        position{line: 561, col: 27, offset: 23302},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 561, col: 36, offset: 23311},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 563, col: 1, offset: 23316},
			expr: &seqExpr{
				pos: position{line: 563, col: 9, offset: 23324},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 563, col: 9, offset: 23324},
						expr: &charClassMatcher{
							pos:        position{line: 563, col: 9, offset: 23324},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 563, col: 17, offset: 23332},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 563, col: 17, offset: 23332},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 563, col: 26, offset: 23341},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 563, col: 35, offset: 23350},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 563, col: 42, offset: 23357},
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
			pos:         position{line: 565, col: 1, offset: 23364},
			expr: &zeroOrMoreExpr{
				pos: position{line: 565, col: 19, offset: 23382},
				expr: &charClassMatcher{
					pos:        position{line: 565, col: 19, offset: 23382},
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
			pos:         position{line: 567, col: 1, offset: 23394},
			expr: &oneOrMoreExpr{
				pos: position{line: 567, col: 31, offset: 23424},
				expr: &charClassMatcher{
					pos:        position{line: 567, col: 31, offset: 23424},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 569, col: 1, offset: 23436},
			expr: &notExpr{
				pos: position{line: 569, col: 8, offset: 23443},
				expr: &anyMatcher{
					line: 569, col: 9, offset: 23444,
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

func (c *current) onColumnDef1(name, dataType, collation, constraint interface{}) (interface{}, error) {
	if dataType == nil {
		return nil, nil
	}
	result := make(map[string]string)
	if err := mergo.Merge(&result, dataType.(map[string]string), mergo.WithOverride); err != nil {
		return nil, err
	}
	if collation != nil {
		result["collation"] = collation.(string)
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
	return p.cur.onColumnDef1(stack["name"], stack["dataType"], stack["collation"], stack["constraint"])
}

func (c *current) onCollate1(collation interface{}) (interface{}, error) {
	return collation, nil
}

func (p *parser) callonCollate1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCollate1(stack["collation"])
}

func (c *current) onCollation1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonCollation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCollation1()
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
