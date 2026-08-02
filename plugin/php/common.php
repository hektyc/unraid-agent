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

// Shared plugin UI styles (all unRAID Agent pages include this file).
// Constrain selects — Unraid's default form CSS stretches them full width.
// Guards: the tab system can include this file twice per request (parent
// page body + active tab body), so everything must be include-safe.
if (!defined('UA_NO_OUTPUT') && !defined('UA_STYLES_PRINTED')) {
    define('UA_STYLES_PRINTED', true);
    echo '<style>
form select { width:auto !important; min-width:140px; }
</style>';
}

// ---------------------------------------------------------------------------
// GraphQL + entity helpers (used by the Permissions pages and API endpoints)
// ---------------------------------------------------------------------------

if (!function_exists('ua_graphql')) {

function ua_graphql($query, $vars = []) {
    global $config;
    $apiUrl = isset($config['UNRAID_API_URL']) ? trim($config['UNRAID_API_URL'], "\"' ") : '';
    $apiKey = isset($config['UNRAID_API_KEY']) ? trim($config['UNRAID_API_KEY'], "\"' ") : '';
    if ($apiUrl === '' || $apiKey === '') {
        return null;
    }
    $verifySsl = ($config['UNRAID_VERIFY_SSL'] ?? 'true') === 'true';
    $allowInsecure = ($config['UNRAID_ALLOW_INSECURE_TLS'] ?? 'false') === 'true';
    $insecure = (!$verifySsl || $allowInsecure);

    $payload = json_encode(['query' => $query, 'variables' => (object)$vars]);
    $cmd = 'curl -s -m 10 ' . ($insecure ? '-k ' : '')
         . '-X POST ' . escapeshellarg($apiUrl)
         . ' -H ' . escapeshellarg('Content-Type: application/json')
         . ' -H ' . escapeshellarg('x-api-key: ' . $apiKey)
         . ' -d ' . escapeshellarg($payload);
    $out = shell_exec($cmd);
    if (!$out) {
        return null;
    }
    $res = json_decode($out, true);
    return isset($res['data']) ? $res['data'] : null;
}

function ua_list_containers() {
    $d = ua_graphql('query { docker { containers { id names image state status labels } } }');
    return ($d && isset($d['docker']['containers'])) ? $d['docker']['containers'] : [];
}

function ua_list_vms() {
    $d = ua_graphql('query { vms { domains { id name state } } }');
    return ($d && isset($d['vms']['domains'])) ? $d['vms']['domains'] : [];
}

function ua_list_plugins() {
    $d = ua_graphql('query { plugins { name version hasApiModule hasCliModule } }');
    return ($d && isset($d['plugins'])) ? $d['plugins'] : [];
}

// Installed plugins: ground truth is the .plg files in /boot/config/plugins/.
// The API's `plugins` query only returns API-module plugins (a subset), so we
// parse each .plg's ENTITY declarations for name/version like plgman does,
// then enrich API/CLI badges from the API-module list where names match.
function ua_list_installed_plugins() {
    $out = [];
    foreach (glob('/boot/config/plugins/*.plg') ?: [] as $file) {
        $content = @file_get_contents($file);
        if ($content === false) {
            continue;
        }
        $name = '';
        $version = '';
        if (preg_match('/<!ENTITY\s+name\s+"([^"]+)"/', $content, $m)) {
            $name = $m[1];
        }
        if (preg_match('/<!ENTITY\s+version\s+"([^"]+)"/', $content, $m)) {
            $version = $m[1];
        }
        if ($name === '' && preg_match('/<PLUGIN\s[^>]*name="([^"]+)"/', $content, $m)) {
            $name = $m[1];
        }
        if ($version === '' && preg_match('/<PLUGIN\s[^>]*version="([^"]+)"/', $content, $m)) {
            $version = $m[1];
        }
        if ($name === '' || strpos($name, '&') !== false) {
            $name = basename($file, '.plg');
        }
        if (strpos($version, '&') !== false) {
            $version = '';
        }
        $out[] = ['name' => $name, 'version' => $version, 'hasApiModule' => false, 'hasCliModule' => false];
    }
    foreach (ua_list_plugins() as $ap) {
        foreach ($out as &$o) {
            if ($o['name'] === ($ap['name'] ?? '')) {
                $o['hasApiModule'] = !empty($ap['hasApiModule']);
                $o['hasCliModule'] = !empty($ap['hasCliModule']);
            }
        }
    }
    unset($o);
    usort($out, function($a, $b) { return strcasecmp($a['name'], $b['name']); });
    return $out;
}

function ua_perms_path() {
    return '/boot/config/plugins/unraid-agent/perms.json';
}

function ua_load_perms() {
    $f = ua_perms_path();
    $empty = ['containers' => [], 'vms' => [], 'plugins' => []];
    if (!file_exists($f)) {
        return $empty;
    }
    $p = json_decode(file_get_contents($f), true);
    if (!is_array($p)) {
        return $empty;
    }
    return array_merge($empty, $p);
}

// Container icon URL from docker labels (net.unraid.docker.icon), or ''.
function ua_container_icon($container) {
    $labels = $container['labels'] ?? [];
    if (is_string($labels)) {
        $decoded = json_decode($labels, true);
        if (is_array($decoded)) {
            $labels = $decoded;
        }
    }
    if (is_array($labels) && !empty($labels['net.unraid.docker.icon'])) {
        return $labels['net.unraid.docker.icon'];
    }
    return '';
}

// VM icon: served via the local icon proxy (validates name + whitelists path).
function ua_vm_icon_url($name) {
    $candidates = [
        "/boot/config/domains/$name/icon.png",
        "/boot/config/plugins/dynamix.vm.manager/templates/images/$name.png",
    ];
    foreach ($candidates as $p) {
        if (is_file($p)) {
            return '/plugins/unraid-agent/php/icon.php?type=vm&name=' . rawurlencode($name);
        }
    }
    return '';
}

// Plugin icon: plgman convention <plugin>/<plugin>.png, via the icon proxy.
// Plugins stash icons in a few places, so check the common candidates and
// fall back to the first PNG found in the plugin dir or images/.
function ua_plugin_icon_url($name) {
    $candidates = [
        "/usr/local/emhttp/plugins/$name/$name.png",
        "/usr/local/emhttp/plugins/$name/icon.png",
        "/usr/local/emhttp/plugins/$name/images/$name.png",
        "/usr/local/emhttp/plugins/$name/images/icon.png",
        "/boot/config/plugins/$name/$name.png",
    ];
    foreach ($candidates as $p) {
        if (is_file($p)) {
            return '/plugins/unraid-agent/php/icon.php?type=plugin&name=' . rawurlencode($name);
        }
    }
    foreach ([glob("/usr/local/emhttp/plugins/$name/*.png") ?: [], glob("/usr/local/emhttp/plugins/$name/images/*.png") ?: []] as $g) {
        if (count($g)) {
            return '/plugins/unraid-agent/php/icon.php?type=plugin&name=' . rawurlencode($name);
        }
    }
    return '';
}

} // end function_exists guard
