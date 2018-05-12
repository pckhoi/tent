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

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func toByteSlice(v interface{}) []byte {
	valsSl := toIfaceSlice(v)
	var result []byte
	for _, val := range valsSl {
		result = append(result, val.([]byte)[0])
	}
	return result
}

var g = &grammar{
	rules: []*rule{
		{
			name: "SQL",
			pos:  position{line: 22, col: 1, offset: 478},
			expr: &actionExpr{
				pos: position{line: 22, col: 8, offset: 485},
				run: (*parser).callonSQL1,
				expr: &labeledExpr{
					pos:   position{line: 22, col: 8, offset: 485},
					label: "stmts",
					expr: &oneOrMoreExpr{
						pos: position{line: 22, col: 14, offset: 491},
						expr: &ruleRefExpr{
							pos:  position{line: 22, col: 14, offset: 491},
							name: "Stmt",
						},
					},
				},
			},
		},
		{
			name: "Stmt",
			pos:  position{line: 26, col: 1, offset: 524},
			expr: &actionExpr{
				pos: position{line: 26, col: 9, offset: 532},
				run: (*parser).callonStmt1,
				expr: &seqExpr{
					pos: position{line: 26, col: 9, offset: 532},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 26, col: 9, offset: 532},
							expr: &ruleRefExpr{
								pos:  position{line: 26, col: 9, offset: 532},
								name: "Comment",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 26, col: 18, offset: 541},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 26, col: 20, offset: 543},
							label: "stmt",
							expr: &choiceExpr{
								pos: position{line: 26, col: 27, offset: 550},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 26, col: 27, offset: 550},
										name: "SetStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 26, col: 37, offset: 560},
										name: "CreateTableStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 26, col: 55, offset: 578},
										name: "CreateExtensionStmt",
									},
									&ruleRefExpr{
										pos:  position{line: 26, col: 77, offset: 600},
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
			pos:  position{line: 40, col: 1, offset: 2177},
			expr: &actionExpr{
				pos: position{line: 40, col: 20, offset: 2196},
				run: (*parser).callonCreateTableStmt1,
				expr: &seqExpr{
					pos: position{line: 40, col: 20, offset: 2196},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 40, col: 20, offset: 2196},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 40, col: 30, offset: 2206},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 40, col: 33, offset: 2209},
							val:        "table",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 40, col: 42, offset: 2218},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 40, col: 45, offset: 2221},
							label: "tablename",
							expr: &ruleRefExpr{
								pos:  position{line: 40, col: 55, offset: 2231},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 40, col: 61, offset: 2237},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 40, col: 63, offset: 2239},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 40, col: 67, offset: 2243},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 40, col: 69, offset: 2245},
							label: "fields",
							expr: &seqExpr{
								pos: position{line: 40, col: 78, offset: 2254},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 40, col: 78, offset: 2254},
										name: "FieldDef",
									},
									&zeroOrMoreExpr{
										pos: position{line: 40, col: 87, offset: 2263},
										expr: &seqExpr{
											pos: position{line: 40, col: 89, offset: 2265},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 40, col: 89, offset: 2265},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 40, col: 91, offset: 2267},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 40, col: 95, offset: 2271},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 40, col: 97, offset: 2273},
													name: "FieldDef",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 40, col: 111, offset: 2287},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 40, col: 113, offset: 2289},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 40, col: 117, offset: 2293},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 40, col: 119, offset: 2295},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 40, col: 123, offset: 2299},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "FieldDef",
			pos:  position{line: 52, col: 1, offset: 2711},
			expr: &actionExpr{
				pos: position{line: 52, col: 13, offset: 2723},
				run: (*parser).callonFieldDef1,
				expr: &seqExpr{
					pos: position{line: 52, col: 13, offset: 2723},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 52, col: 13, offset: 2723},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 52, col: 18, offset: 2728},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 52, col: 24, offset: 2734},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 52, col: 27, offset: 2737},
							label: "dataType",
							expr: &ruleRefExpr{
								pos:  position{line: 52, col: 36, offset: 2746},
								name: "DataType",
							},
						},
						&labeledExpr{
							pos:   position{line: 52, col: 45, offset: 2755},
							label: "notnull",
							expr: &zeroOrOneExpr{
								pos: position{line: 52, col: 53, offset: 2763},
								expr: &seqExpr{
									pos: position{line: 52, col: 55, offset: 2765},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 52, col: 55, offset: 2765},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 52, col: 58, offset: 2768},
											val:        "not",
											ignoreCase: true,
										},
										&ruleRefExpr{
											pos:  position{line: 52, col: 65, offset: 2775},
											name: "_1",
										},
										&litMatcher{
											pos:        position{line: 52, col: 68, offset: 2778},
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
			pos:  position{line: 61, col: 1, offset: 2971},
			expr: &actionExpr{
				pos: position{line: 61, col: 13, offset: 2983},
				run: (*parser).callonDataType1,
				expr: &labeledExpr{
					pos:   position{line: 61, col: 13, offset: 2983},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 61, col: 17, offset: 2987},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 61, col: 17, offset: 2987},
								name: "TimestampT",
							},
							&ruleRefExpr{
								pos:  position{line: 61, col: 30, offset: 3000},
								name: "TimeT",
							},
							&ruleRefExpr{
								pos:  position{line: 61, col: 38, offset: 3008},
								name: "OtherT",
							},
						},
					},
				},
			},
		},
		{
			name: "TimestampT",
			pos:  position{line: 65, col: 1, offset: 3040},
			expr: &actionExpr{
				pos: position{line: 65, col: 15, offset: 3054},
				run: (*parser).callonTimestampT1,
				expr: &seqExpr{
					pos: position{line: 65, col: 15, offset: 3054},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 65, col: 15, offset: 3054},
							val:        "timestamp",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 65, col: 27, offset: 3066},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 65, col: 32, offset: 3071},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 65, col: 45, offset: 3084},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 65, col: 58, offset: 3097},
								expr: &choiceExpr{
									pos: position{line: 65, col: 60, offset: 3099},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 65, col: 60, offset: 3099},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 65, col: 69, offset: 3108},
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
			pos:  position{line: 78, col: 1, offset: 3387},
			expr: &actionExpr{
				pos: position{line: 78, col: 10, offset: 3396},
				run: (*parser).callonTimeT1,
				expr: &seqExpr{
					pos: position{line: 78, col: 10, offset: 3396},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 78, col: 10, offset: 3396},
							val:        "time",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 78, col: 17, offset: 3403},
							label: "prec",
							expr: &ruleRefExpr{
								pos:  position{line: 78, col: 22, offset: 3408},
								name: "SecPrecision",
							},
						},
						&labeledExpr{
							pos:   position{line: 78, col: 35, offset: 3421},
							label: "withTimeZone",
							expr: &zeroOrOneExpr{
								pos: position{line: 78, col: 48, offset: 3434},
								expr: &choiceExpr{
									pos: position{line: 78, col: 50, offset: 3436},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 78, col: 50, offset: 3436},
											name: "WithTZ",
										},
										&ruleRefExpr{
											pos:  position{line: 78, col: 59, offset: 3445},
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
			pos:  position{line: 91, col: 1, offset: 3716},
			expr: &actionExpr{
				pos: position{line: 91, col: 17, offset: 3732},
				run: (*parser).callonSecPrecision1,
				expr: &zeroOrOneExpr{
					pos: position{line: 91, col: 17, offset: 3732},
					expr: &seqExpr{
						pos: position{line: 91, col: 19, offset: 3734},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 91, col: 19, offset: 3734},
								name: "_1",
							},
							&charClassMatcher{
								pos:        position{line: 91, col: 22, offset: 3737},
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
			pos:  position{line: 98, col: 1, offset: 3865},
			expr: &actionExpr{
				pos: position{line: 98, col: 11, offset: 3875},
				run: (*parser).callonWithTZ1,
				expr: &seqExpr{
					pos: position{line: 98, col: 11, offset: 3875},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 98, col: 11, offset: 3875},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 98, col: 14, offset: 3878},
							val:        "with",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 98, col: 21, offset: 3885},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 98, col: 24, offset: 3888},
							val:        "time",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 98, col: 31, offset: 3895},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 98, col: 34, offset: 3898},
							val:        "zone",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "WithoutTZ",
			pos:  position{line: 102, col: 1, offset: 3931},
			expr: &actionExpr{
				pos: position{line: 102, col: 14, offset: 3944},
				run: (*parser).callonWithoutTZ1,
				expr: &zeroOrOneExpr{
					pos: position{line: 102, col: 14, offset: 3944},
					expr: &seqExpr{
						pos: position{line: 102, col: 16, offset: 3946},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 102, col: 16, offset: 3946},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 102, col: 19, offset: 3949},
								val:        "without",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 102, col: 29, offset: 3959},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 102, col: 32, offset: 3962},
								val:        "time",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 102, col: 39, offset: 3969},
								name: "_1",
							},
							&litMatcher{
								pos:        position{line: 102, col: 42, offset: 3972},
								val:        "zone",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "OtherT",
			pos:  position{line: 106, col: 1, offset: 4009},
			expr: &actionExpr{
				pos: position{line: 106, col: 11, offset: 4019},
				run: (*parser).callonOtherT1,
				expr: &choiceExpr{
					pos: position{line: 106, col: 13, offset: 4021},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 106, col: 13, offset: 4021},
							val:        "date",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 22, offset: 4030},
							val:        "integer",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 34, offset: 4042},
							val:        "smallint",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 47, offset: 4055},
							val:        "bigint",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 58, offset: 4066},
							val:        "decimal",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 70, offset: 4078},
							val:        "numeric",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 82, offset: 4090},
							val:        "real",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 91, offset: 4099},
							val:        "smallserial",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 107, offset: 4115},
							val:        "serial",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 118, offset: 4126},
							val:        "bigserial",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 132, offset: 4140},
							val:        "boolean",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CommentExtensionStmt",
			pos:  position{line: 122, col: 1, offset: 5470},
			expr: &actionExpr{
				pos: position{line: 122, col: 25, offset: 5494},
				run: (*parser).callonCommentExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 122, col: 25, offset: 5494},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 122, col: 25, offset: 5494},
							val:        "comment",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 122, col: 36, offset: 5505},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 122, col: 39, offset: 5508},
							val:        "on",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 122, col: 45, offset: 5514},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 122, col: 48, offset: 5517},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 122, col: 61, offset: 5530},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 122, col: 63, offset: 5532},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 122, col: 73, offset: 5542},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 122, col: 79, offset: 5548},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 122, col: 81, offset: 5550},
							val:        "is",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 122, col: 87, offset: 5556},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 122, col: 89, offset: 5558},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 122, col: 97, offset: 5566},
								name: "StringConst",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 122, col: 109, offset: 5578},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 122, col: 111, offset: 5580},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 122, col: 115, offset: 5584},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "CreateExtensionStmt",
			pos:  position{line: 126, col: 1, offset: 5677},
			expr: &actionExpr{
				pos: position{line: 126, col: 24, offset: 5700},
				run: (*parser).callonCreateExtensionStmt1,
				expr: &seqExpr{
					pos: position{line: 126, col: 24, offset: 5700},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 126, col: 24, offset: 5700},
							val:        "create",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 126, col: 34, offset: 5710},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 126, col: 37, offset: 5713},
							val:        "extension",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 126, col: 50, offset: 5726},
							name: "_1",
						},
						&zeroOrOneExpr{
							pos: position{line: 126, col: 53, offset: 5729},
							expr: &seqExpr{
								pos: position{line: 126, col: 55, offset: 5731},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 126, col: 55, offset: 5731},
										val:        "if",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 126, col: 61, offset: 5737},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 126, col: 64, offset: 5740},
										val:        "not",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 126, col: 71, offset: 5747},
										name: "_1",
									},
									&litMatcher{
										pos:        position{line: 126, col: 74, offset: 5750},
										val:        "exists",
										ignoreCase: true,
									},
									&ruleRefExpr{
										pos:  position{line: 126, col: 84, offset: 5760},
										name: "_1",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 126, col: 90, offset: 5766},
							label: "extension",
							expr: &ruleRefExpr{
								pos:  position{line: 126, col: 100, offset: 5776},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 126, col: 106, offset: 5782},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 126, col: 109, offset: 5785},
							val:        "with",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 126, col: 117, offset: 5793},
							name: "_1",
						},
						&litMatcher{
							pos:        position{line: 126, col: 120, offset: 5796},
							val:        "schema",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 126, col: 130, offset: 5806},
							name: "_1",
						},
						&labeledExpr{
							pos:   position{line: 126, col: 133, offset: 5809},
							label: "schema",
							expr: &ruleRefExpr{
								pos:  position{line: 126, col: 140, offset: 5816},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 126, col: 146, offset: 5822},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 126, col: 148, offset: 5824},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 126, col: 152, offset: 5828},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "SetStmt",
			pos:  position{line: 130, col: 1, offset: 5923},
			expr: &actionExpr{
				pos: position{line: 130, col: 12, offset: 5934},
				run: (*parser).callonSetStmt1,
				expr: &seqExpr{
					pos: position{line: 130, col: 12, offset: 5934},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 130, col: 12, offset: 5934},
							val:        "set",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 130, col: 19, offset: 5941},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 130, col: 21, offset: 5943},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 130, col: 25, offset: 5947},
								name: "Key",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 130, col: 29, offset: 5951},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 130, col: 33, offset: 5955},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 130, col: 33, offset: 5955},
									val:        "=",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 130, col: 39, offset: 5961},
									val:        "to",
									ignoreCase: true,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 130, col: 47, offset: 5969},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 130, col: 49, offset: 5971},
							label: "values",
							expr: &ruleRefExpr{
								pos:  position{line: 130, col: 56, offset: 5978},
								name: "CommaSeparatedValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 130, col: 77, offset: 5999},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 130, col: 79, offset: 6001},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 130, col: 83, offset: 6005},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 135, col: 1, offset: 6089},
			expr: &actionExpr{
				pos: position{line: 135, col: 8, offset: 6096},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 135, col: 8, offset: 6096},
					expr: &charClassMatcher{
						pos:        position{line: 135, col: 8, offset: 6096},
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
			pos:  position{line: 150, col: 1, offset: 6936},
			expr: &actionExpr{
				pos: position{line: 150, col: 25, offset: 6960},
				run: (*parser).callonCommaSeparatedValues1,
				expr: &labeledExpr{
					pos:   position{line: 150, col: 25, offset: 6960},
					label: "vals",
					expr: &seqExpr{
						pos: position{line: 150, col: 32, offset: 6967},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 150, col: 32, offset: 6967},
								name: "Value",
							},
							&zeroOrMoreExpr{
								pos: position{line: 150, col: 38, offset: 6973},
								expr: &seqExpr{
									pos: position{line: 150, col: 40, offset: 6975},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 150, col: 40, offset: 6975},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 150, col: 42, offset: 6977},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 150, col: 46, offset: 6981},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 150, col: 48, offset: 6983},
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
			pos:  position{line: 162, col: 1, offset: 7273},
			expr: &choiceExpr{
				pos: position{line: 162, col: 12, offset: 7284},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 162, col: 12, offset: 7284},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 162, col: 21, offset: 7293},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 162, col: 31, offset: 7303},
						name: "StringConst",
					},
					&ruleRefExpr{
						pos:  position{line: 162, col: 45, offset: 7317},
						name: "Ident",
					},
				},
			},
		},
		{
			name: "StringConst",
			pos:  position{line: 164, col: 1, offset: 7326},
			expr: &actionExpr{
				pos: position{line: 164, col: 16, offset: 7341},
				run: (*parser).callonStringConst1,
				expr: &seqExpr{
					pos: position{line: 164, col: 16, offset: 7341},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 164, col: 16, offset: 7341},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 164, col: 20, offset: 7345},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 164, col: 26, offset: 7351},
								expr: &charClassMatcher{
									pos:        position{line: 164, col: 26, offset: 7351},
									val:        "[^'\\n]",
									chars:      []rune{'\'', '\n'},
									ignoreCase: false,
									inverted:   true,
								},
							},
						},
						&litMatcher{
							pos:        position{line: 164, col: 34, offset: 7359},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 168, col: 1, offset: 7411},
			expr: &actionExpr{
				pos: position{line: 168, col: 10, offset: 7420},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 168, col: 10, offset: 7420},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 168, col: 10, offset: 7420},
							val:        "[a-z_]i",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z'},
							ignoreCase: true,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 168, col: 18, offset: 7428},
							expr: &charClassMatcher{
								pos:        position{line: 168, col: 18, offset: 7428},
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
			pos:  position{line: 172, col: 1, offset: 7481},
			expr: &actionExpr{
				pos: position{line: 172, col: 11, offset: 7491},
				run: (*parser).callonNumber1,
				expr: &choiceExpr{
					pos: position{line: 172, col: 13, offset: 7493},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 172, col: 13, offset: 7493},
							val:        "0",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 172, col: 19, offset: 7499},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 172, col: 19, offset: 7499},
									val:        "[1-9]",
									ranges:     []rune{'1', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 172, col: 24, offset: 7504},
									expr: &charClassMatcher{
										pos:        position{line: 172, col: 24, offset: 7504},
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
			name: "Boolean",
			pos:  position{line: 177, col: 1, offset: 7599},
			expr: &actionExpr{
				pos: position{line: 177, col: 12, offset: 7610},
				run: (*parser).callonBoolean1,
				expr: &labeledExpr{
					pos:   position{line: 177, col: 12, offset: 7610},
					label: "value",
					expr: &choiceExpr{
						pos: position{line: 177, col: 20, offset: 7618},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 177, col: 20, offset: 7618},
								name: "BooleanTrue",
							},
							&ruleRefExpr{
								pos:  position{line: 177, col: 34, offset: 7632},
								name: "BooleanFalse",
							},
						},
					},
				},
			},
		},
		{
			name: "BooleanTrue",
			pos:  position{line: 181, col: 1, offset: 7674},
			expr: &actionExpr{
				pos: position{line: 181, col: 16, offset: 7689},
				run: (*parser).callonBooleanTrue1,
				expr: &choiceExpr{
					pos: position{line: 181, col: 18, offset: 7691},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 181, col: 18, offset: 7691},
							val:        "TRUE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 181, col: 27, offset: 7700},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 181, col: 27, offset: 7700},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 181, col: 31, offset: 7704},
									name: "BooleanTrueString",
								},
								&litMatcher{
									pos:        position{line: 181, col: 49, offset: 7722},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 181, col: 55, offset: 7728},
							name: "BooleanTrueString",
						},
					},
				},
			},
		},
		{
			name: "BooleanTrueString",
			pos:  position{line: 185, col: 1, offset: 7774},
			expr: &choiceExpr{
				pos: position{line: 185, col: 24, offset: 7797},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 185, col: 24, offset: 7797},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 185, col: 33, offset: 7806},
						val:        "yes",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 185, col: 41, offset: 7814},
						val:        "on",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 185, col: 48, offset: 7821},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 185, col: 54, offset: 7827},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "BooleanFalse",
			pos:  position{line: 187, col: 1, offset: 7834},
			expr: &actionExpr{
				pos: position{line: 187, col: 17, offset: 7850},
				run: (*parser).callonBooleanFalse1,
				expr: &choiceExpr{
					pos: position{line: 187, col: 19, offset: 7852},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 187, col: 19, offset: 7852},
							val:        "FALSE",
							ignoreCase: false,
						},
						&seqExpr{
							pos: position{line: 187, col: 29, offset: 7862},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 187, col: 29, offset: 7862},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 187, col: 33, offset: 7866},
									name: "BooleanFalseString",
								},
								&litMatcher{
									pos:        position{line: 187, col: 52, offset: 7885},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 187, col: 58, offset: 7891},
							name: "BooleanFalseString",
						},
					},
				},
			},
		},
		{
			name: "BooleanFalseString",
			pos:  position{line: 191, col: 1, offset: 7939},
			expr: &choiceExpr{
				pos: position{line: 191, col: 25, offset: 7963},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 191, col: 25, offset: 7963},
						val:        "false",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 191, col: 35, offset: 7973},
						val:        "no",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 191, col: 42, offset: 7980},
						val:        "off",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 191, col: 50, offset: 7988},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 191, col: 56, offset: 7994},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 203, col: 1, offset: 8509},
			expr: &actionExpr{
				pos: position{line: 203, col: 12, offset: 8520},
				run: (*parser).callonComment1,
				expr: &choiceExpr{
					pos: position{line: 203, col: 14, offset: 8522},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 203, col: 14, offset: 8522},
							name: "SingleLineComment",
						},
						&ruleRefExpr{
							pos:  position{line: 203, col: 34, offset: 8542},
							name: "MultilineComment",
						},
					},
				},
			},
		},
		{
			name: "MultilineComment",
			pos:  position{line: 207, col: 1, offset: 8586},
			expr: &seqExpr{
				pos: position{line: 207, col: 21, offset: 8606},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 207, col: 21, offset: 8606},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 207, col: 26, offset: 8611},
						expr: &seqExpr{
							pos: position{line: 207, col: 28, offset: 8613},
							exprs: []interface{}{
								&zeroOrMoreExpr{
									pos: position{line: 207, col: 28, offset: 8613},
									expr: &anyMatcher{
										line: 207, col: 28, offset: 8613,
									},
								},
								&zeroOrOneExpr{
									pos: position{line: 207, col: 31, offset: 8616},
									expr: &ruleRefExpr{
										pos:  position{line: 207, col: 31, offset: 8616},
										name: "MultilineComment",
									},
								},
								&zeroOrMoreExpr{
									pos: position{line: 207, col: 49, offset: 8634},
									expr: &anyMatcher{
										line: 207, col: 49, offset: 8634,
									},
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 207, col: 55, offset: 8640},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 207, col: 60, offset: 8645},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 209, col: 1, offset: 8650},
			expr: &seqExpr{
				pos: position{line: 209, col: 22, offset: 8671},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 209, col: 22, offset: 8671},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 209, col: 27, offset: 8676},
						expr: &charClassMatcher{
							pos:        position{line: 209, col: 27, offset: 8676},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 36, offset: 8685},
						name: "EOL",
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 211, col: 1, offset: 8690},
			expr: &seqExpr{
				pos: position{line: 211, col: 9, offset: 8698},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 211, col: 9, offset: 8698},
						expr: &charClassMatcher{
							pos:        position{line: 211, col: 9, offset: 8698},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&choiceExpr{
						pos: position{line: 211, col: 17, offset: 8706},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 211, col: 17, offset: 8706},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 211, col: 26, offset: 8715},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 211, col: 35, offset: 8724},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 211, col: 42, offset: 8731},
								val:        "\n",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 211, col: 49, offset: 8738},
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
			pos:         position{line: 213, col: 1, offset: 8744},
			expr: &zeroOrMoreExpr{
				pos: position{line: 213, col: 19, offset: 8762},
				expr: &charClassMatcher{
					pos:        position{line: 213, col: 19, offset: 8762},
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
			pos:         position{line: 215, col: 1, offset: 8774},
			expr: &oneOrMoreExpr{
				pos: position{line: 215, col: 31, offset: 8804},
				expr: &charClassMatcher{
					pos:        position{line: 215, col: 31, offset: 8804},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 217, col: 1, offset: 8816},
			expr: &notExpr{
				pos: position{line: 217, col: 8, offset: 8823},
				expr: &anyMatcher{
					line: 217, col: 9, offset: 8824,
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
	fieldsSlice = append(fieldsSlice, valsSlice[0].(map[string]string))
	restSlice := toIfaceSlice(valsSlice[1])
	for _, v := range restSlice {
		vSlice := toIfaceSlice(v)
		fieldsSlice = append(fieldsSlice, vSlice[3].(map[string]string))
	}
	return parseCreateTableStmt(tablename, fieldsSlice), nil
}

func (p *parser) callonCreateTableStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCreateTableStmt1(stack["tablename"], stack["fields"])
}

func (c *current) onFieldDef1(name, dataType, notnull interface{}) (interface{}, error) {
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

func (c *current) onOtherT1() (interface{}, error) {
	return map[string]string{
		"type": string(c.text),
	}, nil
}

func (p *parser) callonOtherT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOtherT1()
}

func (c *current) onCommentExtensionStmt1(extension, comment interface{}) (interface{}, error) {
	return parseCommentExtensionStmt(extension.(Identifier), comment.(String)), nil
}

func (p *parser) callonCommentExtensionStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCommentExtensionStmt1(stack["extension"], stack["comment"])
}

func (c *current) onCreateExtensionStmt1(extension, schema interface{}) (interface{}, error) {
	return parseCreateExtensionStmt(extension.(Identifier), schema.(Identifier)), nil
}

func (p *parser) callonCreateExtensionStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCreateExtensionStmt1(stack["extension"], stack["schema"])
}

func (c *current) onSetStmt1(key, values interface{}) (interface{}, error) {
	updateSettings(key.(string), toIfaceSlice(values))
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
