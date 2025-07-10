<script setup>
import { mdiRefresh, mdiLink } from "@mdi/js";
import { useRouter } from "vue-router";
import BaseButton from "@/components/BaseButton.vue";
import { ref, onMounted } from "vue";
import { Lookup, Connect } from "@/../wailsjs/go/gui/Sender";
import Layout from "@/layouts/Layout.vue";
import CardBoxDevice from "@/components/CardBoxDevice.vue";
import SectionBanner from "@/components/SectionBanner.vue";
import { gradientBgRedSilver } from "@/colors";
import NotificationBar from "@/components/NotificationBar.vue";
const router = useRouter();

const goToSend = () => {
  router.push("/send");
};

const availableDeviceIpAddresses = ref(null);
const err = ref(null);

const lookupAvailableDevices = () => {
  Lookup()
    .then((data) => {
      availableDeviceIpAddresses.value = data;
    })
    .catch((error) => {
      err.value = error;
    });
};

const connectToDevice = (ip) => {
  // Implement the logic to connect to the device with the given IP address
  Connect(ip)
    .then(() => {
      goToSend();
    })
    .catch((error) => {
      err.value = error;
    });
};
onMounted(() => {
  lookupAvailableDevices();
});
</script>

<template>
  <Layout>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
      <div
        class="flex flex-col justify-between border-sky-400 boder-2"
        v-for="ip in availableDeviceIpAddresses"
      >
        <CardBoxDevice :ip="ip">
          <BaseButton
            :icon="mdiLink"
            label="Connect"
            color="info"
            rounded-full
            @click="connectToDevice(ip)"
          />
        </CardBoxDevice>
      </div>
    </div>

    <div v-if="err">
      <NotificationBar color="danger" :icon="mdiTableOff">
        <template #right></template>
        <b>{{ err }}</b>
      </NotificationBar>
    </div>
    <SectionBanner :class="gradientBgRedSilver">
      <h1 class="text-xl text-black mb-6">Didn't find device?</h1>
      <div>
        <BaseButton
          :icon="mdiRefresh"
          label="Refresh"
          rounded-full
          @click="lookupAvailableDevices"
        />
      </div>
    </SectionBanner>
  </Layout>
</template>
