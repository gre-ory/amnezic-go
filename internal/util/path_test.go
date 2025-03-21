package util_test

import (
	"path/filepath"
	"testing"

	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/stretchr/testify/require"
)

func TestCleanLocalPath(t *testing.T) {
	tests := []struct {
		path        string
		wantCleaned string
		wantErr     error
	}{
		{"../a/b/../..", "", util.ErrInvalidLocalPath},
		{"../../", "", util.ErrInvalidLocalPath},
		{"../..", "", util.ErrInvalidLocalPath},
		{"../", "", util.ErrInvalidLocalPath},
		{"..", "", util.ErrInvalidLocalPath},
		{"", ".", nil},
		{".", ".", nil},
		{"a", "a", nil},
		{"./a", "a", nil},
		{"a/", "a", nil},
		{"./a/", "a", nil},
		{"a/b", "a/b", nil},
		{"./a/b", "a/b", nil},
		{"./a/b/", "a/b", nil},
		{"./a/../../b/", "", util.ErrInvalidLocalPath},
		{"/a/b/c", "/a/b/c", nil},
	}

	for _, tt := range tests {
		t.Run("path["+tt.path+"]", func(t *testing.T) {
			gotCleaned, gotErr := util.CleanLocalPath(tt.path)
			require.Equal(t, tt.wantErr, gotErr)
			if tt.wantCleaned != "" {
				require.Equal(t, filepath.Clean(tt.wantCleaned), gotCleaned)
			} else {
				require.Equal(t, tt.wantCleaned, gotCleaned)
			}
		})
	}
}

func TestCleanExtension(t *testing.T) {
	tests := []struct {
		extension   string
		wantCleaned string
		wantErr     error
	}{
		{"", "", nil},
		{".", "", util.ErrInvalidExtension},
		{"..", "", util.ErrInvalidExtension},
		{"json", ".json", nil},
		{".json", ".json", nil},
		{"..json", "", util.ErrInvalidExtension},
		{".a.json", "", util.ErrInvalidExtension},
		{"a.json", "", util.ErrInvalidExtension},
	}

	for _, tt := range tests {
		t.Run("extension["+tt.extension+"]", func(t *testing.T) {
			gotCleaned, gotErr := util.CleanExtension(tt.extension)
			require.Equal(t, tt.wantErr, gotErr)
			require.Equal(t, tt.wantCleaned, gotCleaned)
		})
	}
}
