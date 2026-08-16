# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to Calendar Versioning (YYYY.MM.DD), the Unraid plugin ecosystem convention.

## [Unreleased]

### Fixed
- Root cause of raw HTML rendering: Unraid's markdown parser processes `.page` file content before PHP execution. Inline `<style>` blocks in the HTML body are not recognized as block-level HTML by the parser, causing subsequent HTML (including `<div class="ua-pane" data-uatab="tools">`) to be rendered as visible text instead of DOM elements. All CSS moved to `common.php`'s PHP `echo` `<style>` block, which executes after markdown parsing.
- Replaced nested sub→sub tabs in Permissions > Tool Access with a clean accordion layout (no nested `<div class="ua-pane" data-uatab="tg-...">` elements). Uses `ua-acc-section`/`ua-acc-header`/`ua-acc-body` with a single form, eliminating the nested HTML structure that triggered the markdown parser issue.

## [0.0.1] - 2026-07-23

### Added
- Initial release of Unraid MCP Server
- Core MCP server with stdio, streamable-http, and SSE transports
- Safety-first architecture with read-only mode and granular action toggles
- Comprehensive tool coverage: array, disk, docker, vm, notification, key, plugin, rclone, setting, connect, customization, oidc, onboarding, user, live, help
- WebSocket subscription system for live telemetry
- Express HTTP server with health endpoint and bearer auth
- GraphQL client with retry logic and redacted logging
- Graceful shutdown on SIGINT/SIGTERM
- Comprehensive test suite (config, guards, client, tools)
- CI/CD with lint, typecheck, build, and test jobs
- Automated release workflow with version bump, changelog update, git tag, GitHub Release
- Documentation including client connection guide for Claude Desktop, Claude Code, Codex, Kilo, OpenCode, VS Code, and Cursor
