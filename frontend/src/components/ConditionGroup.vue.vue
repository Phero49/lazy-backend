<template>
  <div class="condition-group q-pa-sm" :class="{ 'q-pl-md': level > 0 }"
       :style="level ? 'border-left: 2px solid var(--q-primary)' : ''">
    <div class="row items-center q-gutter-sm q-mb-sm">
      <q-btn-toggle v-model="model.logic" dense no-caps
        :options="[{ label: 'AND', value: 'AND' }, { label: 'OR', value: 'OR' }]" />
      <q-btn flat dense icon="add" label="Condition" @click="addCondition" />
      <q-btn flat dense icon="add" label="Group" @click="addGroup" />
      <q-btn v-if="level" flat round dense icon="close" @click="emit('remove')" />
    </div>

    <div v-for="(child, i) in model.children" :key="i" class="q-mb-sm">
      <ConditionGroup
        v-if="child.type === 'group'"
        :level="level + 1"
        v-model="model.children[i]"
        :field-options="fieldOptions"
        @remove="model.children.splice(i, 1)"
      />
      <ConditionRow
        v-else
        :level="level + 1"
        v-model="(model.children[i] as any)"
        :field-options="fieldOptions"
        @remove="model.children.splice(i, 1)"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import ConditionRow, { type ConditionRowModel } from "./ConditionRow.vue";
// no import needed for self-reference: <script setup> auto-registers
// the component under its own filename for recursive use

export interface ConditionGroupModel {
  type: "group";
  logic: "AND" | "OR";
  children: (ConditionRowModel | ConditionGroupModel)[];
}

defineProps<{
  level: number;
  fieldOptions?: (string | number)[];
}>();

const emit = defineEmits<{ (e: "remove"): void }>();

const model = defineModel<ConditionGroupModel>({
  default: () => ({ type: "group", logic: "AND", children: [] }),
});

function addCondition() {
  model.value.children.push({ type: "condition", source: "", field: "", operator: "", value: null });
}
function addGroup() {
  model.value.children.push({ type: "group", logic: "AND", children: [] });
}
</script>