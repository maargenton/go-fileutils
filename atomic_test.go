package fileutils_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

func TestOpenTemp(t *testing.T) {
	assert := asserter.New(t, asserter.AbortOnError())
	dir, err := ioutil.TempDir(".", "testdata-")
	assert.That(err, p.IsNoError())
	defer os.RemoveAll(dir) // clean up

	filename := filepath.Join(dir, "file.txt")
	f, err := fileutils.OpenTemp(filename, "tmp")
	if f != nil {
		defer os.Remove(f.Name())
	}
	defer f.Close()

	assert.That(err, p.IsNoError())
	assert.That(f, p.IsNotNil())
	assert.That(f.Name(), p.StartsWith(filepath.Join(dir, "file")))
}

// ---------------------------------------------------------------------------

type Content struct {
	Seq int `json:"seq,omitempty"`
}

func TestReadWriteFile(t *testing.T) {
	assert := asserter.New(t)
	dir, err := ioutil.TempDir(".", "testdata-")
	assert.That(err, p.IsNoError())
	defer os.RemoveAll(dir) // clean up

	// Write file
	var content = &Content{Seq: 125}
	filename := filepath.Join(dir, "file.txt")
	err = fileutils.WriteFile(filename, func(w io.Writer) error {
		return json.NewEncoder(w).Encode(content)
	})
	assert.That(err, p.IsNoError())

	// Read file
	content = &Content{}
	err = fileutils.ReadFile(filename, func(r io.Reader) error {
		return json.NewDecoder(r).Decode(content)
	})

	assert.That(err, p.IsNoError())
	assert.That(content.Seq, p.Eq(125))
}

func TestWriteFileIsAtomic(t *testing.T) {
	assert := asserter.New(t)
	dir, err := ioutil.TempDir(".", "testdata-")
	assert.That(err, p.IsNoError())
	defer os.RemoveAll(dir) // clean up

	filename := filepath.Join(dir, "file.txt")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)

		go func() {
			content := &Content{Seq: i}
			err := fileutils.WriteFile(filename, func(w io.Writer) error {
				return json.NewEncoder(w).Encode(content)
			})
			assert.That(err, p.IsNoError())
			wg.Done()
		}()
	}

	wg.Wait()
	content := &Content{}
	err = fileutils.ReadFile(filename, func(r io.Reader) error {
		return json.NewDecoder(r).Decode(content)
	})

	assert.That(err, p.IsNoError())
	assert.That(content.Seq, p.Ge(0))
	assert.That(content.Seq, p.Lt(10))
}
