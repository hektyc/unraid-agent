<?php
// Per-entity permission save endpoint.
// POST: csrf_token, entity_type (containers|vms), entity (name), perms (JSON
// object of action -> "inherit"|"allow"|"deny"). Atomic write, whitelisted.

define('UA_NO_OUTPUT', true);
include '/usr/local/emhttp/plugins/unraid-agent/php/common.php';

header('Content-Type: application/json');

// CSRF validation against var.ini
$type = $_POST['entity_type'] ?? '';
$entity = trim($_POST['entity'] ?? '');
$permsRaw = $_POST['perms'] ?? '{}';

if (!in_array($type, ['containers', 'vms', 'plugins'], true)) {
    http_response_code(400);
    echo json_encode(['ok' => false, 'error' => 'Invalid entity_type']);
    exit;
}
if ($entity === '' || strlen($entity) > 128) {
    http_response_code(400);
    echo json_encode(['ok' => false, 'error' => 'Invalid entity name']);
    exit;
}

$perms = json_decode($permsRaw, true);
if (!is_array($perms)) {
    http_response_code(400);
    echo json_encode(['ok' => false, 'error' => 'Invalid perms payload']);
    exit;
}

$allowedActionsMap = [
    'containers' => ['start', 'stop', 'restart', 'pause', 'unpause', 'remove', 'update'],
    'vms'        => ['start', 'stop', 'pause', 'resume', 'force_stop', 'reboot', 'reset'],
    'plugins'    => ['remove'],
];
$allowedActions = $allowedActionsMap[$type];
$allowedValues = ['inherit', 'global', 'allow', 'deny'];

// Self-removal is guarded at every layer: the unraid-agent plugin cannot
// carry a "remove: allow" override (the daemon also hard-blocks it).
if ($type === 'plugins' && $entity === 'unraid-agent') {
    http_response_code(400);
    echo json_encode(['ok' => false, 'error' => 'unraid-agent is protected and cannot take remove overrides']);
    exit;
}

// Validate the entity exists in the live list (prevents arbitrary key injection)
$live = ($type === 'containers') ? ua_list_containers() : (($type === 'vms') ? ua_list_vms() : ua_list_installed_plugins());
$liveNames = [];
foreach ($live as $item) {
    if ($type === 'containers') {
        $names = $item['names'] ?? [];
        if (isset($names[0])) {
            $liveNames[] = ltrim($names[0], '/');
        }
    } else {
        if (isset($item['name'])) {
            $liveNames[] = $item['name'];
        }
    }
}
if (!in_array($entity, $liveNames, true)) {
    http_response_code(400);
    echo json_encode(['ok' => false, 'error' => 'Unknown entity: ' . $entity]);
    exit;
}

// Whitelist-filter the submitted actions/values; "inherit"/"global" entries
// are dropped (absence of a key means use-global).
$clean = [];
foreach ($perms as $action => $value) {
    if (!in_array($action, $allowedActions, true)) {
        continue;
    }
    if (!in_array($value, $allowedValues, true)) {
        continue;
    }
    if ($value === 'inherit' || $value === 'global') {
        continue;
    }
    $clean[$action] = $value;
}

$all = ua_load_perms();
if (empty($clean)) {
    unset($all[$type][$entity]);
} else {
    $all[$type][$entity] = $clean;
}

$file = ua_perms_path();
$tmp = $file . '.tmp';
if (file_put_contents($tmp, json_encode($all, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES)) === false) {
    http_response_code(500);
    echo json_encode(['ok' => false, 'error' => 'Write failed']);
    exit;
}
chmod($tmp, 0600);
rename($tmp, $file);

echo json_encode(['ok' => true, 'entity' => $entity, 'perms' => $clean]);
