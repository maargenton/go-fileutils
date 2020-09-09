package popen

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Command is an additional layer of abstraction over exec.Command aimed at
// simplifying common uses where the output of the process is captured or
// redirected. Unlike exec.Command, all the details of the command to run and
// what to do with its outputs are captured in public fields of the `Command`
// structure. The output streams, stdout and stderr, can be returned as a
// string, redirected to a file or stream-processed through an `io.Reader`.
//
// Except for `StdoutReader` and `StderrReader` which are most likely stateful,
// the command object is stateless and can potentially be `Run()` multiple
// times, concurrently.
type Command struct {
	// Directory if specified will be used as the working directory while
	// running the command, instead of the current working directory.
	Directory string

	// Command is the either the name of the command to execute, or its full
	// path.
	Command string

	// Arguments is the list of arguments to pass to the command. They are
	// passed as is, without any shell expansion.
	Arguments []string

	// Env is an optional list of environment variables formatted as
	// 'ENV=VALUE'. If the same environment is defined multiple time, only the
	// leater definition is retained.
	Env []string

	// OverwriteEnv prevents the target command from inheriting the current
	// process environment
	OverwriteEnv bool

	// Stdin, if not empty, will be written to the stdin input of the command.
	// It is ignored if either `ReadStdinFromFile` or `StdinWriter` are
	// specified.
	// Stdin string

	// ReadStdinFromFile if not empty causes the stdin input of the command to
	// be fed with the cntent of the file specified. If takes precedence over
	// `Stdin`, but is ignored if `StdinWriter` specified.
	// ReadStdinFromFile string

	// StdinWriter is a function expected to write the content of the command
	// stdin input to its argument w. Once it retruns, the command stdin is
	// closed. The command is aborted if an error is returned. StdinWriter takes
	// precedence of both `Stdin` and `ReadStdinFromFile`.
	// StdinWriter func(w io.Writer) error

	// DiscardStdout causes the stdout stream of the child process to not be
	// captured or returned as stdout.
	DiscardStdout bool

	// WriteStdoutToFile if not empty causes the stdout stream of the command to
	// be written to the specified file, instead of being captured into Stdout.
	WriteStdoutToFile string

	// StdoutReader if specified is expected to read the content of the command
	// stdout stream from w. If it returns an error, the command is aborted.
	StdoutReader func(r io.Reader) error

	// DiscardStderr causes the stderr stream of the child process to not be
	// captured or returned as stderr.
	DiscardStderr bool

	// WriteStderrToFile if not empty causes the stderr stream of the command to
	// be written to the specified file, instead of being captured into Stderr.
	WriteStderrToFile string

	// StderrReader if specified is expected to read the content of the command
	// stderr stream from w. If it returns an error, the command is aborted.
	StderrReader func(r io.Reader) error
}

// Run executes the command as specified and returns the captured content of
// stdout and stderr if not discarded. If the process is executed successfully
// but returns a non-zero exit status, the returned error is an exec.ExitError
// that contains the actual status code.
func (c *Command) Run(ctx context.Context) (stdout, stderr string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var cmd = exec.CommandContext(ctx, c.Command, c.Arguments...)
	var closeAfterWait []io.Closer
	var servicers []func()
	var servicerErrors chan error

	// Setup environment
	if c.Directory != "" {
		cmd.Dir = c.Directory
	}

	if len(c.Env) > 0 {
		var env []string
		if !c.OverwriteEnv {
			env = os.Environ()
		}
		cmd.Env = append(env, c.Env...)
	}

	// Configure stdout
	var stdoutStreams []io.Writer
	if c.WriteStdoutToFile != "" {
		f, err := os.OpenFile(
			c.WriteStdoutToFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return "", "", fmt.Errorf(
				"failed to open stdout file '%v': %w",
				c.WriteStdoutToFile, err)
		}
		defer f.Close()
		stdoutStreams = append(stdoutStreams, f)
	}

	if c.StdoutReader != nil {
		r, w := io.Pipe()
		stdoutStreams = append(stdoutStreams, w)
		closeAfterWait = append(closeAfterWait, w)

		handler := c.StdoutReader
		servicers = append(servicers, func() {
			err := handler(r)
			if err != nil && err != io.EOF {
				cancel()
			}
			drain(r)
			servicerErrors <- err
		})
	}

	var stdoutBuf strings.Builder
	if !c.DiscardStdout && c.WriteStdoutToFile == "" {
		stdoutStreams = append(stdoutStreams, &stdoutBuf)
	}

	if len(stdoutStreams) == 1 {
		cmd.Stdout = stdoutStreams[0]
	} else if len(stdoutStreams) > 1 {
		cmd.Stdout = io.MultiWriter(stdoutStreams...)
	}

	// Configure stderr
	var stderrStreams []io.Writer
	if c.WriteStderrToFile != "" {
		f, err := os.OpenFile(
			c.WriteStderrToFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return "", "", fmt.Errorf(
				"failed to open stderr file '%v': %w",
				c.WriteStderrToFile, err)
		}
		defer f.Close()
		stderrStreams = append(stderrStreams, f)
	}

	if c.StderrReader != nil {
		r, w := io.Pipe()
		stderrStreams = append(stderrStreams, w)
		closeAfterWait = append(closeAfterWait, w)

		handler := c.StderrReader
		servicers = append(servicers, func() {
			err := handler(r)
			if err != nil && err != io.EOF {
				cancel()
			}
			drain(r)
			servicerErrors <- err
		})
	}

	var stderrBuf strings.Builder
	if !c.DiscardStderr && c.WriteStderrToFile == "" {
		stderrStreams = append(stderrStreams, &stderrBuf)
	}

	if len(stderrStreams) == 1 {
		cmd.Stderr = stderrStreams[0]
	} else if len(stderrStreams) > 1 {
		cmd.Stderr = io.MultiWriter(stderrStreams...)
	}

	// Configure stdin

	// Start the sub-process
	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf(
			"failed to start command '%v': %w",
			c.Command, err)
	}

	if len(servicers) > 0 {
		servicerErrors = make(chan error, len(servicerErrors))
		for _, servicer := range servicers {
			go servicer()
		}
	}

	err = cmd.Wait()
	for _, c := range closeAfterWait {
		c.Close()
	}

	// Wait for all servicers to complete and capture first error
	var servicerError error
	for range servicers {
		if err := <-servicerErrors; err != nil && servicerError == nil {
			servicerError = err
		}
	}

	if ctx.Err() != nil {
		err = ctx.Err()
	}
	if err == nil || err == context.Canceled {
		if servicerError != nil {
			err = servicerError
		}
	}

	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()
	return
}

func drain(r io.Reader) {
	var p = make([]byte, 4096)
	for {
		n, err := r.Read(p)
		if n == 0 || err != nil {
			break
		}
	}
}
