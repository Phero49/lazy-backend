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
        />
      </q-expansion-item>
    </div>
    <div></div>

    <div
      class="q-ma-md"
      v-for="(value, i) in endPointData.data?.operations || []"
    >
      <DatabaseOperationsUi
        :end-point-data="endPointData"
        :db-operation="value"
        @update-payload="updateLimitPayload"
      />
    </div>
    <div>
      <q-btn
        color="primary"
        icon="bi-database"
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
    </div>
  </div>
  <div v-else>
    <div class="text-center text-subtitle1 q-pa-md">
      Select or create an endpoint to start building api
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, toRaw, watch } from "vue";
import {
  httpMethod,
  type HTTPMethods,
  Payload,
  Routes,
  useAppContextStore,
} from "../store/appContext";
import TableCellEditor from "../components/TableCellEditor.vue";
import { biChevronDown, biDash, biPlus } from "@quasar/extras/bootstrap-icons";
import DatabaseOperationsUi from "../components/DatabaseOperationsUi.vue";
import { useRoute } from "vue-router";
import { date } from "quasar";
import EndpointPayload from "../components/EndpointPayload.vue";

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

function updateLimitPayload(v: boolean) {
  if (v) {
    endPointData.value?.data?.payload.push({
      dataType: "int",
      fieldName: "start_at",
      dataSource: "query",
      required: true,
    });
  } else {
    const index = endPointData.value?.data?.payload.findIndex(
      (v) => v.fieldName == "start_at",
    );
    if (index != undefined && index != -1) {
      endPointData.value?.data?.payload.splice(index, 1);
    }
  }
}

const testOperation = () => {
  null;
};
const deleteOperation = () => {
  null;
};
</script>

<style></style>
