package command

import (
	"strings"
	"testing"
)

func TestCompletionBash(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runCompletion([]string{"bash"}, deps); err != nil {
		t.Fatalf("completion bash: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "complete -F _lcli lcli") {
		t.Error("bash completion missing 'complete -F _lcli lcli' marker")
	}
	if !strings.Contains(out, "_lcli()") {
		t.Error("bash completion missing function definition")
	}
}

func TestCompletionZsh(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runCompletion([]string{"zsh"}, deps); err != nil {
		t.Fatalf("completion zsh: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "#compdef lcli") {
		t.Error("zsh completion missing '#compdef lcli' marker")
	}
	if !strings.Contains(out, "_lcli") {
		t.Error("zsh completion missing function definition")
	}
}

func TestCompletionUnsupportedShell(t *testing.T) {
	deps, _, _ := testDeps()

	err := runCompletion([]string{"fish"}, deps)
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("error = %q", err)
	}
}

func TestCompletionHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runCompletion(nil, deps); err != nil {
		t.Fatalf("completion help: %v", err)
	}
	if !strings.Contains(stdout.String(), "bash") {
		t.Error("help output missing 'bash'")
	}
}
