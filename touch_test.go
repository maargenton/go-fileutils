package fileutil_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

func TestTouchCreatesTargetFile(t *testing.T) {
	assert := asserter.New(t)

	basepath, err := ioutil.TempDir(".", "testdata-")
	assert.That(err, p.IsNoError())
	defer os.RemoveAll(basepath) // clean up

	path := filepath.Join(basepath, "file.txt")
	err = fileutil.Touch(path)
	assert.That(err, p.IsNoError())
	assert.That(fileutil.Exists(path), p.IsTrue())
}
