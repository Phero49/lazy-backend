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
        JOIN <span class="text-bold text-info">{{ data.baseTable }}</span> ON
        <join-condition
          v-for="value in data.conditions"
          :base-table="data.baseTable"
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
import { Join } from "../store/appContext";
import JoinCondition from "./JoinCondition.vue";

const props = defineProps<{
  data?: Join;
  tablesAvailableToJoinWith: string[];
}>();
const joins = ref<Join["joinType"][]>(["INNER", "LEFT", "RIGHT", "FULL"]);
</script>
