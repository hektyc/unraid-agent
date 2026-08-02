package mcp

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

// entityActions maps mutation tool names to the per-entity action key
// stored in perms.json (containers/vms).
var entityActions = map[string]string{
	"docker_start":            "start",
	"docker_stop":             "stop",
	"docker_restart":          "restart",
	"docker_pause":            "pause",
	"docker_unpause":          "unpause",
	"docker_remove_container": "remove",
	"docker_update_container": "update",
	"vm_start":                "start",
	"vm_stop":                 "stop",
	"vm_pause":                "pause",
	"vm_resume":               "resume",
	"vm_force_stop":           "force_stop",
	"vm_reboot":               "reboot",
	"vm_reset":                "reset",
}

// entityPermsFile mirrors /boot/config/plugins/unraid-agent/perms.json.
// Values: "inherit" (use global), "allow", "deny".
type entityPermsFile struct {
	Containers map[string]map[string]string `json:"containers"`
	VMs        map[string]map[string]string `json:"vms"`
	Plugins    map[string]map[string]string `json:"plugins"`
}

type entityCache struct {
	mu    sync.RWMutex
	at    time.Time
	names map[string]map[string]string // kind -> id -> name
}

// loadEntityPerms reads perms.json fresh on every gated call. The file is
// tiny; this keeps overrides effective immediately without a daemon restart.
func (s *Server) loadEntityPerms() *entityPermsFile {
	p := &entityPermsFile{}
	if s.PermsPath == "" {
		return p
	}
	data, err := os.ReadFile(s.PermsPath)
	if err != nil {
		return p
	}
	_ = json.Unmarshal(data, p)
	return p
}

// entityOverride returns "allow", "deny", or "" (inherit / no override) for
// a mutation tool call targeting a specific container or VM.
func (s *Server) entityOverride(ctx context.Context, toolName string, args map[string]interface{}) string {
	action, tracked := entityActions[toolName]
	if !tracked {
		return ""
	}

	var kind, id string
	switch {
	case strings.HasPrefix(toolName, "docker_"):
		kind = "containers"
		id, _ = args["container_id"].(string)
	case strings.HasPrefix(toolName, "vm_"):
		kind = "vms"
		id, _ = args["vm_id"].(string)
	}
	if id == "" {
		return ""
	}

	name := s.resolveEntityName(ctx, kind, id)
	if name == "" {
		return ""
	}

	perms := s.loadEntityPerms()
	var entities map[string]map[string]string
	if kind == "containers" {
		entities = perms.Containers
	} else {
		entities = perms.VMs
	}
	if ent, ok := entities[name]; ok {
		if v, ok := ent[action]; ok && v != "inherit" && v != "global" {
			return v
		}
	}
	return ""
}

// resolveEntityName maps a PrefixedID to the entity's current name using a
// 30-second cached copy of the container/VM list, refreshed on miss/expiry.
func (s *Server) resolveEntityName(ctx context.Context, kind, id string) string {
	s.ecache.mu.RLock()
	if time.Since(s.ecache.at) < 30*time.Second {
		if n, ok := s.ecache.names[kind][id]; ok {
			s.ecache.mu.RUnlock()
			return n
		}
	}
	s.ecache.mu.RUnlock()

	s.refreshEntityCache(ctx, kind)

	s.ecache.mu.RLock()
	defer s.ecache.mu.RUnlock()
	return s.ecache.names[kind][id]
}

func (s *Server) refreshEntityCache(ctx context.Context, kind string) {	var query string
	if kind == "containers" {
		query = `query { docker { containers { id names } } }`
	} else {
		query = `query { vms { domains { id name } } }`
	}
	data, err := s.gql.Query(ctx, query, nil)
	if err != nil {
		return
	}

	fresh := map[string]string{}
	if kind == "containers" {
		docker, _ := data["docker"].(map[string]interface{})
		list, _ := docker["containers"].([]interface{})
		for _, item := range list {
			c, _ := item.(map[string]interface{})
			id, _ := c["id"].(string)
			names, _ := c["names"].([]interface{})
			if len(names) > 0 {
				if n, ok := names[0].(string); ok {
					fresh[id] = strings.TrimPrefix(n, "/")
				}
			}
		}
	} else {
		vms, _ := data["vms"].(map[string]interface{})
		list, _ := vms["domains"].([]interface{})
		for _, item := range list {
			d, _ := item.(map[string]interface{})
			id, _ := d["id"].(string)
			name, _ := d["name"].(string)
			if id != "" && name != "" {
				fresh[id] = name
			}
		}
	}

	s.ecache.mu.Lock()
	if s.ecache.names == nil {
		s.ecache.names = map[string]map[string]string{}
	}
	s.ecache.names[kind] = fresh
	s.ecache.at = time.Now()
	s.ecache.mu.Unlock()
}

// pluginRemoveOverride evaluates per-plugin "remove" overrides for a
// plugin_remove call. names arrive directly in the args (no ID resolution
// needed). Returns ("deny"|"allow"|"") and whether unraid-agent itself was
// targeted. Verdict rules across the names array: any explicit deny blocks
// the call; every name explicitly allowed permits it; otherwise "" falls
// through to the global toggles.
func (s *Server) pluginRemoveOverride(args map[string]interface{}) (string, bool) {
	raw, _ := args["names"].([]interface{})
	if len(raw) == 0 {
		return "", false
	}

	perms := s.loadEntityPerms()
	anyDeny := false
	allAllow := true

	for _, item := range raw {
		name, _ := item.(string)
		if name == "unraid-agent" {
			return "", true
		}
		verdict := ""
		if ent, ok := perms.Plugins[name]; ok {
			if v, ok := ent["remove"]; ok && v != "inherit" && v != "global" {
				verdict = v
			}
		}
		if verdict == "deny" {
			anyDeny = true
		}
		if verdict != "allow" {
			allAllow = false
		}
	}

	if anyDeny {
		return "deny", false
	}
	if allAllow {
		return "allow", false
	}
	return "", false
}
