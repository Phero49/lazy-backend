import { FieldType } from "./../../types.d";
import { defineStore } from "pinia";
import { io } from "../../wailsjs/go/models";
import {
  GetApiSchema,
  GetProjectSchema,
  SaveApiSchema,
  SaveProjectSchema,
} from "../../wailsjs/go/main/App";
import { LocationQuery, useRouter } from "vue-router";
import { nextTick, toRaw } from "vue";
import router from "../routes";
import { Condition } from "..";
export interface WorkspaceTab {
  id: string;
  label: string;
  icon?: string;
  path: string;
  category?: string; // e.g., 'env', 'database', 'project'
  isReadonly?: boolean;
  query?: LocationQuery;
}
export const httpMethod = [
  "post",
  "get",
  "put",
  "delete",
  "patch",
  "head",
  "options",
].map((v) => v.toUpperCase());

export type HTTPMethods =
  | "POST"
  | "GET"
  | "PUT"
  | "DELETE"
  | "PATCH"
  | "HEAD"
  | "OPTIONS";

export type Payload = {
  dataSource:
    | "body"
    | "form"
    | "param"
    | "query"
    | "header"
    | "injected values"
    | "-";
  fieldName: string;
  dataType: string;
  required: boolean;
};

export interface DatabaseQuery {
  tables: {
    table: string;
    join?: Join;
  }[];

  columns: {
    column: string;
    alias?: string;
  }[];
  where: Condition[];
}

export interface DatabaseOperations {
  name: string;
  databaseAction: "SELECT" | "UPDATE" | "DELETE" | "INSERT";
  query: DatabaseQuery;
}
export interface Routes {
  method: HTTPMethods | HTTPMethodsLower;
  path: string;
  data?: {
    payload: Payload[];
    operations: DatabaseOperations[];
  };
}

// If you need the lowercase version for some reason
export type HTTPMethodsLower = Lowercase<HTTPMethods>;
export interface EndPointList {
  prefix: string;
  collapsed?: boolean;
  children: (Routes | EndPointList)[];
}

export interface Join {
  baseTable: string;
  joinType: "INNER" | "LEFT" | "RIGHT" | "FULL";
  conditions: {
    operator: string;
    baseColumn: string;
    targetTable: string;
    targetColumn: string;
  }[];
}

export const useAppContextStore = defineStore("appContext", {
  state: () => ({
    currentProject: null as string | null,
    currentRoute: "/",
    isInspectorOpen: true,
    isExplorerOpen: true,
    activeLeftPanel: "projects",
    openedTabs: [] as WorkspaceTab[],
    activeTabId: null as string | null,
    activeConnectionName: null as string | null,
    databaseSchemas: {} as Record<string, io.IntermediateSchema>,
    apiSchema: [] as EndPointList[],
  }),
  getters: {
    activeTab: (state) =>
      state.openedTabs.find((tab) => tab.id === state.activeTabId) || null,
    activeSchema: (state) =>
      state.activeConnectionName
        ? state.databaseSchemas[state.activeConnectionName]
        : null,
  },
  actions: {
    async createNewTable(name: string) {
      const emptyField: io.Field = {
        isPrimaryKey: true,
        name: "",
        type: "string" as FieldType,
        length: 0,
        nullish: false,
        attributes: "-",
        autoGen: true,
        default: "-",
      } as io.Field;

      const newTable = {
        name: name,
        fields: [emptyField],
        primaryKey: [],
        description: "",
      } as unknown as io.Entity;
      if (
        this.activeConnectionName &&
        this.databaseSchemas[this.activeConnectionName]
      ) {
        this.databaseSchemas[this.activeConnectionName]["schema"].entities.push(
          newTable,
        );
        return;
      }

      this.databaseSchemas[this.activeConnectionName!] = {
        meta: {
          database: {
            name: "none",
            from: "none",
          },
          version: "",
          description: "",
        } as io.Meta,

        schema: {
          entities: [newTable],
          relationships: [],
        } as unknown as io.Schema,
      } as io.IntermediateSchema;
      SaveProjectSchema(this.currentProject!, this.databaseSchemas);
    },

    async getProjectSchema() {
      if (this.currentProject) {
        this.databaseSchemas = await GetProjectSchema(this.currentProject);
      }
    },
    setCurrentProject(projectName: string | null) {
      this.currentProject = projectName;
      this.getProjectSchema();
    },
    setCurrentRoute(route: string) {
      this.currentRoute = route;
    },
    toggleInspector() {
      this.isInspectorOpen = !this.isInspectorOpen;
    },
    toggleExplorer() {
      this.isExplorerOpen = !this.isExplorerOpen;
    },
    openTab(tab: WorkspaceTab) {
      const exists = this.openedTabs.find((t) => t.id === tab.id);
      if (!exists) {
        this.openedTabs.push(tab);
      }
      this.activeTabId = tab.id;
    },
    createNewTab(
      path: string,
      label: string,
      icon: string,
      query: LocationQuery = {},
    ) {
      const exists = this.openedTabs.find((t) => t.path === path);
      if (exists) {
        this.activeTabId = exists.id;
        return;
      }
      const id = `id-${this.openedTabs.length.toString()}`;
      this.openedTabs.push({
        id: id,
        label: label,
        path: path,
        icon: icon,
        query: query,
      });
      this.activeTabId = id;

      router.push({ path, force: true, query });

      //
    },
    closeTab(tabId: string) {
      const index = this.openedTabs.findIndex((t) => t.id === tabId);
      if (index !== -1) {
        this.openedTabs.splice(index, 1);
        if (this.activeTabId === tabId) {
          if (this.openedTabs.length > 0) {
            // Select the neighboring tab
            const nextTab = this.openedTabs[Math.max(0, index - 1)];
            this.activeTabId = nextTab.id;
          } else {
            this.activeTabId = null;
          }
        }
      }
    },
    setActiveTab(tabId: string) {
      this.activeTabId = tabId;
    },
    setActiveConnection(connectionName: string | null) {
      this.activeConnectionName = connectionName;
    },
    setDatabaseSchema(connectionName: string, schema: io.IntermediateSchema) {
      this.databaseSchemas[connectionName] = schema;
    },
    saveApi() {
      SaveApiSchema(
        JSON.stringify(toRaw(this.apiSchema)),
        this.currentProject!,
      );
    },
    async getApiSchema() {
      const api = await GetApiSchema(this.currentProject!);
      this.apiSchema = api;
    },
    getEndPoint(levels: number[]): EndPointList | null {
      let data: EndPointList | null = null;
      const loopthrough = (i: number, routes: EndPointList) => {
        return routes.children[i];
      };
      for (const level of levels) {
        data = this.apiSchema[level];
        if (levels.length > 1) {
          data = loopthrough(level, data) as EndPointList;
        }
      }

      return data;
    },
    deleteEndPont(levels: number[]) {
      if (levels.length === 1) {
        this.apiSchema.splice(levels[0], 1);
      } else {
        let data = this.apiSchema[0].children;
        for (let index = 1; index < levels.length; index++) {
          const i = levels[index];
          if (index >= levels.length - 1) {
            break;
          }
          data = (data[i] as EndPointList).children;
        }
      }
      this.saveApi();
    },
  },
});
