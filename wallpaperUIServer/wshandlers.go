package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ShutdownRequest struct {
	ClientID string
}

type RestartRequest struct {
	Stop bool
	ClientID string
}

var shutdownhub = CreateMessageHub[ShutdownRequest]()
var restarthub = CreateMessageHub[RestartRequest]()

func restartIPC(w http.ResponseWriter, r *http.Request) {
		upgrader2 := CustomUpgrader{}
	upgrader2.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader2.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	clientid := uuid.NewString()
	go func() {
		for {
			msg := restarthub.WaitForMessage()
			if msg.Stop && msg.ClientID == clientid {
				break
			}
			if msg.ClientID == clientid {
				continue;
			}
			if msg.Stop {
				continue;
			}
			err := c.WriteMessage(websocket.TextMessage, []byte("restart"))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}()
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read 2:", err)
			restarthub.SendMessage(RestartRequest{true, clientid})
			break
		}
		restarthub.SendMessage(RestartRequest{false, clientid})
	}

}

func popupIPC(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	go func() {
		for {
			fmt.Println("Waiting for message")
			msg := <-popupchannel
			fmt.Println("Sending message: ", msg)
			if msg.popup_URL == "stop" {
				break
			}
			data := make(map[string]interface{})
			if msg.Type == "popup" {
				data["Type"] = "popup"
				data["popup_URL"] = msg.popup_URL
				data["popup_ClientID"] = msg.popup_ClientID
				data["popup_AppName"] = msg.popup_AppName
				data["popup_Favicon"] = msg.popup_Favicon
				data["popup_Title"] = msg.popup_Title
				data["trackingID"] = msg.trackingID
			} else if msg.Type == "input" {
				data["Type"] = "input"
				data["input_Type"] = msg.input_Type
				data["input_Placeholder"] = msg.input_Placeholder
				data["input_MaxLength"] = msg.input_MaxLength
				data["trackingID"] = msg.trackingID
			} else if msg.Type == "stop" {
				data["Type"] = "stop"
			}
			encoded, err := json.Marshal(data)
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
		fmt.Println("Received message: ", string(message))
		if err != nil {
			popupchannel <- CreatePopupRequest("stop", "stop", "stop", "stop", "stop")
			log.Println("read:", err)
			break
		}
		decoded := string(message)
		parsed := make(map[string]interface{})
		err = json.Unmarshal([]byte(decoded), &parsed)
		if err != nil {
			log.Println("json:", err)
			break
		}
		cancelled := parsed["cancelled"].(bool)
		if parsed["Type"] == "popup" {
			popupresultchannel.SendMessage(CreatePopupResult(parsed["popup_ResultData"].(string), parsed["trackingID"].(string), cancelled))
		} else if parsed["Type"] == "input" {
			popupresultchannel.SendMessage(CreateInputResult(parsed["input_ResultData"].(string), parsed["trackingID"].(string), cancelled))
		}
	}
}

func shutdownHub(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	clientid := uuid.NewString()
	go func() {
		for {
			shutdownhub.WaitForMessage()
			err := c.WriteMessage(websocket.TextMessage, []byte("shutdown"))
			if err != nil {
				log.Println("write:", err)
				return
			}
			os.Exit(0)
		}
	}()
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		request := ShutdownRequest{
			ClientID: clientid,
		}
		shutdownhub.SendMessage(request)
	}
}


func storageHub(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	// Intro message for scope
	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	decoded := string(message)
	clientid := uuid.NewString()
	scope := decoded
	if _, ok := storagechanged[scope]; !ok {
		storagechanged[scope] = CreateMessageHub[StorageChange]()
	}
	if !doesScopeExist(scope) {
		dp, err := os.Create(path.Join("storage", scope))
		if err != nil {
			log.Println("create:", err)
			return
		}
		dp.Close()
	}
	file := path.Join("storage", scope)
	// Read the existing file
	// Write the existing file to the client
	data, err := os.ReadFile(file)
	if err != nil {
		log.Println("read:", err)
		return
	}
	err = c.WriteMessage(mt, data)
	if err != nil {
		log.Println("write:", err)
		return
	}

	go func() {
		for {
			obj := storagechanged[scope];
			msg := obj.WaitForMessage()
			if msg.Data == "stop" && msg.Scope == "stop" && msg.ClientID == clientid {
				break
			}
			if msg.Scope != scope {
				continue
			}

			if msg.ClientID == clientid {
				continue
			}

			err := c.WriteMessage(mt, []byte(msg.Data))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			obj := storagechanged[scope]
			obj.SendMessage(CreateStorageChange("stop", string("stop"), clientid))
			break
		}
		os.WriteFile(file, message, 0644)
		obj := storagechanged[scope]
		obj.SendMessage(CreateStorageChange(scope, string(message), clientid))
	}

}