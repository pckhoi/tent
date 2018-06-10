package main

import (
	"bytes"
	"regexp"
)

type Subrule struct {
	Words []string
}

func (s Subrule) WritePegTo(buffer *bytes.Buffer, context Context) {
	bracket := len(s.Words) > 1 && context.SiblingsCount > 1
	if bracket {
		buffer.WriteString("(")
	}
	for ind, word := range s.Words {
		if token, ok := context.TokenMap[word]; ok {
			buffer.WriteString(token.DisplayString)
		} else {
			buffer.WriteString(word)
		}
		if ind < len(s.Words)-1 {
			buffer.WriteString(" _ ")
		}
	}
	if bracket {
		buffer.WriteString(")")
	}
}

var spaceRe = regexp.MustCompile(`\s+`)

func MakeSubrule(src string) Subrule {
	return Subrule{
		Words: spaceRe.Split(src, -1),
	}
}
