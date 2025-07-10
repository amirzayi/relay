<script setup>
import { useRouter } from "vue-router";
import { mdiDownload, mdiUpload } from "@mdi/js";
import Layout from "@/layouts/Layout.vue";
import BaseButton from "@/components/BaseButton.vue";

const router = useRouter();

const goToReceive = () => {
  router.push("/receive");
};

const goToLookup = () => {
  router.push("/lookup");
};
</script>

<template>
  <Layout>
    <div class="flex flex-col items-center justify-center gap-8 min-h-[calc(100vh-3.5rem)]">
      <div class="flex justify-normal">
        <img
          class="w-auto h-20 md:h-17 rounded-md dark:invert"
          src="@/logo.png"
          alt="logo"
        />
      </div>

      <div class="flex flex-col md:flex-row gap-4 max-w-xs md:max-w-md">
        <BaseButton
          roundedFull
          :icon="mdiUpload"
          label="Send"
          color="success"
          outline
          @click="goToLookup"
        />
        <BaseButton
          roundedFull
          :icon="mdiDownload"
          label="Receive"
          color="info"
          outline
          @click="goToReceive"
        />
      </div>
    </div>
  </Layout>
</template>
