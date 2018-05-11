package postgres

import (
	"github.com/pckhoi/tent/internal/app/storage"
)

func parseCreateExtensionStmt(extension, schema string) storage.DataRow {
	return storage.DataRow{
		TableName: "postgres_extensions",
		ID:        extension,
		Content: map[string]string{
			"name":   extension,
			"schema": schema,
		},
	}
}

func parseCommentExtensionStmt(extension string, comment String) storage.DataRow {
	return storage.DataRow{
		TableName: "postgres_extensions",
		ID:        extension,
		Content: map[string]string{
			"name":    extension,
			"comment": interfaceToString(comment),
		},
	}
}
