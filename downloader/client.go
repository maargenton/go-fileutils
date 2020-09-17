package downloader

import (
	"context"
	"path/filepath"
)

// Client ...
type Client struct {
	OutputPath string
	TokenPool  TokenPool
}

// NewClient creates a new downloader client
func NewClient(outputPath string) (*Client, error) {

	outputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return nil, err
	}
	return &Client{
		OutputPath: outputPath,
	}, nil
}

// Get ...
func (c *Client) Get(ctx context.Context, requestURL string) (r *Request) {
	r = &Request{
		client:     c,
		done:       make(chan struct{}),
		outputPath: c.OutputPath,
		ctx:        ctx,
	}

	r.setup(requestURL)
	if r.err != nil {
		return
	}

	go r.start()
	return
}
