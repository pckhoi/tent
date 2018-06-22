package main

import (
	"bytes"
	// "log"
	"regexp"
)

type Grammar struct {
	Rules []Rule
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
	for _, ruleString := range ruleStrings {
		rule, err := MakeRule(ruleString)
		if err != nil {
			continue
		}
		rule.Inspect()
		if rule.SelfRefAtBeginOnly {
			rule.FixSelfRefAtBeginOnly(false)
			rules = append(rules, rule)
		} else if rule.SelfRefAtEndOnly {
			rule.FixSelfRefAtBeginOnly(true)
			rules = append(rules, rule)
		} else if rule.SelfRefAtBegin {
			replacingRules := rule.SplitSelfRef(nil)
			rules = append(rules, replacingRules...)
		} else {
			rules = append(rules, rule)
		}
	}
	for ind, rule := range rules {
		rule.Simplify()
		rule.ReorderSubrules()
		rules[ind] = rule
	}
	rules = append(rules, keywordRules()...)
	rules = append(rules, miscellaneousRules()...)
	return Grammar{
		Rules: rules,
	}
}

func (g Grammar) WritePegTo(buffer *bytes.Buffer) {
	buffer.WriteString("{\n    package postgres\n}\n\n")

	for _, rule := range g.Rules {
		rule.WritePegTo(buffer)
		buffer.WriteString("\n\n")
	}
}
