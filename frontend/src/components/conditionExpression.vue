<template>
  <div class="row items-center no-wrap q-gutter-x-md">
    <!-- left operand -->
    <span class="cursor-pointer text-primary">
      {{ leftLabel || "match" }}
      <q-menu class="q-px-sm" v-model="leftOpen" @before-show="leftSearch = ''">
        <condition-data-picker
          v-if="expression.left"
          :expression-fragment="expression.left"
        />
      </q-menu>
    </span>

    <!-- operator -->
    <span class="cursor-pointer text-bold">
      {{ model.operator }}
      <q-menu v-model="operatorOpen" @before-show="operatorSearch = ''">
        <q-input
          v-model="operatorSearch"
          dense
          filled
          autofocus
          label="operator"
        />
        <q-separator />
        <q-list style="min-width: 100px">
          <q-item
            v-for="item in filteredOperators"
            :key="item"
            clickable
            v-close-popup
            :active="item === model.operator"
            @click="setOperator(item)"
          >
            <q-item-section>{{ item }}</q-item-section>
          </q-item>
        </q-list>
      </q-menu>
    </span>

    <!-- right operand -->
    <span class="cursor-pointer text-primary">
      {{ rightLabel || "value" }}
      <q-menu v-model="rightOpen" @before-show="rightSearch = ''">
        <q-menu
          class="q-px-sm"
          v-model="rightOpen"
          @before-show="leftSearch = ''"
        >
          <condition-data-picker
            v-if="expression.right"
            :expression-fragment="expression.right"
          />
        </q-menu>
      </q-menu>
    </span>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";
import type { Routes } from "~/store/appContext";
import ConditionDataPicker from "./conditionDataPicker.vue";
//import { getAvailableDataSources } from '~/utils/apiSchema';

// Everything below is derived from the real types, so there is nothing to keep in sync
//type DataSource = ReturnType<typeof getAvailableDataSources>[number];
type Operation = NonNullable<Routes["data"]>["operations"][number];
type ExpressionModel = NonNullable<Operation["condition"]>["data"][number];
type ConditionPart = NonNullable<ExpressionModel["left"]>;

const props = defineProps<{
  /** The expression to edit. It is mutated in place, nothing is emitted. */
  expression: ExpressionModel;
}>();

const operators = ["==", "!=", ">", "<", ">=", "<=", "in", "not in"];

// A computed (not a one-time copy) so the component follows the prop when an
// expression above it is removed and the list shifts.
const model = computed(() => props.expression);

const isListOperator = computed(() =>
  ["in", "not in"].includes(model.value.operator),
);

function partLabel(part: ConditionPart): string {
  if (part.source === "operationBlock") {
    return String(part.value) || `operation #${(part.blockIndex ?? 0) + 1}`;
  }
  if (part.source === "dependencies") return String(part.value) || "dependency";
  return String(part.value);
}

const leftLabel = computed(() =>
  model.value.left ? partLabel(model.value.left) : "",
);
const rightLabel = computed(() =>
  model.value.right ? partLabel(model.value.right) : "",
);

// Typed text becomes a static part. Numbers are stored as numbers so that
// >, <, >= and <= compare numerically (not for `in` / `not in` lists).
function staticPart(text: string): ConditionPart {
  const n = Number(text);
  const isNumber = text !== "" && !Number.isNaN(n) && !isListOperator.value;
  return { source: "static", value: isNumber ? n : text };
}

// ---- search / filtering -------------------------------------------------

const leftOpen = ref(false);
const rightOpen = ref(false);
const operatorOpen = ref(false);
const leftSearch = ref("");
const rightSearch = ref("");
const operatorSearch = ref("");

const filteredOperators = computed(() => {
  const needle = operatorSearch.value.trim();
  return needle ? operators.filter((o) => o.includes(needle)) : operators;
});

// ---- updates ------------------------------------------------------------

function setOperator(operator: string) {
  model.value.operator = operator;
}
</script>

<style></style>
