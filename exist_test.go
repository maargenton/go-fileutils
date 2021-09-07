package fileutils_test

import (
	"os"
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

func TestExistsFunctionsWithFile(t *testing.T) {
	assert := asserter.New(t)

	filename := "exist_test.go"
	assert.That(fileutils.Exists(filename), p.IsTrue())
	assert.That(fileutils.IsFile(filename), p.IsTrue())
	assert.That(fileutils.IsDir(filename), p.IsFalse())
}

func TestExistsFunctionsWithDir(t *testing.T) {
	assert := asserter.New(t)

	path, _ := os.Getwd()
	assert.That(fileutils.Exists(path), p.IsTrue())
	assert.That(fileutils.IsFile(path), p.IsFalse())
	assert.That(fileutils.IsDir(path), p.IsTrue())
}
