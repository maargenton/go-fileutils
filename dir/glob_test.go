package dir_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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
// GlobMatcher.Walk()

func touch(filename string) error {
	dirname := filepath.Dir(filename)
	if err := os.MkdirAll(dirname, 0777); err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	return f.Close()
}

func touchAll(filenames []string) error {

	for _, filename := range filenames {
		if err := touch(filename); err != nil {
			return err
		}
	}
	return nil
}

func setupTestFs(basepath string) error {
	var filenames []string
	for _, n := range []string{"foo", "bar", "aaa", "bbb"} {
		filenames = append(filenames,
			filepath.Join(basepath, "src", n, n+".h"),
			filepath.Join(basepath, "src", n, n+".cpp"),
			filepath.Join(basepath, "src", n, n+"_test.cpp"),
		)
	}
	return touchAll(filenames)
}

func TestGlobMatcherGlob(t *testing.T) {
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
	dirname, err := ioutil.TempDir(".", "testdata-")
	assert.That(err, p.IsNoError())
	defer os.RemoveAll(dirname) // clean up
	err = setupTestFs(dirname)
	assert.That(err, p.IsNoError())

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			assert := asserter.New(t, asserter.AbortOnError())
			assert.That(true, p.IsTrue())

			pattern := path.Join(dirname, tc.pattern)
			m, err := dir.NewGlobMatcher(pattern)
			assert.That(err, p.IsNoError())

			matches, err := m.Glob()
			assert.That(err, p.IsNoError())
			assert.That(matches, p.Length(p.Eq(tc.count)))

		})
	}
}

// GlobMatcher.Walk()
// ---------------------------------------------------------------------------
