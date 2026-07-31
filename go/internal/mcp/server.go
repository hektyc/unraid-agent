package mcp

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hektyc/unraid-mcp-server/go/internal/client"
	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/logger"
)

type ToolHandler func(ctx context.Context, arguments map[string]interface{}) (string, error)

type ToolDef struct {
	Name        string
	Description string
	Query       string
	Params      map[string]string
	Handler     ToolHandler
}

type Server struct {
	config   *config.Config
	gql      *client.Client
	tools    map[string]*ToolDef
	handlers map[string]ToolHandler
	mu       sync.RWMutex
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		config:   cfg,
		tools:    make(map[string]*ToolDef),
		handlers: make(map[string]ToolHandler),
	}
	s.gql = client.New(cfg)
	s.registerTools()
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

func (s *Server) registerTools() {
	// Tool registration is done via init() functions in subpackages
}
