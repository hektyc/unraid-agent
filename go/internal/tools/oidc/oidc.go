package oidc

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "oidc_providers",
		Description: "List providers. Read-only.",
		Query:       `query { oidcProviders { id name clientId issuer authorizationEndpoint tokenEndpoint jwksUri scopes buttonText buttonIcon buttonVariant authorizationRuleMode } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "oidc_create",
		Description: "Create provider.",
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			return "", fmt.Errorf("OIDC provider creation is not available via the Unraid GraphQL API; use the Unraid WebGUI instead")
		},
		Params: map[string]string{
			"name":     "string",
			"issuer":   "string",
			"clientId": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "oidc_delete",
		Description: "Delete provider.",
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			return "", fmt.Errorf("OIDC provider deletion is not available via the Unraid GraphQL API; use the Unraid WebGUI instead")
		},
		Params: map[string]string{
			"provider_id": "string",
		},
	})
}
