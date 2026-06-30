<template>
  <div style="border-left: 1px solid grey">
    <div
      style="cursor: pointer"
      class="row items-center"
      @click="
        endpoint.collapsed
          ? (endpoint.collapsed = false)
          : (endpoint.collapsed = true)
      "
    >
      <!-- Context menu for the folder/prefix node -->
      <q-menu :offset="[10, 0]" context-menu class="no-shadow">
        <q-list bordered style="min-width: 150px" separator>
          <!-- Add Prefix (existing) -->
          <q-item clickable>
            <q-menu v-model="addSubPrefix" class="no-shadow" anchor="top right">
              <q-card bordered flat style="min-width: 200px" class="q-pa-md">
                <div>
                  <q-form @submit="savePrefix()" class="q-gutter-md">
                    <q-input
                      dense
                      hint="eg /api/v1/auth/"
                      v-model="prefix"
                      :rules="[
                        (val: string) =>
                          (val || '').length > 0 || 'Prefix cannot be empty',
                      ]"
                      type="text"
                      label="create prefix"
                    >
                      <template #append>
                        <q-card-actions>
                          <q-btn
                            flat
                            color="white"
                            type="submit"
                            @click.stop="savePrefix()"
                            v-close-popup
                            dense
                            :icon="biSave2Fill"
                            size="sm"
                          />
                        </q-card-actions>
                      </template>
                    </q-input>
                  </q-form>
                </div>
              </q-card>
            </q-menu>
            <q-item-section class="text-capitalize">add Prefix</q-item-section>
          </q-item>

          <!-- Add Route with its own context menu -->
          <q-item
            clickable
            @click.stop="
              () => {
                addRouteMenu = true;
                console.log('clicked add route', addRouteMenu);
              }
            "
          >
            <q-item-section class="text-capitalize">add route </q-item-section>

            <q-menu :offset="[20, 0]" class="no-shadow" anchor="top right">
              <q-card bordered flat style="min-width: 250px" class="q-pa-md">
                <q-form @submit="saveRoute">
                  <q-input
                    dense
                    hint="eg /users/profile"
                    v-model="newRoutePath"
                    type="text"
                    label="route path"
                    class="q-mb-sm"
                  />
                  <q-select
                    dense
                    v-model="newRouteMethod"
                    :options="['GET', 'POST', 'PUT', 'DELETE', 'PATCH']"
                    label="HTTP method"
                    behavior="menu"
                  />
                  <q-card-actions align="center" class="q-mt-sm">
                    <q-btn
                      color="black"
                      @click.stop="saveRoute"
                      v-close-popup
                      dense
                      label="Save"
                    />
                  </q-card-actions>
                </q-form>
              </q-card>
            </q-menu>
          </q-item>

          <q-item clickable v-close-popup>
            <q-item-section>Middlewares</q-item-section>
          </q-item>

          <q-item clickable @click="deleteEndPont()" v-close-popup>
            <q-item-section>Remove route</q-item-section>
          </q-item>

          <!-- Rename with its own context menu -->
          <q-item clickable>
            <q-menu
              v-model="renameMenu"
              @show="editingName = endpoint.prefix"
              class="no-shadow"
              anchor="top right"
            >
              <q-card bordered flat style="min-width: 200px" class="q-pa-md">
                <div>
                  <q-form @submit="saveRename" class="q-gutter-md">
                    <q-input
                      dense
                      v-model="editingName"
                      type="text"
                      label="new prefix name"
                      autofocus
                    >
                      <template #append>
                        <q-btn
                          flat
                          color="white"
                          @click.stop="saveRename"
                          v-close-popup
                          dense
                          type="submit"
                          :icon="biSave2Fill"
                          label="Save"
                          size="sm"
                        />
                      </template>
                    </q-input>
                  </q-form>
                </div>
              </q-card>
            </q-menu>
            <q-item-section>rename route</q-item-section>
          </q-item>
        </q-list>
      </q-menu>

      <div class="row items-center q-pb-sm">
        <div class="q-mr-sm text-grey-3">
          <q-icon :name="endpoint.collapsed ? biChevronDown : biChevronRight" />
          <q-icon :name="biFolder2Open" size="20px" color="primary" />
        </div>
        <div class="text-subtitle2">
          {{ endpoint.prefix }}
        </div>
      </div>
    </div>

    <!-- Recursive children -->
    <div v-if="!endpoint.collapsed" class="q-ml-md">
      <template
        v-for="(value, childIndex) in endpoint.children"
        :key="childIndex"
      >
        <!-- If it has 'children', it's a nested folder - render recursively -->
        <EndpointTreeNode
          v-if="value && typeof value === 'object' && 'children' in value"
          :endpoint="value"
          :i="childIndex"
          :level="[i, childIndex]"
        />

        <!-- If it has 'path', it's a route object - render with its own context menu -->
        <div
          v-else-if="value && typeof value === 'object' && 'path' in value"
          class="row q-py-xs"
          style="cursor: pointer"
        >
          <q-menu
            @hide="appContext.saveApi()"
            :offset="[10, 0]"
            context-menu
            class="no-shadow"
          >
            <q-list bordered style="min-width: 150px" separator>
              <q-item clickable>
                <q-item-section> Edit route path</q-item-section>
                <q-menu class="no-shadow" anchor="top right">
                  <q-card
                    bordered
                    flat
                    style="min-width: 200px"
                    class="q-pa-md"
                  >
                    <q-form>
                      <q-input
                        dense
                        @update:model-value="
                          (v) => {
                            value.path = (v as string).trim().replace(' ', '_');
                          }
                        "
                        v-model="value.path"
                        label="edit path name"
                      >
                        <template #append>
                          <q-btn
                            flat
                            class="bg-white text-black"
                            v-close-popup
                            no-caps
                            size="sm"
                            dense
                            label="Save"
                          />
                        </template>
                      </q-input>
                    </q-form>
                  </q-card>
                </q-menu>
              </q-item>
              <q-item clickable>
                <q-menu
                  @hide="appContext.saveApi()"
                  class="no-shadow"
                  anchor="top right"
                >
                  <q-card
                    bordered
                    flat
                    style="min-width: 200px"
                    class="q-pa-md"
                  >
                    <q-form>
                      <q-select
                        dense
                        v-model="value.method"
                        :options="['GET', 'POST', 'PUT', 'DELETE', 'PATCH']"
                        label="HTTP method"
                        behavior="menu"
                      >
                        <template #append>
                          <q-btn
                            flat
                            class="bg-white text-black"
                            v-close-popup
                            no-caps
                            size="sm"
                            dense
                            label="Save"
                          />
                        </template>
                      </q-select>
                    </q-form>
                  </q-card>
                </q-menu>
                <q-item-section
                  >Change method ({{ value.method }})</q-item-section
                >
              </q-item>
              <q-item clickable @click="deleteRoute(childIndex)" v-close-popup>
                <q-item-section>Remove route</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
          <div
            style="cursor: pointer"
            @click="
              appContext.createNewTab(
                `/api/${endpoint.prefix}/${value.path}`,
                `${value.method} ${value.path}`,
                'bi-code-slash',
                {
                  level: [...level, childIndex].join(','),
                },
              )
            "
            class="text-subtitle2 text-grey"
          >
            <div class="row">
              <div class="">
                <q-badge
                  class="q-mr-sm method-badge"
                  :class="value.method.toLowerCase()"
                  :label="value.method"
                />

                {{ value.path }}
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  biChevronDown,
  biChevronRight,
  biFolder2Open,
  biSave2Fill,
} from "@quasar/extras/bootstrap-icons";
import {
  EndPointList,
  HTTPMethods,
  Routes,
  useAppContextStore,
} from "../store/appContext";
import { ref } from "vue";
import { useQuasar } from "quasar";

defineOptions({
  name: "EndpointTreeNode",
});

const props = defineProps<{
  endpoint: EndPointList;
  i: number;
  level: number[];
}>();
const appContext = useAppContextStore();
const prefix = ref("");
const addSubPrefix = ref(false);
const addRouteMenu = ref(true);
const renameMenu = ref(false);
const newRoutePath = ref("");
const newRouteMethod = ref("GET");
const editingName = ref("");
const $q = useQuasar();
const editRout = ref("");
const savePrefix = (index?: number) => {
  const newEndPoint: EndPointList = {
    prefix: prefix.value,
    children: [],
    collapsed: true,
  };
  const l = index ? [...props.level, index] : props.level;
  const find = appContext
    .getEndPoint(l)
    ?.children.find(
      (v) =>
        (v as EndPointList).prefix.replace("/", "").toLowerCase() ==
        prefix.value.replace("/", "").toLowerCase(),
    );
  if (!find) {
    appContext.getEndPoint(props.level)?.children.push(newEndPoint);
    appContext.saveApi();
  }
  prefix.value = "";
};

const saveRoute = () => {
  const endpoint = appContext.getEndPoint(props.level);
  if (endpoint == undefined) {
    return;
  }

  if (newRoutePath.value.trim()) {
    const normalizePath = newRoutePath.value
      .trim()
      .replace(" ", "_")
      .toLowerCase();
    const exists = endpoint.children.map(
      (v) =>
        (v as Routes).path == normalizePath &&
        newRouteMethod.value == (v as Routes).method,
    );
    if (exists != undefined) {
      $q.notify({ message: "path already exists", type: "negative" });
      return;
    }
    endpoint.children.push({
      path: normalizePath,
      method: newRouteMethod.value as HTTPMethods,
    });
    newRoutePath.value = "";
    newRouteMethod.value = "GET";
  }
  appContext.saveApi();
};

const deleteEndPont = (index?: number) => {
  appContext.deleteEndPont(props.level);
};

const deleteRoute = (index: number) => {
  appContext.getEndPoint(props.level)?.children.splice(index, 1);
  //  props.endpoint.children?.splice(index, 1);
};

const saveRename = () => {
  if (editingName.value.trim()) {
    const endpoint = appContext.getEndPoint(props.level);
    if (endpoint == undefined) {
      return undefined;
    }
    endpoint.prefix = editingName.value.trim();
    props.endpoint.prefix = editingName.value.trim();
    editingName.value = "";
    appContext.saveApi();
  }
};
</script>

<style scoped>
.method-badge.get {
  background-color: #61affe;
  color: white;
}
.method-badge.post {
  background-color: #49cc90;
  color: white;
}
.method-badge.put {
  background-color: #fca130;
  color: white;
}
.method-badge.delete {
  background-color: #f93e3e;
  color: white;
}
.method-badge.patch {
  background-color: #50e3c2;
  color: white;
}
</style>
