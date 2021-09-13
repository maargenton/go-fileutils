package downloader

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/maargenton/go-fileutils"
)

type tracker struct {
	size uint64
	err  error
	mtx  sync.Mutex
	cond *sync.Cond
}

func newTracker() *tracker {
	var t = &tracker{}
	t.cond = sync.NewCond(&t.mtx)
	return t
}

func (t *tracker) set(size uint64) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.size = size
	t.cond.Broadcast()
}

func (t *tracker) setError(err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.err = err
	t.cond.Broadcast()
}

func (t *tracker) add(delta int) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.size += uint64(delta)
	t.cond.Broadcast()
}

func (t *tracker) get(minSize uint64) (size uint64, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	for {
		if t.size >= minSize || t.err != nil {
			return t.size, t.err
		}
		t.cond.Wait()
	}
}

type readbackReader struct {
	f      *os.File
	t      *tracker
	offset uint64
}

func newReadbackReader(filename string, t *tracker) (r *readbackReader, err error) {
	err = fileutils.Touch(filename)
	if err != nil {
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	r = &readbackReader{
		f: f,
		t: t,
	}
	return
}

func (r *readbackReader) Read(p []byte) (n int, err error) {
	var minSize = r.offset + uint64(len(p))
	_, trackerErr := r.t.get(minSize)
	n, err = r.f.Read(p)
	r.offset += uint64(n)

	if errors.Is(err, io.EOF) {
		err = trackerErr
	}
	return
}

func (r *readbackReader) Close() error {
	return r.f.Close()
}

// --------------------------------------------------------------------------

type syncChecker struct {
	mtx  sync.Mutex
	cond *sync.Cond
}

func newSyncChecker() *syncChecker {
	var c = &syncChecker{}
	c.cond = sync.NewCond(&c.mtx)
	return c
}

func (c *syncChecker) singal(f func()) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	f()
	c.cond.Broadcast()
}

func (c *syncChecker) check(f func() bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for {
		if f() {
			return
		}
		c.cond.Wait()
	}
}
