package sql

import (
	// "log"
	"unicode/utf8"
)

type lexerRule struct {
	validStates []lState
	name        string
	matcher     interface{}
	action      func(*sqlSymType, lxr lexer) int
}

func (rule *lexerRule) workInState(state lState) bool {
	if len(rule.validStates) == 0 && state == initial {
		return true
	}
	for _, st := range rule.validStates {
		if st == state {
			return true
		}
	}
	return false
}

func (rule *lexerRule) match(lxr lexer) (bool, int) {
	return match(lxr.data[lxr.start:], rule.matcher)
}

func makeRule(states []lState, matcher func([]byte) (bool, int), action func(*sqlSymType, lexer) int) lexerRule {
	return lexerRule{
		validStates: states,
		matcher:     matcher,
		action:      action,
	}
}

func matchRune(data []byte, val rune) int {
	char, size := utf8.DecodeRune(data)
	if char == val {
		return size
	}
	return 0
}

func matchString(data []byte, val string) int {
	size := len(val)
	if len(data) < size {
		return 0
	}
	for ind := 0; ind < size; ind++ {
		if val[ind] != data[ind] {
			return 0
		}
	}
	return size
}

func match(data []byte, matcher interface{}) (bool, int) {
	switch val := matcher.(type) {
	case rune:
		if size := matchRune(data, val); size > 0 {
			return true, size
		}
	case string:
		if size := matchString(data, val); size > 0 {
			return true, size
		}
	case func([]byte) (bool, int):
		if ok, size := val(data); ok {
			return true, size
		}
	}
	return false, 0
}

func oneOf(matchers ...interface{}) func([]byte) (bool, int) {
	return func(data []byte) (bool, int) {
		for _, m := range matchers {
			if ok, size := match(data, m); ok {
				return true, size
			}
		}
		return false, 0
	}
}

func oneOrMore(matcher interface{}) func([]byte) (bool, int) {
	return func(data []byte) (bool, int) {
		totalLength := 0
		size := 0
		ok := false
		for ind := 0; ind < len(data); ind += size {
			if ok, size = match(data, matcher); ok {
				totalLength += size
			} else {
				break
			}
		}
		return totalLength > 0, totalLength
	}
}

func zeroOrOne(matcher interface{}) func([]byte) (bool, int) {
	return func(data []byte) (bool, int) {
		if ok, size = match(data, matcher); ok {
			return true, size
		}
		return true, 0
	}
}

func concat(matchers ...interface{}) func([]byte) (bool, int) {
	return func(data []byte) (bool, int) {
		matched := true
		totalSize := 0
		for _, m := range matchers {
			if ok, size := match(data[totalSize:], m); ok {
				totalSize += size
			} else {
				matched = false
				break
			}
		}
		if matched {
			return true, totalSize
		}
		return false, 0
	}
}

func anythingBut(matchers ...interface{}) func([]byte) (bool, int) {
	return func(data []byte) (bool, int) {
		char, size := utf8.DecodeRune(data)
		matched := true
		for _, m := range matchers {
			if ok, _ := match(data, m); ok {
				matched := false
				break
			}
		}
		if matched {
			return true, size
		}
		return false, 0
	}
}

func any(matcher interface{}) func([]byte) (bool, int) {
	return func(data []byte) (bool, int) {
		totalLength := 0
		size := 0
		ok := false
		for ind := 0; ind < len(data); ind += size {
			if ok, size = match(data, matcher); ok {
				totalLength += size
			} else {
				break
			}
		}
		return true, totalLength
	}
}

func charRange(startChar rune, endChar rune) func([]byte) (bool, int) {
	return func(data []byte) (bool, int) {
		char, size := utf8.DecodeRune(data)
		if char >= startChar && char <= endChar {
			return true, size
		}
		return false, 0
	}
}

func numChar(matcher interface{}, num ...int) func([]byte) (bool, int) {
	return func(data []byte) (bool, int) {
		totalSize := 0
		matched := true
		if len(num) == 1 {
			for ind := 0; ind < num[0]; ind++ {
				if ok, size = match(data, matcher); ok {
					totalSize += size
				} else {
					matched = false
					break
				}
			}
		} else if len(num) == 2 {
			for ind := 0; ind < num[1]; ind++ {
				if ok, size = match(data, matcher); ok {
					totalSize += size
				} else {
					if ind < num[0] {
						matched = false
					}
					break
				}
			}
		}
		if matched {
			return true, totalSize
		}
		return false, 0
	}
}

func EOF(data []byte) (bool, int) {
    if len(data) == 0 {
        return true, 0
    }
    return false, 0
}
