package main

import (
	"bytes"
	"errors"
	// "log"
	"regexp"
	"strconv"
	"strings"
)

type Rule struct {
	Name               ReferToken
	Expression         TokenPointer
	ReturnsNil         bool
	ReturnsString      bool
	SelfReferencing    bool
	SelfRefAtBegin     bool
	SelfRefAtBeginOnly bool
	SelfRefAtEndOnly   bool
}

var rulemap = map[string]Rule{}

func (r Rule) WritePegTo(buffer *bytes.Buffer) {
	buffer.WriteString(r.Name.String())
	buffer.WriteString("\n    = ")
	r.Expression.WritePegTo(buffer)
	if r.ReturnsNil {
		buffer.WriteString(" {\n        return nil, nil\n    }")
	} else if r.ReturnsString {
		buffer.WriteString(" {\n        return string(c.text), nil\n    }")
	}
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

func (r Rule) SplitMultiChoice() []Rule {
	tokenGroup, ok := r.Expression.(*TokenGroup)
	if !ok || tokenGroup.Type != Choice || len(tokenGroup.Tokens) == 1 {
		return []Rule{r}
	}
	result := []Rule{}

	split := false
	for _, token := range tokenGroup.Tokens {
		group, ok := token.(*TokenGroup)
		if !ok || group.Type != Sequence {
			continue
		}
		for _, childToken := range group.Tokens {
			ref, ok := childToken.(*ReferToken)
			if !ok || ref.Name == "_" {
				continue
			}
			rule, keyOk := rulemap[ref.Name]
			if !keyOk {
				continue
			}
			_, typeOk := rule.Expression.(*TokenGroup)
			if typeOk {
				split = true
				break
			}
		}
		if split {
			break
		}
	}

	if !split {
		return []Rule{r}
	}

	heirTokenGroup := TokenGroup{
		Tokens: []TokenPointer{},
		IsRoot: true,
		Repeat: tokenGroup.Repeat,
		Type:   Choice,
	}
	result = append(result, Rule{
		Name:       r.Name,
		Expression: &heirTokenGroup,
	})

	for ind, token := range tokenGroup.Tokens {
		group, ok := token.(*TokenGroup)
		if !ok {
			(&heirTokenGroup).Tokens = append(heirTokenGroup.Tokens, token)
			continue
		}
		group.IsRoot = true
		ruleName := MakeReferToken(r.Name.Name+"_option"+strconv.Itoa(ind+1), true, 0)
		result = append(result, Rule{
			Name:       *ruleName,
			Expression: token,
		})
		(&heirTokenGroup).Tokens = append(heirTokenGroup.Tokens, MakeReferToken(ruleName.Name, false, 0))
	}

	return result
}

func (r Rule) SplitSelfRef(nameToken *ReferToken) []Rule {
	tokenGroup, ok := r.Expression.(*TokenGroup)
	if !ok || tokenGroup.Type != Choice {
		return []Rule{r}
	}
	var normalGroups []TokenPointer
	var selfRefGroups []TokenPointer
	if nameToken == nil {
		nameToken = &r.Name
	}
	newRuleName := MakeReferToken(nameToken.Name+"_self_ref_split", true, 0)
	for _, token := range tokenGroup.Tokens {
		group, ok := token.(*TokenGroup)
		if !ok {
			normalGroups = append(normalGroups, token)
			continue
		}
		ind := group.IndexOfToken(*nameToken)
		if ind == 0 {
			group.ReplaceToken(*nameToken, *newRuleName)
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
	nameToken := ReferToken{
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
	rulemap[nameToken.Name] = rule
	return rule, nil
}

func MakeToken(src string) TokenPointer {
	if strings.Index(src, " ") != -1 {
		return MakeTokenGroup(src)
	}
	return MakeSimpleToken(src)
}
