<template>
  <div class="row q-gutter-sm items-center">
    <div class="col-2">
      <q-select v-model="model.source" :options="dataSource" label="Source" filled dense
        @update:model-value="model.field = ''" />
    </div>
    <div class="col-3">
      <q-select v-model="model.field" :options="fieldOptions" label="Field" filled dense
        :disable="!model.source" />
    </div>
    <div class="col-2">
      <q-select v-model="model.operator" :options="operatorOptions" label="Operator" filled dense
        @update:model-value="onOperatorChange" />
    </div>
    <div class="col-3" v-if="needsValue">
      <q-input v-model="model.value" label="Value" filled dense />
    </div>
    <div class="col-auto">
      <q-btn flat round dense icon="close" @click="emit('remove')" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";

export interface ConditionRowModel {
  type: "condition";
  source: string;
  field: string;
  operator: string;
  value: string | number | null;
}

defineProps<{
  level: number;
  fieldOptions?: (string | number)[];
}>();

const emit = defineEmits<{ (e: "remove"): void }>();

// self-syncing v-model — no manual watch/emit needed
const model = defineModel<ConditionRowModel>({
  default: () => ({ type: "condition", source: "", field: "", operator: "", value: null }),
});

const dataSource = ["static", "payload"];
const operatorOptions = ["equals", "not_equals", "contains", "is_null", "is_empty", "gt", "lt", "gte", "lte"];
const needsValue = computed(() => !["is_null", "is_empty"].includes(model.value.operator));

function onOperatorChange() {
  if (!needsValue.value) model.value.value = null;
}
</script>