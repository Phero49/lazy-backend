<template>
  <div class="q-pa-md">
    <!-- Header -->
    <div class="row justify-between items-center q-mb-md">
      <div class="text-overline text-grey-7">Plugins</div>
      <q-btn color="primary" icon="more_horiz" flat dense>
        <q-menu>
          <q-list style="min-width: 180px">
            <q-item clickable v-close-popup @click="handleLoadLocal">
              <q-item-section avatar>
                <q-icon name="folder_open" />
              </q-item-section>
              <q-item-section>Install Local Plugin</q-item-section>
            </q-item>
            <q-separator />
            <q-item clickable v-close-popup @click="refreshList">
              <q-item-section avatar>
                <q-icon name="refresh" />
              </q-item-section>
              <q-item-section>Refresh List</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </div>

    <!-- Filter Tabs -->
    <q-select
      v-model="pluginFilter"
      :options="filterOptions"
      label="View"
      filled
      dense
      class="q-mb-md"
      @update:model-value="refreshList"
    />

    <!-- Loading State -->
    <q-inner-loading :showing="loading">
      <q-spinner size="50px" color="primary" />
    </q-inner-loading>

    <!-- Empty State -->
    <div v-if="!loading && filteredPlugins.length === 0" class="text-center q-pa-xl text-grey-6">
      <q-icon name="extension_off" size="48px" class="q-mb-sm" />
      <div class="text-body1">No plugins found</div>
      <div class="text-caption">
        {{ pluginFilter === 'installed' ? 'Install a local plugin to get started' :
           pluginFilter === 'active' ? 'No active plugins in this project' :
           'Marketplace coming soon' }}
      </div>
    </div>

    <!-- Plugin List -->
    <q-list v-if="!loading && filteredPlugins.length > 0" separator bordered class="rounded-borders">
      <q-item v-for="plugin in filteredPlugins" :key="plugin.name" clickable>
        <q-item-section avatar>
          <q-avatar color="primary" text-color="white" font-size="12px">
            {{ plugin.publisher?.charAt(0)?.toUpperCase() || '?' }}
          </q-avatar>
        </q-item-section>

        <q-item-section>
          <q-item-label>{{ plugin.name }}</q-item-label>
          <q-item-label caption>{{ plugin.description || 'No description' }}</q-item-label>
        </q-item-section>

        <q-item-section side>
          <div class="column items-end q-gutter-y-xs">
            <q-badge color="grey-7" :label="plugin.version" />
            <q-badge
              v-if="isActive(plugin.name)"
              color="positive"
              label="Active"
            />
          </div>
        </q-item-section>

        <q-item-section side>
          <q-btn
            v-if="pluginFilter === 'installed' && !isActive(plugin.name)"
            flat dense round icon="add" color="primary"
            @click.stop="handleAddToProject(plugin.name)"
          >
            <q-tooltip>Add to project</q-tooltip>
          </q-btn>
          <q-btn
            v-else-if="pluginFilter === 'active'"
            flat dense round icon="remove" color="negative"
            @click.stop="handleRemoveFromProject(plugin.name)"
          >
            <q-tooltip>Remove from project</q-tooltip>
          </q-btn>
          <q-btn
            v-else-if="pluginFilter === 'installed' && isActive(plugin.name)"          flat dense round icon="check_circle" color="positive" disable
          >
            <q-tooltip>Already active</q-tooltip>
          </q-btn>
        </q-item-section>
      </q-item>
    </q-list>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import {
  LoadLocalPlugin,
  GetPluginsList,
  GetActivePlugins,
  AddPluginToProject,
  RemovePluginFromProject,
} from '#/wailsjs/go/main/App';
import type { io } from '#/wailsjs/go/models';
import { useAppContextStore } from '~/store/appContext';

const $q = useQuasar();
const appContext = useAppContextStore();

const pluginFilter = ref<'installed' | 'market' | 'active'>('installed');
const filterOptions = ['installed', 'market', 'active'];

const loading = ref(false);
const installedPlugins = ref<io.PluginItem[]>([]);
const activePluginNames = ref<Set<string>>(new Set());

const project = computed(() => appContext.currentProject);

// Filtered view based on selected tab
const filteredPlugins = computed(() => {
  switch (pluginFilter.value) {
    case 'active':
      return installedPlugins.value.filter((p) => activePluginNames.value.has(p.name));
    case 'market':
      // Placeholder: market plugins would come from a remote source
      return [];
    case 'installed':
    default:
      return installedPlugins.value;
  }
});

function isActive(pluginName: string): boolean {
  return activePluginNames.value.has(pluginName);
}

async function refreshList() {
  loading.value = true;
  try {
    const [plugins, active] = await Promise.all([
      GetPluginsList(),
      project.value ? GetActivePlugins(project.value) : Promise.resolve([]),
    ]);
    installedPlugins.value = plugins ?? [];
    activePluginNames.value = new Set((active ?? []).map((p) => p.name));
  } catch (err) {
    $q.notify({ type: 'negative', message: String(err) });
  } finally {
    loading.value = false;
  }
}

async function handleLoadLocal() {
  try {
    await LoadLocalPlugin();
    $q.notify({ type: 'positive', message: 'Local plugin installed' });
    await refreshList();
  } catch (err) {
    // Empty string means user cancelled — don't show error
    if (String(err).includes('cancelled')) return;
    $q.notify({ type: 'negative', message: String(err) });
  }
}

async function handleAddToProject(pluginName: string) {
  if (!project.value) {
    $q.notify({ type: 'warning', message: 'No project selected' });
    return;
  }
  try {
    await AddPluginToProject(project.value, pluginName);
    $q.notify({ type: 'positive', message: `Added ${pluginName} to project` });
    await refreshList();
  } catch (err) {
    $q.notify({ type: 'negative', message: String(err) });
  }
}

async function handleRemoveFromProject(pluginName: string) {
  if (!project.value) return;
  try {
    await RemovePluginFromProject(project.value, pluginName);
    $q.notify({ type: 'info', message: `Removed ${pluginName} from project` });
    await refreshList();
  } catch (err) {
    $q.notify({ type: 'negative', message: String(err) });
  }
}

onMounted(refreshList);
</script>