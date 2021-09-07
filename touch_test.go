package fileutils_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

func TestTouchCreatesTargetFile(t *testing.T) {
	assert := asserter.New(t)

	basepath, err := ioutil.TempDir(".", "testdata-")
	assert.That(err, p.IsNoError())
	defer os.RemoveAll(basepath) // clean up

	path := filepath.Join(basepath, "file.txt")
	err = fileutils.Touch(path)
	assert.That(err, p.IsNoError())
	assert.That(fileutils.Exists(path), p.IsTrue())
}
