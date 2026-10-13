<template>
  <div class="row justify-center" v-if="!app.currentProject">
    <div style="max-width: 550px;width: 100%;" >
            <projects-page />

    </div>
  </div>
</template>

<script lang="ts" setup>
import { useAppContextStore } from '~/store/appContext.js';
import ProjectsPage from './leftPanel/ProjectsPage.vue';
const app  = useAppContextStore()

</script>

<style>

</style>