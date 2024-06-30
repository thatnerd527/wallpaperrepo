<script setup lang="ts">
import '@material/web/button/outlined-button.js';
import '@material/web/button/text-button.js';
import '@material/web/icon/icon.js';
import '@material/web/iconbutton/icon-button.js';
import '@material/web/textfield/filled-text-field.js';
import '@material/web/textfield/outlined-text-field.js';
import '@material/web/button/filled-button.js';
import '@material/web/button/filled-tonal-button.js';
import '@material/web/iconbutton/filled-tonal-icon-button.js';
import '@material/web/iconbutton/outlined-icon-button.js';
import '@fontsource/roboto';
import { ref } from 'vue';

const props = defineProps({
  popupurl: String,
  popupclientid: String,
  popupappname: String,
  popupfavicon: String,
  popuptitle: String,
  trackingid: String,
});

const loading = ref(false);

function open() {
    loading.value = true;
    window.location.href = props.popupurl!;
    window.electronAPI.openwebsite({
        url: props.popupurl,
        clientid: props.popupclientid,
        appname: props.popupappname,
        favicon: props.popupfavicon,
        title: props.popuptitle,
        trackingid: props.trackingid,
    });
}

function close() {
    window.electronAPI.popupcancel();
}

</script>

<template>
    <div class="w-full min-h-14 p-4 flex flex-row justify-end items-center align-middle" style="
    background-color: rgba(0, 0, 0, 0.2);
    -webkit-app-region: drag;
  ">
        <div style="
      -webkit-app-region: no-drag;
    ">
            <md-icon-button aria-label="Add new" @click="close">
                <md-icon>close</md-icon>
            </md-icon-button>
        </div>
    </div>
    <div class="w-full h-full flex flex-col text-black dark:text-white align-middle justify-center items-center" style="
        font-size: x-large;
    ">
        <span class="text-center">
            &quot;<b>{{ popupappname }}</b>&quot;
            &nbsp;
            wants to open a window to: &quot;<b>{{ popupurl }}</b>&quot;
        </span>
        <div class="h-16">
        </div>
        <div class="flex flex-row" v-if="!loading">
            <md-outlined-button class="w-1/2" @click="close">
                <md-icon slot="icon">close</md-icon>
                Deny
            </md-outlined-button>
            <div class="w-4"></div>
            <md-filled-tonal-button class="w-1/2" @click="open">
                <md-icon slot="icon">check</md-icon>
                Allow
            </md-filled-tonal-button>
        </div>
        <div class="h-16"></div>
        <div v-if="loading">Loading...</div>
    </div>

</template>