package main

import (
	"encoding/json"
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

type TemplateCustomBackground2 struct {
	LoaderBackgroundID    string
	BackgroundType        string
	BackgroundContent     string
	BackgroundDefaultData string
	ClientID              string
}
/*
type protocol.RuntimeBackground struct {
	protocol.BaseBackground
	PersistentBackgroundID   string
	PersistentBackgroundData string
	Deleted                  bool
	ControlPort              int
} */

/* type ResultResponse2 struct {
	AvailableBackgrounds []protocol.RuntimeBackground
	AvailableTemplateBackgrounds []protocol.BaseBackground
	LastActiveBackgrounds []protocol.RuntimeBackground
	LastActiveBackgroundID string
	DeletedBackgrounds []protocol.RuntimeBackground
}

type SaveRequest struct {
	Backgrounds []protocol.RuntimeBackground
	LastActiveBackgroundID string
} */

var backgroundchangehub = CreateMessageHub[BackgroundDataChangeRequest]()

func LoadAllTemplateBackgrounds(filename string) ([]protocol.BaseBackground, error) {
	backgroundsFile, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return []protocol.BaseBackground{}, err
	}
	backgrounds := []TemplateCustomBackground2{}
	err = json.Unmarshal(backgroundsFile, &backgrounds)
	protoformat := []protocol.BaseBackground{}
	for _, background := range backgrounds {
		protoformat = append(protoformat, protocol.BaseBackground{
			BackgroundType: 	  background.BackgroundType,
			FixedBackgroundID:   background.LoaderBackgroundID,
			BackgroundContent:   background.BackgroundContent,
			BackgroundDefaultData: background.BackgroundDefaultData,
			AppClientID: 		background.ClientID,
		})
	}
	if err != nil {
		log.Println(err)
		return []protocol.BaseBackground{}, err
	}
	return protoformat, nil
}

func LoadAllRuntimeBackgrounds(filename string) ([]protocol.RuntimeBackground, error) {
	backgroundsFile, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return []protocol.RuntimeBackground{}, err
	}
	result := protocol.PersistentRuntimeBackgrounds{}
	err = protojson.Unmarshal(backgroundsFile, &result)
	copied := pie.Map(result.Backgrounds, func(background *protocol.RuntimeBackground) protocol.RuntimeBackground {
		return *proto.Clone(background).(*protocol.RuntimeBackground);
	})
	if err != nil {
		log.Println(err)
		return []protocol.RuntimeBackground{}, err
	}
	return copied, nil
}

func ConvertTemplateBackgroundsToRuntimeBackgrounds(backgrounds []protocol.BaseBackground) []protocol.RuntimeBackground {
	runtimeBackgrounds := []protocol.RuntimeBackground{}
	for _, background := range backgrounds {
		runtimeBackgrounds = append(runtimeBackgrounds, protocol.RuntimeBackground{
			BaseBackground: 		 proto.Clone(&background).(*protocol.BaseBackground),
			UniqueBackgroundID:   GenerateGUID(),
			PersistentData: background.BackgroundDefaultData,
			Deleted:                  false,
		})
	}
	return runtimeBackgrounds
}

func MeshBackgrounds(runtimeBackgrounds []protocol.RuntimeBackground, templateBackgrounds []protocol.BaseBackground) []protocol.RuntimeBackground {
	result := append([]protocol.RuntimeBackground{}, runtimeBackgrounds...)
	for _, templateBackground := range templateBackgrounds {
		index := pie.FindFirstUsing(runtimeBackgrounds, func(bg protocol.RuntimeBackground) bool {
			return bg.BaseBackground.FixedBackgroundID == templateBackground.FixedBackgroundID
		})
		if index == -1 {
			continue
		}
		result[index] = protocol.RuntimeBackground{
			BaseBackground: proto.Clone(&templateBackground).(*protocol.BaseBackground),
			UniqueBackgroundID:   runtimeBackgrounds[index].UniqueBackgroundID,
			PersistentData:  runtimeBackgrounds[index].PersistentData,
			Deleted:                  runtimeBackgrounds[index].Deleted,
		}
	}
	for index, runtimePanel := range runtimeBackgrounds {
		found := false
		for _, templatepanel := range templateBackgrounds {
			if runtimePanel.BaseBackground.FixedBackgroundID == templatepanel.FixedBackgroundID {
				break
			}
		}
		if !found {
			result[index].Deleted = true
		}
	}
	return result
}

func UnstripBackgrounds(templateBackgrounds []protocol.BaseBackground, backgrounds []protocol.RuntimeBackground) []protocol.RuntimeBackground {
	return pie.Map(backgrounds, func(background protocol.RuntimeBackground) protocol.RuntimeBackground {
		switch background.BaseBackground.BackgroundType {
		case "System":
			fromtemplate := templateBackgrounds[pie.FindFirstUsing(templateBackgrounds, func(tp protocol.BaseBackground) bool { return tp.FixedBackgroundID == background.BaseBackground.FixedBackgroundID })]
			background.BaseBackground.BackgroundContent = fromtemplate.BackgroundContent;
			return background
		case "Embedded":
			fromtemplate := templateBackgrounds[pie.FindFirstUsing(templateBackgrounds, func(tp protocol.BaseBackground) bool { return tp.FixedBackgroundID == background.BaseBackground.FixedBackgroundID})]
			background.BaseBackground.BackgroundContent = fromtemplate.BackgroundContent;
			return background
		case "Video", "Image", "Audio":
			return background
		default:
			log.Println("Unknown background type: " + background.BaseBackground.BackgroundType)
			return background
		}
	})
}

func PrepareBackgrounds(runtimeBackgrounds []protocol.RuntimeBackground) []protocol.RuntimeBackground {
	return pie.Map(runtimeBackgrounds, func(bg protocol.RuntimeBackground) protocol.RuntimeBackground {
		switch bg.BaseBackground.BackgroundType {
		case "System":
			data, err := os.ReadFile(path.Join(ADDONS_PATH, bg.BaseBackground.AppClientID, bg.BaseBackground.BackgroundContent));
			if err != nil {
				log.Println(err)
				bg.BaseBackground.BackgroundContent = "Error loading background content"
			}
			bg.BaseBackground.BackgroundContent = string(data)
			return bg;

		case "Embedded":
			url, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", bg.ControlPort))
			if err != nil {
				log.Println(err)
				return bg
			}
			url.Path += "/background"
			added := url.Query()

			added.Add("content", bg.BaseBackground.BackgroundContent)
			added.Add("clientid", bg.BaseBackground.AppClientID)
			added.Add("fixedbackgroundid", bg.BaseBackground.FixedBackgroundID)
			added.Add("uniquebackgroundid", bg.UniqueBackgroundID)
			url.RawQuery = added.Encode()
			bg.BaseBackground.BackgroundContent = fmt.Sprintf("%v", url)

			return bg
		default:
			return bg
		}
	})
}

func DeletedBackgrounds(runtimeBackgrounds []protocol.RuntimeBackground, templateBackgrounds []protocol.BaseBackground) []protocol.RuntimeBackground {
	deleted := []protocol.RuntimeBackground{}
	for _, runtimePanel := range runtimeBackgrounds {
		found := false
		for _, templatepanel := range templateBackgrounds {
			if runtimePanel.BaseBackground.FixedBackgroundID == templatepanel.FixedBackgroundID {
				found = true
				break
			}
		}
		if !found {
			deleted = append(deleted, runtimePanel)
		}
	}
	return deleted
}

func backgroundSystem(w http.ResponseWriter, r *http.Request) {
	upgrader2 := CustomUpgrader{}
	upgrader2.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader2.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	if _, err := os.Stat("lastactivebackground"); err != nil {
		os.WriteFile("lastactivebackground", []byte(""), 0644)
	}
	lastactivebackground, err := os.ReadFile("lastactivebackground")
	if err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, []byte("ERAD: Error loading last active backgrounds"))
		return
	}
	// Upgrade last active backgrounds
	result, err := LoadAllBackgroundsFromAddons()
	result.LastActiveBackgroundID = string(lastactivebackground)
	if err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, []byte("ERAD: Error loading backgrounds"))
		return
	}
	encoded, err := proto.Marshal(&result)
	if err != nil {
		log.Println("json:", err)
		c.WriteMessage(websocket.TextMessage, []byte("ERAD: Error encoding response"))
		return
	}
	fmt.Println("Sending backgrounds")
	//fmt.Println(string(encoded))
	err = c.WriteMessage(websocket.BinaryMessage, encoded)
	if err != nil {
		log.Println("write:", err)
		return
	}
	file := "storage/backgrounds.json"
	clientid := GenerateGUID()
	go func() {
		for {
			change := backgroundchangehub.WaitForMessage()
			if change.Stop && change.ClientID == clientid {
				break
			}
			if (change.ClientID == clientid) {
				continue
			}
			cr := protocol.RuntimeBackgroundChangeRequest{}
			cr.NewBackgrounds = pie.Map(change.NewBackgrounds, func(bg protocol.RuntimeBackground) *protocol.RuntimeBackground {
				return &bg
			})
			cr.NewActiveBackgroundID = change.NewActiveBackground
			marshalled, err := proto.Marshal(&cr)
			if err != nil {
				log.Println("json:", err)
				return
			}
			c.WriteMessage(websocket.BinaryMessage, marshalled)
		}
	}()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			backgroundchangehub.SendMessage(BackgroundDataChangeRequest{
				NewBackgrounds:  nil,
				NewActiveBackground: "",
				ClientID: clientid,
				Stop: false,
			})

			return
		}
		parsedLastActive := protocol.RuntimeBackgroundChangeRequest{}
		proto.Unmarshal(msg, &parsedLastActive)
		if (parsedLastActive.NewActiveBackgroundID != "") {
			addRecentMediaBackground(parsedLastActive.NewActiveBackgroundID, "", parsedLastActive.NewActiveBackgroundID)
		}

		backgroundchangehub.SendMessage(BackgroundDataChangeRequest{
			NewBackgrounds:  pie.Map(parsedLastActive.NewBackgrounds, func(bg *protocol.RuntimeBackground) protocol.RuntimeBackground {
				return *bg
			}),
			NewActiveBackground: parsedLastActive.NewActiveBackgroundID,
			ClientID: clientid,
			Stop: false,
		})

		if err != nil {
			log.Println("json:", err)
			return
		}

		persistent := protocol.PersistentRuntimeBackgrounds{}
		persistent.Backgrounds = pie.Map(UnstripBackgrounds(
			pie.Map(result.TemplateBackgrounds, func(bg *protocol.BaseBackground) protocol.BaseBackground {
				return *bg
			}),
			pie.Map(parsedLastActive.NewBackgrounds, func(bg *protocol.RuntimeBackground) protocol.RuntimeBackground {
				return *bg
			})), func(bg protocol.RuntimeBackground) *protocol.RuntimeBackground {
			return &bg
		})
		saved, err := protojson.Marshal(&persistent)
		if err != nil {
			log.Println("json:", err)
			return
		}
		err = os.WriteFile(file, saved, 0644)
		os.WriteFile("lastactivebackground", []byte(parsedLastActive.NewActiveBackgroundID), 0644)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func LoadAllBackgroundsFromAddons() (protocol.BackgroundSystemInitialResponse, error) {
	if _, err := os.Stat("addons"); err != nil {
		os.Mkdir("addons", 0755)
	}
	addons, err := os.ReadDir(ADDONS_PATH)
	if err != nil {
		log.Println(err)
		return protocol.BackgroundSystemInitialResponse{}, fmt.Errorf("Error reading addons directory: %v", err)
	}
	availableBackgrounds := make([]protocol.BaseBackground, 0)
	convertedAvailableBackgrounds := make([]protocol.RuntimeBackground, 0)
	if _, err := os.Stat("storage/backgrounds.json"); err != nil {
		file, err := os.Create("storage/backgrounds.json")
		file.Write([]byte("{}"))
		if err != nil {
			log.Println(err)
			return protocol.BackgroundSystemInitialResponse{}, fmt.Errorf("Error creating backgrounds.json: %v", err)
		}
		file.Close()
	}
	lastActiveBackgrounds, err := LoadAllRuntimeBackgrounds("storage/backgrounds.json")
	preparedLastActiveBackgrounds := make([]protocol.RuntimeBackground, 0)
	if err != nil {
		log.Println(err)
		return protocol.BackgroundSystemInitialResponse{}, fmt.Errorf("Error loading last active backgrounds: %v", err)
	}
	for i := 0; i < len(addons); i++ {
		addon := addons[i]
		if addon.IsDir() {
			if _, err := os.Stat(ADDONS_PATH + "/" + addon.Name() + "/disabled"); err == nil {
				fmt.Println("Addon " + addon.Name() + " is disabled")
				continue
			}
			if _, err := os.Stat(ADDONS_PATH + "/" + addon.Name() + "/manifest.json"); err != nil {
				log.Println(err)
				fmt.Errorf("Error loading manifest for addon %s: %v", addon.Name(), err)
				continue;
			}
			fmt.Println("Loading addon: " + addon.Name())
			manifest, err := LoadManifest(ADDONS_PATH + "/" + addon.Name() + "/manifest.json")

			if err != nil {
				log.Println(err)
				continue
			}
			declaredBackgrounds, err := LoadAllTemplateBackgrounds(ADDONS_PATH + "/" + addon.Name() + "/backgrounds.json")

			//Load all media as backgrounds
			if _ , err := os.Stat(ADDONS_PATH + "/" + addon.Name() + "/media/backgrounds"); err == nil {
				mediaBackgrounds, err := LoadAllMediaAsTemplateBackgrounds(ADDONS_PATH+"/"+addon.Name()+"/media/backgrounds", manifest)
				if err != nil {
					log.Println(err)
					return protocol.BackgroundSystemInitialResponse{}, fmt.Errorf("Error loading media backgrounds for addon %s: %v", addon.Name(), err)
				}
				declaredBackgrounds = append(declaredBackgrounds, mediaBackgrounds...)
			}

			declaredBackgrounds = pie.Map(declaredBackgrounds, func(bg protocol.BaseBackground) protocol.BaseBackground {
				bg.AppClientID = manifest.ClientID
				return bg
			})
			if err != nil {
				log.Println(err)
				continue
			}
			availableBackgrounds = append(availableBackgrounds, declaredBackgrounds...)
			converted := ConvertTemplateBackgroundsToRuntimeBackgrounds(availableBackgrounds)

			lastActiveBackgrounds = MeshBackgrounds(lastActiveBackgrounds, availableBackgrounds)
			ra, err := manifest.GetRAOrBootstrapAddon(addon.Name())

			backgroundsForAddon := pie.Filter(lastActiveBackgrounds, func(bg protocol.RuntimeBackground) bool {
				return bg.BaseBackground.AppClientID == manifest.ClientID
			})
			templateRuntimeBackgroundsForThisAddon := pie.Filter(converted, func(bg protocol.RuntimeBackground) bool {
				return bg.BaseBackground.AppClientID == manifest.ClientID
			})
			if err != nil {
				if err.Error() == "No bootstrap executable found for addon "+manifest.Name {
					log.Println("No bootstrap executable found for addon " + manifest.Name)
				} else {
					log.Println(err)
				}

			} else {

				for index, _ := range backgroundsForAddon {
					backgroundsForAddon[index].ControlPort = int32(ra.ControlPort)
				}
				for index, _ := range templateRuntimeBackgroundsForThisAddon {
					templateRuntimeBackgroundsForThisAddon[index].ControlPort = int32(ra.ControlPort)
				}

			}
			convertedAvailableBackgrounds = append(convertedAvailableBackgrounds, templateRuntimeBackgroundsForThisAddon...)
			preparedLastActiveBackgrounds = append(preparedLastActiveBackgrounds, backgroundsForAddon...)

		}
	}
	deletedBackgrounds := DeletedBackgrounds(lastActiveBackgrounds, availableBackgrounds)
	for _, bg := range deletedBackgrounds {
		for index, bg2 := range lastActiveBackgrounds {
			if bg2.BaseBackground.FixedBackgroundID == bg.BaseBackground.FixedBackgroundID {
				lastActiveBackgrounds[index].Deleted = true
			}
		}
	}
	preparedLastActiveBackgrounds = PrepareBackgrounds(preparedLastActiveBackgrounds)
	result := protocol.BackgroundSystemInitialResponse{
		InstancedButBlankBackgrounds:  pie.Map(convertedAvailableBackgrounds, func(bg protocol.RuntimeBackground) *protocol.RuntimeBackground {
			return &bg
		}),
		TemplateBackgrounds: 		 pie.Map(availableBackgrounds, func(bg protocol.BaseBackground) *protocol.BaseBackground {
			return &bg
		}),
		InstancedBackgroundsFromStorage: pie.Map(preparedLastActiveBackgrounds, func(bg protocol.RuntimeBackground) *protocol.RuntimeBackground {
			return &bg
		}),
		DeletedInstancedBackgrounds: pie.Map(deletedBackgrounds, func(bg protocol.RuntimeBackground) *protocol.RuntimeBackground {
			return &bg
		}),
	}
	return result, nil
}
