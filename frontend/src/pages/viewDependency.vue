<template>
  <q-page class="doc-page">
    <div v-if="loading" class="flex flex-center q-pa-xl">
      <q-spinner size="40px" color="primary" />
    </div>

    <div v-else-if="!manifest" class="flex flex-center q-pa-xl text-grey">
      Manifest not found for "{{ dependencyName }}"
    </div>

    <div v-else class="q-pa-md" style="max-width: 900px; margin: 0 auto;">
      <!-- Header -->
      <div class="q-mb-lg">
        <div class="row items-center q-gutter-sm">
          <div class="text-h4">{{ dependencyName }}</div>
          <q-badge color="primary" outline>{{ manifest['lang-extension'] }}</q-badge>
          <q-badge color="grey-7" outline>v{{ manifest.version }}</q-badge>
          <q-badge v-if="installed" color="positive" outline>Installed</q-badge>
        </div>
        <div class="text-subtitle1 text-grey-8 q-mt-xs">{{ manifest.description }}</div>
        <div class="text-caption text-grey q-mt-xs">by {{ manifest.publisher }}</div>

        <!-- Action buttons -->
        <div class="row items-center q-gutter-sm q-mt-md">
          <q-btn
            v-if="!installed"
            color="primary"
            icon="download"
            label="Install"
            :loading="installing"
            @click="handleInstall"
          />
          <q-btn
            v-else-if="!usedInProject"
            color="secondary"
            icon="add_task"
            label="Use in Project"
            :loading="adding"
            @click="handleUseInProject"
          />
          <template v-else>
            <q-badge color="positive" outline class="q-pa-sm">
              <q-icon name="check_circle" class="q-mr-xs" />
              In use in this project
            </q-badge>
            <q-btn
              color="negative"
              outline
              icon="remove_circle"
              label="Remove from Project"
              :loading="removing"
              @click="handleRemoveFromProject"
            />
          </template>
        </div>

        <div v-if="manifest.install?.length" class="q-mt-md">
          <div v-for="(cmd, i) in manifest.install" :key="i" class="install-cmd row items-center justify-between q-pa-sm q-mb-xs">
            <code>{{ cmd }}</code>
            <q-btn flat dense round icon="content_copy" size="sm" @click="copy(cmd)" />
          </div>
        </div>
      </div>

      <!-- Search -->
      <q-input
        v-model="search"
        outlined
        dense
        clearable
        placeholder="Search exports by name, description, module, or param..."
        class="q-mb-lg"
      >
        <template #prepend><q-icon name="search" /></template>
      </q-input>

      <div v-if="search && matchCount === 0" class="text-grey text-center q-pa-lg">
        No exports match "{{ search }}"
      </div>

      <!-- Modules -->
      <div v-for="mod in filteredModules" :key="mod.path" class="q-mb-lg">
        <div class="row items-center q-mb-sm">
          <q-icon name="folder_open" class="q-mr-xs" color="grey-7" />
          <span class="text-subtitle2 text-grey-8">{{ mod.path }}</span>
          <q-badge color="grey-4" text-color="grey-9" class="q-ml-sm">
            {{ mod.exports.length }} export{{ mod.exports.length === 1 ? '' : 's' }}
          </q-badge>
        </div>

        <q-list bordered separator class="rounded-borders">
          <q-expansion-item
            v-for="exp in mod.exports"
            :key="exp.name"
            group="exports"
            default-opened
          >
            <template #header>
              <q-item-section>
                <div class="row items-baseline q-gutter-xs">
                  <span class="text-weight-bold code-font">{{ exp.name }}</span>
                  <span class="text-grey code-font">({{ paramSignature(exp.params) }})</span>
                  <q-icon name="arrow_forward" size="14px" class="q-mx-xs text-grey" />
                  <q-chip dense square color="deep-purple-1" text-color="deep-purple-9" class="code-font">
                    {{ exp.returnType }}
                  </q-chip>
                </div>
              </q-item-section>
            </template>

            <q-card flat>
              <q-card-section class="q-pt-none">
                <div v-if="exp.description" class="text-body2 text-grey-9 q-mb-sm">
                  {{ exp.description }}
                </div>

                <div v-if="Object.keys(exp.params || {}).length">
                  <div class="text-caption text-weight-medium text-grey-7 q-mb-xs">PARAMETERS</div>
                  <q-markup-table dense flat bordered class="params-table">
                    <thead>
                      <tr>
                        <th class="text-left">Name</th>
                        <th class="text-left">Type</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(type, pname) in exp.params" :key="pname">
                        <td class="code-font">{{ pname }}</td>
                        <td class="code-font text-primary">{{ type }}</td>
                      </tr>
                    </tbody>
                  </q-markup-table>
                </div>
                <div v-else class="text-caption text-grey">No parameters</div>
              </q-card-section>
            </q-card>
          </q-expansion-item>
        </q-list>
      </div>
    </div>
  </q-page>
</template>

<script lang="ts" setup>
import {
  AddProjectDependency,
  GetAllInstalledDependencies,
  GetDependencyManifest,
  GetProjectUsedDependencies,
  InstallDependencyModule,
  RemoveProjectDependency,
} from "#/wailsjs/go/main/App"
import { io } from "#/wailsjs/go/models"
import { computed, onMounted, ref } from "vue"
import { useRoute } from "vue-router"
import { copyToClipboard, useQuasar } from "quasar"
import { useAppContextStore } from "~/store/appContext"

const $q = useQuasar()
const route = useRoute()
const dependencyName = route.params["id"] as string
const appContext = useAppContextStore()

const manifest = ref<io.DependencySchema | null>(null)
const loading = ref(true)
const search = ref("")

// Two independent states:
// - installed: dependency exists in the app / on disk at all
// - usedInProject: the currently open project has it wired in (e.g. added to its manifest)
const installed = ref(false)
const usedInProject = ref(false)

const installing = ref(false)
const adding = ref(false)
const removing = ref(false)

async function refreshDependencyState() {
  const [manifestData, installedDependencies, projectDependencies] = await Promise.all([
    GetDependencyManifest(dependencyName),
    GetAllInstalledDependencies(),
    appContext.currentProject
      ? GetProjectUsedDependencies(appContext.currentProject)
      : Promise.resolve([]),
  ])

  manifest.value = manifestData
  installed.value = installedDependencies.some((dependency) => dependency.name === dependencyName)
  usedInProject.value = projectDependencies.some((dependency) => dependency.name === dependencyName)
}

onMounted(async () => {
  try {
    await refreshDependencyState()
  } finally {
    loading.value = false
  }
})

async function handleInstall() {
  installing.value = true
  try {
    const dependency = {
      name: dependencyName,
      description: manifest.value?.description ?? "",
      exportsCount: manifest.value?.modules?.reduce((count, module) => count + (module.exports?.length ?? 0), 0) ?? 0,
      version: Number(manifest.value?.version ?? 0),
      publisher: manifest.value?.publisher ?? "",
      location: "",
    } as io.Dependence

    await InstallDependencyModule(dependency)
    installed.value = true
    $q.notify({ message: `${dependencyName} installed`, color: "positive" })
  } catch (err) {
    $q.notify({ message: `Failed to install ${dependencyName}`, color: "negative" })
  } finally {
    installing.value = false
  }
}

async function handleUseInProject() {
  adding.value = true
  try {
    if (!appContext.currentProject) {
      throw new Error("Select a project before adding dependencies")
    }

    const dependency = {
      name: dependencyName,
      description: manifest.value?.description ?? "",
      exportsCount: manifest.value?.modules?.reduce((count, module) => count + (module.exports?.length ?? 0), 0) ?? 0,
      version: Number(manifest.value?.version ?? 0),
      publisher: manifest.value?.publisher ?? "",
      location: "",
    } as io.Dependence

    await AddProjectDependency(appContext.currentProject, dependency)
    usedInProject.value = true
    await refreshDependencyState()
    $q.notify({ message: `${dependencyName} added to project`, color: "positive" })
  } catch (err) {
    $q.notify({ message: `Failed to add ${dependencyName} to project`, color: "negative" })
  } finally {
    adding.value = false
  }
}

async function handleRemoveFromProject() {
  removing.value = true
  try {
    if (!appContext.currentProject) {
      throw new Error("Select a project before removing dependencies")
    }

    await RemoveProjectDependency(appContext.currentProject, dependencyName)
    usedInProject.value = false
    await refreshDependencyState()
    $q.notify({ message: `${dependencyName} removed from project`, color: "positive" })
  } catch (err) {
    $q.notify({ message: `Failed to remove ${dependencyName} from project`, color: "negative" })
  } finally {
    removing.value = false
  }
}

function paramSignature(params: Record<string, any> | undefined): string {
  if (!params) return ""
  return Object.entries(params)
    .map(([name, type]) => `${name}: ${type}`)
    .join(", ")
}

function copy(text: string) {
  copyToClipboard(text).then(() => {
    $q.notify({ message: "Copied", color: "positive", timeout: 800 })
  })
}

const filteredModules = computed(() => {
  const q = search.value.trim().toLowerCase()
  const modules = manifest.value?.modules ?? []
  if (!q) return modules

  return modules
    .map((mod) => {
      const exports = (mod.exports ?? []).filter((exp) => {
        const haystack = [
          exp.name,
          exp.description ?? "",
          exp.returnType,
          mod.path,
          ...Object.keys(exp.params ?? {}),
          ...Object.values(exp.params ?? {}),
        ]
          .join(" ")
          .toLowerCase()
        return haystack.includes(q)
      })
      return { ...mod, exports }
    })
    .filter((mod) => mod.exports.length > 0)
})

const matchCount = computed(() =>
  filteredModules.value.reduce((sum, m) => sum + m.exports.length, 0)
)
</script>

<style scoped>
.code-font {
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 0.9em;
}
.install-cmd {
  background: #1e1e1e;
  color: #d4d4d4;
  border-radius: 6px;
  font-family: 'Fira Code', 'Consolas', monospace;
}
.params-table {
  max-width: 500px;
}
</style>