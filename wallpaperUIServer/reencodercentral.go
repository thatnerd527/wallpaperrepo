package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

var encodertochan = make(map[string]chan string)

func reencodeVideoFile(inputfile string, outputfile string, codec string, quality string, stripAudio bool, controlchannel chan string) {
	// This function is used to reencode video files with ffmpeg
	// inputfile: the path to the video file to reencode
	// outputfile: the path to the reencoded video file
	// codec: the codec to use for reencoding
	// quality: the quality setting for the codec aka CRF value
	// Example: reencodeVideoFile("input.mp4", "output.webm", "libvpx-vp9", "good")
	// The above example will reencode input.mp4 to output.webm using the libvpx-vp9 codec with the good quality setting

	process := exec.Command("tools\\ffmpeg\\bin\\ffmpeg.exe", "-i", inputfile, "-c:v", codec, "-crf", quality, "-b:v", "0", "-profile:v", "0", "-threads", "0","-y")
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
		controlchannel <- "done"
		donechannel <- true
	}

}

func reencodeImageFile(inputfile string, outputfile string, controlchannel chan string) {
	process := exec.Command("tools\\imagemagick\\magick.exe",inputfile,"-define","png:compression-filter=5","-define","png:compression-level=9","-define","png:compression-strategy=2",outputfile)
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
	err := process.Run()
	if err != nil {
		controlchannel <- "error," + err.Error()
	} else {
		controlchannel <- "done"
		donechannel <- true
	}

}

func startReencodingSub(inputfile string, controlchannel chan string, specifictype string) (string, error) {
	switch specifictype {
		case "Video":
			outpath, err := filepath.Abs(path.Join("result", "reencoded.webm"))
			if err != nil {
				return "", err
			}
			go func() {
				reencodeVideoFile(inputfile, outpath, "libsvtav1", "23", true, controlchannel)
			}()
			return outpath, nil
		case "Image":
			outpath, err := filepath.Abs(path.Join("result", "reencoded.png"))
			if err != nil {
				return "", err
			}
			go func() {
				reencodeImageFile(inputfile, outpath,controlchannel)
			}()
			return outpath, nil
		default:
			go func() {
				controlchannel <- "askuser"
				specifiedtype := <-controlchannel
				if (specifiedtype == "cancel") {
					return
				}
				outpath, error := startReencodingSub(inputfile, controlchannel, specifiedtype)
				if error != nil {
					controlchannel <- "error," + error.Error()
				} else {
					controlchannel <- outpath
				}
			}()
			return "", nil
	}
}

func startReencoding(inputfile string, controlchannel chan string) (string, error) {
	filetype := FileNameToMediaType(inputfile)
	if _, ok := os.Stat("result"); os.IsNotExist(ok) {
		os.Mkdir("result", 0755)
	}
	return startReencodingSub(inputfile, controlchannel, filetype)
}
