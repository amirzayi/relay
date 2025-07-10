<script setup>
import { ref, onMounted } from "vue";
import { mdiFolder, mdiFile, mdiPower } from "@mdi/js";
import { useRouter } from "vue-router";
import BaseButton from "@/components/BaseButton.vue";
import Layout from "@/layouts/Layout.vue";
import NotificationBar from "@/components/NotificationBar.vue";
import TableFiles from "@/components/TableFiles.vue";
import { EventsOn } from "@/../wailsjs/runtime/runtime";
import { Close, SendDirectory, SendFiles } from "@/../wailsjs/go/gui/Sender";

const err = ref(null);
const files = ref(null);

const router = useRouter();

const cancel = () => {
  Close()
    .then(() => {
      router.push("/");
    })
    .catch((error) => {
      err.value = error;
    });
};

const selectDirectory = () => {
  SendDirectory()
    .then((data) => {
      const transformedData = data.reduce((acc, obj) => {
        const key = Object.values(obj)[2]; // Get the value of the first property
        acc[key] = obj; // Use that value as the new key
        return acc;
      }, {});
      files.value = { ...files.value, ...transformedData };
    })
    .catch((error) => {
      err.value = error;
    });
};

const selectFiles = () => {
  SendFiles()
    .then((data) => {
      const transformedData = data.reduce((acc, obj) => {
        const key = Object.values(obj)[2]; // Get the value of the first property
        acc[key] = obj; // Use that value as the new key
        return acc;
      }, {});
      files.value = { ...files.value, ...transformedData };
    })
    .catch((error) => {
      err.value = error;
    });
};

const updateProgress = (key, progress) => {
  const item = files.value[key];
  if (!item) return;
  item.Progress = progress;
};

const failProgress = (key, err) => {
  const item = files.value[key];
  if (!item) return;
  item.Err = err;
};

const completeProgress = (key) => {
  const item = files.value[key];
  if (!item) return;
  item.Done = true;
};

onMounted(() => {
  EventsOn("sending", (...data) => {
    const eventName = data[0];
    const fileName = data[1];
    const value = data[2];
    switch (eventName) {
      case "inProgress":
        updateProgress(fileName, value);
        break;
      case "failed":
        failProgress(fileName, value);
        break;
      case "completed":
        completeProgress(fileName);
        break;
      case "communication":
        err.value = value;
        break;
    }
  });
});
</script>

<template>
  <Layout>
    <div v-if="err">
      <NotificationBar color="danger" v-if="err" :icon="mdiTableOff">
        <template #right></template>
        <b>{{ err }}</b>
      </NotificationBar>
    </div>
    <div class="flex flex-col items-center gap-8 p-10 mb-8">
      <div class="flex justify-normal">
        <BaseButton
          class="p-4"
          iconSize="24"
          roundedFull
          :icon="mdiPower"
          color="warning"
          @click="cancel"
        />
      </div>
      <div class="flex flex-col md:flex-row gap-4 max-w-xs md:max-w-md">
        <BaseButton
          :icon="mdiFile"
          label="Select File"
          color="info"
          outline
          @click="selectFiles"
        />
        <BaseButton
          :icon="mdiFolder"
          label="Select Folder"
          color="success"
          outline
          @click="selectDirectory"
        />
      </div>
    </div>
    <TableFiles :items="files" />
  </Layout>
</template>
