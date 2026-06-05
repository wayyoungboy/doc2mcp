package compiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wayyoungboy/doc2mcp/internal/compiler"
)

func TestBuildWritesPortableKnowledgePackage(t *testing.T) {
	out := filepath.Join(t.TempDir(), "demo-docs")
	err := compiler.Build(compiler.BuildOptions{
		SourceDir: filepath.Join("..", "..", "testdata", "docs"),
		OutDir:    out,
		Name:      "demo-docs",
	})
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	for _, name := range []string{"doc2mcp.json", "index.json", filepath.Join("sources", "api.md")} {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Fatalf("expected package file %s: %v", name, err)
		}
	}
}
