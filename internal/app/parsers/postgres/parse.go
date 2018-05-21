package postgres

import (
	"fmt"
	"github.com/pckhoi/tent/internal/app/storage"
	"strconv"
	"strings"
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

func parseCreateTableStmt(tableName interface{}, defs []map[string]string) ([]storage.DataRow, error) {
	table := interfaceToString(tableName)
	results := []storage.DataRow{}
	constraintIndex := 0
	for _, def := range defs {
		if def == nil {
			continue
		}

		var tableName, id string
		if _, ok := def["type"]; ok {
			tableName = fmt.Sprintf("schema/%s", table)
			id = def["name"]
			delete(def, "name")
		} else if _, ok := def["table_constraint"]; ok {
			delete(def, "table_constraint")
			tableName = fmt.Sprintf("constraint/%s", table)
			if val, keyExist := def["constraint_name"]; keyExist {
				id = val
			} else {
				id = strconv.Itoa(constraintIndex)
				constraintIndex++
			}
		}

		results = append(results, storage.DataRow{
			TableName: tableName,
			ID:        id,
			Content:   def,
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

func parseAlterSequenceStmt(name Identifier, owner string) (storage.DataRow, error) {
	return storage.DataRow{
		TableName: "custom/sequence",
		ID:        interfaceToString(name),
		Content: map[string]string{
			"owned_by": owner,
		},
	}, nil
}

func parseTableDotColumn(table, col Identifier) string {
	tableStr := strings.ToLower(interfaceToString(table))
	columnStr := strings.ToLower(interfaceToString(col))
	return strings.Join([]string{tableStr, columnStr}, "/")
}
