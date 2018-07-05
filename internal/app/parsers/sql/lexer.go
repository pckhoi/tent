package sql

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type lState int

const (
	initial lState = iota
	xb
	xc
	xd
	xh
	xe
	xq
	xdolq
	xui
	xuiend
	xus
	xusend
	xeu
)

var space = oneOf(' ', '\t', '\r', '\n', '\f')
var horizSpace = oneOf(' ', '\t', '\f')
var newline = oneOf('\n', '\r')
var nonNewline = anythingBut('\n', '\r')
var comment = concat("--", any(nonNewline))
var whitespace = oneOf(oneOrMore(space), comment)

var specialWhitespace = oneOf(oneOrMore(space), concat(comment, newline))
var horizWhitespace = oneOf(horizSpace, comment)
var whitespaceWithNewline = concat(any(horizWhitespace), newline, any(specialWhitespace))

var quote = '\''
var quotestop = concat('\'', any(whitespace))
var quotecontinue = concat('\'', whitespaceWithNewline, '\'')
var quotefail = concat('\'', any(whitespace), '-')

var xbstart = concat(oneOf('b', 'B'), '\'')
var xbinside = any(anythingBut('\''))

/* Hexadecimal number */
var xhstart = concat(oneOf('x', 'X'), '\'')
var xhinside = any(anythingBut('\''))

/* National character */
var xnstart = concat(any('n', 'N'), '\'')

/* Quoted string that allows backslash escapes */
var xestart = concat(any('e', 'E'), '\'')
var xeinside = oneOrMore(anythingBut('\\', '\''))
var xeescape = concat('\\', anythingBut(charRange('0', '7')))
var xeoctesc = concat('\\', numChar(charRange('0', '7'), 1, 3))
var xehexesc = concat('\\', 'x', numChar(oneOf(charRange('0', '9'), charRange('A', 'F'), charRange('a', 'f')), 1, 2))
var xeunicode = concat('\\', oneOf(
	concat('u', numChar(oneOf(charRange('0', '9'), charRange('A', 'F'), charRange('a', 'f')), 4)),
	concat('U', numChar(oneOf(charRange('0', '9'), charRange('A', 'F'), charRange('a', 'f')), 8)),
))
var xeunicodefail = concat('\\', oneOf(
	concat('u', numChar(oneOf(charRange('0', '9'), charRange('A', 'F'), charRange('a', 'f')), 0, 3)),
	concat('U', numChar(oneOf(charRange('0', '9'), charRange('A', 'F'), charRange('a', 'f')), 0, 7)),
))

/* Extended quote
 * xqdouble implements embedded quote, ''''
 */
var xqstart = '\''
var xqdouble = "''"
var xqinside = oneOrMore(anythingBut('\''))

var dolqStart = oneOf(charRange('A', 'Z'), charRange('a', 'z'), charRange(128, 255), '_')
var dolqCont = oneOf(charRange('A', 'Z'), charRange('a', 'z'), charRange(128, 255), '_', charRange('0', '9'))
var dolqdelim = concat('$', zeroOrOne(concat(dolqStart, any(dolqCont))), '$')
var dolqfailed = concat('$', dolqStart, any(dolqCont))
var dolqinside = oneOrMore(anythingBut('$'))

/* Double quote
 * Allows embedded spaces and other special characters into identifiers.
 */
var dquote = '"'
var xdstart = '"'
var xdstop = '"'
var xddouble = `""`
var xdinside = oneOrMore(anythingBut('"'))

/* Unicode escapes */
var uescapePrefix = concat(
	oneOf('u', 'U'),
	oneOf('e', 'E'),
	oneOf('s', 'S'),
	oneOf('c', 'C'),
	oneOf('a', 'A'),
	oneOf('p', 'P'),
	oneOf('e', 'E'),
)
var uescape = concat(
	uescapePrefix,
	any(whitespace),
	'\'',
	anythingBut('\''),
	'\'',
)

/* error rule to avoid backup */
var uescapefail = oneOf(
	concat(uescapePrefix, any(whitespace), '-'),
	concat(uescapePrefix, any(whitespace), '\'', anythingBut('\'')),
	concat(uescapePrefix, any(whitespace), '\''),
	concat(uescapePrefix, any(whitespace)),
	concat(oneOf('u', 'U'), oneOf('e', 'E'), oneOf('s', 'S'), oneOf('c', 'C'), oneOf('a', 'A'), oneOf('p', 'P')),
	concat(oneOf('u', 'U'), oneOf('e', 'E'), oneOf('s', 'S'), oneOf('c', 'C'), oneOf('a', 'A')),
	concat(oneOf('u', 'U'), oneOf('e', 'E'), oneOf('s', 'S'), oneOf('c', 'C')),
	concat(oneOf('u', 'U'), oneOf('e', 'E'), oneOf('s', 'S')),
	concat(oneOf('u', 'U'), oneOf('e', 'E'), oneOf('u', 'U')),
)

/* Quoted identifier with Unicode escapes */
var xuistart = concat(oneOf('u', 'U'), dquote)

/* Quoted string with Unicode escapes */
var xusstart = concat(oneOf('u', 'U'), quote)

/* Optional UESCAPE after a quoted string or identifier with Unicode escapes. */
var xustop1 = zeroOrOne(uescapefail)
var xustop2 = uescape

/* error rule to avoid backup */
var xufailed = concat(oneOf('u', 'U'), '&')

var digit = charRange('0', '9')
var identStart = oneOf(charRange('A', 'Z'), charRange('a', 'z'), charRange(128, 255), '_')
var identCont = oneOf(charRange('A', 'Z'), charRange('a', 'z'), charRange(128, 255), '_', charRange('0', '9'), '$')

var identifier = concat(identStart, any(identCont))

/* Assorted special-case operators and operator-like tokens */
var typecast = "::"
var dotDot = ".."
var colonEquals = ":="
var equalsGreater = "=>"
var lessEquals = "<="
var greaterEquals = ">="
var lessGreater = "<>"
var notEquals = "!="

var self = oneOf(',', '(', ')', '[', ']', '.', ';', ':', '+', '-', '*', '/', '%', '^', '<', '>', '=')
var opChars = oneOf('~', '!', '@', '#', '^', '&', '|', '`', '?', '+', '-', '*', '/', '%', '<', '>', '=')
var operator = oneOrMore(opChars)

var xcstart = concat('/', '*', any(opChars))
var xcstop = concat(oneOrMore('*'), '/')
var xcinside = oneOrMore(anythingBut('*', '/'))

var integer = oneOrMore(digit)
var decimal = oneOf(concat(any(digit), '.', oneOrMore(digit)), concat(oneOrMore(digit), '.', any(digit)))
var decimalfail = concat(oneOrMore(digit), '.', '.')
var realsuccess = concat(oneOf(integer, decimal), oneOf('E', 'e'), zeroOrOne(oneOf('-', '+')), oneOrMore(digit))
var realfail1 = concat(oneOf(integer, decimal), oneOf('E', 'e'))
var realfail2 = concat(oneOf(integer, decimal), oneOf('E', 'e'), oneOf('-', '+'))

var param = concat('$', integer)

var other = nonNewline

type lexer struct {
	data        []byte
	litslice    []byte
	dolqstart   []byte
	rules       []lexerRule
	state       lState
	terminated  bool
	start       int
	tokenLength int
	xcdepth     int
}

func newLexer(data []byte) lexer {
	return lexer{
		data:  data,
		state: initial,
		start: 0,
		rules: []lexerRule{
			makeRule(
				[]lState{initial},
				whitespace,
				nil,
			),
			makeRule(
				[]lState{initial},
				xcstart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.xcdepth = 0
					lxr.tokenLength = 2
					lxr.begin(xc)
					return 0
				},
			),
			makeRule(
				[]lState{xc},
				xcstart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.xcdepth++
					return 0
				},
			),
			makeRule(
				[]lState{xc},
				xcstop,
				func(yylval *sqlSymType, lxr lexer) int {
					if lxr.xcdepth <= 0 {
						lxr.begin(initial)
					} else {
						lxr.xcdepth--
					}
					return 0
				},
			),
			makeRule(
				[]lState{xc},
				xcinside,
				nil,
			),
			makeRule(
				[]lState{xc},
				opChars,
				nil,
			),
			makeRule(
				[]lState{xc},
				oneOrMore('*'),
				nil,
			),
			makeRule(
				[]lState{xc},
				EOF,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.Error("Unterminated /* comment")
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				xbstart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(xb)
					lxr.startlit()
					return 0
				},
			),
			makeRule(
				[]lState{xb},
				oneOf(quotestop, quotefail),
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 1
					lxr.begin(initial)
					lxr.addToLit()
					yylval.str = lxr.litslice
					return BCONST
				},
			),
			makeRule(
				[]lState{xh},
				xhinside,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xb},
				xbinside,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xh, xb, xq, xe, xus},
				quotecontinue,
				nil,
			),
			makeRule(
				[]lState{xb},
				EOF,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.Error("Unterminated bit string literal")
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				xhstart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(xh)
					lxr.startlit()
					return 0
				},
			),
			makeRule(
				[]lState{xh},
				oneOf(quotestop, quotefail),
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 1
					lxr.begin(initial)
					lxr.addToLit()
					yylval.str = lxr.litslice
					return XCONST
				},
			),
			makeRule(
				[]lState{xh, xq, xe, xus},
				EOF,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.Error("Unterminated hexadecimal string literal")
					return 0
				},
			),
			makeRule(
				[]lState{xq, xe, xus},
				EOF,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.Error("Unterminated quoted string")
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				xnstart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 1
					kw := scanKeywordLookup("nchar")
					if kw != nil {
						yylval.keyword = kw.name
						return kw.value
					}
					yylval.str = string(lxr.currentToken())
					return IDENT
				},
			),
			makeRule(
				[]lState{initial},
				xqstart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(xq)
					lxr.startlit()
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				xestart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(xe)
					lxr.startlit()
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				xestart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(xus)
					lxr.startlit()
					return 0
				},
			),
			makeRule(
				[]lState{xq, xe},
				oneOf(quotestop, quotefail),
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 1
					lxr.begin(initial)
					lxr.addToLit()
					yylval.str = lxr.litslice
					return SCONST
				},
			),
			makeRule(
				[]lState{xus},
				oneOf(quotestop, quotefail),
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 1
					lxr.begin(xusend)
					return 0
				},
			),
			makeRule(
				[]lState{xusend},
				whitespace,
				nil,
			),
			makeRule(
				[]lState{xusend},
				oneOf(EOF, other, xustop1),
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 0
					lxr.begin(initial)
					yylval.str = lxr.litslice
					return SCONST
				},
			),
			makeRule(
				[]lState{xusend},
				xustop2,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(initial)
					lxr.addToLit()
					yylval.str = lxr.litslice
					return SCONST
				},
			),
			makeRule(
				[]lState{xq, xe, xus},
				xqdouble,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xq, xus},
				xqinside,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xe},
				xeinside,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xe},
				xeunicode,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xe},
				xeunicodefail,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.Error("Invalid Unicode escape")
					return 0
				},
			),
			makeRule(
				[]lState{xe},
				xeescape,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xe},
				xeoctesc,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xe},
				xehexesc,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xe},
				other,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				dolqdelim,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.dolqstart = lxr.currentToken()
					lxr.begin(xdolq)
					lxr.startlit()
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				dolqfailed,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 1
					return lxr.currentToken()
				},
			),
			makeRule(
				[]lState{xdolq},
				dolqdelim,
				func(yylval *sqlSymType, lxr lexer) int {
					if lxr.dolqstart == lxr.currentToken() {
						lxr.begin(initial)
						lxr.addToLit()
						yylval.str = lxr.litslice
						return SCONST
					} else {
						lxr.tokenLength--
						lxr.addToLit()
						return 0
					}
				},
			),
			makeRule(
				[]lState{xdolq},
				dolqinside,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xdolq},
				dolqfailed,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xdolq},
				other,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xdolq},
				EOF,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.Error("Unterminated dollar-quoted string")
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				xdstart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(xd)
					lxr.startlit()
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				xuistart,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(xui)
					lxr.startlit()
					return 0
				},
			),
			makeRule(
				[]lState{xd},
				xdstop,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(initial)
					lxr.addToLit()
					yylval.str = lxr.litslice
					return IDENT
				},
			),
			makeRule(
				[]lState{xui},
				dquote,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 1
					lxr.begin(xuiend)
					return 0
				},
			),
			makeRule(
				[]lState{xuiend},
				whitespace,
				nil,
			),
			makeRule(
				[]lState{xuiend},
				oneOf(EOF, other, xustop1),
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 0
					lxr.begin(initial)
					yylval.str = lxr.litslice
					return IDENT
				},
			),
			makeRule(
				[]lState{xuiend},
				oneOf(xustop2),
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.begin(initial)
					lxr.addToLit()
					yylval.str = lxr.litslice
					return IDENT
				},
			),
			makeRule(
				[]lState{xd, xui},
				oneOf(xddouble),
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xd, xui},
				oneOf(xdinside),
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.addToLit()
					return 0
				},
			),
			makeRule(
				[]lState{xd, xui},
				EOF,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.Error("Unterminated quoted identifier")
					return 0
				},
			),
			makeRule(
				[]lState{initial},
				xufailed,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength = 1
					yylval.str = lxr.currentToken()
					return IDENT
				},
			),
			makeRule(
				[]lState{initial},
				typecast,
				func(yylval *sqlSymType, lxr lexer) int {
					return TYPECAST
				},
			),
			makeRule(
				[]lState{initial},
				dotDot,
				func(yylval *sqlSymType, lxr lexer) int {
					return DOT_DOT
				},
			),
			makeRule(
				[]lState{initial},
				colonEquals,
				func(yylval *sqlSymType, lxr lexer) int {
					return COLON_EQUALS
				},
			),
			makeRule(
				[]lState{initial},
				equalsGreater,
				func(yylval *sqlSymType, lxr lexer) int {
					return EQUALS_GREATER
				},
			),
			makeRule(
				[]lState{initial},
				lessEquals,
				func(yylval *sqlSymType, lxr lexer) int {
					return LESS_EQUALS
				},
			),
			makeRule(
				[]lState{initial},
				greaterEquals,
				func(yylval *sqlSymType, lxr lexer) int {
					return GREATER_EQUALS
				},
			),
			makeRule(
				[]lState{initial},
				lessGreater,
				func(yylval *sqlSymType, lxr lexer) int {
					return NOT_EQUALS
				},
			),
			makeRule(
				[]lState{initial},
				notEquals,
				func(yylval *sqlSymType, lxr lexer) int {
					return NOT_EQUALS
				},
			),
			makeRule(
				[]lState{initial},
				self,
				func(yylval *sqlSymType, lxr lexer) int {
					return lxr.currentToken()
				},
			),
			makeRule(
				[]lState{initial},
				operator,
				func(yylval *sqlSymType, lxr lexer) int {
					token := string(lxr.currentToken())
					slashstarInd := strings.Index(token, "/*")
					dashdashInd := strings.Index(token, "--")

					if slashstarInd > -1 && dashdashInd > -1 {
						if slashstarInd > dashdashInd {
							slashstarInd = dashdashInd
						}
					} else if slashstarInd == -1 {
						slashstarInd = dashdashInd
					}
					if slashstarInd > -1 {
						lxr.tokenLength = slashstarInd
					}
					yylval.str = lxr.currentToken()
					return Op
				},
			),
			makeRule(
				[]lState{initial},
				param,
				func(yylval *sqlSymType, lxr lexer) int {
					yylval.str = lxr.currentToken()
					return PARAM
				},
			),
			makeRule(
				[]lState{initial},
				integer,
				func(yylval *sqlSymType, lxr lexer) int {
					yylval.str = lxr.currentToken()
					return ICONST
				},
			),
			makeRule(
				[]lState{initial},
				decimal,
				func(yylval *sqlSymType, lxr lexer) int {
					yylval.str = lxr.currentToken()
					return FCONST
				},
			),
			makeRule(
				[]lState{initial},
				decimalfail,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength -= 2
					yylval.str = lxr.currentToken()
					return ICONST
				},
			),
			makeRule(
				[]lState{initial},
				realsuccess,
				func(yylval *sqlSymType, lxr lexer) int {
					yylval.str = lxr.currentToken()
					return FCONST
				},
			),
			makeRule(
				[]lState{initial},
				realfail1,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength -= 1
					yylval.str = lxr.currentToken()
					return FCONST
				},
			),
			makeRule(
				[]lState{initial},
				realfail2,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.tokenLength -= 2
					yylval.str = lxr.currentToken()
					return FCONST
				},
			),
			makeRule(
				[]lState{initial},
				identifier,
				func(yylval *sqlSymType, lxr lexer) int {
					kw := scanKeywordLookup(lxr.currentToken())
					if kw != nil {
						yylval.keyword = kw.name
						return kw.value
					}
					yylval.str = lxr.currentToken()
					return IDENT
				},
			),
			makeRule(
				[]lState{initial},
				other,
				func(yylval *sqlSymType, lxr lexer) int {
					return lxr.currentToken()
				},
			),
			makeRule(
				[]lState{initial},
				EOF,
				func(yylval *sqlSymType, lxr lexer) int {
					lxr.terminate()
					return 0
				},
			),
		},
	}
}

func (lxr *lexer) shift() {
	lxr.start = lxr.start + lxr.tokenLength
	lxr.tokenLength = 0
}

func (lxr *lexer) terminate() {
	lxr.terminated = true
}

func (lxr *lexer) begin(state lState) {
	lxr.state = state
}

func (lxr *lexer) startlit() {
	lxr.litslice = lxr.data[lxr.start : lxr.start+lxr.tokenLength]
}

func (lxr *lexer) addToLit() {
	lxr.litslice = append(lxr.litslice, lxr.data[lxr.start:lxr.start+lxr.tokenLength])
}

func (lxr *lexer) currentToken() []byte {
	return lxr.data[lxr.start : lxr.start+lxr.tokenLength]
}

func (lxr *lexer) printErrorWithDebugInfo() {
	endInd := lxr.start + 20
	if endInd > len(lxr.data) {
		endInd = len(lxr.data)
	}
	lxr.Error(fmt.Sprintf(
		"Can't match at position %d: \"%s\"", lxr.start,
		lxr.data[lxr.start:endInd]))
}

func (lxr *lexer) Error(str string) {
	log.Printf("Lexer error: %s", str)
}

func (lxr *lexer) Lex(yylval *sqlSymType) int {
	for {
		if lxr.terminated {
			break
		}
		var theRule *lexerRule
		var longestLength int
		lxr.shift()
		for _, rule := range lxr.rules {
			if rule.workInState(lxr.state) {
				ok, matchLength := rule.match(*lxr)
				if ok && matchLength > longestLength {
					longestLength = matchLength
					theRule = &rule
				}
			}
		}

		if lxr.start != len(lxr.data)-1 && (theRule == nil || longestLength == 0) {
			lxr.printErrorWithDebugInfo()
			break
		} else {
			lxr.tokenLength = longestLength
			if theRule.action != nil {
				result = theRule.action(yylval, lxr)
				if result != 0 {
					return result
				}
			}
		}
	}
	return 0
}
