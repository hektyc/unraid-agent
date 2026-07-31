import { type UnraidConfig } from "./types.js";
import { ReadOnlyError, ActionNotAllowedError } from "./errors.js";

export function enforceReadOnly(config: UnraidConfig, action: string, subaction: string): void {
  if (config.readOnly) {
    throw new ReadOnlyError(action, subaction);
  }
}

export function enforceActionToggle(
  config: UnraidConfig,
  action: string,
  subaction: string,
): void {
  if (!isDestructiveAction(action, subaction)) {
    return;
  }

  const toggleKey = getToggleKey(action, subaction);
  const isAllowed = checkDestructiveToggle(config, action, subaction);

  if (!isAllowed) {
    throw new ActionNotAllowedError(action, subaction, toggleKey);
  }
}

export function isDestructiveAction(action: string, subaction: string): boolean {
  const key = `${action}:${subaction}`;
  return DESTRUCTIVE_ACTIONS.has(key);
}

export function getToggleKey(action: string, subaction: string): string {
  const actionMap: Record<string, Record<string, string>> = {
    array: {
      stop: "ALLOW_ARRAY_STOP",
      start: "ALLOW_ARRAY_START",
      add_disk: "ALLOW_ARRAY_ADD_DISK",
      remove_disk: "ALLOW_ARRAY_REMOVE_DISK",
      clear_disk_stats: "ALLOW_ARRAY_CLEAR_STATS",
    },
    docker: {
      stop: "ALLOW_CONTAINER_STOP",
      remove_container: "ALLOW_CONTAINER_REMOVE",
      restart: "ALLOW_CONTAINER_RESTART",
    },
    vm: {
      stop: "ALLOW_VM_STOP",
      force_stop: "ALLOW_VM_FORCE_STOP",
      reset: "ALLOW_VM_RESET",
    },
    plugin: {
      add: "ALLOW_PLUGIN_INSTALL",
      remove: "ALLOW_PLUGIN_REMOVE",
      install: "ALLOW_PLUGIN_INSTALL",
      install_language: "ALLOW_PLUGIN_INSTALL",
    },
    rclone: {
      create_remote: "ALLOW_RCLONE_OPERATIONS",
      delete_remote: "ALLOW_RCLONE_OPERATIONS",
    },
    setting: {
      configure_ups: "ALLOW_SETTING_UPDATES",
      update_ssh: "ALLOW_SSH_UPDATE",
      update_system_time: "ALLOW_TIME_UPDATE",
    },
    connect: {
      sign_in: "ALLOW_CONNECT_ACTIONS",
      sign_out: "ALLOW_CONNECT_ACTIONS",
      pair_device: "ALLOW_CONNECT_ACTIONS",
      remove_device: "ALLOW_CONNECT_ACTIONS",
    },
    notification: {
      delete: "ALLOW_NOTIFICATION_DELETE",
      delete_archived: "ALLOW_NOTIFICATION_DELETE",
    },
    key: {
      create: "ALLOW_API_KEY_CREATE",
      delete: "ALLOW_API_KEY_DELETE",
    },
    disk: {
      flash_backup: "ALLOW_FLASH_BACKUP",
    },
    onboarding: {
      reset: "ALLOW_ONBOARDING_ACTIONS",
      create_internal_boot_pool: "ALLOW_ONBOARDING_ACTIONS",
    },
  };

  return actionMap[action]?.[subaction] || `ALLOW_${action.toUpperCase()}_${subaction.toUpperCase()}`;
}

export function checkDestructiveToggle(config: UnraidConfig, action: string, subaction: string): boolean {
  const key = `${action}:${subaction}`;
  if (config.safetyToggles.allowDestructive) {
    return true;
  }

  const toggleMap: Record<string, (s: typeof config.safetyToggles) => boolean> = {
    "array:stop": (s) => s.allowArrayStop,
    "array:start": (s) => s.allowArrayStart,
    "array:add_disk": (s) => s.allowArrayAddDisk,
    "array:remove_disk": (s) => s.allowArrayRemoveDisk,
    "array:clear_disk_stats": (s) => s.allowArrayClearStats,
    "docker:stop": (s) => s.allowContainerStop,
    "docker:remove_container": (s) => s.allowContainerRemove,
    "docker:restart": (s) => s.allowContainerRestart,
    "vm:stop": (s) => s.allowVmStop,
    "vm:force_stop": (s) => s.allowVmForceStop,
    "vm:reset": (s) => s.allowVmReset,
    "plugin:add": (s) => s.allowPluginInstall,
    "plugin:remove": (s) => s.allowPluginRemove,
    "plugin:install": (s) => s.allowPluginInstall,
    "plugin:install_language": (s) => s.allowPluginInstall,
    "rclone:create_remote": (s) => s.allowRcloneOperations,
    "rclone:delete_remote": (s) => s.allowRcloneOperations,
    "setting:configure_ups": (s) => s.allowSettingUpdates,
    "setting:update_ssh": (s) => s.allowSshUpdate,
    "setting:update_system_time": (s) => s.allowTimeUpdate,
    "connect:sign_in": (s) => s.allowConnectActions,
    "connect:sign_out": (s) => s.allowConnectActions,
    "connect:pair_device": (s) => s.allowConnectActions,
    "connect:remove_device": (s) => s.allowConnectActions,
    "notification:delete": (s) => s.allowNotificationDelete,
    "notification:delete_archived": (s) => s.allowNotificationDelete,
    "key:create": (s) => s.allowApiKeyCreate,
    "key:delete": (s) => s.allowApiKeyDelete,
    "disk:flash_backup": (s) => s.allowFlashBackup,
    "onboarding:reset": (s) => s.allowOnboardingActions,
    "onboarding:create_internal_boot_pool": (s) => s.allowOnboardingActions,
  };

  const checker = toggleMap[key];
  if (!checker) {
    return false;
  }
  return checker(config.safetyToggles);
}

const DESTRUCTIVE_ACTIONS: Set<string> = new Set([
  "array:stop",
  "array:start",
  "array:add_disk",
  "array:remove_disk",
  "array:clear_disk_stats",
  "docker:stop",
  "docker:remove_container",
  "docker:restart",
  "docker:reset_template_mappings",
  "docker:delete_entries",
  "vm:stop",
  "vm:force_stop",
  "vm:reset",
  "notification:delete",
  "notification:delete_archived",
  "key:delete",
  "plugin:add",
  "plugin:remove",
  "plugin:install",
  "plugin:install_language",
  "rclone:delete_remote",
  "setting:configure_ups",
  "setting:update_ssh",
  "setting:update_system_time",
  "connect:sign_in",
  "connect:sign_out",
  "connect:pair_device",
  "connect:remove_device",
  "onboarding:reset",
  "onboarding:create_internal_boot_pool",
  "disk:flash_backup",
]);

export function withSafetyChecks<T>(
  config: UnraidConfig,
  action: string,
  subaction: string,
  fn: () => Promise<T>,
): Promise<T> {
  enforceReadOnly(config, action, subaction);
  enforceActionToggle(config, action, subaction);
  return fn();
}
