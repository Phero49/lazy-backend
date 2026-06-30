<template>
  <div class="q-pa-md">
    <div class="row justify-between">
      <div class="text-overline text-grey-7 q-mb-md">Endpoints</div>
      <div class="text-overline q-mb-md">
        <q-btn
          round
          dense
          flat
          :icon="biPlusCircle"
          color="white"
          size="sm"
          unelevated
          @click="addNewPrefix = !addNewPrefix"
        />
      </div>
    </div>
    <div @keydown.Enter="savePrefix()" v-if="addNewPrefix">
      <q-input
        dense
        hint="eg /api/v1/auth/"
        v-model="prefix"
        type="text"
        label="create prefix "
      >
        <template #append>
          <q-card-actions>
            <q-btn
              flat
              color="white"
              @click="savePrefix()"
              dense
              :icon="biSave2Fill"
              size="sm"
            />
          </q-card-actions>
        </template>
      </q-input>
    </div>
    <q-list dense dark class="q-gutter-y-sm q-mt-md">
      <div
        v-for="(endpoint, i) in endpoints"
        :key="i"
        style="border-left: 1px solid grey"
      >
        <end-point-list-item
          :delete-end-pont="deleteEndPont"
          :add-route="addRoute"
          :key="i"
          :endpoint="endpoint"
          :level="[i]"
          :i="i"
        />
      </div>
    </q-list>
  </div>
</template>

<script setup lang="ts">
import { biPlusCircle, biSave2Fill } from "@quasar/extras/bootstrap-icons";
import { computed, onBeforeMount, ref, watch, watchEffect } from "vue";
import { EndPointList, useAppContextStore } from "../../store/appContext";
import { useQuasar } from "quasar";
import EndPointListItem from "../../components/endPointListItem.vue";
const prefix = ref("");
const addNewPrefix = ref(false);
const appContext = useAppContextStore();
const endpoints = computed(() => appContext.apiSchema);
const addSubPrefix = ref(false);
const method = ref("GET");
const $q = useQuasar();

const savePrefix = (index?: number) => {
  if (index !== undefined) {
    const alreadyExist = appContext.apiSchema[index].children.find(
      (v) => (v as EndPointList).prefix == prefix.value,
    );
    if (alreadyExist) {
      $q.notify({
        type: "negative",
        message: "prefix already exist",
        icon: "error",
        position: "top-left",
      });
      return;
    }

    appContext.apiSchema[index].children.push({
      prefix: prefix.value,
      children: [],
    });
    addSubPrefix.value = false;
  } else {
    const alreadyExist = appContext.apiSchema.find(
      (v) => (v as EndPointList).prefix == prefix.value,
    );
    if (alreadyExist) {
      $q.notify({
        type: "negative",
        message: "prefix already exist",
        icon: "error",
        position: "top-left",
      });
      return;
    }
    appContext.apiSchema.push({
      prefix: prefix.value,
      children: [],
    });
    addNewPrefix.value = false;
  }

  appContext.saveApi();
  prefix.value = "";
};
function deleteEndPont(index: number, index2?: number) {
  if (!index2) {
    appContext.apiSchema.splice(index, 1);
    appContext.saveApi();
  }
}

function addRoute(index: number, index2?: number) {
  if (index2 === undefined) {
    appContext.apiSchema[index].children = [
      { path: "", method: "GET" },
      ...appContext.apiSchema[index].children,
    ];
  } else {
    //  (appContext.apiSchema[index].children[index2] as EndPointList) = [ {path:'',method:'get'}, ...  appContext.apiSchema[index].children[index2]]
  }
}
onBeforeMount(() => {
  appContext.getApiSchema();
});

function addEndPoint() {}
</script>
