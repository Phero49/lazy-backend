<template>
  <span v-if="type == 'checkbox'" class="row justify-center">
    <q-checkbox
      :indeterminate-value="false"
      @update:model-value="updateValue"
      :false-value="false"
      :true-value="true"
      v-model="text"
    />
  </span>
  <span
    v-else
    style="cursor: pointer; width: 100%; display: inline-block"
    class="row item-center"
  >
    {{ text }}
    <q-popup-proxy
      ref="popupProxyRef"
      @hide="updateValue(text)"
      class="no-shadow"
    >
      <q-card class="no-shadow" flat bordered style="min-width: 150px">
        <div class="q-pa-sm" @keydown="onEnter">
          <q-input
            v-model="text as string | number"
            v-if="type == 'input'"
            :type="inputType"
            :label="label"
            :filled="filled"
            autofocus
            :outlined="outlined"
            :dense="dense"
            :clearable="clearable"
            @blur="updateValue(text)"
          />

          <div v-else>
            <q-input
              v-model="filter"
              @update:model-value="(v) => filterOptions(v as string)"
              type="text"
              filled
              label="filter"
              dense
            >
              <template #append>
                <q-icon name="search" size="sm" />
              </template>
            </q-input>

            <q-list v-if="options.length > 0">
              <q-item
                v-close-popup
                @click="
                  text = op;
                  updateValue(op);
                "
                clickable
                v-ripple
                v-for="op in filteredOptions"
                :class="[{ 'text-secondary': text == op }]"
              >
                <q-item-section>{{ op }}</q-item-section>
              </q-item>
            </q-list>
            <div
              v-else
              class="q-pa-md text-caption text-center text-capitalize"
            >
              <div>
                <q-icon name="warning" size="md" color="warning" />
              </div>
              no options available for <br />
              {{ label }}
            </div>
          </div>
          <!-- <q-select 
    v-else-if="type == 'select'" 
    v-model="text" 
    :options="options" 
    
    :label="label"
    :filled="filled"
    :outlined="outlined"
    :dense="dense"
    :clearable="clearable"
    emit-value
    map-options
    @update:model-value="updateValue"
  /> -->
        </div>
      </q-card>
    </q-popup-proxy>
  </span>
</template>

<script lang="ts" setup>
import { QPopupProxy, type QInputProps } from "quasar";
import { computed, ref, watch } from "vue";

interface Props {
  value: string | number | boolean;
  type?: "input" | "select" | "checkbox";
  options?: string[];
  label?: string;
  inputType?: QInputProps["type"];
  filled?: boolean;
  outlined?: boolean;
  dense?: boolean;
  clearable?: boolean;
}
const props = withDefaults(defineProps<Props>(), {
  type: "input",
  inputType: "text",
  options: () => [],
  label: "Label",
  filled: false,
  outlined: false,
  dense: false,
  clearable: false,
});

const computedOptions = computed(() => props.options);
const filteredOptions = ref(computedOptions.value);

const emit = defineEmits<{
  (e: "update:value", value: string | number | boolean): void;
  (e: "change", value: string | number | boolean): void;
}>();
const popupProxyRef = ref<QPopupProxy>();

// Create a local ref for the input value
const text = ref<string | number | boolean>(props.value);
const filter = ref("");
// Watch for external value changes
watch(
  () => props.value,
  (newValue) => {
    if (newValue !== text.value) {
      text.value = newValue || "";
    }
  },
);

// Update parent when local value changes
function updateValue(value: string | number | boolean) {
  if (props.inputType == "number") {
    value = +value;
  }
  console.log("sent", props.type);
  emit("update:value", value);
  if (props.type != "input") {
    emit("change", value);
  }
}

function filterOptions(text?: string) {
  filteredOptions.value = props.options.filter((v) => v.includes(text || ""));
}

function onEnter(e: KeyboardEvent) {
  if (e.key == "Enter") {
    popupProxyRef.value?.hide();
    e.preventDefault();
  }
}
</script>

<style scoped>
/* Add your custom styles here if needed */
</style>
