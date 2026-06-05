# Doc2MCP Threat Model

## Assets

- Source documentation
- Compiled package manifest
- Section checksums and citations
- MCP tools, resources, and prompts exposed to agents

## In Scope

- Accidental broad filesystem exposure
- Unsourced answers caused by weak retrieval boundaries
- Package drift without reproducible build artifacts
- Losing provenance during document chunking

## Out Of Scope For v0.1

- Native PDF/DOCX parser security
- Remote URL crawling
- Sandboxed document conversion
- Hosted registry integrity
- Vector database poisoning detection

## Security Posture

Doc2MCP is a build-time compiler. Agents should connect to the compiled package, not a live crawler over your whole filesystem. The MCP server exposes `search_docs`, `read_doc`, and `cite_source` over known package sections.
