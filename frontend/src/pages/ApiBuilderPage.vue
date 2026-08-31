<template>
  <div v-if="$route.params['endPoint'] && endPointData">
    <div>
      <q-expansion-item
        expand-separator
        expand-icon="bi-chevron-down"
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
    <div></div>

    <div
      class="q-ma-md"
      v-for="(value, i) in endPointData.data?.operations || []"
    >
      <DatabaseOperationsUi
        :index="i"
        :end-point-data="endPointData"
        :db-operation="value"
        @update-payload="endPoint.updateLimitPayload"
      >
        <template #condition>
          <div class="q-px-sm" style="cursor: pointer" @click.stop>
            <q-select
              @click.stop
              dense
              v-model="activeCondition"
              :options="['fff']"
              label="Select Condition"
              filled
            >
              <template #before-options>
                <div>
           
                  <div class="text-center q-py-md">
                    <q-btn
                      color="grey"
                      no-caps
                      class="text-capitalize"
                      unelevated
                      icon="add"
                      label="add new condition"
                      @click="() => null"
                    />
                  </div>
                </div>
              </template>
            </q-select>
          </div>
        </template>
      </DatabaseOperationsUi>
    </div>
    <div class="text-center">
      <q-btn
        color="grey"
        dense
        :icon="symRoundedStepInto"
        @click="
          endPointData.data?.operations.push({
            databaseAction: 'SELECT',
            name: '',
            query: {
              tables: [],
              columns: [],
              where: [],
            },
          });
          appContext.saveApi();
        "
      >
      </q-btn>
      <q-btn color="grey" dense :icon="symRoundedArrowSplit" @click="">
        <q-menu anchor="center middle">
          <q-list style="min-width: 100px">
            <div class="q-pa-md">add branching condition</div>
            <q-item clickable v-close-popup>
              <q-item-section>New tab</q-item-section>
            </q-item>
            <q-separator />
            <q-item clickable v-close-popup>
              <q-item-section>New incognito tab</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </div>
  </div>
  <div v-else>
    <div class="text-center text-subtitle1 q-pa-md">
      Select or create an endpoint to start building api
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, provide, ref, toRaw, watch } from "vue";
import {
  httpMethod,
  type HTTPMethods,
  Payload,
  Routes,
  useAppContextStore,
} from "../store/appContext";
import TableCellEditor from "../components/TableCellEditor.vue";
import { biChevronDown, biDash, biPlus } from "@quasar/extras/bootstrap-icons";
import {
  symRoundedStepInto,
  symRoundedArrowSplit,
} from "@quasar/extras/material-symbols-rounded";
import DatabaseOperationsUi from "../components/DatabaseOperationsUi.vue";
import { useRoute } from "vue-router";
import { date } from "quasar";
import EndpointPayload from "../components/EndpointPayload.vue";
import { useEndPointContext } from "../composable/EndpointBulder";

const route = useRoute();
const level = computed(() =>
  (route.query.level as string).split(",").map((v) => +v),
);
const appContext = useAppContextStore();
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
const activeCondition = ref<string | null>(null);

//TODO: remove the prop architecture
const endPoint = useEndPointContext(endPointData.value);
provide("endPointContext", endPoint);

const testOperation = () => {
  null;
};
const deleteOperation = () => {
  null;
};
</script>

<style></style>
