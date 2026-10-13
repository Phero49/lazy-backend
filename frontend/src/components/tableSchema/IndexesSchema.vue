<template>
  <div class="q-my-md">
    <div class="text-subtitle1 text-bold q-mb-md">Indexes</div>
    <q-markup-table :separator="'cell'">
      <thead>
        <tr>
          <th class="text-center">Actions</th>
          <th class="text-left" v-for="value in Object.keys(tableIndexes[0] || {})" :key="value">
            {{ value }}
          </th>
        </tr>
      </thead>
      <tbody v-if="tableIndexes?.length">
        <tr v-for="(value, i) in tableIndexes" :key="i">
          <td class="text-center">
            <q-btn flat dense color="grey" size="sm" icon="bi-three-dots">
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
            style="cursor: pointer"
            class="text-left"
            v-for="key in Object.keys(value)"
            :key="key"
          >
            <table-cell-editor
              :value="tableIndexes[i][key as keyof io.Index] || ''"
              :label="key"
              :columns-options="key == 'column'"
              v-model="tableIndexes[i][key as keyof io.Index]"
              :type="
                ['indexType', 'column'].includes(key)
                  ? 'select'
                  : key == 'unique'
                    ? 'checkbox'
                    : 'input'
              "
              :options="getIndexColumnOptions(key)"
            />
          </td>
        </tr>
      </tbody>
    </q-markup-table>
  </div>
</template>

<script setup lang="ts">
import { io } from "#/wailsjs/go/models";
import TableCellEditor from "../TableCellEditor.vue";

const props = defineProps<{
  tableIndexes: io.Index[];
  fields: io.Field[];
}>();

const emit = defineEmits<{
  addReference: [];
  deleteReference: [number];
}>();

const getIndexColumnOptions = (key: string) => {
  switch (key) {
    case "column":
      return props.fields?.map((field) => field.name) || [];
    case "indexType":
      return ["default"];
    default:
      return [];
  }
};
</script>
