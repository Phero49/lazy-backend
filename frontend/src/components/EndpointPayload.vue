<template>
  <div>
    <q-list class="q-col-gutter-y-md q-pa-sm">
      <q-item
        v-for="(item, index) in payload || []"
        :key="index"
        class="row items-center q-py-sm rounded-borders bg-dark-page"
        style="border: 1px solid rgba(255, 255, 255, 0.05)"
      >
        <div class="col-auto q-mr-md">
          <TableCellEditor
            :value="item.dataSource"
            :type="'select'"
            :options="getDataSource()"
            @update:value="updatePayload($event, index, 'dataSource')"
            class="text-uppercase text-bold"
          />
        </div>

        <div class="col q-mr-md">
          <TableCellEditor
            type="input"
            placeholder="field_name"
            :value="item.fieldName"
            @update:value="updatePayload($event, index, 'fieldName')"
          />
        </div>

        <div class="col-auto q-mr-lg">
          <TableCellEditor
            type="select"
            :options="['string', 'int', 'float', 'file', 'array', 'object']"
            :value="item.dataType"
            @update:value="updatePayload($event, index, 'dataType')"
          />
        </div>

        <div class="col-auto q-mr-md row items-center">
          <span class="text-caption text-grey-5 q-mr-sm">Required</span>
          <TableCellEditor
            type="checkbox"
            :value="item.required"
            @change="updatePayload($event, index, 'required')"
          />
        </div>

        <div class="col-auto">
          <q-btn
            flat
            round
            color="negative"
            icon="bi-trash"
            size="sm"
            @click="removeSpecificRow(index)"
          />
        </div>
      </q-item>
    </q-list>

    <div class="text-center q-mt-md q-px-sm">
      <q-btn
        color="green-10"
        label="Add Request Field"
        :icon="biPlus"
        no-caps
        @click="modifyRows(true)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { biPlus } from "@quasar/extras/bootstrap-icons";
import { Payload } from "../store/appContext";
import TableCellEditor from "./TableCellEditor.vue";
import { computed } from "vue";

const props = defineProps<{ payload: Payload[]; method: string }>();
function updatePayload(
  p: string | number | boolean,
  index: number,
  column: keyof Payload,
) {
  props.payload[index][column] = p as never;
}
const method = computed(() => props.method);

function removeSpecificRow(index: number) {}
function getDataSource() {
  const sources = ["form", "param", "query"];
  if (["PUT", "POST"].includes(method.value!)) {
    sources.unshift("body");
  }
  return sources;
}

const modifyRows = (add: boolean) => {
  if (add) {
    props.payload.push({
      dataSource: "-",
      dataType: "string",
      fieldName: "-",
      required: false,
    });
  } else {
    props.payload.pop();
  }
};
</script>

<style scoped></style>
