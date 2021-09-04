package dir

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// WalkFunc is a local type alias for `filepath.WalkFunc`
type WalkFunc = filepath.WalkFunc

// ErrRecursiveSymlink is a sentinel error returned during symlink traversal if
// the target path is a symlink that has already been visited.
var ErrRecursiveSymlink = errors.New("Recursive symlink detected")

// ErrAbortWalk is a convenience sentinel error that can be returned by the
// client walkFn to stop the filesystem walk immediately. Note that any error
// returned by walkFn, other than ErrSkipDir will have the same effect in
// aborting the filesystem traversal and will be returned by Walk().
var ErrAbortWalk = errors.New("Recursive symlink detected")

// ErrSkipDir is a local alias for the `filepath.SkipDir` sentinel error that
// stop the descent into a subdirectory, but continues the traversal.
var ErrSkipDir = filepath.SkipDir

// Walk is function that walks through files and folders of a filesystem and
// call the client walkFn for each item discovered. Unlike `filepath.Walk`, it
// also traverses symlinks safely and detects potential recursions.
func Walk(root string, walkFn WalkFunc) error {
	visited := make([]string, 0, 16)
	walkRoot := root
	if walkRoot == "" {
		walkRoot = "."
	}
	return filepath.Walk(walkRoot, MakeSymlinkWalkFunc(visited, root, root, walkFn))
}

// MakeSymlinkWalkFunc returns a wrapped `WalkFunc` that handles symlink
// evaluation and detects potential recursions.
func MakeSymlinkWalkFunc(visited []string, basepath, realpath string, clientFn WalkFunc) WalkFunc {
	var l = len(realpath)
	return func(path string, info os.FileInfo, err error) error {
		walkpath := filepath.Join(basepath, path[l:])
		if err != nil {
			return clientFn(walkpath, info, err)
		}
		if isSymlink(info) {
			realpath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return clientFn(walkpath, info, err)
			}
			for _, v := range visited {
				if strings.HasPrefix(realpath, v) {
					return clientFn(walkpath, info, ErrRecursiveSymlink)
				}
			}

			visited := append(visited, realpath)
			return filepath.Walk(realpath, MakeSymlinkWalkFunc(visited, walkpath, realpath, clientFn))
		}
		return clientFn(walkpath, info, err)
	}
}

func isSymlink(info os.FileInfo) bool {
	return (info.Mode() & os.ModeSymlink) != 0
}
