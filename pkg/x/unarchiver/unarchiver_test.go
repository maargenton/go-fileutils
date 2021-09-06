package unarchiver_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/maargenton/fileutil"
	"github.com/maargenton/fileutil/pkg/x/unarchiver"
)

func printItem(w io.Writer, item unarchiver.ArchiveItem) {
	fmt.Fprintf(w, "%v %8d %v", item.Mode(), item.Size(), item.FullName())
	if item.IsHardLink() || item.Mode()&os.ModeSymlink != 0 {
		fmt.Fprintf(w, " -> %v", item.LinkName())
	}
	fmt.Fprintf(w, "\n")
}

func listArchiveContent(filename string) (string, error) {
	var content strings.Builder
	var err = fileutil.ReadFile(filename, func(r io.Reader) error {
		u, err := unarchiver.New(r)
		if err != nil {
			return err
		}
		for {
			item, err := u.NextItem()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return err
			}
			printItem(&content, item)
		}
		return nil
	})

	return content.String(), err
}

func loadGolden(filename string) string {
	var s strings.Builder
	fileutil.ReadFile(filename, func(r io.Reader) error {
		io.Copy(&s, r)
		return nil
	})
	return s.String()
}

func TestTarUnarchiver(t *testing.T) {
	var tcs = []struct {
		archive string
		content string
	}{
		{"testdata/content.tar", "testdata/content.txt"},
		{"testdata/content.tar.gz", "testdata/content.txt"},
		{"testdata/content.tar.bz2", "testdata/content.txt"},
		{"testdata/content2.tar", "testdata/content2.txt"},
		{"testdata/content2.tar.gz", "testdata/content2.txt"},
		{"testdata/content2.tar.bz2", "testdata/content2.txt"},
	}

	for _, tc := range tcs {
		t.Run(tc.archive, func(t *testing.T) {
			assert := asserter.New(t)

			content, err := listArchiveContent(tc.archive)
			assert.That(err, p.IsNoError())

			fileutil.WriteFile(tc.content, func(w io.Writer) error {
				_, err := fmt.Fprint(w, content)
				return err
			})

			var golden = loadGolden(tc.content + ".golden")

			if content != golden {
				diff := difflib.UnifiedDiff{
					A:        difflib.SplitLines(content),
					B:        difflib.SplitLines(golden),
					FromFile: "Actual",
					ToFile:   "Expected",
					Context:  1,
				}
				text, _ := difflib.GetUnifiedDiffString(diff)
				t.Errorf("\nContent does not match expected content:\n%v", text)
			}
		})
	}
}

func TestZipUnarchiver(t *testing.T) {
	assert := asserter.New(t)
	assert.That(nil, p.IsNil())

	var err = fileutil.ReadFile("testdata/content.zip", func(r io.Reader) error {
		u, err := unarchiver.New(r)
		if err != nil {
			return err
		}

		for {
			item, err := u.NextItem()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return err
			}
			printItem(os.Stdout, item)
		}
		return nil
	})

	assert.That(err, p.IsNoError())
	// t.Fail()
}
