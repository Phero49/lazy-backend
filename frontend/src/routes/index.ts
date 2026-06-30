import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

// Create simple placeholders for areas that don't have pages yet
const Placeholder = (title: string) => ({
  template: `<div class="tw-p-4"><div class="text-overline text-grey-7 q-mb-md">${title}</div><div class="text-caption text-grey-5">Content for ${title} coming soon...</div></div>`,
});

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    component: () => import("../layout/MainLayout.vue"),
    children: [
      {
        path: "",
        components: {
          default: () => import("../pages/IndexPage.vue"),
          
          right: Placeholder("Inspector"),
        },
      },
      {
        path: "project",
        components: {
          default: () => import("../pages/IndexPage.vue"),
          
          right: Placeholder("Project Inspector"),
        },
      },
      {
        path: "database",
        components: {
          default: Placeholder("Database Designer"),
          
          right: Placeholder("DB Inspector"),
        },
      },
      {
        path: "database/schema/:connection/:table",
        components: {
          default: () => import("../pages/TableSchemaPage.vue"),
          
          right: Placeholder("Schema Inspector"),
        },
      },
      {
        path: "database/query/:connection/:table",
        components: {
          default: () => import("../pages/QueryBuilder.vue"),
          
          right: Placeholder("Query Inspector"),
        },
      },
      {
        path: "/api/:endPoint(.*)*",
        components: {
          default: () => import("../pages/ApiBuilderPage.vue") ,
          
          right: Placeholder("API Inspector"),
        },
      },
      {
        path: "plugins",
        components: {
          default: Placeholder("Plugin Market"),
          
          right: Placeholder("Plugin Config"),
        },
      },
      {
        path: "env",
        components: {
          default: Placeholder("Environment Variables"),
          
          right: Placeholder("System Health"),
        },
      },
      {
        path: "env/edit/:category/:name",
        components: {
          default: () => import("../pages/EnvEditorPage.vue"),
          
          right: Placeholder("Config Inspector"),
        },
      },
      {
        name: "update-env",
        path: "env/update/:type",
        components: {
          default: () => import("../pages/UpdateEnvPage.vue"),
          
          right: Placeholder("Update Monitor"),
        },
      },
    ],
  },
  {
    path: "/:catchAll(.*)*",
    component: () => import("../pages/ErrorNotFound.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
