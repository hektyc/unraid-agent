#!/usr/bin/env node

const fs = require("node:fs");
const path = require("node:path");

const packageJsonPath = path.join(process.cwd(), "package.json");
const newVersion = process.argv[2];

if (!newVersion) {
  console.error("Usage: node scripts/bump-version.js <new_version>");
  process.exit(1);
}

function bumpVersion() {
  const pkg = JSON.parse(fs.readFileSync(packageJsonPath, "utf-8"));
  pkg.version = newVersion;
  fs.writeFileSync(packageJsonPath, JSON.stringify(pkg, null, 2) + "\n", "utf-8");
  console.log(`Updated package.json to v${newVersion}`);
}

function bumpChangelog() {
  const changelogPath = path.join(process.cwd(), "CHANGELOG.md");
  const today = new Date().toISOString().split("T")[0]!;
  const header = `## [${newVersion}] - ${today}\n\n### Added\n- Release v${newVersion}\n\n`;

  if (fs.existsSync(changelogPath)) {
    const content = fs.readFileSync(changelogPath, "utf-8");
    fs.writeFileSync(changelogPath, header + content, "utf-8");
  } else {
    const initialContent = `# Changelog\n\nAll notable changes to this project will be documented in this file.\n\n${header}`;
    fs.writeFileSync(changelogPath, initialContent, "utf-8");
  }
  console.log(`Updated CHANGELOG.md for v${newVersion}`);
}

bumpVersion();
bumpChangelog();

console.log("Version bump complete.");
