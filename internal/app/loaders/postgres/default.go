package postgres

import "github.com/pckhoi/tent/internal/app/storage"

var defaultRows = []storage.DataRow{
	storage.DataRow{
		TableName: "schema/postgres_settings",
		ID:        "name",
		Content: map[string]string{
			"type": "string",
		},
	},
	storage.DataRow{
		TableName: "schema/postgres_settings",
		ID:        "setting",
		Content: map[string]string{
			"type": "string",
		},
	},
	storage.DataRow{
		TableName: "schema/postgres_settings",
		ID:        "type",
		Content: map[string]string{
			"type": "string",
		},
	},
}
