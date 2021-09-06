package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/maargenton/fileutil/pkg/x/downloader"
)

var testURL = "https://developer.arm.com/-/media/Files/downloads/gnu-rm/9-2019q4/gcc-arm-none-eabi-9-2019-q4-major-mac.tar.bz2?revision=c2c4fe0e-c0b6-4162-97e6-7707e12f2b6e&amp;la=en&amp;hash=EC9D4B5F5B050267B924F876B306D72CDF3BDDC0"
var checksumStr = "1249f860d4155d9c3ba8f30c19e7a88c5047923cea17e0d08e633f12408f01f0"

func run() error {
	checksum, err := hex.DecodeString(checksumStr)
	if err != nil {
		return fmt.Errorf("invalid checksum: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := downloader.DefaultClient
	return client.Get(ctx, &downloader.Request{
		URL:             testURL,
		OutputDirectory: "testdata/output",
		Hash:            sha256.New(),
		Checksum:        checksum,
		ProgressHandler: progress(),
	})
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
