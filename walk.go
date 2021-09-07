package fileutils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/maargenton/go-errors"
)

// ErrRecursiveSymlink is a sentinel error returned during symlink traversal if
// the target path is a symlink that has already been visited.
var ErrRecursiveSymlink = errors.Sentinel("ErrRecursiveSymlink")

// WalkDir is similar to `filepath.WalkDir()`, but it follows symlinks and takes
// an additional `prefix` argument. The walk starts at '<prefix>/<root>' and the
// paths are reported relative to `prefix`; the starting path is always recursed
// into, but never reported unless an error occurs during filesystem operation.
// Either or both of `prefix` and `root` can be empty, a relative path or an
// absolute path; if `root` is an absolute path, `prefix` is ignored, and all
// reported paths are absolute. In a path refers to a directory, it is reported
// with a trailing path separator.
//
// When symlinks are encountered, the path is reported with the FileInfo of the
// destination, and if it is a directory, it is traversed unless an error is
// returned. Any error that occurs while evaluating the symlink is reported to
// the client. If the symlink points back to a folder along the current path,
// the client is called with ErrRecursiveSymlink and the link evaluation stops
// there.
func WalkDir(prefix, root string, fn fs.WalkDirFunc) error {
	visited := make([]string, 0, 16)
	return walkDir(prefix, root, makeSymlinkWalkFunc(visited, prefix, "", fn))
}

func walkDir(prefix, root string, fn fs.WalkDirFunc) error {
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
				err = clientFn(clientPath, d, ErrRecursiveSymlink)
				if errors.Is(err, fs.SkipDir) {
					// The caller does not know the symlink points to a
					// directory and would skip the rest of the parent directory
					return nil
				}
				return err
			}
		}

		err = clientFn(clientPath, d, err)
		if err != nil {
			if errors.Is(err, fs.SkipDir) {
				// The caller does not know the symlink points to a directory
				// and would skip the rest of the parent directory
				return nil
			}
			return err
		}

		visited := append(visited, realpath)
		return walkDir(realpath, "",
			makeSymlinkWalkFunc(visited, realpath, path, clientFn))
	}
	return f
}

func isSymlink(d fs.DirEntry) bool {
	return (d.Type() & os.ModeSymlink) != 0
}
