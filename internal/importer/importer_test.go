package importer_test

import (
	"path/filepath"
	"testing"

	"github.com/wayyoungboy/doc2mcp/internal/importer"
)

func TestImportDirectoryExtractsMarkdownSectionsWithProvenance(t *testing.T) {
	docs, err := importer.ImportDirectory(filepath.Join("..", "..", "testdata", "docs"))
	if err != nil {
		t.Fatalf("ImportDirectory returned error: %v", err)
	}

	if len(docs) != 2 {
		t.Fatalf("expected 2 docs, got %d", len(docs))
	}

	api := docs[0]
	if api.SourcePath != "api.md" {
		t.Fatalf("expected api.md first, got %q", api.SourcePath)
	}
	if len(api.Sections) != 3 {
		t.Fatalf("expected title plus 2 sections, got %d", len(api.Sections))
	}
	auth := api.Sections[1]
	if auth.ID != "api.md#authentication" {
		t.Fatalf("unexpected section id %q", auth.ID)
	}
	if auth.Title != "Authentication" {
		t.Fatalf("unexpected title %q", auth.Title)
	}
	if auth.Checksum == "" || auth.StartByte < 0 || auth.EndByte <= auth.StartByte {
		t.Fatalf("missing provenance: %#v", auth)
	}
}
