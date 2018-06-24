package main

import (
	"bytes"
)

type LiteralToken struct {
	Literal      string
	Repeat       RepeatCharacteristic
	VariableName string
}

func (token *LiteralToken) WritePegTo(buffer *bytes.Buffer) {
	buffer.WriteString(token.Literal)
	switch token.Repeat {
	case OneOrMany:
		buffer.WriteString("+")
	case OneOrNone:
		buffer.WriteString("?")
	case Any:
		buffer.WriteString("*")
	}
}

func (token *LiteralToken) String() string {
	var buffer bytes.Buffer
	token.WritePegTo(&buffer)
	return buffer.String()
}

func (token *LiteralToken) GetVariableName() string {
	return token.VariableName
}

func (token *LiteralToken) SetRepeat(repeat RepeatCharacteristic) {
	token.Repeat = repeat
}

func (token *LiteralToken) GetRepeat() RepeatCharacteristic {
	return token.Repeat
}
