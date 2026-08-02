package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hektyc/unraid-mcp-server/internal/agentcontent"
	"github.com/hektyc/unraid-mcp-server/internal/client"
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/logger"
)

const protocolVersion = "2024-11-05"

// ServerVersion is injected at build time via
// -ldflags "-X github.com/hektyc/unraid-mcp-server/internal/mcp.ServerVersion=$VERSION"
var ServerVersion = "dev"

type ToolHandler func(ctx context.Context, arguments map[string]interface{}) (string, error)

type ToolDef struct {
	Name        string
	Description string
	Query       string
	Params      map[string]string
	Handler     ToolHandler
	// ReadOnly marks tools that only read state. Read-only tools skip the
	// permission gate entirely, even when they use a Handler.
	ReadOnly bool
}

type Server struct {
	config   *config.Config
	gql      *client.Client
	tools    map[string]*ToolDef
	handlers map[string]ToolHandler
	mu       sync.RWMutex

	// PermsPath points to perms.json (per-container/VM overrides).
	// ecache caches entity ID->name resolution for the override engine.
	PermsPath string
	ecache    entityCache

	// ConfigDir is the plugin config directory (skills/, memory/ live here).
	ConfigDir string
	// clientScope is the sanitized client name from initialize clientInfo,
	// used as the default memory scope for this connection.
	clientScope string
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		config:   cfg,
		tools:    make(map[string]*ToolDef),
		handlers: make(map[string]ToolHandler),
	}
	s.gql = client.New(cfg)
	return s
}

func (s *Server) RegisterTool(t *ToolDef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[t.Name] = t
}

func (s *Server) GraphQLQuery(ctx context.Context, query string, variables map[string]interface{}) (map[string]interface{}, error) {
	return s.gql.Query(ctx, query, variables)
}

// dispatch routes a JSON-RPC request to the appropriate handler.
// Returns nil for notifications (which require no response).
func (s *Server) dispatch(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	switch req.Method {
	case "initialize":
		// Capture the client name for per-client memory scoping.
		if ci, ok := req.Params["clientInfo"].(map[string]interface{}); ok {
			if n, ok := ci["name"].(string); ok {
				s.mu.Lock()
				s.clientScope = agentcontent.SanitizeScope(n)
				s.mu.Unlock()
			}
		}
		return NewSuccessResponse(req.ID, map[string]interface{}{
			"protocolVersion": protocolVersion,
			"capabilities": map[string]interface{}{
				"tools":     map[string]interface{}{},
				"resources": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "unraid-agent",
				"version": ServerVersion,
			},
			"instructions": agentcontent.Instructions(),
		})

	case "ping":
		return &JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]interface{}{}}

	case "tools/list":
		return NewSuccessResponse(req.ID, map[string]interface{}{
			"tools": s.listTools(),
		})

	case "tools/call":
		return s.callTool(ctx, req)

	case "resources/list":
		return NewSuccessResponse(req.ID, map[string]interface{}{
			"resources": s.listResources(),
		})

	case "resources/read":
		uri, _ := req.Params["uri"].(string)
		return s.readResource(req.ID, uri)

	default:
		// Notifications (no id) get no response
		if req.ID == nil || strings.HasPrefix(req.Method, "notifications/") {
			return nil
		}
		return NewErrorResponse(req.ID, -32601, fmt.Sprintf("method not found: %s", req.Method))
	}
}

func (s *Server) listTools() []map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]map[string]interface{}, 0, len(s.tools))
	for _, t := range s.tools {
		properties := map[string]interface{}{}
		required := []string{}
		for name, typ := range t.Params {
			properties[name] = map[string]interface{}{"type": typ}
			required = append(required, name)
		}
		schema := map[string]interface{}{
			"type":       "object",
			"properties": properties,
		}
		if len(required) > 0 {
			schema["required"] = required
		}
		out = append(out, map[string]interface{}{
			"name":        t.Name,
			"description": t.Description,
			"inputSchema": schema,
		})
	}
	return out
}

func (s *Server) callTool(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	name, _ := req.Params["name"].(string)
	if name == "" {
		return NewErrorResponse(req.ID, -32602, "missing tool name")
	}

	s.mu.RLock()
	tool, ok := s.tools[name]
	s.mu.RUnlock()
	if !ok {
		return NewErrorResponse(req.ID, -32602, fmt.Sprintf("unknown tool: %s", name))
	}

	args, _ := req.Params["arguments"].(map[string]interface{})

	// Enforce the permission model before any state-changing execution.
	if err := s.checkPermission(ctx, tool, args); err != nil {
		return toolErrorResult(req.ID, err)
	}

	var text string
	if tool.Handler != nil {
		out, err := tool.Handler(ctx, args)
		if err != nil {
			return toolErrorResult(req.ID, err)
		}
		text = out
	} else if tool.Query != "" {
		data, err := s.GraphQLQuery(ctx, tool.Query, args)
		if err != nil {
			return toolErrorResult(req.ID, err)
		}
		b, _ := json.MarshalIndent(data, "", "  ")
		text = string(b)
	} else {
		text = "ok"
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": text},
			},
		},
	}
}

func toolErrorResult(id interface{}, err error) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": err.Error()},
			},
			"isError": true,
		},
	}
}

// MemoryScope returns the current connection's memory scope (sanitized
// client name from initialize, or "default").
func (s *Server) MemoryScope() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.clientScope == "" {
		return "default"
	}
	return s.clientScope
}

// listResources exposes skills and memory entries as MCP resources.
func (s *Server) listResources() []map[string]interface{} {
	var out []map[string]interface{}
	for _, sk := range agentcontent.ListSkills(s.ConfigDir) {
		out = append(out, map[string]interface{}{
			"uri":         "unraid-agent://skills/" + sk.Name,
			"name":        sk.Name,
			"description": sk.Description,
			"mimeType":    "text/markdown",
		})
	}
	for _, m := range agentcontent.ListMemory(s.ConfigDir, s.MemoryScope()) {
		out = append(out, map[string]interface{}{
			"uri":         fmt.Sprintf("unraid-agent://memory/%s/%s", m.Scope, m.Name),
			"name":        m.Name,
			"description": "memory entry (" + m.Scope + ")",
			"mimeType":    "text/markdown",
		})
	}
	return out
}

// readResource returns a resource's content by URI.
func (s *Server) readResource(id interface{}, uri string) *JSONRPCResponse {
	const skillsPrefix = "unraid-agent://skills/"
	const memoryPrefix = "unraid-agent://memory/"

	switch {
	case strings.HasPrefix(uri, skillsPrefix):
		name := strings.TrimPrefix(uri, skillsPrefix)
		sk, err := agentcontent.GetSkill(s.ConfigDir, name)
		if err != nil {
			return NewErrorResponse(id, -32602, err.Error())
		}
		return NewSuccessResponse(id, map[string]interface{}{
			"contents": []map[string]interface{}{
				{"uri": uri, "mimeType": "text/markdown", "text": sk.Content},
			},
		})

	case strings.HasPrefix(uri, memoryPrefix):
		rest := strings.TrimPrefix(uri, memoryPrefix)
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) != 2 {
			return NewErrorResponse(id, -32602, "invalid memory uri (expected unraid-agent://memory/<scope>/<name>)")
		}
		m, err := agentcontent.GetMemory(s.ConfigDir, parts[0], parts[1])
		if err != nil {
			return NewErrorResponse(id, -32602, err.Error())
		}
		return NewSuccessResponse(id, map[string]interface{}{
			"contents": []map[string]interface{}{
				{"uri": uri, "mimeType": "text/markdown", "text": m.Content},
			},
		})
	}
	return NewErrorResponse(id, -32602, "unknown resource uri: "+uri)
}

// ServeHTTP starts the MCP streamable-http transport.
func (s *Server) ServeHTTP(ctx context.Context, host string, port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleMCP)
	mux.HandleFunc("/mcp", s.handleMCP)

	addr := fmt.Sprintf("%s:%d", host, port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Get().Infof("MCP server listening on http://%s/mcp", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) handleMCP(w http.ResponseWriter, r *http.Request) {
	if !s.authorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, NewErrorResponse(nil, -32700, "parse error"))
			return
		}

		resp := s.dispatch(r.Context(), &req)
		if resp == nil {
			// Notification — accepted, no response body
			w.WriteHeader(http.StatusAccepted)
			return
		}
		writeJSON(w, http.StatusOK, resp)

	case http.MethodGet:
		// SSE listener not supported
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) authorized(r *http.Request) bool {
	if s.config.DisableHTTPAuth {
		return true
	}
	if s.config.BearerToken == "" {
		return true
	}
	auth := r.Header.Get("Authorization")
	return auth == "Bearer "+s.config.BearerToken
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// ServeStdio runs the MCP stdio transport (newline-delimited JSON-RPC).
func (s *Server) ServeStdio(ctx context.Context) error {
	logger.Get().Info("MCP server running on stdio")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	writeCh := make(chan *JSONRPCResponse, 16)
	done := make(chan struct{})

	go func() {
		defer close(done)
		w := bufio.NewWriter(os.Stdout)
		for resp := range writeCh {
			b, err := json.Marshal(resp)
			if err != nil {
				continue
			}
			_, _ = w.Write(b)
			_ = w.WriteByte('\n')
			_ = w.Flush()
		}
	}()

	go func() {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			var req JSONRPCRequest
			if err := json.Unmarshal([]byte(line), &req); err != nil {
				writeCh <- NewErrorResponse(nil, -32700, "parse error")
				continue
			}
			if resp := s.dispatch(ctx, &req); resp != nil {
				writeCh <- resp
			}
		}
		close(writeCh)
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}
	return nil
}
