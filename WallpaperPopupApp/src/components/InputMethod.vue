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
import { ref } from 'vue';

const props = defineProps({
    inputtype: String,
    inputplaceholder: String,
    inputmaxlength: Number,
    trackingid: String,
});



const valuedata = ref("")
function done() {
    window.electronAPI.inputsuccess(valuedata.value);
}
function cancel() {
    window.electronAPI.inputcancel();
}
console.log(props.inputplaceholder)
</script>

<template>
    <div class="w-full min-h-14 p-4 flex flex-row justify-end items-center align-middle" style="
    background-color: rgba(0, 0, 0, 0.2);
    -webkit-app-region: drag;
  ">
        <div style="
      -webkit-app-region: no-drag;
    ">
            <md-icon-button aria-label="Add new" @click="cancel">
                <md-icon>close</md-icon>
            </md-icon-button>
        </div>
    </div>

    <div class="p-8 flex flex-col h-full">
        <md-outlined-text-field class="h-full" v-bind:type="inputtype" @input="(v) => {
                valuedata = v.target.value
            //console.log(v.target.value)
        }" :maxLength="inputmaxlength" :label="inputplaceholder">
        </md-outlined-text-field>
        <div class="min-w-4">
            &nbsp;
        </div>
        <div class="flex flex-row w-full">
            <md-filled-tonal-button class="w-full" @click="done">
                Done
                <md-icon slot="icon">done</md-icon>
            </md-filled-tonal-button>
            <div class="w-4"></div>
            <md-text-button class="w-full" @click="cancel">
                Cancel
                <md-icon slot="icon">close</md-icon>
            </md-text-button>
        </div>

    </div>
</template>