//go:build !windows
// +build !windows

package fileutils_test

import (
	"fmt"
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

// ---------------------------------------------------------------------------
// fileutils.IsAbs

func TestUnixIsAbs(t *testing.T) {
	var tcs = []struct {
		input string
		abs   bool
	}{
		{"path/to/file", false},
		{"/path/to/file", true},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("IsAbs(%#+v)", tc.input), func(t *testing.T) {
			abs := fileutils.IsAbs(tc.input)
			require.That(t, abs).Eq(tc.abs)
		})
	}
}

// fileutils.IsAbs
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.Join

func TestUnixJoin(t *testing.T) {
	var tcs = []struct {
		input  []string
		output string
	}{
		{[]string{"/dev", "tty.usbserial-1240"}, "/dev/tty.usbserial-1240"},
		{[]string{".", "/dev", "tty.usbserial-1240"}, "/dev/tty.usbserial-1240"},
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
// fileutils.ToNative

func TestUnixToNative(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"path/to/file", "path/to/file"},
		{"path/to/dir/", "path/to/dir/"},
		{"/path/to/file", "/path/to/file"},
		{"/path/to/dir/", "/path/to/dir/"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("ToNative(%#+v)", tc.input), func(t *testing.T) {
			output := fileutils.ToNative(tc.input)
			require.That(t, output).Eq(tc.output)
		})
	}
}

// fileutils.ToNative
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.ToSlash

func TestUnixToSlash(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"path/to/file", "path/to/file"},
		{"path/to/dir/", "path/to/dir/"},
		{"/path/to/file", "/path/to/file"},
		{"/path/to/dir/", "/path/to/dir/"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("ToSlash(%#+v)", tc.input), func(t *testing.T) {
			output := fileutils.ToSlash(tc.input)
			require.That(t, output).Eq(tc.output)
		})
	}
}

// fileutils.ToSlash
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.VolumeName

func TestUnixVolumeName(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"path/to/file", ""},
		{"path/to/dir/", ""},
		{"/path/to/file", ""},
		{"/path/to/dir/", ""},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("VolumeName(%#+v)", tc.input), func(t *testing.T) {
			output := fileutils.VolumeName(tc.input)
			require.That(t, output).Eq(tc.output)
		})
	}
}

// fileutils.VolumeName
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// Filename manipulation function that might need to access the underlying
// filesystem to evaluate their result.
// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.Abs
func TestUnixAbs(t *testing.T) {
	var tcs = []struct {
		input, prefix, suffix string
	}{
		{"/tmp/foo", "/tmp/foo", "/tmp/foo"},
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
