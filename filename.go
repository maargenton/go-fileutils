package fileutils

import (
	"fmt"
	"os"
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
	dirname, filename := Split(input)
	extname := Ext(filename)
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
	return Join(dirname, basename+extname)
}

// ExpandPath is similar to ExpandPathRelative with an empty `basepath`;
// relative paths are expanded relative to `$(pwd)`.
func ExpandPath(input string) (output string, err error) {
	return expandPath(input, "")
}

// ExpandPathRelative returns the absolute path for the given input, expanding
// environment variable and handling the special case `~/` referring to the
// current user home directory. If the resulting path after variable expansion
// is relative, it is expanded relative to `basepath`. If the resulting path is
// still relative, it is expanded relative to `$(pwd)`.  The function returns
// and error if one of the underlying calls fails (getting user home or process
// working directory path).
func ExpandPathRelative(input, basepath string) (output string, err error) {
	return expandPath(input, basepath)
}

func expandPath(input, basepath string) (output string, err error) {
	output = Clean(input)
	output = strings.ReplaceAll(output, "~/", "${HOME}/")
	output = os.ExpandEnv(output)
	if !IsAbs(output) {
		output = Join(basepath, output)
	}
	output, err = Abs(output)
	if err != nil {
		return "", fmt.Errorf("failed to expand path '%v', %w", input, err)
	}
	return
}
