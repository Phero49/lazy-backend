import {
  computed,
  ComputedRef,
  onBeforeMount,
  Ref,
  ref,
  watch,
} from "vue";
import { Notify, useQuasar } from "quasar";
import { GetProjectEnv, GetTableList } from "../wailsjs/go/main/App";
import { useAppContextStore } from "../src/store/appContext";
import { mapLabel } from "../src/utils/utills";

/**
 * A database connection definition loaded from the project environment.
 * Used by the shared schema singleton to identify and connect to a configured DB.
 */
export interface DBProfile {
  /** Internal connection name/key from the project config. */
  name: string;
  /** Friendly display label shown in the UI. */
  label: string;
  /** Raw database config values such as host, port, user, password, and database. */
  config: any;
}

/**
 * The public shape of the shared database schema state.
 * Hovering a destructured field from this object will show what each ref represents.
 */
export interface DatabaseSchemaState {
  /** All database profiles discovered in the current project environment. */
  configurations: Ref<DBProfile[]>;
  /** Map of connection name to loaded table names for the currently connected database. */
  activeConnections: Ref<Map<string, string[]>>;
  /** Name of the connection currently being established. */
  connecting: Ref<string | null>;
  /** True while fetching table metadata from the selected database. */
  loadingTables: Ref<boolean>;
  /** Text used to filter the visible table list. */
  tableFilter: Ref<string>;
  /** Name of the table currently being created via the popup form. */
  newTable: Ref<string>;
  /** All table names for the active database from the app schema cache. */
  allTables: ComputedRef<string[]>;
  /** Filtered subset of tables based on the search string. */
  filteredTables: ComputedRef<string[]>;
  /** Loads every database config from the project environment. */
  loadConfigurations: () => Promise<void>;
  /** Refreshes the table list for a specific connection and stores it in activeConnections. */
  refreshTables: (connectionName: string) => Promise<string[] | undefined>;
  /** Connects to a configured database and makes it the active connection. */
  connect: (profile: DBProfile) => Promise<void>;
  /** Connects to the first configured database if one exists, otherwise opens the database config screen. */
  connectDefault: () => Promise<void>;
  /** Switches the active connection if it already exists or connects if it is new. */
  selectDatabase: (name: string) => void;
  /** Removes a connection from the active map and updates the active connection selection. */
  disconnect: (name: string) => void;
}

type DatabaseSchemaInstance = DatabaseSchemaState;

/** Shared singleton instance used across all components to keep connection/table state in sync. */
let singleton: DatabaseSchemaInstance | null = null;

/**
 * Creates the shared database state object used to manage configs, active connections,
 * and table loading across the app.
 */
function createDatabaseSchema(): DatabaseSchemaState {
  const appContext = useAppContextStore();

  /** Configurations discovered in the active project environment. */
  const configurations = ref<DBProfile[]>([]);
  /** Stored table names keyed by connection name. */
  const activeConnections = ref<Map<string, string[]>>(new Map());
  /** Connection currently in progress. */
  const connecting = ref<string | null>(null);
  /** True while we are refreshing the table list. */
  const loadingTables = ref(false);
  /** Search text used for filtering database tables. */
  const tableFilter = ref("");
  /** Name entered in the create-table popup. */
  const newTable = ref("");

  /** Returns table names from the selected connection's schema cache. */
  const allTables = computed<string[]>(() => {
    const connectionName = appContext.activeConnectionName;
    if (!connectionName) return [];

    const schema = appContext.databaseSchemas[connectionName];
    return schema?.schema?.entities?.map((entity: any) => entity.name) || [];
  });

  /** Filters the active database tables based on the current search input. */
  const filteredTables = computed(() => {
    if (!appContext.activeConnectionName) return [];

    const tables = allTables.value || [];
    if (!tableFilter.value) return tables;

    const search = tableFilter.value.toLowerCase();
    return tables.filter((table: string) =>
      table.toLowerCase().includes(search),
    );
  });

  /** Loads database profiles for the active project and stores them in the shared config list. */
  async function loadConfigurations() {
    if (!appContext.currentProject) return;

    try {
      const env = await GetProjectEnv(appContext.currentProject);
      const dbList = env.database || [];
     const mappedDb = dbList.map((item: any) => {
        const value = Object.keys(item)[0];
        const label = mapLabel("database", item);

        return {
          name: value,
          label,
          config: item[value],
        };
      });
     
       configurations.value =  mappedDb
    } catch (error) {
      console.error("Failed to load configurations:", error);
    }
  }

  /** Loads and caches the table names for the given connection. */
  async function refreshTables(connectionName: string) {
 
    // Prefer using a previously-imported/saved schema stored in the app context.
    const saved = appContext.databaseSchemas?.[connectionName];
    if (saved && saved.schema && Array.isArray(saved.schema.entities)) {
      const entities = saved.schema.entities as any[];
      const tables = entities.map((e) => e.name);
      activeConnections.value.set(connectionName, tables);
      return tables;
    }

    // Fall back to querying the DB plugin if no saved schema exists.
    loadingTables.value = true;

  }

  /** Connects to a configured database and makes it the active connection. */
  async function connect(profile: DBProfile) {
    connecting.value = profile.label;
    try {
      appContext.setActiveConnection(profile.label);
      await refreshTables(profile.label);

      Notify.create({
        type: "positive",
        message: `tables loaded for ${profile.label}`,
        timeout: 1000,
        position: "bottom-right",
      });
    } catch (error: any) {
      Notify.create({
        type: "negative",
        message: ` failed to load tables  ${error}`,
      });
    } finally {
      connecting.value = null;
    }
  }

  /** Connects to the first configured database, or opens the config screen when no DB is defined. */
  async function connectDefault() {
    // Ensure configurations are loaded before deciding what to do.
    if (configurations.value.length === 0) {
      await loadConfigurations();
    }
    

    // If still no configurations, open the DB config tab so the user can add one.
    if (configurations.value.length === 0) {
      appContext.createNewTab(
        "/env/update/database",
        "create database",
        "bi-database-add",
      );
      return;
    }

    // If we already have an active connection that matches a known profile,
    // prefer reconnecting to it so the UI remains stable.
    const activeName = appContext.activeConnectionName;
    const activeProfile = configurations.value.find((c) => c.label === activeName);
    if (activeProfile) {

      await connect(activeProfile);
      return;
    }

    // Fallback: connect to the first configured profile.
    await connect(configurations.value[0]);
  }

  /** Selects a database by name. If it is already tracked, it becomes active; otherwise it connects. */
  function selectDatabase(name: string) {
    if (activeConnections.value.has(name)) {
      appContext.setActiveConnection(name);
      return;
    }

    const profile = configurations.value.find((config) => config.label === name);
    if (profile) {
      connect(profile);
    }
  }

  /** Removes a connection from the active list and updates the currently selected connection. */
  function disconnect(name: string) {
    activeConnections.value.delete(name);

    if (appContext.activeConnectionName === name) {
      const next =
        activeConnections.value.size > 0
          ? Array.from(activeConnections.value.keys())[0]
          : null;

      appContext.setActiveConnection(next);
    }
  }

  watch(
    () => appContext.currentProject,
    async () => {
      // Clear previously tracked runtime connections
      activeConnections.value.clear();
      appContext.setActiveConnection(null);

      // Load configurations for the newly selected project
      await loadConfigurations();

      // If the project already has saved database schemas, prefer those
      const schemaKeys = Object.keys(appContext.databaseSchemas || {});
      if (schemaKeys.length > 0) {
        const first = schemaKeys[0];
        // Populate activeConnections from the saved schema so the UI shows counts
        const entities = appContext.databaseSchemas[first]?.schema?.entities || [];
        activeConnections.value.set(first, entities.map((e: any) => e.label));
        appContext.setActiveConnection(first);
        return;
      }

      // Otherwise try to connect to the first configured profile (or open config)
      await connectDefault();
    },
  );

  onBeforeMount(loadConfigurations);
 
  return {
    configurations,
    activeConnections,
    connecting,
    loadingTables,
    tableFilter,  
    newTable,
    allTables,
    filteredTables,
    loadConfigurations,
    refreshTables,
    connect,
    connectDefault,
    selectDatabase,
    disconnect,
  };
}

/**
 * Returns the shared singleton instance of the database schema state.
 * Call this from any component to read or update the same connection/table state everywhere.
 */
export function useDatabaseSchema(): DatabaseSchemaState {
  if (!singleton) {
    singleton = createDatabaseSchema();
    // Ensure configurations are loaded when the first consumer mounts.
    // This avoids lifecycle scoping issues where the original create call
    // registered hooks on a different component and never triggered.
    void singleton.loadConfigurations();
    return singleton;
  }

  // If a new consumer imports the singleton later, make sure configs are present.
  if (singleton.configurations.value.length === 0) {
    void singleton.loadConfigurations();
  }

  return singleton;
}
