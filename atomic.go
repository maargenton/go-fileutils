package fileutils

import (
	"fmt"
	"io"
	"os"
	"time"
)

// OpenTemp returns the name and handle to a newly created temporary file,
// guarantied to not previously exist, located in the same directory as the
// specified file.
func OpenTemp(filename, suffix string) (f *os.File, err error) {
	for {
		tmp := RewriteFilename(filename, &RewriteOpts{
			Suffix: fmt.Sprintf("-%v-%x", suffix, time.Now().Nanosecond()),
		})
		f, err = os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err == nil || !os.IsExist(err) {
			return
		}
	}
}

// ReadFile opens a file for reading ans passes if to the provided reader for
// loading. The file is closed when the function returns
func ReadFile(filename string, reader func(r io.Reader) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return reader(f)
}

// WriteFile creates or atomically replaces the specified file with the content
// written into w by hhe provided function.
func WriteFile(filename string, writer func(w io.Writer) error) error {

	f, err := OpenTemp(filename, "atomic")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if err := writer(f); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(f.Name(), filename)
}
