package cli_test

import (
	"bytes"
	"encoding/json"
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

func TestHelpListsActualFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("help failed: %s", stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{
		"doc2mcp build <docs-dir> [--out <package-dir>] [--name <name>]",
		"doc2mcp search <package-dir> <query> [--limit <n>] [--json]",
		"doc2mcp show <package-dir> <section-id>",
		"doc2mcp serve <package-dir>",
		"doc2mcp version",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("help missing %q\n%s", want, out)
		}
	}
	if strings.Contains(out, "[--json true]") {
		t.Fatalf("help still documents --json as a required value:\n%s", out)
	}
}

func TestBuildHelpDoesNotError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"build", "--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("build --help failed: stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "--out") || !strings.Contains(stdout.String(), "--name") {
		t.Fatalf("expected build usage, got %s", stdout.String())
	}
}

func TestSearchJSONFlagWithoutValue(t *testing.T) {
	out := filepath.Join(t.TempDir(), "demo-docs")
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"build", filepath.Join("..", "..", "testdata", "docs"), "--out", out, "--name", "demo-docs"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("build failed: %s", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = cli.Run([]string{"search", out, "authentication", "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("search --json failed: %s", stderr.String())
	}
	var payload []map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON output, got %s (%v)", stdout.String(), err)
	}
	if len(payload) == 0 {
		t.Fatal("expected search hits")
	}
}

func TestShowWrongExamplePathFailsClearly(t *testing.T) {
	out := filepath.Join(t.TempDir(), "demo-docs")
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"build", filepath.Join("..", "..", "testdata", "docs"), "--out", out, "--name", "demo-docs"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("build failed: %s", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = cli.Run([]string{"show", out, "docs/api.md#authentication"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected missing section for docs/api.md#authentication")
	}
	if !strings.Contains(stderr.String(), "docs/api.md#authentication") {
		t.Fatalf("expected section id in error, got %s", stderr.String())
	}
}
