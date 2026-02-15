package command

import (
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	for _, args := range [][]string{
		{"help"},
		{"-help"},
		{"--help"},
		{"-h"},
		nil,
	} {
		stdout.Reset()
		if err := Run(args, deps); err != nil {
			t.Fatalf("Run(%v) error: %v", args, err)
		}
		if !strings.Contains(stdout.String(), "lcli - LinkedIn CLI tool") {
			t.Errorf("Run(%v) missing usage header", args)
		}
	}
}

func TestRunUnknownCommand(t *testing.T) {
	deps, _, _ := testDeps()

	err := Run([]string{"nonexistent"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("error = %q, want 'unknown command'", err)
	}
}

func TestRunDispatchesVersion(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := Run([]string{"version"}, deps); err != nil {
		t.Fatalf("Run version: %v", err)
	}
	if !strings.Contains(stdout.String(), "lcli") {
		t.Error("version output missing 'lcli'")
	}
}

func TestRunDispatchesCompletion(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := Run([]string{"completion", "bash"}, deps); err != nil {
		t.Fatalf("Run completion bash: %v", err)
	}
	if !strings.Contains(stdout.String(), "complete -F _lcli lcli") {
		t.Error("bash completion missing expected marker")
	}
}
