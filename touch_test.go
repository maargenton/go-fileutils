package fileutils_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-fileutils"
)

func TestTouchCreatesTargetFile(t *testing.T) {
	basepath, err := ioutil.TempDir(".", "testdata-")
	require.That(t, err).IsNil()
	defer os.RemoveAll(basepath) // clean up

	path := filepath.Join(basepath, "file.txt")
	err = fileutils.Touch(path)
	verify.That(t, err).IsNil()
	verify.That(t, fileutils.Exists(path)).IsTrue()
}
