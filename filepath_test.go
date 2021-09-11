package fileutils_test

import (
	"fmt"
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

func TestIsDirectoryName(t *testing.T) {
	var tcs = []struct {
		path     string
		expected bool
	}{
		{"", true},
		{"/", true},
		{".", true},
		{"./", true},
		{"..", true},
		{"../", true},
		{"foo", false},
		{"foo/", true},
		{"foo/bar", false},
		{"foo/bar/", true},
		{"foo/bar/.", true},
		{"foo/bar/..", true},
		{"foo/bar/baz.", false},
		{"foo/bar/baz..", false},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Given %#+v", tc.path), func(t *testing.T) {
			t.Run("when calling IsDirectoryName", func(t *testing.T) {
				dir := fileutils.IsDirectoryName(tc.path)
				t.Run(fmt.Sprintf("then result is %#+v", tc.expected), func(t *testing.T) {
					require.That(t, dir,
						require.Context{Name: "path", Value: tc.path},
					).Eq(tc.expected)
				})
			})
		})
	}
}
