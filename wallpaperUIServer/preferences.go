package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/websocket"
)

type PreferenceUpdate struct {
	prefdict map[string]interface{}
	ClientID string
	Stop     bool
}

var preferenceschannel = CreateMessageHub[PreferenceUpdate]()
var cachedpreferences map[string]interface{} = make(map[string]interface{})

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
	parsed := make(map[string]interface{})
	file, err := os.ReadFile("preferences.json")
	if err != nil {
		log.Println("read:", err)
		return
	}
	err = json.Unmarshal(file, &parsed)
	if err != nil {
		log.Println("json:", err)
		return
	}
	encoded, err := json.Marshal(parsed)
	if err != nil {
		log.Println("json:", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, encoded)
	if err != nil {
		log.Println("write:", err)
		return
	}
	cachedpreferences = parsed

	go func() {
		for {
			fmt.Println("Waiting for message")
			msg := preferenceschannel.WaitForMessage()
			//fmt.Println("Sending message: ", msg)
			fmt.Println("PREFUPDATE")
			if msg.Stop && msg.ClientID == clientid {
				break
			}
			if msg.prefdict == nil {
				continue;
			}
			if msg.ClientID == clientid {
				continue;
			}
			parsed := msg.prefdict
			encoded, err := json.Marshal(parsed)
			if err != nil {
				log.Println("json:", err)
				break
			}
			err = c.WriteMessage(websocket.TextMessage, encoded)
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
			preferenceschannel.SendMessage(PreferenceUpdate{nil, clientid, true})
			break
		}
		var parsed map[string]interface{}
		os.WriteFile("preferences.json", message, 0644)

		err = json.Unmarshal(message, &parsed)
		if isJsonSame(parsed, cachedpreferences) {
			continue
		}
		cachedpreferences = parsed

		if err != nil {
			log.Println("json:", err)
			break
		}
		preferenceschannel.SendMessage(PreferenceUpdate{parsed, clientid, false})
	}
}