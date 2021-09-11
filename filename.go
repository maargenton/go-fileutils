package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------------------------
// Local wrappers of filepath package function, preserving trailing path
// separator, commonly used to indicate a directory.

// // Clean is equivalent to `filepath.Clean()`, but preserves any trailing path
// // separator or appends one for '.' or '..' path fragments.
// func Clean(input string) string {
// 	dir := input == "" || input == "." || input == ".." ||
// 		strings.HasSuffix(input, string(filepath.Separator)+"..") ||
// 		strings.HasSuffix(input, string(filepath.Separator)+".") ||
// 		hasTrailingSeparator(input)

// 	output := filepath.Clean(input)
// 	if dir && !hasTrailingSeparator(output) {
// 		output += string(filepath.Separator)
// 	}
// 	return output
// }

// Rel is equivalent to `filepath.Rel()`, but preserves any trailing path
// separator.
func Rel(basepath, targetpath string) (string, error) {
	output, err := filepath.Rel(basepath, targetpath)
	if err == nil && hasTrailingSeparator(targetpath) && !hasTrailingSeparator(output) {
		output += string(filepath.Separator)
	}
	return output, err
}

// Join provides functionality similar to `filepath.Join()`, but with
// significant differences that make it more convenient to use in common cases.
// It takes any number of path elements and joins them with the path separator
// in between. If any element is an absolute path, all preceding elements are
// discarded and the resulting path is absolute. It also preserves any trailing
// path separator on the last element.
func Join(elem ...string) string {
	var output strings.Builder
	for _, e := range elem {
		if filepath.IsAbs(e) {
			output.Reset()
		}
		if output.Len() > 0 {
			output.WriteRune(filepath.Separator)
		}
		output.WriteString(Clean(e))
	}
	return Clean(output.String())
}

func hasTrailingSeparator(path string) bool {
	l := len(path)
	return l > 0 && path[l-1] == filepath.Separator
}

// ---------------------------------------------------------------------------

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
	if input == "~" || input == "~/" {
		output, err = os.UserHomeDir()
		if err != nil {
			err = fmt.Errorf("failed to expand path '%v', %w", input, err)
		}
		return
	}

	output = input
	if strings.HasPrefix(input, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to expand path '%v', %w", input, err)
		}
		output = filepath.Join(home, input[2:])
	}

	output = os.ExpandEnv(output)
	if !filepath.IsAbs(output) {
		output = filepath.Join(basepath, output)
	}
	output, err = filepath.Abs(output)
	if err != nil {
		return "", fmt.Errorf("failed to expand path '%v', %w", input, err)
	}

	return
}
