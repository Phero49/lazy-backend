<template>
  <block-unit-wrapper :scope="scope" :hide-toggle-btn="blockIndex > 0">
<template #header-extend-left>
<div class=" items-center">
            <q-icon name="http" size="sm" class="q-mr-sm" color="grey-6" />

respond with </div>
</template>
<template #header-extend-right>
  <div>
      <span class="text-bold text-cyan-9 q-mx-ms">
      code : <span class="text-bold text-cyan-9">{{httpResponse.code}}</span>

  </span> 
  </div>
  <div class="q-ml-sm q-mr-md">
    <span>data type </span>:<span
      class="q-ml-xs text-cyan-9 text-bold"
    >
      {{ dataType }}</span
    >

  </div>

  </template>
  <div> 
   
  

        <div class="row">
        <div class="col-4">
          <q-select use-input
            v-model="httpResponse.code"
            :options="filteredStatusCodes"
            label="status code"
            filled
            @filter="filterStatusCodes"
            dense 
            :prefix="`${httpResponse.statusMessage} :`"
          >
            <template #option="scope">
              <q-item @click="()=>{
                httpResponse.statusMessage = scope.opt.label
                scope.toggleOption(scope.opt.value)}">
                <q-item-section>
                  {{ scope.opt.label }} : {{ scope.opt.value }}
                </q-item-section>
              </q-item>
            </template>
          </q-select>
        </div>
        <div class="q-pl-sm col-2">
          <q-select v-model="dataType"  :options="dataTypes" dense label="data type" filled />
        </div>
        <div class="col-6 q-pl-sm">
          <div>
            <template v-if="!returnStaticType && dataType !== 'null'">
  <q-select
            v-model="message"
            :options="getStatuesCodeList"
            label="message"
            filled
            dense
          >
            <template #option="scope">
              <q-item @click="scope.toggleOption(scope.label)">
                <q-item-section>
                  {{ scope.opt.label }} ==> {{ scope.opt.value }}
                </q-item-section>
              </q-item>
            </template>
          </q-select>
            </template>
            <template v-else>
                <q-input dense filled :disable="dataType == 'null'" v-model="message" type="text" label="Label" />
            </template>
          
          <div>
            <q-tooltip >
              you will have  to write the message sent
            </q-tooltip>
            <q-checkbox left-label v-model="returnStaticType" label="static repose" />
          </div>
          </div>
          
        </div>
      </div>
    </div>    </block-unit-wrapper>
    
  
     
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { dataTypes, getStatuesCodeList } from '~/utils/apiSchema';
import { ApiResponse, BlockScope } from '..';
import BlockUnitWrapper from './BlockUnitWrapper.vue';

const props =  defineProps<{
  blockIndex:number,
  httpResponse:ApiResponse,
  scope:BlockScope}>()
const statusCode = ref('')
const dataType = ref<string>('null')
const returnStaticType =  ref(false)
const message  = ref("null")

// Filtered status codes
const filteredStatusCodes = ref([...getStatuesCodeList])

function filterStatusCodes(val: string, update: (fn: () => void) => void) {
  update(() => {
    if (val === '') {
      filteredStatusCodes.value = [...getStatuesCodeList]
    } else {
      const needle = val.toLowerCase()
      filteredStatusCodes.value = getStatuesCodeList.filter(
        (opt: any) =>
          String(opt.label).toLowerCase().includes(needle) ||
          String(opt.value).toLowerCase().includes(needle)
      )
    }
  })
}

onMounted(() => {
  
})

</script>

<style>

</style>