package main

import (
	"bytes"
	// "log"
	"regexp"
	"sort"
)

type ExpressionType int

const (
	Sequence ExpressionType = iota
	Choice
)

type TokenGroup struct {
	Tokens []TokenPointer
	IsRoot bool
	Repeat RepeatCharacteristic
	Type   ExpressionType
}

func (group TokenGroup) GetTokenGroups() func() (int, *TokenGroup) {
	ind := -1
	lenTokens := len(group.Tokens)

	return func() (int, *TokenGroup) {
		for true {
			ind++
			if ind >= lenTokens {
				break
			}
			if val, ok := group.Tokens[ind].(*TokenGroup); ok {
				return ind, val
			}
		}
		return -1, nil
	}
}

func (group TokenGroup) GetSimpleTokens() func() (int, *SimpleToken) {
	ind := -1
	lenTokens := len(group.Tokens)

	return func() (int, *SimpleToken) {
		for true {
			ind++
			if ind >= lenTokens {
				break
			}
			if val, ok := group.Tokens[ind].(*SimpleToken); ok {
				return ind, val
			}
		}
		return -1, nil
	}
}

func (group *TokenGroup) String() string {
	var buffer bytes.Buffer
	group.WritePegTo(&buffer)
	return buffer.String()
}

func (group *TokenGroup) WritePegTo(buffer *bytes.Buffer) {
	bracket := !group.IsRoot || group.Repeat != 0
	if bracket {
		buffer.WriteString("(")
	}
	var separator string
	if group.Type == Sequence {
		separator = " "
	} else if group.IsRoot {
		separator = "\n    / "
	} else {
		separator = " / "
	}
	for ind, token := range group.Tokens {
		token.WritePegTo(buffer)
		if ind < len(group.Tokens)-1 {
			buffer.WriteString(separator)
		}
	}
	if bracket {
		buffer.WriteString(")")
	}
	switch group.Repeat {
	case OneOrMany:
		buffer.WriteString("+")
	case OneOrNone:
		buffer.WriteString("?")
	case Any:
		buffer.WriteString("*")
	}
}

func (group *TokenGroup) SetRepeat(repeat RepeatCharacteristic) {
	group.Repeat = repeat
}

func (group *TokenGroup) GetRepeat() RepeatCharacteristic {
	return group.Repeat
}

func (s *TokenGroup) IndexOfToken(tok SimpleToken) int {
	gen := s.GetSimpleTokens()
	for ind, token := gen(); token != nil; ind, token = gen() {
		if token.Name == tok.Name {
			return ind
		}
	}
	return -1
}

func (s *TokenGroup) AllIndexOfToken(tok SimpleToken) []int {
	result := []int{}
	gen := s.GetSimpleTokens()
	for ind, token := gen(); token != nil; ind, token = gen() {
		if token.Name == tok.Name {
			result = append(result, ind)
		}
	}
	return result
}

func (s *TokenGroup) ReplaceToken(search, replace SimpleToken) {
	gen := s.GetSimpleTokens()
	for ind, val := gen(); val != nil; ind, val = gen() {
		if val.Name == search.Name {
			s.Tokens[ind] = &replace
		}
	}
}

func (group *TokenGroup) DetectSelfReferencing(name SimpleToken) (bool, bool, bool, bool) {
	selfRef := false
	selfRefAtBegin := false
	selfRefAtBeginOnly := false
	selfRefAtEndOnly := false
	selfRefGroups := []*TokenGroup{}
	gen := group.GetTokenGroups()
	for _, val := gen(); val != nil; _, val = gen() {
		inds := val.AllIndexOfToken(name)
		if len(inds) >= 1 {
			selfRef = true
			if inds[0] == 0 {
				selfRefAtBegin = true
			}
			selfRefGroups = append(selfRefGroups, val)
		}
	}
	if selfRef && (len(group.Tokens)-len(selfRefGroups)) <= 1 {
		selfRefAtBeginOnly = true
		selfRefAtEndOnly = true

		for _, group := range selfRefGroups {
			for _, ind := range group.AllIndexOfToken(name) {
				if ind > 0 {
					selfRefAtBeginOnly = false
				}
				if ind < len(group.Tokens)-1 {
					selfRefAtEndOnly = false
				}
			}
		}
	}
	return selfRef, selfRefAtBegin, selfRefAtBeginOnly, selfRefAtEndOnly
}

func (group *TokenGroup) FixSelfRefAtBeginOnly(name SimpleToken, reversed bool) {
	beginTokens := []TokenPointer{}
	var beginToken TokenPointer
	for _, token := range group.Tokens {
		if val, ok := token.(*TokenGroup); ok {
			if val.IndexOfToken(name) == -1 {
				beginTokens = val.Tokens[:]
				beginToken = val
				break
			}
		} else if val, ok := token.(*SimpleToken); ok {
			beginTokens = []TokenPointer{val}
			beginToken = val
			break
		}
	}
	normalTokens := []TokenPointer{}
	processedTokens := []TokenPointer{}
	for _, token := range group.Tokens {
		if token == beginToken {
			continue
		}
		if val, ok := token.(*TokenGroup); ok {
			selfrefInd := 0
			if reversed {
				selfrefInd = len(val.Tokens) - 1
			}
			if val.Type == Choice || val.IndexOfToken(name) != selfrefInd {
				normalTokens = append(normalTokens, val)
				continue
			}
			newTokens := val.Tokens[1:]
			if reversed {
				newTokens = val.Tokens[:selfrefInd]
			}
			childTokenGroup := TokenGroup{
				Tokens: newTokens,
				Type:   Sequence,
			}
			processedTokens = append(processedTokens, &childTokenGroup)
		} else if val, ok := token.(*SimpleToken); ok {
			normalTokens = append(normalTokens, val)
		}
	}

	var anyToken TokenPointer
	if len(processedTokens) > 1 {
		anyToken = &TokenGroup{
			Tokens: processedTokens,
			Type:   Choice,
			Repeat: Any,
		}
	} else {
		anyToken = processedTokens[0]
		anyToken.SetRepeat(Any)
	}
	newTokens := append(beginTokens, anyToken)
	if reversed {
		newTokens = append([]TokenPointer{anyToken}, beginTokens...)
	}
	newTokenGroup := TokenGroup{
		Tokens: newTokens,
		Type:   Sequence,
	}
	group.Tokens = append([]TokenPointer{&newTokenGroup}, normalTokens...)
}

func (group *TokenGroup) ReorderTokens() {
	if group.Type == Sequence {
		return
	}
	sort.Slice(group.Tokens, func(i, j int) bool {
		tokenA := group.Tokens[i]
		tokenB := group.Tokens[j]
		var lenTokenA = 1
		var lenTokenB = 1
		if val, ok := tokenA.(*TokenGroup); ok {
			lenTokenA = len(val.Tokens)
		}
		if val, ok := tokenB.(*TokenGroup); ok {
			lenTokenB = len(val.Tokens)
		}
		return lenTokenA > lenTokenB
	})
}

func (group *TokenGroup) Simplify() TokenPointer {
	result := group.MergeOnlyChild()
	if val, ok := result.(*TokenGroup); ok {
		result = val.Deduplicate()
	}
	return result
}

func (group *TokenGroup) MergeOnlyChild() TokenPointer {
	for ind, token := range group.Tokens {
		if tokenval, ok := token.(*TokenGroup); ok {
			token = tokenval.MergeOnlyChild()
		}
		group.Tokens[ind] = token
	}
	var result TokenPointer
	result = group
	parent, ok := result.(*TokenGroup)
	for ok && len(parent.Tokens) == 1 {
		result = parent.Tokens[0]
		if parent.Repeat > result.GetRepeat() {
			result.SetRepeat(parent.Repeat)
		}
		if val, ok := result.(*TokenGroup); ok {
			val.IsRoot = parent.IsRoot
		}
		parent, ok = result.(*TokenGroup)
	}
	return result
}

func (group *TokenGroup) Deduplicate() TokenPointer {
	if group.Type == Sequence {
		return group
	}
	tokenCombinations := combinations(group.Tokens, 2)
	toRemove := []TokenPointer{}
	for _, combination := range tokenCombinations {
		tokenA := combination[0]
		tokenB := combination[1]
		if tokenEqual(tokenA, tokenB) {
			toRemove = append(toRemove, tokenB)
		}
	}
	for _, pointer := range toRemove {
		for ind, token := range group.Tokens {
			if token == pointer {
				group.Tokens = append(group.Tokens[:ind], group.Tokens[ind+1:]...)
			}
		}
	}
	if val, ok := group.Tokens[0].(*SimpleToken); len(group.Tokens) == 1 && ok {
		val.SetRepeat(group.Repeat)
		return val
	}
	return group
}

func (group *TokenGroup) InsertSpace() {
	gen := group.GetTokenGroups()
	for _, val := gen(); val != nil; _, val = gen() {
		val.InsertSpace()
	}
	if group.Type == Choice {
		return
	}
	newTokens := []TokenPointer{}
	oldlen := len(group.Tokens)
	for ind, token := range group.Tokens {
		newTokens = append(newTokens, token)
		if ind < oldlen-1 {
			spaceToken := MakeSimpleToken("_")
			newTokens = append(newTokens, spaceToken)
		}
	}
	group.Tokens = newTokens
}

func tokenEqual(tokenA, tokenB TokenPointer) bool {
	simA, simAOk := tokenA.(*SimpleToken)
	simB, simBOk := tokenB.(*SimpleToken)
	if simAOk && simBOk {
		return simA.Name == simB.Name
	}
	groupA, groupAOk := tokenA.(*TokenGroup)
	groupB, groupBOk := tokenB.(*TokenGroup)
	if groupAOk && groupBOk {
		if groupA.Type != groupB.Type {
			return false
		}
		if len(groupA.Tokens) != len(groupB.Tokens) {
			return false
		}
		if groupA.Repeat != groupB.Repeat {
			return false
		}
		if groupA.IsRoot != groupB.IsRoot {
			return false
		}
		for ind, val := range groupA.Tokens {
			if !tokenEqual(val, groupB.Tokens[ind]) {
				return false
			}
		}
		return true
	}
	return false
}

var spaceRe = regexp.MustCompile(`\s+`)

func MakeTokenGroup(src string) *TokenGroup {
	words := spaceRe.Split(src, -1)
	tokens := []TokenPointer{}
	for _, word := range words {
		token := MakeSimpleToken(word)
		tokens = append(tokens, token)
	}
	group := TokenGroup{
		Tokens: tokens,
	}
	return &group
}
