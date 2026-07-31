package notification

import (

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_overview",
		Description: "Get notification counts. Read-only.",
		Query:       `query { notification { overview { unread archived } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_list",
		Description: "Get paginated notifications. Read-only.",
		Query:       `query NotificationList($lt: String, $imp: String, $off: Int, $lim: Int) { notification { list(listType: $lt, importance: $imp, offset: $off, limit: $lim) { id title importance } } }`,
		Params: map[string]string{
			"list_type": "string", // required=false
			"importance": "string", // required=false
			"offset": "number", // required=false
			"limit": "number", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_create",
		Description: "Create a notification.",
		Query:       `mutation CreateNotification($title: String!, $subject: String!, $description: String!, $importance: String!) { createNotification(title: $title, subject: $subject, description: $description, importance: $importance) { id } }`,
		Params: map[string]string{
			"title": "string", // required=true
			"subject": "string", // required=true
			"description": "string", // required=true
			"importance": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_archive",
		Description: "Archive a notification.",
		Query:       `mutation ArchiveNotification($id: ID!) { archiveNotification(id: $id) { status } }`,
		Params: map[string]string{
			"notification_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_delete",
		Description: "Delete a notification.",
		Query:       `mutation DeleteNotification($id: ID!, $type: String!) { deleteNotification(id: $id, type: $type) { status } }`,
		Params: map[string]string{
			"notification_id": "string", // required=true
			"notification_type": "string", // required=true
			"confirm": "boolean", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_delete_archived",
		Description: "Delete all archived.",
		Query:       `mutation { deleteArchivedNotifications { status } }`,
		Params: map[string]string{
			"confirm": "boolean", // required=true
		},
	})
}