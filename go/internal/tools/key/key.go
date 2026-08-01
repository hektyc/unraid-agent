package key

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_list",
		Description: "Get all API keys. Read-only.",
		Query:       `query { apiKeys { id name description roles createdAt } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_get",
		Description: "Get single API key details.",
		Query:       `query($key_id: PrefixedID!) { apiKey(id: $key_id) { id name description roles createdAt } }`,
		Params: map[string]string{
			"key_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_create",
		Description: "Create an API key.",
		Query:       `mutation($name: String!) { apiKey { create(input: {name: $name}) { id name key description roles createdAt } } }`,
		Params: map[string]string{
			"name": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_update",
		Description: "Update an API key.",
		Query:       `mutation($key_id: PrefixedID!, $name: String) { apiKey { update(input: {id: $key_id, name: $name}) { id name description } } }`,
		Params: map[string]string{
			"key_id": "string",
			"name":   "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "key_delete",
		Description: "Delete an API key.",
		Query:       `mutation($key_id: PrefixedID!) { apiKey { delete(input: {ids: [$key_id]}) } }`,
		Params: map[string]string{
			"key_id": "string",
		},
	})
}
