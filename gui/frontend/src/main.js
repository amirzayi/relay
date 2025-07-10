import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "@/css/main.css";

const pinia = createPinia();

(async () => {
  try {
    createApp(App).use(pinia).use(router).mount("#app");
  } catch (error) {
    console.error("Failed to initialize app:", error);
    alert(error);
  }
})();

import { useDarkModeStore } from "@/stores/darkMode";

const darkModeStore = useDarkModeStore(pinia);

if (
  (!localStorage["darkMode"] &&
    window.matchMedia("(prefers-color-scheme: dark)").matches) ||
  localStorage["darkMode"] === "1"
) {
  darkModeStore.set(true);
}

const defaultDocumentTitle = "Relay File Sharing";

// Set document title from route meta
router.afterEach((to) => {
  document.title = to.meta?.title
    ? `${to.meta.title} — ${defaultDocumentTitle}`
    : defaultDocumentTitle;
});
