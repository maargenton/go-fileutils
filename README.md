# fileutil

Go filesystem utilities

[![GoDoc](
  https://godoc.org/github.com/maargenton/fileutil?status.svg)](
  https://godoc.org/github.com/maargenton/fileutil)
[![Build Status](
  https://travis-ci.org/maargenton/fileutil.svg?branch=master)](
  https://travis-ci.org/maargenton/fileutil)
[![codecov](
  https://codecov.io/gh/maargenton/fileutil/branch/master/graph/badge.svg)](
  https://codecov.io/gh/maargenton/fileutil)
[![Go Report Card](
  https://goreportcard.com/badge/github.com/maargenton/fileutil)](
  https://goreportcard.com/report/github.com/maargenton/fileutil)


Package `fileutil` is a small collection of utility functions to interact with
the filesystem. The functions provided are slightly higher level than those
provided by the `filepath` package.

## Installation

    go get github.com/maargenton/fileutil

## Key features

### Atomic file operations

- `fileutil.Write()` atomically creates or replaces the destination file with
  the content written into the io.Writer passed to the closure. This guaranties
  that readers of that file will never see an incomplete or partially updated
  content.
- `fileutil.Read()` reads the content of a file through the io.Reader passed to
  the closure.
- `fileutil.OpenTemp()` creates and opens a temporary file, in the same location
  and with the same extension as the target file. The resulting file is
  guarantied to not previously exists, and therefore never steps onto another
  file.

### Filename manipulation

- `fileutil.RewriteFilename()` is a single function that lets you transform a
  filename in many common ways, like replacing either the extension or the
  containing directory, or inserting a prefix or suffix onto the basename of the
  file.
- `fileutil.ExpandPath()` and `fileutil.ExpandPathRelative()` expand an relative
  or absolute path into an absolute path, handling `~/` and environment variable
  expansion, using ether `$(pwd)` or a given `basepath` as base path.

### Filesystem scanning and globing

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

#### Examples

- `src/**/*_test.{c,cc,cpp}` : From `src`, find all files in any sub-directory
  with an `_test` suffix and a `.c`, `.cc` or `.cpp` extension.
- `src/**_test.cpp` is that same as `src/*_test.cpp`; the double star is
  interpreted as two consecutive matches of zero or more.
