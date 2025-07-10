import { createRouter, createWebHashHistory } from "vue-router";
import Home from "@/views/HomeView.vue";

const routes = [
  {
    meta: {
      title: "Home",
    },
    path: "/",
    name: "home",
    component: Home,
  },
  {
    meta: {
      title: "Lookup",
    },
    path: "/lookup",
    name: "lookup",
    component: () => import("@/views/LookupView.vue"),
  },

  {
    meta: {
      title: "Send",
    },
    path: "/send",
    name: "send",
    component: () => import("@/views/SendView.vue"),
  },
  {
    meta: {
      title: "Receive",
    },
    path: "/receive",
    name: "receive",
    component: () => import("@/views/ReceiveView.vue"),
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { top: 0 };
  },
});

export default router;
