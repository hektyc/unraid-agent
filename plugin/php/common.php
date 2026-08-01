<?php
// Shared config helper for all unRAID Agent settings pages.
// Included by each .page tab file.

$plugin = "unraid-agent";
$config_file = "/boot/config/plugins/$plugin/config.cfg";

$config = [];
if (file_exists($config_file)) {
    $parsed = parse_ini_file($config_file, false, INI_SCANNER_RAW);
    if ($parsed !== false) {
        $config = $parsed;
    }
}

$cfg = function($key, $default = '') use ($config) {
    return isset($config[$key]) ? $config[$key] : $default;
};

$running = false;
$pidfile = "/boot/config/plugins/$plugin/server.pid";
if (file_exists($pidfile)) {
    $pid = trim(file_get_contents($pidfile));
    if ($pid !== '') {
        exec("kill -0 $pid 2>/dev/null", $output, $rc);
        if ($rc === 0) {
            $running = true;
        }
    }
}
