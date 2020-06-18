package dir_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maargenton/fileutil/dir"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

// func TestGlobMatcherGlobWithSymlinks(t *testing.T) {
// 	var tcs = []struct {
// 		pattern string
// 		count   int
// 	}{
// 		// {`**/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
// 		{`dst/{foo,bar}/**/*_test.{c,cc,cpp}`, 2},
// 		// {`dst/**/*_test.{c,cc,cpp}`, 4},
// 		// {`dst/**/*_test.{h,hh,hpp}`, 0},
// 		// {`dst/**/*.{h,cpp}`, 12},
// 	}

// 	assert := asserter.New(t, asserter.AbortOnError())
// 	dirname, err := ioutil.TempDir(".", "testdata-")
// 	assert.That(err, p.IsNoError())
// 	defer os.RemoveAll(dirname) // clean up
// 	err = setupTestFs(dirname)
// 	assert.That(err, p.IsNoError())

// 	os.Symlink("src", filepath.Join(dirname, "dst"))
// 	os.Symlink("../src", filepath.Join(dirname, "dst/src"))

// 	for _, tc := range tcs {
// 		t.Run(tc.pattern, func(t *testing.T) {
// 			assert := asserter.New(t, asserter.AbortOnError())
// 			assert.That(true, p.IsTrue())

// 			pattern := path.Join(dirname, tc.pattern)
// 			m, err := dir.NewGlobMatcher(pattern)
// 			assert.That(err, p.IsNoError())

// 			matches, err := m.Glob()
// 			assert.That(err, p.IsNoError())
// 			assert.That(matches, p.Length(p.Eq(tc.count)))
// 		})
// 	}
// }

func decodeMode(mode os.FileMode) string {

	var modes []string
	if (mode & os.ModeDir) != 0 {
		modes = append(modes, "Dir")
	}
	if (mode & os.ModeAppend) != 0 {
		modes = append(modes, "Append")
	}
	if (mode & os.ModeExclusive) != 0 {
		modes = append(modes, "Exclusive")
	}
	if (mode & os.ModeTemporary) != 0 {
		modes = append(modes, "Temporary")
	}
	if (mode & os.ModeSymlink) != 0 {
		modes = append(modes, "Symlink")
	}
	if (mode & os.ModeDevice) != 0 {
		modes = append(modes, "Device")
	}
	if (mode & os.ModeNamedPipe) != 0 {
		modes = append(modes, "NamedPipe")
	}
	if (mode & os.ModeSocket) != 0 {
		modes = append(modes, "Socket")
	}
	if (mode & os.ModeSetuid) != 0 {
		modes = append(modes, "Setuid")
	}
	if (mode & os.ModeSetgid) != 0 {
		modes = append(modes, "Setgid")
	}
	if (mode & os.ModeCharDevice) != 0 {
		modes = append(modes, "CharDevice")
	}
	if (mode & os.ModeSticky) != 0 {
		modes = append(modes, "Sticky")
	}
	if (mode & os.ModeIrregular) != 0 {
		modes = append(modes, "Irregular")
	}

	return strings.Join(modes, " | ")
}

type traversalRecord struct {
	path string
	mode os.FileMode
	err  error
}

func TestWalkDetectsSymlinksRecursion(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	dirname, err := ioutil.TempDir(".", "testdata-")
	assert.That(err, p.IsNoError())
	defer os.RemoveAll(dirname) // clean up
	err = setupTestFs(dirname)
	assert.That(err, p.IsNoError())

	os.Symlink("src", filepath.Join(dirname, "dst"))          // Regular symlink
	os.Symlink("../src", filepath.Join(dirname, "dst/src"))   // Recursive symlink
	os.Symlink("../src2", filepath.Join(dirname, "dst/src3")) // Invalid destination symlink

	var records []traversalRecord
	dir.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		// fmt.Printf("%v - %v - %v\n", path, decodeMode(info.Mode()), err)
		if err != nil {
			relpath, _ := filepath.Rel(dirname, path)
			records = append(records, traversalRecord{
				path: relpath,
				mode: info.Mode() & os.ModeType,
				err:  err,
			})
		}
		return nil
	})

	// for _, r := range records {
	// 	fmt.Printf("%v - %v - %v\n", r.path, decodeMode(r.mode), r.err)
	// }

	assert.That(records, p.Contains([]traversalRecord{{
		path: "dst/src",
		mode: os.ModeSymlink,
		err:  dir.ErrRecursiveSymlink}}),
	)

	assert.That(records, p.Contains([]traversalRecord{{
		path: "src/src/src",
		mode: os.ModeSymlink,
		err:  dir.ErrRecursiveSymlink}}),
	)

	// t.Fail()
}

func TestWalkReportsErrorOnInvalidSymlinks(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	dirname, err := ioutil.TempDir(".", "testdata-")
	assert.That(err, p.IsNoError())
	defer os.RemoveAll(dirname) // clean up
	err = setupTestFs(dirname)
	assert.That(err, p.IsNoError())

	os.Symlink("src", filepath.Join(dirname, "dst"))          // Regular symlink
	os.Symlink("../src", filepath.Join(dirname, "dst/src"))   // Recursive symlink
	os.Symlink("../src2", filepath.Join(dirname, "dst/src3")) // Invalid destination symlink

	var records = make(map[string]traversalRecord)
	dir.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			relpath, _ := filepath.Rel(dirname, path)
			records[relpath] = traversalRecord{
				path: relpath,
				mode: info.Mode() & os.ModeType,
				err:  err,
			}
		}
		return nil
	})

	for _, r := range records {
		fmt.Printf("%v - %v - %v\n", r.path, decodeMode(r.mode), r.err)
	}

	src_src3 := records["src/src3"] // Symlink - lstat testdata-625528602_src2: no such file or directory
	assert.That(os.IsNotExist(src_src3.err), p.IsTrue(), "error", src_src3.err)

	dst_src3 := records["dst/src3"] // Symlink - lstat testdata-625528602_src2: no such file or directory
	assert.That(os.IsNotExist(dst_src3.err), p.IsTrue(), "error", dst_src3.err)

	src_src_src3 := records["src/src/src3"] // Symlink - lstat testdata-625528602/src2: no such file or directory
	assert.That(os.IsNotExist(src_src_src3.err), p.IsTrue(), "error", src_src_src3.err)

	// assert.That(src_src_src3.err, p.Eval(os.IsNotExist, p.IsTrue()))

	// assert.That(src_src_src3.err).Eval(os.IsNotExist).IsTrue()
	// assert.That(src_src_src3).Field("err").Eval(os.IsNotExist).IsTrue()
	// assert.That(src_src_src3).Field("err").Call("funcname").Eq(123)

}

// func SymlinkWalkFn(visited []string, basepath, realpath string, clientFn filepath.WalkFunc) filepath.WalkFunc {
// 	var l = len(realpath)
// 	return func(path string, info os.FileInfo, err error) error {
// 		walkpath := filepath.Join(basepath, path[l:])
// 		if err != nil {
// 			return clientFn(walkpath, info, err)
// 		}
// 		if isSymlink(info) {
// 			realpath, err := filepath.EvalSymlinks(path)
// 			if err != nil {
// 				return clientFn(walkpath, info, err)
// 			}
// 			for _, v := range visited {
// 				if strings.HasPrefix(realpath, v) {
// 					return clientFn(walkpath, info, ErrRecursiveSymlink)
// 				}
// 			}

// 			visited := append(visited, realpath)
// 			return filepath.Walk(realpath, SymlinkWalkFn(visited, walkpath, realpath, clientFn))
// 		}
// 		return clientFn(walkpath, info, err)
// 	}
// }

// func isSymlink(info os.FileInfo) bool {
// 	return (info.Mode() & os.ModeSymlink) != 0
// }

// var ErrRecursiveSymlink = errors.New("Recursive symlink detected")
