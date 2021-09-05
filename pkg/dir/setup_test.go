package dir_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/maargenton/fileutil"
)

func setupTestFolder() (basepath string, cleanup func(), err error) {
	basepath, err = ioutil.TempDir(".", "testdata-")
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
			filepath.Join(basepath, "src", n, n+".h"),
			filepath.Join(basepath, "src", n, n+".cpp"),
			filepath.Join(basepath, "src", n, n+"_test.cpp"),
		)
	}
	err = fileutil.Touch(filenames...)
	return
}

func setupTestFolderWithSymlinks() (basepath string, cleanup func(), err error) {
	basepath, cleanup, err = setupTestFolder()
	if err == nil {
		os.Symlink("src", filepath.Join(basepath, "dst"))          // Regular symlink
		os.Symlink("../src", filepath.Join(basepath, "dst/src"))   // Recursive symlink
		os.Symlink("../src2", filepath.Join(basepath, "dst/src3")) // Invalid destination symlink
	}
	return
}
