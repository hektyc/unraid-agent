package mcp

import (
	"fmt"
	"strings"
)

// Permission resolution order (all evaluated for mutation tools only):
//   1. READ_ONLY=true            -> every mutation blocked
//   2. ALLOW_DESTRUCTIVE=true    -> every mutation allowed
//   3. Category toggle           -> ALLOW_ARRAY_ACTIONS / ALLOW_DOCKER_ACTIONS /
//                                   ALLOW_VM_ACTIONS allow their group
//   4. Specific toggle           -> allows that single tool
//   5. Mutation with no toggle   -> allowed (not considered destructive)
//
// Queries (read-only tools) are always allowed.

// toolPermissions maps a mutation tool name to the config toggles that allow
// it (any-of). Category toggles are included in each tool's list.
var toolPermissions = map[string][]string{
	// Array
	"array_start":         {"ALLOW_ARRAY_START", "ALLOW_ARRAY_ACTIONS"},
	"array_stop":          {"ALLOW_ARRAY_STOP", "ALLOW_ARRAY_ACTIONS"},
	"array_parity_start":  {"ALLOW_ARRAY_ACTIONS"},
	"array_parity_pause":  {"ALLOW_ARRAY_ACTIONS"},
	"array_parity_resume": {"ALLOW_ARRAY_ACTIONS"},
	"array_parity_cancel": {"ALLOW_ARRAY_ACTIONS"},
	"array_add_disk":      {"ALLOW_ARRAY_ADD_DISK", "ALLOW_ARRAY_ACTIONS"},
	"array_remove_disk":   {"ALLOW_ARRAY_REMOVE_DISK", "ALLOW_ARRAY_ACTIONS"},

	// Connect
	"connect_sign_in":  {"ALLOW_CONNECT_ACTIONS"},
	"connect_sign_out": {"ALLOW_CONNECT_ACTIONS"},

	// Customization
	"customization_set_theme": {"ALLOW_SETTING_UPDATES"},

	// Docker
	"docker_start":            {"ALLOW_DOCKER_ACTIONS"},
	"docker_stop":             {"ALLOW_CONTAINER_STOP", "ALLOW_DOCKER_ACTIONS"},
	"docker_restart":          {"ALLOW_CONTAINER_RESTART", "ALLOW_DOCKER_ACTIONS"},
	"docker_pause":            {"ALLOW_DOCKER_ACTIONS"},
	"docker_unpause":          {"ALLOW_DOCKER_ACTIONS"},
	"docker_remove_container": {"ALLOW_CONTAINER_REMOVE", "ALLOW_DOCKER_ACTIONS"},
	"docker_update_container": {"ALLOW_DOCKER_ACTIONS"},
	"docker_create_folder":    {"ALLOW_DOCKER_ACTIONS"},
	"docker_delete_entries":   {"ALLOW_DOCKER_ACTIONS"},
	"docker_refresh_digests":  {"ALLOW_DOCKER_ACTIONS"},

	// API keys
	"key_create": {"ALLOW_API_KEY_CREATE"},
	"key_update": {"ALLOW_SETTING_UPDATES"},
	"key_delete": {"ALLOW_API_KEY_DELETE"},

	// Notifications
	"notification_delete":          {"ALLOW_NOTIFICATION_DELETE"},
	"notification_delete_archived": {"ALLOW_NOTIFICATION_DELETE"},
	"notification_archive":         {}, // state change but non-destructive
	"notification_create":          {}, // non-destructive

	// Onboarding
	"onboarding_complete": {"ALLOW_ONBOARDING_ACTIONS"},
	"onboarding_reset":    {"ALLOW_ONBOARDING_ACTIONS"},

	// Plugins
	"plugin_add":     {"ALLOW_PLUGIN_INSTALL"},
	"plugin_remove":  {"ALLOW_PLUGIN_REMOVE"},
	"plugin_install": {"ALLOW_PLUGIN_INSTALL"},

	// RClone
	"rclone_create_remote": {"ALLOW_RCLONE_OPERATIONS"},
	"rclone_delete_remote": {"ALLOW_RCLONE_OPERATIONS"},

	// Settings
	"setting_update":             {"ALLOW_SETTING_UPDATES"},
	"setting_update_ssh":         {"ALLOW_SSH_UPDATE"},
	"setting_update_system_time": {"ALLOW_TIME_UPDATE"},
	"setting_configure_ups":      {"ALLOW_SETTING_UPDATES"},

	// VMs
	"vm_start":      {"ALLOW_VM_ACTIONS"},
	"vm_stop":       {"ALLOW_VM_STOP", "ALLOW_VM_ACTIONS"},
	"vm_pause":      {"ALLOW_VM_ACTIONS"},
	"vm_resume":     {"ALLOW_VM_ACTIONS"},
	"vm_force_stop": {"ALLOW_VM_FORCE_STOP", "ALLOW_VM_ACTIONS"},
	"vm_reboot":     {"ALLOW_VM_ACTIONS"},
	"vm_reset":      {"ALLOW_VM_RESET", "ALLOW_VM_ACTIONS"},
}

// isMutationTool reports whether executing this tool changes state.
// GraphQL mutations mutate; tools whose Handler self-gates (oidc_*) are
// treated as mutations too so READ_ONLY still short-circuits them.
func isMutationTool(t *ToolDef) bool {
	if t.Handler != nil {
		return true
	}
	q := strings.TrimSpace(t.Query)
	return strings.HasPrefix(q, "mutation")
}

// checkPermission enforces the permission model for a tool call.
// Returns nil when the call is allowed, or a descriptive error.
func (s *Server) checkPermission(t *ToolDef) error {
	if !isMutationTool(t) {
		return nil
	}

	cfg := s.config

	if cfg.ReadOnly {
		return fmt.Errorf("permission denied: %s is a state-changing operation and READ_ONLY mode is enabled (Settings → unRAID Agent → Permissions → Read Only)", t.Name)
	}

	if cfg.AllowDestructive {
		return nil
	}

	toggles, gated := toolPermissions[t.Name]
	if !gated || len(toggles) == 0 {
		// Mutation without a specific toggle — allowed when not READ_ONLY
		return nil
	}

	if anyToggleEnabled(cfg, toggles) {
		return nil
	}

	return fmt.Errorf("permission denied: %s is blocked by configuration — enable one of [%s] in Settings → unRAID Agent → Permissions (or ALLOW_DESTRUCTIVE for all)",
		t.Name, strings.Join(toggles, ", "))
}

func anyToggleEnabled(cfg interface{ permissionValue(string) bool }, toggles []string) bool {
	for _, t := range toggles {
		if cfg.permissionValue(t) {
			return true
		}
	}
	return false
}
