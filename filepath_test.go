package fileutils_test

import (
	"fmt"
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"
)

// ---------------------------------------------------------------------------
// fileutils.IsDirectoryName

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
		t.Run(fmt.Sprintf("IsDirectoryName(%#+v)", tc.path), func(t *testing.T) {
			dir := fileutils.IsDirectoryName(tc.path)
			require.That(t, dir).Eq(tc.expected)
		})
	}
}

// fileutils.IsDirectoryName
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.Clean

func TestClean(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"/", "/"},
		{"//", "/"},
		{"/dev/", "/dev/"},
		{"./abc/", "abc/"},
		{"./abc//def", "abc/def"},
		{"aaa/..", "./"},
		{"aaa/../", "./"},
		{"aaa/.", "aaa/"},
		{"aaa/./", "aaa/"},
		{".", "./"},
		{"./", "./"},
		{"", "./"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Clean(%#+v)", tc.input), func(t *testing.T) {
			output := fileutils.Clean(tc.input)
			require.That(t, output).Eq(tc.output)
		})
	}
}

// fileutils.Clean
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.Base, fileutils.Dir, fileutils.Split

func TestSplitBaeDir(t *testing.T) {
	var tcs = []struct {
		path, dir, base string
	}{
		// Matching filepath.Split()
		{"", "", ""},
		{"/dev", "/", "dev"},
		{"//", "/", ""},
		{"///", "/", ""},
		{"/foo/bar/baz", "/foo/bar/", "baz"},
		{"/foo/bar/baz/..", "/foo/bar/baz/", "../"},

		// Matching filepath.Split(), presered trailing sep
		{"/dev/", "/", "dev/"},
		{"/foo/bar/baz/", "/foo/bar/", "baz/"},
		{"/foo/bar/baz/../", "/foo/bar/baz/", "../"},

		// Relative paths
		{"path/to/file", "path/to/", "file"},
		{"file", "", "file"},

		// Other
		{"/foo/bar/../baz", "/foo/bar/../", "baz"},
		{"/foo/bar/../baz/", "/foo/bar/../", "baz/"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Split(%#+v)", tc.path), func(t *testing.T) {
			dir, base := fileutils.Split(tc.path)
			r := []string{dir, base}
			verify.That(t, r).Eq([]string{tc.dir, tc.base})
		})
	}

	// Dir() and Base() should return the same result as the first or second
	// result of Split(), respectively
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Dir(%#+v)", tc.path), func(t *testing.T) {
			dir := fileutils.Dir(tc.path)
			verify.That(t, dir).Eq(tc.dir)
		})
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Base(%#+v)", tc.path), func(t *testing.T) {
			base := fileutils.Base(tc.path)
			verify.That(t, base).Eq(tc.base)
		})
	}
}

// fileutils.Base, fileutils.Dir, fileutils.Split
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.Ext

func TestExt(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"foo.cpp", ".cpp"},
		{"foo.cpp.o", ".o"},
		{"foo.cpp.d/", ""},
		{"path/to/foo.cpp", ".cpp"},
		{"path/to/foo.cpp.o", ".o"},
		{"path/to/foo.cpp.d/", ""},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Ext(%#+v)", tc.input), func(t *testing.T) {
			output := fileutils.Ext(tc.input)
			require.That(t, output).Eq(tc.output)
		})
	}
}

// fileutils.Ext
// ---------------------------------------------------------------------------
