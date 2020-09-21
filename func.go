package fileutil

import "io"

// ReaderFunc is a function type that implements io.Reader
type ReaderFunc func(p []byte) (n int, err error)

var _ io.Reader = ReaderFunc(nil)

func (r ReaderFunc) Read(p []byte) (n int, err error) {
	return r(p)
}

// WriterFunc is a function type that implements io.Writer
type WriterFunc func(p []byte) (n int, err error)

var _ io.Writer = WriterFunc(nil)

func (r WriterFunc) Write(p []byte) (n int, err error) {
	return r(p)
}
