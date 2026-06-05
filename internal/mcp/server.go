package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/wayyoungboy/doc2mcp/internal/index"
	"github.com/wayyoungboy/doc2mcp/internal/model"
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Server struct {
	pkg model.Package
}

func NewServer(pkg model.Package) *Server {
	return &Server{pkg: pkg}
}

func (s *Server) Handle(req Request) Response {
	switch req.Method {
	case "initialize":
		return ok(req.ID, map[string]interface{}{
			"protocolVersion": "2025-06-18",
			"serverInfo":      map[string]string{"name": "doc2mcp", "version": "0.1.0"},
			"capabilities": map[string]interface{}{
				"tools":     map[string]interface{}{},
				"resources": map[string]interface{}{},
				"prompts":   map[string]interface{}{},
			},
		})
	case "tools/list":
		return ok(req.ID, map[string]interface{}{"tools": []map[string]interface{}{
			{"name": "search_docs", "description": "Search compiled documentation and return cited sections."},
			{"name": "read_doc", "description": "Read a compiled documentation section by id."},
			{"name": "cite_source", "description": "Return citation metadata for a section id."},
		}})
	case "tools/call":
		return s.callTool(req)
	case "resources/list":
		return ok(req.ID, map[string]interface{}{"resources": s.resources()})
	case "resources/read":
		return s.readResource(req)
	case "prompts/list":
		return ok(req.ID, map[string]interface{}{"prompts": []map[string]string{{"name": "answer_with_citations", "description": "Answer using only cited Doc2MCP sections."}}})
	case "prompts/get":
		return ok(req.ID, map[string]interface{}{"description": "Answer using only cited Doc2MCP sections.", "messages": []map[string]interface{}{{"role": "user", "content": map[string]string{"type": "text", "text": "Use search_docs first. Answer only from cited sections and include section IDs."}}}})
	default:
		return fail(req.ID, -32601, "method not found")
	}
}

func (s *Server) callTool(req Request) Response {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return fail(req.ID, -32602, err.Error())
	}
	switch params.Name {
	case "search_docs":
		query, _ := params.Arguments["query"].(string)
		limit := 5
		if raw, ok := params.Arguments["limit"].(float64); ok && raw > 0 {
			limit = int(raw)
		}
		return ok(req.ID, map[string]interface{}{"content": textContent(index.Search(s.pkg, query, limit))})
	case "read_doc":
		id, _ := params.Arguments["section_id"].(string)
		section, found := index.SectionByID(s.pkg, id)
		if !found {
			return fail(req.ID, -32004, "section not found")
		}
		return ok(req.ID, map[string]interface{}{"content": []map[string]string{{"type": "text", "text": section.Text}}})
	case "cite_source":
		id, _ := params.Arguments["section_id"].(string)
		section, found := index.SectionByID(s.pkg, id)
		if !found {
			return fail(req.ID, -32004, "section not found")
		}
		return ok(req.ID, map[string]interface{}{"content": []map[string]string{{"type": "text", "text": fmt.Sprintf("%s (%s, checksum %s)", section.ID, section.SourcePath, section.Checksum)}}})
	default:
		return fail(req.ID, -32602, "unknown tool")
	}
}

func (s *Server) resources() []map[string]string {
	var resources []map[string]string
	for _, doc := range s.pkg.Documents {
		resources = append(resources, map[string]string{
			"uri":         "doc2mcp://" + doc.SourcePath,
			"name":        doc.Title,
			"description": doc.SourcePath,
		})
	}
	return resources
}

func (s *Server) readResource(req Request) Response {
	var params struct {
		URI string `json:"uri"`
	}
	_ = json.Unmarshal(req.Params, &params)
	for _, doc := range s.pkg.Documents {
		if params.URI == "doc2mcp://"+doc.SourcePath {
			var text string
			for _, section := range doc.Sections {
				text += "# " + section.Title + "\n\n" + section.Text + "\n\n"
			}
			return ok(req.ID, map[string]interface{}{"contents": []map[string]string{{"uri": params.URI, "mimeType": "text/markdown", "text": text}}})
		}
	}
	return fail(req.ID, -32004, "resource not found")
}

func textContent(results []model.SearchResult) []map[string]string {
	data, _ := json.MarshalIndent(results, "", "  ")
	return []map[string]string{{"type": "text", "text": string(data)}}
}

func ok(id interface{}, result interface{}) Response {
	return Response{JSONRPC: "2.0", ID: id, Result: result}
}

func fail(id interface{}, code int, message string) Response {
	return Response{JSONRPC: "2.0", ID: id, Error: &Error{Code: code, Message: message}}
}

func Serve(r io.Reader, w io.Writer, server *Server) error {
	scanner := bufio.NewScanner(r)
	encoder := json.NewEncoder(w)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			if err := encoder.Encode(fail(nil, -32700, err.Error())); err != nil {
				return err
			}
			continue
		}
		if err := encoder.Encode(server.Handle(req)); err != nil {
			return err
		}
	}
	return scanner.Err()
}
