package utils

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

func ValidateDirectory(dir string) (string, error) {
	dirPath, err := homedir.Expand(dir)
	if err != nil {
		return "", errors.WithStack(err)
	}
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return "", errors.Wrapf(err, "Directory does not exist: %v\n", dirPath)
	}
	if err != nil {
		return "", errors.Wrapf(err, "Directory error: %v\n", dirPath)

	}
	if !info.IsDir() {
		return "", errors.Errorf("Directory is a file, not a directory: %#v\n", dirPath)
	}
	return dirPath, nil
}

func CheckPathForFile(file string) (string, error) {
	if file == "" {
		panic("File path is empty")
	}
	path := filepath.Dir(file)
	_, err := ValidateDirectory(path)
	if err != nil {
		os.MkdirAll(path, os.ModePerm)
	}
	_, err = ValidateDirectory(path)
	if err != nil {
		panic(err)
	}
	return path, nil
}
