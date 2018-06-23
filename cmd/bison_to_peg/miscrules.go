package main

func miscellaneousRules() []Rule {
	return []Rule{
		Rule{
			Name: *MakeReferToken("ICONST", true, 0),
			Expression: &LiteralToken{
				Literal: "[0-9]",
				Repeat:  OneOrMany,
			},
		},
		Rule{
			Name: *MakeReferToken("Decimal", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&TokenGroup{
						Tokens: []TokenPointer{
							&LiteralToken{
								Literal: "[0-9]",
								Repeat:  Any,
							},
							&StringToken{
								Name: ".",
							},
							&LiteralToken{
								Literal: "[0-9]",
								Repeat:  OneOrMany,
							},
						},
						Type: Sequence,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							&LiteralToken{
								Literal: "[0-9]",
								Repeat:  OneOrMany,
							},
							&StringToken{
								Name: ".",
							},
							&LiteralToken{
								Literal: "[0-9]",
								Repeat:  Any,
							},
						},
						Type: Sequence,
					},
				},
				IsRoot: true,
				Type:   Choice,
			},
		},
		Rule{
			Name: *MakeReferToken("Real", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&TokenGroup{
						Tokens: []TokenPointer{
							MakeReferToken("Decimal", false, 0),
							MakeReferToken("ICONST", false, 0),
						},
						Type: Choice,
					},
					&StringToken{
						Name:        "E",
						Insensitive: true,
					},
					&LiteralToken{
						Literal: "[+-]",
						Repeat:  OneOrNone,
					},
					&LiteralToken{
						Literal: "[0-9]",
						Repeat:  OneOrMany,
					},
				},
				Type:   Sequence,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("FCONST", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					MakeReferToken("Real", false, 0),
					MakeReferToken("Decimal", false, 0),
				},
				Type:   Choice,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("LESS_EQUALS", true, 0),
			Expression: &StringToken{
				Name: "<=",
			},
		},
		Rule{
			Name: *MakeReferToken("COLON_EQUALS", true, 0),
			Expression: &StringToken{
				Name: ":=",
			},
		},
		Rule{
			Name: *MakeReferToken("GREATER_EQUALS", true, 0),
			Expression: &StringToken{
				Name: ">=",
			},
		},
		Rule{
			Name: *MakeReferToken("EQUALS_GREATER", true, 0),
			Expression: &StringToken{
				Name: "=>",
			},
		},
		Rule{
			Name: *MakeReferToken("NOT_EQUALS", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: "<>",
					}, &StringToken{
						Name: "!=",
					},
				},
				Type:   Choice,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("Space", true, 0),
			Expression: &LiteralToken{
				Literal: `[ \t\n\r\f]`,
			},
		},
		Rule{
			Name: *MakeReferToken("NonNewLine", true, 0),
			Expression: &LiteralToken{
				Literal: `[^\n\r]`,
			},
		},
		Rule{
			Name: *MakeReferToken("NewLine", true, 0),
			Expression: &LiteralToken{
				Literal: `[\n\r]`,
			},
		},
		Rule{
			Name: *MakeReferToken("HorizSpace", true, 0),
			Expression: &LiteralToken{
				Literal: `[ \t\f]`,
			},
		},
		Rule{
			Name: *MakeReferToken("SpecialWhiteSpace", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					MakeReferToken("Space", false, OneOrMany),
					&TokenGroup{
						Tokens: []TokenPointer{
							MakeReferToken("Comment", false, 0),
							MakeReferToken("NewLine", false, 0),
						},
						Type: Sequence,
					},
				},
				Type:   Choice,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("HorizWhitespace", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					MakeReferToken("HorizSpace", false, OneOrMany),
					MakeReferToken("Comment", false, 0),
				},
				Type:   Choice,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("WhitespaceWithNewline", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					MakeReferToken("HorizWhitespace", false, Any),
					MakeReferToken("NewLine", false, 0),
					MakeReferToken("SpecialWhiteSpace", false, Any),
				},
				Type:   Sequence,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("Comment", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: "--",
					},
					MakeReferToken("NonNewLine", false, Any),
				},
				Type:   Sequence,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken(`_`, true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					MakeReferToken("Space", false, OneOrMany),
					MakeReferToken("Comment", false, 0),
				},
				IsRoot: true,
				Type:   Choice,
				Repeat: Any,
			},
			ReturnsNil: true,
		},
		Rule{
			Name: *MakeReferToken("PARAM", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: "$",
					},
					MakeReferToken("ICONST", false, 0),
				},
				IsRoot: true,
				Type:   Sequence,
			},
		},
		Rule{
			Name: *MakeReferToken("Identifier", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&LiteralToken{
						Literal: `[A-Za-z\pL_]`,
					},
					&LiteralToken{
						Literal: `[A-Za-z\pL_0-9$]`,
						Repeat:  Any,
					},
				},
				IsRoot: true,
				Type:   Sequence,
			},
			ReturnsString: true,
		},
		Rule{
			Name: *MakeReferToken("DblQuoIdentifier", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: `"`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							&StringToken{
								Name: `""`,
							},
							&LiteralToken{
								Literal: `[^"]`,
							},
						},
						Type:   Choice,
						Repeat: OneOrMany,
					},
					&StringToken{
						Name: `"`,
					},
				},
				IsRoot: true,
				Type:   Sequence,
			},
			ReturnsString: true,
		},
		Rule{
			Name: *MakeReferToken("QuoUniIdentifier", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name:        `U`,
						Insensitive: true,
					},
					&StringToken{
						Name: `&"`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							&StringToken{
								Name: `""`,
							},
							&LiteralToken{
								Literal: `[^"]`,
							},
						},
						Type:   Choice,
						Repeat: OneOrMany,
					},
					&StringToken{
						Name: `"`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							MakeReferToken("_", false, 0),
							&StringToken{
								Name:        `UESCAPE`,
								Insensitive: true,
							},
							MakeReferToken("_", false, 0),
							&StringToken{
								Name: `'`,
							},
							&LiteralToken{
								Literal: `[^']`,
							},
							&StringToken{
								Name: `'`,
							},
						},
						Type:   Sequence,
						Repeat: OneOrNone,
					},
				},
				IsRoot: true,
				Type:   Sequence,
			},
			ReturnsString: true,
		},
		Rule{
			Name: *MakeReferToken("IDENT", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					MakeReferToken("Identifier", false, 0),
					MakeReferToken("DblQuoIdentifier", false, 0),
					MakeReferToken("QuoUniIdentifier", false, 0),
				},
				IsRoot: true,
				Type:   Choice,
			},
		},
		Rule{
			Name: *MakeReferToken("BCONST", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name:        "B",
						Insensitive: true,
					},
					&StringToken{
						Name: `'`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							MakeReferToken("QuoteContinue", false, 0),
							&LiteralToken{
								Literal: `[^']`,
							},
						},
						Type:   Choice,
						Repeat: Any,
					},
					&StringToken{
						Name: `'`,
					},
				},
				IsRoot: true,
				Type:   Sequence,
			},
		},
		Rule{
			Name: *MakeReferToken("QuoteContinue", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: `'`,
					},
					MakeReferToken("WhitespaceWithNewline", false, 0),
					&StringToken{
						Name: `'`,
					},
				},
				Type:   Sequence,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("HexEscape", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: `\x`,
					},
					&LiteralToken{
						Literal: "[0-9A-Fa-f]",
					},
					&LiteralToken{
						Literal: "[0-9A-Fa-f]",
						Repeat:  OneOrNone,
					},
				},
				Type:   Sequence,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("OctalEscape", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: `\`,
					},
					&LiteralToken{
						Literal: "[0-7]",
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							&LiteralToken{
								Literal: "[0-7]",
							},
							&LiteralToken{
								Literal: "[0-7]",
								Repeat:  OneOrNone,
							},
						},
						Type:   Sequence,
						Repeat: OneOrNone,
					},
				},
				Type:   Sequence,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("GenericEscape", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: `\`,
					},
					&LiteralToken{
						Literal: "[^0-7]",
					},
				},
				Type:   Sequence,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("UnicodeEscape", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: `\`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							&TokenGroup{
								Tokens: []TokenPointer{
									&StringToken{
										Name: `u`,
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
								},
								Type: Sequence,
							},
							&TokenGroup{
								Tokens: []TokenPointer{
									&StringToken{
										Name: `U`,
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
									&LiteralToken{
										Literal: "[0-9A-Fa-f]",
									},
								},
								Type: Sequence,
							},
						},
						Type: Choice,
					},
				},
				Type:   Sequence,
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("QuoEscString", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name:        "E",
						Insensitive: true,
					},
					&StringToken{
						Name: `'`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							MakeReferToken("UnicodeEscape", false, 0),
							MakeReferToken("HexEscape", false, 0),
							MakeReferToken("OctalEscape", false, 0),
							MakeReferToken("GenericEscape", false, 0),
							MakeReferToken("QuoteContinue", false, 0),
							&StringToken{
								Name: `''`,
							},
							&LiteralToken{
								Literal: `[^']`,
							},
						},
						Type:   Choice,
						Repeat: OneOrMany,
					},
					&StringToken{
						Name: `'`,
					},
				},
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("QuoUniString", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name:        "U",
						Insensitive: true,
					},
					&StringToken{
						Name: `&'`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							MakeReferToken("QuoteContinue", false, 0),
							&StringToken{
								Name: `''`,
							},
							&LiteralToken{
								Literal: `[^']`,
							},
						},
						Type:   Choice,
						Repeat: OneOrMany,
					},
					&StringToken{
						Name: `'`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							MakeReferToken("_", false, 0),
							&StringToken{
								Name:        `UESCAPE`,
								Insensitive: true,
							},
							MakeReferToken("_", false, 0),
							&StringToken{
								Name: `'`,
							},
							&LiteralToken{
								Literal: `[^']`,
							},
							&StringToken{
								Name: `'`,
							},
						},
						Type:   Sequence,
						Repeat: OneOrNone,
					},
				},
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("QuoString", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name: `'`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							MakeReferToken("QuoteContinue", false, 0),
							&StringToken{
								Name: `''`,
							},
							&LiteralToken{
								Literal: `[^']`,
							},
						},
						Type:   Choice,
						Repeat: OneOrMany,
					},
					&StringToken{
						Name: `'`,
					},
				},
				IsRoot: true,
			},
		},
		Rule{
			Name: *MakeReferToken("SCONST", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					MakeReferToken("QuoUniString", false, 0),
					MakeReferToken("QuoEscString", false, 0),
					MakeReferToken("QuoString", false, 0),
				},
				IsRoot: true,
				Type:   Choice,
			},
		},
		Rule{
			Name: *MakeReferToken("TYPECAST", true, 0),
			Expression: &StringToken{
				Name: "::",
			},
		},
		Rule{
			Name: *MakeReferToken("Op", true, 0),
			Expression: &LiteralToken{
				Literal: "[~!@#^&|`?+-*/%<>=]",
				Repeat:  OneOrMany,
			},
		},
		Rule{
			Name: *MakeReferToken("XCONST", true, 0),
			Expression: &TokenGroup{
				Tokens: []TokenPointer{
					&StringToken{
						Name:        "X",
						Insensitive: true,
					},
					&StringToken{
						Name: `'`,
					},
					&TokenGroup{
						Tokens: []TokenPointer{
							MakeReferToken("QuoteContinue", false, 0),
							&LiteralToken{
								Literal: `[^']`,
							},
						},
						Type:   Choice,
						Repeat: Any,
					},
					&StringToken{
						Name: `'`,
					},
				},
				IsRoot: true,
				Type:   Sequence,
			},
		},
	}
}
