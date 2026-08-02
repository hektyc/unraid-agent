package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hektyc/unraid-mcp-server/internal/agentcontent"
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/logger"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
	"github.com/hektyc/unraid-mcp-server/internal/tools/array"
	agenttools "github.com/hektyc/unraid-mcp-server/internal/tools/agentcontent"
	"github.com/hektyc/unraid-mcp-server/internal/tools/connect"
	"github.com/hektyc/unraid-mcp-server/internal/tools/customization"
	"github.com/hektyc/unraid-mcp-server/internal/tools/docker"
	"github.com/hektyc/unraid-mcp-server/internal/tools/health"
	"github.com/hektyc/unraid-mcp-server/internal/tools/help"
	"github.com/hektyc/unraid-mcp-server/internal/tools/key"
	"github.com/hektyc/unraid-mcp-server/internal/tools/live"
	"github.com/hektyc/unraid-mcp-server/internal/tools/notification"
	"github.com/hektyc/unraid-mcp-server/internal/tools/oidc"
	"github.com/hektyc/unraid-mcp-server/internal/tools/onboarding"
	"github.com/hektyc/unraid-mcp-server/internal/tools/plugin"
	"github.com/hektyc/unraid-mcp-server/internal/tools/rclone"
	"github.com/hektyc/unraid-mcp-server/internal/tools/setting"
	"github.com/hektyc/unraid-mcp-server/internal/tools/system"
	"github.com/hektyc/unraid-mcp-server/internal/tools/user"
	"github.com/hektyc/unraid-mcp-server/internal/tools/vm"
)

// version is injected at build time via -ldflags "-X main.version=$VERSION"
var version = "dev"

func main() {
	configPath := flag.String("config", "/boot/config/plugins/unraid-agent/config.cfg", "Path to config file")
	transportOverride := flag.String("transport", "", "Transport override (stdio|streamable-http)")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("unraid-agent %s\n", version)
		os.Exit(0)
	}

	if err := loadDotEnv(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %s\n", err.Error())
		fmt.Fprintf(os.Stderr, "Set UNRAID_API_URL and UNRAID_API_KEY in environment or config file.\n")
		os.Exit(1)
	}

	// Effective transport: -transport flag overrides config
	transport := cfg.Transport
	if *transportOverride != "" {
		transport = *transportOverride
	}

	// stdio is gated: the binary refuses to serve it unless explicitly
	// enabled in the plugin configuration (default: disabled).
	if transport == "stdio" && !cfg.EnableStdio {
		fmt.Fprintln(os.Stderr, "Configuration error: stdio transport is disabled in plugin configuration (UNRAID_MCP_ENABLE_STDIO=\"false\")")
		fmt.Fprintln(os.Stderr, "Enable it in Settings → unRAID Agent → Enable stdio Transport.")
		os.Exit(1)
	}

	logLevel := os.Getenv("UNRAID_MCP_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.Init(logLevel, "json")

	logger.Get().Infof("Starting unraid-agent %s (Go)", version)
	logger.Get().Infof("Transport: %s, Read-only: %v", transport, cfg.ReadOnly)

	server := mcp.NewServer(cfg)
	configDir := filepath.Dir(*configPath)
	server.PermsPath = filepath.Join(configDir, "perms.json")
	server.ConfigDir = configDir

	// Agent content: sync the embedded skill pack to flash, then generate
	// the server profile memory in the background (first start only).
	if err := agentcontent.SyncDefaults(configDir); err != nil {
		logger.Get().Infof("agent content sync: %v", err)
	}
	go agentcontent.EnsureProfile(configDir, func(q string) (map[string]interface{}, error) {
		return server.GraphQLQuery(context.Background(), q, nil)
	})
	array.RegisterTools(server, cfg)
	agenttools.RegisterTools(server, cfg)
	array.RegisterTools(server, cfg)
	connect.RegisterTools(server, cfg)
	customization.RegisterTools(server, cfg)
	docker.RegisterTools(server, cfg)
	health.RegisterTools(server, cfg)
	help.RegisterTools(server, cfg)
	key.RegisterTools(server, cfg)
	live.RegisterTools(server, cfg)
	notification.RegisterTools(server, cfg)
	oidc.RegisterTools(server, cfg)
	onboarding.RegisterTools(server, cfg)
	plugin.RegisterTools(server, cfg)
	rclone.RegisterTools(server, cfg)
	setting.RegisterTools(server, cfg)
	system.RegisterTools(server, cfg)
	user.RegisterTools(server, cfg)
	vm.RegisterTools(server, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Get().Info("Shutting down...")
		cancel()
	}()

	switch transport {
	case "streamable-http", "http":
		if err := server.ServeHTTP(ctx, cfg.Host, cfg.Port); err != nil {
			logger.Get().Fatal(err)
		}
	default:
		if err := server.ServeStdio(ctx); err != nil {
			logger.Get().Fatal(err)
		}
	}
}

func loadDotEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, `"'`)
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
	return nil
}
