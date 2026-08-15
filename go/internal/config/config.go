package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	APIURL                string
	APIKey                string
	Transport             string
	Host                  string
	Port                  int
	MCPEndpoint           string
	ReadOnly              bool
	VerifySSL             bool
	AllowArrayStop        bool
	AllowArrayStart       bool
	AllowArrayAddDisk     bool
	AllowArrayRemoveDisk  bool
	AllowArrayClearStats  bool
	AllowContainerStop    bool
	AllowContainerRemove  bool
	AllowContainerRestart bool
	AllowContainerStart   bool
	AllowContainerPause   bool
	AllowContainerUnpause bool
	AllowContainerUpdate  bool
	AllowVmForceStop      bool
	AllowVmReset          bool
	AllowVmStop           bool
	AllowVmStart          bool
	AllowVmPause          bool
	AllowVmResume         bool
	AllowVmReboot         bool
	AllowSkillsWrite      bool
	AllowMemoryWrite      bool
	AnonymizeLogs         bool
	AllowPluginInstall    bool
	AllowPluginRemove     bool
	AllowSettingUpdates   bool
	AllowSshUpdate        bool
	AllowTimeUpdate       bool
	AllowNotificationDelete bool
	AllowApiKeyCreate     bool
	AllowApiKeyDelete     bool
	AllowFlashBackup      bool
	AllowRcloneOperations bool
	AllowConnectActions   bool
	AllowOnboardingActions bool
	AllowDockerActions    bool
	AllowVmActions        bool
	AllowArrayActions     bool
	AllowDestructive      bool
	BearerToken           string
	DisableHTTPAuth       bool
	AllowInsecureTLS      bool
	EnableStdio           bool

	// Domains maps tool domain name -> enabled. Disabled domains are not
	// registered, so their tools and schemas never reach clients (token
	// savings). Missing keys default to enabled.
	Domains map[string]bool
}

func Load() (*Config, error) {
	cfg := &Config{
		APIURL:    getEnv("UNRAID_API_URL", ""),
		APIKey:    getEnv("UNRAID_API_KEY", ""),
		Transport: getEnv("TRANSPORT", "streamable-http"),
		Host:      getEnv("UNRAID_MCP_HOST", "0.0.0.0"),
		Port:      getEnvInt("UNRAID_MCP_PORT", 6970),
		MCPEndpoint: getEnv("UNRAID_MCP_ENDPOINT", ""),
		ReadOnly:  getEnvBool("READ_ONLY", true),
		VerifySSL: getEnvBool("UNRAID_VERIFY_SSL", true),
	}

	cfg.AllowArrayStop = getEnvBool("ALLOW_ARRAY_STOP", false)
	cfg.AllowArrayStart = getEnvBool("ALLOW_ARRAY_START", false)
	cfg.AllowArrayAddDisk = getEnvBool("ALLOW_ARRAY_ADD_DISK", false)
	cfg.AllowArrayRemoveDisk = getEnvBool("ALLOW_ARRAY_REMOVE_DISK", false)
	cfg.AllowArrayClearStats = getEnvBool("ALLOW_ARRAY_CLEAR_STATS", false)
	cfg.AllowContainerStop = getEnvBool("ALLOW_CONTAINER_STOP", false)
	cfg.AllowContainerRemove = getEnvBool("ALLOW_CONTAINER_REMOVE", false)
	cfg.AllowContainerRestart = getEnvBool("ALLOW_CONTAINER_RESTART", false)
	cfg.AllowContainerStart = getEnvBool("ALLOW_CONTAINER_START", false)
	cfg.AllowContainerPause = getEnvBool("ALLOW_CONTAINER_PAUSE", false)
	cfg.AllowContainerUnpause = getEnvBool("ALLOW_CONTAINER_UNPAUSE", false)
	cfg.AllowContainerUpdate = getEnvBool("ALLOW_CONTAINER_UPDATE", false)
	cfg.AllowVmForceStop = getEnvBool("ALLOW_VM_FORCE_STOP", false)
	cfg.AllowVmReset = getEnvBool("ALLOW_VM_RESET", false)
	cfg.AllowVmStop = getEnvBool("ALLOW_VM_STOP", false)
	cfg.AllowVmStart = getEnvBool("ALLOW_VM_START", false)
	cfg.AllowVmPause = getEnvBool("ALLOW_VM_PAUSE", false)
	cfg.AllowVmResume = getEnvBool("ALLOW_VM_RESUME", false)
	cfg.AllowVmReboot = getEnvBool("ALLOW_VM_REBOOT", false)
	cfg.AllowSkillsWrite = getEnvBool("ALLOW_SKILLS_WRITE", false)
	cfg.AllowMemoryWrite = getEnvBool("ALLOW_MEMORY_WRITE", false)
	cfg.AnonymizeLogs = getEnvBool("ANONYMIZE_LOGS", false)

	cfg.Domains = map[string]bool{}
	for _, d := range ToolDomains {
		cfg.Domains[d] = getEnvBool("UNRAID_MCP_DOMAIN_"+strings.ToUpper(d), true)
	}
	cfg.AllowPluginInstall = getEnvBool("ALLOW_PLUGIN_INSTALL", false)
	cfg.AllowPluginRemove = getEnvBool("ALLOW_PLUGIN_REMOVE", false)
	cfg.AllowSettingUpdates = getEnvBool("ALLOW_SETTING_UPDATES", false)
	cfg.AllowSshUpdate = getEnvBool("ALLOW_SSH_UPDATE", false)
	cfg.AllowTimeUpdate = getEnvBool("ALLOW_TIME_UPDATE", false)
	cfg.AllowNotificationDelete = getEnvBool("ALLOW_NOTIFICATION_DELETE", false)
	cfg.AllowApiKeyCreate = getEnvBool("ALLOW_API_KEY_CREATE", false)
	cfg.AllowApiKeyDelete = getEnvBool("ALLOW_API_KEY_DELETE", false)
	cfg.AllowFlashBackup = getEnvBool("ALLOW_FLASH_BACKUP", false)
	cfg.AllowRcloneOperations = getEnvBool("ALLOW_RCLONE_OPERATIONS", false)
	cfg.AllowConnectActions = getEnvBool("ALLOW_CONNECT_ACTIONS", false)
	cfg.AllowOnboardingActions = getEnvBool("ALLOW_ONBOARDING_ACTIONS", false)
	cfg.AllowDockerActions = getEnvBool("ALLOW_DOCKER_ACTIONS", false)
	cfg.AllowVmActions = getEnvBool("ALLOW_VM_ACTIONS", false)
	cfg.AllowArrayActions = getEnvBool("ALLOW_ARRAY_ACTIONS", false)
	cfg.AllowDestructive = getEnvBool("ALLOW_DESTRUCTIVE", false)
	cfg.BearerToken = getEnv("UNRAID_MCP_BEARER_TOKEN", "")
	cfg.DisableHTTPAuth = getEnvBool("UNRAID_MCP_DISABLE_HTTP_AUTH", false)
	cfg.AllowInsecureTLS = getEnvBool("UNRAID_ALLOW_INSECURE_TLS", false)
	cfg.EnableStdio = getEnvBool("UNRAID_MCP_ENABLE_STDIO", false)

	return cfg, nil
}

// PermissionValue reports whether a named ALLOW_* toggle is enabled.
func (c *Config) PermissionValue(name string) bool {
	switch name {
	case "ALLOW_ARRAY_ACTIONS":
		return c.AllowArrayActions
	case "ALLOW_ARRAY_STOP":
		return c.AllowArrayStop
	case "ALLOW_ARRAY_START":
		return c.AllowArrayStart
	case "ALLOW_ARRAY_ADD_DISK":
		return c.AllowArrayAddDisk
	case "ALLOW_ARRAY_REMOVE_DISK":
		return c.AllowArrayRemoveDisk
	case "ALLOW_ARRAY_CLEAR_STATS":
		return c.AllowArrayClearStats
	case "ALLOW_CONTAINER_STOP":
		return c.AllowContainerStop
	case "ALLOW_CONTAINER_REMOVE":
		return c.AllowContainerRemove
	case "ALLOW_CONTAINER_RESTART":
		return c.AllowContainerRestart
	case "ALLOW_CONTAINER_START":
		return c.AllowContainerStart
	case "ALLOW_CONTAINER_PAUSE":
		return c.AllowContainerPause
	case "ALLOW_CONTAINER_UNPAUSE":
		return c.AllowContainerUnpause
	case "ALLOW_CONTAINER_UPDATE":
		return c.AllowContainerUpdate
	case "ALLOW_VM_FORCE_STOP":
		return c.AllowVmForceStop
	case "ALLOW_VM_RESET":
		return c.AllowVmReset
	case "ALLOW_VM_STOP":
		return c.AllowVmStop
	case "ALLOW_VM_START":
		return c.AllowVmStart
	case "ALLOW_VM_PAUSE":
		return c.AllowVmPause
	case "ALLOW_VM_RESUME":
		return c.AllowVmResume
	case "ALLOW_VM_REBOOT":
		return c.AllowVmReboot
	case "ALLOW_SKILLS_WRITE":
		return c.AllowSkillsWrite
	case "ALLOW_MEMORY_WRITE":
		return c.AllowMemoryWrite
	case "ALLOW_PLUGIN_INSTALL":
		return c.AllowPluginInstall
	case "ALLOW_PLUGIN_REMOVE":
		return c.AllowPluginRemove
	case "ALLOW_SETTING_UPDATES":
		return c.AllowSettingUpdates
	case "ALLOW_SSH_UPDATE":
		return c.AllowSshUpdate
	case "ALLOW_TIME_UPDATE":
		return c.AllowTimeUpdate
	case "ALLOW_NOTIFICATION_DELETE":
		return c.AllowNotificationDelete
	case "ALLOW_API_KEY_CREATE":
		return c.AllowApiKeyCreate
	case "ALLOW_API_KEY_DELETE":
		return c.AllowApiKeyDelete
	case "ALLOW_FLASH_BACKUP":
		return c.AllowFlashBackup
	case "ALLOW_RCLONE_OPERATIONS":
		return c.AllowRcloneOperations
	case "ALLOW_CONNECT_ACTIONS":
		return c.AllowConnectActions
	case "ALLOW_ONBOARDING_ACTIONS":
		return c.AllowOnboardingActions
	case "ALLOW_DOCKER_ACTIONS":
		return c.AllowDockerActions
	case "ALLOW_VM_ACTIONS":
		return c.AllowVmActions
	case "ALLOW_DESTRUCTIVE":
		return c.AllowDestructive
	}
	return false
}

// ToolDomains enumerates every registerable tool family. Each maps to a
// UNRAID_MCP_DOMAIN_<NAME> config toggle in Permissions -> Tool Access.
var ToolDomains = []string{
	"array", "connect", "customization", "docker", "health", "help", "key",
	"live", "logs", "notification", "oidc", "onboarding", "plugin", "rclone",
	"setting", "system", "user", "vm", "agentcontent",
}

// ToolDomainOf maps a tool name to its domain via prefix rules.
// Keep in sync with the Tool Access UI map in the plugin pages.
func ToolDomainOf(tool string) string {
	if tool == "plugin_logs" {
		return "logs"
	}
	for _, d := range ToolDomains {
		if strings.HasPrefix(tool, d+"_") {
			return d
		}
	}
	switch {
	case tool == "help_full":
		return "help"
	case tool == "user_me":
		return "user"
	case strings.HasPrefix(tool, "skills_"), strings.HasPrefix(tool, "memory_"), tool == "agent_endpoint_log":
		return "agentcontent"
	}
	return "system"
}

// IsToolEnabled resolves registration for a single tool:
// explicit per-tool override (UNRAID_MCP_TOOL_<NAME>=true/false) wins in
// both directions; otherwise the tool follows its domain toggle.
func (c *Config) IsToolEnabled(tool string) bool {
	key := "UNRAID_MCP_TOOL_" + strings.ToUpper(tool)
	v := os.Getenv(key)
	if v != "" && v != "domain" {
		return getEnvBool(key, true)
	}
	return c.IsDomainEnabled(ToolDomainOf(tool))
}

// IsDomainEnabled reports whether a tool domain is enabled (default true).
func (c *Config) IsDomainEnabled(name string) bool {
	if c.Domains == nil {
		return true
	}
	enabled, ok := c.Domains[name]
	if !ok {
		return true
	}
	return enabled
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	v = strings.ToLower(v)
	return !(v == "0" || v == "false" || v == "no" || v == "off")
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
