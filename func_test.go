package fileutil_test

import (
	"io"
	"strings"
	"testing"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

func TestReaderFunc(t *testing.T) {
	assert := asserter.New(t)

	var s = "Hello wonderful world of reader / writer func!"
	var rr = strings.NewReader(s)
	var ww = &strings.Builder{}
	r := fileutil.ReaderFunc(func(p []byte) (n int, err error) {
		return rr.Read(p)
	})

	w := fileutil.WriterFunc(func(p []byte) (n int, err error) {
		return ww.Write(p)
	})
	n, err := io.Copy(w, r)
	assert.That(err, p.IsNoError())
	assert.That(n, p.Eq(len(s)))
	assert.That(ww.String(), p.Eq(s))
}
