<template>
    {{ $route.fullPath }}
      
        <div  class="q-pa-lg">
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

          <div ref="fieldTableContainer"style="overflow-x: scroll; width: 99%" class="custom-scroll">
            <q-markup-table
              bordered
              style=""
              :wrap-cells="true"
              :separator="'cell'"
              class="bg-dark-soft q-mr-md no-shadow"
            >
              <thead class="bg-grey-9">
                <tr>
                  <th
                    class="text-center"
                    :key="key.value"
                    v-for="key in columns"
                  >
                    {{ key.value }}
                    <q-tooltip>
                      {{ key.tooltip }}
                    </q-tooltip>
                  </th>
                </tr>
              </thead>
              <tbody v-if="fields">
                <tr
                  v-for="(value, i) in filteredFields"
                  :key="i"
                  :id="`row-${i}`"
                  :data-field="value.name"
                >
                  <td
                    class="text-center"
                    v-for="{ tooltip, value: key } in columns"
                  >
                    <div v-if="key == 'actions'" style="min-width: 140px">
                      <div class="row q-gutter-x-sm no-wrap">
                        <q-card-actions align="center">
                          <q-btn
                            color="negative"
                            dense
                            icon="delete"
                            size="sm"
                            flat
                            @click="deleteField(i)"
                          />
                          <q-btn
                            flat
                            size="sm"
                            @click="addField(i)"
                            dense
                            :icon="symRoundedAddRowAbove"
                          >
                            <q-tooltip> add row above </q-tooltip>
                          </q-btn>
                          <q-btn
                            flat
                            size="sm"
                            dense
                            @click="addField(i + 1)"
                            :icon="symRoundedAddRowBelow"
                          >
                            <q-tooltip> add row below </q-tooltip>
                          </q-btn>
                        </q-card-actions>
                      </div>
                    </div>

                    <div v-else>
                      <table-cell-editor
                        @change="
                          () => {
                            emits('saveChanges')
                          }
                        "
                        :value="fields[i][key as keyof io.Field]"
                        :read-only="
                          key === 'length' && fields[i]['type'] == 'datetime'
                        "
                        @update:value="
                          (v) => {
                            fields![i][key as keyof io.Field] = v as never;
                          }
                        "
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
                    </div>
                  </td>
                </tr>
              </tbody>
            </q-markup-table>
          </div>
        </div>

</template>

<script lang="ts" setup>
import { io } from '#/wailsjs/go/models';
import { symRoundedAddRowAbove, symRoundedAddRowBelow } from '@quasar/extras/material-symbols-rounded';
import { computed, ref } from 'vue';
import TableCellEditor from "../TableCellEditor.vue";

const props =  defineProps<{fields:io.Field[]}>()
const emits =defineEmits<{saveToEntity:[io.Field[]],saveChanges:[]}>()
const fieldTableContainer = ref<HTMLElement | null>(null);
const filterColumn = ref();
const highlightTimeout = ref<number | null>(null);

defineExpose({
  fieldTableContainer,
  highlightColumn: (columnName: string) => {
    debugger
    if (!fieldTableContainer.value || !columnName) {
      return;
    }

    if (highlightTimeout.value) {
      window.clearTimeout(highlightTimeout.value);
      highlightTimeout.value = null;
    }

    const rows = fieldTableContainer.value.querySelectorAll<HTMLTableRowElement>("tr[data-field]");
    rows.forEach((row) => row.classList.remove("highlight"));

    const targetRow = Array.from(rows).find(
      (row) => row.dataset.field?.toLowerCase() === columnName.toLowerCase(),
    );
console.log(targetRow,'ffff')
    if (!targetRow) {
      return;
    }

    targetRow.classList.add("highlight");
    highlightTimeout.value = window.setTimeout(() => {
      targetRow.classList.remove("highlight");
    }, 5000);
  },
});

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
  if (!props.fields) return [];
  if (!filterColumn.value) return props.fields;
  const search = filterColumn.value.toLowerCase();
  return props.fields.filter((f: { name: string }) =>
    f.name.toLowerCase().includes(search),
  );
});

function deleteField(index: number) {
  if (!props.fields) return;
  props.fields.splice(index, 1);
  emits('saveToEntity',props.fields)
 // addToEntity("fields", fields.value);
}
function addField(index: number) {
  if (!props.fields) return;
  const newField: io.Field = {
    attributes: "-",
    autoGen: false,
    isPrimaryKey: false,
    default: "-",
    length: 0,
    nullish: true,
    name: "new_field_" + props.fields.length,
    type: "string",
  };

 props.fields.splice(index, 0, newField);
    emits('saveToEntity',props.fields)

}

// function saveFieldsChange() {
//   saveSchema();
// }

</script>

<style>

</style>