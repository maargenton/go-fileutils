package dir_test

import (
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/maargenton/go-fileutils"
)

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
