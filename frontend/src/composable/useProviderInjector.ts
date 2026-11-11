import { inject ,provide,type Ref} from "vue"
import { Routes } from "~/store/appContext"
import { EndPointInstance } from "./EndpointBulder"

export type ProviderMap = {
   /**
    * get the whole route object for the current end point tab
    */
   getCurrentRoute: ()=>Routes
 ,
 /**
  * 
  * get the index the of an operation block this method is called from
  */
 getCurrentBlockIndex:()=>number

 /**
  * get endpoint instance for the current end point 
  * endpoint instance includes method to access data
  */
 endPointContext:EndPointInstance

 /**
  * get operation index in the block
  */
 operationIndex:()=>number
} 

export function useProvide<K extends keyof ProviderMap>(key: K, value: any) {
  provide(key, value)          // string key, fine
}

export function useInject<K extends keyof ProviderMap>(key: K): ProviderMap[K] | undefined {
  return inject(key)           // string key, fine
}