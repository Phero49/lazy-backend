<template>
  <block-unit-wrapper :scope="scope" :hide-toggle-btn="blockIndex > 0" :expandable="true">
    <template #header-extend-left>
      <div class="">
        <q-icon name="cable" class="q-mr-sm" color="grey" size="sm" />
       <span> step{{ operationIndex+1 }}</span>_ <span>
          {{ block.name || "step name" }}
        </span>
      </div>
    </template>
    <template #header-extend-right>
      <div class="q-ml-sm q-mr-md">
        <span>output data type </span>:<span
          class="q-ml-xs text-cyan-9 text-bold"
        >
          {{ method?.returnType }}</span
        >
      </div>
    </template>
    <div class="q-ml-md">
      <span>
        {{ block.name }}
        <q-tooltip class="text-body2 bg-black">
          {{ method?.description || "No description available" }}
        </q-tooltip>
        <q-popup-proxy>
          <q-banner>
            <FunctionSelector @select="onMethodSelected" />
          </q-banner>
        </q-popup-proxy>
      </span>
      <span class="text-h6">(</span>
      <span
        style="min-width: 100px"
        :key="index"
        v-for="(param, index) in props.block.params"
      >
        <span>
          <span class="text-grey-8 text-bold">{{ param.value.type }}: </span>
          <span>
            <FunctionParam
              :param="param"
              @update:value="
                (v) => {
                  param.value.data = v?.toString()!;
                }
              "
            />
            <q-tooltip> param name : {{ param.name }} </q-tooltip>
          </span>
        </span>
        <span class="q-mx-xs">,</span>
      </span>
      <span class="text-h6">)</span>
    </div>
  </block-unit-wrapper>
  <div></div>
</template>

<script lang="ts" setup>
import { DependencyBlock } from "~/store/appContext.js";
import FunctionParam from "./functionParam.vue";
import { onMounted, ref } from "vue";
import { GetDependencyExports } from "#/wailsjs/go/main/App.js";
import { io } from "#/wailsjs/go/models.js";
import { BlockScope } from "~/index.js";
import BlockUnitWrapper from "./BlockUnitWrapper.vue";
import FunctionSelector from "./functionSelector.vue";
const props = defineProps<{ block: DependencyBlock;
   scope: BlockScope,blockIndex:number,operationIndex:number }>();
const params = ref<io.Exports["params"]>({});
const method = ref<io.Exports>();
const onMethodSelected = (data: io.Exports, source: string) => {
  props.block.name = data.name;
  props.block.location = source;
  props.block.params = Object.keys(data.params).map((key) => ({
    name: key,
    inputType: "static",
    value: {
      data: "---",
      type: data.params[key],
    },
  }));

  method.value = data;
};
onMounted(async () => {
  GetDependencyExports(`${props.block.location}/manifest.json`)
    .then((data) => {
      const m = data.find((v) => v.name == props.block.name);
      if (m) {
        method.value = m;
      }
    })
    .catch((err) => {
      console.error("Error fetching dependency exports:", err);
    });
});
</script>

<style></style>
