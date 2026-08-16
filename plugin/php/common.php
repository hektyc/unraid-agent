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

/* --- Sub-tab navigation --- */
.ua-subtabs { display:flex; flex-wrap:wrap; border-bottom:1px solid #555; margin-bottom:16px; }
.ua-subtab { padding:10px 22px; font-size:14px; cursor:pointer; border-bottom:2px solid transparent; user-select:none; }
.ua-subtab:hover { opacity:0.8; }
.ua-subtab.active { border-bottom-color:#ff8c2f; font-weight:bold; }
.ua-subtab .fa { font-size:17px; margin-right:7px; vertical-align:-1px; }
.ua-pane { display:none; }
.ua-pane.active { display:block; }

/* Content page uses ac-pane instead of ua-pane for content panes */
.ac-pane { display:none; }
.ac-pane.active { display:block; }

/* --- Card grid --- */
.ua-grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(230px,1fr)); gap:14px; margin-top:12px; }
.ua-card { border:1px solid rgba(128,128,128,0.35); border-radius:10px; padding:12px 14px; cursor:pointer;
           background:rgba(128,128,128,0.06); transition:border-color .15s, transform .1s; }
.ua-card:hover { border-color:#ff8c2f; transform:translateY(-1px); }
.ua-card-top { display:flex; justify-content:space-between; align-items:center; margin-bottom:8px; }
.ua-card-icon { width:40px; height:40px; object-fit:contain; border-radius:6px; }
.ua-card-icon-fa { font-size:36px; line-height:40px; }
.ua-dot { display:inline-block; width:11px; height:11px; border-radius:50%; }
.ua-card-name { font-weight:bold; font-size:14px; margin-bottom:2px; word-break:break-all; }
.ua-card-image { font-size:11px; opacity:0.65; word-break:break-all; margin-bottom:4px; }
.ua-card-status { font-size:11px; opacity:0.8; margin-bottom:8px; }
.ua-card-chips { display:flex; flex-wrap:wrap; gap:4px; }
.ua-chip { font-size:10px; border:1px solid; border-radius:4px; padding:1px 5px; }
.ua-chip-global { border-color:#95a5a6; color:#95a5a6; }

/* --- Modals --- */
.ua-modal-overlay { position:fixed; inset:0; background:rgba(0,0,0,0.55); z-index:9999;
                    display:flex; align-items:center; justify-content:center; }
.ua-modal { background:#222; color:#eee; border:1px solid #555; border-radius:10px;
            width:480px; max-width:92vw; max-height:85vh; overflow-y:auto; padding:18px 20px; }
.ua-modal-head { display:flex; align-items:center; gap:12px; margin-bottom:14px; }
.ua-modal-head img { width:44px; height:44px; object-fit:contain; }
.ua-modal-head .fa { font-size:40px; }
.ua-modal-name { font-size:16px; font-weight:bold; word-break:break-all; }
.ua-perm-row { display:flex; justify-content:center; padding:3px 0; }
.ua-perm-pair { display:flex; align-items:center; gap:14px; padding:7px 2px 9px;
                border-bottom:1px solid rgba(128,128,128,0.25); }
.ua-perm-row:last-of-type .ua-perm-pair { border-bottom:none; }
.ua-perm-label { width:80px; text-align:right; font-size:13px; flex:none; }
.ua-perm-select { width:auto !important; min-width:150px; max-width:180px; flex:none; }
.ua-modal-foot { display:flex; justify-content:flex-end; gap:10px; margin-top:16px; }
.ua-btn { padding:6px 16px; border-radius:5px; border:1px solid #777; background:#333; color:#eee; cursor:pointer; }
.ua-btn:hover { border-color:#ff8c2f; }
.ua-btn-primary { background:#ff8c2f; border-color:#ff8c2f; color:#1a1a1a; font-weight:bold; }
.ua-modal-err { color:#e74c3c; font-size:12px; margin-top:8px; min-height:14px; }

/* --- Content page specific styles --- */
.ac-toolbar { margin:4px 0 12px; }
.ac-btn-new { padding:6px 16px; border-radius:5px; border:1px solid #ff8c2f; background:#ff8c2f; color:#1a1a1a; font-weight:bold; cursor:pointer; }
.ac-btn-new:hover { opacity:0.9; }
.ac-badge { font-size:10px; border:1px solid; border-radius:4px; padding:1px 6px; }
.ac-badge-default { border-color:#4a90d9; color:#4a90d9; }
.ac-badge-custom { border-color:#2ecc71; color:#2ecc71; }
.ac-scope { margin-top:18px; }
.ac-scope-title { font-weight:bold; font-size:14px; margin-bottom:8px; text-transform:capitalize; }
.ac-modal { background:#222; color:#eee; border:1px solid #555; border-radius:10px;
            width:720px; max-width:94vw; max-height:88vh; overflow-y:auto; padding:18px 20px; }
.ac-modal textarea { width:100%; box-sizing:border-box; background:#1b1b1b; color:#ddd;
                     border:1px solid #555; border-radius:6px; padding:10px; font-family:monospace; font-size:12px; }
.ac-modal input[type=text] { background:#1b1b1b; color:#ddd; border:1px solid #555; border-radius:6px; padding:6px 10px; }
.ac-field { margin-bottom:12px; }
.ac-field label { display:block; font-size:12px; opacity:0.75; margin-bottom:4px; }
.ac-name-input { width:280px; }
.ac-desc-input { width:100%; box-sizing:border-box; }

/* --- Tool Access sub-tabs (within Permissions) --- */
.ua-tools-wrap { width:45%; min-width:460px; margin:0 auto; }
.ua-domain { border:1px solid rgba(128,128,128,0.3); border-radius:8px; margin-bottom:8px; }
.ua-domain-head { display:flex; justify-content:space-between; align-items:center; padding:8px 12px; cursor:pointer; user-select:none; }
.ua-domain-head:hover { background:rgba(128,128,128,0.08); }
.ua-domain-title { font-weight:bold; font-size:13px; }
.ua-domain-count { font-size:11px; opacity:0.65; margin-left:8px; }
.ua-domain-caret { margin-right:8px; opacity:0.7; }
.ua-domain-tools { display:none; padding:4px 12px 10px; border-top:1px solid rgba(128,128,128,0.2); }
.ua-tool-row { display:flex; justify-content:space-between; align-items:center; padding:4px 0; border-bottom:1px dashed rgba(128,128,128,0.15); }
.ua-tool-row:last-child { border-bottom:none; }
.ua-tool-name { font-family:monospace; font-size:12px; }
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

if (!function_exists('ua_ep_log')) {

// Append one line to the endpoint request log (WebGUI AJAX diagnostics).
// Rotates at 256KB. Agents can read it via the agent_endpoint_log MCP tool;
// users can see the tail on the Advanced tab under Diagnostics.
function ua_ep_log($action, $detail = '') {
    $dir = '/boot/config/plugins/unraid-agent/logs';
    if (!is_dir($dir)) {
        @mkdir($dir, 0755, true);
    }
    $file = "$dir/endpoints.log";
    if (is_file($file) && filesize($file) > 262144) {
        @rename($file, $file . '.1');
    }
    $detail = trim(preg_replace('/\s+/', ' ', $detail));
    if (strlen($detail) > 300) {
        $detail = substr($detail, 0, 300) . '…';
    }
    $line = date('Y-m-d H:i:s') . ' ' . $action . ($detail !== '' ? ' ' . $detail : '') . "\n";
    @file_put_contents($file, $line, FILE_APPEND | LOCK_EX);
}

}
