package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"

	"github.com/maargenton/fileutil/pkg/downloader"
)

var testURL = "https://github.com/jung-kurt/gofpdf/archive/v2.17.2.tar.gz"
var checksumStr = "54c3dd9d981de133151c4c6a15a081609c19d754d40099eed7bb18645bab5a1d"

func run() error {
	checksum, err := hex.DecodeString(checksumStr)
	if err != nil {
		return fmt.Errorf("invalid checksum: %w", err)
	}
	_ = checksum

	client := downloader.DefaultClient
	err = client.Get(context.Background(), &downloader.Request{
		URL:             testURL,
		OutputDirectory: "testdata/output",
		Hash:            sha256.New(),
		Checksum:        checksum,
		ProgressHandler: progress(),
	})
	return err
}

func progress() func(p float64) {
	var progress = -100.0
	return func(p float64) {
		if math.Abs(p-progress) > 0.01 {
			progress = p
			if progress < 0 {
				fmt.Printf("Progress: done!\n")
			} else {
				fmt.Printf("Progress: %0.1f%%\n", progress*100)
			}
		}
	}
}

func main() {
	fmt.Printf("Downloading '%v' ...\n", testURL)
	if err := run(); err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		fmt.Printf("done\n")
	}
}
