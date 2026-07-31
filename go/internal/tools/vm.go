package vm

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_list",
		Description: "Get all VMs. Read-only.",
		Query:       `query { vm { list { id name state uuid } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_details",
		Description: "Get single VM details.",
		Query:       `query VmDetails($id: ID!) { vm { dom(id: $id) { id name state uuid } } }`,
		Params: map[string]string{
			"vm_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_start",
		Description: "Start a VM.",
		Query:       `mutation StartVm($id: ID!) { startVm(id: $id) { status } }`,
		Params: map[string]string{
			"vm_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_stop",
		Description: "Gracefully stop a VM.",
		Query:       `mutation StopVm($id: ID!) { stopVm(id: $id) { status } }`,
		Params: map[string]string{
			"vm_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_pause",
		Description: "Pause a VM.",
		Query:       `mutation PauseVm($id: ID!) { pauseVm(id: $id) { status } }`,
		Params: map[string]string{
			"vm_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_resume",
		Description: "Resume a VM.",
		Query:       `mutation ResumeVm($id: ID!) { resumeVm(id: $id) { status } }`,
		Params: map[string]string{
			"vm_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_reboot",
		Description: "Reboot a VM.",
		Query:       `mutation RebootVm($id: ID!) { rebootVm(id: $id) { status } }`,
		Params: map[string]string{
			"vm_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_force_stop",
		Description: "Hard power-off a VM.",
		Query:       `mutation ForceStopVm($id: ID!) { forceStopVm(id: $id) { status } }`,
		Params: map[string]string{
			"vm_id": "string", // required=true
			"confirm": "boolean", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "vm_reset",
		Description: "Hard reset a VM.",
		Query:       `mutation ResetVm($id: ID!) { resetVm(id: $id) { status } }`,
		Params: map[string]string{
			"vm_id": "string", // required=true
			"confirm": "boolean", // required=true
		},
	})
}