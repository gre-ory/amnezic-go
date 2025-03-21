package store_test

import (
	"context"
	"testing"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestFileStore(t *testing.T) {
	ctx := context.Background()
	logger := zap.L()

	filter := &model.FileFilter{
		Directory:  "test/image",
		Extensions: []string{".mp3", ".txt"},
	}

	fileStore := store.NewFileStore(logger)

	require.True(t, fileStore.Exists(ctx, filter, "img_1.txt"))
	require.True(t, fileStore.Exists(ctx, filter, "sub/img_2.txt"))
	require.False(t, fileStore.Exists(ctx, filter, "img_9.bad"))
	require.False(t, fileStore.Exists(ctx, filter, "img.mp3"))

	var urls []model.Url
	var err error

	urls, err = fileStore.List(ctx, filter)
	require.NoError(t, err)
	require.Equal(t, []model.Url{
		model.Path("img_1.txt"),
		model.Path("sub", "img_2.txt"),
	}, urls)
}
