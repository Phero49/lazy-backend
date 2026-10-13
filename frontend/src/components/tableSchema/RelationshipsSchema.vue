<template>
  <div class="q-my-md">
    <div class="text-subtitle1 text-bold q-mb-md">Foreign Keys</div>
    <q-markup-table
      bordered
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
      <tbody v-if="foreignKeys?.length">
        <tr v-for="(foreignEntity, i) in foreignKeys" :key="i">
          <td class="text-center">
            <q-btn flat dense color="grey" size="sm" icon="bi-three-dots">
              <q-menu>
                <q-list style="min-width: 100px">
                  <q-item clickable @click="emit('deleteForeignEntity', i)" v-close-popup>
                    <q-item-section>Delete relation</q-item-section>
                  </q-item>
                  <q-separator />
                  <q-item clickable @click="emit('addForeignEntity')" v-close-popup>
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
            @click="
              key === 'localField' && foreignEntity.localField
                ? emit('navigateColumn', key, foreignEntity.localField)
                : key === 'referencedField' && foreignEntity.referencedTable && foreignEntity.referencedField
                  ? emit(
                      'navigateColumn',
                      key,
                      foreignEntity.referencedField,
                      foreignEntity.referencedTable,
                    )
                  : undefined
            "
            style="cursor: pointer;"
          >
            <table-cell-editor
              outlined
              dense
              @change="emit('saveForeignKeys')"
              @update:value="
                (value) => {
                  if (entity?.ForeignKeys[i]) {
                    entity.ForeignKeys[i][key as keyof io.ForeignKeys] =
                      value as string;
                    emit('saveForeignKeys');
                  }
                }
              "
              :value="foreignEntity[key as keyof io.ForeignKeys] || ''"
              :label="key"
              v-model="foreignEntity[key as keyof io.ForeignKeys]"
              :type="key == 'name' ? 'input' : 'select'"
              :options="tabCellEditorReferenceOptions(key, i, foreignEntity)"
            />
          </td>
        </tr>
      </tbody>
    </q-markup-table>

    <div class="q-my-md">
      <div class="text-subtitle1 text-bold q-mb-md">
        Referenced In (Tables that reference this table)
      </div>
      <q-markup-table separator="horizontal" bordered class="no-shadow">
        <thead>
          <tr>
            <th class="text-center">Actions</th>
            <th class="text-left" v-for="item in referencedInColumnList" :key="item">
              {{ item }}
            </th>
          </tr>
        </thead>
        <tbody v-if="referencedIn?.length">
          <tr v-for="(reference, i) in referencedIn" :key="i">
            <td class="text-center">
              <q-btn flat dense color="grey" icon="bi-three-dots">
                <q-menu>
                  <q-list style="min-width: 100px">
                    <q-item clickable @click="emit('deleteReference', i)" v-close-popup>
                      <q-item-section>Delete reference</q-item-section>
                    </q-item>
                    <q-separator />
                    <q-item clickable @click="emit('addReference')" v-close-popup>
                      <q-item-section>Add reference</q-item-section>
                    </q-item>
                  </q-list>
                </q-menu>
              </q-btn>
            </td>
            <td
              class=""
              v-for="(value, key) in reference"
              :key="key"
              @click="
                key === 'localField' && reference.localField
                  ? emit('navigateColumn', key, reference.localField)
                  : key === 'sourceField' && reference.sourceTable && reference.sourceField
                    ? emit('navigateColumn', key, reference.sourceField, reference.sourceTable)
                    : key === 'sourceTable'
                      ? emit('navigateColumn', 'sourceField', reference.sourceField, reference.sourceTable)
                      : undefined
              "
              :style="key === 'sourceField' || key === 'localField' || key === 'sourceTable' ? 'cursor: pointer;' : 'pointer-events: none;'"
            >
              {{ value }}
            </td>
          </tr>
        </tbody>
      </q-markup-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { io } from "#/wailsjs/go/models";
import { useOptions } from "~/store/optionsRegistry";
import TableCellEditor from "../TableCellEditor.vue";

const props = defineProps<{
  entity: io.Entity;
  foreignKeys: io.ForeignKeys[];
  referencedIn: io.Entity["referencedIn"];
  referencedInColumnList: string[];
}>();

const emit = defineEmits<{
  saveForeignKeys: [];
  addForeignEntity: [];
  addReference: [];
  deleteForeignEntity: [number];
  deleteReference: [number];
  navigateColumn: [
    key: string,
    value: string | undefined,
    tableNameOverride?: string,
    columnOverride?: string,
  ];
}>();

const tabCellEditorReferenceOptions = (
  key: string,
  i: number,
  foreignEntity?: io.ForeignKeys,
) => {
  if (key == "onDelete" || key == "onUpdate") {
    return ["cascade", "restrict", "setNull", "noAction"];
  } else if (key == "referencedTable") {
    return useOptions().getTables;
  } else if (key == "localField") {
    return (props.entity?.fields || [])
      .map((v: { name: any }) => v.name)
      .filter((v: string) => v != "" && v != "-");
  } else if (key == "referencedField") {
    const table = foreignEntity?.["referencedTable"];
    if (table) {
      return useOptions().getTableColumns(table);
    }

    const sourceTable = props.entity?.referencedIn?.[i]?.sourceTable;
    if (sourceTable) {
      return useOptions().getTableColumns(sourceTable);
    }
  }
  return [];
};
</script>
