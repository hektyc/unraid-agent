---
name: dockerman-template-authoring
description: Author valid Unraid DockerMan template XML from an application repository or image reference — the full template schema, Config types, appdata conventions, network modes, and a step-by-step repo-to-template workflow. Use whenever the user asks to create, fix, or review a Docker template for Unraid.
---

# DockerMan Template Authoring

## What a DockerMan template is
An XML file Unraid's Docker Manager uses to present a container's configuration UI (ports, paths, variables, icon, description). Users add containers through it via the Docker page → Add Container. Templates are not compose files — they describe ONE container and its settings form.

## Canonical structure
```xml
<?xml version="1.0"?>
<Container version="2">
  <Name>myapp</Name>                          <!-- container name (lowercase, no spaces) -->
  <Repository>org/image:tag</Repository>      <!-- required -->
  <Registry>https://hub.docker.com/r/org/image</Registry>
  <Network>bridge</Network>                   <!-- bridge | host | <custom-bridge-name> -->
  <Privileged>false</Privileged>
  <Support>https://github.com/org/repo/issues</Support>
  <Project>https://github.com/org/repo</Project>
  <Overview>What the app does, 1-3 sentences.</Overview>
  <Category>Tools:</Category>                 <!-- Media: | Network: | Tools: | Productivity: etc. -->
  <WebUI>http://[IP]:[PORT:8080]/</WebUI>     <!-- [IP] and [PORT:x] are Unraid placeholders -->
  <Icon>https://example.com/icon.png</Icon>
  <Config Name="WebUI Port" Target="8080" Default="8080" Mode="tcp" Description="Container port for the web UI" Type="Port" Display="always" Required="true" Mask="false">8080</Config>
  <Config Name="Appdata" Target="/config" Default="/mnt/user/appdata/myapp" Mode="rw" Description="Persistent config/data" Type="Path" Display="always" Required="true" Mask="false">/mnt/user/appdata/myapp</Config>
  <Config Name="Timezone" Target="TZ" Default="Etc/UTC" Description="Timezone variable" Type="Variable" Display="always" Required="false" Mask="false">Etc/UTC</Config>
</Container>
```

## Config entry rules
- **Type="Port"**: container↔host port mapping. Target = container port, value = host port (Default matches Target unless there's a reason to differ). Mode = tcp (default) or udp.
- **Type="Path"**: bind mount. Target = in-container path, Default/value = host path — ALWAYS under /mnt/user/ (appdata convention: /mnt/user/appdata/<name>). Mode = rw or ro.
- **Type="Variable"**: environment variable. Mask="true" for secrets (renders as password field).
- **Type="Label"**: Docker label. **Type="Device"**: device mapping (e.g. /dev/dri for GPU).
- **Display**: always (shown in basic view), advanced (collapsed), hidden.
- **Required="true"** only when the container fails without it.

## Repo → template workflow
1. Read the repo: Dockerfile, docker-compose.yml, README. Extract: image name/tag, EXPOSE'd ports, VOLUME paths, ENV variables, and any required devices/privileged flags.
2. Map every port → Type="Port" Config. Every volume that holds state → Type="Path" under /mnt/user/appdata/<name>. Every meaningful ENV → Type="Variable" (secrets masked).
3. Choose Network: bridge (default, isolated NAT), host (app needs host IP/stack ports), or a custom bridge name (when the app talks to other containers by hostname).
4. Set WebUI with [IP]/[PORT:x] placeholders where x = the Target of the web port. Add Project/Support links, Category, and a square PNG/SVG icon URL.
5. Sanity-check: exactly one <Container version="2"> root, no duplicate Config Names, every Required field has a sensible Default, XML is well-formed (escape & as &amp; etc.).

## Gold-standard rules
- Never put credentials in plaintext Config values — use Mask="true" Variables and tell the user to fill them in.
- Never default paths outside /mnt/user/ — array disks (/mnt/diskX) are not for appdata.
- Ports below 1024 need Privileged or are a red flag for the wrong image.
- If the repo offers an official unraid template or CA listing, tell the user instead of reinventing it.
