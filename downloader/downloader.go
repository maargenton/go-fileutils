package downloader

// TokenPool provides an abstract interface to a shared pool of tokens that can
// be acquired and released
type TokenPool interface {
	Acquire() Token
}

// Token provides interface for releasable tokens
type Token interface {
	Release()
}

type tokenPool chan struct{}
type token struct {
	pool tokenPool
}

// NewTokenPool creates a new pool of given size
func NewTokenPool(n int) TokenPool {
	return make(tokenPool, n)
}

func (tp tokenPool) Acquire() Token {
	tp <- struct{}{}
	return &token{tp}
}

func (t *token) Release() {
	if t.pool == nil {
		panic("token already released")
	}
	<-t.pool
	t.pool = nil
}
