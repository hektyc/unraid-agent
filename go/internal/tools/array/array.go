package array

import (

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_status",
		Description: "Get parity check progress. Read-only.",
		Query:       `query { parity { status speed errors elapsed remaining } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_history",
		Description: "Get past parity results. Read-only.",
		Query:       `query { parity { history { date duration errors status } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_assignable_disks",
		Description: "Get disks not yet in array. Read-only.",
		Query:       `query { array { assignableDisks { id device name size type } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_start",
		Description: "Start a parity check.",
		Query:       `mutation StartParityCheck($correct: Boolean) { startParityCheck(correct: $correct) { status } }`,
		Params: map[string]string{
			"correct": "boolean", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_pause",
		Description: "Pause a parity check.",
		Query:       `mutation { pauseParityCheck { status } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_resume",
		Description: "Resume a parity check.",
		Query:       `mutation { resumeParityCheck { status } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_cancel",
		Description: "Cancel a parity check.",
		Query:       `mutation { cancelParityCheck { status } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_start",
		Description: "Start the Unraid array.",
		Query:       `mutation { startArray { status } }`,
		Params: map[string]string{
			"confirm": "boolean", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_stop",
		Description: "Stop the Unraid array.",
		Query:       `mutation { stopArray { status } }`,
		Params: map[string]string{
			"confirm": "boolean", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_add_disk",
		Description: "Add a disk to the array.",
		Query:       `mutation AddDisk($diskId: ID!, $slot: Int) { addDisk(diskId: $diskId, slot: $slot) { status } }`,
		Params: map[string]string{
			"disk_id": "string", // required=true
			"slot": "number", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_remove_disk",
		Description: "Remove a disk from the array.",
		Query:       `mutation RemoveDisk($diskId: ID!) { removeDisk(diskId: $diskId) { status } }`,
		Params: map[string]string{
			"disk_id": "string", // required=true
			"confirm": "boolean", // required=true
		},
	})
}