<template>
  <div class="q-mt-md">
    <q-fab color="primary" padding="sm" :icon="biPlus" direction="right">
      <q-fab-action color="purple" @click="()=>null" label="dependency" />
      <q-fab-action color="purple" @click="()=>null" label="database block" />
      <q-fab-action color="purple" @click="()=>null" label="response" />
    </q-fab>
  </div>
</template>

<script lang="ts" setup>
import { Routes } from '~/store/appContext';
import { OperationBlock } from '../store/appContext';
import { biPlus } from '@quasar/extras/bootstrap-icons';

const props =  defineProps<{blocks:OperationBlock[]}>()

function addDependency() {
 
    // props.blocks.push({
    //     dependenceBlock:{
    //         name:
    //     }
    // })

}

function addDbOp (){
    props.blocks.push({
        children:[],
        type:'database',
        condition:null,
        databaseBlock:{
            name:'step '+props.blocks.length,
            query:{
                columns:[],
                tables:[],
                where:[]
            },
            databaseAction:'SELECT'
        }
    })
}

function addApiResponse() {
    
}
</script>

<style></style>
