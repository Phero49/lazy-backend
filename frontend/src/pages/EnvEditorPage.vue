<template>
    <div class="col column bg-dark-page" v-if="config">
        <div class="q-pa-md border-bottom bg-dark-soft row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
                <q-icon :name="getIcon(category)" color="primary" />
                <span class="text-caption text-weight-bold text-white">{{ category.toUpperCase() }} / {{ name }}</span>
                <q-badge v-if="store.activeTab?.isReadonly" label="READ-ONLY" color="grey-9" text-color="grey-5"
                    class="q-ml-sm" />
            </div>
            <div class="row items-center q-gutter-x-sm">
                <template v-if="store.activeTab?.isReadonly">
                    <q-btn unelevated color="primary" label="Edit" size="sm" no-caps icon="bi-pencil"
                        @click="editConfig" />
                    <q-btn flat color="red-5" label="Delete" size="sm" no-caps icon="bi-trash" @click="deleteConfig" />
                </template>
                <template v-else>
                    <q-btn flat label="Cancel" size="sm" no-caps color="grey-7" @click="cancelEdit" />
                    <q-btn unelevated color="primary" label="Save Changes" size="sm" no-caps icon="bi-check-lg"
                        @click="saveChanges" />
                </template>
            </div>
        </div>

        <q-scroll-area class="col q-pa-xl">
            <div class="max-width-600 q-mx-auto">
                <div class="column q-gutter-y-lg">
                    <div v-for="(val, key) in config" :key="key" class="column q-gutter-y-xs">
                        <span class="text-overline text-grey-7 size-9">{{ (key as unknown as string).toUpperCase()
                            }}</span>
                        <q-toggle v-if="typeof val === 'boolean'" v-model="config[key]" color="primary" dark
                            :disable="store.activeTab?.isReadonly" />
                        <q-input v-else v-model="config[key]" dark dense standout class="config-input"
                            :type="typeof val === 'number' ? 'number' : 'text'"
                            :disable="store.activeTab?.isReadonly" />
                    </div>
                </div>
            </div>
        </q-scroll-area>
    </div>
    <div v-else class="col column items-center justify-center text-grey-9">
        <q-icon name="bi-sliders" size="64px" class="q-mb-md opacity-20" />
        <div class="text-subtitle2">Select a configuration to edit</div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { GetProjectEnv, SaveProjectEnv } from '../../wailsjs/go/main/App';
import { useAppContextStore } from '../store/appContext';
import { useQuasar } from 'quasar';

const route = useRoute();
const store = useAppContextStore();
const $q = useQuasar();

const category = ref(route.params.category as string);
const name = ref(route.params.name as string);
const config = ref<any>(null);
const originalConfig = ref<any>(null);

const loadConfig = async () => {
    category.value = route.params.category as string;
    name.value = route.params.name as string;

    if (store.currentProject && category.value && name.value) {
        const envData = await GetProjectEnv(store.currentProject);
        const categoryData = envData[category.value];
        if (categoryData) {
            const item = categoryData.find((i: any) => Object.keys(i)[0] === name.value);
            if (item) {
                config.value = JSON.parse(JSON.stringify(item[name.value]));
                originalConfig.value = JSON.parse(JSON.stringify(item[name.value]));
            }
        }
    }
}

onMounted(loadConfig);
watch(() => [route.params.category, route.params.name, store.currentProject], loadConfig);

const getIcon = (cat: string) => {
    switch (cat) {
        case 'database': return 'bi-database';
        case 'ssh': return 'bi-terminal-fill';
        case 'apiKeys': return 'bi-key';
        default: return 'bi-gear';
    }
}

const editConfig = () => {
    if (store.activeTab) {
        store.activeTab.isReadonly = false;
    }
}

const cancelEdit = () => {
    config.value = JSON.parse(JSON.stringify(originalConfig.value));
    if (store.activeTab) {
        store.activeTab.isReadonly = true;
    }
}

const deleteConfig = () => {
    $q.dialog({
        title: 'Confirm Delete',
        message: `Are you sure you want to delete the "${name.value}" configuration from ${category.value}?`,
        cancel: true,
        persistent: true,
        dark: true
    }).onOk(async () => {
        if (!store.currentProject) return;
        try {
            const envData = await GetProjectEnv(store.currentProject);
            const categoryData = envData[category.value];
            if (categoryData) {
                const index = categoryData.findIndex((i: any) => Object.keys(i)[0] === name.value);
                if (index !== -1) {
                    categoryData.splice(index, 1);
                    const res = await SaveProjectEnv(store.currentProject, envData);
                    if (res.success) {
                        $q.notify({ type: 'positive', message: 'Configuration deleted', position: 'bottom-right' });
                        if (store.activeTabId) {
                            store.closeTab(store.activeTabId);
                        }
                    } else {
                        $q.notify({ type: 'negative', message: res.message });
                    }
                }
            }
        } catch (error) {
            $q.notify({ type: 'negative', message: 'Error deleting configuration' });
        }
    });
}

const saveChanges = async () => {
    if (!store.currentProject) return;
    try {
        const envData = await GetProjectEnv(store.currentProject);
        const categoryData = envData[category.value];
        if (categoryData) {
            const item = categoryData.find((i: any) => Object.keys(i)[0] === name.value);
            if (item) {
                item[name.value] = config.value;
                const res = await SaveProjectEnv(store.currentProject, envData);
                if (res.success) {
                    $q.notify({
                        type: 'positive',
                        message: 'Configuration saved successfully',
                        position: 'bottom-right'
                    });
                    originalConfig.value = JSON.parse(JSON.stringify(config.value));
                    if (store.activeTab) {
                        store.activeTab.isReadonly = true;
                    }
                } else {
                    $q.notify({ type: 'negative', message: res.message });
                }
            }
        }
    } catch (error) {
        $q.notify({ type: 'negative', message: 'Error saving configuration' });
    }
}
</script>

<style lang="scss" scoped>
.size-9 {
    font-size: 9px;
}

.max-width-600 {
    max-width: 600px;
}

.bg-dark-page {
    background: #050505;
}

.bg-dark-soft {
    background: #0d0d0d;
}

.border-bottom {
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.opacity-20 {
    opacity: 0.2;
}

.config-input {
    :deep(.q-field__control) {
        background: rgba(255, 255, 255, 0.03) !important;
        border: 1px solid rgba(255, 255, 255, 0.05);
        border-radius: 8px;

        &:hover {
            border-color: rgba(var(--q-primary), 0.3);
        }

        &:focus-within {
            border-color: var(--q-primary);
        }
    }
}
</style>
