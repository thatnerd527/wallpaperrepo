import { app, BrowserWindow, ipcMain, session } from "electron";
import { set } from "lodash-es";
import path from "node:path";
import WebSocket, { WebSocketServer } from "ws";
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
    JSON.stringify({
      Type: "popup",
      trackingID: trackingid,
      popup_ResultData: data || "",
      cancelled: false,
    })
  );
};

const sendPopupCancel = (trackingid) => {
  socket.send(
    JSON.stringify({
      Type: "popup",
      trackingID: trackingid,
      cancelled: true,
      popup_ResultData: "",
    })
  );
};

const sendInputResult = (trackingid, data) => {
  socket.send(
    JSON.stringify({
      Type: "input",
      trackingID: trackingid,
      input_ResultData: data,
      cancelled: false,
    })
  );
};

const sendInputCancel = (trackingid) => {
  socket.send(
    JSON.stringify({
      Type: "input",
      trackingID: trackingid,
      input_ResultData: "",
      cancelled: true,
    })
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
    let message = JSON.parse(data);
    console.log(message);
    switch (message.Type) {
      case "popup":
        createPopup(
          message["popup_URL"],
          message["popup_ClientID"],
          message["popup_AppName"],
          message["popup_Favicon"],
          message["popup_Title"],
          message["trackingID"]
        );
        break;
      case "input":
        createInputWindow(
          message["input_Type"],
          message["input_Placeholder"],
          message["input_MaxLength"],
          message["trackingID"]
        );
        break;
      case "stop":
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
