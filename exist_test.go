package fileutil_test

import (
	"os"
	"testing"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

func TestExistsFunctionsWithFile(t *testing.T) {
	assert := asserter.New(t)

	filename := "exist_test.go"
	assert.That(fileutil.Exists(filename), p.IsTrue())
	assert.That(fileutil.IsFile(filename), p.IsTrue())
	assert.That(fileutil.IsDir(filename), p.IsFalse())
}

func TestExistsFunctionsWithDir(t *testing.T) {
	assert := asserter.New(t)

	path, _ := os.Getwd()
	assert.That(fileutil.Exists(path), p.IsTrue())
	assert.That(fileutil.IsFile(path), p.IsFalse())
	assert.That(fileutil.IsDir(path), p.IsTrue())
}
