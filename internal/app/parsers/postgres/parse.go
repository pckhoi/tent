package postgres

import (
	"fmt"
	"github.com/pckhoi/tent/internal/app/storage"
)

func parseCreateExtensionStmt(extension Identifier, schema Identifier) (storage.DataRow, error) {
	return storage.DataRow{
		TableName: "postgres_extensions",
		ID:        interfaceToString(extension),
		Content: map[string]string{
			"name":   interfaceToString(extension),
			"schema": interfaceToString(schema),
		},
	}, nil
}

func parseCommentExtensionStmt(extension Identifier, comment String) (storage.DataRow, error) {
	return storage.DataRow{
		TableName: "postgres_extensions",
		ID:        interfaceToString(extension),
		Content: map[string]string{
			"name":    interfaceToString(extension),
			"comment": interfaceToString(comment),
		},
	}, nil
}

func parseCreateTableStmt(tableName interface{}, fields []map[string]string) ([]storage.DataRow, error) {
	table := interfaceToString(tableName)
	results := []storage.DataRow{}
	for _, field := range fields {
		if field == nil {
			continue
		}
		fieldName := field["name"]
		delete(field, "name")
		results = append(results, storage.DataRow{
			TableName: fmt.Sprintf("schema/%s", table),
			ID:        fieldName,
			Content:   field,
		})
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results, nil
}

func parseCreateTypeEnumStmt(enum Enum) (storage.DataRow, error) {
	updateSettings("custom_types", []interface{}{interfaceToString(enum.Name)})
	return storage.DataRow{
		TableName: "custom/type",
		ID:        interfaceToString(enum.Name),
		Content: map[string]string{
			"type":   "enum",
			"labels": serializeStringSlice(enum.Labels),
		},
	}, nil
}

func parseCreateSeq(name Identifier, properties map[string]string) (storage.DataRow, error) {
	return storage.DataRow{
		TableName: "custom/sequence",
		ID:        interfaceToString(name),
		Content:   properties,
	}, nil
}

func parseAlterTableStmt(name, owner Identifier) (storage.DataRow, error) {
	return storage.DataRow{
		TableName: "custom/table",
		ID:        interfaceToString(name),
		Content: map[string]string{
			"owner": interfaceToString(owner),
		},
	}, nil
}
