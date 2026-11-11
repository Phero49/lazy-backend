<template>
  <div :class="tintClass(item.type)" >
    <q-menu context-menu>
        <q-list style="min-width: 100px">
            <q-item clickable v-close-popup>
                <q-item-section>delete operation</q-item-section>
            </q-item>
            <q-separator />
            <q-item clickable v-close-popup>
                <q-item-section>delete block</q-item-section>
            </q-item>
        </q-list>
    </q-menu>
    <DatabaseOperationsUi
      v-if="item.type == 'database'"
      :index="operationIndex"
      :blockIndex="blockIndex"
      :end-point-data="endPointData"
      :db-operation="item.databaseBlock"
      @update-payload="endPointInstance?.updateLimitPayload"
      :scope="scope.scope"
    >
    </DatabaseOperationsUi>
    <OperationHttpResponse
      :scope="scope.scope"
      :blockIndex="blockIndex"
      v-else-if="item.type == 'apiResponse'"
      :http-response="item.apiResponse"
    />
    <OperationDependencyUi
      :scope="scope.scope"
      :blockIndex="blockIndex"
      :operation-index="operationIndex"
      :block="item.dependenceBlock"
      v-else-if="item.type == 'dependency'"
    />
  </div>
</template>

<script lang="ts" setup>
import { OperationBlockData, Routes } from "~/store/appContext";
import { BlockScope } from "..";
import DatabaseOperationsUi from "./DatabaseOperationsUi.vue";
import OperationHttpResponse from "./operationHttpResponse.vue";
import OperationDependencyUi from "./operationDependencyUi.vue";
import { useInject, useProvide } from "~/composable/useProviderInjector";

const props = defineProps<{
  scope: { scope: BlockScope };
  item: OperationBlockData;
  blockIndex: number;
  endPointData: Routes;
  operationIndex: number;
}>();

props.scope.scope["operationIndex"] = props.operationIndex;

const endPointInstance = useInject("endPointContext");
useProvide("operationIndex", () => props.operationIndex);
const STEP_TINT: Record<string, string> = {
  database: "tint-data",
  dependency: "tint-fn",
  apiResponse: "tint-http",
};

function tintClass(type: string): string {
  return STEP_TINT[type] ?? "tint-data";
}

const contextMenu = [{
    label:'add'
}]
</script>

<style>
.tint-data,
.tint-fn,
.tint-http {
  background: rgba(255, 255, 255, 0.02); /* barely lifts it off the page */
  border-radius: 6px;
}
.tint-data {
  box-shadow: inset 3px 0 0 #4caf50;
}
.tint-fn {
  box-shadow: inset 3px 0 0 #7c9cf5;
}
.tint-http {
  box-shadow: inset 3px 0 0 #d9a441;
}
</style>
