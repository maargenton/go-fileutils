//go:build windows
// +build windows

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

func TestWindowsExpandPath(t *testing.T) {
	var cleanup = setupTestEnv(map[string]string{
		"FOOBAR":    "foo/bar/foobar",
		"FOOBARABS": "C:/foo/bar/foobar",
	})
	defer cleanup()

	var tcs = []struct {
		input, output string
	}{
		{"C:/.alek", "C:/.alek"},
		{"C:/foo/bar/foobar", "C:/foo/bar/foobar"},

		{"$FOOBARABS/.alek", "C:/foo/bar/foobar/.alek"},
		{"${FOOBARABS}/.alek", "C:/foo/bar/foobar/.alek"},

		{"C:/tmp/$FOOBAR/.alek", "C:/tmp/foo/bar/foobar/.alek"},
		{"C:/tmp/${FOOBAR}/.alek", "C:/tmp/foo/bar/foobar/.alek"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("ExpandPath(%#+v)", tc.input), func(t *testing.T) {
			output, err := fileutils.ExpandPath(tc.input)
			require.That(t, err).IsNil()
			require.That(t, output).Eq(tc.output)
		})
	}
}

func TestWindowsExpandPathFromHome(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
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

func TestWindowsExpandPathFromPwd(t *testing.T) {
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

func TestWindowsExpandPathRelative(t *testing.T) {
	var tcs = []struct {
		input, basepath, output string
	}{
		{".alek", "C:/usr/local/share", "C:/usr/local/share/.alek"},
		{"foo/bar/foobar", "C:/usr/local/share", "C:/usr/local/share/foo/bar/foobar"},
		{"C:/foo/bar/foobar", "C:/usr/local/share", "C:/foo/bar/foobar"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("ExpandPathRelative(%#+v,%#+v)", tc.input, tc.basepath), func(t *testing.T) {
			output, err := fileutils.ExpandPathRelative(tc.input, tc.basepath)
			expected := tc.output

			require.That(t, err).IsNil()
			require.That(t, output).Eq(expected)
		})
	}
}

func TestWindowsExpandPathRelativeFromPwd(t *testing.T) {
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
