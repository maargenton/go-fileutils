# v0.6.4

- Fix issue with `dir.Glob('**/*')` not scanning subdirectories.
- Move `Walk()` function and associated definition to `dir` sub-package.

## Code changes

- Move `Walk()` function to `dir` sub-package ([#13](https://github.com/maargenton/go-fileutils/pull/13))
- Fix issue with `dir.Glob('**/*')` not scanning subdirectories ([#14](https://github.com/maargenton/go-fileutils/pull/14))


# v0.6.3

- `popen.Command` now supports graceful shutdown of the child process or child
  process group on Unix platforms when the associated context is canceled.
- Fix handling of `./` prefix in `dir.GlobMatcher`, `dir.Glob...()` and
  `dir.Scan()`.

## Code changes

- Add graceful shutdown options to popen.Command ([#10](https://github.com/maargenton/go-fileutils/pull/10))
- Fix handling of `./` prefix in glob pattern and filesystem scanning ([#11](https://github.com/maargenton/go-fileutils/pull/11))


# v0.6.2

- `dir.Glob()` and associated function can now match patterns that are an
  explicit filename instead of a glob pattern.

# v0.6.1

- Code cleanup, no changes

# v0.6.0

## Major changes

- Fix all major compatibility issues and inconsistencies with Windows platform.
- Provide full set of functions dealing with filenames, using '/' as path
  separator independently of the platform.
- Remove experimental sub-packages, previously located under `pkg/x`

# v0.5.0

## Major changes

- Rename project to `github.com/maargenton/go-fileutils`
- Switch `dir.Glob()` and all to using new `filepath.WalkDir()` (from go v1.16)
  for efficiency.
- Move rewrite of `dir.Walk()` to main package `fileutils.Walk()`
- Fix inconsistencies in paths returned by walk and glob functions, including
  adding a trailing path separator for directory filenames.
- Move additional packages under `pkg/`
- Move experimental packages under `pkg/x/`
