package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"path/filepath"

	"github.com/elliotchance/pie/v2"
	"github.com/gorilla/websocket"
	"github.com/sqweek/dialog"
)

var MAX_STARTUP_ATTEMPTS = 5
var ALLOW_RESTART = true

func LoadManifest(filename string) (AddonManifest, error) {
	manifestFile, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return AddonManifest{}, err
	}
	manifest := AddonManifest{}
	err = json.Unmarshal(manifestFile, &manifest)
	if err != nil {
		log.Println(err)
		return AddonManifest{}, err
	}
	return manifest, nil
}

func FindUnusedPort() int {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func StartContinousHealthMonitor(RuntimeAddon *RuntimeAddon, handlecrash func()) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%d/keepalive", RuntimeAddon.ControlPort), nil)
	if err != nil {
		log.Println("Failed to connect to keepalive endpoint")
		log.Println(err)
		return
	}
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read keepalive message")
			log.Println("This is likely due to the addon crashing/being killed")
			log.Println(err)
			RuntimeAddon.Disconnected = true
			RuntimeAddon.Socket = nil
			RuntimeAddon.Process = nil
			handlecrash()
			return
		}
	}
}

func TryToConnectToAddon(ra *RuntimeAddon, donecallback func(error)) {
	attempts := 0
	for {
		if attempts > MAX_STARTUP_ATTEMPTS {
			ra.Disconnected = true
			ra.Socket = nil
			ra.Process = nil
			donecallback(fmt.Errorf("Failed to connect to addon"))
			return
		}
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%d/ipc", ra.ControlPort), nil)
		if err != nil {
			attempts++
			log.Println(err)
			time.Sleep(1 * time.Second)
		} else {
			ra.Socket = conn
			ra.Disconnected = false
			donecallback(nil)
			return
		}
	}

}

func BootstrapAddon(foldername string, manifest AddonManifest, donecallback func(error), crashhandler func(error), allowrestarts bool) (*RuntimeAddon, error) {
	// Start the addon process
	runtimeaddon := RuntimeAddon{Manifest: manifest, Disconnected: true}
	ref := &runtimeaddon
	unusedPort := FindUnusedPort()
	runtimeaddon.ControlPort = unusedPort
	resolved , err := filepath.Abs(path.Join(ADDONS_PATH,foldername,manifest.GetBootstrapExecutable()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	addonProcess, err := os.StartProcess(resolved, []string{resolved, fmt.Sprintf("--port=%d",unusedPort)}, &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir: path.Join(ADDONS_PATH,foldername),
	})
	runtimeaddon.Process = addonProcess
	if err != nil {
		log.Println(err)
		return nil, err
	}
	RuntimeAddons = append(RuntimeAddons, ref)
	go TryToConnectToAddon(ref, func(err error) {
		if err != nil {
			log.Default().Printf("Error during addon bootstrap of %s: %s", manifest.Name, err.Error())
			donecallback(err)
		} else {

			fmt.Println("Addon started successfully")
			donecallback(nil)
			go StartContinousHealthMonitor(ref, func() {
				log.Default().Printf("Addon %s has crashed", manifest.Name)
				RuntimeAddons = pie.Filter(RuntimeAddons, func(ra *RuntimeAddon) bool { return ra.Manifest.ClientID != manifest.ClientID })
				crashhandler(fmt.Errorf("Addon %s has crashed", manifest.Name))
				if ALLOW_RESTART && allowrestarts {
					log.Default().Printf("Restarting addon %s", manifest.Name)
					BootstrapAddon(foldername, manifest, func(err error) {
						if err != nil {
							log.Println("Addon failed to restart / will not restart.")
							log.Println(err)
							donecallback(err)
							return
						}
						log.Println("Addon restarted successfully.")
						donecallback(nil)
					}, crashhandler,allowrestarts)
				} else {
					log.Println("Restarting of addons is disabled.")
				}
			});
		}
	})

	return ref, nil
}

func LoadAllTemplatePanels(filename string) ([]TemplateCustomPanel, error) {
	panelsFile, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return []TemplateCustomPanel{}, err
	}
	panels := []TemplateCustomPanel{}
	err = json.Unmarshal(panelsFile, &panels)
	if err != nil {
		log.Println(err)
		return []TemplateCustomPanel{}, err
	}
	return panels, nil
}

func LoadAllRuntimePanels(filename string) ([]RuntimeCustomPanel, error) {
	panelsFile, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return []RuntimeCustomPanel{}, err
	}
	panels := []RuntimeCustomPanel{}
	err = json.Unmarshal(panelsFile, &panels)
	if err != nil {
		log.Println(err)
		return []RuntimeCustomPanel{}, err
	}
	return panels, nil
}

func ConvertTemplatePanelsToRuntimePanels(panels []TemplateCustomPanel) []RuntimeCustomPanel {
	runtimePanels := []RuntimeCustomPanel{}
	for _, panel := range panels {
		runtimePanels = append(runtimePanels, RuntimeCustomPanel{
			PanelType: panel.PanelType,
			LoaderPanelID: panel.LoaderPanelID,
			PersistentPanelID: GenerateGUID(),
			PersistentPanelData: panel.PanelDefaultData,
			PanelTitle: panel.PanelTitle,
			PanelContent: panel.PanelContent,
			PanelRecommendedWidth: panel.PanelRecommendedWidth,
			PanelRecommendedHeight: panel.PanelRecommendedHeight,
			PanelMinWidth: panel.PanelMinWidth,
			PanelMinHeight: panel.PanelMinHeight,
			PanelMaxWidth: panel.PanelMaxWidth,
			PanelMaxHeight: panel.PanelMaxHeight,
			PanelRecommendedX: panel.PanelRecommendedX,
			PanelRecommendedY: panel.PanelRecommendedY,
			PersistentPanelWidth: panel.PanelRecommendedWidth,
			PersistentPanelHeight: panel.PanelRecommendedHeight,
			PersistentPanelX: panel.PanelRecommendedX,
			PersistentPanelY: panel.PanelRecommendedY,
			PanelIcon: panel.PanelIcon,
			ClientID: panel.ClientID,
			Deleted: false,
		})
	}
	return runtimePanels
}

func MeshTemplateAndRuntimePanels(availablePanels []TemplateCustomPanel, lastActivePanels []RuntimeCustomPanel) []RuntimeCustomPanel {
	meshedPanels := append([]RuntimeCustomPanel{}, lastActivePanels...)
	for _, templatepanel := range availablePanels {
		index := pie.FindFirstUsing(meshedPanels, func(panel RuntimeCustomPanel) bool { return panel.LoaderPanelID == templatepanel.LoaderPanelID });
		if index == -1 {
			continue;
		}
		duplicatePanel := meshedPanels[index]
		duplicatePanel.PanelTitle = templatepanel.PanelTitle
		duplicatePanel.PanelContent = templatepanel.PanelContent
		duplicatePanel.PanelRecommendedWidth = templatepanel.PanelRecommendedWidth
		duplicatePanel.PanelRecommendedHeight = templatepanel.PanelRecommendedHeight
		duplicatePanel.PanelMinWidth = templatepanel.PanelMinWidth
		duplicatePanel.PanelMinHeight = templatepanel.PanelMinHeight
		duplicatePanel.PanelMaxWidth = templatepanel.PanelMaxWidth
		duplicatePanel.PanelMaxHeight = templatepanel.PanelMaxHeight
		duplicatePanel.PanelRecommendedX = templatepanel.PanelRecommendedX
		duplicatePanel.PanelRecommendedY = templatepanel.PanelRecommendedY
		duplicatePanel.PanelIcon = templatepanel.PanelIcon
		duplicatePanel.ClientID = templatepanel.ClientID
		duplicatePanel.Deleted = false
		meshedPanels[index] = duplicatePanel
	}
	for index, runtimePanel := range meshedPanels {
		found := false
		for _, templatepanel := range availablePanels {
			if runtimePanel.LoaderPanelID == templatepanel.LoaderPanelID {
				found = true
				break
			}
		}
		if !found {
			meshedPanels[index].Deleted = true
		}
	}

	return meshedPanels
}

func DeletedPanels(lastActivePanels []RuntimeCustomPanel, availablePanels []TemplateCustomPanel) []RuntimeCustomPanel {
	deletedPanels := []RuntimeCustomPanel{}
	for _, runtimePanel := range lastActivePanels {
		found := false
		for _, templatePanel := range availablePanels {
			if runtimePanel.LoaderPanelID == templatePanel.LoaderPanelID {
				found = true
				break
			}
		}
		if !found {
			runtimePanel.Deleted = true
			deletedPanels = append(deletedPanels, runtimePanel)
		}
	}
	return deletedPanels
}

func getAddons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	addons := []AddonManifest{}
	for _, addon := range RuntimeAddons {
		addons = append(addons, addon.GetAddonManifest())
	}
	json.NewEncoder(w).Encode(addons)
}

func getAddonOrigin(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("clientID")
	for _, addon := range RuntimeAddons {
		if addon.Manifest.ClientID == clientID {
			w.Write([]byte(fmt.Sprintf("http://127.0.0.1:%d", addon.ControlPort)))
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)

}

func CopyFolder(source string, dest string) error {
	// Get all files from source
	files, err := os.ReadDir(source)
	if err != nil {
		return err
	}
	for _, file := range files {
		sourceFile := path.Join(source, file.Name())
		destFile := path.Join(dest, file.Name())
		if file.IsDir() {
			err = os.MkdirAll(destFile, os.ModePerm)
			if err != nil {
				return err
			}
			err = CopyFolder(sourceFile, destFile)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(sourceFile, destFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyFile(source string, dest string) error {
	sourceFile, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	err = os.WriteFile(dest, sourceFile, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func installAddon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	installmode := r.URL.Query().Get("installmode")
	folderPath := ""
	if installmode == "zip" {
		zipPath2, err := dialog.File().Title("Select addon Zip").Filter("Zip Files", "zip").Load()
		if err != nil {
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusConflict)

			return
		}
		zipPath, err := smartExtract(zipPath2)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		folderPath = zipPath
	} else if installmode == "folder" {
		folderPath2, err := dialog.Directory().Title("Select addon Folder").Browse()
		if err != nil {
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusConflict)
			return
    	}
		folderPath = folderPath2
	}


    fmt.Println("Selected Folder:", folderPath)
	manifest, err := LoadManifest(path.Join(folderPath, "manifest.json"))
	if (pie.Any(RuntimeAddons,func(ra *RuntimeAddon) bool { return ra.Manifest.ClientID == manifest.ClientID })) {
		addon := pie.Filter(RuntimeAddons, func(ra *RuntimeAddon) bool { return ra.Manifest.ClientID == manifest.ClientID })[0]
		addon.Shutdown()
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	copiedpath := path.Join(ADDONS_PATH, manifest.ClientID)
	copiedname := path.Base(copiedpath)
	err = CopyFolder(folderPath, path.Join(ADDONS_PATH, manifest.ClientID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	_, err = BootstrapAddon(copiedname, manifest, func(err error) {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}, func(err error) {
		w.WriteHeader(http.StatusInternalServerError)
	}, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}