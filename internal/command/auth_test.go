package command

import (
	"strings"
	"testing"
)

func TestAuthDispatchHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runAuth(nil, deps); err != nil {
		t.Fatalf("runAuth help: %v", err)
	}
	if !strings.Contains(stdout.String(), "login") {
		t.Error("help output missing 'login' subcommand")
	}
}

func TestAuthDispatchUnknown(t *testing.T) {
	deps, _, _ := testDeps()

	err := runAuth([]string{"unknown"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}

func TestAuthDispatchHelpFlag(t *testing.T) {
	deps, stdout, _ := testDeps()

	for _, flag := range []string{"-help", "--help", "-h"} {
		stdout.Reset()
		if err := runAuth([]string{flag}, deps); err != nil {
			t.Fatalf("runAuth %s: %v", flag, err)
		}
		if !strings.Contains(stdout.String(), "login") {
			t.Errorf("help output for %s missing 'login'", flag)
		}
	}
}
