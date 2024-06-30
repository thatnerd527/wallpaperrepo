const { contextBridge, ipcRenderer } = require("electron");
const awaiters = [];

//ipcRenderer.invoke("ready");
let trackingID = process.argv
  .filter((arg) => {
    return arg.startsWith("--trackingid=");
  })
  .map((arg) => {
    return arg.split("=")[1];
  })[0];

ipcRenderer.on("data", (event, data) => {
  awaiters.forEach((awaiter) => awaiter(data));
});

contextBridge.exposeInMainWorld("electronAPI", {
  tellready: () => ipcRenderer.send(`${trackingID}-ready`),
  addAwaiter: (awaiter) => awaiters.push(awaiter),
  inputcancel: () => ipcRenderer.send(`${trackingID}-inputcancel`),
  inputsuccess: (data) => ipcRenderer.send(`${trackingID}-inputsuccess`, data),
});