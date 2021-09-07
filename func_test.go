package fileutils_test

import (
	"io"
	"strings"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-fileutils"
)

func TestReaderFunc(t *testing.T) {

	var s = "Hello wonderful world of reader / writer func!"
	var rr = strings.NewReader(s)
	var ww = &strings.Builder{}
	r := fileutils.ReaderFunc(func(p []byte) (n int, err error) {
		return rr.Read(p)
	})

	w := fileutils.WriterFunc(func(p []byte) (n int, err error) {
		return ww.Write(p)
	})
	n, err := io.Copy(w, r)
	verify.That(t, err).IsNil()
	verify.That(t, n).Eq(len(s))
	verify.That(t, ww.String()).Eq(s)
}
