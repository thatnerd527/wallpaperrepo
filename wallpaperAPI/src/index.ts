

import { GetAddons } from "./AddonLoader";
import { PanelManagement } from "./endpoints/panel";
import { PopupDummy } from "./endpoints/popup";
import { SharingAndIPC } from "./endpoints/sharing";
import {Environment, generateGUID, EndpointRegister2} from './IPC';
import { StorageManager } from "./StorageManager";
import { StorageDummy } from "./endpoints/storage";
import { SharingDummy } from "./endpoints/sharing";
import { PanelDummy } from "./endpoints/panel";

import * as Type from "./types";
import { BackgroundDummy } from "./endpoints/background";
import { AddonDummy } from "./endpoints/addon";

StorageDummy();
SharingDummy();
PopupDummy();
PanelDummy();
BackgroundDummy();
AddonDummy();

console.log(Environment.endpointRegister.generateClientSDK());
Environment.controlPort =
  new URL(window.location.href).searchParams.get("controlPort") == null
    ? 8080
    : parseInt(new URL(window.location.href).searchParams.get("controlPort")!);

Environment.addons = [];
console.log("Control Port: " + Environment.controlPort);
console.log(Environment.controlPort);
GetAddons().then((addons) => {
  Environment.addons = addons;
});

declare global {
  interface Window {
    setallowedorigins: (origins: string[]) => void;
  }
}


export default Type;
