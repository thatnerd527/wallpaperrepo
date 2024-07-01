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

	"github.com/getlantern/systray"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/orcaman/concurrent-map/v2"
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
	unused := FindUnusedPort()
	secureport = unused
	secureMux := http.NewServeMux()
	if _, err := os.Stat("storage"); os.IsNotExist(err) {
		os.Mkdir("storage", 0755)
	}

	secureMux.HandleFunc("/storagesocket", storageHub)
	secureMux.HandleFunc("/popupipc",popupIPC)
	secureMux.HandleFunc("/inputrequest",inputRequest)
	secureMux.HandleFunc("/popuprequest",popupRequest)
	secureMux.HandleFunc("/panelsystem", panelSystem)
	secureMux.HandleFunc("/backgroundsystem", backgroundSystem)
	secureMux.HandleFunc("/addons",getAddons)
	secureMux.HandleFunc("/getaddonorigin",getAddonOrigin)
	secureMux.HandleFunc("/mediaregistry",mediaRegistry)
	secureMux.HandleFunc("/preferences",preferencesSystem)
	secureMux.HandleFunc("/restart",restartHandler)
	secureMux.HandleFunc("/restartipc",restartIPC)
	secureMux.HandleFunc("/installaddon",installAddon)

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

	w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")

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
	w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
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

