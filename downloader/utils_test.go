package downloader

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

// ---------------------------------------------------------------------------
// decodeContentRange

func TestDecodeContentRange(t *testing.T) {
	var tcs = []struct {
		value   string
		unit    string
		s, e, l int64
	}{
		{"bytes 2000000-2785937/2785938", "bytes", 2000000, 2785937, 2785938},
		{"bytes 2000000-2785937/*", "bytes", 2000000, 2785937, -1},
		{"bytes */2785938", "bytes", -1, -1, 2785938},
		{"bytes", "", -1, -1, -1},
	}

	for _, tc := range tcs {
		t.Run(tc.value, func(t *testing.T) {
			assert := asserter.New(t)

			unit, s, e, l := decodeContentRange(tc.value)
			assert.That(unit, p.Eq(tc.unit))
			assert.That(s, p.Eq(tc.s))
			assert.That(e, p.Eq(tc.e))
			assert.That(l, p.Eq(tc.l))
		})
	}
}
