package postgres

import (
	parser "github.com/pckhoi/tent/internal/app/parsers/postgres"
	"github.com/pckhoi/tent/internal/app/storage"
	"log"
)

func From(filename string) {
	updates, err := parseFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	writeChannel := storage.CreateWriteChannel()

	for _, row := range defaultRows {
		writeChannel <- storage.Write{Row: row, Overwrite: false}
	}

	for _, row := range updates {
		writeChannel <- storage.Write{Row: row, Overwrite: true}
	}

	close(writeChannel)
}

func parseFile(filename string) ([]storage.DataRow, error) {
	options := parser.Memoize(true)

	got, err := parser.ParseFile(filename, options)
	if err != nil {
		return nil, err
	}
	results := []storage.DataRow{}
	for _, obj := range got.([]interface{}) {
		results = append(results, obj.(storage.DataRow))
	}
	return results, nil
}
