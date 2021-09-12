package fileutils_test

import (
	"fmt"
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

// ---------------------------------------------------------------------------
// fileutils.FromSlash

func TestUnixFromSlash(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"path/to/file", "path/to/file"},
		{"path/to/dir/", "path/to/dir/"},
		{"/path/to/file", "/path/to/file"},
		{"/path/to/dir/", "/path/to/dir/"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("FromSlash(%#+v)", tc.input), func(t *testing.T) {
			output := fileutils.FromSlash(tc.input)
			require.That(t, output).Eq(tc.output)
		})
	}
}

// fileutils.FromSlash
// ---------------------------------------------------------------------------

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
