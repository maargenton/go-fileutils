package fileutil

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// WalkDir is similar to filepath.WalkDir with a few additions to make it more
// convenient to use in common cases.
//
// It takes an additional parameter `prefix` so that the root of the walk is
// <prefix>/<root>, and paths are reported relative to it, e.g.
// '<prefix>/<root>/path/to/file' is reported as '<root>/path/to/file'.
//
// The root of the walk is always recursed into, but never reported unless an
// error occurs during filesystem operation.
//
// If `root` is an absolute path, `prefix` is ignored and all reported paths are
// also absolute.
//
// It reports directory paths with a tailing path separator.
func WalkDir(prefix, root string, fn fs.WalkDirFunc) error {
	walkRoot := Join(prefix, root)
	if filepath.IsAbs(root) {
		walkRoot = Clean(root)
		prefix = ""
	} else {
		if len(prefix) > 0 {
			prefix += string(filepath.Separator)
		}
		prefix = Clean(prefix)
	}

	f := func(path string, d fs.DirEntry, err error) error {
		if path == walkRoot && err == nil {
			return nil
		}
		if strings.HasPrefix(path, prefix) {
			path = path[len(prefix):]
		} else {
			relpath, relerr := Rel(prefix, path)
			if relerr != nil {
				return fmt.Errorf(
					"WalkDir yielded a path '%v' that is not relative to the root path '%v'",
					path, prefix)
			}
			path = relpath
		}
		if len(path) > 0 && d != nil && d.IsDir() {
			path += string(filepath.Separator)
		}
		path = Clean(path)
		return fn(path, d, err)
	}

	return filepath.WalkDir(walkRoot, f)
}
