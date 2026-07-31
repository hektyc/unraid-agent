import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("notification_overview", "Get notification counts. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { notification { overview { unread archived } } }") })),
    make("notification_list", "Get paginated notifications. Read-only.", { type:"object", properties:{ list_type:{type:"string"}, importance:{type:"string"}, offset:{type:"number"}, limit:{type:"number"} } },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query NotificationList($lt: String, $imp: String, $off: Int, $lim: Int) { notification { list(listType: $lt, importance: $imp, offset: $off, limit: $lim) { id title importance } } }", p) })),
    make("notification_create", "Create a notification.", { type:"object", properties:{ title:{type:"string"}, subject:{type:"string"}, description:{type:"string"}, importance:{type:"string"} }, required:["title","subject","description","importance"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation CreateNotification($title: String!, $subject: String!, $description: String!, $importance: String!) { createNotification(title: $title, subject: $subject, description: $description, importance: $importance) { id } }", p) })),
    make("notification_notify_if_unique", "Create notification if unique.", { type:"object", properties:{ title:{type:"string"}, subject:{type:"string"}, description:{type:"string"}, importance:{type:"string"} }, required:["title","subject","description","importance"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation CreateNotificationIfUnique($title: String!, $subject: String!, $description: String!, $importance: String!) { createNotificationIfUnique(title: $title, subject: $subject, description: $description, importance: $importance) { id created } }", p) })),
    make("notification_archive", "Archive a notification.", { type:"object", properties:{ notification_id:{type:"string"} }, required:["notification_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation ArchiveNotification($id: ID!) { archiveNotification(id: $id) { status } }", p) })),
    make("notification_mark_unread", "Move notification to unread.", { type:"object", properties:{ notification_id:{type:"string"} }, required:["notification_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UnarchiveNotification($id: ID!) { unarchiveNotification(id: $id) { status } }", p) })),
    make("notification_recalculate", "Recalculate overview counts.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { recalculateNotificationOverview { status } }") })),
    make("notification_archive_all", "Archive all unread.", { type:"object", properties:{ importance:{type:"string"} } },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation ArchiveAllUnread($imp: String) { archiveAllUnread(importance: $imp) { status } }", p) })),
    make("notification_archive_many", "Archive specific notifications.", { type:"object", properties:{ notification_ids:{type:"array"} }, required:["notification_ids"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation ArchiveMany($ids: [ID!]!) { archiveMany(ids: $ids) { status } }", p) })),
    make("notification_unarchive_many", "Unarchive specific notifications.", { type:"object", properties:{ notification_ids:{type:"array"} }, required:["notification_ids"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UnarchiveMany($ids: [ID!]!) { unarchiveMany(ids: $ids) { status } }", p) })),
    make("notification_unarchive_all", "Unarchive all.", { type:"object", properties:{ importance:{type:"string"} } },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UnarchiveAll($imp: String) { unarchiveAll(importance: $imp) { status } }", p) })),
    make("notification_delete", "Permanently delete notification. Destructive.", { type:"object", properties:{ notification_id:{type:"string"}, notification_type:{type:"string"}, confirm:{type:"boolean"} }, required:["notification_id","notification_type","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation DeleteNotification($id: ID!, $type: String!) { deleteNotification(id: $id, type: $type) { status } }", p) })),
    make("notification_delete_archived", "Delete all archived. Destructive.", { type:"object", properties:{ confirm:{type:"boolean"} }, required:["confirm"] },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { deleteArchivedNotifications { status } }") })),
  ];
}
