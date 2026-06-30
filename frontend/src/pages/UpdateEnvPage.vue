<template>
  <div class="column no-wrap">
    <div class="q-pa-lg full-height column no-wrap">
      <!-- Header with Comment Style -->
      <div class="row items-center justify-between q-mb-md">
        <div class="text-mono text-comment text-weight-bold tracking-tight">
          // {{ envType.toUpperCase() }} Configuration
        </div>
        <q-btn
          flat
          round
          dense
          icon="bi-plus-circle"
          class="text-comment"
          @click="addNewKey"
        >
          <q-tooltip>Add New Key-Value Pair</q-tooltip>
        </q-btn>
      </div>

      <!-- VS Code Style Assignment -->
      <div class="row items-center q-mb-xs q-pl-md">
        <span class="text-mono text-comment q-mr-sm">export </span>
        <q-input
          v-model="configName"
          dark
          dense
          borderless
          placeholder="ConnectionName"
          input-class="text-mono text-key text-weight-bold"
          class="editor-input col-4"
        />
        <span class="text-mono text-punctuation q-ml-sm">=</span>
      </div>

      <!-- JSON Braces Style -->
      <div class="text-mono text-punctuation text-h6 q-mb-sm">{</div>

      <!-- Key-Value List -->
      <q-scroll-area class="q-pr-md" style="height: calc(100vh - 420px)">
        <div class="q-pl-md">
          <div
            v-for="(item, index) in items"
            :key="index"
            class="row items-start q-mb-sm no-wrap group"
          >
            <!-- Indentation Symbol or Space -->
            <div class="q-mr-sm text-grey-9 text-mono select-none">
              {{ (index + 1).toString().padStart(2, "0") }}
            </div>

            <!-- Key Editor -->
            <div class="col-4">
              <q-input
                v-model="item.key"
                dark
                dense
                borderless
                placeholder="KEY"
                input-class="text-mono text-key padding-none"
                class="editor-input"
              />
            </div>

            <!-- Separator -->
            <div class="text-mono text-punctuation q-px-sm">:</div>

            <!-- Value Editor -->
            <div class="col">
              <q-input
                v-model="item.value"
                dark
                dense
                borderless
                placeholder='"VALUE"'
                input-class="text-mono text-value padding-none"
                class="editor-input"
              />
            </div>

            <!-- Delete Action -->
            <q-btn
              flat
              round
              dense
              icon="bi-x"
              color="red-5"
              size="md"
              class="q-ml-sm opacity-0 group-hover-opacity-100 transition-base"
              @click="removeItem(index)"
            />
          </div>
        </div>
      </q-scroll-area>

      <div class="text-mono text-punctuation text-h6 q-mt-sm">}</div>

      <!-- Footer Actions -->
      <div class="row reverse q-mt-lg">
        <q-btn
          unelevated
          color="primary"
          label="Save to Environment"
          class="border-radius-md"
          no-caps
          :loading="saving"
          @click="saveChanges"
        />
        <q-btn
          flat
          label="Discard"
          color="grey-7"
          class="q-mr-sm"
          no-caps
          @click="discardChanges"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { GetProjectEnv, SaveProjectEnv } from "../../wailsjs/go/main/App";
import { useAppContextStore } from "../store/appContext";

const route = useRoute();
const router = useRouter();
const $q = useQuasar();
const store = useAppContextStore();

const envType = ref((route.params.type as string) || "General");
const configName = ref((route.params.type as string) + " Config");
const items = ref<{ key: string; value: string }[]>([]);
const saving = ref(false);

const addNewKey = () => {
  items.value.push({ key: "NEW_KEY", value: '""' });
};

const removeItem = (index: number) => {
  items.value.splice(index, 1);
};

const saveChanges = async () => {
  if (!store.currentProject) {
    $q.notify({ type: "negative", message: "No project selected" });
    return;
  }
  if (!configName.value) {
    $q.notify({ type: "negative", message: "Configuration name is required" });
    return;
  }

  saving.value = true;
  try {
    const envData = await GetProjectEnv(store.currentProject);

    // Construct the new configuration block
    const configBlock: Record<string, any> = {};
    items.value.forEach((item) => {
      let val: any = item.value;
      // Basic type casting
      if (val === "true") val = true;
      else if (val === "false") val = false;
      else if (!isNaN(Number(val)) && val !== "") val = Number(val);
      else if (val.startsWith('"') && val.endsWith('"'))
        val = val.substring(1, val.length - 1);

      configBlock[item.key] = val;
    });

    // Determine category mapping
    let category = envType.value.toLowerCase();
    if (category === "api key") category = "apiKeys";
    if (category === "database") category = "database";
    if (category === "ssh") category = "ssh";

    if (!envData[category]) envData[category] = [];

    // Add new item in format { "Name": { ...config } }
    const newItem = { [configName.value]: configBlock };
    envData[category].push(newItem);

    await SaveProjectEnv(store.currentProject, envData);

    $q.notify({
      type: "positive",
      message: `Added to ${category} configuration`,
      position: "bottom-right",
    });

    router.back();
  } catch (error) {
    $q.notify({ type: "negative", message: "Failed to save configuration" });
  } finally {
    saving.value = false;
  }
};

const discardChanges = () => {
  router.back();
};

onMounted(() => {
  const type = envType.value.toLowerCase();
  // Initialize with default values if specific env type
  if (type === "api keys" || type === "api key") {
    items.value.push({ key: "API_KEY", value: '"your_api_key_here"' });
  } else if (type === "database") {
    items.value.push({ key: "host", value: '"localhost"' });
    items.value.push({ key: "port", value: "3306" });
    items.value.push({ key: "user", value: '"root"' });
    items.value.push({ key: "password", value: '""' });
    items.value.push({ key: "database", value: '""' });
    items.value.push({ key: "timeout", value: "5000" });
    items.value.push({ key: "ssl", value: "false" });
  } else if (type === "ssh") {
    items.value.push({ key: "host", value: '""' });
    items.value.push({ key: "port", value: "22" });
    items.value.push({ key: "user", value: '"root"' });
    items.value.push({ key: "keyPath", value: '"~/.ssh/id_rsa"' });
    items.value.push({ key: "passphrase", value: '""' });
  } else {
    items.value.push({ key: "SERVER_PORT", value: "8080" });
    items.value.push({ key: "DEBUG_MODE", value: "true" });
  }
});
</script>

<style lang="scss" scoped>
.editor-bg {
  background-color: #1e1e1e !important;
}

.text-mono {
  font-family: "Fira Code", "Monaco", "Consolas", monospace !important;
  font-size: 13px;
}

/* VS Code Theming */
.text-comment {
  color: #6a9955;
}

.text-key {
  color: #9cdcfe;
}

.text-value {
  color: #ce9178;
}

.text-punctuation {
  color: #d4d4d4;
}

.padding-none {
  padding: 0 !important;
}

.editor-input {
  :deep(.q-field__control) {
    min-height: unset;
    height: 32px;
  }

  :deep(.q-field__native) {
    padding: 0 4px;

    &::placeholder {
      opacity: 0.3;
    }
  }
}

.border-radius-md {
  border-radius: 8px;
}

.opacity-0 {
  opacity: 0;
}

.group:hover .group-hover-opacity-100 {
  opacity: 1;
}

.transition-base {
  transition: all 0.2s ease;
}

.tracking-tight {
  letter-spacing: -0.02em;
}

/* Hide scrollbar for the line numbers */
.select-none {
  user-select: none;
}
</style>
