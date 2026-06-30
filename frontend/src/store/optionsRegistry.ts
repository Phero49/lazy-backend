import {defineStore} from  'pinia'
import { computed } from 'vue'
import { useAppContextStore } from './appContext'
import { SupportedOperators } from '..'

export const useOptions = defineStore('options',()=>{
  const appContext =  useAppContextStore()
     const getTables = computed(()=>{
       
        return appContext.activeSchema?.schema.entities.map((v)=>v.name) || []
    })
  
    const getTableColumns = (table:string)=>{
        const  entity = appContext.activeSchema?.schema.entities.find((v)=>v.name == table)
    if (entity) {
       return entity.fields.map((v)=>v.name)
    }
    return []
    }


const supportedOperators: SupportedOperators[] = [
  "Equal to",
  "Not equal to",
  ">",
  "<",
  ">=",
  "<=",
  "Like",
  "Not Like",
  "In",
  "Not In",
  "starts with",
  "ends with",
];
    

    return {

        getTables,
        getTableColumns
    }
})

//TODO: fix dynaminc fetching of tables