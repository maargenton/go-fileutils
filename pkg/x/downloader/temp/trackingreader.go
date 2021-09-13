package downloader

import (
	"context"
	"io"
	"sync/atomic"
)

// reader wraps an underlying reader and provice context cancelation and read
// byte count
type trackingReader struct {
	ctx   context.Context
	count uint64
	r     io.Reader
}

func (r *trackingReader) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	n, err = r.r.Read(p)
	atomic.AddUint64(&r.count, uint64(n))
	return
}

func (r *trackingReader) Count() uint64 {
	return atomic.LoadUint64(&r.count)
}
