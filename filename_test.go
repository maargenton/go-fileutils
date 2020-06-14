package fileutil_test

import (
	"testing"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

// ---------------------------------------------------------------------------
// RewriteFilename
// ---------------------------------------------------------------------------

func TestRewriteFilenameFull(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var input = "path/to/file.txt"
	var output = fileutil.RewriteFilename(input, &fileutil.RewriteOpts{
		Dirname: "other/path/to",
		Prefix:  "prefix-",
		Suffix:  "-suffix",
		Extname: ".csv",
	})
	assert.That(output, p.Eq("other/path/to/prefix-file-suffix.csv"))
}

func TestRewriteFilenameNoDotExt(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var input = "path/to/file.txt"
	var output = fileutil.RewriteFilename(input, &fileutil.RewriteOpts{
		Dirname: "other/path/to",
		Prefix:  "prefix-",
		Suffix:  "-suffix",
		Extname: "csv",
	})
	assert.That(output, p.Eq("other/path/to/prefix-file-suffix.csv"))
}
