package fileutil

import (
	"path/filepath"
	"strings"
)

// RewriteOpts contains the options to apply to RewriteFilename to transform the
// input filename
type RewriteOpts struct {
	Dirname string // Replace the path with the specified dirname
	Extname string // Replace the file extension with the specified extname
	Prefix  string // Prefix to prepend on the basename
	Suffix  string // Suffix to append on the basename
}

// RewriteFilename transforms a filename according the the specified options and
// can change dirname, basename, extension or append / prepend a fragment to the
// basename.
func RewriteFilename(input string, opts *RewriteOpts) string {
	dirname, filename := filepath.Split(input)
	extname := filepath.Ext(filename)
	basename := filename[0 : len(filename)-len(extname)]

	basename = opts.Prefix + basename + opts.Suffix
	if len(opts.Extname) != 0 {
		if strings.HasPrefix(opts.Extname, ".") {
			extname = opts.Extname
		} else {
			extname = "." + opts.Extname
		}
	}
	if len(opts.Dirname) != 0 {
		dirname = opts.Dirname
	}
	return filepath.Join(dirname, basename+extname)
}
