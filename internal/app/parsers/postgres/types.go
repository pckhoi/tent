package postgres

type String string

type Identifier string

type DoubleQuotedString string

type Enum struct {
	Name   Identifier
	Labels []String
}
