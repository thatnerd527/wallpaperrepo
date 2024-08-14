package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"wallpaperuiserver/protocol"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type ShutdownRequest struct {
	ClientID string
}

type RestartRequest struct {
	Stop bool
	Message2 string
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
			err := c.WriteMessage(websocket.TextMessage, []byte(msg.Message2))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}()
	for {
		_, p, err := c.ReadMessage()
		if err != nil {
			log.Println("read 2:", err)
			restarthub.SendMessage(RestartRequest{true, string(p),clientid})
			break
		}
		restarthub.SendMessage(RestartRequest{false, string(p), clientid})
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
			buf := protocol.PopupAppControlMessage{}
			if msg.Type == "popup" {
				buf.Type = protocol.PopupAppControlMessage_POPUP
				buf.Message = &protocol.PopupAppControlMessage_PopupRequest{
					PopupRequest: &protocol.PopupRequest{
						URL: msg.popup_URL,
						ClientID: msg.popup_ClientID,
						AppName: msg.popup_AppName,
						Favicon: msg.popup_Favicon,
						Title: msg.popup_Title,
						AutoAuthorize: false,
						RequestID: msg.trackingID,
					},
				}
			} else if msg.Type == "input" {
				buf.Type = protocol.PopupAppControlMessage_INPUT
				buf.Message = &protocol.PopupAppControlMessage_InputRequest{
					InputRequest: &protocol.InputRequest{
						InputType: msg.input_Type,
						InputPlaceholder: msg.input_Placeholder,
						MaxLength: int32(msg.input_MaxLength),
						RequestID: msg.trackingID,
					},
				}
			} else if msg.Type == "stop" {
				buf.Type = protocol.PopupAppControlMessage_SHUTDOWN
				buf.Message = &protocol.PopupAppControlMessage_ShutdownRequest{
					ShutdownRequest: &protocol.ShutdownRequest{

					},
				}
			}
			data2, err := proto.Marshal(&buf)
			if err != nil {
				log.Println("proto:", err)
				break
			}
			err = c.WriteMessage(websocket.TextMessage, data2)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			popupchannel <- CreatePopupRequest("stop", "stop", "stop", "stop", "stop")
			log.Println("read:", err)
			break
		}
		parsed := protocol.PopupAppResponse{}
		err = proto.Unmarshal(message, &parsed)
		if err != nil {
			log.Println("json:", err)
			break
		}
		if parsed.Type == protocol.PopupAppResponse_POPUP {
			popupresultchannel.SendMessage(CreatePopupResult(parsed.GetPopupResponse().ResultData, parsed.GetPopupResponse().GetRequestID(), parsed.Cancelled))
		} else if parsed.Type == protocol.PopupAppResponse_INPUT {
			popupresultchannel.SendMessage(CreateInputResult(parsed.GetInputResponse().ResultData, parsed.GetInputResponse().GetRequestID(), parsed.Cancelled))
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