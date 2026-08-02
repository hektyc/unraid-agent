<?php
// VM icon proxy: serves a VM's icon.png from /boot/config/domains/<name>/
// after validating the VM name against the live list and confining the
// resolved path to the whitelisted directories.

define('UA_NO_OUTPUT', true);
include '/usr/local/emhttp/plugins/unraid-agent/php/common.php';

$type = $_GET['type'] ?? '';
$name = $_GET['name'] ?? '';

if (!in_array($type, ['vm', 'plugin'], true) || $name === '' || strlen($name) > 128 || strpos($name, '..') !== false || strpos($name, '/') !== false) {
    http_response_code(400);
    exit;
}

// Validate against the live list for the requested type
$valid = false;
if ($type === 'vm') {
    foreach (ua_list_vms() as $vm) {
        if (($vm['name'] ?? '') === $name) {
            $valid = true;
            break;
        }
    }
} else {
    foreach (ua_list_installed_plugins() as $p) {
        if (($p['name'] ?? '') === $name) {
            $valid = true;
            break;
        }
    }
}
if (!$valid) {
    http_response_code(404);
    exit;
}

if ($type === 'vm') {
    $whitelist = ['/boot/config/domains/', '/boot/config/plugins/dynamix.vm.manager/templates/images/'];
    $candidates = [
        "/boot/config/domains/$name/icon.png",
        "/boot/config/plugins/dynamix.vm.manager/templates/images/$name.png",
    ];
} else {
    $whitelist = ['/usr/local/emhttp/plugins/', '/boot/config/plugins/'];
    $candidates = [
        "/usr/local/emhttp/plugins/$name/$name.png",
        "/usr/local/emhttp/plugins/$name/icon.png",
        "/usr/local/emhttp/plugins/$name/images/$name.png",
        "/usr/local/emhttp/plugins/$name/images/icon.png",
        "/boot/config/plugins/$name/$name.png",
    ];
    foreach ([glob("/usr/local/emhttp/plugins/$name/*.png") ?: [],
              glob("/usr/local/emhttp/plugins/$name/images/*.png") ?: [],
              glob("/boot/config/plugins/$name/*.png") ?: []] as $g) {
        foreach ($g as $f) {
            $candidates[] = $f;
        }
    }
}

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
