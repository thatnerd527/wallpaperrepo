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

type TemplateCustomBackground struct {
	LoaderBackgroundID    string
	BackgroundType        string
	BackgroundContent     string
	BackgroundDefaultData string
	ClientID              string
}

type RuntimeCustomBackground struct {
	TemplateCustomBackground
	PersistentBackgroundID   string
	PersistentBackgroundData string
	Deleted                  bool
	ControlPort              int
}

type ResultResponse2 struct {
	AvailableBackgrounds []RuntimeCustomBackground
	AvailableTemplateBackgrounds []TemplateCustomBackground
	LastActiveBackgrounds []RuntimeCustomBackground
	LastActiveBackgroundID string
	DeletedBackgrounds []RuntimeCustomBackground
}

type SaveRequest struct {
	Backgrounds []RuntimeCustomBackground
	LastActiveBackgroundID string
}

var backgroundchangehub = CreateMessageHub[BackgroundDataChangeRequest]()

func LoadAllTemplateBackgrounds(filename string) ([]TemplateCustomBackground, error) {
	backgroundsFile, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return []TemplateCustomBackground{}, err
	}
	backgrounds := []TemplateCustomBackground{}
	err = json.Unmarshal(backgroundsFile, &backgrounds)
	if err != nil {
		log.Println(err)
		return []TemplateCustomBackground{}, err
	}
	return backgrounds, nil
}

func LoadAllRuntimeBackgrounds(filename string) ([]RuntimeCustomBackground, error) {
	backgroundsFile, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return []RuntimeCustomBackground{}, err
	}
	backgrounds := []RuntimeCustomBackground{}
	err = json.Unmarshal(backgroundsFile, &backgrounds)
	if err != nil {
		log.Println(err)
		return []RuntimeCustomBackground{}, err
	}
	return backgrounds, nil
}

func ConvertTemplateBackgroundsToRuntimeBackgrounds(backgrounds []TemplateCustomBackground) []RuntimeCustomBackground {
	runtimeBackgrounds := []RuntimeCustomBackground{}
	for _, background := range backgrounds {
		runtimeBackgrounds = append(runtimeBackgrounds, RuntimeCustomBackground{
			TemplateCustomBackground: background,
			PersistentBackgroundID:   GenerateGUID(),
			PersistentBackgroundData: background.BackgroundDefaultData,
			Deleted:                  false,
		})
	}
	return runtimeBackgrounds
}

func MeshBackgrounds(runtimeBackgrounds []RuntimeCustomBackground, templateBackgrounds []TemplateCustomBackground) []RuntimeCustomBackground {
	result := append([]RuntimeCustomBackground{}, runtimeBackgrounds...)
	for _, templateBackground := range templateBackgrounds {
		index := pie.FindFirstUsing(runtimeBackgrounds, func(bg RuntimeCustomBackground) bool {
			return bg.LoaderBackgroundID == templateBackground.LoaderBackgroundID
		})
		if index == -1 {
			continue
		}
		result[index] = RuntimeCustomBackground{
			TemplateCustomBackground: templateBackground,
			PersistentBackgroundID:   runtimeBackgrounds[index].PersistentBackgroundID,
			PersistentBackgroundData: runtimeBackgrounds[index].PersistentBackgroundData,
			Deleted:                  runtimeBackgrounds[index].Deleted,
		}
	}
	for index, runtimePanel := range runtimeBackgrounds {
		found := false
		for _, templatepanel := range templateBackgrounds {
			if runtimePanel.LoaderBackgroundID == templatepanel.LoaderBackgroundID {
				break
			}
		}
		if !found {
			result[index].Deleted = true
		}
	}
	return result
}

func UnstripBackgrounds(templateBackgrounds []TemplateCustomBackground, backgrounds []RuntimeCustomBackground) []RuntimeCustomBackground {
	return pie.Map(backgrounds, func(background RuntimeCustomBackground) RuntimeCustomBackground {
		switch background.BackgroundType {
		case "System":
			fromtemplate := templateBackgrounds[pie.FindFirstUsing(templateBackgrounds, func(tp TemplateCustomBackground) bool { return tp.LoaderBackgroundID == background.LoaderBackgroundID })]
			background.BackgroundContent = fromtemplate.BackgroundContent;
			return background
		case "Embedded":
			fromtemplate := templateBackgrounds[pie.FindFirstUsing(templateBackgrounds, func(tp TemplateCustomBackground) bool { return tp.LoaderBackgroundID == background.LoaderBackgroundID })]
			background.BackgroundContent = fromtemplate.BackgroundContent;
			return background
		case "Video", "Image", "Audio":
			return background
		default:
			log.Println("Unknown background type: " + background.BackgroundType)
			return background
		}
	})
}

func PrepareBackgrounds(runtimeBackgrounds []RuntimeCustomBackground) []RuntimeCustomBackground {
	return pie.Map(runtimeBackgrounds, func(bg RuntimeCustomBackground) RuntimeCustomBackground {
		switch bg.BackgroundType {
		case "System":
			data, err := os.ReadFile(path.Join(ADDONS_PATH, bg.ClientID, bg.BackgroundContent));
			if err != nil {
				log.Println(err)
				bg.BackgroundContent = "Error loading panel content"
			}
			bg.BackgroundContent = string(data)
			return bg;

		case "Embedded":
			url, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", bg.ControlPort))
			if err != nil {
				log.Println(err)
				return bg
			}
			url.Path += "/background"
			added := url.Query()

			added.Add("content", bg.BackgroundContent)
			added.Add("clientid", bg.ClientID)
			added.Add("loaderbackgroundid", bg.LoaderBackgroundID)
			added.Add("persistentbackgroundid", bg.PersistentBackgroundID)
			url.RawQuery = added.Encode()
			bg.BackgroundContent = fmt.Sprintf("%v", url)

			return bg
		default:
			return bg
		}
	})
}

func DeletedBackgrounds(runtimeBackgrounds []RuntimeCustomBackground, templateBackgrounds []TemplateCustomBackground) []RuntimeCustomBackground {
	deleted := []RuntimeCustomBackground{}
	for _, runtimePanel := range runtimeBackgrounds {
		found := false
		for _, templatepanel := range templateBackgrounds {
			if runtimePanel.LoaderBackgroundID == templatepanel.LoaderBackgroundID {
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
	encoded, err := json.Marshal(result)
	if err != nil {
		log.Println("json:", err)
		c.WriteMessage(websocket.TextMessage, []byte("ERAD: Error encoding response"))
		return
	}
	fmt.Println("Sending backgrounds")
	//fmt.Println(string(encoded))
	err = c.WriteMessage(websocket.TextMessage, encoded)
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
			marshalled, err := json.Marshal(BackgroundDataChangeRequest2{
				NewBackgrounds: change.NewBackgrounds,
				NewActiveBackground: change.NewActiveBackground,
			})
			if err != nil {
				log.Println("json:", err)
				return
			}
			c.WriteMessage(websocket.TextMessage, marshalled)
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
		decoded := string(msg)
		parsed := BackgroundDataChangeRequest2{}

		err = json.Unmarshal([]byte(decoded), &parsed)
		backgroundchangehub.SendMessage(BackgroundDataChangeRequest{
			NewBackgrounds:  parsed.NewBackgrounds,
			NewActiveBackground: parsed.NewActiveBackground,
			ClientID: clientid,
			Stop: false,
		})
		if err != nil {
			log.Println("json:", err)
			return
		}
		parsed.NewBackgrounds = UnstripBackgrounds(result.AvailableTemplateBackgrounds, parsed.NewBackgrounds)
		saved, err := json.Marshal(parsed.NewBackgrounds)
		if err != nil {
			log.Println("json:", err)
			return
		}
		err = os.WriteFile(file, saved, 0644)
		os.WriteFile("lastactivebackground", []byte(parsed.NewActiveBackground), 0644)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func LoadAllBackgroundsFromAddons() (ResultResponse2, error) {
	if _, err := os.Stat("addons"); err != nil {
		os.Mkdir("addons", 0755)
	}
	addons, err := os.ReadDir(ADDONS_PATH)
	if err != nil {
		log.Println(err)
		return ResultResponse2{}, fmt.Errorf("Error reading addons directory: %v", err)
	}
	availableBackgrounds := make([]TemplateCustomBackground, 0)
	convertedAvailableBackgrounds := make([]RuntimeCustomBackground, 0)
	if _, err := os.Stat("storage/backgrounds.json"); err != nil {
		file, err := os.Create("storage/backgrounds.json")
		file.Write([]byte("[]"))
		if err != nil {
			log.Println(err)
			return ResultResponse2{}, fmt.Errorf("Error creating backgrounds.json: %v", err)
		}
		file.Close()
	}
	lastActiveBackgrounds, err := LoadAllRuntimeBackgrounds("storage/backgrounds.json")
	preparedLastActiveBackgrounds := make([]RuntimeCustomBackground, 0)
	if err != nil {
		log.Println(err)
		return ResultResponse2{}, fmt.Errorf("Error loading last active backgrounds: %v", err)
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
					return ResultResponse2{}, fmt.Errorf("Error loading media backgrounds for addon %s: %v", addon.Name(), err)
				}
				declaredBackgrounds = append(declaredBackgrounds, mediaBackgrounds...)
			}

			declaredBackgrounds = pie.Map(declaredBackgrounds, func(bg TemplateCustomBackground) TemplateCustomBackground {
				bg.ClientID = manifest.ClientID
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

			backgroundsForAddon := pie.Filter(lastActiveBackgrounds, func(bg RuntimeCustomBackground) bool {
				return bg.ClientID == manifest.ClientID
			})
			templateRuntimeBackgroundsForThisAddon := pie.Filter(converted, func(bg RuntimeCustomBackground) bool {
				return bg.ClientID == manifest.ClientID
			})
			if err != nil {
				if err.Error() == "No bootstrap executable found for addon "+manifest.Name {
					log.Println("No bootstrap executable found for addon " + manifest.Name)
				} else {
					log.Println(err)
				}

			} else {

				for index, _ := range backgroundsForAddon {
					backgroundsForAddon[index].ControlPort = ra.ControlPort
				}
				for index, _ := range templateRuntimeBackgroundsForThisAddon {
					templateRuntimeBackgroundsForThisAddon[index].ControlPort = ra.ControlPort
				}

			}
			convertedAvailableBackgrounds = append(convertedAvailableBackgrounds, templateRuntimeBackgroundsForThisAddon...)
			preparedLastActiveBackgrounds = append(preparedLastActiveBackgrounds, backgroundsForAddon...)

		}
	}
	deletedBackgrounds := DeletedBackgrounds(lastActiveBackgrounds, availableBackgrounds)
	for _, bg := range deletedBackgrounds {
		for index, bg2 := range lastActiveBackgrounds {
			if bg2.LoaderBackgroundID == bg.LoaderBackgroundID {
				lastActiveBackgrounds[index].Deleted = true
			}
		}
	}
	preparedLastActiveBackgrounds = PrepareBackgrounds(preparedLastActiveBackgrounds)
	result := ResultResponse2{
		AvailableBackgrounds:         PrepareBackgrounds(convertedAvailableBackgrounds),
		AvailableTemplateBackgrounds: availableBackgrounds,
		LastActiveBackgrounds:        preparedLastActiveBackgrounds,
		DeletedBackgrounds:           deletedBackgrounds,
	}
	return result, nil
}
