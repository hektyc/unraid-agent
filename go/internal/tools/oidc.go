package tools

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "oidc_providers",
		Description: "List providers. Read-only.",
		Query:       `query { oidc { providers { id name enabled } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "oidc_create",
		Description: "Create provider.",
		Query:       `mutation CreateOidc($name: String!, $issuer: String!, $clientId: String!) { oidc { create(name: $name, issuer: $issuer, clientId: $clientId) { id } } }`,
		Params: map[string]string{
			"name": "string", // required=true
			"issuer": "string", // required=true
			"clientId": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "oidc_delete",
		Description: "Delete provider.",
		Query:       `mutation DeleteOidc($id: ID!) { oidc { delete(id: $id) { status } } }`,
		Params: map[string]string{
			"provider_id": "string", // required=true
		},
	})
}