<template>
    <div
         
          class="row q-py-xs"
          style="cursor: pointer"
        >
          <q-menu
            @hide="
              () => {
                appContext.saveApi();
                updateRouteTab(value);
              }
            "
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
                  <!--change rout name-->
                    <q-form>
                      <q-input
                        dense
                        :rules="[(v: string) => !v.includes('/') || 'name cant contain a /']"
                        v-model="newRouteName"
                        label="edit path name"
                        
                        "
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
                            @click="RenamePath"
                          />
                        </template>
                      </q-input>
                    </q-form>
                  </q-card>
                </q-menu>
              </q-item>
              <q-item clickable>
                <q-menu
                  @hide="
                    () => {
                      appContext.saveApi();
                      updateRouteTab(value);
                    }
                  "
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
                        <!--change HTTP method-->
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
              <q-item
                clickable
                @click="deleteRoute(childIndex, value.path)"
                v-close-popup
              >
                <q-item-section>Remove route</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
          <div
            style="cursor: pointer"
            @click="
              appContext.createNewTab(
                `/api/${value.path}`,
                `${value.method} - ${value.path}`,
                'bi-code-slash',
                undefined,
                {
                  levels: [...level, childIndex],
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

                {{ value.name }}
              </div>
            </div>
          </div>
        </div>
      
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { type Routes, useAppContextStore } from '~/store/appContext';

const props = defineProps<{
  value: Routes;
  childIndex: number;
  level: number[];
}>();
const appContext = useAppContextStore();
const currentEndPointParent = computed(()=> appContext.getEndPoint(props.level));


const deleteRoute = (index: number, path: string) => {
  // appContext.getEndPoint(props.level)?.children.splice(index, 1);
  // console.log(props.endpoint.children[index])
currentEndPointParent.value?.children?.splice(index, 1);
  const tab = appContext.getTabByPath( path);

  if (tab) {
    appContext.closeTab(tab.id);
  }
  appContext.saveApi();
};


function updateRouteTab(value: Routes) {
  appContext.updateTab(value.path, {
    path: `/api/${value.path}`,
    label: `${value.method} - ${value.path}`,
  });
}

const newRouteName = ref(props.value.name);
function RenamePath() {
    if (currentEndPointParent.value) {
       const currentNode = currentEndPointParent.value.children[props.childIndex] as Routes
       if (currentNode) {
        currentNode.name = newRouteName.value.replace(/\s+/g, '-');
       }
       appContext.saveApi();
    }
}
</script>

<style>

</style>