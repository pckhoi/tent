package main

import (
	"bytes"
	"strings"
)

type StringToken struct {
	Name        string
	Repeat      RepeatCharacteristic
	Insensitive bool
}

func (token *StringToken) WritePegTo(buffer *bytes.Buffer) {
	buffer.WriteString(`"`)
	name := strings.Replace(token.Name, `\`, `\\`, -1)
	name = strings.Replace(name, `"`, `\"`, -1)
	buffer.WriteString(name)
	buffer.WriteString(`"`)
	if token.Insensitive {
		buffer.WriteString("i")
	}
	switch token.Repeat {
	case OneOrMany:
		buffer.WriteString("+")
	case OneOrNone:
		buffer.WriteString("?")
	case Any:
		buffer.WriteString("*")
	}
}

func (token *StringToken) String() string {
	var buffer bytes.Buffer
	token.WritePegTo(&buffer)
	return buffer.String()
}

func (token *StringToken) SetRepeat(repeat RepeatCharacteristic) {
	token.Repeat = repeat
}

func (token *StringToken) GetRepeat() RepeatCharacteristic {
	return token.Repeat
}

func MakeStringToken(name string, insensitive bool) *StringToken {
	return &StringToken{
		Name:        name,
		Insensitive: insensitive,
	}
}
