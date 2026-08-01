package config

import (
	"fmt"
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
	AllowVmForceStop      bool
	AllowVmReset          bool
	AllowVmStop           bool
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
}

func Load() (*Config, error) {
	cfg := &Config{
		APIURL:    getEnv("UNRAID_API_URL", ""),
		APIKey:    getEnv("UNRAID_API_KEY", ""),
		Transport: getEnv("TRANSPORT", "stdio"),
		Host:      getEnv("UNRAID_MCP_HOST", "127.0.0.1"),
		Port:      getEnvInt("UNRAID_MCP_PORT", 6970),
		ReadOnly:  getEnvBool("READ_ONLY", false),
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
	cfg.AllowVmForceStop = getEnvBool("ALLOW_VM_FORCE_STOP", false)
	cfg.AllowVmReset = getEnvBool("ALLOW_VM_RESET", false)
	cfg.AllowVmStop = getEnvBool("ALLOW_VM_STOP", false)
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

	if cfg.APIURL == "" || cfg.APIKey == "" {
		return nil, fmt.Errorf("UNRAID_API_URL and UNRAID_API_KEY must be set")
	}
	return cfg, nil
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
