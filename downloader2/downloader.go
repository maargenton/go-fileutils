// Package downloader provides a simple interface to download large immutable
// content over http. Content is fetch from an http or https url and save to the
// local filesystem.
//
// Partial content is saved next to the final output location, in a file with a
// '.download' suffix. If a '.download' file already exists, a resume operation
// is attempted if the server supports it, otherwise the file is discarded and
// downloaded again from the beginning. Once downloaded, the checksum is
// verified if both `Hash` and `Checksum` are specified in the request, and the
// file is moved to its final localtion. If the checksum validation fails, the
// file is deleted and an error is returned.
//
// To monitor progress during the download opeartions, two handlers can be
// registered to received either status update or progress update. In addition,
// when downloading an archive or an intermediate representation of the final
// content, the content can be streamed out (always from the beginning) to a
// streaming unarchiver or a streaming processor, and minimize the overall
// operation time.
package downloader

import (
	"context"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/go-errors"
)

// Client captures the settings common to multiple request
type Client struct {
	HTTPClient      *http.Client
	UserAgent       string
	OutputDirectory string
}

// DefaultClient defines a default client using http.DefaultClient
var DefaultClient = &Client{
	HTTPClient: http.DefaultClient,
	UserAgent:  "fileutil/downloader",
}

// Request captures all request specific settings.
type Request struct {
	URL             string
	OutputFilename  string
	OutputDirectory string
	Hash            hash.Hash
	Checksum        []byte

	StatusUpdateHandler func(status int)
	ProgressHandler     func(progress float64)
	ContentReader       func(r io.Reader) error
}

// ErrInvalidURL is returned by Client.Get() when the target URL specified in
// the request cannot be retrieved due to and http error.
const ErrInvalidURL = errors.Sentinel("ErrInvalidURL")

// ErrLocalFileError is returned by Client.Get() when one of the necessary local
// filesystem operation fails while trying to fulfill the request.
const ErrLocalFileError = errors.Sentinel("ErrLocalError")

// ErrInvalidChecksum is returned when both `Hash` and `Checksum` are specified
// in the request and the actual checksum of the content computed through `Hash`
// does not match the expected `Checksum`.
const ErrInvalidChecksum = errors.Sentinel("ErrInvalidChecksum")

const downloadSuffix = ".download"

// Get fetches and saves the content of the specified url to the local
// filesystem. The function does not return until all necessary operations are
// completed or failed. The request object is modified to capture the details
// that were left unspecified, and to ensure that all paths refer to absolute
// locations.
func (c *Client) Get(ctx context.Context, r *Request) error {
	parsedURL, err := url.Parse(r.URL)
	if err != nil {
		return ErrInvalidURL.Errorf(
			"failed to parse request URL '%v', %w", r.URL, err)
	}

	if r.OutputFilename == "" {
		r.OutputFilename = path.Base(parsedURL.Path)
	}
	err = c.expandOutputPath(r)

	if fileutil.IsFile(r.OutputFilename) {
		if r.ContentReader != nil {
			// TODO: stream content to ContentReader, update progress along the
			// way
		} else if r.ProgressHandler != nil {
			r.ProgressHandler(1.0)
		}
		return nil
	}

	var currentSize uint64
	info, err := os.Stat(r.OutputFilename + downloadSuffix)
	if err == nil && !info.IsDir() {
		currentSize = uint64(info.Size())
	}

	//

	req, err := http.NewRequest("HEAD", r.URL, nil)
	if err != nil {
		return ErrInvalidURL.Wrap(err)
		//  Errorf(
		// 	"failed to create request for URL '%v', %w", r.URL, err)
	}
	if c.UserAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	req.Header.Add("Accept", "*/*")

	head, err := c.HTTPClient.Do(req)
	if err != nil {
		return ErrInvalidURL.Wrap(err)
	}
	head.Body.Close()

	var totalSize = uint64(head.ContentLength)
	var canResume = false
	if head.Header.Get("Accept-Ranges") == "bytes" {
		canResume = true
	}

	//

	var contentRequest = *head.Request
	contentRequest.Method = "GET"
	if canResume && currentSize != 0 {
		contentRequest.Header.Set("Range", fmt.Sprintf("bytes=%d-", currentSize))
	}
	resp, err := c.HTTPClient.Do(&contentRequest)
	if err != nil {
		return ErrInvalidURL.Wrap(err)
	}
	defer resp.Body.Close()

	// Check resp.StatusCode

	reader := &trackingReader{
		ctx:   ctx,
		count: currentSize,
		r:     resp.Body,
	}

	os.MkdirAll(filepath.Dir(r.OutputFilename), os.ModePerm)
	flags := os.O_WRONLY | os.O_CREATE
	if canResume {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	filename := r.OutputFilename + downloadSuffix
	f, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		return fmt.Errorf("failed to open output file '%v', %w", filename, err)
	}
	defer f.Close()

	_, err = io.Copy(f, reader)

	// return
	//

	// _ = currentSize
	_ = totalSize
	// _ = canResume
	return ErrInvalidChecksum
}

func (c *Client) expandOutputPath(r *Request) error {
	var err error
	if filepath.Base(r.OutputFilename) != r.OutputFilename {

		// OutputFilename has a path component, expand relative to current
		// directory
		r.OutputFilename, err = fileutil.ExpandPath(r.OutputFilename)
		if err != nil {
			return ErrLocalFileError.Errorf(
				"failed to expand output filename '%v', %w",
				r.OutputFilename, err)
		}
		r.OutputDirectory = filepath.Dir(r.OutputFilename)
	} else {

		// OutputFilename has a path component, expand relative to output
		// directory, from client or request
		if r.OutputDirectory == "" {
			r.OutputDirectory = c.OutputDirectory
		}
		if r.OutputDirectory == "" {
			r.OutputDirectory = "."
		}
		r.OutputDirectory, err = fileutil.ExpandPath(r.OutputDirectory)
		if err != nil {
			return ErrLocalFileError.Errorf(
				"failed to expand output directory '%v', %w", r.OutputDirectory, err)
		}
		r.OutputFilename = filepath.Join(r.OutputDirectory, r.OutputFilename)
	}
	return nil
}
