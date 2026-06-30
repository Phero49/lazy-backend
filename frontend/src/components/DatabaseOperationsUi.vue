<template>
  <q-expansion-item expand-icon="bi-chevron-up">
    <template #header>
      <div class="row col items-center justify-between">
        <div class="">
          {{ dbOperation.name || "step name" }}
        </div>
        <q-space />

        <div>
          <q-card-section align="right" class="q-gutter-x-sm">
            <q-btn
              color="red-5"
              size="sm"
              dense
              icon="science"
              @click="testOperation"
            />
            <q-btn
              color="red-8"
              dense
              size="sm"
              icon="delete"
              @click="deleteOperation"
            />
          </q-card-section>
        </div>
      </div>
    </template>

    <div
      class="row items-center q-col-gutter-x-md bg-dark-page q-pa-sm rounded-borders"
      style="border: 1px solid rgba(255, 255, 255, 0.05)"
    >
      <div class="col-4">
        <q-input
          v-model="dbOperation.name"
          dense
          filled
          type="text"
          label="Action Name"
          placeholder="e.g., get_user"
          @blur="appContext.saveApi()"
        />
      </div>

      <div class="col-3">
        <q-select
          v-model="dbOperation.databaseAction"
          dense
          filled
          :emit-value="true"
          :map-options="true"
          :options="[
            { label: 'Read Data (SELECT)', value: 'SELECT' },
            { label: 'Update Data (UPDATE)', value: 'UPDATE' },
            { label: 'Delete Data (DELETE)', value: 'DELETE' },
            { label: 'Insert Data (INSERT)', value: 'INSERT' },
          ]"
          label="Database Action"
          @update:model-value="appContext.saveApi()"
        />
      </div>

      <div
        class="col row items-center justify-end q-gutter-x-sm"
        v-if="dbOperation.databaseAction == 'SELECT'"
      >
        <div>
          <q-input
            v-model="startAt"
            style="width: 90px"
            :disable="pagination"
            dense
            filled
            :type="pagination ? 'text' : 'number'"
            label="Skip Items"
          >
            <q-tooltip class="bg-black" v-if="pagination">
              {{ startAt }}
            </q-tooltip>
          </q-input>
        </div>

        <div>
          <q-input
            v-if="!pagination"
            v-model="endAt"
            style="width: 90px"
            dense
            filled
            label="Limit To"
          />
          <q-input
            v-else
            v-model="batch"
            style="width: 90px"
            dense
            filled
            type="number"
            label="Batch Size"
          />
        </div>

        <div>
          <q-checkbox
            @update:model-value="onPagination"
            left-label
            v-model="pagination"
            label="Paginate"
          >
            <q-tooltip>Enable page-by-page loading</q-tooltip>
          </q-checkbox>
        </div>
      </div>
    </div>

    <div
      class="bg-dark-page q-pa-md rounded-borders"
      style="border: 1px solid rgba(255, 255, 255, 0.05)"
    >
      <div class="q-mb-sm">
        <div class="text-caption text-grey-5 text-uppercase text-bold q-mb-xs">
          From Table(s)
        </div>
        <div class="row q-mb-md items-center">
          <div
            v-for="(t, i) in dbOperation.query.tables"
            :key="'table-' + i"
            class="row items-center"
          >
            <q-chip
              :icon-remove="biXCircleFill"
              removable
              dense
              clickable
              @click="removeTable(i)"
              :label="t.table"
            />

            <q-btn
              unelevated
              dense
              class="q-mx-sm"
              padding="0px 3px"
              color="blue-grey-10"
              :icon="mdiSetAll"
              v-if="dbOperation.query.tables[i + 1]"
            >
              <q-popup-proxy>
                <q-card flat class="my-card q-pa-sm" style="min-width: 330px">
                  <div
                    v-if="
                      dbOperation.query.tables[i + 1].join &&
                  
                    "
                  >
                    {{ dbOperation.query.tables[i + 1].join }}sdd
                    <build-join
                      :data="dbOperation.query.tables[i + 1].join"
                      :tableToJoin="dbOperation.query.tables[i + 1].table"
                      :tablesAvailableToJoinWith="
                        dbOperation.query.tables
                          .map((v) => v.table)
                          .slice(0, i + 1)
                      "
                    />
                    <div class="text-center q-py-sm q-pt-md">
                      <q-btn
                        label="add join condition"
                        color="grey-9"
                        size="sm"
                        @click="
                          addJoinCondition(
                            i + 1,
                            dbOperation.query.tables[i + 1].table,
                          )
                        "
                        no-caps
                        unelevated
                      />
                    </div>
                  </div>
                </q-card>
              </q-popup-proxy>
            </q-btn>
          </div>

          <q-btn
            dense
            push
            color="grey-10"
            icon="bi-plus"
            size="sm"
            class="q-px-sm"
            label="Add Table"
          >
            <q-menu>
              <q-list style="min-width: 120px">
                <custom-selector
                  @updated="addTable"
                  :options="options.getTables"
                />
              </q-list>
            </q-menu>
          </q-btn>
        </div>
      </div>

      <q-separator dark class="q-my-sm" style="opacity: 0.1" />

      <div class="">
        <div class="text-caption text-grey-5 text-uppercase text-bold q-mb-xs">
          Get Output Columns
        </div>
        <div class="row items-center q-gutter-sm">
          <div
            v-for="(value, i) in dbOperation.query.columns"
            :key="'col-' + i"
          >
            <q-chip
              flat
              :label="value.alias ? value.alias : value.column"
              dense
              removable
              clickable
              :icon-remove="biXCircleFill"
              @remove="operations.columns.splice(i, 1)"
            >
              <q-menu>
                <q-list style="min-width: 100px">
                  <div
                    class="text-center q-pt-sm text-capitalize text-subtitle1 text-grey"
                  >
                    add function
                  </div>
                  <q-separator spaced />
                  <q-item clickable v-close-popup>
                    <q-item-section>New tab</q-item-section>
                  </q-item>
                  <q-separator />
                  <q-item clickable v-close-popup>
                    <q-item-section>New incognito tab</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-chip>
          </div>

          <q-btn color="grey-8" size="sm" dense icon="bi-plus" class="q-px-sm">
            <q-menu flat class="no-shadow">
              <div class="text-center q-py-sm text-bold">Select Columns</div>
              <select-columns
                @updated="columnsUpdated"
                :tables="dbOperation.query.tables.map((v) => v.table)"
              />
            </q-menu>
            <q-tooltip>Select Columns</q-tooltip>
          </q-btn>
        </div>
      </div>

      <q-separator dark class="q-my-sm" style="opacity: 0.1" />

      <div>
        <div class="text-caption text-grey-5 text-uppercase text-bold q-mb-xs">
          Filters & Conditions
        </div>
        <div class="row items-center q-gutter-sm">
          <div
            v-for="(where, i) in dbOperation.query.where || []"
            :key="'where-' + i"
            class="row items-center"
          >
            <q-badge color="grey-8" class="q-mx-xs" v-if="where.operator">{{
              where.operator
            }}</q-badge>
            <span class="text-light-blue text-white dense">
              <span class="text-weight-light text-grey-4 q-mr-xs"
                >{{ where.column.table }}.</span
              >
              <span class="text-bold">{{ where.column.name }}</span>

              <span class="text-amber-4 text-bold q-mx-sm text-uppercase">{{
                where.comparison
              }}</span>

              <span v-if="where.value.column.table">
                <span class="text-weight-light text-grey-4"
                  >{{ where.value.column.table }}.</span
                >{{ where.value.column.name }}
              </span>
              <span v-else-if="where.value.payload" class="text-green-3">
                payload.{{ where.value.payload }}
              </span>
              <span v-else-if="where.value.value" class="text-italic">
                "{{ where.value.value }}"
              </span>
            </span>
          </div>

          <q-btn color="grey-8" size="sm" dense icon="edit" class="q-px-sm">
            <q-menu style="max-width: 80%" flat class="no-shadow">
              <div class="text-center q-py-sm text-bold">
                Build Filter Condition
              </div>
              <where-builder
                :condition="dbOperation.query.where"
                @close="addWhereCondition"
                v-if="endPointData.data?.payload"
                :tables="dbOperation.query.tables.map((v) => v.table)"
                :payload="endPointData.data?.payload"
              />
            </q-menu>
          </q-btn>
        </div>
      </div>
    </div>
  </q-expansion-item>
</template>

<script lang="ts" setup>
import { biPlus, biX, biXCircleFill } from "@quasar/extras/bootstrap-icons";
import {
  mdiSetLeftCenter,
  mdiSetCenter,
  mdiSetAll,
} from "@quasar/extras/mdi-v7";
import { ref } from "vue";
import { useOptions } from "../store/optionsRegistry";
import CustomSelector from "./customSelector.vue";
import SelectComuns from "./selectComuns.vue";
import SelectColumns from "./selectColumns.vue";
import WhereBuilder from "./WhereBuilder.vue";
import {
  DatabaseOperations,
  DatabaseQuery,
  Routes,
  useAppContextStore,
} from "../store/appContext";
import { Condition } from "..";
import BuildJoin from "./BuildJoin.vue";

const appContext = useAppContextStore();
const operationName = ref("");
type OperationType = "SELECT" | "UPDATE" | "DELETE" | "INSERT";
const operationType = ref<OperationType>("SELECT");
const pagination = ref(false);
const startAt = ref("0");
const endAt = ref("10");
const batch = ref(10);
const props = defineProps<{
  endPointData: Routes;
  dbOperation: DatabaseOperations;
}>();
const filterTables = ref("");
const testOperation = () => {};
const deleteOperation = () => {};
const options = useOptions();
const emits = defineEmits<{ updatePayload: [boolean] }>();

const operations = ref<DatabaseQuery>({
  columns: [],
  tables: [],
  where: [],
});
const addTable = (table: string) => {
  if (props.dbOperation.query.tables == undefined) {
    props.dbOperation.query.tables = [];
  }
  props.dbOperation.query.tables.push({
    table: table,
  });

  const length = props.dbOperation.query.tables.length;
  if (length > 1) {
    const currentEntry = props.dbOperation.query.tables[length - 1];
    const previousEntry = props.dbOperation.query.tables[length - 2];

    if (currentEntry.join == undefined) {
      currentEntry.join = {
        baseTable: previousEntry.table,
        conditions: [],
        joinType: "INNER",
      };
    }

    addJoinCondition(length - 1, table);
  }
  appContext.saveApi();
};

const addJoinCondition = (
  tableIndex: number,

  joinTable: string,
) => {
  props.dbOperation.query.tables[tableIndex]?.join?.conditions.push({
    baseColumn: "---",
    targetColumn: "---",
    operator: "=",
    targetTable: joinTable,
  });
  appContext.saveApi();
};

const addWhereCondition = (condition: Condition[]) => {
  props.dbOperation.query.where = condition;
  appContext.saveApi();
};
function removeTable(index: number) {
  props.dbOperation.query.tables.splice(index, 1);
  appContext.saveApi();
}

function columnsUpdated(v: string, b: boolean) {
  if (props.dbOperation.query.columns == undefined) {
    props.dbOperation.query.columns = [];
  }

  if (b) {
    props.dbOperation.query.columns.push({ column: v, alias: "" });
  } else {
    const index = props.dbOperation.query.columns.findIndex(
      (c) => c.column == v,
    );
    if (index > 0) {
      props.dbOperation.query.columns.splice(index, 1);
    }
  }
  appContext.saveApi();
}

function onPagination(v: boolean) {
  if (v) {
    startAt.value = "payload.query.start_at";
  } else {
    startAt.value = "10";
  }

  emits("updatePayload", v);
  pagination.value = v;
}
</script>

<style></style>
