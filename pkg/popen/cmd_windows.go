//go:build windows
// +build windows

package popen

import (
	"context"
	"os/exec"
	"syscall"
	"time"
)

func (c *Command) configureCommand(cmd *exec.Cmd) {
	// NoProcessGroup options is not supported on windows
}

func (c *Command) wait(cmd *exec.Cmd, ctx context.Context) error {
	var waitError error
	var waitDone = make(chan struct{})

	go func() {
		waitError = cmd.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		return waitError
	case <-ctx.Done():
	}

	cmd.Process.Kill()

	<-waitDone
	return waitError
}
