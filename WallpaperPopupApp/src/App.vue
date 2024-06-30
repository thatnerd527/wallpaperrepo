<script setup lang="ts">

import '@fontsource/roboto';
import '@fontsource-variable/material-symbols-outlined';
import '@fontsource/material-symbols-outlined';
import '@material/web/icon/icon.js';
import '@material/web/button/outlined-button.js';
import PopupConfirmation from './components/PopupConfirmation.vue';
import InputMethod from './components/InputMethod.vue';
import { ref } from 'vue';

const windowtype = ref("waiting")

const popupurl = ref("")
const popupclientid = ref("")
const popupappname = ref("")
const popupfavicon = ref("")
const popuptitle = ref("")


const inputtype = ref("")
const inputplaceholder = ref("")
const inputmaxlength = ref(-1)

const trackingid = ref("")


window.electronAPI.addAwaiter((message) => {
  console.log(message)
  switch (message.type) {
    case 'popup':
      popupurl.value = message.url
      popupclientid.value = message.clientid
      popupappname.value = message.appname
      popupfavicon.value = message.favicon
      popuptitle.value = message.title
      trackingid.value = message.trackingid
      windowtype.value = "popup"
      break;
    case 'input':
      inputtype.value = message.inputtype
      inputplaceholder.value = message.inputplaceholder
      inputmaxlength.value = message.inputmaxlength
      trackingid.value = message.trackingid

      windowtype.value = "input"
      break;
  }
})
window.electronAPI.tellready()

</script>

<template>
  <div class="flex w-full h-full justify-center items-center dark:text-white text-black" v-if="windowtype == 'waiting'">
    Waiting for message from main process....
  </div>
  <PopupConfirmation v-if="windowtype == 'popup'"
    :popupurl="popupurl"
    :popupclientid="popupclientid"
    :popupappname="popupappname"
    :popupfavicon="popupfavicon"
    :popuptitle="popuptitle"
    :trackingid="trackingid"
  />
  <InputMethod v-if="windowtype == 'input'"
    :inputtype="inputtype"
    :inputplaceholder="inputplaceholder"
    :inputmaxlength="inputmaxlength"
    :trackingid="trackingid"/>
</template>

