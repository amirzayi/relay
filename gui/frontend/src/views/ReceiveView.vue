<script setup>
import { ref, onMounted } from "vue";
import { mdiPower } from "@mdi/js";
import { useRouter } from "vue-router";
import BaseButton from "@/components/BaseButton.vue";
import Layout from "@/layouts/Layout.vue";
import { Listen, Close } from "@/../wailsjs/go/gui/Receiver";
import NotificationBar from "@/components/NotificationBar.vue";
import TableFiles from "@/components/TableFiles.vue";
import { EventsOn } from "@/../wailsjs/runtime/runtime";
import { gradientBgGreenCyan } from "@/colors";
import SectionBanner from "@/components/SectionBanner.vue";

const err = ref(null);
const clientIPAddress = ref(null);
const files = ref({});

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

const addFile = (fileName, size) => {
  let name = fileName.split(/[\\/]/).pop();
  let newFile = {
    Name: name,
    Size: size,
    Progress: 0,
    Done: false,
    Err: null,
  };
  files.value[fileName] = newFile;
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
  Listen()
    .then((ip) => {
      clientIPAddress.value = ip;
    })
    .catch((error) => {
      err.value = error;
    });

  EventsOn("receiving", (...data) => {
    const eventName = data[0];
    const fileName = data[1];
    const value = data[2];
    switch (eventName) {
      case "preparing":
        addFile(fileName, value);
        break;
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
    </div>
    <SectionBanner :class="gradientBgGreenCyan" v-if="clientIPAddress">
      <h1 class="text-xl text-white mb-6">
        Connected to {{ clientIPAddress }}
      </h1>
    </SectionBanner>
    <TableFiles :items="files" />
  </Layout>
</template>
