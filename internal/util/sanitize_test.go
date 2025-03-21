package util_test

import (
	"testing"

	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/stretchr/testify/require"
)

func TestSanitizeAlphaLower(t *testing.T) {
	tests := []struct {
		value         string
		wantSanitized string
	}{
		{"", ""},
		{"abc", "abc"},
		{"Aaa  bBB CCC", "aaa-bbb-ccc"},
		{" @ # Aa*a@@bBB && Ccçc eêEË  # ", "aa-a-bbb-ccc-eeee"},
	}

	for _, tt := range tests {
		t.Run("value["+tt.value+"]", func(t *testing.T) {
			gotSanitized := util.SanitizeAlphaLower(tt.value)
			require.Equal(t, tt.wantSanitized, gotSanitized)
		})
	}
}
