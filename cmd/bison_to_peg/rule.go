package main

import (
	"bytes"
	"errors"
	"fmt"
	// "log"
	"regexp"
	"strconv"
	"strings"
)

type Rule struct {
	Name                 ReferToken
	Expression           TokenPointer
	Refers               []*ReferToken
	ReturnsNil           bool
	ReturnsString        bool
	ReturnsComplexString bool
	ReturnsLiteral       string
	SelfReferencing      bool
	SelfRefAtBegin       bool
	SelfRefAtBeginOnly   bool
	SelfRefAtEndOnly     bool
	ExpressionListVar    string
	ListFirstElIndex     int
	ListRestIndex        int
	ListRestElIndex      int
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
	} else if r.ReturnsLiteral != "" {
		buffer.WriteString(fmt.Sprintf(" {\n        return %s\n    }", r.ReturnsLiteral))
	} else if r.ReturnsComplexString {
		buffer.WriteString(fmt.Sprintf(" {\n        return complexString(%s), nil\n    }", r.Expression.GetVariableName()))
	} else if len(r.Refers) > 0 {
		buffer.WriteString(" {\n        return Node{\n            Name: \"")
		buffer.WriteString(r.Name.String())
		buffer.WriteString("\",\n            Props: map[string]interface{}{\n")
		gen := nameWithCount()
		for _, refer := range r.Refers {
			buffer.WriteString("                \"")
			buffer.WriteString(gen(namemap[refer.Name]))
			buffer.WriteString("\": ")
			buffer.WriteString(refer.VariableName)
			buffer.WriteString(",\n")
		}
		buffer.WriteString("            },\n")
		buffer.WriteString("        }, nil\n")
		buffer.WriteString("    }")
	} else if r.ExpressionListVar != "" {
		buffer.WriteString(" {\n        vals := extractList(")
		buffer.WriteString(r.ExpressionListVar)
		buffer.WriteString(", ")
		buffer.WriteString(strconv.Itoa(r.ListFirstElIndex))
		buffer.WriteString(", ")
		buffer.WriteString(strconv.Itoa(r.ListRestIndex))
		buffer.WriteString(", ")
		buffer.WriteString(strconv.Itoa(r.ListRestElIndex))
		buffer.WriteString(")\n        return Node{\n            Name: \"")
		buffer.WriteString(r.Name.String())
		buffer.WriteString("\",\n            Vals: vals,\n")
		buffer.WriteString("        }, nil\n")
		buffer.WriteString("    }")
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

func (r *Rule) DressupList() {
	expressionListVar, listFirstElIndex, listRestIndex, listRestElIndex := r.Expression.(*TokenGroup).DressupList()
	r.ExpressionListVar = expressionListVar
	r.ListFirstElIndex = listFirstElIndex
	r.ListRestIndex = listRestIndex
	r.ListRestElIndex = listRestElIndex
}

func (r *Rule) Dressup() {
	tokenGroup, ok := r.Expression.(*TokenGroup)
	if !ok || tokenGroup.Type != Sequence || len(tokenGroup.Tokens) == 1 {
		return
	}

	r.DressupList()
	if r.ExpressionListVar != "" {
		return
	}

	if tokenGroup.Repeat != 0 {
		return
	}
	refers := []*ReferToken{}
	counter := 0
	for _, token := range tokenGroup.Tokens {
		refer, ok := token.(*ReferToken)
		if !ok {
			continue
		}
		if refer.Name != "_" && refer.Name != strings.ToUpper(refer.Name) {
			counter++
			refer.VariableName = "var" + strconv.Itoa(counter)
			refers = append(refers, refer)
		}
	}

	if len(refers) == 0 {
		return
	}

	r.Refers = refers
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

	referNewRule := *newRuleName
	heirTokenGroup := TokenGroup{
		Tokens: append(
			selfRefGroups,
			&referNewRule,
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
