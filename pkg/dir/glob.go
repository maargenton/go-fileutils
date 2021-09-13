package dir

import (
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"strings"

	"github.com/maargenton/go-fileutils"
	// "golang.org/x/tools/internal/fastwalk"
)

// Glob scans the file tree and returns a list of filenames matching the
// pattern. The pattern must be specified according to the extended glob pattern
// described in the package level documentation.
func Glob(pattern string) (matches []string, err error) {
	m, err := NewGlobMatcher(pattern)
	if err != nil {
		return
	}
	return m.Glob()
}

// GlobFrom scans the file tree starting at basepath and returns a list of
// filenames matching the pattern. The resulting filenames contain the full path
// including the basepath prefix. The pattern should be relative and must be
// specified according to the extended glob pattern described in the package
// level documentation. If the pattern is absolute, the basepath is ignored and
// will not appear as a prefix in the matches.
func GlobFrom(basepath, pattern string) (matches []string, err error) {
	m, err := NewGlobMatcher(pattern)
	if err != nil {
		return
	}
	return m.GlobFrom(basepath)
}

// Scan scans the file tree for filenames matching the pattern and call the
// walkFn function for every match. The pattern must be specified according to
// the extended glob pattern described in the package level documentation.
func Scan(pattern string, walkFn fs.WalkDirFunc) error {
	m, err := NewGlobMatcher(pattern)
	if err != nil {
		return err
	}
	return m.Scan(walkFn)
}

// ScanFrom scans the file tree starting at basepath for filenames matching the
// pattern and call the walkFn function for every match. The resulting filenames
// contain the full path including the basepath prefix. The pattern should be
// relative and must be specified according to the extended glob pattern
// described in the package level documentation. If the pattern is absolute, the
// basepath is ignored and will not appear as a prefix in the matches.
func ScanFrom(basepath, pattern string, walkFn fs.WalkDirFunc) error {
	m, err := NewGlobMatcher(pattern)
	if err != nil {
		return err
	}
	return m.ScanFrom(basepath, walkFn)
}

// ---------------------------------------------------------------------------
// GlobMatcher

// GlobMatcher is a pre-compiled matcher for a glob pattern
type GlobMatcher struct {
	pattern   string
	prefix    string
	fragments []globFragment
}

// NewGlobMatcher compiles an extended glob pattern into a GlobMatcher
func NewGlobMatcher(pattern string) (m *GlobMatcher, err error) {
	var fragments = pattern
	var subdir = false
	var prefix = true

	m = &GlobMatcher{pattern: pattern}
	for fragments != "" {
		var fragment string
		fragment, fragments = splitPath(fragments)
		if isSubdirectoryGlob(fragment) {
			prefix = false
			subdir = true
		} else if isGlobFragment(fragment) {
			prefix = false
			fragment = cleanFragment(fragment)
			re, err := globFragmentToRegexp(fragment)
			if err != nil {
				return nil, err
			}
			m.fragments = append(m.fragments, globFragment{
				subdir: subdir,
				re:     re,
			})
			subdir = false
		} else {
			if prefix {
				m.prefix = fileutils.Join(m.prefix, fragment)
			} else {
				m.fragments = append(m.fragments, globFragment{
					subdir:  subdir,
					literal: cleanFragment(fragment),
				})
				subdir = false
			}
		}
	}

	if m.prefix != "" && !strings.HasSuffix(m.prefix, string(fileutils.Separator)) {
		m.prefix += string(fileutils.Separator)
	}

	return
}

// Match returns true if the provided filename matches the compiled glob
// expressions
func (m *GlobMatcher) Match(filename string) bool {
	if !strings.HasPrefix(filename, m.prefix) {
		return false
	}

	filename = filename[len(m.prefix):]
	return matchFragments(filename, m.fragments)
}

func matchFragments(r string, fn []globFragment) bool {
	if len(fn) == 0 {
		return r == ""
	}

	ff, fn := fn[0], fn[1:]
	var remainders = ff.matchStart(r)
	if len(remainders) == 0 {
		return false
	}

	for _, r := range remainders {
		if matchFragments(r, fn) {
			return true
		}
	}
	return false
}

// PrefixMatch returns true if the provided filename, most likely a directory
// name, is a prefix partial match for the compiled glob expressions. This
// function can be used during scanning to skip over directories that cannot
// math the full pattern.
func (m *GlobMatcher) PrefixMatch(filename string) bool {
	if filename == cleanFragment(m.prefix) {
		return true
	}
	if !strings.HasPrefix(filename, m.prefix) {
		return false
	}

	filename = filename[len(m.prefix):]
	return prefixMatchFragments(filename, m.fragments)
}

func prefixMatchFragments(r string, fn []globFragment) bool {
	if r == "" {
		return true
	}
	if len(fn) != 0 {
		ff, fn := fn[0], fn[1:]
		var remainders = ff.prefixMatchStart(r)
		for _, r := range remainders {
			if prefixMatchFragments(r, fn) {
				return true
			}
		}
	}
	return false
}

// Glob scans the file tree and returns a list of filenames matching the
// pattern. The pattern must be specified according to the extended glob pattern
// described in the package level documentation.
func (m *GlobMatcher) Glob() (matches []string, err error) {
	return m.GlobFrom("")
}

// GlobFrom scans the file tree and returns a list of filenames matching the
// pattern. The pattern must be specified according to the extended glob pattern
// described in the package level documentation.
func (m *GlobMatcher) GlobFrom(basepath string) (matches []string, err error) {
	err = m.ScanFrom(basepath, func(path string, d fs.DirEntry, err error) error {
		if err == nil {
			matches = append(matches, path)
		}
		return nil
	})
	return
}

// Scan scans the file tree for filenames matching the pattern and call the
// walkFn function for every match. The pattern must be specified according to
// the extended glob pattern described in the package level documentation.
func (m *GlobMatcher) Scan(walkFn fs.WalkDirFunc) error {
	return m.ScanFrom("", walkFn)
}

// ScanFrom scans the file tree for filenames matching the pattern and call the
// walkFn function for every match. The pattern must be specified according to
// the extended glob pattern described in the package level documentation.
func (m *GlobMatcher) ScanFrom(basepath string, walkFn fs.WalkDirFunc) error {
	f := func(path string, d fs.DirEntry, err error) error {
		if d != nil && d.IsDir() && !m.PrefixMatch(path) {
			return fileutils.SkipDir
		}
		if m.Match(path) {
			err = walkFn(path, d, err)

			// If path is a match for the full pattern and a directory, there is
			// no need t ogo further in.
			if d != nil && d.IsDir() && err == nil {
				err = fileutils.SkipDir
			}
			return err
		}
		return nil // Ignore any error if no match
	}

	return fileutils.Walk(basepath, m.prefix, f)
}

// GlobMatcher
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// globFragment

type globFragment struct {
	subdir  bool
	literal string
	re      *regexp.Regexp
}

func (f *globFragment) match(fragment string) bool {
	fragment = fileutils.Clean(fragment)
	if hasTrailingSeparator(fragment) {
		fragment = fragment[:len(fragment)-1]
	}
	if f.re != nil {
		return f.re.MatchString(fragment)
	}
	return fragment == f.literal
}

func hasTrailingSeparator(path string) bool {
	return len(path) > 0 && fileutils.IsPathSeparator(path[len(path)-1])
}

// matchStart matches f against the start of a full path and retusn all possible
// reminders of path that follow. A final match return a single empty remainder.
// No remainders are returned if there is no match.
func (f *globFragment) matchStart(path string) (remainders []string) {
	if f.subdir {
		var fragments = path
		for fragments != "" {
			var fragment string
			fragment, fragments = splitPath(fragments)
			if f.match(fragment) {
				remainders = append(remainders, fragments)
			}
		}
	} else {
		fragment, fragments := splitPath(path)
		if f.match(fragment) {
			remainders = append(remainders, fragments)
		}
	}

	return
}

// prefixMatchStart matches f against the start of a partial path that could be
// the prefix of a matching path. If f matches subdirectories, the entire path
// could be matched by the wildcard prefix, resulting in an empty remainder. If
// any fragment matches the non-wildcard part of f, a remainder is added with
// everything that comes after.
func (f *globFragment) prefixMatchStart(path string) (remainders []string) {
	if !f.subdir {
		return f.matchStart(path)
	}

	var fragments = path
	remainders = append(remainders, "")
	for fragments != "" {
		var fragment string
		fragment, fragments = splitPath(fragments)
		if f.match(fragment) && fragments != "" {
			remainders = append(remainders, fragments)
		}
	}

	return
}

// globFragment
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Helper functions for parsing and matching

func splitPath(filepath string) (first, remainder string) {
	for i, c := range filepath {
		if c == '/' {
			return filepath[:i+1], filepath[i+1:]
		}
	}
	return filepath, ""
}

func cleanFragment(s string) string {
	var l = len(s)
	if l > 1 && s[l-1] == '/' {
		return s[:l-1]
	}
	return s
}

func isSubdirectoryGlob(fragment string) bool {
	return path.Clean(fragment) == "**"
}

func isGlobFragment(fragment string) bool {
	for _, c := range fragment {
		switch c {
		case '?', '*', '{', '[':
			return true
		}
	}
	return false
}

func globFragmentToRegexp(glob string) (re *regexp.Regexp, err error) {
	var s strings.Builder
	var escape bool
	var alt int
	s.WriteRune('^')
	for _, c := range glob {
		if escape {
			escape = false
			switch c {
			case '.':
				s.WriteString("\\.")
			case '{':
				s.WriteString("\\{")
			default:
				s.WriteRune(c)
			}
		} else {
			switch c {
			case '{':
				alt++
				s.WriteString("(?:(?:")
			case '}':
				if alt > 0 {
					alt--
					s.WriteString("))")
				} else {
					s.WriteRune(c)
				}
			case ',':
				if alt > 0 {
					s.WriteString(")|(?:")
				} else {
					s.WriteRune(c)
				}

			case '\\':
				escape = true
			case '*':
				s.WriteString(".*")
			case '?':
				s.WriteString(".")
			case '.':
				s.WriteString("\\.")
			default:
				s.WriteRune(c)
			}
		}
	}

	s.WriteRune('$')
	re, err = regexp.Compile(s.String())
	if err != nil {
		err = fmt.Errorf("error tranlating glob pattern '%v': %w",
			glob, err)
	}
	return
}

// Helper functions for parsing and matching
// ---------------------------------------------------------------------------
