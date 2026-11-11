<template>
  
<div  style="min-width: 350px" class="q-py-sm">
  <!--selector to pick data source-->
  <q-select
   
    class="q-mb-sm"
    v-if="expressionFragment.source"
    dense
    filled
    @update:model-value="()=>{
 
    expressionFragment.value = 'value'
  
    }"
    v-model="expressionFragment.source"
    :options="['static', 'payload', 'operationBlock']"
    label="data source"
  />
  <!--handles static values-->
  <!--it includes data type selection and type-->
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
        :options="staticValueDataTypeOptions"
        label="type"
        filled
      />
    </div>
    <div class="col-auto q-pl-sm" style="cursor: pointer">
      <q-input
        square
        :readonly="expressionFragment.valueType == 'null'"
        @keydown="onEnter"
        v-model="search" class="full-width"
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
  <div v-else-if="expressionFragment.source == 'payload'" class="q-pb-md">
    <!--select data source -->
    <q-select
      v-model="expressionFragment.value"
      dense
      :options="endPointInstance?.getPayload() || []"
      label="filter from payload "
      filled
  >
<template #option="scope">
<payload-list-item :payload="scope.opt" @select="()=>{
  const s = scope.opt as Payload
  expressionFragment.value =  `payload.${s.dataSource}.${s.fieldName}`
  

}" />
</template>
</q-select>
  </div>
  <div v-else-if="expressionFragment.source == 'operationBlock'" >
    <q-select dense v-model="expressionFragment.value" :options="endPointInstance?.getEndPointOperationB4Position(currentIndex)" label="select process" filled >
      <template #option="scope" >
<q-item>
  {{ scope }}
</q-item>
      </template>
      <template #no-option>
        <div class="q-pa-md text-center text-grey text-subtitle1">
        no previous blocks found 
        </div>
      </template>
    </q-select>
  </div>


</div>
</template>

<script lang="ts" setup>
import { inject, onMounted, ref } from "vue";
import { ConditionExpressionPart, Payload, Routes } from "~/store/appContext";
import PayloadListItem from "./payloadListItem.vue";
import { EndPointInstance } from "~/composable/EndpointBulder";
import { useInject, useProvide } from "~/composable/useProviderInjector.js";
import { DependencyMethods } from "~/index.js";

const props = defineProps<{ expressionFragment: ConditionExpressionPart }>();
const search = ref(
  Array.isArray(props.expressionFragment.value)
    ? props.expressionFragment.value.join(",")
    : props.expressionFragment.value,
);
type DataOptions = { label: string; dataType: string };
const availableDataForSelectedOp = ref<DataOptions[]>([]);
const filteredAvailableDataForSelectedOp = ref<DataOptions[]>([]);
const getCurrentRoute = useInject("getCurrentRoute");
const endPointInstance = useInject("endPointContext");
const getCurrentIndex  =  useInject('getCurrentBlockIndex')
const currentIndex = getCurrentIndex?.() || -1
const staticValueDataTypeOptions = [
          'null',
          'string',
          'int',
          'list',
          'bool',
          'list (string)',
          'list (int)',
        ]
const setOperand = () => {};
const onEnter = (e: KeyboardEvent) => {
  if (e.key === "Enter") {
    props.expressionFragment.value = search.value;
    search.value = ''
  }
};
const methodList = ref<DependencyMethods[]>([])
const loadDependencies  = async ()=>{
    methodList.value =  await endPointInstance?.getDependencyAllExports() || []

 }

</script>

<style></style>
