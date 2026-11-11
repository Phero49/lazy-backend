<template>
  <div v-if="$route.params['endPoint'] && endPointData">
    <div>
      <q-expansion-item
        expand-separator
        label="Payload"
        header-class="items-center text-subtitle1"
      >
        <EndpointPayload
          :payload="endPointData.data?.payload || []"
          :method="endPointData.method"
        >
        </EndpointPayload>
      </q-expansion-item>
    </div>

    <div
      class="q-ma-md"
      v-for="(value, i) in endPointData.data?.operations || []"
    >
      <div class="row justify-end"></div>
      <OperationBlock
        :previous-block-has-if="
          endPointData.data?.operations[i - 1]?.condition != null
        "
        :isBlock="value.data.length > 0"
        :route="endPointData"
        :index="i"
      >
        <template #main="scope">
          <div v-for="(item, index) in value.data" :key="index">
            <q-menu context-menu class="bg-grey-10">
              <operation-block-add-content
                :top-level="!scope.scope.enableCondition"
                :hide-fab="true"
                :blocks="endPointData.data?.operations || []"
                :block-data="value.data"
              >
                <template #default="actions">
                  <q-list separator style="min-width: 250px">
                    <template
                      v-for="(menuItem, itemIndex) in contextMenu"
                      :key="itemIndex"
                    >
                      <q-item
                        clickable
                        v-if="
                          menuItem.label == 'delete block' &&
                          value.condition?.disableCondition === false
                        "
                      >
                        <q-item-section avatar>
                          <q-icon :name="menuItem.icon" color="grey" />
                        </q-item-section>
                        <q-item-section class="text-capitalize text-grey">{{
                          menuItem.label
                        }}</q-item-section>
                      </q-item>
                      <q-item
                        @click="
                          menuItem.action(
                            index,
                            i,
                            actions.actions[
                              menuItem.key as keyof typeof actions.actions
                            ],
                          )
                        "
                        clickable
                        v-if="itemIndex > 0"
                      >
                        <q-item-section avatar>
                          <q-icon :name="menuItem.icon" color="grey" />
                        </q-item-section>
                        <q-item-section class="text-capitalize text-grey">{{
                          menuItem.label
                        }}</q-item-section>
                      </q-item>
                    </template>
                  </q-list>
                </template>
              </operation-block-add-content>
            </q-menu>
            <EndPointContentWrapper
              :scope="scope"
              :item="item"
              :operation-index="index"
              :block-index="i"
              :end-point-data="endPointData"
            />

            <q-separator spaced />
          </div>
          <operation-block-add-content
            :top-level="false"
            v-if="value.condition?.disableCondition === false"
            :blocks="[]"
            :block-data="value.data"
          />
        </template>
      </OperationBlock>
    </div>
    <OperationBlockAddContent
      :block-data="[]"
      :blocks="endPointData.data?.operations || []"
    />
  </div>
  <div v-else>
    <div class="text-center text-subtitle1 q-pa-md">
      Select or create an endpoint to start building api
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import { Routes, useAppContextStore } from "../store/appContext";
import { useRoute } from "vue-router";
import EndpointPayload from "../components/EndpointPayload.vue";
import { useEndPointContext } from "../composable/EndpointBulder";
import OperationBlock from "~/components/operationBlock.vue";
import OperationBlockAddContent from "~/components/operationBlockAddContent.vue";
import { useProvide } from "~/composable/useProviderInjector";
import EndPointContentWrapper from "~/components/EndPointContentWrapper.vue";
const appContext = useAppContextStore();
const route = useRoute();
const level = computed(
  () => (appContext.tabState[route.path].levels as number[]) ?? [],
);

const endPointData = computed(() => {
  if (level.value.length == 0) {
    return;
  }
  const endpoint = appContext.getEndPoint(
    level.value.slice(0, level.value.length - 1),
  );
  const lastindex = level.value.at(-1);
  if (endpoint && lastindex !== undefined) {
    const route = endpoint.children[lastindex] as Routes;
    if (route.data == undefined) {
      route.data = {
        payload: [],
        operations: [],
      };
    }

    return route;
  }
});

//TODO: remove the prop architecture
const endPoint = useEndPointContext(endPointData.value);

useProvide("endPointContext", endPoint);

type Caller = (
  index?: number | undefined,
  direction?: "after" | "before",
) => void;

interface ContextMenuItem {
  label: string;
  icon: string;
  key: "addDbOp" | "addDependency" | "addApiResponse" | string;
  action: (operationIndex: number, blockIndex: number, caller?: Caller) => void;
}

const contextMenu: ContextMenuItem[] = [
  {
    label: "delete block",
    icon: "delete",
    key: "",
    action: (_, blockIndex) => {
      endPoint.deleteBlockByIndex(blockIndex);
    },
  },
  {
    key: "",
    icon: "delete",
    label: "delete operation",
    action: (operationIndex, blockIndex) => {
      endPoint.deleteOperationByIndex(blockIndex, operationIndex);
    },
  },
  {
    key: "addDependency",
    icon: "cable",
    label: "Add dependency",
    action: (operationIndex, blockIndex, caller) => caller?.(operationIndex),
  },
  {
    icon: "bi-table",
    key: "addDbOp",
    label: "Add database block",
    action: (operationIndex, blockIndex, caller) => caller?.(operationIndex),
  },
  {
    icon: "http",
    key: "addApiResponse",
    label: "Add http response",
    action: (operationIndex, blockIndex, caller) => caller?.(operationIndex),
  },
];
</script>
