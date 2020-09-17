package unarchiver

import (
	"archive/tar"
	"io"
	"os"
	"path"
	"time"
)

type tarUnarchiver struct {
	r *tar.Reader
}

func newTarUnarchiver(r io.Reader) *tarUnarchiver {
	return &tarUnarchiver{
		r: tar.NewReader(r),
	}
}

func (u *tarUnarchiver) NextItem() (item ArchiveItem, err error) {
	hdr, err := u.r.Next()
	if err != nil {
		return nil, err
	}
	return tarItem{
		hdr:     hdr,
		content: u.r,
	}, nil
}

// ---------------------------------------------------------------------------

type tarItem struct {
	hdr     *tar.Header
	content io.Reader
}

func (i tarItem) Name() string       { return path.Base(i.hdr.Name) }
func (i tarItem) Size() int64        { return i.hdr.Size }
func (i tarItem) Mode() os.FileMode  { return i.modePerm() + i.modeType() }
func (i tarItem) ModTime() time.Time { return i.hdr.ModTime }
func (i tarItem) IsDir() bool        { return i.hdr.Typeflag == tar.TypeDir }
func (i tarItem) Sys() interface{}   { return nil }

func (i tarItem) FullName() string   { return i.hdr.Name }
func (i tarItem) LinkName() string   { return i.hdr.Linkname }
func (i tarItem) IsHardLink() bool   { return i.hdr.Typeflag == tar.TypeLink }
func (i tarItem) Content() io.Reader { return i.content }

func (i tarItem) modePerm() os.FileMode {
	return os.FileMode(i.hdr.Mode) & os.ModePerm
}

func (i tarItem) modeType() os.FileMode {
	switch i.hdr.Typeflag {
	case tar.TypeReg, tar.TypeRegA:
		return 0
	case tar.TypeLink:
		return os.ModeIrregular
	case tar.TypeSymlink:
		return os.ModeSymlink
	case tar.TypeChar:
		return os.ModeCharDevice
	case tar.TypeBlock:
		return os.ModeDevice
	case tar.TypeDir:
		return os.ModeDir
	case tar.TypeFifo:
		return os.ModeNamedPipe
	default:
		return os.ModeIrregular
	}
}
