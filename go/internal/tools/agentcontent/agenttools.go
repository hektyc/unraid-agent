// Package agenttools registers the agent skills + memory tools.
package agenttools

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	s.RegisterTool(&mcp.ToolDef{
		Name:        "agent_endpoint_log",
		Description: "Read recent lines of the plugin's endpoint request log (WebGUI action diagnostics: content reads, saves, deletes, exports). Read-only.",
		ReadOnly:    true,
		Params:      map[string]string{"tail": "number"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			tail := 50
			if v, ok := args["tail"].(float64); ok && v > 0 {
				tail = int(v)
			}
			if tail > 200 {
				tail = 200
			}
			data, err := os.ReadFile(filepath.Join(s.ConfigDir, "logs", "endpoints.log"))
			if err != nil {
				return "no endpoint log entries yet", nil
			}
			lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
			if len(lines) > tail {
				lines = lines[len(lines)-tail:]
			}
			return strings.Join(lines, "\n"), nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "agent_page_check",
		Description: "Diagnostic: checksum and inspect the installed UnraidAgentPermissions.page on disk (md5, size, and the bytes around the tools pane tag) to verify file integrity. Read-only.",
		ReadOnly:    true,
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			path := "/usr/local/emhttp/plugins/unraid-agent/UnraidAgentPermissions.page"
			data, err := os.ReadFile(path)
			if err != nil {
				return "", fmt.Errorf("read page file: %w", err)
			}
			sum := md5.Sum(data)
			idx := bytes.Index(data, []byte(`data-tab="tools"`))
			region := "(tag not found)"
			if idx >= 0 {
				start := idx - 150
				if start < 0 {
					start = 0
				}
				end := idx + 120
				if end > len(data) {
					end = len(data)
				}
				region = string(data[start:end])
			}
			return fmt.Sprintf("md5=%x\nsize=%d bytes\nexpected_md5(v2026.08.02.2207)=e4b5ee0eb44d2b4cd1a554dfbdd6907e\n--- region ---\n%s", sum, len(data), region), nil
		},
	})
}
