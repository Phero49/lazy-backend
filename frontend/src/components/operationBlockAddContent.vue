<template>
  <div v-if="!hideFab" class="q-mt-md q-mt-md q-ml-md">
    <q-fab color="grey-9" padding="xs" :icon="biPlus" direction="right">
      <q-fab-action color="purple" @click="()=>addDependency()" label="dependency" />
      <q-fab-action color="purple" @click="()=>addDbOp()" label="database block" />
      <q-fab-action color="purple" @click="()=>addApiResponse()" label="response" />
    </q-fab>
  </div> 

  <slot  :actions="{
    addDependency,addDbOp,addApiResponse
  }" ></slot>
</template>

<script lang="ts" setup>
import { Routes, useAppContextStore } from "~/store/appContext";
import { OperationBlock } from "../store/appContext";
import { biPlus } from "@quasar/extras/bootstrap-icons";
import { defaultCondition } from "~/utils/utills";

const props = withDefaults(
  defineProps<{
    blockData: OperationBlock["data"];
    blocks: OperationBlock[];
    topLevel?: boolean;
    hideFab?: boolean;
  }>(),
  { topLevel: true },
);
const appContext = useAppContextStore();
function save() {
  appContext.saveApi();
}

/** Direction of insertion relative to the given index. */
type Direction = "after" | "before";

/**
 * Insert `item` into `list`.
 * - If `index` is omitted → push at the end (original behaviour).
 * - Else insert before/after `index` based on `direction` (default "after").
 */
function insertBlock<T>(list: T[], item: T, index?: number, direction: Direction = "after") {
  if (index == null) {
    list.push(item);
    return;
  }
  // Clamp index so out-of-range values don't silently drop the item.
  const safeIndex = Math.max(0, Math.min(index, list.length));
  // For "after" we insert at index+1; for "before" at index.
  const insertAt = direction === "before" ? safeIndex : safeIndex + 1;
  list.splice(insertAt, 0, item);
}

/** Build a new top-level block wrapper (used when topLevel is true). */
function makeTopLevelBlock(data: OperationBlock["data"][number]): OperationBlock {
  return { condition: { ...defaultCondition }, data: [data] };
}

function addDependency(index?: number, direction: Direction = "after") {
  const data: OperationBlock["data"][number] = {
    type: "dependency",
    dependenceBlock: {
      location: '',
      name: "function_name",
      params: [],
    },
  };

  if (props.topLevel) {
    insertBlock(props.blocks, makeTopLevelBlock(data), index, direction);
  } else {
    insertBlock(props.blockData, data, index, direction);
  }
}

function addDbOp(index?: number, direction: Direction = "after") {
  const data: OperationBlock["data"][number] = {
    type: "database",
    databaseBlock: {
      name: "step " + props.blocks.length,
      query: { columns: [], tables: [], where: [] },
      databaseAction: "SELECT",
    },
  };

  if (props.topLevel) {
    insertBlock(props.blocks, makeTopLevelBlock(data), index, direction);
  } else {
    insertBlock(props.blockData, data, index, direction);
  }
  save();
}

function addApiResponse(index?: number, direction: Direction = "after") {
  const data: OperationBlock["data"][number] = {
    type: 'apiResponse',
    apiResponse: {
      code: 200,
      statusMessage: 'ok',
      userMessage: { dataType: 'null', message: null },
    },
  };

  if (props.topLevel) {
    insertBlock(props.blocks, makeTopLevelBlock(data), index, direction);
  } else {
    insertBlock(props.blockData, data, index, direction);
  }
  save();
}
</script>

<style></style>
