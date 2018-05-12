package postgres

import (
	"fmt"
	"github.com/pckhoi/tent/internal/app/storage"
)

func parseCreateExtensionStmt(extension Identifier, schema Identifier) storage.DataRow {
	return storage.DataRow{
		TableName: "postgres_extensions",
		ID:        interfaceToString(extension),
		Content: map[string]string{
			"name":   interfaceToString(extension),
			"schema": interfaceToString(schema),
		},
	}
}

func parseCommentExtensionStmt(extension Identifier, comment String) storage.DataRow {
	return storage.DataRow{
		TableName: "postgres_extensions",
		ID:        interfaceToString(extension),
		Content: map[string]string{
			"name":    interfaceToString(extension),
			"comment": interfaceToString(comment),
		},
	}
}

func parseCreateTableStmt(tableName interface{}, fields []map[string]string) []storage.DataRow {
	table := interfaceToString(tableName)
	results := []storage.DataRow{}
	for _, field := range fields {
		fieldName := field["name"]
		delete(field, "name")
		results = append(results, storage.DataRow{
			TableName: fmt.Sprintf("schema/%s", table),
			ID:        fieldName,
			Content:   field,
		})
	}
	return results
}
