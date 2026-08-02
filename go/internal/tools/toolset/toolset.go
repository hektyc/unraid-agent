// Package toolset registers the agent-toolset management tools. These tools
// are ALWAYS registered (never filtered) so the toolset can never disable
// itself out of reach.
package toolset

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

// TS is the shared file-backed toolset state, wired up in main.
var TS *config.Toolset

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "toolset_list",
		Description: "List the agent toolset: every domain with its enabled state and any per-tool overrides. Read-only.",
		ReadOnly:    true,
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			domains, overrides := TS.State()
			var domList []string
			for d := range domains {
				domList = append(domList, d)
			}
			sort.Strings(domList)
			out := map[string]interface{}{
				"domains":   domains,
				"overrides": overrides,
				"allDomains": config.ToolDomains,
				"note": "toolset_set target=domain|tool name=<name> value=enabled|disabled|default — applies instantly, no restart needed",
			}
			_ = domList
			b, _ := json.MarshalIndent(out, "", "  ")
			return string(b), nil
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "toolset_set",
		Description: "Enable, disable, or reset a tool domain or individual tool. Applies instantly without a daemon restart. Requires Allow Setting Updates.",
		Params:      map[string]string{"target": "string", "name": "string", "value": "string"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			target, _ := args["target"].(string)
			name, _ := args["name"].(string)
			value, _ := args["value"].(string)
			if err := TS.ApplySet(target, name, value); err != nil {
				return "", err
			}
			return "toolset updated: " + target + " " + name + " -> " + value, nil
		},
	})
}
