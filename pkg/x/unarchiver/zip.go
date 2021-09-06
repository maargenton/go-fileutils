package unarchiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"
)

// zipFile define the minimal interface for an object to be used as the source
// of a zip unarchiver.
type zipFile interface {
	ReadAt(b []byte, off int64) (n int, err error)
	Stat() (os.FileInfo, error)
}

type zipUnarchiver struct {
	r     *zip.Reader
	index int

	currentFile     *zip.File
	currentFileInfo os.FileInfo
	currentLinkName string
	currentContent  io.ReadCloser
}

func newZipUnarchiver(r io.Reader) (*zipUnarchiver, error) {
	zf, ok := r.(zipFile)
	if !ok {
		return nil, fmt.Errorf("cannot decode ZIP archive from '%T'", r)
	}
	info, err := zf.Stat()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to determine size of ZIP archive file, %w", err)
	}
	zr, err := zip.NewReader(zf, info.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to open ZIP archive, %w", err)
	}

	return &zipUnarchiver{
		r: zr,
	}, nil
}

func (u *zipUnarchiver) NextItem() (item ArchiveItem, err error) {
	if u.currentContent != nil {
		u.currentContent.Close()
		u.currentContent = nil
	}

	if u.index >= len(u.r.File) {
		return nil, io.EOF
	}

	u.currentFile = u.r.File[u.index]
	u.currentFileInfo = u.currentFile.FileInfo()

	if u.currentFileInfo.Mode()&os.ModeSymlink != 0 {
		r, err := u.currentFile.Open()
		if err != nil {
			return zipItem{}, err
		}
		var buf strings.Builder
		_, err = io.Copy(&buf, r)
		if err != nil {
			return zipItem{}, err
		}

		err = r.Close()
		if err != nil {
			return zipItem{}, err
		}
		u.currentLinkName = buf.String()
	} else {
		u.currentLinkName = ""
	}

	item = zipItem{
		FileInfo: u.currentFileInfo,
		u:        u,
	}

	u.index++
	return
}

// ---------------------------------------------------------------------------

type zipItem struct {
	os.FileInfo
	u *zipUnarchiver
}

var _ ArchiveItem = zipItem{}

func (i zipItem) FullName() string { return i.u.currentFile.Name }
func (i zipItem) LinkName() string { return i.u.currentLinkName }
func (i zipItem) IsHardLink() bool { return false }

func (i zipItem) Content() io.Reader {
	if i.u.currentContent != nil {
		return i.u.currentContent
	}

	i.u.currentContent, _ = i.u.currentFile.Open()
	return nil
}
