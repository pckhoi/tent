package postgres

import "github.com/pckhoi/tent/internal/app/storage"

var defaultRows = []storage.DataRow{
	storage.DataRow{
		TableName: "schema/postgres_extensions",
		ID:        "name",
		Content: map[string]string{
			"type": "string",
		},
	},
	storage.DataRow{
		TableName: "schema/postgres_extensions",
		ID:        "schema",
		Content: map[string]string{
			"type": "string",
		},
	},
	storage.DataRow{
		TableName: "schema/postgres_extensions",
		ID:        "comment",
		Content: map[string]string{
			"type": "string",
		},
	},
}
