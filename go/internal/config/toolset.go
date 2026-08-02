package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Toolset resolves tool registration from the config FILE (not process env),
// reloading automatically when the file changes. That lets toolset_set apply
// instantly from any MCP client — no daemon restart required.
type Toolset struct {
	path    string
	mu      sync.RWMutex
	modTime time.Time
	domains map[string]bool
	tools   map[string]bool // explicit per-tool overrides
}

func NewToolset(configPath string) *Toolset {
	t := &Toolset{path: configPath, domains: map[string]bool{}, tools: map[string]bool{}}
	t.reload()
	return t
}

func (t *Toolset) maybeReload() {
	info, err := os.Stat(t.path)
	if err != nil {
		return
	}
	t.mu.RLock()
	fresh := info.ModTime().After(t.modTime)
	t.mu.RUnlock()
	if fresh {
		t.reload()
	}
}

func (t *Toolset) reload() {
	data, err := os.ReadFile(t.path)
	if err != nil {
		return
	}
	domains := map[string]bool{}
	tools := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "UNRAID_MCP_") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		b := val == "true"
		switch {
		case strings.HasPrefix(key, "UNRAID_MCP_DOMAIN_"):
			domains[strings.ToLower(strings.TrimPrefix(key, "UNRAID_MCP_DOMAIN_"))] = b
		case strings.HasPrefix(key, "UNRAID_MCP_TOOL_"):
			tools[strings.TrimPrefix(key, "UNRAID_MCP_TOOL_")] = b
		}
	}
	info, _ := os.Stat(t.path)
	t.mu.Lock()
	t.domains = domains
	t.tools = tools
	if info != nil {
		t.modTime = info.ModTime()
	}
	t.mu.Unlock()
}

// Enabled reports whether a tool should be registered, with per-tool
// overrides beating the domain toggle in both directions.
func (t *Toolset) Enabled(tool string) bool {
	t.maybeReload()
	t.mu.RLock()
	defer t.mu.RUnlock()
	upper := strings.ToUpper(tool)
	if v, ok := t.tools[upper]; ok {
		return v
	}
	if v, ok := t.domains[ToolDomainOf(tool)]; ok {
		return v
	}
	return true
}

// State returns a snapshot for toolset_list.
func (t *Toolset) State() (map[string]bool, map[string]bool) {
	t.maybeReload()
	t.mu.RLock()
	defer t.mu.RUnlock()
	d := map[string]bool{}
	for k, v := range t.domains {
		d[k] = v
	}
	tools := map[string]bool{}
	for k, v := range t.tools {
		tools[k] = v
	}
	return d, tools
}

// ApplySet writes one toolset change to the config file and reloads.
// value: "enabled" (true), "disabled" (false), "default" (remove override).
func (t *Toolset) ApplySet(target, name, value string) error {
	var key string
	switch target {
	case "domain":
		found := false
		for _, d := range ToolDomains {
			if d == strings.ToLower(name) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown domain: %s", name)
		}
		key = "UNRAID_MCP_DOMAIN_" + strings.ToUpper(name)
	case "tool":
		key = "UNRAID_MCP_TOOL_" + strings.ToUpper(name)
	default:
		return fmt.Errorf("target must be domain or tool")
	}

	var newVal *string
	switch value {
	case "enabled":
		s := "true"
		newVal = &s
	case "disabled":
		s := "false"
		newVal = &s
	case "default":
		newVal = nil // remove the line
	default:
		return fmt.Errorf("value must be enabled, disabled, or default")
	}

	data, err := os.ReadFile(t.path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	replaced := false
	var out []string
	for _, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			replaced = true
			if newVal != nil {
				out = append(out, fmt.Sprintf(`%s="%s"`, key, *newVal))
			}
			continue
		}
		out = append(out, line)
	}
	if !replaced && newVal != nil {
		out = append(out, fmt.Sprintf(`%s="%s"`, key, *newVal))
	}
	tmp := t.path + ".tmp"
	if err := os.WriteFile(tmp, []byte(strings.Join(out, "\n")), 0600); err != nil {
		return err
	}
	if err := os.Rename(tmp, t.path); err != nil {
		return err
	}
	t.reload()
	return nil
}
