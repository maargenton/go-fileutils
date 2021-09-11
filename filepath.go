package fileutils

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Separator is a fixed default path separator and is always '/'
var Separator = '/'

// OSSeparator is the os specific path separator, usually either '/' on unix
// latforms or '\\` on windows.
var OSSeparator = os.PathSeparator

// WalkFunc is the type of the function called by Walk to visit each file or
// directory. Since all walk functions in this package are based on go 1.16
// WalkDir(), it uses `fs.DirEntry` as second argument to capture the file
// information.
type WalkFunc = fs.WalkDirFunc

// IsPathSeparator returns true for runes that are either the default path
// separator or the os native path separator.
func IsPathSeparator(c rune) bool {
	return c == Separator || c == OSSeparator
}

// IsDirectoryPath returns true is the input can be inferred to be a directory
// based on its name only, without having to access the filesystem to check.
// This include paths with trailing separator, an empty path that resolved to
// "./", or paths the end with either a "." or ".." path fragment that are
// always directories.
func IsDirectoryName(path string) bool {
	return path == "" || path == "." || path == ".." ||
		IsPathSeparator(rune(path[len(path)-1])) ||
		strings.HasSuffix(path, ".") && IsPathSeparator(rune(path[len(path)-2])) ||
		strings.HasSuffix(path, "..") && IsPathSeparator(rune(path[len(path)-3]))
}

func Base(path string) string {
	return filepath.Base(path)
}

// Clean is equivalent to `filepath.Clean()`, but preserves any trailing path
// separator or appends one for '.' or '..' path fragments.
func Clean(input string) string {
	dir := input == "" || input == "." || input == ".." ||
		strings.HasSuffix(input, string(filepath.Separator)+"..") ||
		strings.HasSuffix(input, string(filepath.Separator)+".") ||
		hasTrailingSeparator(input)

	output := filepath.Clean(input)
	if dir && !hasTrailingSeparator(output) {
		output += string(filepath.Separator)
	}
	return output
}

// func Clean(path string) string {
// 	return filepath.Clean(path)
// }

func Dir(path string) string {
	return filepath.Dir(path)
}

func Ext(path string) string {
	return filepath.Ext(path)
}

func FromSlash(path string) string {
	return filepath.FromSlash(path)
}

// func HasPrefix(p, prefix string) bool //DEPRECATED {
// 	return filepath.HasPrefix(path)
// }

func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

// func Join(elem ...string) string {
// 	return filepath.Join(elem...)
// }

func Split(path string) (dir, file string) {
	return filepath.Split(path)
}

func SplitList(path string) []string {
	return filepath.SplitList(path)
}

func ToSlash(path string) string {
	return filepath.ToSlash(path)
}

func VolumeName(path string) string {
	return filepath.VolumeName(path)
}

// func Abs(path string) (string, error)
// func EvalSymlinks(path string) (string, error)
// func Glob(pattern string) (matches []string, err error)
// func Match(pattern, name string) (matched bool, err error)
// func Rel(basepath, targpath string) (string, error)
// func Walk(root string, fn WalkFunc) error
// func WalkDir(root string, fn fs.WalkDirFunc) error
