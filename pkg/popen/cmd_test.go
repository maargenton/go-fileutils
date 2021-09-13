package popen_test

import (
	"bufio"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-fileutils/pkg/popen"
)

// ---------------------------------------------------------------------------
// Directory and environment

func TestCommandDirectory(t *testing.T) {
	var tmp = tempDir(t)
	var dir = fileutils.Join(tmp, "a", "b", "c")
	os.MkdirAll(dir, 0777)

	var cmd = popen.Command{
		Directory: dir,
		Command:   "pwd",
	}
	if runtime.GOOS == "windows" {
		cmd.Command = "echo $PWD"
	}

	stdout, _, err := cmd.Run(context.Background())
	stdout = strings.TrimSpace(stdout)
	path := fileutils.Clean(stdout)
	verify.That(t, err).IsNil()
	verify.That(t, path).Eq(dir)

}

func TestCommandEnv(t *testing.T) {
	var cmd = popen.Command{
		Command: "env",
		Env: []string{
			"aaa=bbb",
			"ccc=bbb",
		},
	}
	stdout, _, err := cmd.Run(context.Background())
	verify.That(t, err).IsNil()
	var env = parseEnv(stdout)

	verify.That(t, env).MapKeys().IsSupersetOf([]string{"PATH", "aaa", "ccc"})
}

func TestCommandEnvOverride(t *testing.T) {

	var cmd = popen.Command{
		Command: "env",
		Env: []string{
			"aaa=bbb",
			"ccc=bbb",
		},
		OverwriteEnv: true,
	}
	stdout, _, err := cmd.Run(context.Background())
	verify.That(t, err).IsNil()
	var env = parseEnv(stdout)

	verify.That(t, env).MapKeys().IsSupersetOf([]string{"aaa", "ccc"})
	verify.That(t, env).MapKeys().IsDisjointSetFrom([]string{"PATH"})
}

// Directory and environment
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Stdout

func TestCommandStdoutReader(t *testing.T) {

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
	verify.That(t, err).IsNil()
	verify.That(t, buf.String()).Eq(stdout)
}

func TestCommandStdoutReaderErrorAbortsCommand(t *testing.T) {

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
	verify.That(t, err).IsError(expectedError)
}

func TestCommandStdoutToFile(t *testing.T) {

	var tmp = tempDir(t)

	var stdoutFile = fileutils.Join(tmp, "stdout.txt")
	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World; done",
		},
		WriteStdoutToFile: stdoutFile,
	}

	stdout, _, err := cmd.Run(context.Background())
	verify.That(t, err).IsNil()
	verify.That(t, stdout).Eq("") // Discarded when written to file

	var linecnt = 0
	fileutils.ReadFile(stdoutFile, func(r io.Reader) error {
		var scanner = bufio.NewScanner(r)
		for scanner.Scan() {
			linecnt++
		}
		return scanner.Err()
	})
	// content, _ := ioutil.ReadFile(stdoutFile)
	// fmt.Println(string(content))
	verify.That(t, linecnt).Eq(10)
}

func TestCommandStdoutToInvalidPathFile(t *testing.T) {

	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World; done",
		},
		WriteStdoutToFile: "__invalid_path__/stdout.txt",
	}

	stdout, _, err := cmd.Run(context.Background())
	verify.That(t, err).IsNotNil()
	verify.That(t, os.IsNotExist(errors.Unwrap(err))).IsTrue()
	verify.That(t, stdout).Eq("") // Discarded when written to file
}

// Stdout
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Stderr

func TestCommandStderrReader(t *testing.T) {

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
	verify.That(t, err).IsNil()
	verify.That(t, buf.String()).Eq(stderr)
}

func TestCommandStderrReaderErrorAbortsCommand(t *testing.T) {

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
	verify.That(t, err).IsError(expectedError)
}

func TestCommandStderrToFile(t *testing.T) {

	var tmp = tempDir(t)

	var stderrFile = fileutils.Join(tmp, "stderr.txt")
	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World 1>&2; done",
		},
		WriteStderrToFile: stderrFile,
	}

	_, stderr, err := cmd.Run(context.Background())
	verify.That(t, err).IsNil()
	verify.That(t, stderr).Eq("") // Discarded when written to file

	var linecnt = 0
	fileutils.ReadFile(stderrFile, func(r io.Reader) error {
		var scanner = bufio.NewScanner(r)
		for scanner.Scan() {
			linecnt++
		}
		return scanner.Err()
	})
	verify.That(t, linecnt).Eq(10)
}

func TestCommandStderrToInvalidPathFile(t *testing.T) {

	var cmd = popen.Command{
		Command: "bash",
		Arguments: []string{
			"-c",
			"for i in {1..10}; do echo Hello World 1> &2; done",
		},
		WriteStderrToFile: "__invalid_path__/stderr.txt",
	}

	_, stderr, err := cmd.Run(context.Background())
	verify.That(t, err).IsNotNil()
	verify.That(t, os.IsNotExist(errors.Unwrap(err))).IsTrue()
	verify.That(t, stderr).Eq("") // Discarded when written to file
}

// Stderr
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Invalid command

func TestCommandWithInvalidCommand(t *testing.T) {

	var cmd = popen.Command{
		Command: "__invalid__command__",
	}

	_, _, err := cmd.Run(context.Background())
	verify.That(t, err).IsNotNil()
	verify.That(t, err).ToString().Contains("executable file not found")
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
	tempDir, err = fileutils.Abs(tempDir)
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
