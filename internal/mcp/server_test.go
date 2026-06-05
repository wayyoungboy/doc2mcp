package mcp_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/wayyoungboy/doc2mcp/internal/compiler"
	"github.com/wayyoungboy/doc2mcp/internal/index"
	"github.com/wayyoungboy/doc2mcp/internal/mcp"
)

func TestHandleToolsCallSearchDocs(t *testing.T) {
	out := filepath.Join(t.TempDir(), "demo-docs")
	if err := compiler.Build(compiler.BuildOptions{SourceDir: filepath.Join("..", "..", "testdata", "docs"), OutDir: out, Name: "demo-docs"}); err != nil {
		t.Fatal(err)
	}
	pkg, err := index.Load(out)
	if err != nil {
		t.Fatal(err)
	}
	server := mcp.NewServer(pkg)

	resp := server.Handle(mcp.Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"search_docs","arguments":{"query":"authentication token","limit":2}}`),
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %#v", resp.Error)
	}
	data, _ := json.Marshal(resp.Result)
	if !stringsContains(string(data), "api.md#authentication") {
		t.Fatalf("expected search result in response: %s", string(data))
	}
}

func TestHandleListsResourcesAndPrompts(t *testing.T) {
	out := filepath.Join(t.TempDir(), "demo-docs")
	if err := compiler.Build(compiler.BuildOptions{SourceDir: filepath.Join("..", "..", "testdata", "docs"), OutDir: out, Name: "demo-docs"}); err != nil {
		t.Fatal(err)
	}
	pkg, err := index.Load(out)
	if err != nil {
		t.Fatal(err)
	}
	server := mcp.NewServer(pkg)

	resources := server.Handle(mcp.Request{JSONRPC: "2.0", ID: 1, Method: "resources/list"})
	if resources.Error != nil {
		t.Fatalf("resources/list error: %#v", resources.Error)
	}
	prompts := server.Handle(mcp.Request{JSONRPC: "2.0", ID: 2, Method: "prompts/list"})
	if prompts.Error != nil {
		t.Fatalf("prompts/list error: %#v", prompts.Error)
	}
}

func stringsContains(value, part string) bool {
	return len(part) == 0 || (len(value) >= len(part) && (value == part || stringsContains(value[1:], part) || value[:len(part)] == part))
}
