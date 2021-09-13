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
	"bytes"
	"context"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/maargenton/go-errors"
	"github.com/maargenton/go-fileutils"
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

	ProgressHandler func(progress float64)
	ContentReader   func(r io.Reader) error
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

	// Determine output filename based on request arguments and URL
	parsedURL, err := url.Parse(r.URL)
	if err != nil {
		return ErrInvalidURL.Errorf(
			"failed to parse request URL, %w", err)
	}
	if r.OutputFilename == "" {
		r.OutputFilename = path.Base(parsedURL.Path)
	}
	err = c.expandOutputPath(r)

	// Handle case where the final file is already there and no download is
	// necessary
	if fileutils.IsFile(r.OutputFilename) {
		if r.ContentReader != nil {
			// TODO: stream content to ContentReader, update progress along the
			// way
		} else if r.ProgressHandler != nil {
			r.ProgressHandler(-1.0)
		}
		return nil
	}

	// Setup background task to prime the hash function with the partial local
	// content while sending the initial HEAD request.
	var filename = r.OutputFilename + downloadSuffix
	var hashPrimeWG sync.WaitGroup
	var hashPrimeError error
	if fileutils.Exists(filename) && len(r.Checksum) != 0 && r.Hash != nil {
		hashPrimeWG.Add(1)
		go func() {
			defer hashPrimeWG.Done()
			hashPrimeError = hashFileContent(filename, r.Hash)
		}()
	}

	// HEAD request to fetch target information and server capability
	req, err := http.NewRequestWithContext(ctx, "HEAD", r.URL, nil)
	if err != nil {
		return ErrInvalidURL.Wrap(err)
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
	hashPrimeWG.Wait()

	if head.StatusCode >= 400 {
		return ErrInvalidURL.Errorf("unexpected http response: '%v'", head.Status)
	}
	if hashPrimeError != nil {
		return ErrLocalFileError.Wrap(hashPrimeError)
	}

	var downloader = &downloader{
		client:      c.HTTPClient,
		head:        head,
		destination: filename,
		hash:        r.Hash,
	}

	if r.ProgressHandler != nil && r.ContentReader == nil {
		downloader.progress = r.ProgressHandler
	}
	os.MkdirAll(fileutils.Dir(filename), os.ModePerm)
	if err := downloader.resume(); err != nil {
		return err
	}
	if downloader.targetSize != 0 && downloader.currentSize != downloader.targetSize {
		return io.ErrUnexpectedEOF
	}

	// Finalize download
	if r.Hash != nil && len(r.Checksum) > 0 {
		// Verify checksum
		if !bytes.Equal(r.Checksum, r.Hash.Sum(nil)) {
			if downloader.targetSize != 0 {
				if err := os.Remove(filename); err != nil {
					return ErrInvalidChecksum.Errorf("failed to remove target file, %w", err)
				}
			}
			return ErrInvalidChecksum.Errorf("content checksum mismatch")
		}
	}

	if err := os.Rename(filename, r.OutputFilename); err != nil {
		return ErrLocalFileError.Errorf(
			"failed to move downloaded content to its final location, %w", err)
	}

	if downloader.progress != nil {
		downloader.progress(-1)
	}

	return nil
}

func hashFileContent(filename string, h hash.Hash) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for checksum, %w", err)
	}
	defer f.Close()
	_, err = io.Copy(h, f)
	if err != nil {
		return fmt.Errorf("failed to read content for checksum, %w", err)
	}
	return nil
}

func (c *Client) expandOutputPath(r *Request) error {
	var err error
	if fileutils.Base(r.OutputFilename) != r.OutputFilename {

		// OutputFilename has a path component, expand relative to current
		// directory
		r.OutputFilename, err = fileutils.ExpandPath(r.OutputFilename)
		if err != nil {
			return ErrLocalFileError.Errorf(
				"failed to expand output filename '%v', %w",
				r.OutputFilename, err)
		}
		r.OutputDirectory = fileutils.Dir(r.OutputFilename)
	} else {

		// OutputFilename has a path component, expand relative to output
		// directory, from client or request
		if r.OutputDirectory == "" {
			r.OutputDirectory = c.OutputDirectory
		}
		if r.OutputDirectory == "" {
			r.OutputDirectory = "."
		}
		r.OutputDirectory, err = fileutils.ExpandPath(r.OutputDirectory)
		if err != nil {
			return ErrLocalFileError.Errorf(
				"failed to expand output directory '%v', %w", r.OutputDirectory, err)
		}
		r.OutputFilename = fileutils.Join(r.OutputDirectory, r.OutputFilename)
	}
	return nil
}

// ---------------------------------------------------------------------------

// downloader captures the result of the initial HEAD request to the target URL
// and performs the actual download to the destination file while updating the
// progress tracker. The downloader never sets the error on the tracker and
// instead returns it from the resume function.
type downloader struct {
	client      *http.Client
	head        *http.Response
	destination string
	currentSize uint64
	targetSize  uint64
	hash        hash.Hash
	progress    func(float64)
}

// func (d *downloader) canResume() bool {
// 	return d.head.Header.Get("Accept-Ranges") == "bytes"
// }

func (d *downloader) resume() error {

	// Check current size of the destination file
	info, err := os.Stat(d.destination)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if info != nil {
		d.currentSize = uint64(info.Size())
	}

	if d.head.ContentLength > 0 {
		d.targetSize = uint64(d.head.ContentLength)
	}
	if d.currentSize != 0 && d.currentSize == d.targetSize {
		if d.progress != nil {
			d.progress(1)
		}
		return nil
	}

	// Send http request to fetch content, resuming if supported
	var canResume = d.head.Header.Get("Accept-Ranges") == "bytes"
	var resume = canResume && d.currentSize > 0
	var request = *d.head.Request
	request.Method = "GET"
	if resume {
		request.Header.Set("Range", fmt.Sprintf("bytes=%d-", d.currentSize))
	}
	resp, err := d.client.Do(&request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return ErrInvalidURL.Errorf("unexpected http response: '%v'", resp.Status)
	}
	if resume && resp.StatusCode != http.StatusPartialContent {
		resume = false
	}

	// Last chance to set targetSize based on response headers
	if d.targetSize == 0 {
		if resume {
			_, _, _, l := decodeContentRange(resp.Header.Get("Content-Range"))
			if l > 0 {
				d.targetSize = uint64(l)
			}
		} else {
			if resp.ContentLength > 0 {
				d.targetSize = uint64(resp.ContentLength)
			}
		}
	}

	// Open destination file for writing, truncating it if resume failed
	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if !resume {
		flags |= os.O_TRUNC
	}
	f, err := os.OpenFile(d.destination, flags, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	var w = fileutils.WriterFunc(func(p []byte) (n int, err error) {
		n, err = f.Write(p)
		d.currentSize += uint64(n)

		if d.hash != nil {
			d.hash.Write(p)
		}
		if d.progress != nil {
			d.progress(float64(d.currentSize) / float64(d.targetSize))
		}

		return
	})
	_, err = io.Copy(w, resp.Body)

	return err
}
