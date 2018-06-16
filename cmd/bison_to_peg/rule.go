package main

import (
	"bytes"
	"errors"
	// "log"
	"regexp"
	"strings"
)

type Rule struct {
	Name               SimpleToken
	Expression         TokenPointer
	SelfReferencing    bool
	SelfRefAtBegin     bool
	SelfRefAtBeginOnly bool
	SelfRefAtEndOnly   bool
}

func (r Rule) WritePegTo(buffer *bytes.Buffer) {
	buffer.WriteString(r.Name.String())
	buffer.WriteString("\n    = ")
	r.Expression.WritePegTo(buffer)
}

func (r Rule) String() string {
	var buffer bytes.Buffer
	r.WritePegTo(&buffer)
	return buffer.String()
}

func (r *Rule) Inspect() {
	if val, ok := r.Expression.(*TokenGroup); ok {
		selfRef, selfRefAtBegin, selfRefAtBeginOnly, selfRefAtEndOnly := val.DetectSelfReferencing(r.Name)
		r.SelfReferencing = selfRef
		r.SelfRefAtBegin = selfRefAtBegin
		r.SelfRefAtBeginOnly = selfRefAtBeginOnly
		r.SelfRefAtEndOnly = selfRefAtEndOnly
	}
}

func (r *Rule) FixSelfRefAtBeginOnly(reversed bool) {
	if val, ok := r.Expression.(*TokenGroup); ok {
		val.FixSelfRefAtBeginOnly(r.Name, reversed)
	}
}

func (r *Rule) ReorderSubrules() {
	if val, ok := r.Expression.(*TokenGroup); ok {
		val.ReorderTokens()
	}
}

func (r Rule) SplitSelfRef() []Rule {
	tokenGroup, ok := r.Expression.(*TokenGroup)
	if !ok {
		return []Rule{r}
	}
	var normalGroups []TokenPointer
	var selfRefGroups []TokenPointer
	newRuleName := MakeSimpleToken(r.Name.Name + "1")
	newRuleName.MarkAsRuleName()
	for _, token := range tokenGroup.Tokens {
		group, ok := token.(*TokenGroup)
		if !ok {
			normalGroups = append(normalGroups, token)
			continue
		}
		ind := group.IndexOfToken(r.Name)
		if ind == 0 {
			group.ReplaceToken(r.Name, *newRuleName)
			selfRefGroups = append(selfRefGroups, group)
		} else {
			normalGroups = append(normalGroups, token)
		}
	}

	heirTokenGroup := TokenGroup{
		Tokens: append(
			selfRefGroups,
			newRuleName,
		),
		IsRoot: true,
		Repeat: tokenGroup.Repeat,
		Type:   Choice,
	}

	childTokenGroup := TokenGroup{
		Tokens: normalGroups,
		IsRoot: true,
		Type:   Choice,
	}

	return []Rule{
		Rule{
			Name:       r.Name,
			Expression: &heirTokenGroup,
		},
		Rule{
			Name:       *newRuleName,
			Expression: &childTokenGroup,
		},
	}
}

func (r *Rule) Simplify() {
	val, ok := r.Expression.(*TokenGroup)
	if !ok {
		return
	}
	r.Expression = val.Simplify()
}

var nameRe = regexp.MustCompile(`(\w+) *:\s+`)
var pipeRe = regexp.MustCompile(`\s*\|\s*`)

func MakeRule(src string) (Rule, error) {
	var rule Rule
	indicies := nameRe.FindStringSubmatchIndex(src)
	if indicies == nil {
		return rule, errors.New("Cant find rule name")
	}

	ruleName := src[indicies[2]:indicies[3]]
	nameToken := SimpleToken{
		Name: ruleName,
	}
	nameToken.MarkAsRuleName()

	var expression TokenPointer
	src = src[indicies[1]:]
	subruleStrings := pipeRe.Split(src, -1)
	tokens := []TokenPointer{}
	var exprRepeat RepeatCharacteristic
	for _, v := range subruleStrings {
		if v == "" {
			exprRepeat = OneOrNone
			continue
		}
		tokens = append(tokens, MakeToken(v))
	}
	if len(tokens) == 1 {
		expression = tokens[0]
	} else {
		expression = &TokenGroup{
			Tokens: tokens[:],
			IsRoot: true,
			Type:   Choice,
		}
	}
	expression.SetRepeat(exprRepeat)

	if val, ok := expression.(*TokenGroup); ok {
		val.InsertSpace()
	}

	rule = Rule{
		Name:       nameToken,
		Expression: expression,
	}
	return rule, nil
}

func MakeToken(src string) TokenPointer {
	if strings.Index(src, " ") != -1 {
		return MakeTokenGroup(src)
	}
	return MakeSimpleToken(src)
}
