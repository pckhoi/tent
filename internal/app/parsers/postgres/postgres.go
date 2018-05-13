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
								name: "BitVarT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 67, offset: 2933},
								name: "BitT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 74, offset: 2940},
								name: "IntT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 81, offset: 2947},
								name: "OtherT",
							},
							&ruleRefExpr{
								pos:  position{line: 56, col: 90, offset: 2956},
								name: "CustomT",
							},
						},
					},
				},
			},
		},
		{
			name: "TimestampT",
			pos:  position{line: 60, col: 1, offset: 2989},
			expr: &actionExpr{
				pos: position{line: 60, col: 15, offset: 3003},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 60, col: 15, offset: 3003},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 60, col: 15, offset: 3003},
							val:        "timestamp",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 60, col: 28, offset: 3016},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 60, col: 33, offset: 3021},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 60, col: 46, offset: 3034},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 60, col: 59, offset: 3047},
								expr: &choiceExpr{
									pos: position{line: 60, col: 61, offset: 3049},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 60, col: 61, offset: 3049},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 60, col: 70, offset: 3058},
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
			pos:  position{line: 73, col: 1, offset: 3337},
			expr: &actionExpr{
				pos: position{line: 73, col: 10, offset: 3346},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 73, col: 10, offset: 3346},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 73, col: 10, offset: 3346},
							val:        "time",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 73, col: 18, offset: 3354},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 73, col: 23, offset: 3359},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 73, col: 36, offset: 3372},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 73, col: 49, offset: 3385},
								expr: &choiceExpr{
									pos: position{line: 73, col: 51, offset: 3387},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 73, col: 51, offset: 3387},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 73, col: 60, offset: 3396},
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
			pos:  position{line: 86, col: 1, offset: 3667},
			expr: &actionExpr{
				pos: position{line: 86, col: 17, offset: 3683},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 86, col: 17, offset: 3683},
					expr: &seqExpr{
						pos: position{line: 86, col: 19, offset: 3685},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 86, col: 19, offset: 3685},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 86, col: 22, offset: 3688},
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
			pos:  position{line: 93, col: 1, offset: 3816},
			expr: &actionExpr{
				pos: position{line: 93, col: 11, offset: 3826},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 93, col: 11, offset: 3826},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 93, col: 11, offset: 3826},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 93, col: 14, offset: 3829},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 93, col: 22, offset: 3837},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 93, col: 25, offset: 3840},
							val:        "time",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 93, col: 33, offset: 3848},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 93, col: 36, offset: 3851},
							val:        "zone",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 97, col: 1, offset: 3885},
			expr: &actionExpr{
				pos: position{line: 97, col: 14, offset: 3898},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 97, col: 14, offset: 3898},
					expr: &seqExpr{
						pos: position{line: 97, col: 16, offset: 3900},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 97, col: 16, offset: 3900},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 97, col: 19, offset: 3903},
								val:        "without",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 97, col: 30, offset: 3914},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 97, col: 33, offset: 3917},
								val:        "time",
								ignoreCase: true,
							},
							&ruleRefExpr{
								pos:  position{line: 97, col: 41, offset: 3925},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 97, col: 44, offset: 3928},
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
			pos:  position{line: 101, col: 1, offset: 3966},
			expr: &actionExpr{
				pos: position{line: 101, col: 10, offset: 3975},
				run: (*parser).callonCharT1,
				expr: &seqExpr{
					pos: position{line: 101, col: 10, offset: 3975},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 101, col: 12, offset: 3977},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 101, col: 12, offset: 3977},
									val:        "character",
									ignoreCase: true,
								},
								&litMatcher{
									pos:        position{line: 101, col: 27, offset: 3992},
									val:        "char",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 101, col: 37, offset: 4002},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 101, col: 44, offset: 4009},
								expr: &seqExpr{
									pos: position{line: 101, col: 46, offset: 4011},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 101, col: 46, offset: 4011},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 101, col: 50, offset: 4015},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 101, col: 61, offset: 4026},
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
			pos:  position{line: 113, col: 1, offset: 4281},
			expr: &actionExpr{
				pos: position{line: 113, col: 13, offset: 4293},
				run: (*parser).callonVarcharT1,
				expr: &seqExpr{
					pos: position{line: 113, col: 13, offset: 4293},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 113, col: 15, offset: 4295},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 113, col: 17, offset: 4297},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 113, col: 17, offset: 4297},
											val:        "character",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 113, col: 30, offset: 4310},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 113, col: 33, offset: 4313},
											val:        "varying",
											ignoreCase: true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 113, col: 48, offset: 4328},
									val:        "varchar",
									ignoreCase: true,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 113, col: 61, offset: 4341},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 113, col: 68, offset: 4348},
								expr: &seqExpr{
									pos: position{line: 113, col: 70, offset: 4350},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 113, col: 70, offset: 4350},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 113, col: 74, offset: 4354},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 113, col: 85, offset: 4365},
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
			pos:  position{line: 124, col: 1, offset: 4600},
			expr: &actionExpr{
				pos: position{line: 124, col: 9, offset: 4608},
				run: (*parser).callonBitT1,
				expr: &seqExpr{
					pos: position{line: 124, col: 9, offset: 4608},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 124, col: 9, offset: 4608},
							val:        "bit",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 124, col: 16, offset: 4615},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 124, col: 23, offset: 4622},
								expr: &seqExpr{
									pos: position{line: 124, col: 25, offset: 4624},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 124, col: 25, offset: 4624},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 124, col: 29, offset: 4628},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 124, col: 40, offset: 4639},
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
			pos:  position{line: 136, col: 1, offset: 4893},
			expr: &actionExpr{
				pos: position{line: 136, col: 12, offset: 4904},
				run: (*parser).callonBitVarT1,
				expr: &seqExpr{
					pos: position{line: 136, col: 12, offset: 4904},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 136, col: 12, offset: 4904},
							val:        "bit",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 136, col: 19, offset: 4911},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 136, col: 22, offset: 4914},
							val:        "varying",
							ignoreCase: true,
						},
						&labeledExpr{
							pos:   position{line: 136, col: 33, offset: 4925},
							label: "length",
							expr: &zeroOrOneExpr{
								pos: position{line: 136, col: 40, offset: 4932},
								expr: &seqExpr{
									pos: position{line: 136, col: 42, offset: 4934},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 136, col: 42, offset: 4934},
											val:        "(",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 136, col: 46, offset: 4938},
											name: "NonZNumber",
										},
										&litMatcher{
											pos:        position{line: 136, col: 57, offset: 4949},
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
			pos:  position{line: 147, col: 1, offset: 5183},
			expr: &actionExpr{
				pos: position{line: 147, col: 9, offset: 5191},
				run: (*parser).callonIntT1,
				expr: &choiceExpr{
					pos: position{line: 147, col: 11, offset: 5193},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 147, col: 11, offset: 5193},
							val:        "integer",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 147, col: 24, offset: 5206},
							val:        "int",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 153, col: 1, offset: 5288},
			expr: &actionExpr{
				pos: position{line: 153, col: 11, offset: 5298},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 153, col: 13, offset: 5300},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 153, col: 13, offset: 5300},
							val:        "date",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 23, offset: 5310},
							val:        "smallint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 37, offset: 5324},
							val:        "bigint",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 49, offset: 5336},
							val:        "decimal",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 62, offset: 5349},
							val:        "numeric",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 75, offset: 5362},
							val:        "real",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 85, offset: 5372},
							val:        "smallserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 102, offset: 5389},
							val:        "serial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 114, offset: 5401},
							val:        "bigserial",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 129, offset: 5416},
							val:        "boolean",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 142, offset: 5429},
							val:        "text",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 152, offset: 5439},
							val:        "money",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 163, offset: 5450},
							val:        "bytea",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 174, offset: 5461},
							val:        "point",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 185, offset: 5472},
							val:        "line",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 195, offset: 5482},
							val:        "lseg",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 205, offset: 5492},
							val:        "box",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 214, offset: 5501},
							val:        "path",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 224, offset: 5511},
							val:        "polygon",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 237, offset: 5524},
							val:        "circle",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 249, offset: 5536},
							val:        "cidr",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 259, offset: 5546},
							val:        "inet",
							ignoreCase: true,
						},
						&litMatcher{
							pos:        position{line: 153, col: 269, offset: 5556},
							val:        "macaddr",
							ignoreCase: true,
						},
					},
				},
			},
		},
		{
			name: "CustomT",
			pos:  position{line: 159, col: 1, offset: 5664},
			expr: &actionExpr{
				pos: position{line: 159, col: 13, offset: 5676},
				run: (*parser).callonCustomT1,
				expr: &ruleRefExpr{
					pos:  position{line: 159, col: 13, offset: 5676},
					name: "Ident",
				},
			},
		},
		{
			name: "CreateTypeStmt",
			pos:  position{line: 180, col: 2, offset: 7206},
			expr: &actionExpr{
				pos: position{line: 180, col: 20, offset: 7224},
				run: (*parser).callonCreateTypeStmt1,
				expr: &seqExpr{
					pos: position{line: 180, col: 20, offset: 7224},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 180, col: 20, offset: 7224},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 180, col: 30, offset: 7234},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 180, col: 33, offset: 7237},
							val:        "type",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 180, col: 41, offset: 7245},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 180, col: 44, offset: 7248},
							label: "typename",
							expr: &ruleRefExpr{
								pos:  position{line: 180, col: 53, offset: 7257},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 180, col: 59, offset: 7263},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 180, col: 62, offset: 7266},
							val:        "as",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 180, col: 68, offset: 7272},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 180, col: 71, offset: 7275},
							label: "typedef",
							expr: &ruleRefExpr{
								pos:  position{line: 180, col: 79, offset: 7283},
								name: "EnumDef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 180, col: 87, offset: 7291},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 180, col: 89, offset: 7293},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 180, col: 93, offset: 7297},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "EnumDef",
			pos:  position{line: 186, col: 2, offset: 7414},
			expr: &actionExpr{
				pos: position{line: 186, col: 13, offset: 7425},
				run: (*parser).callonEnumDef1,
				expr: &seqExpr{
					pos: position{line: 186, col: 13, offset: 7425},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 186, col: 13, offset: 7425},
							val:        "ENUM",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 20, offset: 7432},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 186, col: 22, offset: 7434},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 26, offset: 7438},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 186, col: 28, offset: 7440},
							label: "vals",
							expr: &seqExpr{
								pos: position{line: 186, col: 35, offset: 7447},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 186, col: 35, offset: 7447},
										name: "StringConst",
									},
									&zeroOrMoreExpr{
										pos: position{line: 186, col: 47, offset: 7459},
										expr: &seqExpr{
											pos: position{line: 186, col: 49, offset: 7461},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 186, col: 49, offset: 7461},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 186, col: 51, offset: 7463},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 186, col: 55, offset: 7467},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 186, col: 57, offset: 7469},
													name: "StringConst",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 75, offset: 7487},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 186, col: 77, offset: 7489},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 211, col: 1, offset: 9093},
			expr: &actionExpr{
				pos: position{line: 211, col: 25, offset: 9117},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 211, col: 25, offset: 9117},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 211, col: 25, offset: 9117},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 36, offset: 9128},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 211, col: 39, offset: 9131},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 45, offset: 9137},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 211, col: 48, offset: 9140},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 61, offset: 9153},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 211, col: 63, offset: 9155},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 211, col: 73, offset: 9165},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 79, offset: 9171},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 211, col: 81, offset: 9173},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 87, offset: 9179},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 211, col: 89, offset: 9181},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 211, col: 97, offset: 9189},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 109, offset: 9201},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 211, col: 111, offset: 9203},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 115, offset: 9207},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 215, col: 1, offset: 9295},
			expr: &actionExpr{
				pos: position{line: 215, col: 24, offset: 9318},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 215, col: 24, offset: 9318},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 215, col: 24, offset: 9318},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 215, col: 34, offset: 9328},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 215, col: 37, offset: 9331},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 215, col: 50, offset: 9344},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 215, col: 53, offset: 9347},
							expr: &seqExpr{
								pos: position{line: 215, col: 55, offset: 9349},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 215, col: 55, offset: 9349},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 215, col: 61, offset: 9355},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 215, col: 64, offset: 9358},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 215, col: 71, offset: 9365},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 215, col: 74, offset: 9368},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 215, col: 84, offset: 9378},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 215, col: 90, offset: 9384},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 215, col: 100, offset: 9394},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 215, col: 106, offset: 9400},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 215, col: 109, offset: 9403},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 215, col: 117, offset: 9411},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 215, col: 120, offset: 9414},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 215, col: 130, offset: 9424},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 215, col: 133, offset: 9427},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 215, col: 140, offset: 9434},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 215, col: 146, offset: 9440},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 215, col: 148, offset: 9442},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 215, col: 152, offset: 9446},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 219, col: 1, offset: 9536},
			expr: &actionExpr{
				pos: position{line: 219, col: 12, offset: 9547},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 219, col: 12, offset: 9547},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 219, col: 12, offset: 9547},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 219, col: 19, offset: 9554},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 219, col: 21, offset: 9556},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 219, col: 25, offset: 9560},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 219, col: 29, offset: 9564},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 219, col: 33, offset: 9568},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 219, col: 33, offset: 9568},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 219, col: 39, offset: 9574},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 219, col: 47, offset: 9582},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 219, col: 49, offset: 9584},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 219, col: 56, offset: 9591},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 219, col: 77, offset: 9612},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 219, col: 79, offset: 9614},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 219, col: 83, offset: 9618},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 224, col: 1, offset: 9699},
			expr: &actionExpr{
				pos: position{line: 224, col: 8, offset: 9706},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 224, col: 8, offset: 9706},
					expr: &charClassMatcher{
						pos:        position{line: 224, col: 8, offset: 9706},
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
			pos:  position{line: 239, col: 1, offset: 10546},
			expr: &actionExpr{
				pos: position{line: 239, col: 25, offset: 10570},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 239, col: 25, offset: 10570},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 239, col: 32, offset: 10577},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 239, col: 32, offset: 10577},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 239, col: 38, offset: 10583},
								expr: &seqExpr{
									pos: position{line: 239, col: 40, offset: 10585},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 239, col: 40, offset: 10585},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 239, col: 42, offset: 10587},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 239, col: 46, offset: 10591},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 239, col: 48, offset: 10593},
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
			pos:  position{line: 251, col: 1, offset: 10883},
			expr: &choiceExpr{
				pos: position{line: 251, col: 12, offset: 10894},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 251, col: 12, offset: 10894},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 251, col: 21, offset: 10903},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 251, col: 31, offset: 10913},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 251, col: 45, offset: 10927},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 253, col: 1, offset: 10936},
			expr: &actionExpr{
				pos: position{line: 253, col: 16, offset: 10951},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 253, col: 16, offset: 10951},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 253, col: 16, offset: 10951},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 253, col: 20, offset: 10955},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 253, col: 26, offset: 10961},
								expr: &choiceExpr{
									pos: position{line: 253, col: 27, offset: 10962},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 253, col: 27, offset: 10962},
											val:        "[^'\\n]",
											chars:      []rune{'\'', '\n'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 253, col: 36, offset: 10971},
											val:        "''",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 253, col: 43, offset: 10978},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 257, col: 1, offset: 11030},
			expr: &actionExpr{
				pos: position{line: 257, col: 10, offset: 11039},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 257, col: 10, offset: 11039},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 257, col: 10, offset: 11039},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 257, col: 18, offset: 11047},
							expr: &charClassMatcher{
								pos:        position{line: 257, col: 18, offset: 11047},
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
			pos:  position{line: 261, col: 1, offset: 11100},
			expr: &actionExpr{
				pos: position{line: 261, col: 11, offset: 11110},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 261, col: 13, offset: 11112},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 261, col: 13, offset: 11112},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 261, col: 19, offset: 11118},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 261, col: 19, offset: 11118},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 261, col: 24, offset: 11123},
									expr: &charClassMatcher{
										pos:        position{line: 261, col: 24, offset: 11123},
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
			pos:  position{line: 266, col: 1, offset: 11218},
			expr: &actionExpr{
				pos: position{line: 266, col: 15, offset: 11232},
				run: (*parser).callonNonZNumber1,
				expr: &seqExpr{
					pos: position{line: 266, col: 15, offset: 11232},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 266, col: 15, offset: 11232},
							val:        "[1-9]",
							ranges:     []rune{'1', '9'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 266, col: 20, offset: 11237},
							expr: &charClassMatcher{
								pos:        position{line: 266, col: 20, offset: 11237},
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
			pos:  position{line: 271, col: 1, offset: 11330},
			expr: &actionExpr{
				pos: position{line: 271, col: 12, offset: 11341},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 271, col: 12, offset: 11341},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 271, col: 20, offset: 11349},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 271, col: 20, offset: 11349},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 271, col: 34, offset: 11363},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 275, col: 1, offset: 11405},
			expr: &actionExpr{
				pos: position{line: 275, col: 16, offset: 11420},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 275, col: 18, offset: 11422},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 275, col: 18, offset: 11422},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 275, col: 27, offset: 11431},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 275, col: 27, offset: 11431},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 275, col: 31, offset: 11435},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 275, col: 49, offset: 11453},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 55, offset: 11459},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 279, col: 1, offset: 11505},
			expr: &choiceExpr{
				pos: position{line: 279, col: 24, offset: 11528},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 279, col: 24, offset: 11528},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 279, col: 33, offset: 11537},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 279, col: 41, offset: 11545},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 279, col: 48, offset: 11552},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 279, col: 54, offset: 11558},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 281, col: 1, offset: 11565},
			expr: &actionExpr{
				pos: position{line: 281, col: 17, offset: 11581},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 281, col: 19, offset: 11583},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 281, col: 19, offset: 11583},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 281, col: 29, offset: 11593},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 281, col: 29, offset: 11593},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 281, col: 33, offset: 11597},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 281, col: 52, offset: 11616},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 281, col: 58, offset: 11622},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 285, col: 1, offset: 11670},
			expr: &choiceExpr{
				pos: position{line: 285, col: 25, offset: 11694},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 285, col: 25, offset: 11694},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 285, col: 35, offset: 11704},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 285, col: 42, offset: 11711},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 285, col: 50, offset: 11719},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 285, col: 56, offset: 11725},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 297, col: 1, offset: 12240},
			expr: &actionExpr{
				pos: position{line: 297, col: 12, offset: 12251},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 297, col: 14, offset: 12253},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 297, col: 14, offset: 12253},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 297, col: 34, offset: 12273},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 301, col: 1, offset: 12317},
			expr: &seqExpr{
				pos: position{line: 301, col: 21, offset: 12337},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 301, col: 21, offset: 12337},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 301, col: 26, offset: 12342},
						expr: &seqExpr{
							pos: position{line: 301, col: 28, offset: 12344},
							exprs: []interface{}{
								&zeroOrMoreExpr{
									pos: position{line: 301, col: 28, offset: 12344},
									expr: &anyMatcher{
										line: 301, col: 28, offset: 12344,
									},
								},
								&zeroOrOneExpr{
									pos: position{line: 301, col: 31, offset: 12347},
									expr: &ruleRefExpr{
										pos:  position{line: 301, col: 31, offset: 12347},
										name: "MultilineComment",
									},
								},
								&zeroOrMoreExpr{
									pos: position{line: 301, col: 49, offset: 12365},
									expr: &anyMatcher{
										line: 301, col: 49, offset: 12365,
									},
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 301, col: 55, offset: 12371},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 301, col: 60, offset: 12376},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 303, col: 1, offset: 12381},
			expr: &seqExpr{
				pos: position{line: 303, col: 22, offset: 12402},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 303, col: 22, offset: 12402},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 303, col: 27, offset: 12407},
						expr: &charClassMatcher{
							pos:        position{line: 303, col: 27, offset: 12407},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 303, col: 36, offset: 12416},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 305, col: 1, offset: 12421},
			expr: &seqExpr{
				pos: position{line: 305, col: 9, offset: 12429},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 305, col: 9, offset: 12429},
						expr: &charClassMatcher{
							pos:        position{line: 305, col: 9, offset: 12429},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 305, col: 17, offset: 12437},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 305, col: 17, offset: 12437},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 305, col: 26, offset: 12446},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 305, col: 35, offset: 12455},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 305, col: 42, offset: 12462},
								val:        "\n",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 305, col: 49, offset: 12469},
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
			pos:         position{line: 307, col: 1, offset: 12475},
			expr: &zeroOrMoreExpr{
				pos: position{line: 307, col: 19, offset: 12493},
				expr: &charClassMatcher{
					pos:        position{line: 307, col: 19, offset: 12493},
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
			pos:         position{line: 309, col: 1, offset: 12505},
			expr: &oneOrMoreExpr{
				pos: position{line: 309, col: 31, offset: 12535},
				expr: &charClassMatcher{
					pos:        position{line: 309, col: 31, offset: 12535},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 311, col: 1, offset: 12547},
			expr: &notExpr{
				pos: position{line: 311, col: 8, offset: 12554},
				expr: &anyMatcher{
					line: 311, col: 9, offset: 12555,
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
