<?php
// Save or delete agent content (skills and memory). POST + CSRF, whitelisted
// names, atomic writes. Skills always land in custom/; memory in its scope dir.

define('UA_NO_OUTPUT', true);
include '/usr/local/emhttp/plugins/unraid-agent/php/common.php';

header('Content-Type: application/json');

$base = '/boot/config/plugins/unraid-agent';
$kind = $_POST['kind'] ?? '';
$action = $_POST['action'] ?? 'save';
$name = trim($_POST['name'] ?? '');
$content = $_POST['content'] ?? '';

if (!in_array($kind, ['skill', 'memory'], true)) {
    http_response_code(400);
    echo json_encode(['ok' => false, 'error' => 'Invalid kind']);
    exit;
}

// Deletes are file-based and handled before name validation — the entry's
// real path is the authority, confined to the writable dirs:
// skills only in custom/, memory in any scope EXCEPT defaults/.
if ($action === 'delete') {
    $file = $_POST['file'] ?? '';
    $real = realpath($file);
    if ($file === '' || $real === false || !is_file($real)) {
        http_response_code(404);
        echo json_encode(['ok' => false, 'error' => 'Entry not found']);
        exit;
    }
    if ($kind === 'skill') {
        if (strpos($real, "$base/skills/custom/") !== 0) {
            http_response_code(400);
            echo json_encode(['ok' => false, 'error' => 'Default pack skills cannot be deleted']);
            exit;
        }
    } else {
        if (strpos($real, "$base/memory/") !== 0 || strpos($real, "$base/memory/defaults/") === 0) {
            http_response_code(400);
            echo json_encode(['ok' => false, 'error' => 'This memory entry cannot be deleted']);
            exit;
        }
    }
    ua_ep_log('save-content', "$kind delete " . basename($real) . ' ok');
    unlink($real);
    echo json_encode(['ok' => true, 'deleted' => basename($real, '.md')]);
    exit;
}

if (!preg_match('/^[a-z][a-z0-9-]{0,63}$/', $name)) {
    http_response_code(400);
    echo json_encode(['ok' => false, 'error' => 'Invalid name (lowercase letters, digits, hyphens)']);
    exit;
}

if ($kind === 'skill') {
    $dir = "$base/skills/custom";
    $desc = trim($_POST['description'] ?? '');
} else {
    $scope = preg_replace('/[^a-zA-Z0-9_-]/', '-', strtolower(trim($_POST['scope'] ?? 'custom')));
    $scope = trim($scope, '-');
    if ($scope === '' ) $scope = 'custom';
    if ($scope === 'defaults') {
        http_response_code(400);
        echo json_encode(['ok' => false, 'error' => 'The defaults scope is read-only']);
        exit;
    }
    $dir = "$base/memory/$scope";
}
if (!is_dir($dir)) {
    mkdir($dir, 0755, true);
}
$target = "$dir/$name.md";

if (strlen($content) > 65536) {
    http_response_code(400);
    echo json_encode(['ok' => false, 'error' => 'Content exceeds 64KB']);
    exit;
}

if ($kind === 'skill') {
    if ($desc === '' || strlen($desc) > 1024) {
        http_response_code(400);
        echo json_encode(['ok' => false, 'error' => 'Description required, max 1024 characters']);
        exit;
    }
    // Rebuild with proper SKILL.md frontmatter
    $body = $content;
    if (preg_match('/^---\s*\n(.*?)\n---\s*\n?(.*)$/s', $content, $m)) {
        $body = $m[2];
    }
    $out = "---\nname: $name\ndescription: $desc\n---\n\n" . ltrim($body);
} else {
    $out = $content;
}

ua_ep_log('save-content', "$kind $action $name");
$tmp = $target . '.tmp';
if (file_put_contents($tmp, $out) === false) {
    http_response_code(500);
    echo json_encode(['ok' => false, 'error' => 'Write failed']);
    exit;
}
chmod($tmp, 0644);
rename($tmp, $target);

echo json_encode(['ok' => true, 'saved' => $name]);
