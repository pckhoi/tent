package storage

import (
	"fmt"
	"strings"
)

type DataRow struct {
	TableName string
	ID        string
	Content   map[string]string
}

func (row *DataRow) getFileName() string {
	return fmt.Sprintf("%s/%s", row.TableName, row.ID)
}

func (row *DataRow) getFileContent(fieldNames []string) string {
	var result []string
	for _, key := range fieldNames {
		result = append(result, row.Content[key])
	}
	return strings.Join(result, "\n")
}

func (row *DataRow) getFieldNames() []string {
	var result []string
	for key := range row.Content {
		result = append(result, key)
	}
	return result
}
