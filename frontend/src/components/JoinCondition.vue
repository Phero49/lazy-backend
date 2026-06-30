<template>
  <span class="text-bold text-info"
    >{{ baseTable }}.{{ data.baseColumn
    }}<q-menu>
      <q-list style="min-width: 100px">
        <q-item
          clickable
          @click="data.baseColumn = item"
          v-for="(item, i) in availabColums(baseTable)"
          :key="i"
        >
          <q-item-section>{{ item }}</q-item-section>
        </q-item>
        <q-separator />
      </q-list> </q-menu
  ></span>
  <span class="q-px-sm"
    >{{ data.operator
    }}<q-menu>
      <q-list style="min-width: 100px">
        <q-item
          clickable
          @click="data.operator = item"
          v-for="(item, i) in operators"
          :key="i"
        >
          <q-item-section>{{ item }}</q-item-section>
        </q-item>
        <q-separator />
      </q-list> </q-menu
  ></span>
  <span class="text-bold text-orange"
    >{{ data.targetTable
    }}<q-menu>
      <q-list style="min-width: 100px">
        <q-item
          clickable
          @click="data.targetTable = item"
          v-for="(item, i) in tablesAvailableToJoinWith"
          :key="i"
        >
          <q-item-section>{{ item }}</q-item-section>
        </q-item>
        <q-separator />
      </q-list> </q-menu></span
  >.<span
    >{{ data.targetColumn
    }}<q-menu>
      <q-list style="min-width: 100px">
        <q-item
          clickable
          @click="data.targetColumn = item"
          v-for="(item, i) in availabColums(data.targetTable)"
          :key="i"
        >
          <q-item-section>{{ item }}</q-item-section>
        </q-item>
        <q-separator />
      </q-list> </q-menu
  ></span>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { Join } from "../store/appContext";
import { useOptions } from "../store/optionsRegistry";

defineProps<{
  data: Join["conditions"][0];
  baseTable: string;
  tablesAvailableToJoinWith: string[];
}>();
const operators = ref(["=", "!=", ">", "<", ">=", "<="]);
const options = useOptions();
const availabColums = (table: string) => options.getTableColumns(table);
</script>

<style scoped></style>
