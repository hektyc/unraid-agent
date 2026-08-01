package array

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_status",
		Description: "Get parity check progress. Read-only.",
		Query:       `query { array { state parityCheckStatus { date duration speed status errors progress correcting paused running } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_history",
		Description: "Get past parity results. Read-only.",
		Query:       `query { parityHistory { date duration speed status errors correcting } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_assignable_disks",
		Description: "Get disks not yet in array. Read-only.",
		Query:       `query { assignableDisks { id device name vendor size type interfaceType smartStatus temperature serialNum isSpinning } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_start",
		Description: "Start the Unraid array.",
		Query:       `mutation { array { setState(input: {desiredState: START}) { state } } }`,
		Params:      map[string]string{},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_stop",
		Description: "Stop the Unraid array.",
		Query:       `mutation { array { setState(input: {desiredState: STOP}) { state } } }`,
		Params:      map[string]string{},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_start",
		Description: "Start a parity check.",
		Query:       `mutation($correct: Boolean) { parityCheck { start(correct: $correct) } }`,
		Params: map[string]string{
			"correct": "boolean",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_pause",
		Description: "Pause a parity check.",
		Query:       `mutation { parityCheck { pause } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_resume",
		Description: "Resume a parity check.",
		Query:       `mutation { parityCheck { resume } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_parity_cancel",
		Description: "Cancel a parity check.",
		Query:       `mutation { parityCheck { cancel } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_add_disk",
		Description: "Add a disk to the array.",
		Query:       `mutation($disk_id: PrefixedID!, $slot: Int) { array { addDiskToArray(input: {id: $disk_id, slot: $slot}) { state } } }`,
		Params: map[string]string{
			"disk_id": "string",
			"slot":    "number",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "array_remove_disk",
		Description: "Remove a disk from the array.",
		Query:       `mutation($disk_id: PrefixedID!) { array { removeDiskFromArray(input: {id: $disk_id}) { state } } }`,
		Params: map[string]string{
			"disk_id": "string",
		},
	})
}
