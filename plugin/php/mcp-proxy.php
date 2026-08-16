<?php
// Reverse proxy: forwards HTTP requests from the Unraid WebGUI (port 443)
// to the unraid-agent MCP server on localhost (port 6970). This lets any
// reverse proxy in front of Unraid (Pangolin, nginx, Cloudflare, etc.)
// reach the MCP server without additional proxy configuration — the
// WebGUI path /plugins/unraid-agent/mcp.php is already forwarded.

// Suppress Unraid's normal page chrome (styles, etc.)
define('UA_NO_OUTPUT', true);
include '/usr/local/emhttp/plugins/unraid-agent/php/common.php';

$mcpPort = isset($config['UNRAID_MCP_PORT']) ? trim($config['UNRAID_MCP_PORT']) : '6970';
if (!ctype_digit($mcpPort) || (int)$mcpPort < 1 || (int)$mcpPort > 65535) {
    $mcpPort = '6970';
}

// Bearer token auth: pass through if configured
$bearer = isset($config['UNRAID_MCP_BEARER_TOKEN']) ? trim($config['UNRAID_MCP_BEARER_TOKEN']) : '';
$disableAuth = (isset($config['UNRAID_MCP_DISABLE_HTTP_AUTH']) && $config['UNRAID_MCP_DISABLE_HTTP_AUTH'] === 'true');

$target = "http://127.0.0.1:" . $mcpPort . $_SERVER['REQUEST_URI'];

// Read request body
$body = '';
$contentType = '';
if (in_array($_SERVER['REQUEST_METHOD'], ['POST', 'PUT', 'PATCH', 'DELETE'])) {
    $contentType = $_SERVER['CONTENT_TYPE'] ?? '';
    $body = file_get_contents('php://input');
}

// Build headers for the upstream request
$headers = [];

// Forward important headers
$forwardHeaders = [
    'Content-Type', 'Authorization', 'Accept', 'Mcp-Session-Id',
    'Last-Event-ID', 'X-Forwarded-Proto', 'X-Forwarded-Host',
    'X-Forwarded-Port', 'User-Agent', 'Cookie'
];
foreach ($forwardHeaders as $h) {
    $key = 'HTTP_' . strtoupper(str_replace('-', '_', $h));
    if (isset($_SERVER[$key])) {
        $headers[] = $h . ': ' . $_SERVER[$key];
    }
}

// Add Host header
$headers[] = 'Host: 127.0.0.1:' . $mcpPort;

// Set a reasonable timeout
$timeout = 30;
if (isset($_SERVER['HTTP_ACCEPT']) &&
    strpos($_SERVER['HTTP_ACCEPT'], 'text/event-stream') !== false) {
    // SSE connections need a long timeout and unbuffered output
    $timeout = 0; // no timeout for SSE
}

// Use cURL to forward the request
$ch = curl_init($target);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, false);
curl_setopt($ch, CURLOPT_HEADER, false);
curl_setopt($ch, CURLOPT_FOLLOWLOCATION, false);
curl_setopt($ch, CURLOPT_FAILONERROR, false);
curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 5);
curl_setopt($ch, CURLOPT_TIMEOUT, $timeout);
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $_SERVER['REQUEST_METHOD']);

if ($body !== '') {
    curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
}

if (!empty($headers)) {
    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
}

// For SSE, we need to stream the response
if ($timeout === 0) {
    curl_setopt($ch, CURLOPT_WRITEFUNCTION, function($curl, $data) {
        echo $data;
        return strlen($data);
    });
}

// Also stream headers
$responseHeaders = [];
curl_setopt($ch, CURLOPT_HEADERFUNCTION, function($curl, $header) use (&$responseHeaders) {
    $len = strlen($header);
    $parts = explode(':', $header, 2);
    if (count($parts) === 2) {
        $responseHeaders[trim($parts[0])] = trim($parts[1]);
    }
    return $len;
});

// Disable output buffering for streaming responses
if ($timeout === 0) {
    if (ob_get_level() > 0) {
        ob_end_clean();
    }
    // Don't set content-type here — let the upstream server set it
}

$exec = curl_exec($ch);
$curlErr = curl_error($ch);
$httpCode = curl_getinfo($ch, CURL_HTTP_CODE);
$contentType = curl_getinfo($ch, CURL_CONTENT_TYPE);

if ($curlErr && $timeout !== 0) {
    http_response_code(502);
    header('Content-Type: text/plain');
    echo 'Bad Gateway: ' . $curlErr;
    curl_close($ch);
    exit;
}

curl_close($ch);

// Output response headers
foreach ($responseHeaders as $name => $value) {
    // Skip hop-by-hop headers and some PHP-incompatible ones
    $lower = strtolower($name);
    if (in_array($lower, ['connection', 'transfer-encoding', 'content-length', 'content-encoding', 'content-type'])) {
        continue;
    }
    header($name . ': ' . $value);
}

// Set content type from upstream
if ($contentType) {
    header('Content-Type: ' . $contentType);
}

http_response_code($httpCode);

// For non-SSE, output the response body
if ($timeout !== 0 && $exec) {
    echo $exec;
}
