<template>
  <div style="width: 100%" class="q-pa-md column full-height">
    <div class="row items-center justify-between q-mb-md">
      <div class="text-overline text-grey-7">Projects</div>
      <q-btn
        flat
        round
        dense
        icon="bi-plus-lg"
        size="xs"
        color="primary"
        @click="addNewProject = true"
      >
        <q-tooltip>Create New Project</q-tooltip>
      </q-btn>
    </div>

    <!-- Skeleton Loader -->
    <template v-if="projects === undefined">
      <q-list dense dark class="q-gutter-y-sm">
        <q-item v-for="i in 5" :key="i" class="q-px-none">
          <q-item-section avatar class="min-width-auto q-pr-sm">
            <q-skeleton type="QBtn" size="14px" dark />
          </q-item-section>
          <q-item-section>
            <q-skeleton type="text" width="80%" dark />
          </q-item-section>
        </q-item>
      </q-list>
    </template>

    <!-- Empty State -->
    <template v-else-if="projects.length === 0">
      <div class="column items-center justify-center q-pa-xl text-grey-8">
        <q-icon name="bi-folder2-open" size="48px" class="q-mb-md" />
        <div class="text-caption">No projects found</div>
        <q-btn
          flat
          no-caps
          label="Create one core"
          color="primary"
          class="q-mt-sm"
          @click="addNewProject = true"
        />
      </div>
    </template>
    <!-- Projects List -->
    <q-scroll-area v-else style="height: 80vh; width: 100%">
      <q-list dense dark class="q-gutter-y-xs">
        <q-item
          clickable
          v-for="project in projects"
          :key="project.name"
          class="project-item"
          :active="store.currentProject === project.name"
          active-class="active-project"
          @click="selectProject(project)"
        >
          <q-item-section avatar class="min-width-auto q-pr-sm">
            <q-icon
              :name="
                store.currentProject === project.name
                  ? 'bi-folder-fill'
                  : 'bi-folder'
              "
              size="14px"
              :color="
                store.currentProject === project.name ? 'primary' : 'grey-7'
              "
            />
          </q-item-section>
          <q-item-section class="text-caption">{{
            project.name
          }}</q-item-section>

          <q-item-section side v-if="store.currentProject === project.name">
            <div class="active-dot" />
          </q-item-section>
          <q-item-section side v-else>
            <q-btn
              color="grey-7"
              icon="bi-trash"
              size="xs"
              flat
              round
              dense
              @click="deleteProject(project)"
            >
              <q-tooltip>Delete Project</q-tooltip></q-btn
            >
          </q-item-section>
        </q-item>
      </q-list>
    </q-scroll-area>

    <!-- Add Project Dialog -->
    <q-dialog v-model="addNewProject" persistent>
      <q-card
        square
        class="bg-dark-soft border-lighter q-pa-md"
        style="min-width: 400px"
      >
        <q-card-section class="q-pb-none">
          <div class="text-h6 text-white text-weight-bold">New Project</div>
          <div class="text-caption text-grey-6">
            Initialize a new lazy-backend workspace.
          </div>
        </q-card-section>

        <q-card-section>
          <q-input
            dark
            dense
            standout
            v-model="newProjectName"
            placeholder="e.g. ecommerce-api"
            label="Project Name"
            autofocus
            @keyup.enter="handleCreateProject"
            :rules="[(val) => !!val || 'Field is required']"
          >
            <template v-slot:prepend>
              <q-icon name="bi-plus-square" size="xs" color="primary" />
            </template>
          </q-input>
        </q-card-section>

        <q-card-actions align="right" class="q-pt-md">
          <q-btn flat label="Cancel" color="grey-7" v-close-popup no-caps />
          <q-btn
            unelevated
            label="Create Project"
            color="primary"
            :loading="creating"
            @click="handleCreateProject"
            no-caps
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, onMounted, ref } from "vue";
import {
  ListProjects,
  CreateNewProject,
  DeleteProject,
} from "../../../wailsjs/go/main/App";
import { useAppContextStore } from "../../store/appContext";
import { useQuasar } from "quasar";

const $q = useQuasar();
const store = useAppContextStore();
const projects = ref<{ name: string; path: string }[] | undefined>();
const addNewProject = ref(false);
const newProjectName = ref("");
const creating = ref(false);

const loadProjects = async () => {
  try {
    const list = await ListProjects();
    console.log("Loaded projects:", list);
    if (!list) {
      
      projects.value = [];
      return;
    }
    projects.value = list.map((item) => {
      return {
        name: item.split(/[/\\]/).pop() || item,
        path: item,
      };
    });
  } catch (error) {
    console.error("Failed to load projects:", error);
    projects.value = [];
  }
};

onMounted(loadProjects);

const handleCreateProject = async () => {
  if (!newProjectName.value) return;

  creating.value = true;
  try {
    const res = await CreateNewProject(newProjectName.value);
    if (res.created) {
      $q.notify({
        type: "positive",
        message: res.message,
        timeout: 2000,
      });
      newProjectName.value = "";
      addNewProject.value = false;
      await loadProjects();
    } else {
      $q.notify({
        type: "negative",
        message: res.message,
      });
    }
  } catch (error) {
    $q.notify({
      type: "negative",
      message: "Failed to create project",
    });
  } finally {
    creating.value = false;
  }
};

const selectProject = (project: { name: string; path: string }) => {
  store.setCurrentProject(project.name);
  $q.notify({
    message: `Switched to project: ${project.name}`,
    color: "dark",
    icon: "bi-check2-circle",
    textColor: "primary",
    position: "bottom-right",
    timeout: 1500,
  });
};

async function deleteProject(project: { name: string; path: string }) {
  $q.dialog({
    title: "Confirm Delete",
    message: `Are you sure you want to delete project "${project.name}"? This action cannot be undone.`,
    cancel: true,
    persistent: true,
    dark: true,
  }).onOk(async () => {
    try {
      const res = await DeleteProject(project.path);
      if (res.success) {
        $q.notify({
          type: "positive",
          message: res.message,
          timeout: 2000,
        });
        await loadProjects();
        if (store.currentProject === project.name) {
          store.setCurrentProject(null);
        }
      } else {
        $q.notify({
          type: "negative",
          message: res.message,
        });
      }
    } catch (error) {
      $q.notify({
        type: "negative",
        message: "Error deleting project",
      });
    }
  });
}
</script>

<style lang="scss" scoped>
.min-width-auto {
  min-width: unset !important;
}

.project-item {
  border-radius: 8px;
  transition: all 0.2s;
  color: #888;

  &:hover {
    background: rgba(255, 255, 255, 0.03);
    color: #eee;
  }
}

.active-project {
  color: #fff !important;
  background: rgba(var(--q-primary), 0.1) !important;
}

.active-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--q-primary);
  box-shadow: 0 0 8px var(--q-primary);
}

.bg-dark-soft {
  background: #111 !important;
}

.border-lighter {
  border: 1px solid rgba(255, 255, 255, 0.05);
}
</style>
