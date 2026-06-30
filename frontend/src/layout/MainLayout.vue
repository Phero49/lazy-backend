<template>
  <q-layout view="hHh LpR fFf" class="bg-dark-page text-grey-5">
    <!-- Sleek Top Header -->
    <q-header class="bg-dark-soft border-bottom">
      <q-toolbar class="q-px-lg">
        <div class="row items-center q-gutter-x-md">
          <div class="app-logo">
            <q-icon name="bi-lightning-fill" color="primary" size="20px" />
          </div>
          <div class="column">
            <span
              class="text-white text-weight-bolder text-uppercase tracking-m size-12"
              >Lazy Backend</span
            >
            <span
              class="text-primary text-weight-bold size-9 text-uppercase tracking-xl"
              >System V2</span
            >
          </div>
        </div>

        <q-space />

        <!-- Minimal Search -->
        <q-input
          dark
          dense
          standout
          v-model="search"
          placeholder="Search components..."
          class="search-input gt-sm"
        >
          <template v-slot:prepend>
            <q-icon name="bi-search" size="14px" />
          </template>
        </q-input>

        <q-space />

        <div class="row items-center q-gutter-x-sm">
          <q-btn
            flat
            round
            dense
            icon="bi-layout-sidebar-inset"
            @click="store.toggleExplorer()"
            :color="store.isExplorerOpen ? 'primary' : 'grey'"
            size="sm"
          />
          <q-btn
            flat
            round
            dense
            icon="bi-layout-sidebar-inset-reverse"
            @click="store.toggleInspector()"
            :color="store.isInspectorOpen ? 'primary' : 'grey'"
            size="sm"
          />
          <q-separator vertical dark class="q-mx-sm bg-grey-9" inset />
          <q-btn flat round dense icon="bi-bell" />
          <q-separator vertical dark class="q-mx-sm bg-grey-9" inset />
          <q-btn flat round class="q-pa-none">
            <q-avatar size="32px" class="border-grey">
              <img
                src="https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"
              />
            </q-avatar>
          </q-btn>
        </div>
      </q-toolbar>
    </q-header>
    <q-drawer
      side="left"
      flat
      v-model="store.isExplorerOpen"
      :width="$q.screen.width * 0.25"
      :breakpoint="500"
      content-class="bg-grey-3"
    >
      <MainNav />
    </q-drawer>
    <q-page-container>
      <q-page class="">
        <div
          style="z-index: 1000; margin-left: 5px"
          class="fixed q-mb-lg bg-dark-soft border-bottom"
        >
          <div
            class="row overflow-hidden no-wrap wails-drag"
            style="height: 36px"
          >
            <!-- Important: Put a nodrag class on interactive items like tabs -->
            <div
              v-for="tab in store.openedTabs"
              :key="tab.id"
              class="workspace-tab row items-center q-px-md transition-base wails-nodrag"
              :class="{ 'active-tab': store.activeTabId === tab.id }"
              @click="onTabClick(tab.id, tab.path)"
            >
              <q-icon
                :name="tab.icon || 'bi-file-earmark-code'"
                color="primary"
                size="14px"
                class="q-mr-sm"
              />
              <span
                class="text-caption text-weight-medium"
                :class="
                  store.activeTabId === tab.id ? 'text-white' : 'text-grey-7'
                "
              >
                {{ tab.label }}
              </span>
              <q-icon
                name="bi-x"
                size="18px"
                class="q-ml-sm text-grey-7 hover-white cursor-pointer close-icon"
                @click.stop="store.closeTab(tab.id)"
              />
            </div>
          </div>
        </div>
        <!-- Center Workspace -->
        <div style="padding-top: 35px">
          <!-- Workspace Tabs A rea (MAKE THIS THE DRAG HANDLE) -->

          <!-- Quasar Scroll Area (STRICTLY NO DRAG) -->
          <router-view />
        </div>
      </q-page>
    </q-page-container>

    <!-- Professional Status Bar -->
    <q-footer class="bg-dark-soft border-top q-py-xs q-px-md">
      <div class="row items-center no-wrap size-10 text-grey-7">
        <div class="row items-center q-gutter-x-md">
          <div class="row items-center q-gutter-x-xs">
            <div class="status-dot online" />
            <span class="text-primary text-weight-bold"
              >Go Service: Online</span
            >
          </div>
          <q-separator vertical dark class="bg-grey-9" />
          <div class="row items-center q-gutter-x-xs">
            <q-icon name="bi-git" size="12px" />
            <span>main*</span>
          </div>
        </div>
        <q-space />
        <div class="row items-center q-gutter-x-md">
          <span>Go 1.21.4</span>
          <span>UTF-8</span>
          <q-icon name="bi-bell-fill" size="12px" color="primary" />
        </div>
      </div>
    </q-footer>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import MainNav from "../components/MainNav.vue";
import { useAppContextStore } from "../store/appContext";

const router = useRouter();
const route = useRoute();
const store = useAppContextStore();
const search = ref("");

const onTabClick = (tabId: string, path: string) => {
  store.setActiveTab(tabId);
  router.push({ path, force: true });
};

// Keep active tab ID in sync with the current URL path
watch(
  () => route.fullPath,
  (newPath) => {
    const matchingTab = store.openedTabs.find((tab) => tab.path === newPath);
    if (matchingTab && store.activeTabId !== matchingTab.id) {
      store.setActiveTab(matchingTab.id);
    }
  },
  { immediate: true },
);
</script>

<style lang="scss" scoped>
.bg-dark-page {
  background: #050505;
}

.bg-dark-soft {
  background: #0d0d0d;
}

.border-bottom {
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.border-right {
  border-right: 1px solid rgba(255, 255, 255, 0.05);
}

.border-left {
  border-left: 1px solid rgba(255, 255, 255, 0.05);
}

.border-top {
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.app-logo {
  width: 32px;
  height: 32px;
  background: rgba(var(--q-primary), 0.1);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(var(--q-primary), 0.2);
}

.tracking-m {
  letter-spacing: 0.1em;
}

.tracking-xl {
  letter-spacing: 0.2em;
}

.size-12 {
  font-size: 12px;
}

.size-9 {
  font-size: 9px;
}

.size-10 {
  font-size: 10px;
}

.search-input {
  width: 340px;

  :deep(.q-field__control) {
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.03) !important;
    border: 1px solid rgba(255, 255, 255, 0.05);

    &:hover {
      border-color: rgba(var(--q-primary), 0.3);
    }

    &:focus-within {
      border-color: var(--q-primary);
    }
  }

  :deep(.q-field__native) {
    font-size: 11px;
  }
}

.workspace-tab {
  border-right: 1px solid rgba(255, 255, 255, 0.05);
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: rgba(255, 255, 255, 0.02);
  }

  &.active-tab {
    background: #050505;
    border-top: 2px solid var(--q-primary);

    span {
      color: #fff;
    }
  }
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;

  &.online {
    background: var(--q-primary);
    box-shadow: 0 0 8px var(--q-primary);
  }
}

.border-grey {
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.hover-white:hover {
  color: #fff !important;
}
</style>
