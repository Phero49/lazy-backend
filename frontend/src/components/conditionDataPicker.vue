<template>
  <q-select
    style="min-width: 300px"
    class="q-mb-sm"
    v-if="expressionFragment.source"
    dense
    @update:model-value="loadOptions"
    v-model="expressionFragment.source"
    :options="['static', 'dependencies', 'payload', 'operationBlock']"
    label="data source"
  />
  <!--handles static values-->
  <div class="row" v-if="expressionFragment.source == 'static'">
    <div class="col-4">
      <q-select
        dense
        v-model="expressionFragment.valueType"
        @update:model-value="
          (v) => {
            console.log(v)
            if (v == 'null') {
              expressionFragment.value = 'null';
              search = 'null'
            }
          }
        "
        :options="[
          'null',
          'string',
          'int',
          'list',
          'bool',
          'list (string)',
          'list (int)',
        ]"
        label="type"
        filled
      />
    </div>
    <div class="col-auto" style="cursor: pointer">
      <q-input
        square
        :readonly="expressionFragment.valueType == 'null'"
        @keydown="onEnter"
        v-model="search"
        dense
        filled
        autofocus
        :label="
          expressionFragment.source == 'static'
            ? `write your ${expressionFragment.valueType}`
            : ` filter values in ${expressionFragment?.source}`
        "
        :hint="`Press enter to save. ${expressionFragment.valueType?.includes('list') ? 'separate values by ,' : ''}`"
      />
    </div>
  </div>
  <div v-else class="q-pb-md">
    <q-select
      v-model="search"
      dense
      :options="availableDataForSelectedOp"
      :label="`filter from ${expressionFragment.source}`"
      filled
    />
  </div>
  <q-list
    style="min-width: 220px; max-height: 280px"
    v-if="expressionFragment?.source != 'static'"
  >
  <template v-if="filteredAvailableDataForSelectedOp.length > 0">
  <q-item
      v-for="(item, i) in filteredAvailableDataForSelectedOp"
      :key="i"
      clickable class="text-subtitle2"
      v-close-popup @click="()=>{
        expressionFragment.value = item.label
      }"
    >
      <q-item-section>{{ item.label }}</q-item-section>
      <q-item-section side>
        <q-item-label caption>{{ item.dataType }}</q-item-label>
      </q-item-section>
    </q-item>
  </template>
  <template v-else>  <q-item >
      <q-item-section class="text-grey">Nothing available</q-item-section>
    </q-item></template>
  
  
  </q-list>
</template>

<script lang="ts" setup>
import { inject, onMounted, ref } from "vue";
import { ConditionExpressionPart, Routes } from "~/store/appContext";

const props = defineProps<{ expressionFragment: ConditionExpressionPart }>();
const search = ref(
  Array.isArray(props.expressionFragment.value)
    ? props.expressionFragment.value.join(",")
    : props.expressionFragment.value,
);
type DataOptions = { label: string; dataType: string };
const availableDataForSelectedOp = ref<DataOptions[]>([]);
const filteredAvailableDataForSelectedOp = ref<DataOptions[]>([]);
const getCurrentRoute = inject<() => Routes>("getCurrentRoute");
const setOperand = () => {};
const onEnter = (e: KeyboardEvent) => {
  if (e.key === "Enter") {
    props.expressionFragment.value = search.value;
  }
};

const loadPayloadValues = () => {
  const route = getCurrentRoute?.();
  if (route) {
    availableDataForSelectedOp.value =
      route.data?.payload.map((v) => {
        return {
          label: v.dataSource + "." + v.fieldName,
          dataType: v.dataType,
        };
      }) || [];
      filteredAvailableDataForSelectedOp.value =  availableDataForSelectedOp.value
  }
};
const loadOptions = ()=>{
    if (props.expressionFragment.source == 'static') {
        return
    }
    loadPayloadValues()
}
onMounted(() => {
   loadOptions()
})
</script>

<style></style>
