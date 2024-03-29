package fileutils_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-fileutils"
)

func TestOpenTemp(t *testing.T) {
	dir, err := ioutil.TempDir(".", "testdata-")
	require.That(t, err).IsNil()
	defer os.RemoveAll(dir) // clean up

	filename := fileutils.Join(dir, "file.txt")
	f, err := fileutils.OpenTemp(filename, "tmp")
	if f != nil {
		defer os.Remove(f.Name())
	}
	defer f.Close()

	require.That(t, err).IsNil()
	require.That(t, f).IsNotNil()
	require.That(t, f.Name()).StartsWith(fileutils.Join(dir, "file"))
}

// ---------------------------------------------------------------------------

type Content struct {
	Seq int `json:"seq,omitempty"`
}

func TestReadWriteFile(t *testing.T) {
	dir, err := ioutil.TempDir(".", "testdata-")
	verify.That(t, err).IsNil()
	defer os.RemoveAll(dir) // clean up

	// Write file
	var content = &Content{Seq: 125}
	filename := fileutils.Join(dir, "file.txt")
	err = fileutils.WriteFile(filename, func(w io.Writer) error {
		return json.NewEncoder(w).Encode(content)
	})
	verify.That(t, err).IsNil()

	// Read file
	content = &Content{}
	err = fileutils.ReadFile(filename, func(r io.Reader) error {
		return json.NewDecoder(r).Decode(content)
	})

	verify.That(t, err).IsNil()
	verify.That(t, content.Seq).Eq(125)
}

func TestWriteFileIsAtomic(t *testing.T) {
	dir, err := ioutil.TempDir(".", "testdata-")
	verify.That(t, err).IsNil()
	defer os.RemoveAll(dir) // clean up

	filename := fileutils.Join(dir, "file.txt")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)

		go func() {
			content := &Content{Seq: i}
			err := fileutils.WriteFile(filename, func(w io.Writer) error {
				return json.NewEncoder(w).Encode(content)
			})

			if runtime.GOOS != "windows" {
				// Windows is tripping all over itself where there is any kind
				// of contention on the filesystem, so it often returns spurious
				// bogus permissions error that we shall ignore here
				verify.That(t, err).IsNil()
			}
			wg.Done()
		}()
	}

	wg.Wait()
	content := &Content{}
	err = fileutils.ReadFile(filename, func(r io.Reader) error {
		return json.NewDecoder(r).Decode(content)
	})

	verify.That(t, err).IsNil()
	verify.That(t, content.Seq).Ge(0)
	verify.That(t, content.Seq).Lt(10)
}
