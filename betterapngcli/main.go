package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/thatnerd/betterapng"
)

type NameAndIndex struct {
	Name string
	Index int
}

func sortFiles(files []string) []string {
	// First step, is to remove everything that is not a number from the filename
	// Second step, is to convert the filename to an integer
	// Third step, is to sort the files

	// First step and second step
	var sortedFiles []NameAndIndex
	for _, file := range files {
		// Remove everything that is not a number from the filename
		var filename string
		for _, char := range file {
			if char >= '0' && char <= '9' {
				filename += string(char)
			}
		}
		// Convert the filename to an integer
		index, err := strconv.Atoi(filename)
		if err != nil {
			continue
		}
		sortedFiles = append(sortedFiles, NameAndIndex{file, index})
	}

	// Third step
	var sorted []string
	for i := 0; i < len(sortedFiles); i++ {
		for j := 0; j < len(sortedFiles); j++ {
			if sortedFiles[j].Index == i {
				sorted = append(sorted, sortedFiles[j].Name)
			}
		}
	}
	return sorted
}

func nameToCodec(name string) string {
	if (strings.HasSuffix(strings.ToLower(name), ".png")) {
		return "PNG"
	}
	if (strings.HasSuffix(strings.ToLower(name), ".avif")) {
		return "AVIF"
	}
	if (strings.HasSuffix(strings.ToLower(name), ".webp")) {
		return "WEBP"
	}
	if (strings.HasSuffix(strings.ToLower(name), ".jpg") || strings.HasSuffix(strings.ToLower(name), ".jpeg")) {
		return "JPEG"
	}
	return "UNKNOWN"
}

func readFolder(folder string) ([]string, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, file := range files {
		paths = append(paths, folder+"/"+file.Name())
	}
	return paths, nil
}

func main() {
	// Operation mode
	// 	encode: Encode a series of images into a BAPNG
	// 	decode: Decode a BAPNG into a series of images
	mode := flag.String("mode", "encode", "Operation mode (encode / decode)")
	input := flag.String("input", "", "Input folder / file")
	output := flag.String("output", "", "Output folder / file")
	fps := flag.Uint64("fps", 30, "Frames per second")
	flag.Parse()

	if *mode == "encode" {
		// Read the folder
		files, err := readFolder(*input)
		if err != nil {
			panic(err)
		}
		if len(files) == 0 {
			panic("No files found in the folder")
		}
		// Sort the files
		files = sortFiles(files)
		output, err := os.Create(*output)
		if err != nil {
			panic(err)
		}
		firstfile, err := os.Open(files[0])
		if err != nil {
				panic(err)
			}
		img, err := png.Decode(firstfile)
		if err != nil {
			panic(err)
		}


		config := betterapng.BAPNGConfig{
			Width: img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
			DesiredFPS: int(*fps),
			NumberOfFrames: uint64(len(files)),
		}
		firstfile.Close()
		// Create a new BAPNG
		bapng := betterapng.NewBAPNGWriter(output, config)
		bapng.WriteHeader()

		// Add the images to the BAPNG
		for i, filen := range files {
			file, err := os.ReadFile(filen)
			if err != nil {
				panic(err)
			}
			fc := betterapng.BAPNGFrame{
				FrameNumber: uint64(i),
				Codec: nameToCodec(filen),
			}
			err = bapng.WriteNextFrameRAW(file,fc)
			if err != nil {
				panic(err)
			}
		}
		// Close the BAPNG
		err = bapng.Close()
		if err != nil {
			panic(err)
		}
	} else if *mode == "decode" {
		// Read the BAPNG
		fmt.Println("Decoding BAPNG")
		file, err := os.Open(*input)
		bapng := betterapng.NewBAPNG(file)
		bapng.Open()
		if err != nil {
			panic(err)
		}
		// Read all the frames
		frames, configs, err := bapng.ReadAllFrames()
		fmt.Println("Decoded", len(frames), "frames")
		if err != nil {
			panic(err)
		}
		// Write the frames to the output folder
		for i, frame := range frames {
			file, err := os.Create(*output + "/frame" + strconv.Itoa(i) + "." + configs[i].Codec)
			png.Encode(file, frame)
			if err != nil {
				panic(err)
			}
		}
	}
}