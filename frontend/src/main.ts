import { createApp } from "vue";
import App from "./App.vue";
import { Quasar, Notify, Dialog } from "quasar";
import quasarIconSet from "quasar/icon-set/fontawesome-v6";
import router from "./routes";
import { createPinia } from "pinia";

// Import icon libraries
import "@quasar/extras/material-icons/material-icons.css";
// ..required because of selected iconSet:
import "@quasar/extras/bootstrap-icons/bootstrap-icons.css";

// Import Quasar css
import "quasar/src/css/index.sass";

import  './css/style.css'
const myApp = createApp(App);
const pinia = createPinia();

myApp.use(pinia);
myApp.use(router);
myApp.use(Quasar, {
  config: {
    dark: true,
  },
  plugins: { Notify, Dialog }, // import Quasar plugins and add here
  iconSet: quasarIconSet,
});

// Assumes you have a <div id="app"></div> in your index.html
myApp.mount("#app");
