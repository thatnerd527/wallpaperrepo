package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
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

	process := exec.Command("ffmpeg", "-i", inputfile, "-c:v", codec, "-crf", quality, "-b:v", "0", "-profile:v", "0", "-threads", "0")
	if stripAudio {
		process.Args = append(process.Args, "-an")
	}
	process.Args = append(process.Args, outputfile)
	donechannel := make(chan bool)
	go func() {
		select {
		case <-donechannel:
			controlchannel <- "done"
		case message := <-controlchannel:
			if message == "stop" {
				process.Process.Kill()
			}
		}
	}()
	process.Run()
	donechannel <- true

}

func reencodeImageFile(inputfile string, outputfile string, controlchannel chan string) {
	process := exec.Command("magick",inputfile,"-define","png:compression-filter=5","-define","png:compression-level=9","-define","png:compression-strategy=2",outputfile)
	donechannel := make(chan bool)
	go func() {
		select {
		case <-donechannel:
			controlchannel <- "done"
		case message := <-controlchannel:
			if message == "stop" {
				process.Process.Kill()
			}
		}
	}()
	process.Run()
	donechannel <- true

}

func startReencodingSub(inputfile string, controlchannel chan string, specifictype string) (string, error) {
	switch specifictype {
		case "Video":
			outpath, err := filepath.Abs(path.Join("result", "reencoded.webm"))
			if err != nil {
				return "", err
			}
			go func() {
				reencodeVideoFile(inputfile, outpath, "libaom-av1", "23", true, controlchannel)
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
				startReencodingSub(inputfile, controlchannel, specifiedtype)
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
