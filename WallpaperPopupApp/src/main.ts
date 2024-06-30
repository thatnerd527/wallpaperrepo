import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import {
    fluentButton,
    fluentTextField,
    fluentTextArea,
  provideFluentDesignSystem,
} from "@fluentui/web-components";

provideFluentDesignSystem().register(fluentButton(), fluentTextField(), fluentTextArea());
createApp(App).mount('#app')
