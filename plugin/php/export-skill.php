<?php
// Download a skill as a SKILL.md file (valid for Kilo Code, Claude Code,
// Cursor, Gemini CLI, and other SKILL.md-compatible IDEs).

$base = '/boot/config/plugins/unraid-agent/skills';
$name = $_GET['name'] ?? '';
$source = $_GET['source'] ?? 'default';

if (!preg_match('/^[a-z][a-z0-9-]{0,63}$/', $name) || !in_array($source, ['default', 'custom'], true)) {
    http_response_code(400);
    exit;
}

$file = "$base/$source/$name.md";
$real = realpath($file);
if ($real === false || strpos($real, "$base/") !== 0 || !is_file($real)) {
    http_response_code(404);
    exit;
}

header('Content-Type: text/markdown; charset=utf-8');
header('Content-Disposition: attachment; filename="SKILL.md"');
header('Content-Length: ' . filesize($real));
readfile($real);
