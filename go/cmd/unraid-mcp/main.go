package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/logger"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

const version = "0.0.1"

func main() {
	configPath := flag.String("config", "/boot/config/plugins/unraid-mcp/config.cfg", "Path to config file")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("unraid-mcp-server v%s\n", version)
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

	logLevel := os.Getenv("UNRAID_MCP_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.Init(logLevel, "json")

	logger.Get().Infof("Starting unraid-mcp-server v%s (Go)", version)
	logger.Get().Infof("Transport: %s, Read-only: %v", cfg.Transport, cfg.ReadOnly)

	server := mcp.NewServer(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Get().Info("Shutting down...")
		cancel()
	}()

	switch cfg.Transport {
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
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
	return nil
}
