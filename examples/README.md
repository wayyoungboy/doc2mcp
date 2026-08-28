# Examples

MCP client configs must use an **absolute** path to a compiled package. The placeholder `/absolute/path/to/dist/demo-docs` is not a working path until you replace it.

Build the demo package from this repo:

```bash
doc2mcp build testdata/docs --out dist/demo-docs --name demo-docs
realpath dist/demo-docs
```

Then substitute that absolute path in:

- `claude-code.mcp.json` — copy into a project `.mcp.json` for Claude Code
- `codex-config.toml` — merge into `~/.codex/config.toml` for Codex

Section IDs from this demo are like `api.md#authentication`, not `docs/api.md#authentication`.
