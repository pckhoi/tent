package storage

import (
	"fmt"
	"github.com/pckhoi/tent/internal/app/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func updateFieldNames(row DataRow) ([]string, error) {
	rowFieldNames := row.getFieldNames()
	existingFieldNames, getErr := getFieldNames(row.TableName)
	if getErr != nil {
		return nil, getErr
	}
	fieldNames := utils.MergeStringSlices(existingFieldNames, rowFieldNames)
	err := setFieldNames(row.TableName, fieldNames)
	if err != nil {
		return nil, err
	}
	return fieldNames, nil
}

var fieldNamesMap = make(map[string][]string)
var fieldNamesMapMutex = &sync.Mutex{}

func getFieldNames(tableName string) ([]string, error) {
	if fieldNames, ok := fieldNamesMap[tableName]; ok {
		return fieldNames, nil
	}
	names, err := readFieldNamesFromFile(tableName)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read field names from file")
	}
	fieldNamesMapMutex.Lock()
	fieldNamesMap[tableName] = names
	fieldNamesMapMutex.Unlock()
	return names, nil
}

func setFieldNames(tableName string, fieldNames []string) error {
	if utils.StringSliceEqual(fieldNamesMap[tableName], fieldNames) {
		return nil
	}

	fieldNamesMapMutex.Lock()
	fieldNamesMap[tableName] = fieldNames
	fieldNamesMapMutex.Unlock()
	err := writeFieldNamesToFile(tableName, fieldNames)
	return err
}

func readFieldNamesFromFile(tableName string) ([]string, error) {
	path, err := getFieldNamesFilePath(tableName)
	if err != nil {
		return nil, err
	}
	file, openErr := os.Open(path)
	if openErr == nil {
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, errors.Wrap(err, "Can't read field names from file")
		}
		return strings.Split(string(bytes), "\n"), nil
	}
	return []string{}, nil
}

func writeFieldNamesToFile(tableName string, fieldNames []string) error {
	path, err := getFieldNamesFilePath(tableName)
	if err != nil {
		return err
	}
	log.Printf("folder to create %s\n", filepath.Dir(path))
	mkdirErr := os.MkdirAll(filepath.Dir(path), 0755)
	if mkdirErr != nil {
		return errors.Wrap(mkdirErr, "Can't make parent dir for field names file")
	}
	_, createErr := os.Create(path)
	if createErr != nil {
		return errors.Wrap(createErr, "Can't open file to write")
	}
	content := strings.Join(fieldNames, "\n")
	return ioutil.WriteFile(path, []byte(content), 0755)
}

func getFieldNamesFilePath(tableName string) (string, error) {
	tableDir, err := getFilePath(tableName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.fieldnames", tableDir), nil
}
