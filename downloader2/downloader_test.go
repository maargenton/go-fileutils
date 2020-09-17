package downloader_test

import (
	"context"
	"crypto/sha256"
	"io"
	"testing"

	downloader "github.com/maargenton/fileutil/downloader2"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

var testURL1 = "https://developer.arm.com/-/media/Files/downloads/gnu-rm/9-2019q4/gcc-arm-none-eabi-9-2019-q4-major-mac.tar.bz2?revision=c2c4fe0e-c0b6-4162-97e6-7707e12f2b6e&amp;la=en&amp;hash=EC9D4B5F5B050267B924F876B306D72CDF3BDDC0"
var testURL2 = "https://github.com/jung-kurt/gofpdf/archive/v2.17.2.tar.gz"

func Test(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	outputPath := ""
	client := downloader.DefaultClient
	err := client.Get(context.Background(), &downloader.Request{
		URL:             testURL2,
		OutputDirectory: "testdata/output",
		Hash:            sha256.New(),
		Checksum:        []byte{1, 2, 3},

		ContentReader: func(r io.Reader) error {
			return extractArchiveContent(r, outputPath)
		},
	})

	assert.That(err, p.IsNoError())
}

func extractArchiveContent(r io.Reader, path string) error {
	return nil
}

// func TestFilepathBase(t *testing.T) {
// 	assert := asserter.New(t)
// 	assert.That(nil, p.IsNil())

// 	filename := "./aaa/bbb/archive.tar.bz2"

// 	assert.That(filepath.Base(filename), p.Eq(filename))
// 	assert.That(filepath.Dir(filename), p.Eq(""))

// }
