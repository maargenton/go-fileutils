package fileutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

// ---------------------------------------------------------------------------
// RewriteFilename

func TestRewriteFilenameFull(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var input = "path/to/file.txt"
	var output = fileutil.RewriteFilename(input, &fileutil.RewriteOpts{
		Dirname: "other/path/to",
		Prefix:  "prefix-",
		Suffix:  "-suffix",
		Extname: ".csv",
	})
	assert.That(output, p.Eq("other/path/to/prefix-file-suffix.csv"))
}

func TestRewriteFilenameNoDotExt(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var input = "path/to/file.txt"
	var output = fileutil.RewriteFilename(input, &fileutil.RewriteOpts{
		Dirname: "other/path/to",
		Prefix:  "prefix-",
		Suffix:  "-suffix",
		Extname: "csv",
	})
	assert.That(output, p.Eq("other/path/to/prefix-file-suffix.csv"))
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
		for k, _ := range env {
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
			assert := asserter.New(t)

			output, err := fileutil.ExpandPath(tc.input)
			assert.That(err, p.IsNoError())
			assert.That(output, p.Eq(tc.output))
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
			assert := asserter.New(t)

			output, err := fileutil.ExpandPath(tc.input)
			expected := filepath.Join(home, tc.output)

			assert.That(err, p.IsNoError())
			assert.That(output, p.Eq(expected))
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
			assert := asserter.New(t)

			output, err := fileutil.ExpandPath(tc.input)
			expected := filepath.Join(pwd, tc.output)

			assert.That(err, p.IsNoError())
			assert.That(output, p.Eq(expected))
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
			assert := asserter.New(t)

			output, err := fileutil.ExpandPathRelative(tc.input, tc.basepath)
			expected := tc.output //filepath.Join(pwd, tc.output)

			assert.That(err, p.IsNoError())
			assert.That(output, p.Eq(expected))
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
			assert := asserter.New(t)

			output, err := fileutil.ExpandPathRelative(tc.input, tc.basepath)
			expected := filepath.Join(pwd, tc.output)

			assert.That(err, p.IsNoError())
			assert.That(output, p.Eq(expected))
		})
	}
}

// ExpandPathRelative
// ---------------------------------------------------------------------------
