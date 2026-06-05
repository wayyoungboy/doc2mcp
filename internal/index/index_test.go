package index_test

import (
	"path/filepath"
	"testing"

	"github.com/wayyoungboy/doc2mcp/internal/compiler"
	"github.com/wayyoungboy/doc2mcp/internal/index"
)

func TestSearchReturnsCitedSections(t *testing.T) {
	out := filepath.Join(t.TempDir(), "demo-docs")
	if err := compiler.Build(compiler.BuildOptions{SourceDir: filepath.Join("..", "..", "testdata", "docs"), OutDir: out, Name: "demo-docs"}); err != nil {
		t.Fatal(err)
	}

	pkg, err := index.Load(out)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	results := index.Search(pkg, "authentication token scope", 3)

	if len(results) == 0 {
		t.Fatal("expected search results")
	}
	if results[0].SectionID != "api.md#authentication" {
		t.Fatalf("expected authentication section first, got %#v", results[0])
	}
	if results[0].Citation.SourcePath != "api.md" {
		t.Fatalf("missing citation source: %#v", results[0])
	}
}

func TestReadSectionByID(t *testing.T) {
	out := filepath.Join(t.TempDir(), "demo-docs")
	if err := compiler.Build(compiler.BuildOptions{SourceDir: filepath.Join("..", "..", "testdata", "docs"), OutDir: out, Name: "demo-docs"}); err != nil {
		t.Fatal(err)
	}
	pkg, err := index.Load(out)
	if err != nil {
		t.Fatal(err)
	}

	section, ok := index.SectionByID(pkg, "api.md#rate-limits")
	if !ok {
		t.Fatal("expected section")
	}
	if section.Title != "Rate Limits" {
		t.Fatalf("unexpected section: %#v", section)
	}
}
