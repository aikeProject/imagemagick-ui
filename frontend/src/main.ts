import { createApp } from "vue";
import App from "App.vue";
import router from "router";
import store from "store";
import * as Wails from "@wailsapp/runtime";
import installElementPlus from "plugins/element";
// import initAnd from "plugins/ant";
import "assets/css/tailwind.css";
import "assets/css/main.css";

function RunApp() {
  const app = createApp(App);
  installElementPlus(app);
  // initAnd(app);
  app.use(router);
  app.use(store);
  app.mount("#app");
}

if (process.env.VUE_APP_NO_WAILS) {
  RunApp();
} else {
  Wails.Init(RunApp);
}
