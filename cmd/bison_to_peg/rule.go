package main

import (
	"bytes"
	"errors"
	"github.com/iancoleman/strcase"
	"regexp"
)

type Rule struct {
	Name     string
	Subrules []Subrule
	Optional bool
}

func (r Rule) WritePegTo(buffer *bytes.Buffer, context Context) {
	buffer.WriteString(context.TokenMap[r.Name].DisplayString)
	buffer.WriteString("\n    <- ")
	if r.Optional {
		buffer.WriteString("(")
	}
	context.SiblingsCount = len(r.Subrules)
	for ind, subrule := range r.Subrules {
		subrule.WritePegTo(buffer, context)
		if ind < len(r.Subrules)-1 {
			buffer.WriteString("\n    / ")
		}
	}
	if r.Optional {
		if len(r.Subrules) > 1 {
			buffer.WriteString("\n    ")
		}
		buffer.WriteString(")?")
	}
}

var nameRe = regexp.MustCompile(`(\w+) *:\s+`)
var pipeRe = regexp.MustCompile(`\s*\|\s*`)

func MakeRule(src string, tokenMap map[string]Token) (Rule, error) {
	var rule Rule
	indicies := nameRe.FindStringSubmatchIndex(src)
	if indicies == nil {
		return rule, errors.New("Cant find rule name")
	}

	ruleName := src[indicies[2]:indicies[3]]
	tokenMap[ruleName] = Token{
		IsRuleName:    true,
		DisplayString: strcase.ToCamel(ruleName),
	}

	src = src[indicies[1]:]
	subruleStrings := pipeRe.Split(src, -1)
	subrules := []Subrule{}
	optional := false
	for _, v := range subruleStrings {
		if v == "" {
			optional = true
			continue
		}
		subrules = append(subrules, MakeSubrule(v))
	}
	rule = Rule{
		Name:     ruleName,
		Subrules: subrules,
		Optional: optional,
	}
	return rule, nil
}
