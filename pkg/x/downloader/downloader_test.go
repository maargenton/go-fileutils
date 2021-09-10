package downloader_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-fileutils/pkg/x/downloader"
)

// ---------------------------------------------------------------------------
// Test Request.OutputFilename determination / expansion

func TestGetWithInvalidUrl(t *testing.T) {
	var client = transportErrorClient()
	err := client.Get(context.Background(), &downloader.Request{
		URL: ":http://localhost:8080/aaa/bbb.tar.gz",
	})

	verify.That(t, err).IsError(downloader.ErrInvalidURL)
}

func TestGetWithBareFilename(t *testing.T) {
	var client = transportErrorClient()
	var request = &downloader.Request{
		URL:             "http://localhost:8080/aaa/bbb.tar.gz",
		OutputFilename:  "bbb.tar.gz",
		OutputDirectory: "__output__",
	}
	err := client.Get(context.Background(), request)

	verify.That(t, err).IsError(mockError)
	verify.That(t, request.OutputFilename).Contains("__output__")
}

func TestGetWithPathSpecifiedFilename(t *testing.T) {
	var client = transportErrorClient()
	client.OutputDirectory = "/tmp/__client__/__output__"

	var request = &downloader.Request{
		URL:             "http://localhost:8080/aaa/bbb.tar.gz",
		OutputFilename:  "__path__/bbb.tar.gz",
		OutputDirectory: "__output__",
	}
	err := client.Get(context.Background(), request)

	verify.That(t, err).IsError(mockError)
	verify.That(t, request.OutputFilename).Contains("__path__")

	var expectedPath = filepath.Join(pwd, "__path__/bbb.tar.gz")
	verify.That(t, request.OutputFilename).Eq(expectedPath)
}

func TestGetUsingClientOutputDirectoryAndURLFilename(t *testing.T) {
	var client = transportErrorClient()
	client.OutputDirectory = "/tmp/__client__/__output__"

	var request = &downloader.Request{
		URL: "http://localhost:8080/aaa/bbb.tar.gz",
	}
	err := client.Get(context.Background(), request)

	verify.That(t, err).IsError(mockError)
	verify.That(t, request.OutputFilename).Eq("/tmp/__client__/__output__/bbb.tar.gz")
}

func TestGetUsingRequestOutputDirectoryAndURLFilename(t *testing.T) {
	var client = transportErrorClient()
	client.OutputDirectory = "/tmp/__client__/__output__"

	var request = &downloader.Request{
		URL:             "http://localhost:8080/aaa/bbb.tar.gz",
		OutputDirectory: "/tmp/__request__/__output__",
	}
	err := client.Get(context.Background(), request)

	verify.That(t, err).IsError(mockError)
	verify.That(t, request.OutputFilename).Eq("/tmp/__request__/__output__/bbb.tar.gz")
}

func TestGetWithNoOutputDirectoryUsesWorkingDirectory(t *testing.T) {
	var client = transportErrorClient()
	var request = &downloader.Request{
		URL: "http://localhost:8080/aaa/bbb.tar.gz",
	}
	err := client.Get(context.Background(), request)

	verify.That(t, err).IsError(mockError)
	var expected = filepath.Join(pwd, "bbb.tar.gz")
	verify.That(t, request.OutputFilename).Eq(expected)
}

// ---------------------------------------------------------------------------
// Test Request.OutputFilename determination / expansion
