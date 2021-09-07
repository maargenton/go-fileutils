// Package dir exposes globing and scanning functions that walk the filesystem
// and find files and folder that match specific patterns. Patterns are defined
// using an extended glob pattern that can match deep subdirectories.
//
// Extended glob patterns always use `/` as path separator and `\` as escape
// character, regardless of the OS native filename format.
//
// Extended glob patterns are interpreted as a sequence of one or more path
// fragments. Each path fragment can be matched against either a literal
// sequence of characters or a glob pattern.
//
// `*` matches zero or more occurrences of any character within a path fragment,
// `?` matches one occurrence of any character within a path fragment,
// `[<range>]` matches one occurrence of any listed character within a path
// fragment, and  `{foo,bar}` matches one occurrence of either `foo` or `bar`
// within a path fragment
//
// `**/` allows the subsequent fragment to be matched anywhere within the
// directory tree. It should always be followed by another fragment matching
// expression.
package dir
