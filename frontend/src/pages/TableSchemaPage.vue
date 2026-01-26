<template>
    <div class="col column bg-dark-page" v-if="entity">
        <div class="q-pa-md border-bottom bg-dark-soft row items-center justify-between sticky-top">
            <div class="row items-center q-gutter-x-sm">
                <q-icon name="bi-table" color="primary" />
                <span class="text-caption text-weight-bold text-white">{{ entity.name }}</span>
                <q-badge color="grey-10" text-color="grey-6" class="q-ml-sm" label="Intermediate Structure" />
            </div>
            <div class="row items-center q-gutter-x-sm">
                <q-btn flat round dense icon="bi-arrow-clockwise" size="sm" color="grey-7" @click="loadTableSchema" />
            </div>
        </div>

        <q-scroll-area class="col">
            <div class="q-pa-lg">
                <!-- Field Table -->
                <q-table flat bordered dark :rows="entity.fields" :columns="columns" row-key="name" hide-bottom
                    class="bg-dark-soft border-lighter" :pagination="{ rowsPerPage: 0 }">
                    <template v-slot:body-cell-name="props">
                        <q-td :props="props">
                            <div class="row items-center no-wrap">
                                <q-icon v-if="entity.primaryKey.includes(props.value)" name="bi-key-fill"
                                    color="warning" size="12px" class="q-mr-xs" />
                                <span class="text-weight-medium text-white">{{ props.value }}</span>
                            </div>
                        </q-td>
                    </template>

                    <template v-slot:body-cell-type="props">
                        <q-td :props="props">
                            <q-badge outline :label="props.value" color="primary" class="text-uppercase"
                                style="font-size: 10px" />
                        </q-td>
                    </template>

                    <template v-slot:body-cell-length="props">
                        <q-td :props="props">
                            <span v-if="props.value" class="text-grey-6 text-mono">{{ props.value }}</span>
                            <span v-else class="text-grey-9 text-italic">-</span>
                        </q-td>
                    </template>

                    <template v-slot:body-cell-nullable="props">
                        <q-td :props="props">
                            <q-icon :name="props.row.nulish ? 'bi-check-circle' : 'bi-dash-circle'"
                                :color="props.row.nulish ? 'positive' : 'grey-9'" size="14px" />
                        </q-td>
                    </template>

                    <template v-slot:body-cell-default="props">
                        <q-td :props="props">
                            <code v-if="props.value"
                                class="bg-grey-10 q-px-xs border-radius-sm text-grey-5">{{ props.value }}</code>
                            <span v-else class="text-grey-9">-</span>
                        </q-td>
                    </template>
                </q-table>

                <!-- Primary Key Info -->
                <div class="q-mt-lg row items-center q-gutter-x-md text-grey-8">
                    <div class="text-overline">Primary Key:</div>
                    <div class="row q-gutter-xs">
                        <q-chip v-for="pk in entity.primaryKey" :key="pk" dense color="warning-opacity"
                            text-color="warning" class="text-weight-bold">
                            {{ pk }}
                        </q-chip>
                    </div>
                </div>
            </div>
        </q-scroll-area>
    </div>

    <div v-else-if="loading" class="col column items-center justify-center">
        <q-spinner-grid color="primary" size="lg" />
        <div class="text-grey-7 q-mt-md">Fetching intermediate structure...</div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { GetProjectEnv, GetTableSchema } from '../../wailsjs/go/main/App';
import { useAppContextStore } from '../store/appContext';
import { useQuasar } from 'quasar';

const route = useRoute();
const store = useAppContextStore();
const $q = useQuasar();

const connectionName = ref(route.params.connection as string);
const tableName = ref(route.params.table as string);
const entity = ref<any>(null);
const loading = ref(true);

const columns: any[] = [
    { name: 'name', label: 'Field', field: 'name', align: 'left' },
    { name: 'type', label: 'Type', field: 'type', align: 'left' },
    { name: 'length', label: 'Length', field: 'length', align: 'left' },
    { name: 'nullable', label: 'Nullable', field: 'nulish', align: 'center' },
    { name: 'default', label: 'Default', field: 'default', align: 'left' },
];

async function loadTableSchema() {
    if (!store.currentProject || !connectionName.value || !tableName.value) return;

    loading.value = true;
    try {
        const env = await GetProjectEnv(store.currentProject);
        const dbConfigs = env.database || [];
        const profile = dbConfigs.find((c: any) => Object.keys(c)[0] === connectionName.value);

        if (profile) {
            const config = profile[connectionName.value];
            entity.value = await GetTableSchema(
                config.host,
                config.port,
                config.user,
                config.password,
                config.database,
                tableName.value
            );
        }
    } catch (error: any) {
        $q.notify({ type: 'negative', message: `Failed to load schema: ${error}` });
    } finally {
        loading.value = false;
    }
}

onMounted(loadTableSchema);
watch(() => [route.params.connection, route.params.table], () => {
    connectionName.value = route.params.connection as string;
    tableName.value = route.params.table as string;
    loadTableSchema();
});
</script>

<style lang="scss" scoped>
.editor-bg {
    background-color: #1e1e1e !important;
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

.border-lighter {
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.border-radius-md {
    border-radius: 8px;
}

.sticky-top {
    position: sticky;
    top: 0;
    z-index: 10;
}

.warning-opacity {
    background: rgba(var(--q-warning), 0.1) !important;
}
</style>
