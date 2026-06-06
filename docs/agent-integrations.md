# Agent Integrations

Doc2MCP serves compiled documentation over stdio MCP, so it works well with agents that can launch local MCP servers.

## Prepare a Package

Build the documentation package first:

```bash
doc2mcp build ./docs --out ./dist/product-docs --name product-docs
```

Use an absolute package path in agent config:

```bash
realpath ./dist/product-docs
```

## Claude Code

Claude Code supports project-scoped MCP servers through a repository `.mcp.json`. Add a server from the project root:

```bash
claude mcp add product-docs --scope project -- doc2mcp serve /absolute/path/to/dist/product-docs
```

Or write the config directly:

```json
{
  "mcpServers": {
    "product-docs": {
      "type": "stdio",
      "command": "doc2mcp",
      "args": ["serve", "/absolute/path/to/dist/product-docs"]
    }
  }
}
```

Claude Code asks for trust approval before using project-scoped MCP servers. Keep generated packages narrow and rebuild them when source docs change.

## Codex

Codex MCP launchers are configured from the Codex config layer. Add the Doc2MCP server to `~/.codex/config.toml`:

```toml
[mcp_servers.product-docs]
command = "doc2mcp"
args = ["serve", "/absolute/path/to/dist/product-docs"]
```

Then start Codex in the repository and ask it to use the `product-docs` MCP tools:

```text
Use product-docs to answer with citations from the package.
```

## Recommended Workflow

1. Convert non-Markdown sources upstream if needed.
2. Run `doc2mcp build`.
3. Inspect `doc2mcp.json`, `index.json`, and `sources/`.
4. Register the package with Claude Code or Codex.
5. Rebuild and restart the MCP session after documentation changes.
