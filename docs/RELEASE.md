# Release Notes

## v0.1.0

Initial open-source MVP.

### Included

- `build` command for Markdown/text directories
- `search` command for cited retrieval
- `show` command for source section display
- `serve` command for stdio MCP
- Tools: `search_docs`, `read_doc`, `cite_source`
- Resources and prompt listing
- Package manifest and local index
- GitHub Actions CI

### Verification

```bash
go test ./...
go build -o doc2mcp ./cmd/doc2mcp
./doc2mcp build testdata/docs --out dist/demo-docs --name demo-docs
./doc2mcp search dist/demo-docs "authentication"
./doc2mcp search dist/demo-docs "authentication" --json
./doc2mcp show dist/demo-docs api.md#authentication
```
