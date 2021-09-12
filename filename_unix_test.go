//go:build !windows

package fileutils_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-fileutils"
)

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

func TestUnixExpandPath(t *testing.T) {
	var cleanup = setupTestEnv(map[string]string{
		"FOOBAR":    "foo/bar/foobar",
		"FOOBARABS": "/foo/bar/foobar",
	})
	defer cleanup()

	var tcs = []struct {
		input, output string
	}{
		{"/.alek", "/.alek"},
		{"/foo/bar/foobar", "/foo/bar/foobar"},

		{"$FOOBARABS/.alek", "/foo/bar/foobar/.alek"},
		{"${FOOBARABS}/.alek", "/foo/bar/foobar/.alek"},

		{"/tmp/$FOOBAR/.alek", "/tmp/foo/bar/foobar/.alek"},
		{"/tmp/${FOOBAR}/.alek", "/tmp/foo/bar/foobar/.alek"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("ExpandPath(%#+v)", tc.input), func(t *testing.T) {
			output, err := fileutils.ExpandPath(tc.input)
			require.That(t, err).IsNil()
			require.That(t, output).Eq(tc.output)
		})
	}
}

func TestUnixExpandPathFromHome(t *testing.T) {
	var tcs = []struct{ input, output string }{
		{"~", ""},
		{"~/", ""},
		{"~/.alek", ".alek"},
	}

	var home, _ = os.UserHomeDir()
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("ExpandPath(%#+v)", tc.input), func(t *testing.T) {
			output, err := fileutils.ExpandPath(tc.input)
			expected := fileutils.Join(home, tc.output)

			require.That(t, err).IsNil()
			require.That(t, output).Eq(expected)
		})
	}
}

func TestUnixExpandPathFromPwd(t *testing.T) {
	var tcs = []struct{ input, output string }{
		{".alek", ".alek"},
		{"foo/bar/foobar", "foo/bar/foobar"},
	}

	var pwd, _ = os.Getwd()
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("ExpandPath(%#+v)", tc.input), func(t *testing.T) {
			output, err := fileutils.ExpandPath(tc.input)
			expected := fileutils.Join(pwd, tc.output)

			require.That(t, err).IsNil()
			require.That(t, output).Eq(expected)
		})
	}
}

// ExpandPath
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// ExpandPathRelative

func TestUnixExpandPathRelative(t *testing.T) {
	var tcs = []struct {
		input, basepath, output string
	}{
		{".alek", "/usr/local/share", "/usr/local/share/.alek"},
		{"foo/bar/foobar", "/usr/local/share", "/usr/local/share/foo/bar/foobar"},
		{"/foo/bar/foobar", "/usr/local/share", "/foo/bar/foobar"},
	}

	// var pwd, _ = os.Getwd()
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("ExpandPathRelative(%#+v,%#+v)", tc.input, tc.basepath), func(t *testing.T) {
			output, err := fileutils.ExpandPathRelative(tc.input, tc.basepath)
			expected := tc.output //fileutils.Join(pwd, tc.output)

			require.That(t, err).IsNil()
			require.That(t, output).Eq(expected)
		})
	}
}

func TestUnixExpandPathRelativeFromPwd(t *testing.T) {
	var tcs = []struct {
		input, basepath, output string
	}{
		{".alek", "build/darwin-amd64", "build/darwin-amd64/.alek"},
		{"foo/bar/foobar", "build/darwin-amd64", "build/darwin-amd64/foo/bar/foobar"},
	}

	var pwd, _ = os.Getwd()
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("ExpandPathRelative(%#+v,%#+v)", tc.input, tc.basepath), func(t *testing.T) {
			output, err := fileutils.ExpandPathRelative(tc.input, tc.basepath)
			expected := fileutils.Join(pwd, tc.output)

			require.That(t, err).IsNil()
			require.That(t, output).Eq(expected)
		})
	}
}

// ExpandPathRelative
// ---------------------------------------------------------------------------
