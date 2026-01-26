<template>
    <div class="q-pa-md column full-height no-wrap">
        <!-- Explorer Header -->
        <div class="row items-center justify-between q-mb-md">
            <div class="text-overline text-grey-7">Database Explorer</div>
            <q-btn flat round dense icon="bi-plus-lg" size="xs" color="primary"
                :to="{ name: 'update-env', params: { type: 'database' } }">
                <q-tooltip>Create New Database</q-tooltip>
            </q-btn>
        </div>

        <!-- Connection Manager (Expansion) -->
        <q-list padding class="q-mb-sm">
            <q-expansion-item class="border-lighter border-radius-md overflow-hidden bg-dark-soft"
                header-class="text-weight-bold" dense-toggle default-opened>
                <template v-slot:header>
                    <q-item-section avatar class="min-width-auto q-pr-sm">
                        <q-icon name="bi-hdd-network" color="primary" size="16px" />
                    </q-item-section>
                    <q-item-section>
                        <div class="text-caption">
                            {{ activeConnectionName ? `Active: ${activeConnectionName}` : 'Manage Connections' }}
                        </div>
                    </q-item-section>
                    <q-item-section side>
                        <q-badge :label="`${activeConnections.size} connected`" color="grey-10" text-color="grey-6"
                            style="font-size: 9px" />
                    </q-item-section>
                </template>

                <q-card dark class="bg-transparent border-top-lighter">
                    <q-list dense padding>
                        <q-item v-for="config in configurations" :key="config.name" clickable
                            class="q-mx-sm border-radius-sm" :active="activeConnectionName === config.name"
                            active-class="active-db-item" @click="selectDatabase(config.name)">
                            <q-item-section avatar class="min-width-auto q-pr-sm">
                                <q-icon :name="activeConnections.has(config.name) ? 'bi-database-fill' : 'bi-database'"
                                    :color="activeConnections.has(config.name) ? 'primary' : 'grey-8'" size="13px" />
                            </q-item-section>
                            <q-item-section class="text-caption">{{ config.name }}</q-item-section>
                            <q-item-section side>
                                <q-spinner-dots v-if="connecting === config.name" color="primary" size="xs" />
                                <template v-else>
                                    <q-btn v-if="!activeConnections.has(config.name)" flat round dense
                                        icon="bi-play-fill" size="xs" color="grey-8" @click.stop="connect(config)" />
                                    <q-btn v-else flat round dense icon="bi-stop-fill" size="xs" color="red-5"
                                        @click.stop="disconnect(config.name)" />
                                </template>
                            </q-item-section>
                        </q-item>
                    </q-list>
                </q-card>
            </q-expansion-item>
        </q-list>

        <!-- Main Context Area (Tables) -->
        <div v-if="activeConnectionName" class="col column no-wrap">
            <!-- Filter & Toolbar (Sticky-like Header) -->
            <div class="column q-mb-sm">
                <div class="row items-center justify-between q-py-xs">
                    <div class="text-caption text-grey-9 uppercase-tracking text-weight-bold">
                        TABLES ({{ filteredTables.length }})
                    </div>
                    <q-btn flat round dense icon="bi-arrow-clockwise" size="xs" color="grey-9"
                        @click="refreshTables(activeConnectionName)" />
                </div>
                <q-input v-model="tableFilter" placeholder="Filter tables..." dark dense standout
                    class="table-filter-input">
                    <template v-slot:prepend>
                        <q-icon name="bi-search" size="12px" color="grey-8" />
                    </template>
                </q-input>
            </div>

            <!-- Table List (Scrollable) -->

            <div v-if="loadingTables" class="column items-center q-pa-xl">
                <q-spinner-grid color="primary" size="md" />
            </div>
            <q-list v-else dense dark class="q-gutter-y-xs">
                <q-item v-for="table in filteredTables" :key="table" clickable class="border-radius-sm table-item"
                    @click="openTableSchema(activeConnectionName!, table)">
                    <q-item-section avatar class="min-width-auto q-pr-sm">
                        <q-icon name="bi-table" size="14px" color="primary" />
                    </q-item-section>
                    <q-item-section class="text-caption text-grey-4">{{ table }}</q-item-section>
                </q-item>
                <div v-if="filteredTables.length === 0" class="text-center q-pa-xl text-grey-9 text-caption italic">
                    {{ tableFilter ? 'No tables match your filter.' : 'No tables found in this database.' }}
                </div>
            </q-list>

        </div>

        <!-- Empty States -->
        <div v-else class="column items-center justify-center col text-grey-9 text-center">
            <q-icon name="bi-database" size="56px" class="q-mb-md opacity-20" />
            <div class="text-subtitle2" v-if="configurations.length > 0">Select or connect to a database to explore
                tables.
            </div>
            <div v-else>
                <div class="text-subtitle2">No database profiles configured.</div>
                <q-btn flat no-caps label="Create Connection" color="primary" class="q-mt-sm"
                    :to="{ name: 'update-env', params: { type: 'database' } }" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref, computed, watch, reactive } from 'vue';
import { GetProjectEnv, GetTableList } from '../../wailsjs/go/main/App';
import { useAppContextStore, WorkspaceTab } from '../store/appContext';
import { useQuasar } from 'quasar';

const appContext = useAppContextStore();
const $q = useQuasar();

interface DBProfile {
    name: string;
    config: any;
}

const configurations = ref<DBProfile[]>([]);
const activeConnections = ref<Map<string, string[]>>(new Map()); // Name -> Tables
const activeConnectionName = ref<string | null>(null);
const connecting = ref<string | null>(null);
const loadingTables = ref(false);
const tableFilter = ref('');

const filteredTables = computed(() => {
    if (!activeConnectionName.value) return [];
    const tables = activeConnections.value.get(activeConnectionName.value) || [];
    if (!tableFilter.value) return tables;
    const search = tableFilter.value.toLowerCase();
    return tables.filter((t: string) => t.toLowerCase().includes(search));
});

async function loadConfigurations() {
    if (!appContext.currentProject) return;
    try {
        const env = await GetProjectEnv(appContext.currentProject);
        const dbList = env.database || [];
        configurations.value = dbList.map((item: any) => {
            const name = Object.keys(item)[0];
            return {
                name,
                config: item[name]
            };
        });
    } catch (error) {
        console.error('Failed to load configurations:', error);
    }
}

async function connect(profile: DBProfile) {
    connecting.value = profile.name;
    try {
        const tables = await GetTableList(
            profile.config.host,
            profile.config.port,
            profile.config.user,
            profile.config.password,
            profile.config.database
        );
        activeConnections.value.set(profile.name, tables);
        activeConnectionName.value = profile.name;
        $q.notify({
            type: 'positive',
            message: `Connected to ${profile.name}`,
            timeout: 1000,
            position: 'bottom-right'
        });
    } catch (error: any) {
        $q.notify({
            type: 'negative',
            message: `Connection failed: ${error}`,
        });
    } finally {
        connecting.value = null;
    }
}

function selectDatabase(name: string) {
    if (activeConnections.value.has(name)) {
        activeConnectionName.value = name;
    } else {
        const profile = configurations.value.find(c => c.name === name);
        if (profile) connect(profile);
    }
}

function disconnect(name: string) {
    activeConnections.value.delete(name);
    if (activeConnectionName.value === name) {
        activeConnectionName.value = activeConnections.value.size > 0
            ? Array.from(activeConnections.value.keys())[0]
            : null;
    }
}

function openTableSchema(connectionName: string, tableName: string) {
    const tab: WorkspaceTab = {
        id: `schema-${connectionName}-${tableName}`,
        label: `${tableName}`,
        icon: 'bi-layout-text-sidebar-reverse',
        path: `/database/schema/${connectionName}/${tableName}`,
        category: 'database',
        isReadonly: true
    };
    appContext.openTab(tab);
}

async function refreshTables(connectionName: string) {
    const profile = configurations.value.find(c => c.name === connectionName);
    if (!profile) return;

    loadingTables.value = true;
    try {
        const tables = await GetTableList(
            profile.config.host,
            profile.config.port,
            profile.config.user,
            profile.config.password,
            profile.config.database
        );
        activeConnections.value.set(connectionName, tables);
    } catch (error) {
        $q.notify({ type: 'negative', message: 'Failed to refresh tables' });
    } finally {
        loadingTables.value = false;
    }
}

watch(() => appContext.currentProject, () => {
    activeConnections.value.clear();
    activeConnectionName.value = null;
    loadConfigurations();
});

onBeforeMount(loadConfigurations);

</script>

<style lang="scss" scoped>
.min-width-auto {
    min-width: unset !important;
}

.bg-primary-opacity {
    background: rgba(var(--q-primary), 0.1) !important;
}

.border-lighter {
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.border-radius-md {
    border-radius: 8px;
}

.active-db-item {
    background: rgba(var(--q-primary), 0.1);
    color: var(--q-primary);
}

.uppercase-tracking {
    text-transform: uppercase;
    letter-spacing: 0.1em;
    font-size: 10px;
}

.table-filter-input {
    :deep(.q-field__control) {
        background: rgba(255, 255, 255, 0.03) !important;
        border: 1px solid rgba(255, 255, 255, 0.05);
        border-radius: 6px;
        min-height: 32px;
        height: 32px;
        font-size: 11px;
    }
}

.table-item {
    &:hover {
        background: rgba(255, 255, 255, 0.03);
    }
}

.opacity-20 {
    opacity: 0.2;
}
</style>
