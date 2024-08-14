package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
	"wallpaperuiserver/protocol"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type PreferenceUpdate struct {
	prefdict protocol.AppSettings
	ClientID string
	Stop     bool
}

var preferenceschannel = CreateMessageHub[PreferenceUpdate]()
var cachedpreferences Preferences = Preferences{prefs: protocol.AppSettings{}}

// DeepEqual compares two maps deeply
func DeepEqual(map1, map2 map[string]interface{}) bool {
	return deepEqual(map1, map2)
}

// deepEqual is a helper function that handles the actual deep comparison
func deepEqual(val1, val2 interface{}) bool {
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return false
	}

	switch v1 := val1.(type) {
	case map[string]interface{}:
		v2, ok := val2.(map[string]interface{})
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for k, v := range v1 {
			if !deepEqual(v, v2[k]) {
				return false
			}
		}
		return true
	case []interface{}:
		v2, ok := val2.([]interface{})
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for i := range v1 {
			if !deepEqual(v1[i], v2[i]) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(val1, val2)
	}
}


func isJsonSame(a map[string]interface{}, b map[string]interface{}) bool {
	comp1, err := json.Marshal(a)
	if err != nil {
		return false
	}
	comp2, err := json.Marshal(b)
	if err != nil {
		return false
	}
	fmt.Println("SAME: ", string(comp1) == string(comp2))
	return string(comp1) == string(comp2)
}

func preferencesSystem(w http.ResponseWriter, r *http.Request) {
		upgrader2 := CustomUpgrader{}
	upgrader2.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader2.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	clientid := GenerateGUID()
	if _ , ok := os.Stat("preferences.json"); os.IsNotExist(ok) {
		file, err := os.Create("preferences.json")
		if err != nil {
			log.Println("create:", err)
			return
		}
		_, err = file.WriteString("{}")
		if err != nil {
			log.Println("write:", err)
			return
		}
		file.Close()
	}
	parsed := protocol.AppSettings{}

	file, err := os.ReadFile("preferences.json")
	if err != nil {
		log.Println("read:", err)
		return
	}
	err = protojson.Unmarshal(file, &parsed)
	cachedpreferences.Write(func(prefs protocol.AppSettings) protocol.AppSettings {return parsed})
	if err != nil {
		log.Println("json:", err)
		return
	}
	encoded, err := proto.Marshal(&parsed)
	if err != nil {
		log.Println("json:", err)
		return
	}
	err = c.WriteMessage(websocket.BinaryMessage, encoded)
	if err != nil {
		log.Println("write:", err)
		return
	}


	go func() {
		for {
			fmt.Println("Waiting for message")
			msg := preferenceschannel.WaitForMessage()
			//fmt.Println("Sending message: ", msg)
			fmt.Println("PREFUPDATE")
			if msg.Stop && msg.ClientID == clientid {
				break
			}
			if msg.ClientID == clientid {
				continue;
			}
			parsed := msg.prefdict
			encoded, err := proto.Marshal(&parsed)
			if err != nil {
				log.Println("json:", err)
				break
			}
			err = c.WriteMessage(websocket.BinaryMessage, encoded)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			preferenceschannel.SendMessage(PreferenceUpdate{protocol.AppSettings{}, clientid, true})
			break
		}
		parsed := protocol.AppSettings{}
		//os.WriteFile("preferences.json", message, 0644)

		err = proto.Unmarshal(message, &parsed)
		if proto.Equal(&parsed,referenceify(cachedpreferences.Read())) {
			continue
		}
		cachedpreferences.Write(func(prefs protocol.AppSettings) protocol.AppSettings {return parsed})

		if err != nil {
			log.Println("json:", err)
			break
		}
		preferenceschannel.SendMessage(PreferenceUpdate{parsed, clientid, false})
	}
}

type Preferences struct {
	prefs protocol.AppSettings
	writeHandlers []func(protocol.AppSettings) protocol.AppSettings
	writeLock bool
	mut sync.Mutex
}

func (p *Preferences) AddWriteHandler(handler func(protocol.AppSettings) protocol.AppSettings) {
	p.writeHandlers = append(p.writeHandlers, handler)
}

func (p *Preferences) Write(modifier func(protocol.AppSettings) protocol.AppSettings) {
	p.mut.Lock()
	cloned := proto.Clone(&p.prefs).(*protocol.AppSettings)
	*cloned = p.PreventNulls(*cloned)
	operated := modifier(*cloned)
	p.prefs = *proto.Clone(&operated).(*protocol.AppSettings)
	for _, handler := range p.writeHandlers {
		cloned := proto.Clone(&p.prefs).(*protocol.AppSettings)
		operated := handler(*cloned)
		p.prefs = *proto.Clone(&operated).(*protocol.AppSettings)
	}
	p.mut.Unlock()
}

func (p *Preferences) Read() protocol.AppSettings {
	return p.PreventNulls(*proto.Clone(&p.prefs).(*protocol.AppSettings))
}

func (p *Preferences) PreventNulls(p2 protocol.AppSettings) protocol.AppSettings {
	if p2.RecentBackgroundStore == nil {
		p2.RecentBackgroundStore = &protocol.RecentBackgroundStore{RecentBackgrounds: map[string]*protocol.RecentBackground{}}
	}
	if p2.RecentBackgroundStore.RecentBackgrounds == nil {
		p2.RecentBackgroundStore.RecentBackgrounds = map[string]*protocol.RecentBackground{}
	}
	if p2.SimpleBackgroundsSystem == nil {
		p2.SimpleBackgroundsSystem = &protocol.SimpleBackgroundsSystem{}
	}
	if p2.SimpleBackgroundsSystem.SimpleBackgrounds == nil {
		p2.SimpleBackgroundsSystem.SimpleBackgrounds = []*protocol.SimpleBackground{}
	}
	if p2.RecentColorSystem == nil {
		p2.RecentColorSystem = &protocol.RecentColorSystem{RecentColors: []*protocol.RecentColor{}}
	}
	if p2.RecentColorSystem.RecentColors == nil {
		p2.RecentColorSystem.RecentColors = []*protocol.RecentColor{}
	}
	return p2
}