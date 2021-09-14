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
	return path == "" || path == "." || path == ".." || path == "~" ||
		IsPathSeparator(path[len(path)-1]) ||
		strings.HasSuffix(path, ".") && IsPathSeparator(path[len(path)-2]) ||
		strings.HasSuffix(path, "..") && IsPathSeparator(path[len(path)-3])
}

// Base returns the second part of Split().
func Base(path string) string {
	_, base := Split(path)
	return base
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

// Dir returns the first part of Split().
func Dir(path string) string {
	dir, _ := Split(path)
	return dir
}

// Ext return the filename extension of path if any.
func Ext(path string) string {
	return filepath.Ext(path)
}

// func HasPrefix(p, prefix string) bool //DEPRECATED {
// 	return filepath.HasPrefix(path)
// }

// IsAbs reports whether the path is absolute.
func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

// Join joins multiple path fragments into a single path, preserving a trailing
// separator if any, and handling any intermediate absolute path as the new root
// of the resulting path.
func Join(elem ...string) string {
	var output strings.Builder
	for _, e := range elem {
		if filepath.IsAbs(e) {
			output.Reset()
		}
		if output.Len() > 0 {
			output.WriteRune(filepath.Separator)
		}
		output.WriteString(Clean(e))
	}
	return Clean(output.String())
}

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

// ToNative converts path to its platform native representation
func ToNative(path string) string {
	return filepath.FromSlash(path)
}

// ToSlash converts path into a path using '/' as path separator, consistent
// with other functions in this package.
func ToSlash(path string) string {
	return filepath.ToSlash(path)
}

// VolumeName returns the name of the volume if specified in path. The result is
// always empty on unix platforms.
func VolumeName(path string) string {
	return ToSlash(filepath.VolumeName(path))
}

// ---------------------------------------------------------------------------
// Filename manipulation function that might need to access the underlying
// filesystem to evaluate their result.
// ---------------------------------------------------------------------------

// Abs is equivalent to `filepath.Abs()`, but preserves any trailing path
// separator.
func Abs(path string) (string, error) {
	output, err := filepath.Abs(path)
	if err == nil {
		output = Clean(output)
		if IsDirectoryName(path) && !hasTrailingSeparator(output) {
			output += string(Separator)
		}
	}
	return output, err
}

// EvalSymlinks is equivalent to `filepath.EvalSymlinks()`, but preserves any
// trailing path separator.
func EvalSymlinks(path string) (string, error) {
	output, err := filepath.EvalSymlinks(path)
	if err == nil {
		output = Clean(output)
		if IsDirectoryName(path) && !hasTrailingSeparator(output) {
			output += string(Separator)
		}
	}
	return output, err
}

// Rel is equivalent to `filepath.Rel()`, but preserves any trailing path
// separator.
func Rel(basepath, targetpath string) (string, error) {
	output, err := filepath.Rel(basepath, targetpath)
	if err == nil {
		output = Clean(output)
		if IsDirectoryName(targetpath) && !hasTrailingSeparator(output) {
			output += string(Separator)
		}
	}
	return output, err
}

// ---------------------------------------------------------------------------
// Functions from "path/filepath" that don't have a direct equivalent in this
// package.
// ---------------------------------------------------------------------------

// func Match(pattern, name string) (matched bool, err error)
// func Glob(pattern string) (matches []string, err error)
// func Walk(root string, fn WalkFunc) error
// func WalkDir(root string, fn fs.WalkDirFunc) error
