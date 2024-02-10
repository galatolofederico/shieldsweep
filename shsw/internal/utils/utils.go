package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/fatih/color"
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

func CopyFile(src string, dst string) {
	data, err := os.ReadFile(src)
	if err != nil {
		panic(errors.Wrapf(err, "Error reading file: %v\n", src))
	}
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		panic(errors.Wrapf(err, "Error writing file: %v\n", dst))
	}
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
	if date == "never" {
		return "never"
	}
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

func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		panic(errors.Wrap(err, "Error getting current user"))
	}
	return currentUser.Username == "root"
}

func Get(httpc http.Client, path string) []byte {
	response, err := httpc.Get(path)
	if err != nil {
		color.Red("[!] Error: %v\n", err)
		color.Red("[!] Is the daemon running?\n")
		os.Exit(1)
	}
	if response.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(response.Body)
		panic(errors.Errorf("Error: %s\n%s", response.Status, resBody))
	}
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(errors.Wrapf(err, "Error reading response body from %s\n", path))
	}
	return resBody
}

func GetUnixClient(sock string) http.Client {
	return http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sock)
			},
		},
	}
}
