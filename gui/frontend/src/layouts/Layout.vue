<script setup>
import { useRouter } from "vue-router";
import menuNavBar from "@/menuNavBar.js";
import { useDarkModeStore } from "@/stores/darkMode.js";
import NavBar from "@/components/NavBar.vue";

const darkModeStore = useDarkModeStore();

const router = useRouter();

const menuClick = (event, item) => {
  if (item.isToggleLightDark) {
    darkModeStore.set();
  }
};
</script>

<template>
  <div>
    <div
      class="min-h-screen w-screen transition-position lg:w-auto bg-gray-50 dark:bg-slate-800 dark:text-slate-100"
    >
      <NavBar :menu="menuNavBar" @menu-click="menuClick"> </NavBar>
      <div class="pt-14" />
      <slot />
    </div>
  </div>
</template>
