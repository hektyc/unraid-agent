package mcp

import (
	"context"
	"sync"

	"github.com/hektyc/unraid-mcp-server/internal/client"
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/tools"
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
	tools.RegisterarrayTools(s, s.config)
	tools.RegisterconnectTools(s, s.config)
	tools.RegistercustomizationTools(s, s.config)
	tools.RegisterdockerTools(s, s.config)
	tools.RegisterhealthTools(s, s.config)
	tools.RegisterhelpTools(s, s.config)
	tools.RegisterkeyTools(s, s.config)
	tools.RegisterliveTools(s, s.config)
	tools.RegisternotificationTools(s, s.config)
	tools.RegisteroidcTools(s, s.config)
	tools.RegisteronboardingTools(s, s.config)
	tools.RegisterpluginTools(s, s.config)
	tools.RegisterrcloneTools(s, s.config)
	tools.RegistersettingTools(s, s.config)
	tools.RegistersystemTools(s, s.config)
	tools.RegisteruserTools(s, s.config)
	tools.RegistervmTools(s, s.config)
}

func (s *Server) ServeHTTP(ctx context.Context, host string, port int) error {
	return nil
}

func (s *Server) ServeStdio(ctx context.Context) error {
	return nil
}
