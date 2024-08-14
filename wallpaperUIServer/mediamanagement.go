package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"wallpaperuiserver/protocol"

	"github.com/elliotchance/pie/v2"
	"github.com/sqweek/dialog"
	"google.golang.org/protobuf/proto"
)

var guidtofilepath = make(map[string]string)

func FileNameToMediaType(filename string) string {
	switch filepath.Ext(filename) {
	case ".mp4", ".webm", ".mov", ".avi", ".mkv", ".flv", ".wmv", ".mpg", ".mpeg", ".m4v", ".3gp", ".3g2":
		return "Video"
	case ".jpg", ".jpeg", ".png", ".avif", ".webp", ".gif":
		return "Image"
	default:
		return "unknown"
	}
}

func setBackgroundFromCache(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		if !r.URL.Query().Has("delete") {
			backgrounds := cachedpreferences.Read().SimpleBackgroundsSystem.SimpleBackgrounds
			backgroundFiltered := pie.Filter(referenceToArray(backgrounds), func(s protocol.SimpleBackground) bool {
					return s.BackgroundID == r.URL.Query().Get("backgroundid")
			})
			background := backgroundFiltered[0]
			if _, err := os.Stat(background.FilePath); err == nil {
				cachedpreferences.Write(func(as protocol.AppSettings) protocol.AppSettings {
					as.SimpleBackgroundsSystem.IsSimpleBackgroundEnabled = true
					as.SimpleBackgroundsSystem.ActiveSimpleBackgroundID = r.URL.Query().Get("backgroundid")
					return as
				})
				backgroundid := r.URL.Query().Get("backgroundid")

				addRecentMediaBackground(backgroundid, path.Base(background.FilePath), "")
				fmt.Println("Set simple background to " + background.FilePath)

			} else {
				log.Println("File not found")
				w.WriteHeader(http.StatusNotFound)
			}
			return
		} else {
			idtofilename := pie.Filter(referenceToArray(cachedpreferences.Read().SimpleBackgroundsSystem.SimpleBackgrounds), func(s protocol.SimpleBackground) bool {
				return s.BackgroundID == r.URL.Query().Get("backgroundid")
			})[0].FilePath
			if _, err := os.Stat(idtofilename); err == nil {
				err := os.Remove(idtofilename)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

func simpleBackgroundHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		background := cachedpreferences.Read().SimpleBackgroundsSystem.SimpleBackgrounds
		filtered := pie.Filter(referenceToArray(background), func(s protocol.SimpleBackground) bool {
			return s.BackgroundID == r.URL.Query().Get("backgroundid")
		})
		w.Header().Add("Access-Control-Allow-Origin", "*")
		if len(filtered) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			http.ServeFile(w, r, filtered[0].FilePath)
			return
		}

	} else if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parsed := protocol.SimpleBackgroundRequest{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	proto.Unmarshal(data, &parsed)
	if parsed.RequestType == protocol.SimpleBackgroundRequest_DELETE {
		cachedpreferences.Write(func(as protocol.AppSettings) protocol.AppSettings {
			as.SimpleBackgroundsSystem.IsSimpleBackgroundEnabled = false
			as.SimpleBackgroundsSystem.ActiveSimpleBackgroundID = ""
			return as
		})
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	if parsed.RequestType == protocol.SimpleBackgroundRequest_ADD {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		filename := parsed.ResultFileName
		zipPath2, err := dialog.File().Title("Select simple background file").
			Filter("Image files", "jpg", "jpeg", "png", "avif", "webp", "gif").
			Filter("Video Files", "mp4", "webm", "mov", "avi", "mkv", "flv", "wmv", "mpg", "mpeg", "m4v", "3gp", "3g2").
			Filter("Other files", "*").Load()
		if err != nil {
			log.Println(err)
			tmp := protocol.SimpleBackgroundResponse{}
			clone := err.Error()
			tmp.ErrorMessage = clone
			tmp.ResponseType = protocol.SimpleBackgroundResponse_CANCELLED
			tmp.EncodingTicketID = ""
			tmp.StatusMessage = "cancel"
			data, err := proto.Marshal(&tmp)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(data)
			return
		}

		encodingguid := GenerateGUID()
		encodertochan[encodingguid] = make(chan string)
		resultfile, err := startReencoding(zipPath2, encodertochan[encodingguid], filename,parsed)
		if err != nil {
			log.Println(err)
			tmp := protocol.SimpleBackgroundResponse{}
			clone := err.Error()
			tmp.ErrorMessage = clone
			tmp.ResponseType = protocol.SimpleBackgroundResponse_FAILURE
			tmp.EncodingTicketID = ""
			tmp.StatusMessage = "errorprestart"
			data, err := proto.Marshal(&tmp)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write(data)
			return
		}
		tmp := protocol.SimpleBackgroundResponse{}
		tmp.ErrorMessage = ""
		tmp.ResponseType = protocol.SimpleBackgroundResponse_SUCCESS
		tmp.EncodingTicketID = encodingguid
		tmp.ResultFilePath = resultfile
		tmp.StatusMessage = "started"
		tmp.OriginalFilePath = zipPath2
		data, err := proto.Marshal(&tmp)
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(data)
		return
	}
}

func getEncodingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "GET" {
		if val, ok := encodertochan[r.URL.Query().Get("guid")]; ok {
			if r.URL.Query().Get("action") == "status" {
				result := <-val
				_, err := w.Write([]byte(result))
				if err != nil {
					log.Println("GETENCODINGSTATUS " + err.Error())
					val <- result
				}

			} else {
				var tmp = make(map[string]string)
				data2 := <-val
				tmp["status"] = data2
				data, err := json.Marshal(tmp)
				if err != nil {
					log.Println("GETENCODINGSTATUS " + err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, err = w.Write(data)
				if err != nil {
					log.Println("GETENCODINGSTATUS " + err.Error())
					val <- data2
				}
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	if r.Method == "POST" {
		if val, ok := encodertochan[r.URL.Query().Get("guid")]; ok {
			if r.URL.Query().Has("action") {
				val <- r.URL.Query().Get("action")

			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func LoadAllMediaAsPanels(panelspath string, manifest AddonManifest) ([]protocol.BasePanel, error) {
	files, err := os.ReadDir(panelspath)
	if err != nil {
		log.Println(err)
		return []protocol.BasePanel{}, err
	}
	panels := []protocol.BasePanel{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		guid := GenerateGUID()
		filename := file.Name()
		panel := protocol.BasePanel{}
		panel.PanelType = FileNameToMediaType(filename)
		panel.FixedPanelID = filename + "_panel_" + manifest.Name + "_" + manifest.ClientID
		url, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v", secureport))
		if err != nil {
			log.Println(err)
			return []protocol.BasePanel{}, err
		}
		url.Path = "/mediaregistry"
		q := url.Query()
		q.Add("guid", guid)
		url.RawQuery = q.Encode()
		panel.PanelContent = url.String()
		panel.PanelTitle = filename
		guidtofilepath[guid] = filepath.Join(panelspath, filename)
		panels = append(panels, panel)
	}
	return panels, nil
}

func LoadAllMediaAsTemplateBackgrounds(backgroundspath string, manifest AddonManifest) ([]protocol.BaseBackground, error) {
	files, err := os.ReadDir(backgroundspath)
	if err != nil {
		log.Println(err)
		return []protocol.BaseBackground{}, err
	}
	backgrounds := []protocol.BaseBackground{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		guid := GenerateGUID()
		filename := file.Name()
		background := protocol.BaseBackground{}
		background.BackgroundType = FileNameToMediaType(filename)
		fmt.Println("MEDIADEBG " + filename + "_background_" + manifest.Name + "_" + manifest.ClientID)
		background.FixedBackgroundID = filename + "_background_" + manifest.Name + "_" + manifest.ClientID
		url, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v", secureport))
		if err != nil {
			log.Println(err)
			return []protocol.BaseBackground{}, err
		}
		url.Path = "/mediaregistry"
		q := url.Query()
		q.Add("guid", guid)
		url.RawQuery = q.Encode()
		background.BackgroundContent = url.String()
		guidtofilepath[guid] = filepath.Join(backgroundspath, filename)
		backgrounds = append(backgrounds, background)
	}
	return backgrounds, nil
}
