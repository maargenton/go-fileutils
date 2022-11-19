package dir_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/subexpr"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-fileutils/pkg/dir"
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
