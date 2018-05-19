package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

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
	bytes, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		log.Fatal(readErr)
	}
	src := string(bytes)

	re := regexp.MustCompile(`\s*\{\n(?:(?:    .+)?\n)+\}`)
	src = re.ReplaceAllString(src, "")

	re = regexp.MustCompile(`(\w+)(\s+"[^"]+")?\s+<-\s+([^\n]+)`)
	src = re.ReplaceAllString(src, "$1$2 =\n  $3")

	re = regexp.MustCompile("`([^`]+)`")
	src = re.ReplaceAllString(src, "\"$1\"")

	fmt.Println(string(src))
}
