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
		{"./abc", "abc"},
		{"./abc/", "abc/"},
		{"./abc//def", "abc/def"},
		{"aaa/..", "./"},
		{"aaa/../", "./"},
		{"aaa/.", "aaa/"},
		{"aaa/./", "aaa/"},
		{".", "./"},
		{"./", "./"},
		{"", "./"},
		{"~", "~/"},

		{"./abc/def/*.go", "abc/def/*.go"},
		{"./**/file.go", "**/file.go"},
		{"./**/*.go", "**/*.go"},
		{"./*.go", "*.go"},
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

// ---------------------------------------------------------------------------
// fileutils.Join

func TestJoin(t *testing.T) {
	var tcs = []struct {
		input  []string
		output string
	}{
		{[]string{"aaa/bbb", "ccc"}, "aaa/bbb/ccc"},
		{[]string{"aaa/bbb/", "ccc"}, "aaa/bbb/ccc"},
		{[]string{"aaa/bbb/", "ccc/"}, "aaa/bbb/ccc/"},
		{[]string{"", ""}, "./"},
		{[]string{"aaa/bbb", ""}, "aaa/bbb/"},
		{[]string{"aaa/bbb", "../ccc", "../ddd"}, "aaa/ddd"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Join(%#+v)", tc.input), func(t *testing.T) {
			output := fileutils.Join(tc.input...)
			require.That(t, output).Eq(tc.output)
		})
	}
}

// fileutils.Join
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// Filename manipulation function that might need to access the underlying
// filesystem to evaluate their result.
// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.Abs
func TestAbs(t *testing.T) {
	var tcs = []struct {
		input, prefix, suffix string
	}{
		{"testdata/src", "", "testdata/src"},
		{"testdata/src/", "", "testdata/src/"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Abs(%#+v)", tc.input), func(t *testing.T) {
			output, err := fileutils.Abs(tc.input)
			require.That(t, err).IsNil()
			require.That(t, output).StartsWith(tc.prefix)
			require.That(t, output).EndsWith(tc.suffix)
		})
	}
}

// fileutils.Abs
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.EvalSymlinks
func TestEvalSymlinks(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"testdata/src", "testdata/src"},
		{"testdata/src/", "testdata/src/"},
		{"testdata/dst", "testdata/src"},
		{"testdata/dst/", "testdata/src/"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("EvalSymlinks(%#+v)", tc.input), func(t *testing.T) {
			output, err := fileutils.EvalSymlinks(tc.input)
			require.That(t, err).IsNil()
			require.That(t, output).Eq(tc.output)
		})
	}
}

// fileutils.EvalSymlinks
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.Rel
func TestRel(t *testing.T) {
	var tcs = []struct {
		basepath, targetpath, output string
	}{
		{"testdata", "testdata/src", "src"},
		{"testdata/", "testdata/src", "src"},
		{"testdata/", "testdata/src/", "src/"},
		{"/", "/testdata/src/", "testdata/src/"},
		{"/testdata", "/testdata/src/", "src/"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Rel(%#+v,%#+v)", tc.basepath, tc.targetpath), func(t *testing.T) {
			output, err := fileutils.Rel(tc.basepath, tc.targetpath)
			require.That(t, err).IsNil()
			require.That(t, output).Eq(tc.output)
		})
	}
}

// fileutils.Rel
// ---------------------------------------------------------------------------
