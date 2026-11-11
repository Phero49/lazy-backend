<template>
 <q-item
                clickable
                v-ripple
                dense
                v-close-popup
                @click="$emit('select')"
              >
                <q-item-section class="text-body2">
                  {{ payload.fieldName }}
                </q-item-section>
                <q-item-section side>
                  <div class="row items-center q-gutter-x-xs">
                    <q-badge color="grey-4" text-color="grey-9" outline>
                      {{ payload.dataSource }}
                    </q-badge>
                    <q-badge color="blue-4" text-color="blue-9" outline>
                      {{ payload.dataType }}
                    </q-badge>
                  </div>
                </q-item-section>
              </q-item>
</template>

<script lang="ts" setup>
import { Payload } from '~/store/appContext';

defineProps<{
  payload: Payload;
}>();

defineEmits<{
  (e: 'select'): void;
}>();

</script>

<style>
</style>