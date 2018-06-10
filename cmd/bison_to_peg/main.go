package main

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type Rule struct {
	Name     string
	Subrules []Subrule
	Optional bool
}

type Subrule struct {
	Words []string
}

type Token struct {
	IsRuleName    bool
	DisplayString string
}

var tokenMap = map[string]Token{}

func ruleNameToTokenMap(word string) {
	tokenMap[word] = Token{
		IsRuleName:    true,
		DisplayString: strcase.ToCamel(word),
	}
}

func concatTokens(words []string) string {
	var transformedWords []string
	for _, word := range words {
		if token, ok := tokenMap[word]; ok {
			transformedWords = append(transformedWords, token.DisplayString)
		} else {
			transformedWords = append(transformedWords, word)
		}
	}
	return strings.Join(transformedWords, " _ ")
}

func main() {
	if len(os.Args) < 1 {
		log.Fatal(fmt.Errorf("too few arguments, pass input file as first argument"))
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	bytesSlice, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		log.Fatal(readErr)
	}

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

	// Split rules
	re = regexp.MustCompile(`\s*\n *;\n\s*`)
	nameRe := regexp.MustCompile(`(\w+) *:\s+`)
	pipeRe := regexp.MustCompile(`\s*\|\s*`)
	spaceRe := regexp.MustCompile(`\s+`)
	rulesSlice := re.Split(src, -1)
	log.Printf("There are %d rules\n", len(rulesSlice))
	ruleObjsSlice := []Rule{}
	for _, rule := range rulesSlice {
		indicies := nameRe.FindStringSubmatchIndex(rule)
		if indicies == nil {
			continue
		}
		ruleName := rule[indicies[2]:indicies[3]]
		ruleNameToTokenMap(ruleName)
		rule = rule[indicies[1]:]
		subruleStrings := pipeRe.Split(rule, -1)
		subrules := []Subrule{}
		optional := false
		for _, v := range subruleStrings {
			if v == "" {
				optional = true
				continue
			}
			subrules = append(subrules, Subrule{
				Words: spaceRe.Split(v, -1),
			})
		}
		ruleObjsSlice = append(ruleObjsSlice, Rule{
			Name:     ruleName,
			Subrules: subrules,
			Optional: optional,
		})
	}

	// present the rules again
	var buffer bytes.Buffer
	for _, obj := range ruleObjsSlice {
		buffer.WriteString(tokenMap[obj.Name].DisplayString)
		buffer.WriteString("\n    <- ")
		if obj.Optional {
			buffer.WriteString("(")
		}
		var subrules []string
		for _, subrule := range obj.Subrules {
			resultString := concatTokens(subrule.Words)
			if len(subrule.Words) > 1 && len(obj.Subrules) > 1 {
				resultString = "(" + resultString + ")"
			}
			subrules = append(subrules, resultString)
		}
		buffer.WriteString(strings.Join(subrules, "\n    / "))
		if obj.Optional {
			if len(obj.Subrules) > 1 {
				buffer.WriteString("\n    ")
			}
			buffer.WriteString(")?")
		}
		buffer.WriteString("\n\n")
	}
	fmt.Print(buffer.String())
}
