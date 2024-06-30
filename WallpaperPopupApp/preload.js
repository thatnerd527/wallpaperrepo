const { contextBridge, ipcRenderer } = require('electron')
const awaiters = [];

let trackingID = process.argv
  .filter((arg) => {
    return arg.startsWith("--trackingid=");
  })
  .map((arg) => {
    return arg.split("=")[1];
  })[0];

ipcRenderer.on('data', (event, data) => {
    awaiters.forEach((awaiter) => awaiter(data));
});

contextBridge.exposeInMainWorld("electronAPI", {
  tellready: () => ipcRenderer.send(`${trackingID}-ready`),
  addAwaiter: (awaiter) => awaiters.push(awaiter),
  popupcancel: () => ipcRenderer.send(`${trackingID}-popupcancel`),
  openwebsite: (data) => ipcRenderer.send(`${trackingID}-openwebsite`, data),
});