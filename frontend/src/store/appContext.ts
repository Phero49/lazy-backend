import { defineStore } from "pinia";

export interface WorkspaceTab {
  id: string;
  label: string;
  icon?: string;
  path: string;
  category?: string; // e.g., 'env', 'database', 'project'
  isReadonly?: boolean;
}

export const useAppContextStore = defineStore("appContext", {
  state: () => ({
    currentProject: null as string | null,
    currentRoute: "/",
    isInspectorOpen: true,
    isExplorerOpen: true,
    openedTabs: [] as WorkspaceTab[],
    activeTabId: null as string | null,
  }),
  getters: {
    activeTab: (state) =>
      state.openedTabs.find((tab) => tab.id === state.activeTabId) || null,
  },
  actions: {
    setCurrentProject(projectName: string | null) {
      this.currentProject = projectName;
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
  },
});
