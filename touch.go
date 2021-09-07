package fileutils

import (
	"os"
	"path/filepath"
)

// Touch create a new files or update files mtime at the specified locations. If
// the destination directory does not exist, it is created as well. An errors is
// returned if any destination directory does not exist and cannot be created,
// or any of the specified files cannot be opened.
func Touch(filenames ...string) (err error) {
	for _, filename := range filenames {
		dirname := filepath.Dir(filename)
		err = os.MkdirAll(dirname, 0777)
		if err == nil {
			var f *os.File
			f, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0666)
			if err == nil {
				err = f.Close()
			}
		}
		if err != nil {
			break
		}
	}
	return
}
