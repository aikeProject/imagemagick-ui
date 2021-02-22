import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";
import Home from "views/Home.vue";
import About from "views/About.vue";
import Setting from "views/Setting.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: Home
  },
  {
    path: "/setting",
    name: "Setting",
    component: Setting
  },
  {
    path: "/about",
    name: "About",
    component: About
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
});

export default router;
