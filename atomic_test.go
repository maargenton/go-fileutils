package fileutil_test

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

func TestOpenTemp(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	f, err := fileutil.OpenTemp("testdata/file.txt", "tmp")
	defer os.Remove(f.Name())
	defer f.Close()

	assert.That(err, p.IsNoError())
	assert.That(f, p.IsNotNil())
	assert.That(f.Name(), p.StartsWith("testdata/file-tmp"))
}

// ---------------------------------------------------------------------------

type Content struct {
	Seq int `json:"seq,omitempty"`
}

func TestWriteFile(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var content = &Content{Seq: 125}
	err := fileutil.WriteFile("testdata/file.txt", func(w io.Writer) error {
		return json.NewEncoder(w).Encode(content)
	})

	assert.That(err, p.IsNoError())
}

func TestReadFile(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var content = &Content{}
	err := fileutil.ReadFile("testdata/file.txt", func(r io.Reader) error {
		return json.NewDecoder(r).Decode(content)
	})

	assert.That(err, p.IsNoError())
	assert.That(content.Seq, p.Eq(125))
}
