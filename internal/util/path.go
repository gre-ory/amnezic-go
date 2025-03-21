package util

import (
	"fmt"
	"path/filepath"
	"strings"
)

// //////////////////////////////////////////////////
// helper

var (
	ErrInvalidLocalPath = fmt.Errorf("invalid local path")
	ErrInvalidExtension = fmt.Errorf("invalid extension")
)

func CleanLocalPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if strings.Contains(path, "..") {
		return "", ErrInvalidLocalPath
	}
	path = filepath.Clean(path)
	if path == "" {
		path = "."
	}
	return path, nil
}

func CleanExtension(extension string) (string, error) {
	extension = strings.TrimSpace(extension)
	if extension == "." {
		return "", ErrInvalidExtension
	}
	extension = strings.TrimPrefix(extension, ".")
	if extension == "" {
		return "", nil
	}
	if strings.Contains(extension, ".") {
		return "", ErrInvalidExtension
	}
	extension = strings.ToLower(extension)
	return fmt.Sprintf(".%s", extension), nil
}
