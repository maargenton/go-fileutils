package dir

import (
	"regexp"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"
)

// ---------------------------------------------------------------------------
// globFragment.match()

func TestLiteralGlobFragmentMatch(t *testing.T) {
	var f = &globFragment{
		literal: "src",
	}
	verify.That(t, f.match("src")).IsTrue()
	verify.That(t, f.match("src/")).IsTrue()
	verify.That(t, f.match("not-src")).IsFalse()
	verify.That(t, f.match("not-src/")).IsFalse()
}

func TestRegexpGlobFragmentMatch(t *testing.T) {
	var f = &globFragment{
		re: regexp.MustCompile(".*src.*"),
	}
	verify.That(t, f.match("src")).IsTrue()
	verify.That(t, f.match("src/")).IsTrue()
	verify.That(t, f.match("not-src")).IsTrue()
	verify.That(t, f.match("not-src/")).IsTrue()

	verify.That(t, f.match("dst")).IsFalse()
	verify.That(t, f.match("dst/")).IsFalse()
}

// globFragment.match()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// globFragment.matchStart()

func TestGlobFragmentMatchStart(t *testing.T) {
	var f = &globFragment{
		literal: "src",
	}

	verify.That(t, f.matchStart("src/aaa/bbb")).IsEqualSet([]string{"aaa/bbb"})
	verify.That(t, f.matchStart("aaa/src/bbb")).IsEqualSet([]string{})
}

func TestSubdirGlobFragmentMatchStart(t *testing.T) {
	var f = &globFragment{
		subdir:  true,
		literal: "src",
	}

	verify.That(t, f.matchStart("src/aaa/bbb")).IsEqualSet([]string{"aaa/bbb"})
	verify.That(t, f.matchStart("aaa/src/bbb")).IsEqualSet([]string{"bbb"})
	verify.That(t, f.matchStart("aaa/src/bbb/src/ccc/ddd")).IsEqualSet([]string{
		"bbb/src/ccc/ddd",
		"ccc/ddd",
	})
}

// globFragment.matchStart()
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// globFragment.prefixMatchStart()

func TestGlobFragmentPrefixMatchStart(t *testing.T) {
	var f = &globFragment{
		literal: "src",
	}

	verify.That(t, f.prefixMatchStart("src/aaa/bbb")).Eq([]string{"aaa/bbb"})
	verify.That(t, f.prefixMatchStart("aaa/src/bbb")).Eq([]string{})

}

func TestSubdirGlobFragmentPrefixMatchStart(t *testing.T) {
	var f = &globFragment{
		subdir:  true,
		literal: "src",
	}

	verify.That(t, f.prefixMatchStart("src/aaa/bbb")).Eq([]string{"", "aaa/bbb"})
	verify.That(t, f.prefixMatchStart("aaa/src/bbb")).Eq([]string{"", "bbb"})
	verify.That(t, f.prefixMatchStart("aaa/src/bbb/src/ccc/ddd")).Eq(
		[]string{"", "bbb/src/ccc/ddd", "ccc/ddd"})
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
			a, b := splitPath(tc.input)
			verify.That(t, a).Eq(tc.fragment)
			verify.That(t, b).Eq(tc.remainder)
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
			verify.That(t, cleanFragment(tc.input)).Eq(tc.output)
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
			re, err := globFragmentToRegexp(tc.input)
			require.That(t, err).IsNil()
			require.That(t, re).IsNotNil()
			require.That(t, re.String()).Eq(tc.re)
			require.That(t, re.MatchString(tc.match)).IsTrue()

			// require.That(t, re.MatchString(tc.match),
			// 	predicate.ContextValue{Name: "re", Value: re.String()},
			// 	predicate.ContextValue{Name: "match", Value: tc.match},
			// ).IsTrue()
		})
	}
}

func TestGlobFragmentToRegexpError(t *testing.T) {
	var pattern = `*.{a,b`
	re, err := globFragmentToRegexp(pattern)
	verify.That(t, err).IsNotNil()
	verify.That(t, re).IsNil()
}

// globFragmentToRegexp()
// ---------------------------------------------------------------------------
