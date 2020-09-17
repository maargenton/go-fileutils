package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/maargenton/fileutil"
)

// ---------------------------------------------------------------------------

// Request ...
type Request struct {
	client *Client
	ctx    context.Context
	done   chan struct{}
	status Status
	mtx    sync.Mutex
	err    error

	// values set during request initialization -- no sync requires
	requestURL     *url.URL
	outputPath     string
	outputFilename string

	// values set during request handling -- require sync
	downloadURL *url.URL
	currentSize uint64
	totalSize   uint64
	canResume   bool
}

// Wait waits for the request to complete and returns an error if any
func (r *Request) Wait() error {
	<-r.done
	return r.err
}

// RequestURL returns the original request URL
func (r *Request) RequestURL() *url.URL {
	return r.requestURL
}

// OutputPath retruns the output path defined for the client at the time the
// request was created
func (r *Request) OutputPath() string {
	return r.outputPath
}

// OutputFilename returns the computed output filename based on the request URL
// an potentially the redirect URLs being followed.
func (r *Request) OutputFilename() string {
	return r.outputFilename
}

// DownloadURL returns the URL the content is downloaded from, which might be
// either the original request URL or the final target of a list of redirects.
// The value is initially nil, and is set to its final value at the end of the
// initial phase of the request.
func (r *Request) DownloadURL() *url.URL {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.downloadURL
}

// Status returns the current status of the request
func (r *Request) Status() Status {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return StatusInitial
}

// ---------------------------------------------------------------------------
// Request internal worker implementation

func (r *Request) setup(requestURL string) {
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		err = fmt.Errorf("failed to parse requestURL '%v', %w", requestURL, err)
		r.fail(err)
		return
	}
	r.requestURL = parsedURL
	r.outputFilename = filepath.Join(r.outputPath, path.Base(r.requestURL.Path))
}

func (r *Request) start() {
	err := r.run()
	if err != nil {
		r.fail(err)
	}
}

func (r *Request) run() error {
	r.checkLocalFile()
	if r.status == StatusCompleted {
		return nil
	}

	resp, err := r.getHead()
	if err != nil {
		return err
	}

	err = r.downloadContent(resp.Request)
	if err != nil {
		return err
	}

	return errors.New("Unimplemented")
}

const downloadSuffix = ".download"
const failedSuffix = ".failed"

func (r *Request) checkLocalFile() error {
	if fileutil.IsFile(r.outputFilename) {
		// destination file exists and is therefore assumed valid
		r.setStatus(StatusCompleted)
		return nil
	}

	if fileutil.IsFile(r.outputFilename + failedSuffix) {
		var s strings.Builder
		err := fileutil.ReadFile(r.outputFilename+failedSuffix, func(r io.Reader) error {
			_, err := io.Copy(&s, r)
			return err
		})
		if err != nil {
			return err
		}
		return errors.New(s.String())
	}

	info, err := os.Stat(r.outputFilename + downloadSuffix)
	if err == nil && !info.IsDir() {
		r.currentSize = uint64(info.Size())
	}

	return nil
}

func (r *Request) getHead() (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", r.requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "*/*")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentLength, err := strconv.ParseUint(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return resp, fmt.Errorf("invalid Content-Length header, %w", err)
	}

	var canResume = false
	if resp.Header.Get("Accept-Ranges") == "bytes" {
		canResume = true
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.totalSize = contentLength
	r.downloadURL = resp.Request.URL
	r.canResume = canResume

	return resp, nil
}

func (r *Request) downloadContent(req *http.Request) (err error) {
	contentRequest := *req
	contentRequest.Method = "GET"

	if r.canResume && r.currentSize != 0 {
		contentRequest.Header.Set("Range", fmt.Sprintf("bytes=%d-", r.currentSize))
	}
	resp, err := http.DefaultClient.Do(&contentRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := &trackingReader{
		ctx:   r.ctx,
		count: r.currentSize,
		r:     resp.Body,
	}
	_ = reader

	os.MkdirAll(filepath.Dir(r.outputFilename), os.ModePerm)
	flags := os.O_WRONLY | os.O_CREATE
	if r.canResume {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	filename := r.outputFilename + downloadSuffix
	f, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		return fmt.Errorf("failed to open output file '%v', %w", filename, err)
	}
	defer f.Close()

	r.setStatus(StatusContent)
	_, err = io.Copy(f, reader)

	return
}

func (r *Request) fail(err error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.err = err
	r.status = StatusFailed
	close(r.done)
}

func (r *Request) setStatus(status Status) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.status != status {
		r.status = status
		if status == StatusCompleted {
			close(r.done)
		}
	}
}
