package dir_test

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/subexpr"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-fileutils/pkg/dir"
)

// ---------------------------------------------------------------------------
// dir.NewGlobMatcher()

func TestNewGlobMatcherError(t *testing.T) {
	pattern := `**/src/*.{c,cc,cpp`
	m, err := dir.NewGlobMatcher(pattern)
	require.That(t, err).IsNotNil()
	require.That(t, m).IsNil()
}

// dir.NewGlobMatcher()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// GlobMatcher.Match()

func TestGlobMatcherMatch(t *testing.T) {
	pattern := `content/**/src/**/*.{c,cc,cpp,h,hh,hpp}`
	m, err := dir.NewGlobMatcher(pattern)
	require.That(t, err).IsNil()
	require.That(t, m).IsNotNil()
	require.That(t, m.Match("aaa/bbb/src/ccc/ddd/something.cpp")).IsFalse()
	require.That(t, m.Match("content/aaa/bbb/src/ccc/ddd/something.cpp")).IsTrue()
	require.That(t, m.Match("content/aaa/bbb/src/ccc/ddd/")).IsFalse()
}

func TestGlobMatcherMatchWithWildcardStart(t *testing.T) {
	pattern := `**/src/**/*.{c,cc,cpp,h,hh,hpp}`
	m, err := dir.NewGlobMatcher(pattern)
	require.That(t, err).IsNil()
	require.That(t, m).IsNotNil()
	require.That(t, m.Match("aaa/bbb/src/ccc/ddd/something.cpp")).IsTrue()
	require.That(t, m.Match("content/aaa/bbb/src/ccc/ddd/something.cpp")).IsTrue()
	require.That(t, m.Match("content/aaa/bbb/src/ccc/ddd/")).IsFalse()
}

// GlobMatcher.Match()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// GlobMatcher.PrefixMatch()

func TestGlobMatcherPrefixMatch(t *testing.T) {
	pattern := `content/**/src/**/*.{c,cc,cpp,h,hh,hpp}`
	m, err := dir.NewGlobMatcher(pattern)
	require.That(t, err).IsNil()
	require.That(t, m).IsNotNil()

	require.That(t, m.PrefixMatch("content")).IsTrue()
	require.That(t, m.PrefixMatch("content/")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/bbb")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/bbb/")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/bbb/src")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/bbb/src/")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/bbb/src/ccc")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/bbb/src/ccc/")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/bbb/src/ccc/ddd")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/bbb/src/ccc/ddd/")).IsTrue()
	require.That(t, m.PrefixMatch("content/aaa/bbb/src/ccc/ddd/something.cpp")).IsTrue()
}

func TestGlobMatcherPrefixMatchNoMatch(t *testing.T) {
	pattern := `src`
	m, err := dir.NewGlobMatcher(pattern)
	require.That(t, err).IsNil()
	require.That(t, m).IsNotNil()

	require.That(t, m.PrefixMatch("src")).IsTrue()
	require.That(t, m.PrefixMatch("src/")).IsTrue()
	require.That(t, m.PrefixMatch("src/ddd")).IsFalse()
	require.That(t, m.PrefixMatch("src/ddd/")).IsFalse()
	require.That(t, m.PrefixMatch("dst/ddd")).IsFalse()
	require.That(t, m.PrefixMatch("dst/ddd/")).IsFalse()
}

func TestGlobMatcherPrefixMatchWithLeadingWildcardAlwaysMatch(t *testing.T) {
	pattern := `**/src`
	m, err := dir.NewGlobMatcher(pattern)
	require.That(t, err).IsNil()
	require.That(t, m).IsNotNil()

	require.That(t, m.PrefixMatch("src")).IsTrue()
	require.That(t, m.PrefixMatch("src/")).IsTrue()
	require.That(t, m.PrefixMatch("src/ddd")).IsTrue()
	require.That(t, m.PrefixMatch("src/ddd/")).IsTrue()
}

// GlobMatcher.Match()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// dir.Glob()

func TestGlob(t *testing.T) {
	var tcs = []struct {
		pattern string
		count   int
	}{
		{`**/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
		{`src/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
		{`src/**/*_test.{c,cc,cpp}`, 4},
		{`src/**/*_test.{h,hh,hpp}`, 0},
		{`src/**/*.{h,cpp}`, 12},
	}

	basepath, cleanup, err := setupTestFolder()
	require.That(t, err).IsNil()
	defer cleanup()

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			pattern := path.Join(basepath, tc.pattern)
			matches, err := dir.Glob(pattern)

			require.That(t, err).IsNil()
			require.That(t, matches).Length().Eq(tc.count)
			require.That(t, matches).All(
				subexpr.Value().StartsWith(basepath))
		})
	}
}

func TestGlobStarFromCurrentDirectory(t *testing.T) {
	basepath, cleanup, err := setupTestFolder()
	require.That(t, err).IsNil()
	defer cleanup()

	pwd, _ := os.Getwd()
	os.Chdir(basepath)
	defer func() {
		os.Chdir(pwd)
	}()

	matches, err := dir.Glob("*")
	require.That(t, err).IsNil()
	require.That(t, matches).IsEqualSet([]string{"src/"})
}

func TestGlobFromSystemRoot(t *testing.T) {
	matches, err := dir.Glob("/d*")
	require.That(t, err).IsNil()
	require.That(t, matches).Contains([]string{"/dev/"})
}

// func TestGlobFromSystemRoot2(t *testing.T) {
// 	matches, err := dir.Glob("/dev/std*")
// 	require.That(t, err).IsNil()
// 	require.That(t, matches).Contains([]string{"/dev/"})
// }

// dir.Glob()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// dir.GlobFrom()

func TestGlobMatcherGlobFrom(t *testing.T) {
	var tcs = []struct {
		pattern string
		count   int
	}{
		{`**/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
		{`src/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
		{`src/**/*_test.{c,cc,cpp}`, 4},
		{`src/**/*_test.{h,hh,hpp}`, 0},
		{`src/**/*.{h,cpp}`, 12},
	}
	basepath, cleanup, err := setupTestFolder()
	require.That(t, err).IsNil()
	defer cleanup()

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			matches, err := dir.GlobFrom(basepath, tc.pattern)

			require.That(t, err).IsNil()
			require.That(t, matches).Length().Eq(tc.count)
			require.That(t, matches).All(
				subexpr.Value().StartsWith("src"))
		})
	}
}

// dir.GlobFrom()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// dir.Glob()

func TestScan(t *testing.T) {
	var tcs = []struct {
		pattern string
		count   int
	}{
		{`**/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
		{`src/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
		{`src/**/*_test.{c,cc,cpp}`, 4},
		{`src/**/*_test.{h,hh,hpp}`, 0},
		{`src/**/*.{h,cpp}`, 12},
	}
	basepath, cleanup, err := setupTestFolder()
	require.That(t, err).IsNil()
	defer cleanup()

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			var count = 0
			var countingWalk = func(path string, d fs.DirEntry, err error) error {
				count++
				return nil
			}

			err = dir.Scan(path.Join(basepath, tc.pattern), countingWalk)

			require.That(t, err).IsNil()
			require.That(t, count).Eq(tc.count)
		})
	}
}

// dir.Scan()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// GlobMatcher.ScanFrom()

func TestGlobMatcherScanFrom(t *testing.T) {
	var tcs = []struct {
		pattern string
		count   int
	}{
		{`**/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
		{`src/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
		{`src/**/*_test.{c,cc,cpp}`, 4},
		{`src/**/*_test.{h,hh,hpp}`, 0},
		{`src/**/*.{h,cpp}`, 12},
	}
	basepath, cleanup, err := setupTestFolder()
	require.That(t, err).IsNil()
	defer cleanup()

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			var count = 0
			var countingWalk = func(path string, d fs.DirEntry, err error) error {
				count++
				return nil
			}

			err = dir.ScanFrom(basepath, tc.pattern, countingWalk)

			require.That(t, err).IsNil()
			require.That(t, count).Eq(tc.count)
		})
	}
}

// GlobMatcher.ScanFrom()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Test error path for dir.Glob, dir.GlobFrom, dir.Scan, dir.ScanFrom

func TestGlobFunctionsErrorWithBadPattern(t *testing.T) {
	basepath, cleanup, err := setupTestFolder()
	require.That(t, err).IsNil()
	defer cleanup()

	var dummyWalk = func(path string, d fs.DirEntry, err error) error {
		return err
	}

	_, err = dir.Glob(path.Join(basepath, "src/**/*/{"))
	require.That(t, err).IsNotNil()

	_, err = dir.GlobFrom(basepath, "src/**/*/{")
	require.That(t, err).IsNotNil()

	err = dir.Scan(path.Join(basepath, "src/**/*/{"), dummyWalk)
	require.That(t, err).IsNotNil()

	err = dir.ScanFrom(basepath, "src/**/*/{", dummyWalk)
	require.That(t, err).IsNotNil()
}

// Test error path for dir.Glob, dir.GlobFrom, dir.Scan, dir.ScanFrom
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func setupTestFolder() (basepath string, cleanup func(), err error) {
	basepath, err = ioutil.TempDir(".", "testdata-")
	basepath = fileutils.Clean(basepath)
	cleanup = func() {
		if basepath != "" {
			os.RemoveAll(basepath)
		}
	}
	if err != nil {
		return
	}

	var filenames []string
	for _, n := range []string{"foo", "bar", "aaa", "bbb"} {
		filenames = append(filenames,
			fileutils.Join(basepath, "src", n, n+".h"),
			fileutils.Join(basepath, "src", n, n+".cpp"),
			fileutils.Join(basepath, "src", n, n+"_test.cpp"),
		)
	}
	err = fileutils.Touch(filenames...)
	return
}
