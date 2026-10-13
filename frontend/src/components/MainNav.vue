<template>
  <div class="row fit">
    <div
      class="col-auto column full-height items-center q-py-md bg-dark-soft border-right"
    >
      <!-- Main Navigation Items -->
      <div class="column items-center q-gutter-y-sm full-width">
        <q-btn
          v-for="item in navItems"
          :key="item.label"
          flat
          round
          :disable="store.currentProject == null"
          dense
          :icon="item.icon"
          @click="store.activeLeftPanel = item.to"
          class="nav-btn transition-base"
          active-class="active-nav-btn"
        >
          <q-tooltip
            anchor="center right"
            self="center left"
            :offset="[14, 0]"
            class="bg-grey-10 text-primary text-weight-bold text-uppercase tracking-widest border-primary"
          >
            {{ item.label }}
          </q-tooltip>
        </q-btn>
      </div>

      <q-space />

      <!-- Bottom Actions -->
      <div class="column items-center q-gutter-y-sm q-pb-md">
        <q-btn
          flat
          round
          dense
          icon="bi-question-circle"
          class="text-grey-7 hover-white"
        />
        <q-btn
          flat
          round
          dense
          icon="bi-gear"
          to="/settings"
          class="text-grey-7 hover-white"
        />
      </div>
    </div>
    <div
      v-if="store.isExplorerOpen"
      class="col-grow bg-dark-soft border-right column overflow-hidden transition-base"
    >
      <q-scroll-area
        v-if="store.currentProject"
        class="col"
        style="height: calc(100% - 90px); overflow-x: hidden"
      >
        <component
          :is="sidePanels[store.activeLeftPanel as keyof c]"
        ></component>
      </q-scroll-area>
      <div v-else class="text-center text-subtitle1 q-py-xl">
        open a project to start working
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { watch } from "vue";
import { useRoute } from "vue-router";
import { useAppContextStore } from "../store/appContext";
import ProjectsPage from "../pages/leftPanel/ProjectsPage.vue";
import DatabaseTablesPage from "../pages/leftPanel/DatabaseTablesPage.vue";
import EndpointsPage from "../pages/leftPanel/EndpointsPage.vue";
import PluginsPage from "../pages/leftPanel/PluginsPage.vue";
import ProjectEnvPage from "../pages/leftPanel/ProjectEnvPage.vue";
import Dependencies from "../pages/leftPanel/Dependencies.vue";

const route = useRoute();
const store = useAppContextStore();

const navItems = [
  { icon: "bi-git", label: "changes", to: "changes" },
  { icon: "bi-table", label: "tables", to: "tables" },
  { icon: "bi-diagram-3", label: "api Endpoints", to: "api" },
  { icon: "bi-plugin", label: "Plugins", to: "plugins" },
  { icon: "bi-layers", label: "Dependencies", to: "dependencies" },
  { icon: "bi-code", label: "Env", to: "env" },
];

const sidePanels = {
  changes: ProjectsPage,
  tables: DatabaseTablesPage,
  api: EndpointsPage,
  plugins: PluginsPage,
  env: ProjectEnvPage,
  dependencies: Dependencies,
};
type c = typeof sidePanels;
// Sync current route to store
watch(
  () => route.path,
  (newPath) => {
    store.setCurrentRoute(newPath);
  },
  { immediate: true },
);
</script>

<style lang="scss" scoped>
.bg-dark-soft {
  background: #0d0d0d;
}

.border-right {
  border-right: 1px solid rgba(255, 255, 255, 0.05);
}

.transition-base {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.nav-btn {
  color: #666;
  width: 44px;
  height: 44px;
  border-radius: 12px;

  &:hover {
    color: #fff;
    background: rgba(255, 255, 255, 0.05);
  }
}

.active-nav-btn {
  color: var(--q-primary) !important;
  background: rgba(var(--q-primary), 0.1) !important;
  box-shadow: inset 0 0 10px rgba(var(--q-primary), 0.05);

  &::before {
    content: "";
    position: absolute;
    left: -12px;
    top: 20%;
    bottom: 20%;
    width: 3px;
    background: var(--q-primary);
    border-radius: 0 4px 4px 0;
    box-shadow: 0 0 10px var(--q-primary);
  }
}

.tracking-widest {
  letter-spacing: 0.15em;
  font-size: 10px;
}

.border-primary {
  border: 1px solid rgba(var(--q-primary), 0.2);
}

.hover-white:hover {
  color: #fff !important;
}
</style>
