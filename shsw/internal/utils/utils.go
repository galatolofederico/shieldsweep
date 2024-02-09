package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

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

func CheckPathForFile(file string) string {
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
	return path
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err == nil {
		return !info.IsDir()
	}
	return false
}

func SHA256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	sha := h.Sum(nil)
	return hex.EncodeToString(sha), nil
}

func ParseDate(date string) string {
	ret, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "date-parse-error"
	}
	return ret.Format("2006-01-02 15:04:05")
}

func DaysAgo(date string) string {
	now := time.Now()
	prev, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "date-parse-error"
	}
	diff := now.Sub(prev)
	hours := diff.Hours()
	if hours < 24 {
		return fmt.Sprintf("%v hours ago", int(hours))
	} else {
		days := int(hours / 24)
		return fmt.Sprintf("%v days and %v hours ago", days, int(hours)%24)
	}
}
