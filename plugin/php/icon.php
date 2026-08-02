<?php
// VM icon proxy: serves a VM's icon.png from /boot/config/domains/<name>/
// after validating the VM name against the live list and confining the
// resolved path to the whitelisted directories.

define('UA_NO_OUTPUT', true);
include '/usr/local/emhttp/plugins/unraid-agent/php/common.php';

$type = $_GET['type'] ?? '';
$name = $_GET['name'] ?? '';

if ($type !== 'vm' || $name === '' || strlen($name) > 128 || strpos($name, '..') !== false || strpos($name, '/') !== false) {
    http_response_code(400);
    exit;
}

// Validate against the live VM list
$valid = false;
foreach (ua_list_vms() as $vm) {
    if (($vm['name'] ?? '') === $name) {
        $valid = true;
        break;
    }
}
if (!$valid) {
    http_response_code(404);
    exit;
}

$whitelist = ['/boot/config/domains/', '/boot/config/plugins/dynamix.vm.manager/templates/images/'];
$candidates = [
    "/boot/config/domains/$name/icon.png",
    "/boot/config/plugins/dynamix.vm.manager/templates/images/$name.png",
];

foreach ($candidates as $path) {
    $real = realpath($path);
    if ($real === false || !is_file($real)) {
        continue;
    }
    $ok = false;
    foreach ($whitelist as $dir) {
        if (strpos($real, $dir) === 0) {
            $ok = true;
            break;
        }
    }
    if (!$ok) {
        continue;
    }
    header('Content-Type: image/png');
    header('Cache-Control: private, max-age=300');
    readfile($real);
    exit;
}

http_response_code(404);
