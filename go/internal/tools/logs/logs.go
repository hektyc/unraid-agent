package logs

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "logs_list",
		Description: "List all Unraid system log files with name, path, size, and modified time (syslog, docker.log, graphql-api.log, plugin logs, etc.). Read-only.",
		Query:       `query { logFiles { name path size modifiedAt } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "logs_read",
		Description: "Read an Unraid system log file by path with line pagination: lines = number of lines to return (from the end by default), startLine = line to start from. Use logs_list to discover paths. Read-only.",
		Query:       `query($path: String!, $lines: Int, $startLine: Int) { logFile(path: $path, lines: $lines, startLine: $startLine) { path content totalLines startLine } }`,
		Params: map[string]string{
			"path":      "string",
			"lines":     "number",
			"startLine": "number",
		},
	})
}
