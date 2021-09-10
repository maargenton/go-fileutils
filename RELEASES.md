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