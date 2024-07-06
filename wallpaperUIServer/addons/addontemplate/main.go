package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/websocket"
)

type CustomUpgrader struct {
	websocket.Upgrader
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("public", "wwwroot", r.URL.Path))
}

func keepAlive(w http.ResponseWriter, r *http.Request) {
	upgrader2 := CustomUpgrader{}
	upgrader2.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader2.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		time.Sleep(20 * time.Second)
		err := c.WriteMessage(websocket.PingMessage, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func ipcHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Received message: %s\n", message)
	}
}

func shutdown(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Shutting down")
	os.Exit(0)
}

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()
	http.HandleFunc("/", fileHandler)
	http.HandleFunc("/keepalive", keepAlive)
	http.HandleFunc("/ipc", ipcHandler)
	http.HandleFunc("/shutdown", shutdown)
	fmt.Printf("Listening on port %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d",*port), nil)
	if err != nil {
		fmt.Println(err)
	}
}