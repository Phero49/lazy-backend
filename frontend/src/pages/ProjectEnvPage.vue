<template>
    <div class="bg-dark-page">
        <div>
            <!-- Left List (Categories) -->
            <div class="col-4 border-right bg-dark-soft ">
                <div class="q-pa-md row items-center justify-between border-bottom">
                    <span class="text-overline text-grey-7 text-weight-bold">Env Variables</span>
                    <q-btn flat round dense icon="bi-plus-lg" size="xs" color="primary">
                        <q-menu dark anchor="bottom right" self="top right" class="bg-grey-10 border-lighter">
                            <q-list dense style="min-width: 150px">
                                <q-item clickable :to="{
                                    name: 'update-env',

                                    params: {
                                        type: type.toLocaleLowerCase()
                                    },

                                    replace: true
                                }" v-close-popup v-for="type in ['Database', 'API Key', 'SSH']" :key="type">
                                    <q-item-section class="text-caption">Add {{ type }}</q-item-section>
                                </q-item>
                            </q-list>
                        </q-menu>
                    </q-btn>
                </div>

                <q-tabs v-model="tab" dense class="text-grey-7 border-bottom" active-color="primary"
                    indicator-color="primary" align="justify" narrow-indicator no-caps>
                    <q-tab name="dev" label="Development" />
                    <q-tab name="production" label="Production" />
                </q-tabs>


                <q-tab-panels v-model="tab" animated class="bg-transparent col">
                    <q-tab-panel name="dev" class="q-pa-none">

                        <div v-if="!envData" class="column items-center q-pa-xl text-center">
                            <q-icon name="bi-terminal" size="48px" color="grey-9" />
                            <div class="text-caption text-grey-8 q-mt-sm">Select a project to view environment</div>
                        </div>
                        <q-list v-else dense padding class="q-gutter-y-xs">
                            <q-expansion-item v-for="(items, category) in envData" :key="category"
                                :label="(category as unknown as string).toUpperCase()"
                                header-class="text-overline text-grey-6 text-weight-bold" class="border-bottom-lighter"
                                default-opened>
                                <q-list dense padding>
                                    <q-item clickable v-for="(item, index) in items" :key="index"
                                        class="q-mx-sm border-radius-sm"
                                        @click="openConfigTab(category as unknown as string, item)"
                                        :active="store.activeTabId === `${category}-${Object.keys(item)[0]}`"
                                        active-class="bg-primary-opacity text-primary">
                                        <q-item-section avatar class="min-width-auto q-pr-sm">
                                            <q-icon :name="getIcon(category as unknown as string)" size="14px" />
                                        </q-item-section>
                                        <q-item-section class="text-caption text-weight-medium">
                                            {{ Object.keys(item)[0] }}
                                        </q-item-section>
                                    </q-item>
                                    <div v-if="items.length === 0"
                                        class="q-pa-md text-center text-grey-9 text-italic text-caption">
                                        No {{ category }} configured
                                    </div>
                                </q-list>
                            </q-expansion-item>
                        </q-list>
                    </q-tab-panel>

                    <q-tab-panel name="production" class="q-pa-none column">
                        <div class="column items-center justify-center col text-grey-8">
                            <q-icon name="bi-shield-lock" size="32px" class="q-mb-md" />
                            <div class="text-caption">Production settings are encrypted by default</div>
                            <q-btn outline color="primary" label="Unlock" size="sm" class="q-mt-md" no-caps />
                        </div>
                    </q-tab-panel>
                </q-tab-panels>
            </div>


        </div>
    </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref, watch } from 'vue';
import { GetProjectEnv } from '../../wailsjs/go/main/App';
import { useAppContextStore, WorkspaceTab } from '../store/appContext';

const store = useAppContextStore()
const envData = ref<any>()
const tab = ref('dev')

const loadEnv = async () => {
    if (store.currentProject) {
        envData.value = await GetProjectEnv(store.currentProject)
    } else {
        envData.value = null
    }
}

onBeforeMount(loadEnv)
watch(() => store.currentProject, loadEnv)

const openConfigTab = (category: string, item: any) => {
    const label = Object.keys(item)[0]
    const tab: WorkspaceTab = {
        id: `${category}-${label}`,
        label: label,
        icon: getIcon(category),
        path: `/env/edit/${category}/${label}`,
        category: 'env',
        isReadonly: true
    }
    store.openTab(tab)
}

const getIcon = (category: string) => {
    switch (category) {
        case 'database': return 'bi-database';
        case 'ssh': return 'bi-terminal-fill';
        case 'apiKeys': return 'bi-key';
        default: return 'bi-gear';
    }
}
</script>

<style lang="scss" scoped>
.bg-dark-page {
    background: #050505;
}

.bg-dark-soft {
    background: #0d0d0d;
}

.border-right {
    border-right: 1px solid rgba(255, 255, 255, 0.05);
}

.border-bottom {
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.border-left {
    border-left: 1px solid rgba(255, 255, 255, 0.05);
}

.border-lighter {
    border: 1px solid rgba(255, 255, 255, 0.1);
}

.border-bottom-lighter {
    border-bottom: 1px solid rgba(255, 255, 255, 0.02);
}

.bg-primary-opacity {
    background: rgba(var(--q-primary), 0.1);
}

.min-width-auto {
    min-width: unset !important;
}

.size-9 {
    font-size: 9px;
}

.max-width-600 {
    max-width: 600px;
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

.opacity-20 {
    opacity: 0.2;
}
</style>
