package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/elliotchance/pie/v2"
	"github.com/gorilla/websocket"
)

var ADDONS_PATH = "addons"
var RuntimeAddons = []*RuntimeAddon{}

type ResultResponse struct {
	AvailablePanels         []RuntimeCustomPanel
	AvailableTemplatePanels []TemplateCustomPanel
	LastActivePanels        []RuntimeCustomPanel
	DeletedPanels           []RuntimeCustomPanel
}

type RuntimeCustomPanel struct {
	PanelType string

	// Unique per-panel on file
	LoaderPanelID string

	// Unique per-panel instance
	PersistentPanelID      string
	PanelTitle             string
	PanelContent           string
	PanelRecommendedWidth  int
	PanelRecommendedHeight int
	PanelMinWidth          int
	PanelMinHeight         int
	PanelMaxWidth          int
	PanelMaxHeight         int
	PersistentPanelWidth   int
	PersistentPanelHeight  int
	PanelRecommendedX      int
	PanelRecommendedY      int
	PersistentPanelX       int
	PersistentPanelY       int
	ControlPort            int
	Deleted                bool
	PersistentPanelData    string

	PanelIcon string
	ClientID  string
}

type TemplateCustomPanel struct {
	PanelType string

	// Unique per-panel on file
	LoaderPanelID          string
	PanelTitle             string
	PanelContent           string
	PanelRecommendedWidth  int
	PanelRecommendedHeight int
	PanelMinWidth          int
	PanelMinHeight         int
	PanelMaxWidth          int
	PanelMaxHeight         int
	PanelRecommendedX      int
	PanelRecommendedY      int
	PanelDefaultData       string

	PanelIcon string
	ClientID  string
}

type AddonManifest struct {
	Name               string `json:"name"`
	Version            string `json:"version"`
	Author             string `json:"author"`
	Description        string `json:"description"`
	ClientID           string `json:"clientID"`
	BootstapExecutable string `json:"bootstrapExecutable"`
	EnableRestart      bool   `json:"enableRestart"`
}

type RuntimeAddon struct {
	Manifest     AddonManifest
	Socket       *websocket.Conn
	Process      *os.Process
	Disconnected bool
	ControlPort  int
}

type PanelDataChangeRequest struct {
	NewData []RuntimeCustomPanel
	ClientID string
	Stop bool
}



type BackgroundDataChangeRequest struct {
	NewBackgrounds []RuntimeCustomBackground
	NewActiveBackground string
	ClientID string
	Stop bool
}

type BackgroundDataChangeRequest2 struct {
	NewBackgrounds []RuntimeCustomBackground
	NewActiveBackground string
}

var panelchangehub = CreateMessageHub[PanelDataChangeRequest]()


func (am *AddonManifest) GetBootstrapExecutable() string {
	return am.BootstapExecutable
}

func (ra *RuntimeAddon) GetAddonManifest() AddonManifest {
	return ra.Manifest
}

func (ra *RuntimeAddon) GetSocket() *websocket.Conn {
	return ra.Socket
}

func (ra *RuntimeAddon) Shutdown() {
	if ra.Manifest.IsBootstrapExecutable() {
		http.Get(fmt.Sprintf("http://127.0.0.1:%d/shutdown", ra.ControlPort))
	}
}

func (ra *RuntimeAddon) GetState() string {
	if ra.Manifest.IsBootstrapExecutable() {
		if ra.Disconnected {
			return "Disconnected"

		} else {
			return "Running"
		}
	} else {
		return "Running"
	}
}

func (am *AddonManifest) GetRAOrBootstrapAddon(foldername string) (*RuntimeAddon, error) {
	if pie.Any(RuntimeAddons, func(ra *RuntimeAddon) bool { return ra.Manifest.ClientID == am.ClientID }) {
		fmt.Println("Returning existing RA")
		return pie.Filter(RuntimeAddons, func(ra *RuntimeAddon) bool { return ra.Manifest.ClientID == am.ClientID })[0], nil
	}
	if am.BootstapExecutable == "" {
		return nil, fmt.Errorf("No bootstrap executable found for addon %s", am.Name)
	}
	return BootstrapAddon(foldername, *am, func(err error) {
		// TODO: Handle start failure
		if err != nil {
			log.Println("Addon failed to start: " + am.Name)
			log.Println(err)
		}
	}, func(err error) {
		if err != nil {
			log.Println("Addon crashed: " + am.Name)
			log.Println(err)
		}
	}, am.EnableRestart)
}

func (am *AddonManifest) IsBootstrapExecutable() bool {
	return am.BootstapExecutable != ""
}

func CreateTemplateCustomPanel(panelType string, loaderPanelID string, panelTitle string, panelContent string, panelRecommendedWidth int, panelRecommendedHeight int, panelMinWidth int, panelMinHeight int, panelMaxWidth int, panelMaxHeight int, panelRecommendedX int, panelRecommendedY int, panelIcon string, clientID string, defaultData string) TemplateCustomPanel {
	return TemplateCustomPanel{PanelType: panelType, LoaderPanelID: loaderPanelID, PanelTitle: panelTitle, PanelContent: panelContent, PanelRecommendedWidth: panelRecommendedWidth, PanelRecommendedHeight: panelRecommendedHeight, PanelMinWidth: panelMinWidth, PanelMinHeight: panelMinHeight, PanelMaxWidth: panelMaxWidth, PanelMaxHeight: panelMaxHeight, PanelRecommendedX: panelRecommendedX, PanelRecommendedY: panelRecommendedY, PanelIcon: panelIcon, ClientID: clientID, PanelDefaultData: defaultData}
}

type CustomUpgrader struct {
	websocket.Upgrader
}

func PreparePanels(panels []RuntimeCustomPanel) []RuntimeCustomPanel {
	return pie.Map(panels, func(panel RuntimeCustomPanel) RuntimeCustomPanel {
		switch panel.PanelType {
		case "System":
			data, err := os.ReadFile(path.Join(ADDONS_PATH, panel.ClientID, panel.PanelContent));
			//fmt.Println(string(data));
			if err != nil {
				log.Println(err)
				//fmt.Println("Error loading panel content")
				panel.PanelContent = fmt.Sprintf("Error loading panel content troubleshooting info: %v., Addon Client ID: %v, Panel File Name: %v", err, panel.ClientID,panel.PanelContent)
				return panel
			}
			panel.PanelContent = string(data)
			return panel
		case "Embedded":
			url, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", panel.ControlPort))
			if err != nil {
				log.Println(err)
				return panel
			}
			added := url.Query()
			added.Add("content", panel.PanelContent)
			added.Add("clientid", panel.ClientID)
			added.Add("loaderpanelid", panel.LoaderPanelID)
			added.Add("persistentpanelid", panel.PersistentPanelID)
			url.RawQuery = added.Encode()
			panel.PanelContent = fmt.Sprintf("%v", url)
			return panel
		case "Video", "Image", "Audio":
			return panel
		default:
			log.Println("Unknown panel type: " + panel.PanelType)
			return panel
		}
	})
}

func UnstripPanels(templatepanels []TemplateCustomPanel, panels []RuntimeCustomPanel) []RuntimeCustomPanel {
	return pie.Map(panels, func(panel RuntimeCustomPanel) RuntimeCustomPanel {
		switch panel.PanelType {
		case "System":
			fromtemplate := templatepanels[pie.FindFirstUsing(templatepanels, func(tp TemplateCustomPanel) bool { return tp.LoaderPanelID == panel.LoaderPanelID })]
			panel.PanelContent = fromtemplate.PanelContent;
			return panel
		case "Embedded":
			fromtemplate := templatepanels[pie.FindFirstUsing(templatepanels, func(tp TemplateCustomPanel) bool { return tp.LoaderPanelID == panel.LoaderPanelID })]
			panel.PanelContent = fromtemplate.PanelContent;
			return panel
		case "Video", "Image", "Audio":
			return panel
		default:
			log.Println("Unknown panel type: " + panel.PanelType)
			return panel
		}
	})
}

func panelSystem(w http.ResponseWriter, r *http.Request) {
	upgrader2 := CustomUpgrader{}
	upgrader2.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader2.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	fmt.Println("Panel system connected")
	// Upgrading last active panels with new panels
	result, error := LoadAllPanelsFromAddons()
	if error != nil {
		log.Println(error)
		c.WriteMessage(websocket.TextMessage, []byte("ERAD: Error loading panels"))
		return
	}

	marshalled, err := json.Marshal(result)
	clientid := GenerateGUID()

	if err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, []byte("ERAD: Error marshalling result"))
		return
	}
	err = c.WriteMessage(websocket.TextMessage, marshalled)
	if err != nil {
		log.Println("write:", err)
		return
	}
	file := "storage/panels.json"
	go func() {
		for {
			change := panelchangehub.WaitForMessage()
			if change.Stop && change.ClientID == clientid {
				break
			}
			if change.ClientID == clientid {
				continue
			}

			marshalled, err := json.Marshal(change.NewData)
			if err != nil {
				log.Println("json:", err)
				return
			}
			c.WriteMessage(websocket.TextMessage, marshalled)
		}
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			panelchangehub.SendMessage(PanelDataChangeRequest{NewData: nil, ClientID: clientid, Stop: true})
			return
		}
		decoded := string(message)
		parsedLastActive := make([]RuntimeCustomPanel, 0)

		err = json.Unmarshal([]byte(decoded), &parsedLastActive)
		panelchangehub.SendMessage(PanelDataChangeRequest{NewData: parsedLastActive, ClientID: clientid, Stop: false})
		if err != nil {
			log.Println("json:", err)
			return
		}
		parsedLastActive = UnstripPanels(result.AvailableTemplatePanels, parsedLastActive)
		saved, err := json.Marshal(parsedLastActive)
		if err != nil {
			log.Println("json:", err)
			return
		}
		err = os.WriteFile(file, saved, 0644)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}

}

func GetManifestsAndPaths() ([]AddonManifest, []string, error) {
	addons, err := os.ReadDir(ADDONS_PATH)
	if err != nil {
		log.Println(err)
		return []AddonManifest{}, []string{}, fmt.Errorf("Error reading addons directory: %v", err)
	}
	manifests := make([]AddonManifest, 0)
	paths := make([]string, 0)
	for i := 0; i < len(addons); i++ {
		addon := addons[i]
		manifest, err := LoadManifest(ADDONS_PATH + "/" + addon.Name() + "/manifest.json")
		if err != nil {
			log.Println(err)
			return []AddonManifest{}, []string{}, fmt.Errorf("Error loading manifest for addon %s: %v", addon.Name(), err)
		}
		manifests = append(manifests, manifest)
		paths = append(paths, ADDONS_PATH+"/"+addon.Name())
	}
	return manifests, paths, nil
}

func LoadAllPanelsFromAddons() (ResultResponse, error) {
	if _, err := os.Stat("addons"); err != nil {
		os.Mkdir("addons", 0755)
	}
	addons, err := os.ReadDir(ADDONS_PATH)
	if err != nil {
		log.Println(err)
		return ResultResponse{}, fmt.Errorf("Error reading addons directory: %v", err)
	}
	availablePanels := make([]TemplateCustomPanel, 0)
	convertedAvailablePanels := make([]RuntimeCustomPanel, 0)
	if _, err := os.Stat("storage/panels.json"); err != nil {
		file, err := os.Create("storage/panels.json")
		file.Write([]byte("[]"))
		if err != nil {
			log.Println(err)
			return ResultResponse{}, fmt.Errorf("Error creating panels.json: %v", err)
		}
		file.Close()
	}
	lastActivePanels, err := LoadAllRuntimePanels("storage/panels.json")
	preparedLastActivePanels := make([]RuntimeCustomPanel, 0)
	if err != nil {
		log.Println(err)
		return ResultResponse{}, fmt.Errorf("Error loading last active panels: %v", err)
	}
	for i := 0; i < len(addons); i++ {
		addon := addons[i]
		if _, err := os.Stat(ADDONS_PATH + "/" + addon.Name() + "/disabled"); err == nil {
				fmt.Println("Addon " + addon.Name() + " is disabled")
				continue
			}
			if _, err := os.Stat(ADDONS_PATH + "/" + addon.Name() + "/manifest.json"); err != nil {
				log.Println(err)
				fmt.Errorf("Error loading manifest for addon %s: %v", addon.Name(), err)
				continue;
			}

		manifest, err := LoadManifest(ADDONS_PATH + "/" + addon.Name() + "/manifest.json")

		if err != nil {
			log.Println(err)
			return ResultResponse{}, fmt.Errorf("Error loading manifest for addon %s: %v", addon.Name(), err)
		}
		declaredPanels, err := LoadAllTemplatePanels(ADDONS_PATH + "/" + addon.Name() + "/panels.json")

		//Load all media as panels
		if _ , err := os.Stat(ADDONS_PATH + "/" + addon.Name() + "/media/panels"); err == nil {
			mediaPanels, err := LoadAllMediaAsPanels(ADDONS_PATH+"/"+addon.Name()+"/media/panels", manifest)
			if err != nil {
				log.Println(err)
				return ResultResponse{}, fmt.Errorf("Error loading media panels for addon %s: %v", addon.Name(), err)
			}
			declaredPanels = append(declaredPanels, mediaPanels...)
		}

		declaredPanels = pie.Map(declaredPanels, func(panel TemplateCustomPanel) TemplateCustomPanel {
			panel.ClientID = manifest.ClientID
			return panel
		})
		if err != nil {
			log.Println(err)
			return ResultResponse{}, fmt.Errorf("Error loading panels for addon %s: %v", addon.Name(), err)
		}
		availablePanels = append(availablePanels, declaredPanels...)
		converted := ConvertTemplatePanelsToRuntimePanels(availablePanels)

		lastActivePanels = MeshTemplateAndRuntimePanels(declaredPanels, lastActivePanels)
		ra, err := manifest.GetRAOrBootstrapAddon(addon.Name())
		panelsForThisAddon := pie.Filter(lastActivePanels, func(ra RuntimeCustomPanel) bool { return ra.ClientID == manifest.ClientID })
		templateRuntimePanelsForThisAddon2 := pie.Filter(converted, func(ra RuntimeCustomPanel) bool { return ra.ClientID == manifest.ClientID })
		if err != nil {
			if err.Error() == "No bootstrap executable found for addon "+manifest.Name {
				log.Println("No bootstrap executable found for addon " + manifest.Name)
			} else {
				log.Println(err)
			}

		} else {

			for index, _ := range panelsForThisAddon {
				panelsForThisAddon[index].ControlPort = ra.ControlPort
			}
			for index, _ := range templateRuntimePanelsForThisAddon2 {
				templateRuntimePanelsForThisAddon2[index].ControlPort = ra.ControlPort
			}

		}
		convertedAvailablePanels = append(convertedAvailablePanels, templateRuntimePanelsForThisAddon2...)
		preparedLastActivePanels = append(preparedLastActivePanels, panelsForThisAddon...)

	}
	deleted := DeletedPanels(preparedLastActivePanels, availablePanels)
	for _, panel := range deleted {
		for index, ra := range preparedLastActivePanels {
			if panel.LoaderPanelID == ra.LoaderPanelID {
				preparedLastActivePanels[index].Deleted = true
			}
		}
	}

	finallastActivePanels := PreparePanels(preparedLastActivePanels)
	result := ResultResponse{
		AvailablePanels:         PreparePanels(convertedAvailablePanels),
		AvailableTemplatePanels: availablePanels,
		LastActivePanels:        finallastActivePanels,
		DeletedPanels:           deleted,
	}
	return result, nil
}
