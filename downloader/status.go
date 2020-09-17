package downloader

import (
	"fmt"
	"strings"
)

// Status represents the current status of a download request
type Status int

const (
	// StatusInitial covers the initial phase of the download request where the
	// initial request is made to the server to follow redirects, retrieve
	// meta-data (filename and size) and resume capability.
	StatusInitial Status = iota

	// StatusContent covers the main phase of the download request where the
	// content of the file is downloaded.
	StatusContent

	// StatusFinal covers the final phase of the download request where the
	// content is eventually verified against a checksum and the downloaded file
	// is moved to its final location.
	StatusFinal

	// StatusCompleted indicates that the download request has successfully
	// completed and is available locally as a file.
	StatusCompleted

	// StatusFailed indicates that the download request is no longer in progress
	// but could not complete successfully.
	StatusFailed
)

func (s Status) String() string {
	switch s {
	case StatusInitial:
		return "Initial"
	case StatusContent:
		return "Content"
	case StatusFinal:
		return "Final"
	case StatusCompleted:
		return "Completed"
	case StatusFailed:
		return "Failed"
	}
	return "InvalidStatus"
}

// Set implements flags.Value interface
func (s *Status) Set(value string) error {
	switch strings.ToLower(value) {
	case "initial":
		*s = StatusInitial
	case "content":
		*s = StatusContent
	case "final":
		*s = StatusFinal
	case "completed":
		*s = StatusCompleted
	case "failed":
		*s = StatusFailed
	default:
		return fmt.Errorf("invalid Status value '%v'", value)
	}
	return nil
}
