package main

import (
	"bytes"
	"regexp"
)

type Token struct {
	IsRuleName    bool
	DisplayString string
}

type Context struct {
	TokenMap      map[string]Token
	SiblingsCount int
}

type Grammar struct {
	Rules    []Rule
	TokenMap map[string]Token
}

func cleanInput(bytesSlice []byte) string {
	// Remove everything keep only the rules declaration
	ruleSep := []byte("\n%%\n")
	beginOfRules := bytes.Index(bytesSlice, ruleSep) + 4
	endOfRules := bytes.LastIndex(bytesSlice, []byte("\n%%\n"))
	bytesSlice = bytesSlice[beginOfRules:endOfRules]

	src := string(bytesSlice)

	// Remove code blocks
	re := regexp.MustCompile(`\n +\{[^{}]+\}\n`)
	for re.Match([]byte(src)) {
		src = re.ReplaceAllString(src, "\n")
	}

	// Remove inline code blocks
	re = regexp.MustCompile(`( +)\{[^{}]*\}\n`)
	src = re.ReplaceAllString(src, "\n")

	// Remove trailing comment
	re = regexp.MustCompile(`(\w) +/\*.*\*/\n`)
	src = re.ReplaceAllString(src, "$1\n")

	// Remove prec directive
	re = regexp.MustCompile(` +\%prec +\w+\n`)
	src = re.ReplaceAllString(src, "\n")

	// Remove inline code blocks with trailing semi
	re = regexp.MustCompile(`\s+\{[^{}]+\};\n`)
	src = re.ReplaceAllString(src, "\n        ;\n")

	// Remove all inline comments
	re = regexp.MustCompile(` */\*.+\*/`)
	src = re.ReplaceAllString(src, "")

	// Remove all inline comments
	re = regexp.MustCompile(`\n */\*.*\n( *\*.*\n)+`)
	src = re.ReplaceAllString(src, "\n")

	// Remove all trailing _P
	re = regexp.MustCompile(`_(?:P|LA)(\s)`)
	src = re.ReplaceAllString(src, "$1")

	return src
}

func MakeFromBison(bytes []byte) Grammar {
	src := cleanInput(bytes)
	re := regexp.MustCompile(`\s*\n *;\n\s*`)
	ruleStrings := re.Split(src, -1)
	rules := []Rule{}
	tokenMap := map[string]Token{}
	for _, ruleString := range ruleStrings {
		rule, err := MakeRule(ruleString, tokenMap)
		if err != nil {
			continue
		}
		rules = append(rules, rule)
	}
	return Grammar{
		TokenMap: tokenMap,
		Rules:    rules,
	}
}

func (g Grammar) WritePegTo(buffer *bytes.Buffer) {
	context := Context{
		TokenMap:      g.TokenMap,
		SiblingsCount: 0,
	}
	for _, rule := range g.Rules {
		rule.WritePegTo(buffer, context)
		buffer.WriteString("\n\n")
	}
}
