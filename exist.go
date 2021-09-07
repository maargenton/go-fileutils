package fileutils

import "os"

// Exists returns true if a file or directory exists at the specified location
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsFile returns true if a file exists at the specified location
func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// IsDir returns true if a directory exists at the specified location
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
