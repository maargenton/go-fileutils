package fileutil_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

// ---------------------------------------------------------------------------
// fileutil.Clean

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
		{".", "./"},
		{"./", "./"},
		{"", "./"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Given %v", tc.input), func(t *testing.T) {
			t.Run("when calling Clean", func(t *testing.T) {
				output := fileutil.Clean(tc.input)
				t.Run("then output match expected", func(t *testing.T) {
					require.That(t, output).Eq(tc.output)
				})
			})
		})
	}
}

// fileutil.Clean
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutil.Rel
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
		t.Run(fmt.Sprintf("Given %v + %v", tc.basepath, tc.targetpath), func(t *testing.T) {
			t.Run("when calling Clean", func(t *testing.T) {
				output, err := fileutil.Rel(tc.basepath, tc.targetpath)
				require.That(t, err).IsNil()
				t.Run("then output match expected", func(t *testing.T) {
					require.That(t, output).Eq(tc.output)
				})
			})
		})
	}
}

// fileutil.Rel
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutil.Join

func TestJoin(t *testing.T) {
	var tcs = []struct {
		input  []string
		output string
	}{
		{[]string{"aaa/bbb", "ccc"}, "aaa/bbb/ccc"},
		{[]string{"aaa/bbb/", "ccc"}, "aaa/bbb/ccc"},
		{[]string{"aaa/bbb/", "ccc/"}, "aaa/bbb/ccc/"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Given %v", tc.input), func(t *testing.T) {
			t.Run("when calling Join", func(t *testing.T) {
				output := fileutil.Join(tc.input...)
				t.Run("then output match expected", func(t *testing.T) {
					require.That(t, output).Eq(tc.output)
				})
			})
		})
	}
}

// fileutil.Join
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// RewriteFilename

func TestRewriteFilenameFull(t *testing.T) {
	var input = "path/to/file.txt"
	var output = fileutil.RewriteFilename(input, &fileutil.RewriteOpts{
		Dirname: "other/path/to/",
		Prefix:  "prefix-",
		Suffix:  "-suffix",
		Extname: ".csv",
	})
	require.That(t, output).Eq("other/path/to/prefix-file-suffix.csv")
}

func TestRewriteFilenameNoDotExt(t *testing.T) {
	var input = "path/to/file.txt"
	var output = fileutil.RewriteFilename(input, &fileutil.RewriteOpts{
		Dirname: "other/path/to",
		Prefix:  "prefix-",
		Suffix:  "-suffix",
		Extname: "csv",
	})
	require.That(t, output).Eq("other/path/to/prefix-file-suffix.csv")
}

// RewriteFilename
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// ExpandPath

func setupTestEnv(env map[string]string) func() {
	for k, v := range env {
		os.Setenv(k, v)
	}
	return func() {
		for k := range env {
			os.Unsetenv(k)
		}
	}
}

func TestExpandPath(t *testing.T) {
	var cleanup = setupTestEnv(map[string]string{
		"FOOBAR":    "foo/bar/foobar",
		"FOOBARABS": "/foo/bar/foobar",
	})
	defer cleanup()

	var tcs = []struct{ input, output string }{
		{"/.alek", "/.alek"},
		{"/foo/bar/foobar", "/foo/bar/foobar"},

		{"$FOOBARABS/.alek", "/foo/bar/foobar/.alek"},
		{"${FOOBARABS}/.alek", "/foo/bar/foobar/.alek"},

		{"/tmp/$FOOBAR/.alek", "/tmp/foo/bar/foobar/.alek"},
		{"/tmp/${FOOBAR}/.alek", "/tmp/foo/bar/foobar/.alek"},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			output, err := fileutil.ExpandPath(tc.input)
			require.That(t, err).IsNil()
			require.That(t, output).Eq(tc.output)
		})
	}
}

func TestExpandPathFromHome(t *testing.T) {
	var tcs = []struct{ input, output string }{
		{"~", ""},
		{"~/", ""},
		{"~/.alek", ".alek"},
	}

	var home, _ = os.UserHomeDir()
	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			output, err := fileutil.ExpandPath(tc.input)
			expected := filepath.Join(home, tc.output)

			require.That(t, err).IsNil()
			require.That(t, output).Eq(expected)
		})
	}
}

func TestExpandPathFromPwd(t *testing.T) {
	var tcs = []struct{ input, output string }{
		{".alek", ".alek"},
		{"foo/bar/foobar", "foo/bar/foobar"},
	}

	var pwd, _ = os.Getwd()
	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			output, err := fileutil.ExpandPath(tc.input)
			expected := filepath.Join(pwd, tc.output)

			require.That(t, err).IsNil()
			require.That(t, output).Eq(expected)
		})
	}
}

// ExpandPath
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// ExpandPathRelative

func TestExpandPathRelative(t *testing.T) {
	var tcs = []struct{ input, basepath, output string }{
		{".alek", "/usr/local/share", "/usr/local/share/.alek"},
		{"foo/bar/foobar", "/usr/local/share", "/usr/local/share/foo/bar/foobar"},
		{"/foo/bar/foobar", "/usr/local/share", "/foo/bar/foobar"},
	}

	// var pwd, _ = os.Getwd()
	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			output, err := fileutil.ExpandPathRelative(tc.input, tc.basepath)
			expected := tc.output //filepath.Join(pwd, tc.output)

			require.That(t, err).IsNil()
			require.That(t, output).Eq(expected)
		})
	}
}

func TestExpandPathRelativeFromPwd(t *testing.T) {
	var tcs = []struct{ input, basepath, output string }{
		{".alek", "build/darwin-amd64", "build/darwin-amd64/.alek"},
		{"foo/bar/foobar", "build/darwin-amd64", "build/darwin-amd64/foo/bar/foobar"},
	}

	var pwd, _ = os.Getwd()
	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			output, err := fileutil.ExpandPathRelative(tc.input, tc.basepath)
			expected := filepath.Join(pwd, tc.output)

			require.That(t, err).IsNil()
			require.That(t, output).Eq(expected)
		})
	}
}

// ExpandPathRelative
// ---------------------------------------------------------------------------
