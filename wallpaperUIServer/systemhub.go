package main

import (
	"fmt"
	"net/http"
	"time"
)

func restartHandler(w http.ResponseWriter, r *http.Request) {
	// Restart the popup app

	// Send a message to the popup app to restart
	popupchannel <- PopupRequest{Type: "stop"}
	fmt.Println("Sent stop message")
	time.Sleep(5 * time.Second)
	startPopupApp()
	fmt.Println("Started popup app")


	// Shutdown all addons
	for _, addon := range RuntimeAddons {
		addon.Shutdown()
	}
	fmt.Println("Shutted down all addons")
	w.Header().Set("Access-Control-Allow-Origin", "*")



}