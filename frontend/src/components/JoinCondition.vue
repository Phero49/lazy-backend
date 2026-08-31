<template>
<span v-if="data" > 
    <span v-if="data.logical" class="q-px-md">
    {{ data.logical }}
    <q-menu>
      <q-list style="min-width: 100px">
        <q-item clickable @click="data.logical = 'AND'" >
          <q-item-section>AND</q-item-section>
        </q-item>
        <q-separator />
        <q-item @click="data.logical = 'OR'" clickable >
          <q-item-section>OR</q-item-section>
        </q-item>
      </q-list>
    </q-menu>
  </span>

  <span class="text-bold text-orange"
    >{{ data.leftTable
    }}<q-menu v-if="conditionIndex > 0">
      <q-list style="min-width: 100px">
        <q-item
          clickable
          v-close-popup
          :class="{ 'selected-item': data.leftTable === item }"
          @click="data.leftTable = item"
          v-for="(item, i) in tablesAvailableToJoinWith"
          :key="i"
        >
          <q-item-section>{{ item }}</q-item-section>
        </q-item>
      </q-list> </q-menu
  ></span>
  .
  <span class="text-bold q-px-sm text-info"
    >{{ data.leftColumn
    }}<q-menu>
      <q-list style="min-width: 100px">
        <q-item
          clickable
          v-close-popup
          :class="{ 'selected-item': data.leftColumn === item }"
          @click="data.leftColumn = item"
          v-for="(item, i) in availabColums(data.leftTable)"
          :key="i"
        >
          <q-item-section>{{ item }}</q-item-section>
        </q-item>
      </q-list> </q-menu
  ></span>
  <span class="q-px-sm"
    >{{ data.operator
    }}<q-menu>
      <q-list style="min-width: 100px">
        <q-item
          clickable
          v-close-popup
          :class="{ 'selected-item': data.operator === item }"
          @click="data.operator = item"
          v-for="(item, i) in operators"
          :key="i"
        >
          <q-item-section>{{ item }}</q-item-section>
        </q-item>
      </q-list> </q-menu
  ></span>

  <span>
    <q-btn
      class="absolute"
      style="top: -5px"
      color="grey"
      icon="more_horiz"
      flat
      padding="0px 3px"
    >
      <q-menu class="q-pa-md">
        <div class="q-mb-md">
          <div class="q-pa-sm text-center">Select condition data source</div>
          <div class="q-gutter-x-md items-center text-caption text-capitalize">
          <div
              class="inline-block"
              v-for="value in ['column', 'payload', 'static', 'computed']"
              :key="value"
            >
              <q-radio
                dense 
                v-model="data.right.source"
                :val="value"
                :label="value"
              />
            </div>
          </div>
        </div>
        <q-separator spaced />

        <!-- COLUMN -->
        <div v-if="data.right.source == 'column'">
          <div class="row no-wrap" style="width: 350px">
            <div class="col-5">
              <div class="text-caption text-weight-medium text-grey-7 q-px-sm q-pt-xs">
                Tables
              </div>
              <q-scroll-area style="height: 180px">
                <q-list dense>
                  <q-item
                    v-for="(item, i) in tablesAvailableToJoinWith"
                    :key="i"
                    clickable
                    v-ripple
                    dense
                    :class="{ 'selected-item': data.right.rightTable === item }"
                    @click="data.right.rightTable = item"
                  >
                    <q-item-section class="text-body2">{{ item }}</q-item-section>
                  </q-item>
                </q-list>
              </q-scroll-area>
            </div>

            <q-separator vertical />

            <div class="col-7">
              <div class="text-caption text-weight-medium text-grey-7 q-px-sm q-pt-xs">
                Column
              </div>
              <q-scroll-area style="height: 180px">
                <q-list dense v-if="data.right.rightTable">
                  <q-item
                    v-for="(item, i) in availabColums(data.right.rightTable)"
                    :key="i"
                    clickable
                    v-ripple
                    dense
                    v-close-popup
                    :class="{ 'selected-item': data.right.rightColumn === item }"
                    @click="data.right.rightColumn = item"
                  >
                    <q-item-section class="text-body2">{{ item }}</q-item-section>
                  </q-item>
                </q-list>
                <div v-else class="text-caption text-grey-5 text-italic q-pa-sm">
                  Select a table
                </div>
              </q-scroll-area>
            </div>
          </div>
        </div>

        <!-- PAYLOAD -->
        <div v-else-if="data.right.source == 'payload' && payloadData">
          <q-list dense style="min-width: 280px">
            <q-item
              v-for="(value, i) in payloadData"
              :key="i"
              clickable
              v-ripple
              dense
              v-close-popup
              :class="{ 'selected-item': data.right.value === value.fieldName }"
              @click="data.right.value = value.fieldName"
            >
              <q-item-section class="text-body2">
                {{ value.fieldName }}
              </q-item-section>
              <q-item-section side>
                <div class="row items-center q-gutter-x-xs">
                  <q-badge color="grey-4" text-color="grey-9" outline>
                    {{ value.dataSource }}
                  </q-badge>
                  <q-badge color="blue-4" text-color="blue-9" outline>
                    {{ value.dataType }}
                  </q-badge>
                </div>
              </q-item-section>
            </q-item>

            <q-item v-if="!payloadData.length">
              <q-item-section class="text-caption text-grey-5 text-italic">
                No payload fields available
              </q-item-section>
            </q-item>
          </q-list>
        </div>

        <!-- STATIC -->
        <div v-else-if="data.right.source == 'static'" style="min-width: 280px">
          <q-input
            v-model="data.right.value"
            dense
            outlined
            autofocus
            label="Static value"
            placeholder="Enter a value..."
          />
        </div>

        <!-- COMPUTED -->
        <div
          v-else-if="data.right.source == 'computed'"
          class="text-center text-grey-5 text-italic q-pa-md"
          style="min-width: 280px"
        >
          Coming soon
        </div>
      </q-menu>
    </q-btn>
    <span v-if="data.right.source == 'column'">
        <span class="text-bold text-orange">{{ data.right.rightTable }}</span
    >.<span>{{ data.right.rightColumn }}</span>
    </span><span v-else-if="data.right.source=='payload'"  > <span class="text-bold text-yellow">
      payload
    </span> .{{ data.right.value }}</span>
    <span v-else-if="data.right.source=='static'"  >"{{ data.right.value }}"</span>
  
  </span> 
</span>

</template>

<script setup lang="ts">
import { inject, ref,isReactive, computed } from "vue";
import { Join, Payload } from "../store/appContext";
import { useOptions } from "../store/optionsRegistry";
import { UseOperationInfo } from "..";
import { EndPointInstance } from "../composable/EndpointBulder";

const props = defineProps<{
  conditionIndex: number;
   operationIndex:number  
  tableIndex:number
  tablesAvailableToJoinWith: string[];
}>();

const operators = ref(["=", "!=", ">", "<", ">=", "<="]);
const options = useOptions();
const availabColums = (table: string) => options.getTableColumns(table);
const endPointInstance = inject<EndPointInstance>("endPointContext");
const data = computed(()=>endPointInstance?.getTableJoin({optIndex:props.operationIndex,tableIndex:props.tableIndex})?.conditions[props.conditionIndex])
const payloadData = computed(() => {
 
  return endPointInstance?.getPayload ()
});
</script>

<style scoped>
.selected-item {
  background: rgba(25, 118, 210, 0.18);
  color: #64b5f6;
  font-weight: 600;
}
</style>