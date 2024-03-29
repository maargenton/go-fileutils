//go:build !windows
// +build !windows

package popen

import (
	"context"
	"os/exec"
	"syscall"
	"time"
)

func (c *Command) configureCommand(cmd *exec.Cmd) {
	if !c.NoProcessGroup {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
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

	if c.ShutdownGracePeriod != 0 {
		if c.ShutdownSignal == 0 {
			c.ShutdownSignal = syscall.SIGINT
		}
		c.kill(cmd, c.ShutdownSignal)

		select {
		case <-waitDone:
			return waitError
		case <-time.After(c.ShutdownGracePeriod):
		}
	}

	cmd.Process.Kill()

	// Kill process after potential grace period; ignore error -- process
	// already exited
	c.kill(cmd, syscall.SIGKILL)

	<-waitDone
	return waitError
}

func (c *Command) kill(cmd *exec.Cmd, signal syscall.Signal) error {
	if c.NoProcessGroup {
		return cmd.Process.Signal(signal)
	}
	return syscall.Kill(-cmd.Process.Pid, signal)
}
