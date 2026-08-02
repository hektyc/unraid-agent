<?php
// Read a skills/memory file's content. POST + CSRF. The file must resolve
// inside the plugin's skills/ or memory/ directories on flash.

define('UA_NO_OUTPUT', true);
include '/usr/local/emhttp/plugins/unraid-agent/php/common.php';

header('Content-Type: application/json');

$base = '/boot/config/plugins/unraid-agent';
$file = $_POST['file'] ?? '';
$real = realpath($file);
if ($real === false || !is_file($real) ||
    (strpos($real, "$base/skills/") !== 0 && strpos($real, "$base/memory/") !== 0)) {
    ua_ep_log('read-content', "FAIL invalid file: " . basename($file));
    http_response_code(400);
    echo json_encode(['ok' => false, 'error' => 'Invalid file']);
    exit;
}
ua_ep_log('read-content', basename($real) . ' ok');

echo json_encode(['ok' => true, 'content' => file_get_contents($real)]);
