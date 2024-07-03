<script setup lang="ts">
import './style.css'
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import '@fontsource/roboto/900.css';
import '@material/web/progress/linear-progress.js';
import { ref } from 'vue';
const loaded = ref(false);
const receivedReady = ref(false);
const src = ref("")
const embedKey = "<EMBEDKEY.THIS WILL BE REPLACED DURING INSTALLATION>";

window.addEventListener("message", (event) => {
  if (event.data == "ready") {

    receivedReady.value = true;
  }
  if (event.data == "reload") {
    (document.getElementById("mainframe") as HTMLIFrameElement).src += "";
    src.value = urlGenerator();
  }
});

function tryConnect() {
  fetch("http://127.0.0.1:8080/generate_200").then((req2) => {
    if (req2.status == 200) {
      loaded.value = true;
      src.value = urlGenerator();
    } else {
      setTimeout(tryConnect, 1000);
    }
  }).catch(() => {
    setTimeout(tryConnect, 1000);
  })
}

function urlGenerator() {
  let url = new URL("http://127.0.0.1:8080");
  url.searchParams.append("mode", "wallpaper");
  url.searchParams.append("embedkey", embedKey);
  return url.toString();
}

tryConnect()
</script>

<template>
  <iframe :src="src" v-show="loaded && receivedReady" class="w-full h-full" id="mainframe">

  </iframe>
  <div style="height: 100%; background-color: transparent;" class="text-white" v-show="!loaded || !receivedReady">
    <img src="./assets/background.png" class="absolute object-cover w-full h-full">
    </img>
    <div class="absolute card right-0 bottom-0 m-4 mb-16 p-4" style="
      width: 256px;
      height: 80px;

    ">
      <div class="font-bold">Loading wallpaper..</div>
      <div class="mt-4">

      </div>
      <md-linear-progress indeterminate></md-linear-progress>
    </div>
  </div>
</template>

<style scoped></style>
