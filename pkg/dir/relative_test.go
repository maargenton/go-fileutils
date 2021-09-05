package dir_test

import (
	"errors"
	"os"
	"testing"

	"github.com/maargenton/fileutil/dir"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

func TestMakeRelativeWalkFn(t *testing.T) {
	var testErr = errors.New("test error")
	var relErr = errors.New("Rel: can't make /aaa/bbb/ccc/ddd relative to aaa/bbb")
	var tcs = []struct {
		path          string
		err           error
		expectedPath  string
		expectedError error
	}{
		{"aaa/bbb/ccc/ddd", nil, "ccc/ddd", nil},
		{"aaa/bbb/ccc/ddd", testErr, "ccc/ddd", testErr},
		{"/aaa/bbb/ccc/ddd", nil, "/aaa/bbb/ccc/ddd", relErr},
		{"/aaa/bbb/ccc/ddd", testErr, "/aaa/bbb/ccc/ddd", testErr},
	}

	for _, tc := range tcs {
		t.Run(tc.path, func(t *testing.T) {
			assert := asserter.New(t)
			assert.That(true, p.IsTrue())

			var recordedPath string
			var recordedError error
			walkFn := dir.MakeRelativeWalkFunc("aaa/bbb",
				func(path string, info os.FileInfo, err error) error {
					recordedPath = path
					recordedError = err
					return err
				})

			walkFn(tc.path, nil, tc.err)
			assert.That(recordedPath, p.Eq(tc.expectedPath))
			assert.That(recordedError, p.Eq(tc.expectedError))
		})
	}
}
