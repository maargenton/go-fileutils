package dir_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maargenton/fileutil/dir"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

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
	basepath, cleanup, err := setupTestFolderWithSymlinks()
	assert.That(err, p.IsNoError())
	defer cleanup()

	var records []traversalRecord
	dir.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			relpath, _ := filepath.Rel(basepath, path)
			records = append(records, traversalRecord{
				path: relpath,
				mode: info.Mode() & os.ModeType,
				err:  err,
			})
		}
		return nil
	})

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
}

func TestWalkReportsErrorOnInvalidSymlinks(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	basepath, cleanup, err := setupTestFolderWithSymlinks()
	assert.That(err, p.IsNoError())
	defer cleanup()

	var records = make(map[string]traversalRecord)
	dir.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			relpath, _ := filepath.Rel(basepath, path)
			records[relpath] = traversalRecord{
				path: relpath,
				mode: info.Mode() & os.ModeType,
				err:  err,
			}
		}
		return nil
	})

	// for _, r := range records {
	// 	fmt.Printf("%v - %v - %v\n", r.path, decodeMode(r.mode), r.err)
	// }

	// Symlink - lstat testdata-625528602_src2: no such file or directory
	src_src3 := records["src/src3"]
	assert.That(os.IsNotExist(src_src3.err), p.IsTrue(), "error", src_src3.err)

	// Symlink - lstat testdata-625528602_src2: no such file or directory
	dst_src3 := records["dst/src3"]
	assert.That(os.IsNotExist(dst_src3.err), p.IsTrue(), "error", dst_src3.err)

	// Symlink - lstat testdata-625528602/src2: no such file or directory
	src_src_src3 := records["src/src/src3"]
	assert.That(os.IsNotExist(src_src_src3.err), p.IsTrue(), "error", src_src_src3.err)
}
