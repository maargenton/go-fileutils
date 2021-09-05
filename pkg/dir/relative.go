package dir

import (
	"os"
	"path/filepath"
)

// MakeRelativeWalkFunc return a wrapping WalkFunc that forward the calls with a
// path made relative to the provided basepath
func MakeRelativeWalkFunc(basepath string, clientFn WalkFunc) WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if path == basepath {
			return nil
		}
		relpath, relerr := filepath.Rel(basepath, path)
		if relerr != nil {
			relpath = path
			if err == nil {
				err = relerr
			}
		}
		return clientFn(relpath, info, err)
	}
}
