package dir_test

import (
	"io/fs"
	"os"
	"runtime"
	"testing"

	"github.com/maargenton/go-fileutils/pkg/dir"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/subexpr"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-fileutils"
)

func TestWalkWithPrefix(t *testing.T) {
	var records []string
	var f = makeWalkDirPathRecorder(&records, nil)

	err := dir.Walk("testdata", "", f)
	verify.That(t, err).IsError(nil)

	verify.That(t, records).IsEqualSet([]string{
		"src/",
		"src/foo.cpp",
		"src/foo.h",
		"dst/",
		"dst/foo.cpp",
		"dst/foo.h",
	})
}

func TestWalkWithRoot(t *testing.T) {
	var records []string
	var f = makeWalkDirPathRecorder(&records, nil)

	err := dir.Walk("", "testdata", f)
	verify.That(t, err).IsError(nil)
	verify.That(t, records).IsEqualSet([]string{
		"testdata/src/",
		"testdata/src/foo.cpp",
		"testdata/src/foo.h",
		"testdata/dst/",
		"testdata/dst/foo.cpp",
		"testdata/dst/foo.h",
	})
}

func TestWalkWithNoRootNoPrefix(t *testing.T) {
	var records []string
	var f = makeWalkDirPathRecorder(&records, nil)

	err := dir.Walk("", "", f)
	verify.That(t, err).IsError(nil)
	verify.That(t, records).IsSupersetOf([]string{
		"testdata/",
		"testdata/src/",
		"testdata/src/foo.cpp",
		"testdata/src/foo.h",
		"testdata/dst/",
		"testdata/dst/foo.cpp",
		"testdata/dst/foo.h",
	})
	verify.That(t, records).IsDisjointSetFrom([]string{
		"./",
	})
}

func TestWalkFromFsRoot(t *testing.T) {
	var records []string
	var f = makeWalkDirPathRecorder(&records, skipDirFunc)

	if runtime.GOOS == "windows" {
		err := dir.Walk("", "C:/", f)
		verify.That(t, err).IsError(nil)
		verify.That(t, records).IsSupersetOf([]string{
			"C:/Documents and Settings/",
			"C:/Program Files/",
		})
		verify.That(t, records).IsDisjointSetFrom([]string{
			"C:/",
		})

	} else {
		err := dir.Walk("", "/", f)
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
}

func TestSymlinks(t *testing.T) {
	recursive, broken := false, false
	basepath, cleanup, err := setupTestFolderWithSymlinks(recursive, broken)
	require.That(t, err).IsNil()
	defer cleanup()

	var records []string
	var f = makeWalkDirPathRecorder(&records, nil)

	err = dir.Walk(basepath, "", f)
	verify.That(t, err).IsNil()
	verify.That(t, records).Any(subexpr.Value().StartsWith("src/foo/"))
	verify.That(t, records).Any(subexpr.Value().StartsWith("src/bar/"))
	verify.That(t, records).Any(subexpr.Value().StartsWith("dst/"))
	verify.That(t, records).Any(subexpr.Value().StartsWith("dst/foo/"))
	verify.That(t, records).Any(subexpr.Value().StartsWith("dst/bar/"))
}

func TestSymlinksRecursion(t *testing.T) {
	recursive, broken := true, false
	basepath, cleanup, err := setupTestFolderWithSymlinks(recursive, broken)
	require.That(t, err).IsNil()
	defer cleanup()

	var records []walkErrorRecord
	var f = makeWalkDirErrorRecorder(&records, nil)

	err = dir.Walk(basepath, "", f)
	verify.That(t, err).IsNil()
	verify.That(t, records).Field("Err").All(
		subexpr.Value().IsError(dir.ErrRecursiveSymlink),
	)
	verify.That(t, records).Field("Path").IsEqualSet([]string{
		"dst/src/",
		"src/src/src/",
	})
}

func TestSymlinksBroken(t *testing.T) {
	recursive, broken := false, true
	basepath, cleanup, err := setupTestFolderWithSymlinks(recursive, broken)
	require.That(t, err).IsNil()
	defer cleanup()

	var records []walkErrorRecord
	var f = makeWalkDirErrorRecorder(&records, nil)

	err = dir.Walk(basepath, "", f)
	verify.That(t, err).IsNil()
	verify.That(t, records).Field("Path").IsEqualSet([]string{
		"src/src3",
		"dst/src3",
	})
	verify.That(t, os.IsNotExist(records[0].Err)).IsTrue()
	verify.That(t, os.IsNotExist(records[1].Err)).IsTrue()
}

// ---------------------------------------------------------------------------
// Helpers

// func setupTestFolder() (basepath string, cleanup func(), err error) {
// 	basepath, err = ioutil.TempDir(".", "testdata-")
// 	cleanup = func() {
// 		if basepath != "" {
// 			os.RemoveAll(basepath)
// 		}
// 	}
// 	if err != nil {
// 		return
// 	}

// 	var filenames []string
// 	for _, n := range []string{"foo", "bar"} {
// 		filenames = append(filenames,
// 			fileutils.Join(basepath, "src", n, n+".h"),
// 			fileutils.Join(basepath, "src", n, n+".cpp"),
// 			fileutils.Join(basepath, "src", n, n+"_test.cpp"),
// 		)
// 	}
// 	err = fileutils.Touch(filenames...)
// 	return
// }

func setupTestFolderWithSymlinks(recursive, broken bool) (basepath string, cleanup func(), err error) {
	basepath, cleanup, err = setupTestFolder()
	if err == nil {
		os.Symlink("src", fileutils.Join(basepath, "dst")) // Regular symlink
		if recursive {
			os.Symlink("../src", fileutils.Join(basepath, "dst/src")) // Recursive symlink
		}
		if broken {
			os.Symlink("../src2", fileutils.Join(basepath, "dst/src3")) // Invalid destination symlink
		}
	}
	return
}

func makeWalkDirPathRecorder(records *[]string, clientFn fs.WalkDirFunc) fs.WalkDirFunc {
	f := func(path string, d fs.DirEntry, err error) error {
		*records = append(*records, path)
		if clientFn != nil {
			return clientFn(path, d, err)
		}
		return nil
	}
	return f
}

type walkErrorRecord struct {
	Path string
	Err  error
}

func makeWalkDirErrorRecorder(records *[]walkErrorRecord, clientFn fs.WalkDirFunc) fs.WalkDirFunc {
	f := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			*records = append(*records, walkErrorRecord{path, err})
		}
		if clientFn != nil {
			return clientFn(path, d, err)
		}
		return nil
	}
	return f
}

func skipDirFunc(path string, d fs.DirEntry, err error) error {
	if err == nil && d.IsDir() {
		return fs.SkipDir
	}
	return nil
}
