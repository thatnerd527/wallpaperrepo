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

func isJsonSame(a map[string]interface{}, b map[string]interface{}) bool {
	return reflect.DeepEqual(a, b)
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
			fmt.Println("Sending message: ", msg)
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