const { contextBridge, ipcRenderer } = require("electron");

//ipcRenderer.invoke("ready");
let trackingID = process.argv.filter((arg) => {
    return arg.startsWith("--trackingid=");
}).map((arg) => { return arg.split("=")[1]; })[0];

contextBridge.exposeInMainWorld("wallpaperAPI", {
    popupcancel: () => ipcRenderer.send(`${trackingID}-popupcancel`),
    popupsuccess: (data) => ipcRenderer.send(`${trackingID}-popupsuccess`, data),
});