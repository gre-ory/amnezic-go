package service

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// file service

type FileService interface {
	List(ctx context.Context, filter *model.FileFilter) ([]model.Url, error)
}

func NewFileService(logger *zap.Logger, fileStore store.FileStore) FileService {
	return &fileService{
		logger:    logger,
		fileStore: fileStore,
	}
}

type fileService struct {
	logger    *zap.Logger
	fileStore store.FileStore
}

// //////////////////////////////////////////////////
// list

func (s *fileService) List(ctx context.Context, filter *model.FileFilter) ([]model.Url, error) {
	return s.fileStore.List(ctx, filter)
}
