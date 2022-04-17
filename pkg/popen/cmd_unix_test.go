//go:build !windows
// +build !windows

package popen_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-fileutils/pkg/popen"
)

// ---------------------------------------------------------------------------
// KillGracePeriod -- unix only

func TestCommandTimeout(t *testing.T) {
	var cmd = popen.Command{
		Command: "sleep",
		Arguments: []string{
			"3",
		},
		KillGracePeriod: 200 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, _, err := cmd.Run(ctx)

	verify.That(t, err).ToString().Eq("signal: interrupt")
}

func TestCommandTimeoutSigint(t *testing.T) {
	var cmd = popen.Command{
		Command: "sleep",
		Arguments: []string{
			"3",
		},
		KillGracePeriod: 200 * time.Millisecond,
		PreKillSignal:   syscall.SIGTERM,
		NoProcessGroup:  true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, _, err := cmd.Run(ctx)

	verify.That(t, err).ToString().Eq("signal: terminated")
}

func TestCommandTimeoutTrapSigtermWait(t *testing.T) {
	var cmd = popen.Command{
		Command: "go",
		Arguments: []string{
			"run",
			"./test-child-process",
		},
		KillGracePeriod: 300 * time.Millisecond,
		PreKillSignal:   syscall.SIGTERM,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	stdout, _, err := cmd.Run(ctx)
	verify.That(t, err).IsNotNil()

	// child process takes 200ms after SIGINT / SIGTERM to shutdown; we wait
	// longer that that is the grace period, so the process should exit on its
	// own.
	verify.That(t, stdout).EndsWith("exiting\n")
}

func TestCommandTimeoutTrapSigtermKill(t *testing.T) {
	var cmd = popen.Command{
		Command: "go",
		Arguments: []string{
			"run",
			"./test-child-process",
		},
		KillGracePeriod: 100 * time.Millisecond,
		PreKillSignal:   syscall.SIGTERM,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	stdout, _, err := cmd.Run(ctx)
	verify.That(t, err).IsNotNil()

	// child process takes 200ms after SIGINT / SIGTERM to shutdown; we wait
	// only 100ms before sending SIGKILL, so the process should be killed while
	// shutting down.
	verify.That(t, stdout).EndsWith("shuting down ...\n")
}

// Context done
// ---------------------------------------------------------------------------
