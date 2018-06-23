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

type ReferToken struct {
	Name         string
	VariableName string
	Repeat       RepeatCharacteristic
}

var namemap = map[string]string{}
var seenRepr = []string{}
var rulenames = []string{}

func indexOfString(slice []string, find string) int {
	for ind, val := range slice {
		if val == find {
			return ind
		}
	}
	return -1
}

func newRepr(name string) string {
	if name == "_" {
		return name
	}
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

func (token *ReferToken) WritePegTo(buffer *bytes.Buffer) {
	if token.VariableName != "" {
		buffer.WriteString(token.VariableName)
		buffer.WriteString(":")
	}
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

func (token *ReferToken) SetRepeat(repeat RepeatCharacteristic) {
	token.Repeat = repeat
}

func (token *ReferToken) GetRepeat() RepeatCharacteristic {
	return token.Repeat
}

func (token *ReferToken) MarkAsRuleName() {
	namemap[token.Name] = newRepr(token.Name)
	rulenames = append(rulenames, token.Name)
}

func (token *ReferToken) String() string {
	var buffer bytes.Buffer
	token.WritePegTo(&buffer)
	return buffer.String()
}

func MakeSimpleToken(name string) TokenPointer {
	if name[0] == '\'' && name[len(name)-1] == '\'' {
		return MakeStringToken(name[1:len(name)-1], false)
	}
	return MakeReferToken(name, false, 0)
}

func MakeReferToken(name string, root bool, repeat RepeatCharacteristic) *ReferToken {
	if _, ok := namemap[name]; !ok {
		namemap[name] = name
	}
	token := ReferToken{
		Name:   name,
		Repeat: repeat,
	}
	if root {
		token.MarkAsRuleName()
	}
	return &token
}
