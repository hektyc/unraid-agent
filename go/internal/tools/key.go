package key

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_list",
		Description: "Get all API keys. Read-only.",
		Query:       `query { key { list { id name roles permissions } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_get",
		Description: "Get single API key details.",
		Query:       `query KeyGet($id: ID!) { key { get(id: $id) { id name roles permissions } } }`,
		Params: map[string]string{
			"key_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_create",
		Description: "Create an API key.",
		Query:       `mutation CreateKey($name: String!) { key { create(name: $name) { id name token } } }`,
		Params: map[string]string{
			"name": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_update",
		Description: "Update an API key.",
		Query:       `mutation UpdateKey($id: ID!, $name: String) { key { update(id: $id, name: $name) { id name } } }`,
		Params: map[string]string{
			"key_id": "string", // required=true
			"name": "string", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_delete",
		Description: "Delete an API key.",
		Query:       `mutation DeleteKey($id: ID!) { key { delete(id: $id) { status } } }`,
		Params: map[string]string{
			"key_id": "string", // required=true
			"confirm": "boolean", // required=true
		},
	})
}