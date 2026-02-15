package command

import (
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	deps, stdout, _ := testDeps()

	BuildVersion = "1.2.3"
	BuildCommit = "abc1234"
	BuildDate = "2025-01-01T00:00:00Z"
	defer func() {
		BuildVersion = "dev"
		BuildCommit = "unknown"
		BuildDate = "unknown"
	}()

	if err := runVersion(deps); err != nil {
		t.Fatalf("runVersion: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "1.2.3") {
		t.Errorf("output missing version: %s", out)
	}
	if !strings.Contains(out, "abc1234") {
		t.Errorf("output missing commit: %s", out)
	}
	if !strings.Contains(out, "2025-01-01T00:00:00Z") {
		t.Errorf("output missing date: %s", out)
	}
}
