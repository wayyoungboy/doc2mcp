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
go run ./cmd/doc2mcp build testdata/docs --out /tmp/doc2mcp-demo --name demo-docs
go run ./cmd/doc2mcp search /tmp/doc2mcp-demo "authentication"
go run ./cmd/doc2mcp show /tmp/doc2mcp-demo api.md#authentication
```
