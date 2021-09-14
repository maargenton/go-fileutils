// Package fileutils is a collection of filename manipulation and filesystem
// utilities including directory traversal with symlinks support, finding file
// and folders with extended glob pattern, and atomic file operations.
//
// To help support non-unix platforms, it also includes ad set of functions that
// are similar to those found in package "path/filepath", but but using '/' as
// path separator, and preserving trailing separator for directory filenames.
package fileutils
