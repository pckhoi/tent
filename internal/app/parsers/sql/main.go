//go:generate goyacc -o sql.go -p "sql" sql.y

package sql

func Parse(data []byte) interface{} {
	lex := newLexer(data)
	parser := sqlNewParser()
	parser.Parse(&lex)
	return parser.(*sqlParserImpl).lval
}
