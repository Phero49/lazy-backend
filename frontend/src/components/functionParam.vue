<template>
  <span>
    {{ param.value.data  == '' || param.value.data == null ? '____' : param.value.data }}
    <q-popup-proxy>

      <q-banner style="min-width: 400px">
        <div class="row items-center">
          <div class="col-5 q-pr-xs">
            <div>
              <q-select
                dense
                v-model="param.inputType"
                :options="options"
                label="data source"
                outlined
              />
            </div>
          </div>
          <div class="col-7">
            <div v-if="param.inputType == 'static'">
              <q-input
                dense
                v-model="param.value.data"
                outlined
                label="write value"
              />
            </div>
            <div v-else-if="param.inputType == 'payload'">
              <q-select
                v-model="param.value.data"
                @update:model-value="emit('update:value', $event)"
                dense
                outlined
                :options="endPointContext?.getPayload()"
                label="select value"
              >
                <template v-slot:no-option>
                  <q-item>
                    <q-item-section
                      class="text-caption text-grey-5 text-italic"
                    >
                      No payload fields available
                    </q-item-section>
                  </q-item>
                </template>
                <template v-slot:option="scope">
                  <payload-list-item
                    :payload="scope.opt"
                    @select="
                      () => {
                        const s = scope.opt as Payload;
                        param.value.data = `payload.${s.dataSource}.${s.fieldName}`;
                        param.value.type = s.dataType;
                      }
                    "
                  />
                </template>
              </q-select>
            </div>
            <div v-else >
              <q-select
                v-model="param.value.data"
                @update:model-value="emit('update:value', $event)"
                dense
                outlined
                :options="endPointContext?.getEndPointOperationB4Position(getCurrentBlockIndex?.()||-1)"
                label="select value"
              >
                <template v-slot:no-option>
                  <q-item>
                    <q-item-section
                      class="text-caption text-grey-5 text-italic"
                    >
                      No operations available
                    </q-item-section>
                  </q-item>
                </template>
                <template #option="scope" >
                  <q-item @click="()=>{
                    param.value.data = `operationBlock.${scope.opt.name}`;
                   // param.value.type = scope.opt.returnType;
                  }">
                    <q-item-section class="text-body2">
                      {{ scope.opt.name }}
                    </q-item-section>
                  </q-item>
                </template>
              </q-select>
            </div>
          </div>
        </div>
      </q-banner> </q-popup-proxy
  ></span>
</template>

<script lang="ts" setup>
import { inject, ref } from "vue";
import { AvailableDataSources } from "..";
import { useEndPointContext } from "~/composable/EndpointBulder";
import PayloadListItem from "./payloadListItem.vue";
import { DependencyBlock, Payload } from "~/store/appContext";

const props = defineProps<{
  param: DependencyBlock["params"][number];
}>();
const model = ref<string>("static");
const emit = defineEmits<{
  (e: "update:value", value: string | number | null): void;
}>();
const options = ref<AvailableDataSources[]>([
  "static",
  "payload",
  "operationBlock",
]);
const getCurrentBlockIndex = inject<() => number>("getCurrentBlockIndex");
const endPointContext =
  inject<ReturnType<typeof useEndPointContext>>("endPointContext");
</script>

<style></style>
