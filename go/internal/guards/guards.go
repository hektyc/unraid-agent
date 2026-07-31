package guards

import (
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
)

type SafetyError struct {
	Action    string
	Subaction string
	Message   string
}

func (e *SafetyError) Error() string {
	return e.Message
}

func EnforceReadOnly(cfg *config.Config, action, subaction string) error {
	if !cfg.ReadOnly {
		return nil
	}
	return &SafetyError{
		Action:    action,
		Subaction: subaction,
		Message:   fmt.Sprintf("READ_ONLY mode is enabled. Action '%s.%s' is blocked.", action, subaction),
	}
}

func EnforceToggle(cfg *config.Config, action, subaction string) error {
	if !isDestructive(action, subaction) {
		return nil
	}
	toggle := getToggleKey(action, subaction)
	if !isToggleEnabled(cfg, action, subaction) {
		return &SafetyError{
			Action:    action,
			Subaction: subaction,
			Message:   fmt.Sprintf("Action '%s.%s' requires '%s=true' in environment.", action, subaction, toggle),
		}
	}
	return nil
}

func IsDestructive(action, subaction string) bool {
	key := action + ":" + subaction
	_, ok := destructiveActions[key]
	return ok
}

func GetDestructiveWarning(action, subaction string) string {
	warnings := map[string]string{
		"array:stop":                   "Stops the Unraid array, unmounting all shares and disks.",
		"array:start":                  "Starts the Unraid array, mounting all shares and disks.",
		"array:remove_disk":            "Removes a disk from the array. Data loss may occur.",
		"array:clear_disk_stats":       "Clears all I/O statistics. Irreversible.",
		"docker:stop":                  "Stops a running Docker container.",
		"docker:remove_container":      "Permanently removes a Docker container.",
		"docker:restart":               "Restarts a Docker container.",
		"vm:stop":                      "Gracefully stops a virtual machine.",
		"vm:force_stop":                "Hard power-off. Data loss is possible.",
		"vm:reset":                     "Hard reset without graceful shutdown.",
		"notification:delete":          "Permanently deletes a notification.",
		"notification:delete_archived": "Permanently deletes all archived notifications.",
		"key:delete":                   "Revokes an API key immediately.",
		"plugin:add":                   "Installs a plugin. Runs code as root.",
		"plugin:remove":                "Uninstalls a plugin.",
		"plugin:install":               "Installs a .plg file from a URL. Runs code as root.",
		"rclone:delete_remote":         "Removes an rclone remote configuration.",
		"setting:configure_ups":        "Updates UPS configuration.",
		"setting:update_ssh":           "Updates SSH daemon settings.",
		"setting:update_system_time":   "Updates system clock and NTP.",
		"connect:sign_in":              "Signs in to Unraid Connect.",
		"connect:sign_out":             "Signs out of Unraid Connect.",
		"disk:flash_backup":            "Initiates an rclone backup of the flash drive.",
	}
	if msg, ok := warnings[action+":"+subaction]; ok {
		return msg
	}
	return fmt.Sprintf("Potentially destructive action: %s.%s", action, subaction)
}

func isDestructive(action, subaction string) bool {
	key := action + ":" + subaction
	_, ok := destructiveActions[key]
	return ok
}

func getToggleKey(action, subaction string) string {
	toggles := map[string]string{
		"array:stop":                  "ALLOW_ARRAY_STOP",
		"array:start":                 "ALLOW_ARRAY_START",
		"array:add_disk":              "ALLOW_ARRAY_ADD_DISK",
		"array:remove_disk":           "ALLOW_ARRAY_REMOVE_DISK",
		"array:clear_disk_stats":      "ALLOW_ARRAY_CLEAR_STATS",
		"docker:stop":                 "ALLOW_CONTAINER_STOP",
		"docker:remove_container":     "ALLOW_CONTAINER_REMOVE",
		"docker:restart":              "ALLOW_CONTAINER_RESTART",
		"vm:stop":                     "ALLOW_VM_STOP",
		"vm:force_stop":               "ALLOW_VM_FORCE_STOP",
		"vm:reset":                    "ALLOW_VM_RESET",
		"plugin:add":                  "ALLOW_PLUGIN_INSTALL",
		"plugin:remove":               "ALLOW_PLUGIN_REMOVE",
		"plugin:install":              "ALLOW_PLUGIN_INSTALL",
		"plugin:install_language":     "ALLOW_PLUGIN_INSTALL",
		"rclone:delete_remote":        "ALLOW_RCLONE_OPERATIONS",
		"setting:configure_ups":       "ALLOW_SETTING_UPDATES",
		"setting:update_ssh":          "ALLOW_SSH_UPDATE",
		"setting:update_system_time":  "ALLOW_TIME_UPDATE",
		"connect:sign_in":             "ALLOW_CONNECT_ACTIONS",
		"connect:sign_out":            "ALLOW_CONNECT_ACTIONS",
		"connect:pair_device":         "ALLOW_CONNECT_ACTIONS",
		"connect:remove_device":       "ALLOW_CONNECT_ACTIONS",
		"notification:delete":         "ALLOW_NOTIFICATION_DELETE",
		"notification:delete_archived":"ALLOW_NOTIFICATION_DELETE",
		"key:create":                  "ALLOW_API_KEY_CREATE",
		"key:delete":                  "ALLOW_API_KEY_DELETE",
		"disk:flash_backup":           "ALLOW_FLASH_BACKUP",
		"onboarding:reset":            "ALLOW_ONBOARDING_ACTIONS",
		"onboarding:create_internal_boot_pool": "ALLOW_ONBOARDING_ACTIONS",
	}
	if t, ok := toggles[action+":"+subaction]; ok {
		return t
	}
	if action == "array" {
		return "ALLOW_ARRAY_ACTIONS"
	}
	if action == "docker" {
		return "ALLOW_DOCKER_ACTIONS"
	}
	if action == "vm" {
		return "ALLOW_VM_ACTIONS"
	}
	return "ALLOW_DESTRUCTIVE"
}

func isToggleEnabled(cfg *config.Config, action, subaction string) bool {
	if cfg.AllowDestructive {
		return true
	}
	key := action + ":" + subaction
	switch key {
	case "array:stop": return cfg.AllowArrayStop
	case "array:start": return cfg.AllowArrayStart
	case "array:add_disk": return cfg.AllowArrayAddDisk
	case "array:remove_disk": return cfg.AllowArrayRemoveDisk
	case "array:clear_disk_stats": return cfg.AllowArrayClearStats
	case "docker:stop": return cfg.AllowContainerStop
	case "docker:remove_container": return cfg.AllowContainerRemove
	case "docker:restart": return cfg.AllowContainerRestart
	case "vm:stop": return cfg.AllowVmStop
	case "vm:force_stop": return cfg.AllowVmForceStop
	case "vm:reset": return cfg.AllowVmReset
	case "plugin:add": return cfg.AllowPluginInstall
	case "plugin:remove": return cfg.AllowPluginRemove
	case "plugin:install": return cfg.AllowPluginInstall
	case "plugin:install_language": return cfg.AllowPluginInstall
	case "rclone:delete_remote": return cfg.AllowRcloneOperations
	case "setting:configure_ups": return cfg.AllowSettingUpdates
	case "setting:update_ssh": return cfg.AllowSshUpdate
	case "setting:update_system_time": return cfg.AllowTimeUpdate
	case "connect:sign_in": return cfg.AllowConnectActions
	case "connect:sign_out": return cfg.AllowConnectActions
	case "connect:pair_device": return cfg.AllowConnectActions
	case "connect:remove_device": return cfg.AllowConnectActions
	case "notification:delete": return cfg.AllowNotificationDelete
	case "notification:delete_archived": return cfg.AllowNotificationDelete
	case "key:create": return cfg.AllowApiKeyCreate
	case "key:delete": return cfg.AllowApiKeyDelete
	case "disk:flash_backup": return cfg.AllowFlashBackup
	case "onboarding:reset": return cfg.AllowOnboardingActions
	case "onboarding:create_internal_boot_pool": return cfg.AllowOnboardingActions
	}
	if action == "array" {
		return cfg.AllowArrayActions
	}
	if action == "docker" {
		return cfg.AllowDockerActions
	}
	if action == "vm" {
		return cfg.AllowVmActions
	}
	return false
}

var destructiveActions = map[string]bool{
	"array:stop": true,
	"array:start": true,
	"array:add_disk": true,
	"array:remove_disk": true,
	"array:clear_disk_stats": true,
	"docker:stop": true,
	"docker:remove_container": true,
	"docker:restart": true,
	"docker:reset_template_mappings": true,
	"docker:delete_entries": true,
	"vm:stop": true,
	"vm:force_stop": true,
	"vm:reset": true,
	"notification:delete": true,
	"notification:delete_archived": true,
	"key:delete": true,
	"plugin:add": true,
	"plugin:remove": true,
	"plugin:install": true,
	"plugin:install_language": true,
	"rclone:delete_remote": true,
	"setting:configure_ups": true,
	"setting:update_ssh": true,
	"setting:update_system_time": true,
	"connect:sign_in": true,
	"connect:sign_out": true,
	"connect:pair_device": true,
	"connect:remove_device": true,
	"onboarding:reset": true,
	"onboarding:create_internal_boot_pool": true,
	"disk:flash_backup": true,
}
