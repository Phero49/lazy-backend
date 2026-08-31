<template>
  <q-card flat style="max-width: 500px; width: 500px" class="my-card q-px-md">
    <q-input v-model="alias" dense type="text" label="alias" />
    <q-markup-table separator="cell" class="q-mt-md">
      <thead>
        <tr>
          <th class="text-left">Attribute</th>
          <th class="text-right">values</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(attr, i) in column.attributes" :key="i">
          <td class="text-left">
            {{ attr.name }}
            <q-tooltip class="bg-black text-body2">
              {{ getFullAttribute(attr.name)?.description }}
            </q-tooltip>
          </td>
          <td>
            <span v-if="attr.isFunction"> ( </span>
            <span
              v-for="(param, index) in getFullAttribute(attr.name)?.params"
              :key="index"
            >
              <q-menu context-menu>
                <q-list style="min-width: 100px">
                  <q-item
                    clickable
                    @click="attr.values[index].source = 'static'"
                    v-close-popup="0"
                  >
                    <q-item-section>Static value</q-item-section>
                    <q-item-section side>
                      <q-icon
                        :name="
                          attr.values[index].source === 'static'
                            ? symRoundedCheckBox
                            : symRoundedSquare
                        "
                        color="white"
                      />
                    </q-item-section>
                  </q-item>
                  <q-separator />
                  <q-item
                    clickable
                    @click="attr.values[index].source = 'column'"
                    v-close-popup="0"
                  >
                    <q-item-section>Column value</q-item-section>
                    <q-item-section side>
                      <q-icon
                        :name="
                          attr.values[index].source === 'column'
                            ? symRoundedCheckBox
                            : symRoundedSquare
                        "
                        color="white"
                      />
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
              <!--param values-->
              <span class="q-pa-xs">
                <span v-if="attr.values[index].source === 'column'"
                  ><span
                    contenteditable
                    @focus="()=>{
                      editingOn = index
                    }"
                    @blur="
                      (e) => {
                        const t = e.target as HTMLSpanElement;
                        attr.values[index].table = t.textContent;
                      }
                    "
                    style="outline-color: #c8c8c8"
                    >{{
                      attr.values[index].table
                        ? attr.values[index].table
                        : "Table"
                    }}</span
                  >.</span
                >
                <span
                  contenteditable
                  @blur="
                    (e) => {
                      const t = e.target as HTMLSpanElement;
                      attr.values[index].value = t.textContent;
                    }
                  "
                  style="outline-color: #c8c8c8"
                >
                  {{
                    attr.values[index].value
                      ? attr.values[index].value
                      : param.name
                  }}
                </span>
              </span>
              <span v-if="!param.optional" class="text-red text-subtitle1"
                >*</span
              >,
              <q-tooltip class="bg-black text-body2" v-if="editingOn !== index" >
               <div >
                {{ param.description }}
               </div> 
               
              </q-tooltip>
              <q-menu hover   v-else>
             <div class="q-pa-md" style="min-width: 300px;" >
                <div class="text-center text-bold text-grey" >Select table</div>
                <q-separator spaced  />
                <q-list >
                  <q-item clickable @click="attr.values[index].table = table" v-ripple v-for="(table, tIndex) in data.dbOperation.query.tables.map((v)=>v.table)" :key="tIndex">
                  <q-item-section side avatar>
                    <q-avatar size="sm" color="white">
                                    <q-icon color="black" name="bi-table" />
      
                    </q-avatar>
                  </q-item-section>
                    <q-item-section>
                      <q-item-label lines="1">
{{ table }}
                      </q-item-label>
                      </q-item-section>
                      <q-item-section side>
               <q-icon
                        :name="
                         attr.values[index].table === table
                            ? symRoundedCheckBox
                            : symRoundedSquare
                        "
                        color="white"
                      />
                    </q-item-section>
                  </q-item>
                </q-list>
               </div>
               </q-menu>
            </span>
            <span v-if="attr.isFunction"> ) </span>
          </td>
        </tr>
      </tbody>
    </q-markup-table>

    <div class="text-center q-mt-md">
      <q-btn label="add attribute" no-caps icon="add">
        <q-popup-proxy>
          <q-banner>
            <q-card-section>
              <q-input
                v-model="searchQuery"
                placeholder="Search attributes..."
                dense
                outlined
                clearable
                :loading="isLoading"
              >
                <template v-slot:append>
                  <q-icon v-if="searchQuery !== ''" name="search" />
                  <q-spinner v-if="isLoading" color="primary" size="1.5em" />
                </template>
              </q-input>
            </q-card-section>

            <!-- Error Display -->
            <q-card-section v-if="error">
              <q-banner inline-actions class="text-white bg-negative">
                <template v-slot:avatar>
                  <q-icon name="error" color="white" />
                </template>
                {{ error }}
                <template v-slot:action>
                  <q-btn
                    flat
                    color="white"
                    label="Retry"
                    @click="fetchAttributes"
                  />
                </template>
              </q-banner>
            </q-card-section>

            <!-- Loading Spinner -->
            <q-card-section v-if="isLoading" class="flex flex-center q-pa-xl">
              <div class="text-center">
                <q-spinner color="primary" size="3em" :thickness="4" />
                <div class="q-mt-sm text-grey-6">Loading attributes...</div>
              </div>
            </q-card-section>

            <!-- Attributes List -->
            <q-card-section v-else>
              <div
                style="height: 270px; overflow-y: auto"
                class="native-scroll"
              >
                <q-list separator>
                  <!-- No Results -->
                  <q-item
                    v-if="filteredAttributes.length === 0"
                    class="text-center"
                  >
                    <q-item-section>
                      <q-icon name="search_off" size="2em" color="grey-6" />
                      <div class="text-grey-6 q-mt-sm">
                        {{
                          searchQuery
                            ? "No matching attributes found"
                            : "No attributes available"
                        }}
                      </div>
                    </q-item-section>
                  </q-item>

                  <!-- Attribute Items -->
                  <q-item
                    v-for="(item, index) in filteredAttributes"
                    :key="index"
                    clickable
                    v-close-popup="true"
                    v-ripple
                    @click="addAttribute(item)"
                  >
                    <q-tooltip
                      class="bg-grey-10"
                      hover
                      anchor="bottom end"
                      :offset="[18, 18]"
                      self="top start"
                      :max-width="'400px'"
                    >
                      <div class="q-pa-md" style="min-width: 300px">
                        <!-- Name - highlighted with left border -->
                        <div
                          style="
                            display: flex;
                            padding: 8px 12px;
                            background: rgba(100, 200, 255, 0.08);
                            border-left: 3px solid #64c8ff;
                            border-radius: 4px;
                            margin-bottom: 8px;
                          "
                        >
                          <strong style="color: #64c8ff; min-width: 80px"
                            >Name</strong
                          >
                          <span
                            style="
                              color: #ffffff;
                              font-weight: 500;
                              word-break: break-word;
                            "
                            >{{ item.name }}</span
                          >
                        </div>

                        <!-- Type - highlighted with left border -->
                        <div
                          style="
                            display: flex;
                            padding: 8px 12px;
                            background: rgba(255, 200, 100, 0.08);
                            border-left: 3px solid #ffc864;
                            border-radius: 4px;
                            margin-bottom: 8px;
                          "
                        >
                          <strong style="color: #ffc864; min-width: 80px"
                            >Type</strong
                          >
                          <span
                            style="color: #e0e0e0; word-break: break-word"
                            >{{ item.type || "N/A" }}</span
                          >
                        </div>

                        <!-- Description - highlighted with left border -->
                        <div
                          style="
                            display: flex;
                            padding: 8px 12px;
                            background: rgba(100, 255, 150, 0.06);
                            border-left: 3px solid #64ff96;
                            border-radius: 4px;
                            margin-bottom: 8px;
                          "
                        >
                          <strong style="color: #64ff96; min-width: 80px"
                            >Description</strong
                          >
                          <span
                            style="color: #d0d0d0; word-break: break-word"
                            >{{ item.description || "N/A" }}</span
                          >
                        </div>

                        <!-- Parameters Section -->
                        <div
                          v-if="item.params && item.params.length"
                          style="margin-top: 12px"
                        >
                          <div
                            style="
                              padding: 10px 12px;
                              background: rgba(200, 150, 255, 0.06);
                              border-left: 3px solid #c896ff;
                              border-radius: 4px;
                            "
                          >
                            <!-- Parameters header -->
                            <div
                              style="
                                color: #c896ff;
                                font-weight: 600;
                                font-size: 0.95em;
                                margin-bottom: 10px;
                                letter-spacing: 0.5px;
                              "
                            >
                              ⚙️ Parameters
                            </div>

                            <!-- Individual parameters with proper spacing -->
                            <div
                              v-for="(param, idx) in item.params"
                              :key="idx"
                              style="
                                padding: 8px 10px;
                                margin-bottom: 6px;
                                background: rgba(255, 255, 255, 0.03);
                                border-radius: 4px;
                                border: 1px solid rgba(255, 255, 255, 0.05);
                              "
                            >
                              <!-- Parameter name -->
                              <div
                                style="
                                  display: flex;
                                  align-items: center;
                                  gap: 8px;
                                  margin-bottom: 4px;
                                  flex-wrap: wrap;
                                "
                              >
                                <span style="color: #8cb8ff; font-weight: 600">
                                  • {{ param.name
                                  }}<span
                                    v-if="param.variadic"
                                    style="color: #8cb8ff; opacity: 0.7"
                                    >, ...</span
                                  >
                                </span>
                                <span
                                  v-if="param.variadic"
                                  style="
                                    color: #ffcc66;
                                    font-size: 0.75em;
                                    padding: 1px 8px;
                                    border-radius: 10px;
                                    background: rgba(255, 204, 102, 0.15);
                                  "
                                >
                                  1+ args
                                </span>
                                <span
                                  v-if="param.optional !== undefined"
                                  :style="{
                                    color: param.optional
                                      ? '#80c880'
                                      : '#ff8a80',
                                    fontSize: '0.75em',
                                    padding: '1px 8px',
                                    borderRadius: '10px',
                                    background: param.optional
                                      ? 'rgba(128, 200, 128, 0.15)'
                                      : 'rgba(255, 138, 128, 0.15)',
                                  }"
                                >
                                  {{ param.optional ? "Optional" : "Required" }}
                                </span>
                              </div>

                              <!-- Parameter details -->
                              <div
                                style="
                                  padding-left: 18px;
                                  color: #b0b0b0;
                                  font-size: 0.9em;
                                  display: flex;
                                  flex-wrap: wrap;
                                  gap: 6px 16px;
                                "
                              >
                                <span v-if="param.value" style="color: #a0d8a0">
                                  <span style="color: #80c880">Value:</span>
                                  {{ param.value
                                  }}{{ param.variadic ? "[]" : "" }}
                                </span>
                                <span
                                  v-if="param.description"
                                  style="color: #c8c8c8"
                                >
                                  <span style="color: #a0a0a0">{{
                                    param.description
                                  }}</span>
                                </span>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </q-tooltip>
                    <q-item-section>
                      <q-item-label>{{ item.name }}</q-item-label>
                      <q-item-label caption
                        >Type: {{ item.type || "N/A" }}</q-item-label
                      >
                    </q-item-section>
                    <q-item-section side>
                      <q-badge
                        :color="getTypeColor(item.type)"
                        text-color="white"
                        outline
                      >
                        {{ item.type || "unknown" }}
                      </q-badge>
                    </q-item-section>
                  </q-item>
                </q-list>
              </div>

              <!-- Results Count -->
              <div class="q-mt-sm text-grey-6 text-caption text-right">
                Showing {{ filteredAttributes.length }} of
                {{ attributes?.length || 0 }} attributes
              </div>
            </q-card-section>
          </q-banner>
        </q-popup-proxy>
      </q-btn>
    </div>
  </q-card>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref, computed, reactive } from "vue";
import { DatabaseAttribute, OperationInfo } from "..";
import { importDatabaseColumnAttributes } from "~/client/database";
import { DatabaseQuery } from "../store/appContext";
import {
  symRoundedCheckBox,
  symRoundedSquare,
} from "@quasar/extras/material-symbols-rounded";

const attributes = ref<DatabaseAttribute[]>([]);

const error = ref<string | null>(null);
const isLoading = ref(false);
const searchQuery = ref("");
const props = defineProps<{
  data: OperationInfo;
  column: DatabaseQuery["columns"][0];
  colIndex: number;
}>();
const alias = ref("");
const includeColumnValue = reactive<Record<string, string>>({});
const editingOn = ref();
function addAttribute(attr: DatabaseAttribute) {
  if (props.column.attributes == undefined) {
    props.column.attributes = [];
  }

  props.column.attributes.push({
    name: attr.name,
    isFunction: attr.type == "function",
    values:
      attr.params?.map(() => ({
        source: "static",
        value: "",
        table: undefined,
      })) || [],
  });
}

const getFullAttribute = (name: string) => {
  return attributes.value.find((v) => v.name == name);
};
// Filtered attributes based on search query - client-side filtering only
const filteredAttributes = computed(() => {
  // If no attributes or empty array, return empty array
  if (!attributes.value || attributes.value.length === 0) {
    return [];
  }

  // If no search query, return all attributes
  if (!searchQuery.value || searchQuery.value.trim() === "") {
    return attributes.value;
  }

  const query = searchQuery.value.toLowerCase().trim();

  // Filter attributes client-side
  return attributes.value.filter((item) => {
    // Search in name
    if (item.name?.toLowerCase().includes(query)) return true;

    // Search in description
    if (item.description?.toLowerCase().includes(query)) return true;

    // Search in type
    if (item.type?.toLowerCase().includes(query)) return true;

    // Search in params
    if (item.params && item.params.length > 0) {
      return item.params.some(
        (param) =>
          param.name?.toLowerCase().includes(query) ||
          (param.value && String(param.value).toLowerCase().includes(query)) ||
          (param.description &&
            param.description.toLowerCase().includes(query)),
      );
    }

    return false;
  });
});

// Helper function for badge colors based on type
const getTypeColor = (type: string): string => {
  if (!type) return "grey";
  const typeLower = type.toLowerCase();
  if (
    typeLower.includes("string") ||
    typeLower.includes("text") ||
    typeLower.includes("varchar")
  )
    return "blue";
  if (
    typeLower.includes("number") ||
    typeLower.includes("int") ||
    typeLower.includes("float") ||
    typeLower.includes("decimal")
  )
    return "green";
  if (typeLower.includes("boolean") || typeLower.includes("bool"))
    return "orange";
  if (
    typeLower.includes("date") ||
    typeLower.includes("time") ||
    typeLower.includes("timestamp") ||
    typeLower.includes("datetime")
  )
    return "purple";
  if (
    typeLower.includes("array") ||
    typeLower.includes("list") ||
    typeLower.includes("set")
  )
    return "red";
  if (
    typeLower.includes("object") ||
    typeLower.includes("json") ||
    typeLower.includes("map")
  )
    return "teal";
  if (
    typeLower.includes("id") ||
    typeLower.includes("key") ||
    typeLower.includes("uuid")
  )
    return "indigo";
  return "grey";
};

// Fetch attributes function - fetches ALL attributes once
const fetchAttributes = async () => {
  isLoading.value = true;
  error.value = null;

  try {
    const result = await importDatabaseColumnAttributes();
    attributes.value = Array.isArray(result) ? result : [];
  } catch (e) {
    const errorMessage =
      e instanceof Error ? e.message : "Failed to load attributes";
    error.value = errorMessage;
    attributes.value = [];
    console.error("Error fetching attributes:", e);
  } finally {
    isLoading.value = false;
  }
};

// Load attributes on component mount
onBeforeMount(async () => {
  await fetchAttributes();
});
</script>

<style scoped>
.my-card {
  margin: 0 auto;
}
.required-param {
  text-decoration: underline red;
}

.q-banner {
  border-radius: 4px;
  margin-bottom: 8px;
}

/* Tooltip styling */
:deep(.q-tooltip) {
  font-size: 0.85em;
  max-height: 400px;
  overflow-y: auto;
}

/* Scrollable tooltip content */
:deep(.q-tooltip)::-webkit-scrollbar {
  width: 4px;
}
:deep(.q-tooltip)::-webkit-scrollbar-track {
  background: #2a2a2a;
}
:deep(.q-tooltip)::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 2px;
}

/* Loading animation */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.q-card-section {
  animation: fadeIn 0.3s ease-out;
}

/* Responsive design */
@media (max-width: 600px) {
  .q-item-section {
    font-size: 0.9em;
  }

  :deep(.q-tooltip) {
    max-width: 250px !important;
  }
}
.native-scroll {
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.2) transparent;
}
.native-scroll::-webkit-scrollbar {
  width: 8px;
}
.native-scroll::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}
</style>
