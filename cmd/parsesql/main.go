package main

import (
	"github.com/alecthomas/repr"
	"github.com/pckhoi/tent/internal/app/parsers/sql"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	in := os.Stdin
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
	}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal(err)
	}

	got := sql.Parse(b)
	repr.Print(got)
}
