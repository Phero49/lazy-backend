<template>
  <div v-if="data">
    <div>
      <div class="text-subtitle1">
        <span style="cursor: pointer" class="text-bold text-primary"
          >{{ data.joinType }}
          <q-menu>
            <q-list dense style="min-width: 100px">
              <q-item
                @click="data.joinType = item"
                clickable
                v-for="(item, i) in joins"
                :key="i"
              >
                <q-item-section>{{ item }}</q-item-section>
              </q-item>
            </q-list>
          </q-menu></span
        >
        JOIN
        <span class="text-bold text-info">{{
          data.conditions[0].leftTable
        }}</span>
        ON   
        <join-condition
          :condition-index="i"
          :table-index="tableIndex"
          :operation-index="operationIndex"
          v-for="(value, i) in data.conditions"
          :data="value"
          :tables-available-to-join-with="tablesAvailableToJoinWith"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useOptions } from "../store/optionsRegistry";
import { Join, Payload, Routes } from "../store/appContext";
import JoinCondition from "./JoinCondition.vue";

const props = defineProps<{
  data?: Join;
  operationIndex:number
  tableIndex:number
  tablesAvailableToJoinWith: string[];

}>();

const joins = ref<Join["joinType"][]>(["INNER", "LEFT", "RIGHT", "FULL"]);
</script>
