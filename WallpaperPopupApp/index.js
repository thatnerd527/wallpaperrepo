import { app, BrowserWindow, ipcMain, session } from "electron";
import { set } from "lodash-es";
import path from "path";
import WebSocket, { WebSocketServer } from "ws";
import * as protocol from "./src/protocol/protocol.js"
const __dirname = import.meta.dirname;
let socket = null;

const DEV =
  process.argv.filter((arg) => arg.includes("--controlport=")).length == 0;

let port = () => {
  if (DEV) {
    return 8080;
  } else {
    return process.argv
      .filter((arg) => {
        return arg.startsWith("--controlport=");
      })
      .map((arg) => {
        return arg.split("=")[1];
      })[0];
  }
};

connectWebSocket();

const load = (win) => {
  if (DEV) {
    win.loadURL("http://localhost:5173/");
  } else {
    win.loadFile("./dist/index.html");
  }
};

const createHiddenWindow = () => {
  const window = new BrowserWindow({
    show: false,
  });
  load(window);
};

const createWindow = () => {
  //createSecurePopupWindow("https://www.google.com", "clientid", "appname", "favicon", "title", "trackingid");
  createInputWindow("textarea", "Enter your name", 50, "trackingid");
  //createPopup("https://www.google.com", "clientid", "appname", "favicon", "title", "trackingid");
};

const sendPopupResult = (trackingid, data) => {
  socket.send(
    protocol.PopupAppResponse.encode(
      new protocol.PopupAppResponse({
        type: protocol.PopupAppResponse.MessageType.POPUP,
        cancelled: false,
        popupResponse: {
          requestID: trackingid,
          resultData: data || "",
        },
      })
    ).finish()
  );
};

const sendPopupCancel = (trackingid) => {
  socket.send(
    protocol.PopupAppResponse.encode(
      new protocol.PopupAppResponse({
        type: protocol.PopupAppResponse.MessageType.POPUP,
        cancelled: true,
        popupResponse: {
          requestID: trackingid,
          resultData: "",
        },
      })
    ).finish()
  );
};

const sendInputResult = (trackingid, data) => {
  socket.send(
    protocol.PopupAppResponse.encode(
      new protocol.PopupAppResponse({
        type: protocol.PopupAppResponse.MessageType.INPUT,
        cancelled: false,
        inputResponse: {
          requestID: trackingid,
          resultData: data || "",
        },
      })
    ).finish()
  );
};

const sendInputCancel = (trackingid) => {
  socket.send(
    protocol.PopupAppResponse.encode(
      new protocol.PopupAppResponse({
        type: protocol.PopupAppResponse.MessageType.INPUT,
        cancelled: true,
        inputResponse: {
          requestID: trackingid,
          resultData: "",
        },
      })
    ).finish()
  );
};

const createSecurePopupWindow = (
  url,
  clientid,
  appname,
  favicon,
  title,
  trackingid
) => {
  const window = new BrowserWindow({
    width: 650,
    height: 450,
    webPreferences: {
      preload: path.join(__dirname, "preload_popup.js"),
      additionalArguments: [`--trackingid=${trackingid}`],
      session: session.fromPartition(`persist:popup_${clientid}`),
    },
  });
  let url2 = new URL(url);
  let cancelled = false;
  url2.searchParams.append("trackingid", trackingid);
  ipcMain.once(`${trackingid}-popupcancel`, () => {
    if (!cancelled) {
      sendPopupCancel(trackingid);
      cancelled = true;
      window.close();
    }
  });
  ipcMain.once(`${trackingid}-popupsuccess`, (event, data) => {
    if (!cancelled) {
      sendPopupResult(trackingid, data);
      cancelled = true;
    }
  });
  window.webContents.on("close", () => {
    if (!cancelled) {
      sendPopupCancel(trackingid);
      cancelled = true;
    }
  });
  window.loadURL(url2.href);
  window.setTitle(title);
  window.setMenu(null);
  try {
    window.setIcon(favicon);
  } catch (e) {
    console.error(e);
  }
};

const createPopup = (url, clientid, appname, favicon, title, trackingid) => {
  if (clientid == "system") {
    createSecurePopupWindow(url, clientid, appname, favicon, title, trackingid);
    return;
  }
  const win = new BrowserWindow({
    width: 650,
    height: 450,
    backgroundMaterial: "acrylic",

    transparent: true,
    frame: false,
    webPreferences: {
      preload: path.join(__dirname, "preload.js"),
      additionalArguments: [`--trackingid=${trackingid}`],
    },
  });
  let cancelled = false;
  win.on("close", () => {
    if (!cancelled) {
      sendPopupCancel(trackingid);
      cancelled = true;
    }
  });
  ipcMain.once(`${trackingid}-ready`, () => {
    win.webContents.send("data", {
      url: url,
      clientid: clientid,
      appname: appname,
      favicon: favicon,
      title: title,
      trackingid: trackingid,
      type: "popup",
    });

    ipcMain.once(`${trackingid}-openwebsite`, (event, data) => {
      createSecurePopupWindow(
        data.url,
        data.clientid,
        data.appname,
        data.favicon,
        data.title,
        data.trackingid
      );
      if (!cancelled) {
        cancelled = true;
        win.close();
      }
    });
    ipcMain.once(`${trackingid}-popupcancel`, () => {
      if (!cancelled) {
        sendPopupCancel(trackingid);
        cancelled = true;
        win.close();
      }
    });
  });

  load(win);
};

const createInputWindow = (inputtype, placeholder, maxlength, trackingid) => {
  const win = new BrowserWindow({
    width: 650,
    height: 350,
    backgroundMaterial: "acrylic",
    transparent: true,
    frame: false,
    webPreferences: {
      preload: path.join(__dirname, "preload_input.js"),
      additionalArguments: [`--trackingid=${trackingid}`],
    },
  });
  let cancelled = false;
  win.on("close", () => {
    if (!cancelled) {
      sendInputCancel(trackingid);
      cancelled = true;
      win.close();
    }
  });
  ipcMain.once(`${trackingid}-ready`, () => {
    win.webContents.send("data", {
      inputtype: inputtype,
      inputplaceholder: placeholder,
      inputmaxlength: maxlength,
      trackingid: trackingid,
      type: "input",
    });
    ipcMain.once(`${trackingid}-inputsuccess`, (event, data) => {
      sendInputResult(trackingid, data);
      if (!cancelled) {
        cancelled = true;
        win.close();
      }
    });
    ipcMain.once(`${trackingid}-inputcancel`, () => {
      if (!cancelled) {
        sendInputCancel(trackingid);
        cancelled = true;
        win.close();
      }
    });
  });
  load(win);
};

async function connectWebSocket() {
  try {
    console.log("Connecting to WebSocket");
    socket = new WebSocket(`ws://localhost:${port()}/popupipc`);
    let dontopen = false;
    socket.on("error", (err) => {
      dontopen = true;
      setTimeout(async () => {
        connectWebSocket();
      }, 3000);
    });
    socket.on("open", () => {
      if (dontopen) return;
      console.log("Connected to WebSocket");
      setupWebSocket();
    });
    return;
  } catch (e) {
    console.log("Error connecting to WebSocket");
    console.error(e);
    setTimeout(async () => {
      connectWebSocket();
    }, 1500);
  }
}

function setupWebSocket() {
  socket.on("close", async () => {
    console.log("Socket closed");
    connectWebSocket();
  });
  socket.on("open", () => {
    console.log("Socket connected");
  });
  socket.on("error", (err) => {
    console.error(err);
  });
  socket.on("message", (data) => {
    let message = protocol.PopupAppControlMessage.decode(data);
    switch (message.type) {
      case protocol.PopupAppControlMessage.MessageType.POPUP:
        createPopup(
          message.popupRequest.URL,
          message.popupRequest.ClientID,
          message.popupRequest.AppName,
          message.popupRequest.Favicon,
          message.popupRequest.Title,
          message.popupRequest.requestID
        );
        break;
      case protocol.PopupAppControlMessage.MessageType.INPUT:
        createInputWindow(
          message.inputRequest.InputType,
          message.inputRequest.InputPlaceholder,
          message.inputRequest.MaxLength,
          message.inputRequest.requestID
        );
        break;
      case protocol.PopupAppControlMessage.MessageType.SHUTDOWN:
        app.quit();
        break;
    }
  });
}

app.whenReady().then(() => {
  //createWindow();
  createHiddenWindow();
  session.defaultSession.webRequest.onHeadersReceived((details, callback) => {
    callback({
      responseHeaders: {
        ...details.responseHeaders,
        "Content-Security-Policy": ["script-src 'self'"],
      },
    });
  });
});
