<template>
  <div class="col column bg-dark-page" v-if="fields && entity">
    <q-tabs
      v-model="tab" :align="'left'"
     inline-label no-caps shrink
    >
      <q-tab name="table" icon="bi-table" label="Table" />
      <q-tab name="indexes" :icon="biBookHalf" label="Indexes" />
      <q-tab name="relationships" :icon="biNodePlus" label="Relationships" />
    </q-tabs>
  <q-tab-panels v-model="tab" animated>
    <q-tab-panel name="table">
     <div ref="fieldTableContainer" class="q-pa-lg">
      <div class="q-mb-md">
        <q-input
          dense
          v-model="filterColumn"
          placeholder="Filter fields by name..."
          dark
          standout
          class="table-filter-input"
        >
          <template v-slot:prepend>
            <q-icon name="bi-search" size="12px" color="grey-8" />
          </template>
        </q-input>
      </div>
      <q-markup-table
        bordered
        style=""
        :wrap-cells="true"
        :separator="'cell'"
        class="bg-dark-soft q-mr-md no-shadow"
      >
        <thead class="bg-grey-9">
          <tr>
            <th class="text-center" :key="key.value" v-for="key in columns">
              {{ key.value }}
              <q-tooltip>
                {{ key.tooltip }}
              </q-tooltip>
            </th>
          </tr>
        </thead>
        <tbody v-if="fields">
          <tr v-for="(value, i) in filteredFields" :key="i" :id="`row-${i}`">
            <td class="text-center" v-for="{ tooltip, value: key } in columns">
              <div v-if="key == 'actions'">
                <div class="row q-gutter-x-sm no-wrap">
                  <q-btn
                    color="negative"
                    dense
                    rounded
                    icon="delete"
                    size="xs"
                    flat
                    @click="deleteField(i)"
                  />
                  <q-btn-dropdown
                    dense
                    rounded
                    dropdown-icon="add_circle"
                    flat
                    size="xs"
                  >
                    <q-list dense>
                      <q-item clickable v-close-popup @click="addField(i + 1)">
                        <q-item-section>
                          <q-item-label>Add row below</q-item-label>
                        </q-item-section>
                      </q-item>
                      <q-item clickable v-close-popup @click="addField(i)">
                        <q-item-section>
                          <q-item-label>Add row above</q-item-label>
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-btn-dropdown>
                </div>
              </div>
               
              <table-cell-editor
              @change="()=>{
                console.log(fields![i][key as keyof io.Field])
              saveFieldsChange()
              }"
                :value="
                 fields[i][key as keyof io.Field]
                "
                @update:value="(v)=>{
                  fields![i][key as keyof io.Field] = v as never
                }"
                :label="key"
                :input-type="key == 'length' ? 'number' : 'text'"
                :type="
                  (() => {
                    if (
                      ['nullish', 'isprimarykey', 'autogen'].includes(
                        key.toLowerCase(),
                      )
                    ) {
                      return 'checkbox';
                    } else if (['attributes', 'type'].includes(key)) {
                      return 'select';
                    } else {
                      return 'input';
                    }
                  })()
                "
                :options="
                  (() => {
                    if (key == 'type') {
                      return availableDataTypes;
                    } else if (key == 'attributes') {
                      //TODO fix inject available attributes
                      return [];
                    }
                  })()
                "
              >
              </table-cell-editor>
            </td>
          </tr>
        </tbody>
      </q-markup-table>
    </div>
    </q-tab-panel>
    <q-tab-panel name="indexes">
     <div class="q-my-md">
        <div class="text-subtitle1 text-bold q-mb-md">Indexes</div>
        <q-markup-table :separator="'cell'">
          <thead>
            <tr>
              <th class="text-center">Actions</th>
              <th
                class="text-left"
                v-for="value in Object.keys(tableIndexes![0] || {})"
              >
                {{ value }}
              </th>
            </tr>
          </thead>
          <tbody v-if="tableIndexes">
            <tr v-for="(value, i) in tableIndexes || []" :key="i">
              <td class="text-center">
                <q-btn flat dense color="grey"  size="sm" icon="bi-three-dots">
                  <q-menu>
                    <q-list style="min-width: 100px">
                      <q-item
                        clickable
                        @click="deleteReference(i)"
                        v-close-popup
                      >
                        <q-item-section>Delete reference</q-item-section>
                      </q-item>
                      <q-separator />
                      <q-item clickable @click="addReference" v-close-popup>
                        <q-item-section>Add reference</q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
              </td>
              <td
                style="cursor: pointer"
                class="text-left"
                v-for="key in Object.keys(value)"
                :key="key"
              >
                <table-cell-editor
                  :value="tableIndexes![i][key as keyof io.Index] || ''"
                  :label="key"
                  v-model="tableIndexes![i][key as keyof io.Index]"
                  :type="
                    key == 'indexType'
                      ? 'select'
                      : key == 'unique'
                        ? 'checkbox'
                        : 'input'
                  "
                  :options="
                    key == 'indexType'
                      ? ['default', 'btree', 'hash', 'gin', 'gist']
                      : []
                  "
                />
              </td>
            </tr>
          </tbody>
        </q-markup-table>
      </div>
    </q-tab-panel>
    <q-tab-panel name="movies">
        <div class="q-my-md">
        <div class="text-subtitle1 text-bold q-mb-md">Foreign Keys</div>
        <q-markup-table
          bordered
          style=""
          :wrap-cells="true"
          class="no-shadow"
          :separator="'cell'"
        >
          <thead>
            <tr>
              <th class="text-center">Actions</th>
              <th
                class="text-left"
                v-for="key in Object.keys(foreignKeys[0] || {})"
                :key="key"
              >
                {{ key }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(foreignEntity, i) in foreignKeys || []">
              <td class="text-center">
                <q-btn flat dense color="grey" size="sm" icon="bi-three-dots">
                  <q-menu>
                    <q-list style="min-width: 100px">
                      <q-item
                        clickable
                        @click="deleteForeignEntity(i)"
                        v-close-popup
                      >
                        <q-item-section>Delete relation</q-item-section>
                      </q-item>
                      <q-separator />
                      <q-item clickable @click="addForeignEntity" v-close-popup>
                        <q-item-section>Add relation</q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
              </td>
              <td
                class="text-left"
                v-for="key in Object.keys(foreignEntity)"
                :key="key"
              >
                <table-cell-editor
                  outlined
                  dense
                  :value="foreignEntity[key as keyof io.ForeignKeys] || ''"
                  :label="key"
                  v-model="foreignEntity[key as keyof io.ForeignKeys]"
                  :type="key == 'name' ? 'input' : 'select'"
                  :options="
                    (() => {
                      if (key == 'onDelete' || key == 'onUpdate') {
                        return ['cascade', 'restrict', 'setNull', 'noAction'];
                      } else if (key == 'referencedTable') {
                        return useOptions().getTables;
                      } else if (key == 'localField') {
                        return entity.fields
                          .map((v) => v.name)
                          .filter((v) => v != '' && v != '-');
                      } else if (key == 'referencedField') {
                        return useOptions().getTableColumns(
                          foreignEntity['referencedTable'],
                        );
                      }
                      return [];
                    })()
                  "
                />
              </td>
            </tr>
          </tbody>
        </q-markup-table>
      </div>
      <div class="q-my-md">
        <div class="text-subtitle1 text-bold q-mb-md">
          Referenced In (Tables that reference this table)
        </div>
        <q-markup-table separator="cell" bordered class="no-shadow">
          <thead>
            <tr>
              <th class="text-center">Actions</th>

              <th
                class="text-left"
                v-for="item in Object.keys(
                  (entity.referencedIn || [])[0] || {},
                )"
              >
                {{ item }}
              </th>
            </tr>
          </thead>
          <tbody v-if="entity.referencedIn">
            <tr v-for="(reference, i) in entity?.referencedIn || []" :key="i">
              <td class="text-center">
                <q-btn flat dense color="grey" icon="bi-three-dots">
                  <q-menu>
                    <q-list style="min-width: 100px">
                      <q-item
                        clickable
                        @click="deleteReference(i)"
                        v-close-popup
                      >
                        <q-item-section>Delete reference</q-item-section>
                      </q-item>
                      <q-separator />
                      <q-item clickable @click="addReference" v-close-popup>
                        <q-item-section>Add reference</q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
              </td>
              <td class="" v-for="key in Object.keys(reference)" :key="key">
                <table-cell-editor
                  :value="entity.referencedIn[i][key as keyof io.Reference]"
                  v-model="entity.referencedIn[i][key as keyof io.Reference]"
                  :type="'constraintName' == key ? 'input' : 'select'"
                  :options="
                    (() => {
                      if (key == 'onDelete' || key == 'onUpdate') {
                        return ['cascade', 'restrict', 'setNull', 'noAction'];
                      } else if (key == 'referencedTable') {
                        return useOptions().getTables;
                      } else if (key == 'localField') {
                        return entity.fields
                          .map((v) => v.name)
                          .filter((v) => v != '' && v != '-');
                      } else if (key == 'referencedField') {
                        return useOptions().getTableColumns(
                          entity.referencedIn[i]['sourceTable'],
                        );
                      }
                      return [];
                    })()
                  "
                >
                </table-cell-editor>
              </td>
            </tr>
          </tbody>
        </q-markup-table>
      </div>
    </q-tab-panel>
  </q-tab-panels>
  
  </div>

  <div v-else-if="loading" class="col column items-center justify-center">
    <q-spinner-grid color="primary" size="lg" />
    <div class="text-grey-7 q-mt-md">Fetching intermediate structure...</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import { useRoute } from "vue-router";
import {
  GetProjectEnv,
  GetTableSchema,
  GetProjectSchema,
  SaveProjectSchema,
} from "../../wailsjs/go/main/App";
import { useAppContextStore } from "../store/appContext";
import { useQuasar } from "quasar";
import { io } from "../../wailsjs/go/models";
import TableSchemaGraph from "../components/TableSchemaGraph.vue";
import { forInStatement } from "@babel/types";
import TableCellEditor from "../components/TableCellEditor.vue";
import { useOptions } from "../store/optionsRegistry";
import { biBookHalf, biNodePlus, biNodePlusFill } from "@quasar/extras/bootstrap-icons";

const route = useRoute();
const store = useAppContextStore();
const $q = useQuasar();

const connectionName = ref(route.params.connection as string);
const tableName = ref(route.params.table as string);
const fields = ref<io.Field[] | null | undefined>(null);
const tableIndexes = ref<io.Index[]>();
const entity = ref<io.Entity>();
const loading = ref(true);
const filterColumn = ref();
const fieldTableContainer = ref<HTMLDivElement>();
const tab = ref('table')
const dTypeFilter = ref("");
const availableDataTypes = ref([
  "string",
  "integer",
  "float",
  "boolean",
  "datetime",
  "time",
  "date",
  "json",
  "array",
  "enum",
]);

const columns = [
  { value: "actions", tooltip: "table actions" },
  { value: "name", tooltip: "Column name identifier" },
  { value: "type", tooltip: "Data type" },
  {
    value: "length",
    tooltip: "For string: max characters, For float: precision digits",
  },
  {
    value: "isPrimaryKey",
    tooltip: "Set as primary key for unique row identification",
  },
  {
    value: "autoGen",
    tooltip: "Auto-incrementing value (typically for ID columns)",
  },
  { value: "nullish", tooltip: "NULL or NOT NULL - allows empty values" },
  { value: "default", tooltip: "Default value when none provided" },
  { value: "attributes", tooltip: "Additional database-specific attributes" },
];

const filteredFields = computed(() => {
  if (!fields.value) return [];
  if (!filterColumn.value) return fields.value;
  const search = filterColumn.value.toLowerCase();
  return fields.value.filter((f) => f.name.toLowerCase().includes(search));
});

const foreignKeys = ref<io.ForeignKeys[]>([]);

async function loadTableSchema() {
  if (!store.currentProject || !connectionName.value || !tableName.value)
    return;

  loading.value = true;
  try {
    const schema = store.databaseSchemas[connectionName.value];
    if (schema?.schema?.entities) {
      entity.value = schema.schema.entities.find(
        (entity) => entity.name === tableName.value,
      );
    }

    if (!entity.value) {
      const dbSchema = store.databaseSchemas[connectionName.value];
      if (dbSchema?.schema?.entities) {
        const table = dbSchema.schema.entities.find(
          (entity) => entity.name === tableName.value,
        );
        if (table == null) {
          return;
        }
        //add empty first row for index

        //  entity.value = table;
        //   console.log(entity.value?.indexes);
      }
    }
    foreignKeys.value = entity.value?.ForeignKeys || [
      {
        name: "none",
        localField: "none",
        referencedTable: "none",
        referencedField: "none",
        onDelete: "none",
        onUpdate: "none",
      },
    ];
    tableIndexes.value = entity.value?.indexes || [
      {
        column: "-",
        indexType: "default",
        name: "-",
        unique: true,
      },
    ];

    entity.value!.referencedIn ??= [
      {
        constraintName: "none",
        localField: "none",

        sourceField: "none",
        sourceTable: "none",
        onDelete: "none",
        onUpdate: "none",
      },
    ];

    fields.value = entity.value?.fields.map((v)=>{
      const isbool = (t:any)=> typeof t === 'boolean'
      if (!isbool(v.nullish)) {
        v.nullish  =  false
      }
      if (!isbool(v.autoGen)  ) {
        v.autoGen = false
      }
      if (!isbool(v.isPrimaryKey)) {
        v.isPrimaryKey = false
      }
      if (!v.attributes) {
        v.attributes = '-'
      } 
        if (!v.default) {
        v.default = '-'
      } 
      if (!v.name) {
        v.name = '-'
      }
      
      if (!v.length) {
        v.length = 0
      }
      return v
    });
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

function addField(index: number) {
  if (!fields.value) return;
  const newField:io.Field =  {
    attributes:'-',
    autoGen:false,
    isPrimaryKey:false,
    default:'-',
    length:0,
    nullish:true,
    name:"new_field_" + fields.value.length,
    type:'string'
  }


  fields.value.splice(index, 0, newField);
  addToEntity("fields", fields.value);
}

function deleteField(index: number) {
  if (!fields.value) return;
  fields.value.splice(index, 1);
  addToEntity("fields", fields.value);
}

function saveFieldsChange (){
   saveSchema()
}


function deleteForeignEntity(index: number) {
  if (!entity.value) return;
}

function goToField(field: string) {
  if (fieldTableContainer.value) {
    const found = fieldTableContainer.value.querySelector(
      `[col-value="${field}"]`,
    );
    if (found) {
      const rowIndex = found.getAttribute("row-index");
      if (rowIndex) {
        const row = fieldTableContainer.value.querySelector(`#row-${rowIndex}`);
        if (row) {
          row.classList.add("highlight");
          setTimeout(() => {
            row.classList.remove("highlight");
          }, 1000);
          row.scrollIntoView({ behavior: "smooth", block: "center" });
        }
      }
    }
  }
}

function addForeignEntity() {}

function deleteReference(index: number) {
  if (!entity.value) return;
  entity.value.referencedIn?.splice(index, 1);
}

function addReference() {
  if (!entity.value) return;
  if (!entity.value.referencedIn) {
    entity.value.referencedIn = [];
  }
  entity.value.referencedIn.push(new io.Reference());
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
      (e) => e.name === tableName.value,
    );
    if (entityIndex !== -1) {
      schema.schema.entities[entityIndex] = entity.value;
    }

    const res = await SaveProjectSchema(
      store.currentProject,
      store.databaseSchemas,
    );
    if (res.success) {
      $q.notify({ color:'green-8', message: "Schema saved successfully",position:'bottom-right',icon:'check' });
    } else {
      $q.notify({ type: "negative", message: res.message });
    }
  } catch (error) {
    $q.notify({ type: "negative", message: `Failed to save schema: ${error}` });
  }
}

function filterTable(val: string, update: Function, abort: Function) {
  update(() => {
    const needle = val.toLowerCase();
    return getTableList().filter((v) => v.toLowerCase().includes(needle));
  });
}

function getTableList() {
  return (
    store.databaseSchemas[store.activeConnectionName || ""].schema.entities.map(
      (e) => e.name,
    ) || []
  ).sort();
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
