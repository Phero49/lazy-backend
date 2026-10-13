<template>
  <div class="col column bg-dark-page" v-if="fields && entity">
    <q-tabs
      @update:model-value="(v) => {
        if (v === 'relationships') {
          loadReferencedIn();
        }
      }"
      v-model="tab"
      :align="'left'"
      inline-label
      no-caps
      shrink
    >
      <q-tab name="table" icon="bi-table" label="Table" />
      <q-tab name="indexes" :icon="biBookHalf" label="Indexes" />
      <q-tab name="relationships" :icon="biNodePlus" label="Relationships" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="table">
        <ColumnsSchema
          ref="columnsSchemaRef"
          :fields="fields ?? []"
          @save-changes="saveFieldsChange"
          @save-to-entity="(value) => addToEntity('fields', value)"
        />
      </q-tab-panel>

      <q-tab-panel name="indexes">
        <IndexesSchema
          :table-indexes="tableIndexes ?? []"
          :fields="fields ?? []"
          @add-reference="addReference"
          @delete-reference="deleteReference"
        />
      </q-tab-panel>

      <q-tab-panel name="relationships">
        <RelationshipsSchema
          :entity="entity"
          :foreign-keys="foreignKeys"
          :referenced-in="referencedIn"
          :referenced-in-column-list="referencedInColumnList"
          @save-foreign-keys="saveForeignKeys"
          @add-foreign-entity="addForeignEntity"
          @delete-foreign-entity="deleteForeignEntity"
          @add-reference="addReference"
          @delete-reference="deleteReference"
          @navigate-column="navigateRelationshipColumn"
        />
      </q-tab-panel>
    </q-tab-panels>
  </div>

  <div v-else-if="loading" class="col column items-center justify-center">
    <q-spinner-grid color="primary" size="lg" />
    <div class="text-grey-7 q-mt-md">Fetching intermediate structure...</div>
  </div>
</template>

<script setup lang="ts">
import {
  ref,
  onMounted,
  watch,
  toRaw,
  onActivated,
  nextTick,
} from "vue";
import { useRoute } from "vue-router";
import { SaveProjectSchema } from "../../wailsjs/go/main/App";
import { useAppContextStore } from "../store/appContext";
import { useQuasar } from "quasar";
import { io } from "../../wailsjs/go/models";
import {
  biBookHalf,
  biNodePlus,
} from "@quasar/extras/bootstrap-icons";
import {
  computeReferencedIn,
  viewColumnColumnSource,
} from "~/utils/databaseUtills";
import ColumnsSchema from "~/components/tableSchema/ColumnsSchema.vue";
import IndexesSchema from "~/components/tableSchema/IndexesSchema.vue";
import RelationshipsSchema from "~/components/tableSchema/RelationshipsSchema.vue";

const route = useRoute();
const store = useAppContextStore();
const $q = useQuasar();
const connectionName = ref(route.params.connection as string);
const tableName = ref(route.params.table as string);
const fields = ref<io.Field[] | null | undefined>(null);
const tableIndexes = ref<io.Index[]>();
const entity = ref<io.Entity>();
const loading = ref(true);
const columnsSchemaRef = ref<InstanceType<typeof ColumnsSchema> | null>(null);
const tab = ref("table");
const referencedInColumnList = [
  "constraintName",
  "localField",
  "sourceField",
  "sourceTable",
  "onDelete",
  "onUpdate",
];

const foreignKeys = ref<io.ForeignKeys[]>([]);
const referencedIn = ref<io.Entity["referencedIn"]>([]);

const navigateRelationshipColumn = (
  key: string,
  value: string | undefined,
  tableNameOverride?: string,
  columnOverride?: string,
) => {
  const targetColumn = columnOverride || value;
  const targetTable = tableNameOverride || tableName.value;

  if (!targetColumn || !targetTable) {
    return;
  }

  tab.value = "table";
  viewColumnColumnSource(

    `/database/schema/${connectionName.value}/${targetTable}`,
    targetTable,
    targetColumn,
  );
};

const loadReferencedIn = () => {
  const entities = store.databaseSchemas[connectionName.value].schema.entities;
  referencedIn.value = computeReferencedIn(entities, tableName.value);
};
onActivated(loadReferencedIn);

async function loadTableSchema() {
  if (!store.currentProject || !connectionName.value || !tableName.value)
    return;

  loading.value = true;
  try {
    const schema = store.databaseSchemas[connectionName.value];
    const table = schema.schema.entities.find(
      (entity: { name: string }) => entity.name === tableName.value,
    );

    if (table == undefined) {
      return;
    }

    if (!entity.value) {
      const dbSchema = store.databaseSchemas[connectionName.value];
      if (dbSchema?.schema?.entities) {
        const table = dbSchema.schema.entities.find(
          (entity: { name: string }) => entity.name === tableName.value,
        );
        if (table == null) {
          return;
        }
        //add empty first row for index

        //  entity.value = table;
        //   console.log(entity.value?.indexes);
      }
    }
    if (table?.ForeignKeys == null) {
      table!.ForeignKeys = [
        {
          name: "none",
          localField: "none",
          referencedTable: "none",
          referencedField: "none",
          onDelete: "none",
          onUpdate: "none",
        },
      ];
    }
    if (table && table?.indexes == undefined) {
      table.indexes = [
        {
          column: "-",
          indexType: "default",
          name: "-",
          unique: true,
        },
      ];
    }

    foreignKeys.value = table?.ForeignKeys || [];
    tableIndexes.value = table?.indexes || [];
    entity.value = table;

    fields.value = entity.value?.fields.map(
      (v: {
        type: string;
        nullish: boolean;
        autoGen: boolean;
        isPrimaryKey: boolean;
        attributes: string;
        default: string;
        name: string;
        length: number;
      }) => {
        const isbool = (t: any) => typeof t === "boolean";
        if (!isbool(v.nullish)) {
          v.nullish = false;
        }
        if (!isbool(v.autoGen)) {
          v.autoGen = false;
        }
        if (!isbool(v.isPrimaryKey)) {
          v.isPrimaryKey = false;
        }
        if (!v.attributes) {
          v.attributes = "-";
        }
        if (!v.default) {
          v.default = "-";
        }
        if (!v.name) {
          v.name = "-";
        }

        if (!v.length) {
          v.length = 0;
        }
        return v;
      },
    );
  } catch (error: any) {
    $q.notify({ type: "negative", message: `Failed to load schema: ${error}` });
  } finally {
    loading.value = false;
  }
}

function addToEntity<T extends keyof io.Entity>(
  key: keyof io.Entity,
  value: io.Entity[T],
) {
  if (!entity.value) {
    return;
  }
  entity.value[key] = value as any;
  saveSchema();
}

function saveFieldsChange() {
  saveSchema();
}

function deleteForeignEntity(index: number) {
  if (!entity.value) return;
  if (entity.value.ForeignKeys) {
    entity.value.ForeignKeys.splice(index, 1);
    saveForeignKeys();
  }
}

function addForeignEntity() {
  if (!entity.value) return;

  if (!entity.value.ForeignKeys) {
    entity.value.ForeignKeys = [];
  }

  entity.value.ForeignKeys.push({
    name: "none",
    localField: "none",
    referencedTable: "none",
    referencedField: "none",
    onDelete: "none",
    onUpdate: "none",
  });

  foreignKeys.value = entity.value.ForeignKeys;
  saveForeignKeys();
}

function deleteReference(index: number) {
  if (!entity.value) return;
  entity.value.referencedIn?.splice(index, 1);
  saveSchema();
}

function addReference() {
  if (!entity.value) return;
  if (!entity.value.referencedIn) {
    entity.value.referencedIn = [];
  }
  entity.value.referencedIn.push(new io.Reference());
  saveSchema();
}

async function saveForeignKeys() {
  try {
    const schema = store.databaseSchemas[connectionName.value].schema;
    const schemaEntity = schema.entities.find((v) => v.name == tableName.value);
    if (schemaEntity == undefined) {
      return;
    }
    console.log(entity.value?.ForeignKeys, "this is what is being saved");
    schemaEntity.ForeignKeys = entity.value?.ForeignKeys || [];

    const res = await SaveProjectSchema(
      store.currentProject || "",
      toRaw(store.databaseSchemas),
    );
    if (res.success) {
      $q.notify({
        color: "green-8",
        message: "Schema saved successfully",
        position: "bottom-right",
        icon: "check",
      });
    } else {
      $q.notify({ type: "negative", message: res.message });
    }
  } catch (error) {
    $q.notify({ type: "negative", message: `Failed to save schema: ${error}` });
  }
}

async function saveSchema() {
  if (!store.currentProject || !entity.value || !fields.value) return;

  try {
    const schema = store.databaseSchemas[connectionName.value];

    if (!schema.schema) {
      schema.schema = new io.Schema({ entities: [], relationships: [] });
    }

    if (!schema.schema.entities) {
      schema.schema.entities = [];
    }

    const entityIndex = schema.schema.entities.findIndex(
      (e: { name: string }) => e.name === tableName.value,
    );
    if (entityIndex !== -1) {
      schema.schema.entities[entityIndex] = entity.value;
    }

    const res = await SaveProjectSchema(
      store.currentProject,
      store.databaseSchemas,
    );
    if (res.success) {
      $q.notify({
        color: "green-8",
        message: "Schema saved successfully",
        position: "bottom-right",
        icon: "check",
      });
    } else {
      $q.notify({ type: "negative", message: res.message });
    }
  } catch (error) {
    $q.notify({ type: "negative", message: `Failed to save schema: ${error}` });
  }
}

onMounted(loadTableSchema);
watch(
  () => [route.params.connection, route.params.table],
  () => {
    connectionName.value = route.params.connection as string;
    tableName.value = route.params.table as string;
    fields.value = null;
    loadTableSchema();
  },
);
watch(
  () => route.query.highlightColumn,
  (column) => {
    const highlightColumn = Array.isArray(column) ? column[0] : column;
    if (!highlightColumn) {
      return;
    }

    nextTick(() => {
      tab.value = "table";
      columnsSchemaRef.value?.highlightColumn(String(highlightColumn));
    });
  },
  { immediate: true },
);
</script>

<style lang="scss" scoped>
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

.highlight {
  background: #19db7aa9 !important;
  animation: highlight 1s ease-in-out;
}

@keyframes highlight {
  0% {
    background-color: #19db7aa9;
  }

  100% {
    background-color: transparent;
  }
}

.underline {
  text-decoration: underline;
  text-decoration-style: dotted;
}
</style>
