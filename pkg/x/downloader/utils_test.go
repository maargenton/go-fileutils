package downloader

// import (
// 	"testing"

// 	"github.com/maargenton/go-testpredicate/pkg/verify"
// )

// // ---------------------------------------------------------------------------
// // decodeContentRange

// func TestDecodeContentRange(t *testing.T) {
// 	var tcs = []struct {
// 		value   string
// 		unit    string
// 		s, e, l int64
// 	}{
// 		{"bytes 2000000-2785937/2785938", "bytes", 2000000, 2785937, 2785938},
// 		{"bytes 2000000-2785937/*", "bytes", 2000000, 2785937, -1},
// 		{"bytes */2785938", "bytes", -1, -1, 2785938},
// 		{"bytes", "", -1, -1, -1},
// 	}

// 	for _, tc := range tcs {
// 		t.Run(tc.value, func(t *testing.T) {
// 			unit, s, e, l := decodeContentRange(tc.value)
// 			verify.That(t, unit).Eq(tc.unit)
// 			verify.That(t, s).Eq(tc.s)
// 			verify.That(t, e).Eq(tc.e)
// 			verify.That(t, l).Eq(tc.l)
// 		})
// 	}
// }
