<script setup>
import { mdiCheck, mdiClose } from "@mdi/js";
import BaseIcon from "./BaseIcon.vue";

defineProps({
  items: Object,
});

const mbSize = (bytes) => {
  if (bytes === 0) return "0 Bytes";

  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const value = parseFloat(bytes / Math.pow(k, i)).toFixed(0);

  return `${value} ${sizes[i]}`;
};
</script>

<template>
  <table>
    <thead>
      <tr>
        <th class="w-4" />
        <th>Name</th>
        <th>Size</th>
        <th />
      </tr>
    </thead>
    <tbody>
      <tr v-for="item in items">
        <td>
          <BaseIcon v-if="item.Err" :path="mdiClose" />
          <BaseIcon v-if="item.Done" :path="mdiCheck" />
        </td>
        <td data-label="Name">
          {{ item.Name }}{{ item.Err ? ` (error: ${item.Err})` : "" }}
        </td>
        <td data-label="Size">{{ mbSize(item.Size) }}</td>
        <td data-label="Progress" class="lg:w-52">
          <progress
            class="flex w-2/5 self-center lg:w-full"
            max="100"
            :value="item.Progress || 0"
          >
            {{ item.Progress || 0 }}%
          </progress>
        </td>
      </tr>
    </tbody>
  </table>
</template>
