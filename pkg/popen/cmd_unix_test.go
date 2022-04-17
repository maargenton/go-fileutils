//go:build !windows
// +build !windows

package popen_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-fileutils/pkg/popen"
)

// ---------------------------------------------------------------------------
// ShutdownGracePeriod -- unix only

func TestCommandTimeout(t *testing.T) {
	var cmd = popen.Command{
		Command: "sleep",
		Arguments: []string{
			"3",
		},
		ShutdownGracePeriod: 200 * time.Millisecond,
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
		ShutdownGracePeriod: 200 * time.Millisecond,
		ShutdownSignal:      syscall.SIGTERM,
		NoProcessGroup:      true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, _, err := cmd.Run(ctx)

	verify.That(t, err).ToString().Eq("signal: terminated")
}

func TestCommandShutdownGracePeriod(t *testing.T) {
	// Pre-built test child process to avoid long timeouts
	var build = popen.Command{
		Command: "go",
		Arguments: []string{
			"run",
			"./test-child-process",
		},
	}
	build.Run(nil)

	bdd.Given(t, "a command that takes time to shutdown", func(t *bdd.T) {
		var cmd = popen.Command{
			Command: "go",
			Arguments: []string{
				"run",
				"./test-child-process",
				"200ms",
			},
		}

		t.When("setting grace period longer than shutdown time", func(t *bdd.T) {
			cmd.ShutdownGracePeriod = 300 * time.Millisecond

			t.Then("the command should exit on its own", func(t *bdd.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				defer cancel()
				stdout, _, err := cmd.Run(ctx)
				verify.That(t, err).IsNotNil()
				verify.That(t, stdout).EndsWith("exiting\n")
			})
		})
		t.When("setting grace period shorter than shutdown time", func(t *bdd.T) {
			cmd.ShutdownGracePeriod = 100 * time.Millisecond

			t.Then("the command should be killed while shutting down", func(t *bdd.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				defer cancel()
				stdout, _, err := cmd.Run(ctx)
				verify.That(t, err).IsNotNil()
				verify.That(t, stdout).EndsWith("shuting down ...\n")
			})
		})
	})
}

// ShutdownGracePeriod -- unix only
// ---------------------------------------------------------------------------
