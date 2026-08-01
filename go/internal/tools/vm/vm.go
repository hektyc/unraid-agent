package vm

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_list",
		Description: "Get all VMs. Read-only.",
		Query:       `query { vms { domains { id name state } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_details",
		Description: "Get single VM details (returns all domains; match by id).",
		Query:       `query { vms { domains { id name state } } }`,
		Params: map[string]string{
			"vm_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_start",
		Description: "Start a VM.",
		Query:       `mutation($vm_id: PrefixedID!) { vm { start(id: $vm_id) } }`,
		Params: map[string]string{
			"vm_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_stop",
		Description: "Gracefully stop a VM.",
		Query:       `mutation($vm_id: PrefixedID!) { vm { stop(id: $vm_id) } }`,
		Params: map[string]string{
			"vm_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_pause",
		Description: "Pause a VM.",
		Query:       `mutation($vm_id: PrefixedID!) { vm { pause(id: $vm_id) } }`,
		Params: map[string]string{
			"vm_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_resume",
		Description: "Resume a VM.",
		Query:       `mutation($vm_id: PrefixedID!) { vm { resume(id: $vm_id) } }`,
		Params: map[string]string{
			"vm_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_force_stop",
		Description: "Hard power-off a VM.",
		Query:       `mutation($vm_id: PrefixedID!) { vm { forceStop(id: $vm_id) } }`,
		Params: map[string]string{
			"vm_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_reboot",
		Description: "Reboot a VM.",
		Query:       `mutation($vm_id: PrefixedID!) { vm { reboot(id: $vm_id) } }`,
		Params: map[string]string{
			"vm_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_reset",
		Description: "Hard reset a VM.",
		Query:       `mutation($vm_id: PrefixedID!) { vm { reset(id: $vm_id) } }`,
		Params: map[string]string{
			"vm_id": "string",
		},
	})
}
