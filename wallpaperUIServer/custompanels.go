package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"wallpaperuiserver/protocol"

	"github.com/elliotchance/pie/v2"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var ADDONS_PATH = "addons"
var RuntimeAddons = []*RuntimeAddon{}
/*
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
} */

type TemplateCustomPanel2 struct {
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
	NewData []protocol.RuntimePanel
	ClientID string
	Stop bool
}



type BackgroundDataChangeRequest struct {
	NewBackgrounds []protocol.RuntimeBackground
	NewActiveBackground string
	ClientID string
	Stop bool
}

type BackgroundDataChangeRequest2 struct {
	NewBackgrounds []protocol.RuntimeBackground
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

type CustomUpgrader struct {
	websocket.Upgrader
}

func PreparePanels(panels []protocol.RuntimePanel) []protocol.RuntimePanel {
	return pie.Map(panels, func(panel protocol.RuntimePanel) protocol.RuntimePanel {
		switch panel.BasePanel.PanelType {
		case "System":
			data, err := os.ReadFile(path.Join(ADDONS_PATH, panel.BasePanel.AppClientID, panel.BasePanel.PanelContent));
			//fmt.Println(string(data));
			if err != nil {
				log.Println(err)
				//fmt.Println("Error loading panel content")
				panel.BasePanel.PanelContent = fmt.Sprintf("Error loading panel content troubleshooting info: %v., Addon Client ID: %v, Panel File Name: %v", err, panel.BasePanel.AppClientID,panel.BasePanel.PanelContent)
				return panel
			}
			panel.BasePanel.PanelContent = string(data)
			return panel
		case "Embedded":
			url, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", panel.ControlPort))
			if err != nil {
				log.Println(err)
				return panel
			}
			added := url.Query()
			added.Add("content", panel.BasePanel.PanelContent)
			added.Add("clientid", panel.BasePanel.AppClientID)
			added.Add("loaderpanelid", panel.BasePanel.FixedPanelID)
			added.Add("persistentpanelid", panel.UniquePanelID)
			url.RawQuery = added.Encode()
			panel.BasePanel.PanelContent = fmt.Sprintf("%v", url)
			return panel
		case "Video", "Image", "Audio":
			return panel
		default:
			log.Println("Unknown panel type: " + panel.BasePanel.PanelType)
			return panel
		}
	})
}

func UnstripPanels(templatepanels []protocol.BasePanel, panels []protocol.RuntimePanel) []protocol.RuntimePanel {
	return pie.Map(panels, func(panel protocol.RuntimePanel) protocol.RuntimePanel {
		switch panel.BasePanel.PanelType {
		case "System":
			fromtemplate := templatepanels[pie.FindFirstUsing(templatepanels, func(tp protocol.BasePanel) bool { return tp.FixedPanelID == panel.BasePanel.FixedPanelID })]
			panel.BasePanel.PanelContent = fromtemplate.PanelContent;
			return panel
		case "Embedded":
			fromtemplate := templatepanels[pie.FindFirstUsing(templatepanels, func(tp protocol.BasePanel) bool { return tp.FixedPanelID == panel.BasePanel.FixedPanelID })]
			panel.BasePanel.PanelContent = fromtemplate.PanelContent;
			return panel
		case "Video", "Image", "Audio":
			return panel
		default:
			log.Println("Unknown panel type: " + panel.BasePanel.PanelType)
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

	marshalled, err := proto.Marshal(&result)
	clientid := GenerateGUID()

	if err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, []byte("ERAD: Error marshalling result"))
		return
	}
	err = c.WriteMessage(websocket.BinaryMessage, marshalled)
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
			cr := protocol.RuntimePanelChangeRequest{}
			cr.PanelChanges = pie.Map(change.NewData, func(panel protocol.RuntimePanel) *protocol.RuntimePanel { return &panel })
			marshalled, err := proto.Marshal(&cr)
			if err != nil {
				log.Println("json:", err)
				return
			}
			c.WriteMessage(websocket.BinaryMessage, marshalled)
		}
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			panelchangehub.SendMessage(PanelDataChangeRequest{NewData: nil, ClientID: clientid, Stop: true})
			return
		}
		parsedLastActive := protocol.RuntimePanelChangeRequest{}
		proto.Unmarshal(message, &parsedLastActive)

		panelchangehub.SendMessage(PanelDataChangeRequest{
			NewData: pie.Map(parsedLastActive.PanelChanges, func(panel *protocol.RuntimePanel) protocol.RuntimePanel { return *panel }),
			ClientID: clientid,
			Stop: false,
		})
		if err != nil {
			log.Println("json:", err)
			return
		}
		persistent := protocol.PersistentRuntimePanels{}

		persistent.Panels = pie.Map(UnstripPanels(
			pie.Map(result.TemplatePanels,func(a *protocol.BasePanel) protocol.BasePanel {return *a}),
			pie.Map(parsedLastActive.PanelChanges, func(panel *protocol.RuntimePanel) protocol.RuntimePanel { return *panel }),
		), func(panel protocol.RuntimePanel) *protocol.RuntimePanel { return &panel })
		parsedLastActive.PanelChanges = persistent.Panels;
		saved, err := protojson.Marshal(&persistent)
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

func LoadAllPanelsFromAddons() (protocol.PanelSystemInitialResponse, error) {
	if _, err := os.Stat("addons"); err != nil {
		os.Mkdir("addons", 0755)
	}
	addons, err := os.ReadDir(ADDONS_PATH)
	if err != nil {
		log.Println(err)
		return protocol.PanelSystemInitialResponse{}, fmt.Errorf("Error reading addons directory: %v", err)
	}
	availablePanels := make([]protocol.BasePanel, 0)
	convertedAvailablePanels := make([]protocol.RuntimePanel, 0)
	if _, err := os.Stat("storage/panels.json"); err != nil {
		file, err := os.Create("storage/panels.json")
		file.Write([]byte("{}"))
		if err != nil {
			log.Println(err)
			return protocol.PanelSystemInitialResponse{}, fmt.Errorf("Error creating panels.json: %v", err)
		}
		file.Close()
	}
	lastActivePanels, err := LoadAllRuntimePanels("storage/panels.json")
	preparedLastActivePanels := make([]protocol.RuntimePanel, 0)
	if err != nil {
		log.Println(err)
		return protocol.PanelSystemInitialResponse{}, fmt.Errorf("Error loading last active panels: %v", err)
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
			return protocol.PanelSystemInitialResponse{}, fmt.Errorf("Error loading manifest for addon %s: %v", addon.Name(), err)
		}
		declaredPanels, err := LoadAllBasePanels(ADDONS_PATH + "/" + addon.Name() + "/panels.json")

		//Load all media as panels
		if _ , err := os.Stat(ADDONS_PATH + "/" + addon.Name() + "/media/panels"); err == nil {
			mediaPanels, err := LoadAllMediaAsPanels(ADDONS_PATH+"/"+addon.Name()+"/media/panels", manifest)
			if err != nil {
				log.Println(err)
				return protocol.PanelSystemInitialResponse{}, fmt.Errorf("Error loading media panels for addon %s: %v", addon.Name(), err)
			}
			declaredPanels = append(declaredPanels, mediaPanels...)
		}

		declaredPanels = pie.Map(declaredPanels, func(panel protocol.BasePanel) protocol.BasePanel {
			panel.AppClientID = manifest.ClientID
			return panel
		})
		if err != nil {
			log.Println(err)
			return protocol.PanelSystemInitialResponse{}, fmt.Errorf("Error loading panels for addon %s: %v", addon.Name(), err)
		}
		availablePanels = append(availablePanels, declaredPanels...)
		converted := ConvertTemplatePanelsToRuntimePanels(availablePanels)

		lastActivePanels = MeshTemplateAndRuntimePanels(declaredPanels, lastActivePanels)
		ra, err := manifest.GetRAOrBootstrapAddon(addon.Name())
		panelsForThisAddon := pie.Filter(lastActivePanels, func(ra protocol.RuntimePanel) bool { return ra.BasePanel.AppClientID == manifest.ClientID })
		templateRuntimePanelsForThisAddon2 := pie.Filter(converted, func(ra protocol.RuntimePanel) bool { return ra.BasePanel.AppClientID  == manifest.ClientID })
		if err != nil {
			if err.Error() == "No bootstrap executable found for addon "+manifest.Name {
				log.Println("No bootstrap executable found for addon " + manifest.Name)
			} else {
				log.Println(err)
			}

		} else {

			for index, _ := range panelsForThisAddon {
				panelsForThisAddon[index].ControlPort = int32(ra.ControlPort)
			}
			for index, _ := range templateRuntimePanelsForThisAddon2 {
				templateRuntimePanelsForThisAddon2[index].ControlPort = int32(ra.ControlPort)
			}

		}
		convertedAvailablePanels = append(convertedAvailablePanels, templateRuntimePanelsForThisAddon2...)
		preparedLastActivePanels = append(preparedLastActivePanels, panelsForThisAddon...)

	}
	deleted := DeletedPanels(preparedLastActivePanels, availablePanels)
	for _, panel := range deleted {
		for index, ra := range preparedLastActivePanels {
			if panel.BasePanel.FixedPanelID == ra.BasePanel.FixedPanelID {
				preparedLastActivePanels[index].Deleted = true
			}
		}
	}

	finallastActivePanels := PreparePanels(preparedLastActivePanels)
	result := protocol.PanelSystemInitialResponse{

		InstancedButBlankPanels:         pie.Map(PreparePanels(convertedAvailablePanels), func(panel protocol.RuntimePanel) *protocol.RuntimePanel { return &panel }),
		TemplatePanels: pie.Map(availablePanels, func(panel protocol.BasePanel) *protocol.BasePanel { return &panel }),
		InstancedPanelsFromStorage:        pie.Map(finallastActivePanels, func(panel protocol.RuntimePanel) *protocol.RuntimePanel { return &panel }),
		DeletedInstancedPanels:           pie.Map(deleted, func(panel protocol.RuntimePanel) *protocol.RuntimePanel { return &panel }),
	}
	return result, nil
}
