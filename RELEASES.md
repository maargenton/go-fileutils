# v0.4.4

## Major changes

- Switch `dir.Glob()` and all to using new `filepath.WalkDir()` (from go v1.16)
  for efficiency.
- Move rewrite of `dir.Walk()` to main package `filetuils.Walk()`
- Fix inconsistencies in paths returned by walk and glob functions, including
  adding a trailing path separator for directory filenames.
