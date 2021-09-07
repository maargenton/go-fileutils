package popen_test

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-fileutils/pkg/popen"

	"github.com/maargenton/go-testpredicate/pkg/asserter"
	"github.com/maargenton/go-testpredicate/pkg/p"
)

// ---------------------------------------------------------------------------
// Directory and environment

func TestCommandDirectory(t *testing.T) {
	assert := asserter.New(t)

	var tmp = tempDir(t)
	var dir = filepath.Join(tmp, "a", "b", "c")
	os.MkdirAll(dir, 0777)

	var cmd = popen.Command{
		Directory: dir,
		Command:   "pwd",
	}

	stdout, _, err := cmd.Run(context.Background())
	assert.That(err, p.IsNoError())
	assert.That(stdout, p.StartsWith(dir)) // Ignore trailing LF
}

func TestCommandEnv(t *testing.T) {
	assert := asserter.New(t)

	var cmd = popen.Command{
		Command: "env",
		Env: []string{
			"aaa=bbb",
			"ccc=bbb",
		},
	}
	stdout, _, err := cmd.Run(context.Background())
	assert.That(err, p.IsNoError())
	var env = parseEnv(stdout)

	assert.That(env, p.MapKeys(p.IsSupersetOf([]string{"PATH", "aaa", "ccc"})))
}

func TestCommandEnvOverride(t *testing.T) {
	assert := asserter.New(t)

	var cmd = popen.Command{
		Command: "env",
		Env: []string{
			"aaa=bbb",
			"ccc=bbb",
		},
		OverwriteEnv: true,
	}
	stdout, _, err := cmd.Run(context.Background())
	assert.That(err, p.IsNoError())
	var env = parseEnv(stdout)

	assert.That(env, p.MapKeys(p.IsSupersetOf([]string{"aaa", "ccc"})))
	assert.That(env, p.MapKeys(p.IsDisjointSetFrom([]string{"PATH"})))
}

// Directory and environment
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Stdout

func TestCommandStdoutReader(t *testing.T) {
	assert := asserter.New(t)

	var buf strings.Builder
	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World; done",
		},
		StdoutReader: func(r io.Reader) error {
			_, err := io.Copy(&buf, r)
			return err
		},
	}

	stdout, _, err := cmd.Run(context.Background())
	assert.That(err, p.IsNoError())
	assert.That(buf.String(), p.Eq(stdout))
}

func TestCommandStdoutReaderErrorAbortsCommand(t *testing.T) {
	assert := asserter.New(t)

	var expectedError = errors.New("stdout reader error")
	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"echo Hello;sleep 10;echo Bye",
		},
		StdoutReader: func(r io.Reader) error {
			// // Uncomment should time out test after 5 sec
			// io.Copy(ioutil.Discard, r)
			return expectedError
		},
	}

	_, _, err := cmd.Run(context.Background())
	assert.That(err, p.IsError(expectedError))
}

func TestCommandStdoutToFile(t *testing.T) {
	assert := asserter.New(t)
	var tmp = tempDir(t)

	var stdoutFile = filepath.Join(tmp, "stdout.txt")
	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World; done",
		},
		WriteStdoutToFile: stdoutFile,
	}

	stdout, _, err := cmd.Run(context.Background())
	assert.That(err, p.IsNoError())
	assert.That(stdout, p.Eq("")) // Discarded when written to file

	var linecnt = 0
	fileutils.ReadFile(stdoutFile, func(r io.Reader) error {
		var scanner = bufio.NewScanner(r)
		for scanner.Scan() {
			linecnt++
		}
		return scanner.Err()
	})
	content, _ := ioutil.ReadFile(stdoutFile)
	fmt.Println(string(content))
	assert.That(linecnt, p.Eq(10))
}

func TestCommandStdoutToInvalidPathFile(t *testing.T) {
	assert := asserter.New(t)

	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World; done",
		},
		WriteStdoutToFile: "__invalid_path__/stdout.txt",
	}

	stdout, _, err := cmd.Run(context.Background())
	assert.That(err, p.IsNotNil())
	assert.That(os.IsNotExist(errors.Unwrap(err)), p.IsTrue())
	assert.That(stdout, p.Eq("")) // Discarded when written to file
}

// Stdout
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Stderr

func TestCommandStderrReader(t *testing.T) {
	assert := asserter.New(t)

	var buf strings.Builder
	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World 1>&2; done",
		},
		StderrReader: func(r io.Reader) error {
			_, err := io.Copy(&buf, r)
			return err
		},
	}

	_, stderr, err := cmd.Run(context.Background())
	assert.That(err, p.IsNoError())
	assert.That(buf.String(), p.Eq(stderr))
}

func TestCommandStderrReaderErrorAbortsCommand(t *testing.T) {
	assert := asserter.New(t)

	var expectedError = errors.New("stderr reader error")
	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"echo Hello;sleep 10;echo Bye",
		},
		StderrReader: func(r io.Reader) error {
			// // Uncomment should time out test after 5 sec
			// io.Copy(ioutil.Discard, r)
			return expectedError
		},
	}

	_, _, err := cmd.Run(context.Background())
	assert.That(err, p.IsError(expectedError))
}

func TestCommandStderrToFile(t *testing.T) {
	assert := asserter.New(t)
	var tmp = tempDir(t)

	var stderrFile = filepath.Join(tmp, "stderr.txt")
	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World 1>&2; done",
		},
		WriteStderrToFile: stderrFile,
	}

	_, stderr, err := cmd.Run(context.Background())
	assert.That(err, p.IsNoError())
	assert.That(stderr, p.Eq("")) // Discarded when written to file

	var linecnt = 0
	fileutils.ReadFile(stderrFile, func(r io.Reader) error {
		var scanner = bufio.NewScanner(r)
		for scanner.Scan() {
			linecnt++
		}
		return scanner.Err()
	})
	assert.That(linecnt, p.Eq(10))
}

func TestCommandStderrToInvalidPathFile(t *testing.T) {
	assert := asserter.New(t)

	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World 1> &2; done",
		},
		WriteStderrToFile: "__invalid_path__/stderr.txt",
	}

	_, stderr, err := cmd.Run(context.Background())
	assert.That(err, p.IsNotNil())
	assert.That(os.IsNotExist(errors.Unwrap(err)), p.IsTrue())
	assert.That(stderr, p.Eq("")) // Discarded when written to file
}

// Stderr
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Invalid command

func TestCommandWithInvalidCommand(t *testing.T) {
	assert := asserter.New(t)

	var cmd = popen.Command{
		Command: "__invalid__command__",
	}

	_, _, err := cmd.Run(context.Background())
	assert.That(err, p.IsNotNil())
	assert.That(err, p.ToString(p.Contains("executable file not found")))
}

// Invalid command
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Helpers

func tempDir(t *testing.T) string {
	tempDir, err := ioutil.TempDir(".", "testdata-")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	tempDir, err = filepath.Abs(tempDir)
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	return tempDir
}

func parseEnv(env string) map[string]string {
	var result = make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(env))
	for scanner.Scan() {
		var line = scanner.Text()
		var parts = strings.SplitN(line, "=", 2)
		var key, value = parts[0], parts[1]
		result[key] = value
	}
	return result
}
