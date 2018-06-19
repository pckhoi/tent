package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	bytesSlice, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		log.Fatal(readErr)
	}

	log.Println("Reading from Bison...")
	grammar := MakeFromBison(bytesSlice)
	var buffer bytes.Buffer
	log.Println("Writing Peg...")
	grammar.WritePegTo(&buffer)
	fmt.Print(buffer.String())
	for k, _ := range namemap {
		cont := false
		for _, rulename := range rulenames {
			if rulename == k {
				cont = true
				break
			}
		}
		if cont {
			continue
		}
		// for _, keyword := range keywords {
		// 	if keyword == k {
		// 		cont = true
		// 		break
		// 	}
		// }
		// if cont {
		// 	continue
		// }
		log.Println(k)
	}
}
