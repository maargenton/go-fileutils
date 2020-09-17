package downloader_test

import (
	"context"
	"testing"

	"github.com/maargenton/fileutil/downloader"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

var testURL = "https://developer.arm.com/-/media/Files/downloads/gnu-rm/9-2019q4/gcc-arm-none-eabi-9-2019-q4-major-mac.tar.bz2?revision=c2c4fe0e-c0b6-4162-97e6-7707e12f2b6e&amp;la=en&amp;hash=EC9D4B5F5B050267B924F876B306D72CDF3BDDC0"

// var testURL = "https://github.com/jung-kurt/gofpdf/archive/v2.17.2.tar.gz"

func TestGet(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())

	var client, err = downloader.NewClient("testdata/output")
	assert.That(err, p.IsNoError())

	var request = client.Get(context.Background(), testURL)
	err = request.Wait()
	assert.That(err, p.IsNoError())
}

// func TestGrab(t *testing.T) {
// 	assert := asserter.New(t)
// 	assert.That(nil, p.IsNil())

// 	resp, err := grab.Get("testdata/output", testURL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Download saved to", resp.Filename)
// }
