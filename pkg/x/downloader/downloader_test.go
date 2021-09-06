package downloader_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"

	"github.com/maargenton/fileutil/pkg/x/downloader"
)

// ---------------------------------------------------------------------------
// Test Request.OutputFilename determination / expansion

func TestGetWithInvalidUrl(t *testing.T) {
	assert := asserter.New(t)

	var client = transportErrorClient()
	err := client.Get(context.Background(), &downloader.Request{
		URL: ":http://localhost:8080/aaa/bbb.tar.gz",
	})

	assert.That(err, p.IsError(downloader.ErrInvalidURL))
}

func TestGetWithBareFilename(t *testing.T) {
	assert := asserter.New(t)

	var client = transportErrorClient()
	var request = &downloader.Request{
		URL:             "http://localhost:8080/aaa/bbb.tar.gz",
		OutputFilename:  "bbb.tar.gz",
		OutputDirectory: "__output__",
	}
	err := client.Get(context.Background(), request)

	assert.That(err, p.IsError(mockError))
	assert.That(request.OutputFilename, p.Contains("__output__"))
}

func TestGetWithPathSpecifiedFilename(t *testing.T) {
	assert := asserter.New(t)

	var client = transportErrorClient()
	client.OutputDirectory = "/tmp/__client__/__output__"

	var request = &downloader.Request{
		URL:             "http://localhost:8080/aaa/bbb.tar.gz",
		OutputFilename:  "__path__/bbb.tar.gz",
		OutputDirectory: "__output__",
	}
	err := client.Get(context.Background(), request)

	assert.That(err, p.IsError(mockError))
	assert.That(request.OutputFilename, p.Contains("__path__"))

	var expectedPath = filepath.Join(pwd, "__path__/bbb.tar.gz")
	assert.That(request.OutputFilename, p.Eq(expectedPath))
}

func TestGetUsingClientOutputDirectoryAndURLFilename(t *testing.T) {
	assert := asserter.New(t)

	var client = transportErrorClient()
	client.OutputDirectory = "/tmp/__client__/__output__"

	var request = &downloader.Request{
		URL: "http://localhost:8080/aaa/bbb.tar.gz",
	}
	err := client.Get(context.Background(), request)

	assert.That(err, p.IsError(mockError))
	assert.That(request.OutputFilename, p.Eq("/tmp/__client__/__output__/bbb.tar.gz"))
}

func TestGetUsingRequestOutputDirectoryAndURLFilename(t *testing.T) {
	assert := asserter.New(t)

	var client = transportErrorClient()
	client.OutputDirectory = "/tmp/__client__/__output__"

	var request = &downloader.Request{
		URL:             "http://localhost:8080/aaa/bbb.tar.gz",
		OutputDirectory: "/tmp/__request__/__output__",
	}
	err := client.Get(context.Background(), request)

	assert.That(err, p.IsError(mockError))
	assert.That(request.OutputFilename, p.Eq("/tmp/__request__/__output__/bbb.tar.gz"))
}

func TestGetWithNoOutputDirectoryUsesWorkingDirectory(t *testing.T) {
	assert := asserter.New(t)

	var client = transportErrorClient()
	var request = &downloader.Request{
		URL: "http://localhost:8080/aaa/bbb.tar.gz",
	}
	err := client.Get(context.Background(), request)

	assert.That(err, p.IsError(mockError))
	var expected = filepath.Join(pwd, "bbb.tar.gz")
	assert.That(request.OutputFilename, p.Eq(expected))
}

// ---------------------------------------------------------------------------
// Test Request.OutputFilename determination / expansion
