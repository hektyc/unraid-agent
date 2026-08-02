// Package agentcontent implements the agent skills + memory engine.
//
// Layout on flash (persistent, survives updates):
//
//	/boot/config/plugins/unraid-agent/skills/defaults/   plugin-owned, refreshed from the embedded pack
//	/boot/config/plugins/unraid-agent/skills/custom/     user-owned, never touched
//	/boot/config/plugins/unraid-agent/memory/defaults/   generated server profile
//	/boot/config/plugins/unraid-agent/memory/<scope>/    per-client memory entries
package agentcontent

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

//go:embed skills instructions.md
var packFS embed.FS

const (
	defaultSkillsSubdir = "defaults"
	customSkillsSubdir  = "custom"
	defaultMemorySubdir = "defaults"
	maxEntryBytes       = 128 * 1024
)

var nameRE = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)

type Skill struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Source      string    `json:"source"` // "default" | "custom"
	Content     string    `json:"content,omitempty"`
	ModTime     time.Time `json:"modTime"`
}

type MemoryEntry struct {
	Name    string    `json:"name"`
	Scope   string    `json:"scope"`
	Content string    `json:"content,omitempty"`
	ModTime time.Time `json:"modTime"`
}

func SkillsDir(configDir string) string { return filepath.Join(configDir, "skills") }
func MemoryDir(configDir string) string { return filepath.Join(configDir, "memory") }

// Instructions returns the embedded initialize-handshake primer.
func Instructions() string {
	data, err := packFS.ReadFile("instructions.md")
	if err != nil {
		return ""
	}
	return string(data)
}

// SyncDefaults writes every embedded skill into <configDir>/skills/defaults/.
// The defaults dir is plugin-owned and refreshed whenever content changes;
// skills/custom/ is never touched. Idempotent and safe on every daemon start.
func SyncDefaults(configDir string) error {
	defDir := filepath.Join(SkillsDir(configDir), defaultSkillsSubdir)
	custDir := filepath.Join(SkillsDir(configDir), customSkillsSubdir)
	if err := os.MkdirAll(defDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(custDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(MemoryDir(configDir), defaultMemorySubdir), 0755); err != nil {
		return err
	}

	return fs.WalkDir(packFS, "skills", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := packFS.ReadFile(path)
		if err != nil {
			return err
		}
		target := filepath.Join(defDir, filepath.Base(path))
		existing, err := os.ReadFile(target)
		if err == nil && sha256hex(existing) == sha256hex(data) {
			return nil // unchanged
		}
		return os.WriteFile(target, data, 0644)
	})
}

// EnsureProfile generates memory/defaults/server-profile.md once from live
// introspection. Skips silently if it exists or the API is unreachable.
func EnsureProfile(configDir string, query func(q string) (map[string]interface{}, error)) {
	dir := filepath.Join(MemoryDir(configDir), defaultMemorySubdir)
	target := filepath.Join(dir, "server-profile.md")
	if _, err := os.Stat(target); err == nil {
		return
	}
	if query == nil {
		return
	}

	var b strings.Builder
	b.WriteString("# Server Profile (auto-generated)\n\n")
	b.WriteString("Facts about this Unraid server, captured by unraid-agent on first start.\n\n")

	if data, err := query(`query { info { os { platform distro release kernel arch hostname uptime } cpu { brand cores threads } machineId } metrics { memory { total } } }`); err == nil {
		if info, ok := data["info"].(map[string]interface{}); ok {
			b.WriteString("## System\n")
			if osi, ok := info["os"].(map[string]interface{}); ok {
				fmt.Fprintf(&b, "- OS: %v %v (kernel %v, %v)\n", osi["platform"], osi["release"], osi["kernel"], osi["arch"])
				fmt.Fprintf(&b, "- Hostname: %v\n", osi["hostname"])
			}
			if cpu, ok := info["cpu"].(map[string]interface{}); ok {
				fmt.Fprintf(&b, "- CPU: %v (%v cores / %v threads)\n", cpu["brand"], cpu["cores"], cpu["threads"])
			}
		}
		if m, ok := data["metrics"].(map[string]interface{}); ok {
			if mem, ok := m["memory"].(map[string]interface{}); ok {
				if t, ok := mem["total"].(float64); ok {
					fmt.Fprintf(&b, "- RAM: %.0f GiB\n", t/1073741824)
				}
			}
		}
		b.WriteString("\n")
	}

	if data, err := query(`query { array { state capacity { kilobytes { total used free } } parities { name status } disks { name status } caches { name status } } }`); err == nil {
		if arr, ok := data["array"].(map[string]interface{}); ok {
			b.WriteString("## Array\n")
			fmt.Fprintf(&b, "- State: %v\n", arr["state"])
			if cap, ok := arr["capacity"].(map[string]interface{}); ok {
				if kb, ok := cap["kilobytes"].(map[string]interface{}); ok {
					fmt.Fprintf(&b, "- Capacity: %v used / %v total (KiB)\n", kb["used"], kb["total"])
				}
			}
			writeDiskLine := func(key, label string) {
				if list, ok := arr[key].([]interface{}); ok {
					for _, d := range list {
						if dm, ok := d.(map[string]interface{}); ok {
							fmt.Fprintf(&b, "- %s %v: %v\n", label, dm["name"], dm["status"])
						}
					}
				}
			}
			writeDiskLine("parities", "Parity")
			writeDiskLine("disks", "Data")
			writeDiskLine("caches", "Cache")
			b.WriteString("\n")
		}
	}

	if data, err := query(`query { docker { containers { names state } } }`); err == nil {
		if docker, ok := data["docker"].(map[string]interface{}); ok {
			if list, ok := docker["containers"].([]interface{}); ok {
				fmt.Fprintf(&b, "## Docker (%d containers)\n", len(list))
				for _, c := range list {
					if cm, ok := c.(map[string]interface{}); ok {
						if names, ok := cm["names"].([]interface{}); ok && len(names) > 0 {
							fmt.Fprintf(&b, "- %v (%v)\n", names[0], cm["state"])
						}
					}
				}
				b.WriteString("\n")
			}
		}
	}

	if data, err := query(`query { network { accessUrls { type name ipv4 } } }`); err == nil {
		if net, ok := data["network"].(map[string]interface{}); ok {
			if urls, ok := net["accessUrls"].([]interface{}); ok {
				b.WriteString("## Access URLs\n")
				for _, u := range urls {
					if um, ok := u.(map[string]interface{}); ok {
						fmt.Fprintf(&b, "- %v: %v %v\n", um["type"], um["name"], um["ipv4"])
					}
				}
			}
		}
	}

	_ = os.MkdirAll(dir, 0755)
	_ = os.WriteFile(target, []byte(b.String()), 0644)
}

func sha256hex(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// ParseFrontmatter extracts name and description from SKILL.md content.
func ParseFrontmatter(content string) (name, description string) {
	if !strings.HasPrefix(content, "---") {
		return "", ""
	}
	rest := strings.TrimPrefix(content, "---")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return "", ""
	}
	for _, line := range strings.Split(rest[:end], "\n") {
		if v, ok := strings.CutPrefix(strings.TrimSpace(line), "name:"); ok {
			name = strings.TrimSpace(v)
		}
		if v, ok := strings.CutPrefix(strings.TrimSpace(line), "description:"); ok {
			description = strings.TrimSpace(v)
		}
	}
	return name, description
}

func readSkillDir(dir, source string, out *[]Skill) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		name, desc := ParseFrontmatter(string(data))
		if name == "" {
			name = strings.TrimSuffix(e.Name(), ".md")
		}
		info, _ := e.Info()
		var mt time.Time
		if info != nil {
			mt = info.ModTime()
		}
		*out = append(*out, Skill{Name: name, Description: desc, Source: source, ModTime: mt})
	}
}

// ListSkills merges defaults + custom (custom wins on name collisions).
func ListSkills(configDir string) []Skill {
	var out []Skill
	readSkillDir(filepath.Join(SkillsDir(configDir), defaultSkillsSubdir), "default", &out)
	readSkillDir(filepath.Join(SkillsDir(configDir), customSkillsSubdir), "custom", &out)
	seen := map[string]int{}
	var merged []Skill
	for _, sk := range out {
		if i, dup := seen[sk.Name]; dup {
			if sk.Source == "custom" {
				merged[i] = sk
			}
			continue
		}
		seen[sk.Name] = len(merged)
		merged = append(merged, sk)
	}
	sort.Slice(merged, func(i, j int) bool { return merged[i].Name < merged[j].Name })
	return merged
}

// GetSkill returns one skill with content (custom wins over default).
func GetSkill(configDir, name string) (*Skill, error) {
	for _, source := range []string{customSkillsSubdir, defaultSkillsSubdir} {
		dir := filepath.Join(SkillsDir(configDir), source)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				continue
			}
			n, desc := ParseFrontmatter(string(data))
			if n == "" {
				n = strings.TrimSuffix(e.Name(), ".md")
			}
			if n == name {
				return &Skill{Name: n, Description: desc, Source: source, Content: string(data)}, nil
			}
		}
	}
	return nil, fmt.Errorf("skill not found: %s", name)
}

// WriteSkill creates or updates a custom skill. Returns a validation error
// for bad names/oversized content; the file is written atomically.
func WriteSkill(configDir, name, description, content string) error {
	if !nameRE.MatchString(name) {
		return fmt.Errorf("invalid skill name %q (lowercase letters, digits, hyphens)", name)
	}
	if description == "" || len(description) > 1024 {
		return fmt.Errorf("description required, max 1024 characters")
	}
	if len(content) > maxEntryBytes {
		return fmt.Errorf("content exceeds %d bytes", maxEntryBytes)
	}
	full := content
	if !strings.HasPrefix(content, "---") {
		full = fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\n%s", name, description, content)
	}
	dir := filepath.Join(SkillsDir(configDir), customSkillsSubdir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	tmp := filepath.Join(dir, name+".md.tmp")
	if err := os.WriteFile(tmp, []byte(full), 0644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, name+".md"))
}

// DeleteSkill removes a custom skill. Default-pack skills cannot be deleted
// (they refresh from the embedded pack anyway).
func DeleteSkill(configDir, name string) error {
	if !nameRE.MatchString(name) {
		return fmt.Errorf("invalid skill name %q", name)
	}
	target := filepath.Join(SkillsDir(configDir), customSkillsSubdir, name+".md")
	if _, err := os.Stat(target); err != nil {
		if _, derr := os.Stat(filepath.Join(SkillsDir(configDir), defaultSkillsSubdir, name+".md")); derr == nil {
			return fmt.Errorf("cannot delete default skill %q (it ships with the plugin)", name)
		}
		return fmt.Errorf("skill not found: %s", name)
	}
	return os.Remove(target)
}

var scopeRE = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,63}$`)

// SanitizeScope normalizes a client name (from initialize clientInfo) into a
// safe memory scope.
func SanitizeScope(clientName string) string {
	s := strings.ToLower(strings.TrimSpace(clientName))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "default"
	}
	if len(out) > 64 {
		out = out[:64]
	}
	return out
}

func memoryScopeDirs(configDir, scope string) []string {
	dirs := []string{filepath.Join(MemoryDir(configDir), scope)}
	if scope != defaultMemorySubdir {
		dirs = append(dirs, filepath.Join(MemoryDir(configDir), defaultMemorySubdir))
	}
	return dirs
}

// ListMemory returns entries visible to a scope: its own dir plus defaults.
func ListMemory(configDir, scope string) []MemoryEntry {
	var out []MemoryEntry
	for _, dir := range memoryScopeDirs(configDir, scope) {
		sc := filepath.Base(dir)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			info, _ := e.Info()
			var mt time.Time
			if info != nil {
				mt = info.ModTime()
			}
			out = append(out, MemoryEntry{
				Name:    strings.TrimSuffix(e.Name(), ".md"),
				Scope:   sc,
				ModTime: mt,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Scope != out[j].Scope {
			return out[i].Scope < out[j].Scope
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// GetMemory returns one entry's content (scope dir first, then defaults).
func GetMemory(configDir, scope, name string) (*MemoryEntry, error) {
	if !nameRE.MatchString(name) && !strings.Contains(name, "-") {
		return nil, fmt.Errorf("invalid memory name %q", name)
	}
	for _, dir := range memoryScopeDirs(configDir, scope) {
		target := filepath.Join(dir, name+".md")
		data, err := os.ReadFile(target)
		if err == nil {
			return &MemoryEntry{Name: name, Scope: filepath.Base(dir), Content: string(data)}, nil
		}
	}
	return nil, fmt.Errorf("memory not found: %s", name)
}

// WriteMemory upserts an entry into the scope's own dir (never defaults).
func WriteMemory(configDir, scope, name, content string) error {
	if !scopeRE.MatchString(scope) {
		return fmt.Errorf("invalid scope %q", scope)
	}
	if scope == defaultMemorySubdir {
		return fmt.Errorf("the defaults scope is read-only")
	}
	if !nameRE.MatchString(name) {
		return fmt.Errorf("invalid memory name %q (lowercase letters, digits, hyphens)", name)
	}
	if len(content) > maxEntryBytes {
		return fmt.Errorf("content exceeds %d bytes", maxEntryBytes)
	}
	dir := filepath.Join(MemoryDir(configDir), scope)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	tmp := filepath.Join(dir, name+".md.tmp")
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, name+".md"))
}

// DeleteMemory removes an entry from the scope's own dir only.
func DeleteMemory(configDir, scope, name string) error {
	if scope == defaultMemorySubdir {
		return fmt.Errorf("the defaults scope is read-only")
	}
	if !nameRE.MatchString(name) {
		return fmt.Errorf("invalid memory name %q", name)
	}
	target := filepath.Join(MemoryDir(configDir), scope, name+".md")
	if _, err := os.Stat(target); err != nil {
		return fmt.Errorf("memory not found: %s", name)
	}
	return os.Remove(target)
}
