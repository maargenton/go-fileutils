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

func TestWindowsIsDirectoryName(t *testing.T) {
	var tcs = []struct {
		path     string
		expected bool
	}{
		{"", true},
		{"\\", true},
		{".", true},
		{".\\", true},
		{"..", true},
		{"..\\", true},
		{"foo", false},
		{"foo\\", true},
		{"foo\\bar", false},
		{"foo\\bar\\", true},
		{"foo\\bar\\.", true},
		{"foo\\bar\\..", true},
		{"foo\\bar\\baz.", false},
		{"foo\\bar\\baz..", false},
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

// fileutils.IsDirectoryName
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.Clean

func TestWindowsClean(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"\\", "/"},
		{"\\\\", "/"},
		{"\\dev\\", "/dev/"},
		{".\\abc\\", "abc/"},
		{".\\abc\\\\def", "abc/def"},
		{"aaa\\..", "./"},
		{"aaa\\..\\", "./"},
		{"aaa\\.", "aaa/"},
		{"aaa\\.\\", "aaa/"},
		{".", "./"},
		{".\\", "./"},
		{"", "./"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Given %v", tc.input), func(t *testing.T) {
			t.Run("when calling Clean", func(t *testing.T) {
				output := fileutils.Clean(tc.input)
				t.Run("then output match expected", func(t *testing.T) {
					require.That(t, output).Eq(tc.output)
				})
			})
		})
	}
}

// fileutils.Clean
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// fileutils.Base, fileutils.Dir, fileutils.Split

func TestWindowsSplit(t *testing.T) {
	var tcs = []struct {
		path, dir, base string
	}{
		{"C:\\Program Files\\Microsoft\\", "C:/Program Files/", "Microsoft/"},
		{"\\\\hostname\\volume\\path\\to\\file", "//hostname/volume/path/to/", "file"},
		{"\\\\hostname\\volume\\path\\to\\dir\\", "//hostname/volume/path/to/", "dir/"},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Test Split(%#+v)", tc.path), func(t *testing.T) {
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
// fileutils.IsAbs

func TestWindowsIsAbs(t *testing.T) {
	var tcs = []struct {
		input string
		abs   bool
	}{
		{"path/to/file", false},
		{"/path/to/file", false},
		{"c:/path/to/file", true},
		{"//hostname/volume/path/to/file", true},
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

func TestWindowsJoin(t *testing.T) {
	var tcs = []struct {
		input  []string
		output string
	}{
		{[]string{"aaa\\bbb", "ccc"}, "aaa/bbb/ccc"},
		{[]string{"aaa\\bbb\\", "ccc"}, "aaa/bbb/ccc"},
		{[]string{"aaa\\bbb\\", "ccc\\"}, "aaa/bbb/ccc/"},
		{[]string{"c:\\dev", "tty.usbserial-1240"}, "c:/dev/tty.usbserial-1240"},
		{[]string{"prefix\\path\\to\\dir\\", "c:\\dev", "tty.usbserial-1240"}, "c:/dev/tty.usbserial-1240"},
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

func TestWindowsToNative(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"path/to/file", "path\\to\\file"},
		{"path/to/dir/", "path\\to\\dir\\"},
		{"/path/to/file", "\\path\\to\\file"},
		{"/path/to/dir/", "\\path\\to\\dir\\"},
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

func TestWindowsToSlash(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"path\\to\\file", "path/to/file"},
		{"path\\to\\dir\\", "path/to/dir/"},
		{"\\path\\to\\file", "/path/to/file"},
		{"\\path\\to\\dir\\", "/path/to/dir/"},
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

func TestWindowsVolumeName(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"c:/path/to/file", "c:"},
		{"c:\\path\\to\\file", "c:"},
		{"//hostname/volume/path/to/file", "//hostname/volume"},
		{"\\\\hostname\\volume\\path\\to\\file", "//hostname/volume"},
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
