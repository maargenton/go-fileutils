package fileutils_test

// import (
// 	"os"
// 	"testing"

// 	"github.com/maargenton/go-testpredicate/pkg/verify"

// 	"github.com/maargenton/go-fileutils"
// )

// func TestExistsFunctionsWithFile(t *testing.T) {
// 	filename := "exist_test.go"
// 	verify.That(t, fileutils.Exists(filename)).IsTrue()
// 	verify.That(t, fileutils.IsFile(filename)).IsTrue()
// 	verify.That(t, fileutils.IsDir(filename)).IsFalse()
// }

// func TestExistsFunctionsWithDir(t *testing.T) {
// 	path, _ := os.Getwd()
// 	verify.That(t, fileutils.Exists(path)).IsTrue()
// 	verify.That(t, fileutils.IsFile(path)).IsFalse()
// 	verify.That(t, fileutils.IsDir(path)).IsTrue()
// }
