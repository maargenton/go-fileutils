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
func IsPathSeparator(c uint8) bool {
	return rune(c) == Separator || rune(c) == OSSeparator
}

// IsDirectoryName returns true is the input can be inferred to be a directory
// based on its name only, without having to access the filesystem to check.
// This include paths with trailing separator, an empty path that resolved to
// "./", or paths the end with either a "." or ".." path fragment that are
// always directories.
func IsDirectoryName(path string) bool {
	return path == "" || path == "." || path == ".." ||
		IsPathSeparator(path[len(path)-1]) ||
		strings.HasSuffix(path, ".") && IsPathSeparator(path[len(path)-2]) ||
		strings.HasSuffix(path, "..") && IsPathSeparator(path[len(path)-3])
}

func Base(path string) string {
	return Clean(filepath.Base(path))
}

// Clean returns a lexically equivalent path, using '/' as separator, removing
// any discardable '/' or "./", and collapsing any intermediate "../". It
// preserves a trailing separator for directory names.
func Clean(input string) string {
	dir := IsDirectoryName(input)
	output := filepath.ToSlash(filepath.Clean(input))
	if dir && !hasTrailingSeparator(output) {
		output += string(Separator)
	}
	return output
}

func hasTrailingSeparator(path string) bool {
	return len(path) > 0 && IsPathSeparator(path[len(path)-1])
}

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

// Split splits the last path fragment of `path` from everything that precedes,
// so that path = dir+file.
func Split(path string) (dir, base string) {
	isDir := IsDirectoryName(path)
	for hasTrailingSeparator(path) && len(path) > 1 {
		path = path[:len(path)-1]
	}
	dir, base = filepath.Split(path)
	dir = filepath.ToSlash(dir)
	base = filepath.ToSlash(base)
	if isDir && len(base) > 0 && !hasTrailingSeparator(base) {
		base += string(Separator)
	}
	return
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