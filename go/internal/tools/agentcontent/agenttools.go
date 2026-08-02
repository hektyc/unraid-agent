// Package agenttools registers the agent skills + memory tools.
package agenttools

import (
	"context"
	"encoding/json"

	ac "github.com/hektyc/unraid-mcp-server/internal/agentcontent"
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func scopeArg(s *mcp.Server, args map[string]interface{}) string {
	if v, ok := args["scope"].(string); ok && v != "" {
		return ac.SanitizeScope(v)
	}
	return s.MemoryScope()
}

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "skills_list",
		Description: "List available agent skills (workflow playbooks) with name, description, and source (default/custom). Read-only.",
		ReadOnly:    true,
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			skills := ac.ListSkills(s.ConfigDir)
			for i := range skills {
				skills[i].Content = ""
			}
			out, _ := json.MarshalIndent(map[string]interface{}{"skills": skills}, "", "  ")
			return string(out), nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "skills_get",
		Description: "Fetch a skill's full playbook content by name. Read-only.",
		ReadOnly:    true,
		Params:      map[string]string{"name": "string"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			name, _ := args["name"].(string)
			sk, err := ac.GetSkill(s.ConfigDir, name)
			if err != nil {
				return "", err
			}
			out, _ := json.MarshalIndent(sk, "", "  ")
			return string(out), nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "skills_create",
		Description: "Create a custom agent skill (playbook). Requires ALLOW_SKILLS_WRITE.",
		Params:      map[string]string{"name": "string", "description": "string", "content": "string"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			name, _ := args["name"].(string)
			desc, _ := args["description"].(string)
			content, _ := args["content"].(string)
			if err := ac.WriteSkill(s.ConfigDir, name, desc, content); err != nil {
				return "", err
			}
			return "skill created: " + name, nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "skills_update",
		Description: "Update a custom agent skill's content. Requires ALLOW_SKILLS_WRITE.",
		Params:      map[string]string{"name": "string", "description": "string", "content": "string"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			name, _ := args["name"].(string)
			desc, _ := args["description"].(string)
			content, _ := args["content"].(string)
			if desc == "" {
				if sk, err := ac.GetSkill(s.ConfigDir, name); err == nil {
					desc = sk.Description
				}
			}
			if err := ac.WriteSkill(s.ConfigDir, name, desc, content); err != nil {
				return "", err
			}
			return "skill updated: " + name, nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "skills_delete",
		Description: "Delete a custom agent skill. Requires ALLOW_SKILLS_WRITE.",
		Params:      map[string]string{"name": "string"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			name, _ := args["name"].(string)
			if err := ac.DeleteSkill(s.ConfigDir, name); err != nil {
				return "", err
			}
			return "skill deleted: " + name, nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "memory_list",
		Description: "List memory entries visible to you (your scope plus defaults). Read-only.",
		ReadOnly:    true,
		Params:      map[string]string{"scope": "string"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			out, _ := json.MarshalIndent(map[string]interface{}{"memory": ac.ListMemory(s.ConfigDir, scopeArg(s, args))}, "", "  ")
			return string(out), nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "memory_get",
		Description: "Read a memory entry's content by name. Read-only.",
		ReadOnly:    true,
		Params:      map[string]string{"name": "string", "scope": "string"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			name, _ := args["name"].(string)
			m, err := ac.GetMemory(s.ConfigDir, scopeArg(s, args), name)
			if err != nil {
				return "", err
			}
			out, _ := json.MarshalIndent(m, "", "  ")
			return string(out), nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "memory_write",
		Description: "Create or update a memory entry in your scope. Requires ALLOW_MEMORY_WRITE.",
		Params:      map[string]string{"name": "string", "content": "string", "scope": "string"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			name, _ := args["name"].(string)
			content, _ := args["content"].(string)
			scope := scopeArg(s, args)
			if err := ac.WriteMemory(s.ConfigDir, scope, name, content); err != nil {
				return "", err
			}
			return "memory written: " + scope + "/" + name, nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "memory_delete",
		Description: "Delete a memory entry from your scope. Requires ALLOW_MEMORY_WRITE.",
		Params:      map[string]string{"name": "string", "scope": "string"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			name, _ := args["name"].(string)
			scope := scopeArg(s, args)
			if err := ac.DeleteMemory(s.ConfigDir, scope, name); err != nil {
				return "", err
			}
			return "memory deleted: " + scope + "/" + name, nil
		},
	})
}
