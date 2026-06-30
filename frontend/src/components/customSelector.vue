<template>
    <div>
   <q-input  v-model="filter"  
 @update:model-value="(v)=>filterOptions(v as string)"
type="text" filled label="filter"  dense>
  <template #append>
<q-icon name="search" size="sm" />
  </template>
</q-input>
    <q-list v-if="options.length > 0" >

      <q-item  @click=" updateValue(op)" clickable
       v-ripple v-for="op,i in filteredOptions"  >
       <q-item-section avatar>
        <q-checkbox @update:model-value="(v)=>{
          if (v) {
           updateValue(op) 
          }
        }" left-label v-model="selectedOptions[i]" />
       </q-item-section>
      
        <q-item-section v-close-popup>  {{ op }}</q-item-section>
      </q-item>
    </q-list>
    <div v-else class="q-pa-md text-caption  text-center text-capitalize
    ">
   <div>
     <q-icon name="warning" size="md" color="warning" />
   </div>
  no options available for <br/>
  {{ label ? label :'' }}
    </div>      
    </div>

</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';

const props =  defineProps<{options:string[],multiple?:boolean, label?:string,validate?:(input:string)=>boolean}>()
const emits   =  defineEmits<{'updated':[string]}>()
const filter = ref('')

const filteredOptions = ref<string[]>(computed(()=>props.options).value)
const selectedOptions = ref(Array.from({length:props.options.length},()=>false) )
const  updateValue = (v:string)=>{
   emits('updated',v)
}
function filterOptions(text?:string) {
filteredOptions.value = props.options.filter((v)=>v.includes(text||''))  
}


</script>

<style>

</style>