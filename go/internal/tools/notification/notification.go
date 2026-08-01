package notification

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_overview",
		Description: "Get notification counts. Read-only.",
		Query:       `query { notifications { overview { unread { info warning alert total } archive { info warning alert total } } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_list",
		Description: "Get paginated notifications. Read-only.",
		Query:       `query($offset: Int! = 0, $limit: Int! = 25, $list_type: NotificationType, $importance: NotificationImportance) { notifications { list(filter: {offset: $offset, limit: $limit, type: $list_type, importance: $importance}) { id title subject description importance type timestamp link } } }`,
		Params: map[string]string{
			"offset":     "number",
			"limit":      "number",
			"list_type":  "string",
			"importance": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_archive",
		Description: "Archive a notification.",
		Query:       `mutation($notification_id: PrefixedID!) { archiveNotification(id: $notification_id) { id title importance } }`,
		Params: map[string]string{
			"notification_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_create",
		Description: "Create a notification.",
		Query:       `mutation($title: String!, $subject: String!, $description: String!, $importance: NotificationImportance!) { createNotification(input: {title: $title, subject: $subject, description: $description, importance: $importance}) { id title } }`,
		Params: map[string]string{
			"title":       "string",
			"subject":     "string",
			"description": "string",
			"importance":  "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_delete",
		Description: "Delete a notification.",
		Query:       `mutation($notification_id: PrefixedID!, $notification_type: NotificationType!) { deleteNotification(id: $notification_id, type: $notification_type) { unread { total } } }`,
		Params: map[string]string{
			"notification_id":   "string",
			"notification_type": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "notification_delete_archived",
		Description: "Delete all archived.",
		Query:       `mutation { deleteArchivedNotifications { archive { total } } }`,
	})
}
