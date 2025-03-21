package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/model"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// file store

type FileStore interface {
	List(ctx context.Context, filter *model.FileFilter) ([]model.Url, error)
	Exists(ctx context.Context, filter *model.FileFilter, path string) bool
	PathValidator(ctx context.Context, filter *model.FileFilter) model.PathValidator
}

// //////////////////////////////////////////////////
// dummy file store

func NewFileStore(logger *zap.Logger) FileStore {
	return &fileStore{
		logger: logger,
	}
}

type fileStore struct {
	logger *zap.Logger
}

func (s *fileStore) List(_ context.Context, filter *model.FileFilter) ([]model.Url, error) {
	directory := filepath.Clean(filter.Directory)
	urls := make([]model.Url, 0)
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !filter.MatchExtension(info.Name()) {
				return nil
			}
			s.logger.Info(fmt.Sprintf(" (+) file %q", path))
			prefix := directory + string(filepath.Separator)
			path = strings.TrimPrefix(path, prefix)
			urls = append(urls, model.Url(path))
			return nil
		})
	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] list files on directory %q", directory), zap.Object("filter", filter), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] list %d files on directory %q", len(urls), directory), zap.Object("filter", filter))
	sort.Slice(urls, func(i, j int) bool {
		return urls[i] < urls[j]
	})
	return urls, nil
}

func (s *fileStore) Exists(_ context.Context, filter *model.FileFilter, path string) bool {
	if !filter.MatchExtension(path) {
		return false
	}
	path = filepath.Clean(path)
	path = filepath.Join(filter.Directory, path)
	_, err := os.Stat(path)
	return err == nil
}

func (s *fileStore) PathValidator(ctx context.Context, filter *model.FileFilter) model.PathValidator {
	return func(path string) error {
		if !s.Exists(ctx, filter, path) {
			return model.ErrPathNotFound(path)
		}
		return nil
	}
}
