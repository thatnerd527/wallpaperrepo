let panelsystem = null
let backgroundsystem = null;
let restartipc = null;
let preferencessystem = null;
let controlPort = new URL(window.location.href).searchParams.get("controlPort");
let mode = new URL(window.location.href).searchParams.get("mode");

window.testanimation = function () {
    let contentdiv = document.getElementById("contentdiv");
    let expanderdiv = document.getElementById("expander");

    contentdiv.classList.remove("max-h-0")
    contentdiv.classList.add("max-h-max")
    contentdiv.style.maxHeight = "0px";
    console.log("Started.")
    requestAnimationFrame(() => {
        let height = contentdiv.getBoundingClientRect().height;
        contentdiv.classList.add("max-h-0");
        contentdiv.classList.remove("max-h-max");
        contentdiv.classList.add("opacity-0");

        requestAnimationFrame(() => {
            contentdiv.classList.remove("max-h-0");
            contentdiv.classList.remove("opacity-0");
            contentdiv.style.maxHeight = height + "px";
            contentdiv.classList.add("opacity-100");
        });
    });
}

console.log("Registered hooks")

window.mode = function () {
    let mode = new URL(window.location.href).searchParams.get("mode");
    if (mode == null) {
        return "wallpaper"
    }
    return mode
}

window.reloadui = function () {
    window.top.postMessage("reload", "*")
}

function setupPanelConnection() {
    return new Promise((resolve, reject) => {
        let connected = false;
        let initialConnectionPass = false;
        let sentInitialData = false;
        console.log("Setting up panel connection")

        panelsystem = new WebSocket(`ws://localhost:${controlPort}/panelsystem`);
        panelsystem.addEventListener("open", function () {
            connected = true;
            DotNet.invokeMethod("WallpaperUI", "SetConnectionStatus", true);
            initialConnectionPass = true;
            DotNet.invokeMethodAsync('WallpaperUI', 'PassPanelWebsocket', DotNet.createJSObjectReference(panelsystem));
        });
        panelsystem.addEventListener("close", function () {
            connected = false;
            DotNet.invokeMethod("WallpaperUI", "SetConnectionStatus", false);
            try {
                panelsystem.close();
                panelsystem = null;
            } catch (e) {
                console.log(e);
            }
            if (!initialConnectionPass) {
                return
            }
            setupPanelConnection();
        });
        panelsystem.addEventListener("message", function (event) {
            if (!connected) {
                return;
            }
            if (!sentInitialData) {
                DotNet.invokeMethod("WallpaperUI", "LoadPanelData", event.data);
                sentInitialData = true;
                resolve();
            } else {
                DotNet.invokeMethod("WallpaperUI", "UpdatePanelData", event.data);
            }
        });
        panelsystem.addEventListener("error", function (event) {
            console.log("Error");
            if (!initialConnectionPass) {
                console.log("Failed to connect")
                DotNet.invokeMethod("WallpaperUI", "ConnectionStartFailure", event.data);
            } else {
                setupPanelConnection();
            }
        });
    });
}

function setupBackgroundConnection() {
    return new Promise((resolve, reject) => {
        let connected = false;
        let initialConnectionPass = false;
        let sentInitialData = false;

        backgroundsystem = new WebSocket(`ws://localhost:${controlPort}/backgroundsystem`);
        backgroundsystem.addEventListener("open", function () {
            connected = true;
            DotNet.invokeMethod("WallpaperUI", "SetConnectionStatus", true);
            initialConnectionPass = true;
            DotNet.invokeMethodAsync('WallpaperUI', 'PassBackgroundWebsocket', DotNet.createJSObjectReference(backgroundsystem));
        });
        backgroundsystem.addEventListener("close", function () {
            connected = false;
            DotNet.invokeMethod("WallpaperUI", "SetConnectionStatus", false);
            try {
                backgroundsystem.close();
                backgroundsystem = null;
            } catch (e) {
                console.log(e);
            }
            if (!initialConnectionPass) {
                return
            }
            setupBackgroundConnection();
        });
        backgroundsystem.addEventListener("message", function (event) {
            if (!connected) {
                return;
            }
            if (!sentInitialData) {
                DotNet.invokeMethod("WallpaperUI", "LoadBackgroundData", event.data);
                sentInitialData = true;
                resolve();
            } else {
                DotNet.invokeMethod("WallpaperUI", "UpdateBackgroundData", event.data);
            }
        });
        backgroundsystem.addEventListener("error", function (event) {
            console.log("Error");
            if (!initialConnectionPass) {
                console.log("Failed to connect")
                DotNet.invokeMethod("WallpaperUI", "ConnectionStartFailure", event.data);
            } else {
                setupBackgroundConnection();
            }
        });
    });
}

function setupPreferencesConnection() {
    return new Promise((resolve, reject) => {
        let connected = false;
        let initialConnectionPass = false;
        let sentInitialData = false;

        preferencessystem = new WebSocket(`ws://localhost:${controlPort}/preferences`);
        preferencessystem.addEventListener("open", function () {
            connected = true;
            DotNet.invokeMethod("WallpaperUI", "SetConnectionStatus", true);
            initialConnectionPass = true;
            DotNet.invokeMethodAsync('WallpaperUI', 'PassPreferencesWebsocket', DotNet.createJSObjectReference(preferencessystem));
        });
        preferencessystem.addEventListener("close", function () {
            connected = false;
            DotNet.invokeMethod("WallpaperUI", "SetConnectionStatus", false);
            try {
                preferencessystem.close();
                preferencessystem = null;
            } catch (e) {
                console.log(e);
            }
            if (!initialConnectionPass) {
                return
            }
            setupPreferencesConnection();
        });
        preferencessystem.addEventListener("message", function (event) {
            if (!connected) {
                return;
            }
            DotNet.invokeMethod("WallpaperUI", "UpdatePreferences", event.data);
            resolve();
        });
        preferencessystem.addEventListener("error", function (event) {
            console.log("Error");
            if (!initialConnectionPass) {
                console.log("Failed to connect")
                DotNet.invokeMethod("WallpaperUI", "ConnectionStartFailure", event.data);
            } else {
                setupPreferencesConnection();
            }
        });
    });
}

function setupRestartConnection() {
    return new Promise((resolve, reject) => {
        let connected = false;
        let initialConnectionPass = false;
        let sentInitialData = false;

        restartipc = new WebSocket(`ws://localhost:${controlPort}/restartipc`);
        restartipc.addEventListener("open", function () {
            connected = true;
            DotNet.invokeMethod("WallpaperUI", "SetConnectionStatus", true);
            initialConnectionPass = true;
            DotNet.invokeMethodAsync('WallpaperUI', 'PassRestartWebsocket', DotNet.createJSObjectReference(restartipc));
            resolve();
        });
        restartipc.addEventListener("close", function () {
            connected = false;
            DotNet.invokeMethod("WallpaperUI", "SetConnectionStatus", false);
            try {
                restartipc.close();
                restartipc = null;
            } catch (e) {
                console.log(e);
            }
            if (!initialConnectionPass) {
                return
            }
            setupRestartConnection();
        });
        restartipc.addEventListener("message", function (event) {
            if (!connected) {
                return;
            }
            console.log("DECLARED MESSAGE " + event.data)
            if (event.data == "restart") {
                console.log("DECLARED MESSAGE 2 " + event.data)
                window.restart();
            } else if (event.data == "exit") {
                console.log("DECLARED MESSAGE 3 " + event.data)
                window.quit();
            }
            

        });
        restartipc.addEventListener("error", function (event) {
            console.log("Error");
            if (!initialConnectionPass) {
                console.log("Failed to connect")
                DotNet.invokeMethod("WallpaperUI", "ConnectionStartFailure", event.data);
            } else {
                setupRestartConnection();
            }
        });
    })
}

if (controlPort == null) {
    let toredirect = new URL("http://localhost:8081/redirect")
    toredirect.searchParams.append("url", window.location.href);
    window.location.href = toredirect.toString();
        
} else {
    var newurl = new URL(window.location.href)
    newurl.searchParams.delete("controlPort")
    newurl.searchParams.delete("embedkey")
    window.restart = async function () {
        let request = new URL(`http://127.0.0.1:${controlPort}/restart`)
        DotNet.invokeMethod("WallpaperUI", "SetRestart");
        try {
            await fetch(request.toString())
        } catch (e) {
            console.log(e)
        }
        window.top.postMessage("reload", "*")
    }

    window.quit = async function () {
        let request = new URL(`http://127.0.0.1:${controlPort}/shutdown`)
        DotNet.invokeMethod("WallpaperUI", "SetExit");
        try {
            await fetch(request.toString())
        } catch (e) {
            console.log(e)
        }
        window.top.postMessage("quit", "*")
    }

    window.remoterestart = function () {
        console.log("Restarting")
        window.wallpaperAPI.popupsuccess("");
        restartipc.send("restart");
    }

    window.remoteexit = function () {
        console.log("Exiting")
        window.wallpaperAPI.popupsuccess("");
        restartipc.send("exit");
    }

    window.history.replaceState({ path: newurl.toString() }, '', newurl.toString());
    window.dotnetready = function () {
        DotNet.invokeMethod("WallpaperUI", "PassControlPort", Number.parseInt(controlPort));
        setupBackgroundConnection().then(() => {
            setupPanelConnection().then(() => {
                setupPreferencesConnection().then(() => {
                    setupRestartConnection();
                });
            })
        })
        

        fetch("http://localhost:" + controlPort + "/addons").then(x => {
            x.text().then(y => {

                DotNet.invokeMethod("WallpaperUI", "UpdateAddonData",y);
            })
        });
        fetch("http://localhost:" + controlPort + "/disabledaddons").then(x => {
            x.text().then(y => {

                DotNet.invokeMethod("WallpaperUI", "UpdateDisabledAddonData", y);
            })
        });
        window.top.postMessage("ready", "*");
    }

    window.opensettings = function () {
        let url = new URL(window.location.href);
        url.searchParams.set("mode", "settings");

        let request = new URL(`http://localhost:${controlPort}/popuprequest`)
        request.searchParams.append("popup_URL", url.toString());
        request.searchParams.append("popup_ClientID", "system");
        request.searchParams.append("popup_AppName", "Wallpaper System");
        request.searchParams.append("popup_Favicon", "none");
        request.searchParams.append("popup_Title", "Wallpaper Settings");
        var guid = btoa(Math.random().toString()).replaceAll(/[^A-z0-9]/g, "")
        request.searchParams.append("trackingID", guid);
        fetch(request.toString())
    }

    window.refreshaddons = function () {
        fetch("http://localhost:" + controlPort + "/addons").then(x => {
            x.text().then(y => {

                DotNet.invokeMethod("WallpaperUI", "UpdateAddonData", y);
            })
        });
        fetch("http://localhost:" + controlPort + "/disabledaddons").then(x => {
            x.text().then(y => {

                DotNet.invokeMethod("WallpaperUI", "UpdateDisabledAddonData", y);
            })
        });
    }

    //window.savePanelData = function (data) {
    //    panelsystem.send(data);
    //}

    //window.saveBackgroundData = function (data) {
    //    backgroundsystem.send(data);
    //}
}
