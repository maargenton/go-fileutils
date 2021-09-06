package fileutil

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/maargenton/go-errors"
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

// ErrRecursiveSymlink is a sentinel error returned during symlink traversal if
// the target path is a symlink that has already been visited.
var ErrRecursiveSymlink = errors.Sentinel("ErrRecursiveSymlink")

// WalkDirSymlink is similar to `WalkDir()`, but also follows symlinks and
// detects potential recursions. When evaluating symlinks, any error evaluating
// the link is passed to the client, but obviously the link is not followed. If
// the link points to a directory, the client is called with wrapped FileInfo of
// that directory. If the link points back to a folder along the current path,
// the client is called with ErrRecursiveSymlink and the link evaluation stops
// there.
func WalkDirSymlink(prefix, root string, fn fs.WalkDirFunc) error {
	visited := make([]string, 0, 16)
	return WalkDir(prefix, root, makeSymlinkWalkFunc(visited, prefix, "", fn))
}

func makeSymlinkWalkFunc(visited []string, basepath, clientPrefix string, clientFn fs.WalkDirFunc) fs.WalkDirFunc {
	f := func(path string, d fs.DirEntry, err error) error {
		clientPath := Join(clientPrefix, path)
		if err != nil {
			return clientFn(clientPath, d, err)
		}

		if !isSymlink(d) {
			return clientFn(clientPath, d, err)
		}

		linkpath := Join(basepath, path)
		realpath, err := filepath.EvalSymlinks(linkpath)
		if err != nil {
			return clientFn(clientPath, d, err)
		}
		info, err := os.Lstat(realpath)
		if err != nil {
			return clientFn(clientPath, d, err)
		}
		d = fs.FileInfoToDirEntry(info)
		if info.IsDir() {
			clientPath = Join(clientPath, "")
		}

		// Check if visited and recurse
		for _, v := range visited {
			if strings.HasPrefix(realpath, v) {
				return clientFn(clientPath, d, ErrRecursiveSymlink)
			}
		}

		err = clientFn(clientPath, d, err)
		if err != nil {
			return err
		}

		visited := append(visited, realpath)
		return WalkDir(realpath, "",
			makeSymlinkWalkFunc(visited, realpath, path, clientFn))
	}
	return f
}

func isSymlink(d fs.DirEntry) bool {
	return (d.Type() & os.ModeSymlink) != 0
}
