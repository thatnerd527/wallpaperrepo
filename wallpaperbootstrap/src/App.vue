<script setup lang="ts">
import './style.css'
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import '@fontsource/roboto/900.css';
import '@material/web/progress/linear-progress.js';
import { ref } from 'vue';
const loaded = ref(false);

while (true) {
  try {
    let req = await fetch("http://127.0.0.1:8080/generate_200");
    if (req.status == 200) {
      window.location.href = "http://127.0.0.1:8080?mode=wallpaper";
      break;
    }
  } catch (e) {
    console.log(e);
  }
}

</script>

<template>
  <div style="height: 100%; background-color: transparent;" class="text-white" v-show="!loaded">
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

<style scoped>
</style>
