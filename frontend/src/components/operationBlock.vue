<template>
  <div
    :class="{ 'block-boundary q-py-sm': condition?.disableCondition === false }"
    class="q-pl-sm"
  >
    <div
      v-if="condition && enableCondition"
      class="full-width row items-center"
    >
      <!-- if / else if / else. Only switchable when a previous block already has an `if` -->
      <span
        class="text-bold text-red q-mr-md"
        :class="{ 'cursor-pointer': previousBlockHasIf }"
      >
        {{ condition.type }}
        <q-menu v-if="previousBlockHasIf">
          <q-list style="min-width: 100px">
            <q-item
              v-for="item in availableConditions"
              :key="item"
              clickable
              v-close-popup
              :active="item === condition.type"
              @click="setType(item)"
            >
              <q-item-section>{{ item }}</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </span>

      <!-- `else` has no expression -->
      <template v-if="condition.type !== 'else'">
        <span
          v-for="(expression, i) in condition.data"
          :key="i"
          class="row items-center no-wrap"
        >
          <!-- join with the previous expression -->
          <span v-if="i > 0" class="q-mx-sm text-bold cursor-pointer">
            {{ expression.joinClause ?? "AND" }}
            <q-menu>
              <q-list style="min-width: 100px">
                <q-item
                  v-for="item in joinClauses"
                  :key="item"
                  clickable
                  v-close-popup
                  :active="item === (expression.joinClause ?? 'AND')"
                  @click="expression.joinClause = item"
                >
                  <q-item-section>{{ item }}</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </span>

          <conditional-expression :expression="expression" />

          <q-btn
            flat
            round
            dense
            size="xs"
            color="grey"
            icon="close"
            class="q-ml-xs"
            @click="removeExpression(i)"
          >
            <q-tooltip>Remove expression</q-tooltip>
          </q-btn>
        </span>

        <q-btn
          color="primary"
          round
          dense
          size="sm"
          icon="add"
          flat
          @click="addExpression()"
        />
      </template>
    </div>

    <div class="row justify-end">
      <q-toggle
        size="sm"
        @update:model-value="toggleCondition"
        v-model="enableCondition"
        left-label
        label="condition"
        color="green"
      />
    </div>
  
    <slot name="main"></slot> 
    
    <div class="q-mt-md bg-grey-10 q-pa-md">
  <operation-http-response />
    </div>

    <operation-block-add-content  :blocks="props.route.data?.operations||[]" />
  </div>
</template>

<script lang="ts" setup>
import conditionalExpression from "./conditionExpression.vue";
import { computed, provide, ref, watch } from "vue";
import { Routes } from "~/store/appContext";
import { getStatuesCodeList, statusCodes } from "~/utils/apiSchema.js";
import OperationHttpResponse from "./operationHttpResponse.vue";
import { biPlus } from "@quasar/extras/bootstrap-icons/index.js";
import OperationBlockAddContent from "./operationBlockAddContent.vue";

const props = defineProps<{
  index: number;
  route: Routes;
  previousBlockHasIf: boolean;
}>();

const availableConditions = ["if", "else if", "else"] as const;
const joinClauses = ["AND", "OR"] as const;
const statusCode = ref(null);
const returnStaticType = ref(true)
const operation = computed(() => props.route.data?.operations[props.index]);
const condition = computed(() => operation.value?.condition);

// Recomputed so operations added or removed above are picked up

// `disableCondition` is the inverse of "enabled", and nothing is enabled
// until a condition exists
const enableCondition = ref(
  !!condition.value && !condition.value.disableCondition,
);

const toggleCondition = (enabled: boolean) => {
  const op = operation.value;
  if (!op) return;

  if (enabled) {
    if (!op.condition) {
      op.condition = { data: [], type: "if", disableCondition: false };
    } else {
      op.condition.disableCondition = false; // re-enable an existing condition
    }
    // start with one empty expression so there is something to fill in
    if (op.condition.type !== "else" && op.condition.data.length === 0) {
      addExpression();
    }
  } else if (op.condition) {
    op.condition.disableCondition = true; // keep the data, just switch it off
  }
};

const setType = (type: (typeof availableConditions)[number]) => {
  if (!condition.value) return;
  condition.value.type = type;

  if (type === "else") {
    condition.value.data = [];
  } else if (condition.value.data.length === 0) {
    addExpression();
  }
};

// The first block of a chain can only be `if`
watch(
  () => props.previousBlockHasIf,
  (hasIf) => {
    if (!hasIf && condition.value && condition.value.type !== "if") {
      setType("if");
    }
  },
);

const addExpression = () => {
  if (!condition.value) return;
  condition.value.data.push({
    left: {
      source: "static",
      value: "value",
      valueType: "string",
    },
    right: {
      source: "static",
      valueType: "string",
      value: "value",
    },

    operator: "==",
    // every expression after the first joins the one before it
    joinClause: condition.value.data.length > 0 ? "AND" : undefined,
  });
};

const removeExpression = (i: number) => {
  if (!condition.value) return;
  condition.value.data.splice(i, 1);

  // the new first expression has nothing to join to
  if (i === 0 && condition.value.data.length > 0) {
    condition.value.data[0].joinClause = undefined;
  }
};

provide("getCurrentRoute", () => {
  return props.route;
});

</script>

<style scoped>
.block-boundary {
  border-left: 2px solid red;
  border-radius: 15px;
}
</style>
