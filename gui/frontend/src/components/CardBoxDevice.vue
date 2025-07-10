<script setup>
import { mdiDevices } from "@mdi/js";
import CardBox from "@/components/CardBox.vue";
import BaseLevel from "@/components/BaseLevel.vue";
import IconRounded from "@/components/IconRounded.vue";

const props = defineProps({
  ip: {
    type: String,
    required: true,
  },
});
</script>

<template>
  <CardBox class="mb-6 last:mb-0">
    <BaseLevel>
      <BaseLevel type="justify-start">
        <IconRounded :icon="mdiDevices" color="success" class="md:mr-6" />
        <div class="text-center space-y-1 md:text-left md:mr-6">
          <h4 class="text-xl">{{ ip }}</h4>
        </div>
      </BaseLevel>
      <slot />
    </BaseLevel>
  </CardBox>
</template>
