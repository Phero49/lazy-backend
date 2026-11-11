<template>
  <div v-if="!hideFab" class="q-mt-md q-mt-md q-ml-md">
    <q-fab color="grey-9" padding="xs" :icon="biPlus" direction="right">
      <q-fab-action color="purple" @click="addDependency" label="dependency" />
      <q-fab-action color="purple" @click="addDbOp" label="database block" />
      <q-fab-action color="purple" @click="addApiResponse" label="response" />
    </q-fab>
  </div> 

  <slot  :scope="{
    addDependency,addDbOp,addApiResponse
  }" ></slot>
</template>

<script lang="ts" setup>
import { Routes, useAppContextStore } from "~/store/appContext";
import { OperationBlock } from "../store/appContext";
import { biPlus } from "@quasar/extras/bootstrap-icons";
import { defaultCondition } from "~/utils/utills";

const props = withDefaults(
  defineProps<{
    blockData: OperationBlock["data"];
    blocks: OperationBlock[];
    topLevel?: boolean;
    hideFab?:boolean
  }>(),
  { topLevel: true },
);    
const appContext  =  useAppContextStore()
function save() {
     appContext.saveApi()

}
function addDependency() {
  if (props.topLevel) {
    props.blocks.push({
      condition: {...defaultCondition},
      data: [
        {
          type: "dependency",
          
          dependenceBlock: {
            location:'',
            name: "function_name",
            params: [],
          },
        },
      ],
    });
  } else {
    props.blockData.push({
      type: "dependency",
      dependenceBlock: {
        
        location:'',
        name: "function_name",
        params: [],
      },
    });
  }

}


function addDbOp() {
  if (props.topLevel) {
    props.blocks.push({
      condition: {...defaultCondition},
      data: [
        {
          type: "database",
          databaseBlock: {
            name: "step " + props.blocks.length,
            query: {
              columns: [],
              tables: [],
              where: [],
            },
            databaseAction: "SELECT",
          },
        },
      ],
    });
  } else {
    props.blockData.push({
      type: "database",
      databaseBlock: {
        name: "step " + props.blocks.length,
        query: {
          columns: [],
          tables: [],
          where: [],
        },
        databaseAction: "SELECT",
      },
    });
  }
 save()
}

function addApiResponse() {
    if (props.topLevel) {
          props.blocks.push({
      condition: {...defaultCondition},
      data: [{
        type:'apiResponse',
        apiResponse:{
            code:200,
            statusMessage:'ok',
            userMessage:{
                dataType:'null',
                message:null
            }
        }
      }]})
    }else{
        props.blockData.push({
        type:'apiResponse',
        apiResponse:{
            code:200,
            statusMessage:'ok',
            userMessage:{
                dataType:'null',
                message:null
            }
        }
      })
    }
save()
}
</script>

<style></style>
