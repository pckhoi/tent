package main

import (
	"fmt"
	"github.com/alecthomas/repr"
	"github.com/pckhoi/tent/internal/app/parsers/postgres"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	in := os.Stdin
	nm := "stdin"
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		in = f
		nm = os.Args[1]
	}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal(err)
	}

	got, err := postgres.Parse(nm, b)
	if err != nil {
		fmt.Println(err, string(b))
		os.Exit(1)
	}
	repr.Print(got)
}
