package cli_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wayyoungboy/doc2mcp/internal/cli"
)

func TestBuildSearchAndShowCommands(t *testing.T) {
	out := filepath.Join(t.TempDir(), "demo-docs")
	var stdout, stderr bytes.Buffer

	code := cli.Run([]string{"build", filepath.Join("..", "..", "testdata", "docs"), "--out", out, "--name", "demo-docs"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("build failed: stdout=%s stderr=%s", stdout.String(), stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = cli.Run([]string{"search", out, "authentication token"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("search failed: stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "api.md#authentication") {
		t.Fatalf("expected cited section, got %s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = cli.Run([]string{"show", out, "api.md#authentication"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("show failed: stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "Authorization") {
		t.Fatalf("expected section text, got %s", stdout.String())
	}
}

func TestVersionCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("version failed: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Doc2MCP") {
		t.Fatalf("unexpected version output %q", stdout.String())
	}
}
