package unarchiver

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
)

// ArchiveItem represent one item within an archive. The name and attributes are
// provided as an os.FileInfo, and the content if any is accessible through an
// io.Reader.
type ArchiveItem interface {
	os.FileInfo
	FullName() string
	LinkName() string
	IsHardLink() bool // Special item type not represented in os.FileInfo.Mode
	Content() io.Reader
}

// Unarchiver defines an common interface for accessing the content of an
// archive in a sequential fashion (for extraction purposes).
type Unarchiver interface {
	// NextItem returns the next item in the archive. The item returned, and
	// specifically its content,  may become invalid on the next call to
	// NextItem() and should not be used after that point.
	NextItem() (item ArchiveItem, err error)
}

var gzipHeader = []byte("\x1f\x8b")
var bz2Header = []byte("BZh")
var zipHeader = []byte("\x04\x03\x4b\x50")

// New creates and returns a new unarchiver capable of extracting the content of
// the stream r. The format is determined automatically based on the first few
// bytes of the stream. Supported formats include zip, tar.gz and tar.bz2, each
// of which as a recognizable signature prefix. If prefix is not recognized, the
// format is assumed to be a raw tar archive, the only supported format without
// a recognizable signature.
func New(r io.Reader) (u Unarchiver, err error) {
	var header = make([]byte, 1024)
	n, err := r.Read(header)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	header = header[:n]
	r = io.MultiReader(bytes.NewReader(header), r)

	if bytes.HasPrefix(header, zipHeader) {
		return nil, fmt.Errorf("unimplemented ZIP archive format")
	}

	if bytes.HasPrefix(header, gzipHeader) {
		r, err = gzip.NewReader(r)
		if err != nil {
			return nil, fmt.Errorf("failed to create GZIP reader, %w", err)
		}
	} else if bytes.HasPrefix(header, bz2Header) {
		r = bzip2.NewReader(r)
	}

	return &tarUnarchiver{
		r: tar.NewReader(r),
	}, nil
}
