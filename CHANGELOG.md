# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to Calendar Versioning (YYYY.MM.DD), the Unraid plugin ecosystem convention.

## [Unreleased]

### Fixed
- Nested tabs (Permissions, Content) now work correctly: inline CSS `<style>` blocks restored in both child pages with the missing `.ua-subtab` styling added to the Content page
- Inline `<script>` blocks restored in both child pages (event handlers bound directly to elements which are present on AJAX load)
- Removed CSS/JS file verification from `.plg` install script — the external CSS/JS files are optional fallbacks only, not required for core UI
- Removed external `ua-script.js` `<script src>` from parent page to prevent duplicate event handlers conflicting with child page inline JS

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
