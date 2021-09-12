package fileutils_test

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-fileutils"
)

// ---------------------------------------------------------------------------
// RewriteFilename

func TestRewriteFilenameFull(t *testing.T) {
	var input = "path/to/file.txt"
	var output = fileutils.RewriteFilename(input, &fileutils.RewriteOpts{
		Dirname: "other/path/to/",
		Prefix:  "prefix-",
		Suffix:  "-suffix",
		Extname: ".csv",
	})
	require.That(t, output).Eq("other/path/to/prefix-file-suffix.csv")
}

func TestRewriteFilenameNoDotExt(t *testing.T) {
	var input = "path/to/file.txt"
	var output = fileutils.RewriteFilename(input, &fileutils.RewriteOpts{
		Dirname: "other/path/to",
		Prefix:  "prefix-",
		Suffix:  "-suffix",
		Extname: "csv",
	})
	require.That(t, output).Eq("other/path/to/prefix-file-suffix.csv")
}

// RewriteFilename
// ---------------------------------------------------------------------------
