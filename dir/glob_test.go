package dir_test

import (
	"os"
	"path"
	"testing"

	"github.com/maargenton/fileutil/dir"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

// ---------------------------------------------------------------------------
// dir.NewGlobMatcher()

func TestNewGlobMatcherError(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	assert.That(nil, p.IsNil())

	pattern := `**/src/*.{c,cc,cpp`
	g, err := dir.NewGlobMatcher(pattern)
	assert.That(err, p.IsNotNil())
	assert.That(g, p.IsNil())
}

// dir.NewGlobMatcher()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// GlobMatcher.Match()

func TestGlobMatcherMatch(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	assert.That(nil, p.IsNil())

	pattern := `content/**/src/**/*.{c,cc,cpp,h,hh,hpp}`
	g, err := dir.NewGlobMatcher(pattern)
	assert.That(err, p.IsNoError())
	assert.That(g, p.IsNotNil())
	assert.That(g.Match("aaa/bbb/src/ccc/ddd/something.cpp"), p.IsFalse())
	assert.That(g.Match("content/aaa/bbb/src/ccc/ddd/something.cpp"), p.IsTrue())
	assert.That(g.Match("content/aaa/bbb/src/ccc/ddd/"), p.IsFalse())
}

func TestGlobMatcherMatchWithWildcardStart(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	assert.That(nil, p.IsNil())

	pattern := `**/src/**/*.{c,cc,cpp,h,hh,hpp}`
	g, err := dir.NewGlobMatcher(pattern)
	assert.That(err, p.IsNoError())
	assert.That(g, p.IsNotNil())
	assert.That(g.Match("aaa/bbb/src/ccc/ddd/something.cpp"), p.IsTrue())
	assert.That(g.Match("content/aaa/bbb/src/ccc/ddd/something.cpp"), p.IsTrue())
	assert.That(g.Match("content/aaa/bbb/src/ccc/ddd/"), p.IsFalse())
}

// GlobMatcher.Match()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// GlobMatcher.PrefixMatch()

func TestGlobMatcherPrefixMatch(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	assert.That(nil, p.IsNil())

	pattern := `content/**/src/**/*.{c,cc,cpp,h,hh,hpp}`
	g, err := dir.NewGlobMatcher(pattern)
	assert.That(err, p.IsNoError())
	assert.That(g, p.IsNotNil())

	assert.That(g.PrefixMatch("content"), p.IsTrue())
	assert.That(g.PrefixMatch("content/"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/bbb"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/bbb/"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/bbb/src"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/bbb/src/"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/bbb/src/ccc"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/bbb/src/ccc/"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/bbb/src/ccc/ddd"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/bbb/src/ccc/ddd/"), p.IsTrue())
	assert.That(g.PrefixMatch("content/aaa/bbb/src/ccc/ddd/something.cpp"), p.IsTrue())
}

func TestGlobMatcherPrefixMatchNoMatch(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	assert.That(nil, p.IsNil())

	pattern := `src`
	g, err := dir.NewGlobMatcher(pattern)
	assert.That(err, p.IsNoError())
	assert.That(g, p.IsNotNil())

	assert.That(g.PrefixMatch("src"), p.IsTrue())
	assert.That(g.PrefixMatch("src/"), p.IsTrue())
	assert.That(g.PrefixMatch("src/ddd"), p.IsFalse())
	assert.That(g.PrefixMatch("src/ddd/"), p.IsFalse())
	assert.That(g.PrefixMatch("dst/ddd"), p.IsFalse())
	assert.That(g.PrefixMatch("dst/ddd/"), p.IsFalse())
}

func TestGlobMatcherPrefixMatchWithLeadingWildcardAlwaysMatch(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	assert.That(nil, p.IsNil())

	pattern := `**/src`
	g, err := dir.NewGlobMatcher(pattern)
	assert.That(err, p.IsNoError())
	assert.That(g, p.IsNotNil())

	assert.That(g.PrefixMatch("src"), p.IsTrue())
	assert.That(g.PrefixMatch("src/"), p.IsTrue())
	assert.That(g.PrefixMatch("src/ddd"), p.IsTrue())
	assert.That(g.PrefixMatch("src/ddd/"), p.IsTrue())
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

	assert := asserter.New(t, asserter.AbortOnError())
	basepath, cleanup, err := setupTestFolder()
	assert.That(err, p.IsNoError())
	defer cleanup()

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			assert := asserter.New(t, asserter.AbortOnError())
			assert.That(true, p.IsTrue())

			pattern := path.Join(basepath, tc.pattern)
			matches, err := dir.Glob(pattern)

			assert.That(err, p.IsNoError())
			assert.That(matches, p.Length(p.Eq(tc.count)))
			assert.That(matches, p.All(p.StartsWith(basepath)))
		})
	}
}

// dir.Glob()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// GlobMatcher.GlobFrom()

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

	assert := asserter.New(t, asserter.AbortOnError())
	basepath, cleanup, err := setupTestFolder()
	assert.That(err, p.IsNoError())
	defer cleanup()

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			assert := asserter.New(t, asserter.AbortOnError())
			assert.That(true, p.IsTrue())

			matches, err := dir.GlobFrom(basepath, tc.pattern)

			assert.That(err, p.IsNoError())
			assert.That(matches, p.Length(p.Eq(tc.count)))
			assert.That(matches, p.All(p.StartsWith("src")))
		})
	}
}

// GlobMatcher.GlobFrom()
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

	assert := asserter.New(t, asserter.AbortOnError())
	basepath, cleanup, err := setupTestFolder()
	assert.That(err, p.IsNoError())
	defer cleanup()

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			assert := asserter.New(t, asserter.AbortOnError())
			assert.That(true, p.IsTrue())

			var count = 0
			var countingWalk = func(path string, info os.FileInfo, err error) error {
				count++
				return nil
			}

			err = dir.Scan(path.Join(basepath, tc.pattern), countingWalk)

			assert.That(err, p.IsNoError())
			assert.That(count, p.Eq(tc.count))
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

	assert := asserter.New(t, asserter.AbortOnError())
	basepath, cleanup, err := setupTestFolder()
	assert.That(err, p.IsNoError())
	defer cleanup()

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			assert := asserter.New(t, asserter.AbortOnError())
			assert.That(true, p.IsTrue())

			var count = 0
			var countingWalk = func(path string, info os.FileInfo, err error) error {
				count++
				return nil
			}

			err = dir.ScanFrom(basepath, tc.pattern, countingWalk)

			assert.That(err, p.IsNoError())
			assert.That(count, p.Eq(tc.count))
		})
	}
}

// GlobMatcher.ScanFrom()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Test error path for dir.Glob, dir.GlobFrom, dir.Scan, dir.ScanFrom

func TestGlobFunctionsErrorWithBadPattern(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	basepath, cleanup, err := setupTestFolder()
	assert.That(err, p.IsNoError())
	defer cleanup()

	var dummyWalk = func(path string, info os.FileInfo, err error) error {
		return err
	}

	_, err = dir.Glob(path.Join(basepath, "src/**/*/{"))
	assert.That(err, p.IsNotNil())

	_, err = dir.GlobFrom(basepath, "src/**/*/{")
	assert.That(err, p.IsNotNil())

	err = dir.Scan(path.Join(basepath, "src/**/*/{"), dummyWalk)
	assert.That(err, p.IsNotNil())

	err = dir.ScanFrom(basepath, "src/**/*/{", dummyWalk)
	assert.That(err, p.IsNotNil())
}

// Test error path for dir.Glob, dir.GlobFrom, dir.Scan, dir.ScanFrom
// ---------------------------------------------------------------------------
