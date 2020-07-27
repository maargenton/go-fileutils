package dir

import (
	"regexp"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

// ---------------------------------------------------------------------------
// globFragment.match()

func TestLiteralGlobFragmentMatch(t *testing.T) {
	assert := asserter.New(t)

	var f = &globFragment{
		literal: "src",
	}
	assert.That(f.match("src"), p.IsTrue(), "fragment", "src")
	assert.That(f.match("src/"), p.IsTrue(), "fragment", "src/")
	assert.That(f.match("not-src"), p.IsFalse(), "fragment", "not-src")
	assert.That(f.match("not-src/"), p.IsFalse(), "fragment", "not-src/")
}

func TestRegexpGlobFragmentMatch(t *testing.T) {
	assert := asserter.New(t)

	var f = &globFragment{
		re: regexp.MustCompile(".*src.*"),
	}
	assert.That(f.match("src"), p.IsTrue(), "fragment", "src")
	assert.That(f.match("src/"), p.IsTrue(), "fragment", "src/")
	assert.That(f.match("not-src"), p.IsTrue(), "fragment", "not-src")
	assert.That(f.match("not-src/"), p.IsTrue(), "fragment", "not-src/")

	assert.That(f.match("dst"), p.IsFalse(), "fragment", "dst")
	assert.That(f.match("dst/"), p.IsFalse(), "fragment", "dst/")
}

// globFragment.match()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// globFragment.matchStart()

func TestGlobFragmentMatchStart(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var f = &globFragment{
		literal: "src",
	}

	assert.That(f.matchStart("src/aaa/bbb"), p.IsEqualSet([]string{"aaa/bbb"}))
	assert.That(f.matchStart("aaa/src/bbb"), p.IsEqualSet([]string{}))

}

func TestSubdirGlobFragmentMatchStart(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var f = &globFragment{
		subdir:  true,
		literal: "src",
	}

	assert.That(f.matchStart("src/aaa/bbb"), p.IsEqualSet([]string{"aaa/bbb"}))
	assert.That(f.matchStart("aaa/src/bbb"), p.IsEqualSet([]string{"bbb"}))
	assert.That(f.matchStart("aaa/src/bbb/src/ccc/ddd"), p.IsEqualSet([]string{
		"bbb/src/ccc/ddd",
		"ccc/ddd",
	}))

}

// globFragment.matchStart()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// globFragment.prefixMatchStart()

func TestGlobFragmentPrefixMatchStart(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var f = &globFragment{
		literal: "src",
	}

	assert.That(f.prefixMatchStart("src/aaa/bbb"), p.Eq([]string{"aaa/bbb"}))
	assert.That(f.prefixMatchStart("aaa/src/bbb"), p.Eq([]string{}))

}

func TestSubdirGlobFragmentPrefixMatchStart(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var f = &globFragment{
		subdir:  true,
		literal: "src",
	}

	assert.That(f.prefixMatchStart("src/aaa/bbb"), p.Eq([]string{"", "aaa/bbb"}))
	assert.That(f.prefixMatchStart("aaa/src/bbb"), p.Eq([]string{"", "bbb"}))
	assert.That(f.prefixMatchStart("aaa/src/bbb/src/ccc/ddd"), p.Eq([]string{"", "bbb/src/ccc/ddd", "ccc/ddd"}))
}

// globFragment.matchStart()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// splitPath()

func TestSplitPath(t *testing.T) {
	var tcs = []struct {
		input               string
		fragment, remainder string
	}{
		{"aaa/bbb/ccc", "aaa/", "bbb/ccc"},
		{"/aaa/bbb/ccc", "/", "aaa/bbb/ccc"},
		{"aaa", "aaa", ""},
		{"", "", ""},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			assert := asserter.New(t)
			a, b := splitPath(tc.input)
			assert.That(a, p.Eq(tc.fragment))
			assert.That(b, p.Eq(tc.remainder))
		})
	}
}

// splitPath()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// cleanFragment()

func TestCleanFragment(t *testing.T) {
	var tcs = []struct {
		input, output string
	}{
		{"/", "/"},
		{"aaa/bbb/", "aaa/bbb"},
		{"aaa/bbb", "aaa/bbb"},
		{"/aaa/bbb/", "/aaa/bbb"},
		{"/aaa/bbb", "/aaa/bbb"},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			assert := asserter.New(t)
			assert.That(cleanFragment(tc.input), p.Eq(tc.output))

		})
	}
}

// cleanFragment()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// globFragmentToRegexp()

func TestGlobFragmentToRegexp(t *testing.T) {
	var tcs = []struct {
		input, re, match string
	}{
		{`aaa.cpp`, `^aaa\.cpp$`, "aaa.cpp"},
		{`*.cpp`, `^.*\.cpp$`, "aaa.cpp"},
		{`aaa*.cpp`, `^aaa.*\.cpp$`, "aaabbb.cpp"},
		{`*a*.cpp`, `^.*a.*\.cpp$`, "babb.cpp"},
		{`?a*.cpp`, `^.a.*\.cpp$`, "babb.cpp"},
		{`[a-z]*.cpp`, `^[a-z].*\.cpp$`, "zbbb.cpp"},

		{`{a,b}`, `^(?:(?:a)|(?:b))$`, "a"},
		{`\{a,b}`, `^\{a,b}$`, "{a,b}"},
		{`*_test.{c,cc,cpp}`, `^.*_test\.(?:(?:c)|(?:cc)|(?:cpp))$`, "foo_test.cc"},
		{`\a\b\c\{\.`, `^abc\{\.$`, "abc{."},
		{`{,*_}main.cpp`, `^(?:(?:)|(?:.*_))main\.cpp$`, "main.cpp"},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			assert := asserter.New(t, asserter.AbortOnError())
			assert.That(true, p.IsTrue())

			re, err := globFragmentToRegexp(tc.input)
			assert.That(err, p.IsNoError())
			assert.That(re, p.IsNotNil())
			assert.That(re.String(), p.Eq(tc.re))
			assert.That(
				re.MatchString(tc.match), p.IsTrue(),
				"re", re.String(), "match", tc.match,
			)
		})
	}
}

func TestGlobFragmentToRegexpError(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var pattern = `*.{a,b`
	re, err := globFragmentToRegexp(pattern)
	assert.That(err, p.IsNotNil())
	assert.That(re, p.IsNil())
}

// globFragmentToRegexp()
// ---------------------------------------------------------------------------
