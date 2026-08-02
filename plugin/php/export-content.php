<?php
// Download agent content: a skill as SKILL.md (IDE-compatible) or a memory
// entry as markdown. Read-only; names are strictly validated and the file
// must resolve inside the plugin's skills/ or memory/ directories.

$base = '/boot/config/plugins/unraid-agent';
$kind = $_GET['kind'] ?? 'skill';
$name = $_GET['name'] ?? '';

if (!preg_match('/^[a-z][a-z0-9-]{0,63}$/', $name)) {
    http_response_code(400);
    exit;
}

if ($kind === 'skill') {
    $source = $_GET['source'] ?? 'defaults';
    if ($source === 'default') {
        $source = 'defaults'; // badge label vs directory name
    }
    if (!in_array($source, ['defaults', 'custom'], true)) {
        http_response_code(400);
        exit;
    }
    $file = "$base/skills/$source/$name.md";
    $prefix = "$base/skills/";
    $download = "$name.SKILL.md";
} elseif ($kind === 'memory') {
    $scope = $_GET['scope'] ?? '';
    if (!preg_match('/^[a-zA-Z0-9][a-zA-Z0-9_-]{0,63}$/', $scope)) {
        http_response_code(400);
        exit;
    }
    $file = "$base/memory/$scope/$name.md";
    $prefix = "$base/memory/";
    $download = "$name.md";
} else {
    http_response_code(400);
    exit;
}

$real = realpath($file);
if ($real === false || strpos($real, $prefix) !== 0 || !is_file($real)) {
    http_response_code(404);
    exit;
}

header('Content-Type: text/markdown; charset=utf-8');
header('Content-Disposition: attachment; filename="' . $download . '"');
header('Content-Length: ' . filesize($real));
readfile($real);
