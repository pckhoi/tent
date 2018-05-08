package storage

import (
	"fmt"
	"github.com/pckhoi/tent/internal/app/settings"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getRootDir() (string, error) {
	s := settings.Get()
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if s.RootDir != "" {
		return filepath.Join(wd, s.RootDir), nil
	}
	return wd, nil
}

func getFilePath(tableName string) (string, error) {
	rootDir, err := getRootDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", rootDir, tableName), nil
}

func fileExist(fileName string) bool {
	path, err := getFilePath(fileName)
	if err != nil {
		return false
	}
	_, openErr := os.Open(path)
	if openErr != nil {
		return false
	}
	return true
}

func ensureTableDir(tableName string) (string, error) {
	tableDir, err := getFilePath(tableName)
	if err != nil {
		return "", err
	}
	chDirErr := os.Chdir(tableDir)
	if chDirErr != nil {
		mkdirErr := os.MkdirAll(tableDir, 0755)
		if mkdirErr != nil {
			return "", errors.Wrap(mkdirErr, "Can't create table dir")
		}
	}
	return tableDir, nil
}

func writeRow(fileName string, content string) error {
	filePath, err := getFilePath(fileName)
	if err != nil {
		return err
	}
	writeErr := ioutil.WriteFile(filePath, []byte(content), 0755)
	return writeErr
}
