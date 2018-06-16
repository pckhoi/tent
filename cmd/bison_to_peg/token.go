package main

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
)

type RepeatCharacteristic int

const (
	// OneOrNone represent "?" repeater
	OneOrNone RepeatCharacteristic = iota + 1
	// Any represent "*" repeater
	Any
	// OneOrMany represent "+" repeater
	OneOrMany
)

type TokenPointer interface {
	WritePegTo(*bytes.Buffer)
	SetRepeat(RepeatCharacteristic)
	String() string
	GetRepeat() RepeatCharacteristic
}

type SimpleToken struct {
	Name   string
	Repeat RepeatCharacteristic
}

var namemap = map[string]string{}
var seenRepr = []string{}

func indexOfString(slice []string, find string) int {
	for ind, val := range slice {
		if val == find {
			return ind
		}
	}
	return -1
}

func newRepr(name string) string {
	repr := strcase.ToCamel(name)
	if indexOfString(seenRepr, repr) == -1 {
		seenRepr = append(seenRepr, repr)
		return repr
	}
	cnt := 1
	var altrepr string
	for true {
		altrepr = fmt.Sprintf("%s%d", repr, cnt)
		if indexOfString(seenRepr, altrepr) == -1 {
			seenRepr = append(seenRepr, altrepr)
			break
		}
	}
	return altrepr
}

func (token *SimpleToken) WritePegTo(buffer *bytes.Buffer) {
	buffer.WriteString(namemap[token.Name])
	switch token.Repeat {
	case OneOrMany:
		buffer.WriteString("+")
	case OneOrNone:
		buffer.WriteString("?")
	case Any:
		buffer.WriteString("*")
	}
}

func (token *SimpleToken) SetRepeat(repeat RepeatCharacteristic) {
	token.Repeat = repeat
}

func (token *SimpleToken) GetRepeat() RepeatCharacteristic {
	return token.Repeat
}

func (token *SimpleToken) MarkAsRuleName() {
	namemap[token.Name] = newRepr(token.Name)
}

func (token *SimpleToken) String() string {
	return namemap[token.Name]
}

func MakeSimpleToken(name string) *SimpleToken {
	if _, ok := namemap[name]; !ok {
		namemap[name] = name
	}
	token := SimpleToken{
		Name: name,
	}
	return &token
}
