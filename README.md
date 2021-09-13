# go-fileutils

A collection of filesystem utilities for Go

[![Latest](
  https://img.shields.io/github/v/tag/maargenton/go-fileutils?color=blue&label=latest&logo=go&logoColor=white&sort=semver)](
  https://pkg.go.dev/github.com/maargenton/go-fileutils)
[![Build](
  https://img.shields.io/github/workflow/status/maargenton/go-fileutils/build?label=build&logo=github&logoColor=aaaaaa)](
  https://github.com/maargenton/go-fileutils/actions?query=branch%3Amaster)
[![Codecov](
  https://img.shields.io/codecov/c/github/maargenton/go-fileutils?label=codecov&logo=codecov&logoColor=aaaaaa&token=fVZ3ZMAgfo)](
  https://codecov.io/gh/maargenton/go-fileutils)
[![Go Report Card](
  https://goreportcard.com/badge/github.com/maargenton/go-fileutils)](
  https://goreportcard.com/report/github.com/maargenton/go-fileutils)


---------------------------

Package `fileutils` is a collection of filename manipulation and filesystem
utilities including directory traversal with symlinks support, finding file and
folders with extended glob pattern, and atomic file operations.

To help support non-unix platforms, it also includes ad set of functions that
are similar to those found in package "path/filepath", but represents all path
using '/' as a separator, and preserves a trailing path separator commonly used
to represent directory names.


## Installation

    go get github.com/maargenton/go-fileutils

## Key features

### Filenames with consistent '/' separator

All unix platforms use '/' as path separator, and while windows recommends using
`\\`, it also accepts paths with regular forward slash as path separators. For
that reason, this package takes the stance of always using forward slash as path
separator. The immediate benefit is that all relative paths become platform
agnostic, freeing the cross-platform client code from having to deal with
special cases for windows.

The notion of absolute path remains different across platforms, but they can
still be manipulated safely and consistently without having to deal with
platform-specific special cases in most instances.


### Atomic file operations

- `fileutils.Write()` atomically creates or replaces the destination file with
  the content written into the io.Writer passed to the closure. This guaranties
  that readers of that file will never see an incomplete or partially updated
  content.
- `fileutils.Read()` reads the content of a file through the io.Reader passed to
  the closure.
- `fileutils.OpenTemp()` creates and opens a temporary file, in the same location
  and with the same extension as the target file. The resulting file is
  guarantied to not previously exists, and therefore never steps onto another
  file.

### Filename manipulation

- `fileutils.RewriteFilename()` is a single function that lets you transform a
  filename in many common ways, like replacing either the extension or the
  containing directory, or inserting a prefix or suffix onto the basename of the
  file.
- `fileutils.ExpandPath()` and `fileutils.ExpandPathRelative()` expand an relative
  or absolute path into an absolute path, handling `~/` and environment variable
  expansion, using ether `$(pwd)` or a given `basepath` as base path.
- `fileutils.Clean()`, `fileutils.Rel()` and `fileutils.Join()` are equivalent to
  their `filepath` counterpart, but preserve any trailing path separator,
  commonly used to indicate a directory. In addition, `fileutils.Join()` properly
  handles the case where one of the elements is an absolute path, resulting in
  an absolute path with all preceding elements ignored.

### Filesystem scanning and globing

`fileutils.Walk()` implements an enhanced version of `filepath.WalkDir()` that
follows symlinks safely and adds some flexibility in the way paths are reported.

`dir.Glob()` and `dir.Scan()` are convenient functions to locate and
enumerate files matching a particular pattern. The pattern is specified as an
extended glob pattern that can match deep subdirectories and alternative
patterns:

- Extended glob patterns always use `/` as path separator and `\` as escape
  character, regardless of the OS native filename format.
- Extended glob patterns are interpreted as a sequence of one or more path
  fragments. Each path fragment can be matched against a literal sequence of
  characters or a glob pattern.
- `*` matches zero or more occurrences of any character within a path fragment
- `?` matches one occurrence of any character within a path fragment
- `[<range>]`: matches one occurrence of any listed character within a path
  fragment
- `{foo,bar}` matches one occurrence of either `foo` or `bar` within a path
  fragment
- `**/` allows the subsequent fragment to be matched anywhere within the
  directory tree. It should always be followed by another fragment matching
  expression.

Symbolic links are followed safely as needed, emitting an `ErrRecursiveSymlink`
each time a filesystem location is visited again.

#### Examples

- `src/**/*_test.{c,cc,cpp}` : From `src`, find all files in any sub-directory
  with an `_test` suffix and a `.c`, `.cc` or `.cpp` extension.
- `src/**_test.cpp` is that same as `src/*_test.cpp`; the double star is
  interpreted as two consecutive matches of zero or more.


### Sub-process execution

`popen.Command` is an additional layer of abstraction over exec.Command aimed at
simplifying common uses where the output of the process is captured or
redirected. Unlike `exec.Command`, all the details of the command to run and
what to do with its outputs are captured in public fields of the `Command`
structure. The output streams, stdout and stderr, can be returned as a string,
redirected to a file or stream-processed through an `io.Reader`. If the process
is executed successfully but returns a non-zero exit status, the returned error
is an exec.ExitError that contains the actual status code.

The behavior of stdout and stderr is controlled by 3 similar variables:

- When `WriteStdoutToFile` is set to the path of a destination file for the
  content of the command stdout, `DiscardStdout` is ignored and the returned
  stdout string is always empty. If needed, the output of the command can be
  read back from that file.
- When `StdoutReader` is set, the raw output of the command is still captured
  and returned in the stdout string, unless `DiscardStdout` is set to `true`.
- `WriteStdoutToFile` and `StdoutReader` can both be set, in which case the
  output of the command is sent to both and the returned stdout string is empty.

Except for `StdoutReader` and `StderrReader` which are most likely stateful, the
command object is stateless and can potentially be `Run()` multiple times,
concurrently.
