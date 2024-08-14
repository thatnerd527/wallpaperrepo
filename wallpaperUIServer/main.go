package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	"wallpaperuiserver/protocol"

	"github.com/getlantern/systray"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/orcaman/concurrent-map/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var upgrader = websocket.Upgrader{}
var storagechanged = make(map[string]MessageHub[StorageChange])
var popupchannel = make(chan PopupRequest)
var popupprocess *os.Process = nil
var popupresultchannel = CreateMessageHub[PopupResult]()
var awaitingpopup = cmap.New[PopupRequest]()
var secureport = 8080
var allowedredirectport = 0
var ELECTRON_PATH = "popupapp"
var embedKey = "factory"

type StorageChange struct {
	Scope string
	Data  string
	ClientID string
}


func CreateStorageChange(scope string, data string, clientid string) StorageChange {
	return StorageChange{Scope: scope, Data: data, ClientID: clientid}
}

func GenerateGUID() string {
    id := uuid.New()
	return fmt.Sprint(id.String())
}

func startPopupApp() {
	resolved , err := filepath.Abs(path.Join(ELECTRON_PATH,"wallpaperpopupapp.exe"))
	if err != nil {
		log.Println(err)
	}
	addonProcess, err := os.StartProcess(resolved, []string{resolved, fmt.Sprintf("--controlport=%d",secureport)}, &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir: path.Join(ELECTRON_PATH),
	})
	if err != nil {
		log.Println(err)
	}
	popupprocess = addonProcess
}


func main() {
	fmt.Println("Starting SystemHub")
	if _, err := os.Stat(".pid"); err == nil {
		data, err := os.ReadFile(".pid")
		if err != nil {
			log.Println(err)
		}
		pid, err := strconv.Atoi(string(data))
		if err != nil {
			log.Println(err)
		}
		proc, err := os.FindProcess(pid)

		if proc != nil && err == nil {
			err = proc.Signal(syscall.Signal(0))
			fmt.Println(err)
		}
		if err == nil {
			fmt.Println("Already running")
			fmt.Println("Found previous instance with pid", pid, "name: ",proc.Kill())
			os.Exit(0)
			return
		} else {
			fmt.Println("Previous instance may have crashed, resetting pid file")
		}
	}
	os.WriteFile(".pid", []byte(fmt.Sprint(os.Getpid())), 0644)
	if _, err := os.Stat("embedkey"); os.IsNotExist(err) {
		os.WriteFile("embedkey", []byte(embedKey), 0644)
	} else {
		data, err := os.ReadFile("embedkey")
		if err != nil {
			log.Println(err)
		}
		embedKey = string(data)
	}

	unused := FindUnusedPort()
	secureport = unused
	secureMux := http.NewServeMux()
	if _, err := os.Stat("storage"); os.IsNotExist(err) {
		os.Mkdir("storage", 0755)
	}

	cachedpreferences.AddWriteHandler(func(prefs protocol.AppSettings) protocol.AppSettings {
		marshalled, _ := protojson.Marshal(&prefs)
		os.WriteFile("preferences.json", marshalled, 0755)
		preferenceschannel.SendMessage(PreferenceUpdate{prefs, GenerateGUID(), false})
		return prefs
	})


	secureMux.HandleFunc("/storagesocket", storageHub)
	secureMux.HandleFunc("/popupipc",popupIPC)
	secureMux.HandleFunc("/inputrequest",inputRequest)
	secureMux.HandleFunc("/popuprequest",popupRequest)
	secureMux.HandleFunc("/panelsystem", panelSystem)
	secureMux.HandleFunc("/backgroundsystem", backgroundSystem)
	secureMux.HandleFunc("/addons",getActiveAddons)
	secureMux.HandleFunc("/disabledaddons", getDisabledAddons)
	secureMux.HandleFunc("/getaddonorigin",getAddonOrigin)
	secureMux.HandleFunc("/mediaregistry",mediaRegistry)
	secureMux.HandleFunc("/preferences",preferencesSystem)
	secureMux.HandleFunc("/restart",restartHandler)
	secureMux.HandleFunc("/shutdown",shutdownHandler)
	secureMux.HandleFunc("/restartipc",restartIPC)
	secureMux.HandleFunc("/installaddon",installAddon)
	secureMux.HandleFunc("/simplebackground",simpleBackgroundHandler)
	secureMux.HandleFunc("/getencodingstatus",getEncodingStatus)
	secureMux.HandleFunc("/autostart",autoStartHandler)
	secureMux.HandleFunc("/addonchanges",applyAddonChanges)
	secureMux.HandleFunc("/getpreviewfile",previewFileHandler)
	secureMux.HandleFunc("/setbackgroundfromcache",setBackgroundFromCache)

	defaultMux := http.NewServeMux()
	defaultMux.HandleFunc("/", fileHandler)
	defaultMux.HandleFunc("/generate_200", generate_200)

	redirectMux := http.NewServeMux()
	redirectMux.HandleFunc("/redirect", redirectHandler)

	go func() {
		fmt.Println("Listening on port 8080")
		http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v",8080), defaultMux)

	}()

	go func() {
		fmt.Println("Listening on port 8081")
		http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v",8081), redirectMux)

	}()
	startPopupApp()
	fmt.Println("Listening on port", unused)
	go func() {
		systray.Run(createTray, func() {})
	}()
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v",unused), secureMux)
	shutdown()
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	shutdown()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Shutting down"))
}

func actuallyExitFr() {
	go func() {
		os.Exit(0)
	}()

}

func shutdown() {
	popupchannel <- PopupRequest{Type: "stop"}
	restarthub.SendMessage(RestartRequest{
		Stop: false,
		ClientID: uuid.NewString(),
		Message2: "exit",
	})
	if popupprocess != nil {
		popupprocess.Kill()
	}
	for _, addon := range RuntimeAddons {
		addon.Shutdown()
	}
	os.Remove(".pid")
	time.Sleep(5 * time.Second)
	actuallyExitFr()

}

func mediaRegistry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if path, ok := guidtofilepath[r.URL.Query().Get("guid")]; ok {
		http.ServeFile(w, r, path)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if (!r.URL.Query().Has("embedkey") || r.URL.Query().Get("embedkey") != embedKey) {
		fmt.Println("EXPECTED: ", embedKey + " GOT: ", r.URL.Query().Get("embedkey"))
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
		w.Header().Set("Cache-Control", "no-store")
	}
		fmt.Println("Redirecting to ", r.URL.Query().Get("url"))
		url2 := r.URL.Query().Get("url")
		newurl, err := url.Parse(url2)
		if err != nil {
			fmt.Println(err)
		}

		// Additional maybe not needed security check, but just in case to prevent redirecting to other ports.
		// Not that escaping an Iframe is even possible, but just in case.
		if (allowedredirectport != 0 && newurl.Port() != strconv.Itoa(allowedredirectport)){
			w.WriteHeader(http.StatusForbidden)
			return
		}
		allowedredirectport, _ = strconv.Atoi(newurl.Port())
		query := newurl.Query()
		query.Add("controlPort", strconv.Itoa(secureport))
		newurl.RawQuery = query.Encode()
		http.Redirect(w, r, newurl.String(), http.StatusFound)


}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if (!r.URL.Query().Has("embedkey") || r.URL.Query().Get("embedkey") != embedKey) {
		fmt.Println("EXPECTED: " + embedKey + " GOT: " + r.URL.Query().Get("embedkey"))
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
	}
	w.Header().Set("Cache-Control", "no-store")
	http.ServeFile(w, r, path.Join("public", "wwwroot", r.URL.Path))


}


func popupRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	trackingID := r.URL.Query().Get("trackingID")
	if _, ok := awaitingpopup.Get(trackingID); ok {
		for {
			msg := popupresultchannel.WaitForMessage()
			if msg.trackingID == trackingID {
				if msg.cancelled {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.Write([]byte(msg.popup_ResultData))
				}
				// Delete from awaitingpopup
				awaitingpopup.Remove(trackingID)
				return
				break
			}
		}
	}
	request := CreatePopupRequest(r.URL.Query().Get("popup_URL"), r.URL.Query().Get("popup_ClientID"), r.URL.Query().Get("popup_AppName"), r.URL.Query().Get("popup_Favicon"), r.URL.Query().Get("popup_Title"));
	request.trackingID = trackingID
	awaitingpopup.Set(trackingID, request)
	popupchannel <- request
	for {
		msg := popupresultchannel.WaitForMessage()
		if msg.trackingID == trackingID {
			if msg.cancelled {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.Write([]byte(msg.popup_ResultData))
				}
			// Delete from awaitingpopup
			awaitingpopup.Remove(trackingID)
			break
		}
	}
}

func inputRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	MaxLength, _ := strconv.Atoi(r.URL.Query().Get("input_MaxLength"))
	trackingID := r.URL.Query().Get("trackingID")
	fmt.Println("Sent request2")
	if _, ok := awaitingpopup.Get(trackingID); ok {
		fmt.Println("Sent request3")
		for {
			msg := popupresultchannel.WaitForMessage()
			if msg.trackingID == trackingID {
				if msg.cancelled {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.Write([]byte(msg.input_ResultData))
				}
				// Delete from awaitingpopup
				awaitingpopup.Remove(trackingID)
				return
				break
			}
		}
	}
	request := CreateTypingRequest(r.URL.Query().Get("input_Type"), r.URL.Query().Get("input_Placeholder"), MaxLength);
	request.trackingID = trackingID
	awaitingpopup.Set(trackingID, request)
	fmt.Println("Sent request")
	popupchannel <- request

	for {
		msg := popupresultchannel.WaitForMessage()
		if msg.trackingID == trackingID {
			if msg.cancelled {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.Write([]byte(msg.input_ResultData))
				}
			// Delete from awaitingpopup
			awaitingpopup.Remove(trackingID)
			break
		}
	}
}


func generate_200(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Hello, World!"))
}

func doesScopeExist(scope string) bool {
	stat, err := os.Stat(path.Join("storage", scope))
	return err == nil && stat.Mode().IsRegular()
}

