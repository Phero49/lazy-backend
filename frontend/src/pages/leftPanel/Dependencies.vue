<template>
  <div style="overflow-x: hidden">
    <div class="row justify-between q-px-sm">
      <div class="text-overline text-grey-7">Dependencies</div>
      <div>
        <q-btn size="sm" icon="bi-three-dots" flat>
          <q-menu>
            <q-list style="min-width: 100px">
              <q-item clickable v-close-popup @click="loadLocalDependencies">
                <q-item-section>upload local dependencies </q-item-section>
              </q-item>
              <q-separator />
              <q-item clickable v-close-popup @click="showManifestDialog = true">
                <q-item-section>Generate dependency manifest</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>
    </div>
    <div class="col q-px-sm">
      <q-input
        v-model="search"
        dense
        type="text"
        filled
        label="search for dependencies"
      />
    </div>
    <q-separator />
    <q-tabs dense v-model="tab" no-caps>
      <q-tab name="market-place" icon="shop" label="market place" />

      <q-tab name="installed" icon="download" label="installed" />
    </q-tabs>
    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="market-place">
        <div class="text-h6">Mails</div>
        Lorem ipsum dolor sit amet consectetur adipisicing elit.
      </q-tab-panel>
      <q-tab-panel name="installed">
        <q-list>
          <q-item
            clickable
            @click="openDependency(item.name)"
            v-ripple
            v-for="(item, index) in appContext.installedDependencies"
            :key="index"
          >
            <q-item-section>
              <q-item-label>
                {{ item.name }} - v{{ item.version }}
              </q-item-label>

              <q-item-label caption :lines="1">
                {{ item.description }}
              </q-item-label>
              <q-item-label class="text-body2 text-grey text-bold">
                <div class="row justify-between">
                  <div>@pheromkuka</div>
                  <div>
                    <q-badge color="primary" label="use" />
                  </div>
                </div>
              </q-item-label>
            </q-item-section>
            <q-item-section side> </q-item-section>
          </q-item>
        </q-list>
      </q-tab-panel>
    </q-tab-panels>
  </div>

  <q-dialog v-model="loadingError" persistent>
    <q-card>
      <q-card-section>
        {{ errorMsg }}
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat label="Cancel" color="primary" v-close-popup />
        <q-btn flat label="Turn on Wifi" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <q-dialog v-model="showManifestDialog">
    <q-card style="min-width: 320px">
      <q-card-section>
        <div class="text-h6">Generate dependency manifest</div>
      </q-card-section>
      <q-card-section>
        Coming soon: this dialog will let you pick a generation script for your dependency manifest.
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat label="Close" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import {
  GetAllInstalledDependencies,
  LoadLocalDependencies,
} from "#/wailsjs/go/main/App";
import { useQuasar } from "quasar";
import { useAppContextStore } from "~/store/appContext";
const search = ref("");
const loadingError = ref(false);
const errorMsg = ref<string | unknown>("");
const showManifestDialog = ref(false);
const tab = ref("installed");
const $q = useQuasar();
const appContext = useAppContextStore();

function openDependency(name: string) {
  appContext.createNewTab(`/view-dependency/${name}`, " x module", "bi-layer");
}

async function loadLocalDependencies() {
  try {
    const dependency = await LoadLocalDependencies();
    appContext.installedDependencies.push(dependency);
    $q.notify({
      message: `${dependency.name} v${dependency.version} has been installed you can now 
  include in your project`,
      position: "top-left",
      color: "green",
    });
  } catch (error) {
    $q.notify({ message: "failed to load schema error" });
    loadingError.value = true;
    errorMsg.value = error;
  }
}
onMounted(() => {
  appContext.getInstalledDependencies();
});
</script>

<style></style>
