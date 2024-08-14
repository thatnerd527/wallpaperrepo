package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
	"time"
	"wallpaperuiserver/protocol"
)

var encodertochan = make(map[string]chan string)

type RecentBackground struct {
	CachedPath             string `json:"cachedpath"`
	PersistentBackgroundID string `json:"persistentbackgroundid"`
	TimestampAddedNanos    string  `json:"timestampaddednanos"`
}

func previewFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "GET" {
		filename := r.URL.Query().Get("filename")
		http.ServeFile(w, r, path.Join("result", filename))
	}
}

func hashCode(s string) int {
	h := 0
	for _, char := range s {
		h = 31*h + int(char)
	}
	return h
}

func addRecentMediaBackground(backgroundid string, filename string, persistentbackgroundid string) {
	fmt.Println("Addded:" + backgroundid + ", " + filename + ", " + persistentbackgroundid)
	bkgtype := protocol.RecentBackground_SIMPLE
	if persistentbackgroundid != "" {
		bkgtype = protocol.RecentBackground_INSTANCED
	}

	cachedpreferences.Write(func(as protocol.AppSettings) protocol.AppSettings {
		id := fmt.Sprint(hashCode(backgroundid))
		as.RecentBackgroundStore.RecentBackgrounds[id] = &protocol.RecentBackground{
			BackgroundType: 		bkgtype,
			TimestampAdded: 		time.Now().UnixNano(),
		}

		if bkgtype == protocol.RecentBackground_INSTANCED {
			as.RecentBackgroundStore.RecentBackgrounds[id].Background = &protocol.RecentBackground_InstancedBackgroundID{
				InstancedBackgroundID: persistentbackgroundid,
			}
		} else {
			as.RecentBackgroundStore.RecentBackgrounds[id].Background = &protocol.RecentBackground_SimpleBackgroundID{
				SimpleBackgroundID: backgroundid,
			}
		}
		return as
	})
}

func removeRecentMediaBackground(backgroundid string) {
	cachedpreferences.Write(func(as protocol.AppSettings) protocol.AppSettings {
		delete(as.RecentBackgroundStore.RecentBackgrounds, backgroundid)
		return as
	})
}

func reencodeVideoFile(inputfile string, outputfile string, codec string, quality string, stripAudio bool, controlchannel chan string, request protocol.SimpleBackgroundRequest) {
	// This function is used to reencode video files with ffmpeg
	// inputfile: the path to the video file to reencode
	// outputfile: the path to the reencoded video file
	// codec: the codec to use for reencoding
	// quality: the quality setting for the codec aka CRF value
	// Example: reencodeVideoFile("input.mp4", "output.webm", "libvpx-vp9", "good")
	// The above example will reencode input.mp4 to output.webm using the libvpx-vp9 codec with the good quality setting

	process := exec.Command("tools\\ffmpeg\\bin\\ffmpeg.exe", "-i", inputfile, "-c:v", codec, "-crf", quality, "-b:v", "0", "-profile:v", "0", "-threads", "0", "-y")
	if stripAudio {
		process.Args = append(process.Args, "-an")
	}
	process.Args = append(process.Args, outputfile)
	donechannel := make(chan bool)
	go func() {
		select {
		case <-donechannel:
		case message := <-controlchannel:
			if message == "stop" {
				process.Process.Kill()
			} else {
				controlchannel <- message
			}
		}
	}()
	process.Stdout = os.Stdout
	process.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := process.Run()

	if err != nil {
		controlchannel <- "error," + err.Error()
	} else {
		addRecentMediaBackground(request.SimpleBackgroundID, filepath.Base(outputfile), "")
		controlchannel <- "done"
		donechannel <- true
	}

}

func reencodeImageFile(inputfile string, outputfile string, controlchannel chan string, request protocol.SimpleBackgroundRequest) {
	process := exec.Command("tools\\imagemagick\\magick.exe", inputfile, "-define", "png:compression-filter=5", "-define", "png:compression-level=9", "-define", "png:compression-strategy=2", outputfile)
	donechannel := make(chan bool)
	go func() {
		select {
		case <-donechannel:

		case message := <-controlchannel:
			if message == "stop" {
				process.Process.Kill()
			} else {
				controlchannel <- message
			}
		}
	}()
	process.Stdout = os.Stdout
	process.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := process.Run()
	if err != nil {
		controlchannel <- "error," + err.Error()
	} else {
		addRecentMediaBackground(request.SimpleBackgroundID, filepath.Base(outputfile), "")
		controlchannel <- "done"
		donechannel <- true
	}

}

func startReencodingSub(inputfile string, filename string, controlchannel chan string, specifictype string, request protocol.SimpleBackgroundRequest) (string, error) {
	switch specifictype {
	case "Video":
		outpath, err := filepath.Abs(path.Join("result", filename+".webm"))
		if err != nil {
			return "", err
		}
		go func() {
			reencodeVideoFile(inputfile, outpath, "libsvtav1", "23", true, controlchannel,request)
		}()
		return outpath, nil
	case "Image":
		outpath, err := filepath.Abs(path.Join("result", filename+".png"))
		if err != nil {
			return "", err
		}
		go func() {
			reencodeImageFile(inputfile, outpath, controlchannel,request)
		}()
		return outpath, nil
	default:
		go func() {
			controlchannel <- "askuser"
			specifiedtype := <-controlchannel
			if specifiedtype == "cancel" {
				return
			}
			outpath, error := startReencodingSub(inputfile, filename, controlchannel, specifiedtype,request)
			if error != nil {
				controlchannel <- "error," + error.Error()
			} else {
				controlchannel <- outpath
			}
		}()
		return "", nil
	}
}

func startReencoding(inputfile string, controlchannel chan string, filename string, request protocol.SimpleBackgroundRequest) (string, error) {
	filetype := FileNameToMediaType(inputfile)
	if _, ok := os.Stat("result"); os.IsNotExist(ok) {
		os.Mkdir("result", 0755)
	}
	return startReencodingSub(inputfile, filename, controlchannel, filetype,request)
}
