package model

import (
	"strings"

	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// file filter

type FileFilter struct {
	Directory  string
	Extensions []string
}

func NewFileFilter(directory string, extensions []string) (*FileFilter, error) {
	cleanedDirectory, err := util.CleanLocalPath(directory)
	if err != nil {
		return nil, err
	}
	cleanedExtensions := make([]string, 0, len(extensions))
	for _, extension := range extensions {
		cleanedExtension, err := util.CleanExtension(extension)
		if err != nil {
			return nil, err
		}
		cleanedExtensions = append(cleanedExtensions, cleanedExtension)
	}
	return &FileFilter{
		Directory:  cleanedDirectory,
		Extensions: cleanedExtensions,
	}, nil
}

func (o *FileFilter) MatchExtension(filename string) bool {
	if len(o.Extensions) == 0 {
		return true
	}
	filename = strings.ToLower(filename)
	for _, extension := range o.Extensions {
		if strings.HasSuffix(filename, strings.ToLower(extension)) {
			return true
		}
	}
	return false
}

func (o *FileFilter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if o.Directory != "" {
		enc.AddString("directory", o.Directory)
	}
	if o.Extensions != nil {
		enc.AddString("extensions", util.Join(o.Extensions, ","))
	}
	return nil
}
