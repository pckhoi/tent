package postgres

import (
	parser "github.com/pckhoi/tent/internal/app/parsers/postgres"
	"github.com/pckhoi/tent/internal/app/storage"
	"log"
)

func From(filename string) {
	writeChannel := storage.CreateWriteChannel()

	for _, row := range defaultRows {
		writeChannel <- storage.Write{Row: row, Overwrite: false}
	}

	parseFile(filename, writeChannel)
}

func parseFile(filename string, writeChannel chan storage.Write) {
	got, err := parser.ParseFile(
		filename,
		parser.Memoize(true),
		parser.Debug(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, obj := range got.([]interface{}) {
		if obj == nil {
			continue
		} else if slice, ok := obj.([]storage.DataRow); ok {
			for _, row := range slice {
				log.Printf("emit %s\n", row)
				writeChannel <- storage.Write{Row: row, Overwrite: true}
			}
		} else {
			writeChannel <- storage.Write{Row: obj.(storage.DataRow), Overwrite: true}
		}
	}

	close(writeChannel)
}
