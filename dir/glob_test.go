package dir_test

import (
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
// GlobMatcher.Walk()

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
	basepath, cleanup, err := setupTestFolder()
	assert.That(err, p.IsNoError())
	defer cleanup()

	for _, tc := range tcs {
		t.Run(tc.pattern, func(t *testing.T) {
			assert := asserter.New(t, asserter.AbortOnError())
			assert.That(true, p.IsTrue())

			pattern := path.Join(basepath, tc.pattern)
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
