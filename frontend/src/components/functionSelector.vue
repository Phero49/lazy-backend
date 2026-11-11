<template>
  <div style="min-width: 300px; min-height: 400px">
    <div style="position: sticky; z-index: 99; top: 0">
      <div class="bg-grey-10">
        <q-input
          @update:model-value="filter"
          class="full-width q-py-xs"
          v-model="search"
          type="text"
          filled
          dense
          label="search function"
        >
          <template #append><q-icon name="search" /></template>
        </q-input>
      </div>
    </div>
    <q-list>
        <q-item
          clickable
          v-close-popup
          @click="emits('select', item, item.location)"
          v-ripple
          v-for="(item, index) in filtered"
          :key="index"
        >
          <q-item-section
            >{{ item.name }}
            <q-item-label caption lines="2">
              {{ item.description }}
            </q-item-label>
          </q-item-section>
        </q-item>
    </q-list>
  </div>
</template>

<script lang="ts" setup>
import {
  GetDependencyManifest,
  GetProjectUsedDependencies,
} from "#/wailsjs/go/main/App";
import { io } from "#/wailsjs/go/models";
import { onMounted, ref } from "vue";
import { useAppContextStore } from "~/store/appContext";
import { DependencyMethods } from "..";
import { getDependencyAllExports } from "~/composable/EndpointBulder";
import { useInject } from "~/composable/useProviderInjector";
const appContext = useAppContextStore();
const exports = ref<FunctionListItem[]>([]);
const filtered = ref<FunctionListItem[]>([]);
interface  FunctionListItem extends io.Exports {
  location:string
}
const search = ref<string>("");
const emits = defineEmits<{ select: [io.Exports, string] }>();
const endPointContext = useInject("endPointContext");
const activeDependencies = ref();
const filter = () => {
 filtered.value = exports.value.filter((v)=>v.name.includes(search.value))
};
onMounted(async () => {
  const data = await endPointContext?.getDependencyAllExports();
const flatten = data?.map((v)=>{
 const d =  v.data.map((v2)=>({...v2,location:v.source}))
return d
}).flat()
  exports.value = flatten || [];
  filtered.value = flatten || [];
});
</script>

<style></style>
