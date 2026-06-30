<template>
  <q-card
    flat
    bordered
    class="q-pa-sm"
    style="min-height: 380px; min-width: 650px"
  >
    <q-markup-table bordered separator="cell" class="no-shadow">
      <thead>
        <tr>
          <th class="text-left">operator</th>
          <th class="text-left">Column</th>
          <th class="text-center">comparison</th>
          <th class="text-center">Value</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, i) in conditions">
          <td>
            <q-select
              v-if="item.operator"
              v-model="item.operator"
              :options="['AND', 'OR']"
              label=""
              borderless
            />
          </td>
          <td>
            <div class="q-gutter-y-xs">
              <q-select
                style="min-width: 130px"
                v-model="item.column.table"
                dense
                :options="tables"
                label="select table"
                filled
              /><q-select
                style="min-width: 130px"
                v-model="item.column.name"
                dense
                label="select column"
                :options="options.getTableColumns(item.column.table)"
                filled
              />
            </div>
          </td>
          <td class="text-right" style="max-width: 150px">
            <q-select
              style="min-width: 130px"
              v-model="item.comparison"
              dense
              :options="supportedOperators"
              filled
            />
          </td>
          <td style="min-width: 150px">
            <div class="q-gutter-y-xs" v-if="compareBy == 'column'">
              <q-select
                style="min-width: 130px"
                v-model="item.value.column.table"
                dense
                :options="tables"
                label="select table"
                filled
              >
                <template #append>
                  <q-btn
                    color="red"
                    dense
                    flat
                    icon="close"
                    size="sm"
                    @click="resetValue(i)"
                  />
                </template> </q-select
              ><q-select
                style="min-width: 130px"
                v-model="item.value.column.name"
                dense
                label="select column"
                :options="options.getTableColumns(item.value.column.table)"
                filled
              >
              </q-select>
            </div>
            <div v-else-if="compareBy == 'payload'">
              <q-select
                v-model="item.value.payload"
                dense 
                :options="
                  payload.map((v) => {
                    return v.fieldName;
                  })
                "
                label="payload."
                filled
              >
                <template #after>
                  <q-btn
                    color="red"
                    dense
                    flat
                    icon="close"
                    size="sm"
                    @click="resetValue(i)"
                  />
                </template>
                <template #no-option>
                  <div class="text-center text-subtitle2 q-pa-md">
                    payload data schema not defined
                  </div>
                </template>
              </q-select>
            </div>
            <div v-else-if="compareBy == 'value'">
              <q-input
                dense
                v-model="item.value.value"
                type="text"
                label="write here"
                borderless
                ><template #after>
                  <q-btn
                    color="red"
                    dense
                    flat
                    icon="close"
                    size="sm"
                    @click="resetValue(i)"
                  /> </template
              ></q-input>
            </div>
            <div v-else>
              <q-select
                v-model="compareBy"
                dense
                :options="compareByOptions"
                label="value source"
                borderless
              >
              </q-select>
            </div>
          </td>
        </tr>
      </tbody>
    </q-markup-table>
    <div class="text-center q-gutter-x-sm q-mt-sm">
      <q-btn
        color="primary"
        dense
        icon="add"
        padding="0px 10px"
        @click="addCondition"
      />
      <q-btn
        color="red"
        dense
        :icon="biDash"
        padding="0px 10px"
        @click="removeCondition"
      />
    </div>
  </q-card>
</template>

<script lang="ts" setup>
import { onMounted, onBeforeUnmount, ref } from "vue";
import { useOptions } from "../store/optionsRegistry";
import { biDash } from "@quasar/extras/bootstrap-icons";
import { Payload } from "../store/appContext";
import type { Condition, SupportedOperators } from "..";

const props = defineProps<{
  tables: string[];
  payload: Payload[];
  condition: Condition[];
}>();
const options = useOptions();
const emits = defineEmits<{ close: [Condition[]] }>();
const compareByOptions = ["column", "payload", "value"] as const;
const compareBy = ref<(typeof compareByOptions)[number]>();
const supportedOperators: SupportedOperators[] = [
  "Equal to",
  "Not equal to",
  ">",
  "<",
  ">=",
  "<=",
  "Like",
  "Not Like",
  "In",
  "Not In",
  "starts with",
  "ends with",
];
function addCondition() {
  const payload: Condition = {
    operator: conditions.value.length > 0 ? "AND" : undefined,
    column: {
      table: "",
      name: "",
    },
    comparison: "Equal to",
    value: {
      payload: null,
      value: undefined,
      column: {
        table: "",
        name: "",
      },
    },
  };

  conditions.value.push(payload);
}
function resetValue(i: number) {
  conditions.value[i].value = {
    payload: null,
    value: undefined,
    column: {
      table: "",
      name: "",
    },
  };
  compareBy.value = undefined;
}
function removeCondition() {
  conditions.value.pop();
}
const conditions = ref<Condition[]>(props.condition || []);
onMounted(() => {
  addCondition();
});
onBeforeUnmount(() => {
  emits("close", conditions.value);
});
</script>

<style></style>
