package fileutil_test

import (
	"io/fs"
	"testing"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/go-testpredicate/pkg/verify"
)

func TestWalkDirWithPrefix(t *testing.T) {
	var records []string
	var f = makeWalkDirPathRecorder(&records, nil)

	err := fileutil.WalkDir("testdata", "", f)
	verify.That(t, err).IsError(nil)

	verify.That(t, records).IsEqualSet([]string{
		"src/",
		"src/foo.cpp",
		"src/foo.h",
	})
}

func TestWalkDirWithRoot(t *testing.T) {
	var records []string
	var f = makeWalkDirPathRecorder(&records, nil)

	err := fileutil.WalkDir("", "testdata", f)
	verify.That(t, err).IsError(nil)
	verify.That(t, records).IsEqualSet([]string{
		"testdata/src/",
		"testdata/src/foo.cpp",
		"testdata/src/foo.h",
	})
}

func TestWalkDirWithNoRootNoPrefix(t *testing.T) {
	var records []string
	var f = makeWalkDirPathRecorder(&records, nil)

	err := fileutil.WalkDir("", "", f)
	verify.That(t, err).IsError(nil)
	verify.That(t, records).IsSupersetOf([]string{
		"testdata/",
		"testdata/src/",
		"testdata/src/foo.cpp",
		"testdata/src/foo.h",
	})
	verify.That(t, records).IsDisjointSetFrom([]string{
		"./",
	})
}

func TestWalkDirFromFsRoot(t *testing.T) {
	var records []string
	var f = makeWalkDirPathRecorder(&records, skipDirFunc)

	err := fileutil.WalkDir("", "/", f)
	verify.That(t, err).IsError(nil)
	verify.That(t, records).IsSupersetOf([]string{
		"/bin/",
		"/dev/",
		"/sbin/",
		"/usr/",
	})
	verify.That(t, records).IsDisjointSetFrom([]string{
		"/",
	})
}

// ---------------------------------------------------------------------------
// Helpers

func makeWalkDirPathRecorder(records *[]string, clientFn fs.WalkDirFunc) fs.WalkDirFunc {
	f := func(path string, d fs.DirEntry, err error) error {
		*records = append(*records, path)
		if clientFn != nil {
			return clientFn(path, d, err)
		}
		return err
	}
	return f
}

func skipDirFunc(path string, d fs.DirEntry, err error) error {
	if err == nil && d.IsDir() {
		return fs.SkipDir
	}
	return err
}
