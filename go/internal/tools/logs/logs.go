package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

var anon = newAnonymizer()

func numArg(args map[string]interface{}, key string, def int) int {
	if v, ok := args[key].(float64); ok && v > 0 {
		return int(v)
	}
	return def
}

// scrub applies the user's anonymize setting to log text.
func scrub(ctx context.Context, s *mcp.Server, cfg *config.Config, text string) string {
	if cfg.AnonymizeLogs {
		return anon.scrub(ctx, s, text)
	}
	return text
}

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "logs_list",
		Description: "List all Unraid system log files with name, path, size, and modified time (syslog, docker.log, graphql-api.log, plugin logs, etc.). Read-only.",
		Query:       `query { logFiles { name path size modifiedAt } }`,
	})

	s.RegisterTool(&mcp.ToolDef{
		Name:        "logs_read",
		Description: "Read an Unraid system log file by path with line pagination: lines = number of lines to return (from the end by default), startLine = line to start from. Use logs_list to discover paths. Honors the Anonymize Logs setting. Read-only.",
		ReadOnly:    true,
		Params:      map[string]string{"path": "string", "lines": "number", "startLine": "number"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			if path == "" {
				return "", fmt.Errorf("path is required (see logs_list)")
			}
			vars := map[string]interface{}{"path": path}
			if v, ok := args["lines"].(float64); ok && v > 0 {
				vars["lines"] = int(v)
			}
			if v, ok := args["startLine"].(float64); ok && v > 0 {
				vars["startLine"] = int(v)
			}
			data, err := s.GraphQLQuery(ctx, `query($path: String!, $lines: Int, $startLine: Int) { logFile(path: $path, lines: $lines, startLine: $startLine) { path content totalLines startLine } }`, vars)
			if err != nil {
				return "", err
			}
			if lf, ok := data["logFile"].(map[string]interface{}); ok {
				if c, ok := lf["content"].(string); ok {
					lf["content"] = scrub(ctx, s, cfg, c)
				}
				out, _ := json.MarshalIndent(map[string]interface{}{"logFile": lf}, "", "  ")
				return string(out), nil
			}
			return "", fmt.Errorf("log file not found: %s", path)
		},
	})

	s.RegisterTool(&mcp.ToolDef{
		Name:        "plugin_logs",
		Description: "Find logs related to a specific plugin by name: reads the plugin's own log files when present and scans recent syslog lines mentioning it. Best-effort — plugins have no unified log system. Honors the Anonymize Logs setting. Read-only.",
		ReadOnly:    true,
		Params:      map[string]string{"plugin": "string", "lines": "number", "syslog_lines": "number"},
		Handler: func(ctx context.Context, args map[string]interface{}) (string, error) {
			plugin, _ := args["plugin"].(string)
			if plugin == "" {
				return "", fmt.Errorf("plugin name is required")
			}
			lines := numArg(args, "lines", 100)
			syslogLines := numArg(args, "syslog_lines", 2000)
			needle := strings.ToLower(plugin)

			result := map[string]interface{}{"plugin": plugin, "files": []interface{}{}, "syslog_mentions": []string{}}

			// 1. Plugin's own log files (name match in logFiles)
			data, err := s.GraphQLQuery(ctx, `query { logFiles { name path size } }`, nil)
			if err != nil {
				return "", err
			}
			matched := 0
			if files, ok := data["logFiles"].([]interface{}); ok {
				for _, f := range files {
					fm, _ := f.(map[string]interface{})
					name, _ := fm["name"].(string)
					path, _ := fm["path"].(string)
					if name == "" || path == "" || !strings.Contains(strings.ToLower(name), needle) {
						continue
					}
					if matched >= 3 {
						break
					}
					matched++
					fd, err := s.GraphQLQuery(ctx, `query($path: String!, $lines: Int) { logFile(path: $path, lines: $lines) { path content totalLines startLine } }`, map[string]interface{}{"path": path, "lines": lines})
					entry := map[string]interface{}{"name": name, "path": path, "size": fm["size"]}
					if err == nil {
						if lf, ok := fd["logFile"].(map[string]interface{}); ok {
							if c, ok := lf["content"].(string); ok {
								lf["content"] = scrub(ctx, s, cfg, c)
							}
							entry["tail"] = lf
						}
					} else {
						entry["error"] = err.Error()
					}
					result["files"] = append(result["files"].([]interface{}), entry)
				}
			}

			// 2. Recent syslog lines mentioning the plugin
			sd, err := s.GraphQLQuery(ctx, `query($path: String!, $lines: Int) { logFile(path: $path, lines: $lines) { content totalLines startLine } }`, map[string]interface{}{"path": "/var/log/syslog", "lines": syslogLines})
			if err == nil {
				if lf, ok := sd["logFile"].(map[string]interface{}); ok {
					if content, ok := lf["content"].(string); ok {
						var matches []string
						for _, line := range strings.Split(content, "\n") {
							if strings.Contains(strings.ToLower(line), needle) {
								matches = append(matches, line)
								if len(matches) >= lines {
									break
								}
							}
						}
						if len(matches) > 0 {
							result["syslog_mentions"] = strings.Split(scrub(ctx, s, cfg, strings.Join(matches, "\n")), "\n")
						}
					}
					result["syslog_scanned"] = map[string]interface{}{"lines": syslogLines, "totalLines": lf["totalLines"]}
				}
			} else {
				result["syslog_error"] = err.Error()
			}

			if len(result["files"].([]interface{})) == 0 && len(result["syslog_mentions"].([]string)) == 0 {
				result["note"] = "no plugin-specific log files or recent syslog mentions found — the plugin may log elsewhere or not at all recently"
			}
			out, _ := json.MarshalIndent(result, "", "  ")
			return string(out), nil
		},
	})
}
