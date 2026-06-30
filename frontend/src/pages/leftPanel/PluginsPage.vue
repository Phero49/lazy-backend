<template>
    <div class="q-pa-md">
        <div class="text-overline text-grey-7 q-mb-md">Plugins</div>
        <q-list dense dark class="q-gutter-y-xs">
            <q-item clickable v-for="i in 3" :key="i" class="border-radius-sm">
                <q-item-section avatar>
                    <q-icon name="bi-plugin" size="14px" color="primary" />
                </q-item-section>
                <q-item-section class="text-caption">Plugin Module {{ i }}</q-item-section>
            </q-item>
        </q-list>
    </div>
</template>

<script setup lang="ts">
</script>
