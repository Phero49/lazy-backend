<template>
   
  <div v-if="tables.length > 0">
    <q-markup-table bordered separator="vertical" class="no-shadow">
        <thead>
            <tr>
                <th class="text-center" v-for="item in tables" ><q-icon :name="biTable" color="green" class="q-mr-sm" /> {{ item }}</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="row in maxRows" >
                <td v-for="t in tables"  >
                    <div  v-if="columns[t] && columns[t][row-1]">
                            <q-checkbox left-label @update:model-value="(bool)=>{
                               updateValue( t+'.'+ columns[t][row-1].name,bool)
                            }" v-model="columns[t][row-1].selected" />{{ 
                   columns[t] ==  undefined ? '-'  : columns[t][row-1].name
                    }}
                    </div>
                </td>
                
            </tr>
        </tbody>
    </q-markup-table>
  </div>
  <div v-else class="text-subtitle1 text-center q-pa-md" >
    <q-icon name="warning" color="warning"  size="md" /> <br>
    Select tables first b4 you select columns</div>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { useOptions } from '../store/optionsRegistry';
import { biTable } from '@quasar/extras/bootstrap-icons';

const props = defineProps<{tables:string[]}>()
const options  = useOptions()
const columns = reactive<Record<string,{name:string,selected:boolean}[]>>({})
const maxRows = ref<number>(1)
const emits   =  defineEmits<{'updated':[string,boolean]}>()
const  updateValue = (v:string,remove:boolean)=>{
   emits('updated',v,remove)
}
onMounted(() => {
    new Set(props.tables).forEach((t)=>{
     const c  =  options.getTableColumns(t)
     maxRows.value = (c.length > maxRows.value) ? c.length : maxRows.value

      columns[t] = c.map((n)=>({
        name:n,
        selected:false
      }))
    })
})


</script>

<style>

</style>