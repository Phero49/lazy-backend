<template>
    <div> 
      <div class="text-subtitle2 text-capitalize">
          Respond with: <q-icon size="sm"  name="bi-reply" />
      </div>
      <q-separator spaced   />

        <div class="row">
        <div class="col-4">
          <q-select
            v-model="statusCode"
            :options="getStatuesCodeList"
            label="status code"
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
      </div></div>
     
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { dataTypes, getStatuesCodeList } from '~/utils/apiSchema';

const statusCode = ref('')
const dataType = ref<string>('null')
const returnStaticType =  ref(false)
const message  = ref("null")
</script>

<style>

</style>