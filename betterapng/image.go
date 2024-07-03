package betterapng

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strconv"
)

var pluginstarted = false

type BAPNGConfig struct {
	Width  int
	Height int
	DesiredFPS int
	NumberOfFrames uint64
	framesread uint64
}

type BAPNGFrame struct {
	FrameNumber uint64
	Codec string
}

type BAPNG struct {
	imageStream io.Reader
	imageWriter io.Writer
	config BAPNGConfig
}

func NewBAPNG(imageStream io.Reader) *BAPNG {
	return &BAPNG{imageStream,nil,BAPNGConfig{0,0,0,0,0}}
}

func NewBAPNGWriter(imageWriter io.Writer, config BAPNGConfig) *BAPNG {
	return &BAPNG{nil,imageWriter,config}
}

func decodeHeader(header []byte) (BAPNGConfig, error) {
	decoded := make(map[string]interface{})
	err := json.Unmarshal(header, &decoded)
	if err != nil {
		return BAPNGConfig{}, err
	}
	config := BAPNGConfig{}
	config.Width = int(decoded["Width"].(float64))
	config.Height = int(decoded["Height"].(float64))
	config.DesiredFPS = int(decoded["DesiredFPS"].(float64))
	config.NumberOfFrames = uint64(decoded["NumberOfFrames"].(float64))
	return config, nil
}

func (bapng *BAPNG) Open() (error) {
	// Read the protocol header

	header, err := readAllUntilByte(bapng.imageStream, 0x0A)
	if err != nil {
		return err
	}
	base64Header := string(header)
	decodedHeader, err := base64.StdEncoding.DecodeString(base64Header)
	if err != nil {
		return err
	}
	// Decode the header
	bapng.config, err = decodeHeader(decodedHeader)
	if err != nil {
		return err
	}
	return nil
}

func (bapng *BAPNG) WriteHeader() (error) {
	// Encode the header
	header, err := json.Marshal(bapng.config)
	if err != nil {
		return err
	}
	encodedHeader := base64.StdEncoding.EncodeToString(header)
	_, err = (bapng.imageWriter).Write([]byte(encodedHeader))
	if err != nil {
		return err
	}
	_, err = (bapng.imageWriter).Write([]byte("\x0A"))
	if err != nil {
		return err
	}
	return nil
}

func (bapng *BAPNG) WriteNextFrameAsPNG(img image.Image) (error) {
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, img)
	if err != nil {
		return err
	}
	config := BAPNGFrame{FrameNumber: uint64(bapng.config.NumberOfFrames), Codec: "PNG"}
	return bapng.WriteNextFrameRAW(buffer.Bytes(),config)
}

func (bapng *BAPNG) WriteNextFrameAsJPEG(img image.Image, quality int) (error) {
	buffer := new(bytes.Buffer)
	jpgeoptions := &jpeg.Options{
		Quality: quality,
	}
	err := EncodeJPEG(buffer, img, jpgeoptions)
	if err != nil {
		return err
	}
	config := BAPNGFrame{FrameNumber: uint64(bapng.config.NumberOfFrames), Codec: "JPEG"}
	return bapng.WriteNextFrameRAW(buffer.Bytes(),config)
}

func (bapng *BAPNG) WriteNextFrameRAW(image []byte, frameconfig BAPNGFrame) (error) {
	// Encode the image
	buffer := new(bytes.Buffer)
	// Write the frame config
	frameconfigjson, err := json.Marshal(frameconfig)
	if err != nil {
		return err
	}
	encodedFrameConfig := base64.StdEncoding.EncodeToString(frameconfigjson)
	_, err = buffer.Write([]byte(encodedFrameConfig))
	if err != nil {
		return err
	}
	_, err = buffer.Write([]byte("\x0A"))
	if err != nil {
		return err
	}
	_, err = buffer.Write(image)
	if err != nil {
		return err
	}
	// Write the image length
	_, err = (bapng.imageWriter).Write([]byte(strconv.FormatUint(uint64(buffer.Len()), 10)))
	if err != nil {
		return err
	}
	_, err = (bapng.imageWriter).Write([]byte("\x0A"))
	if err != nil {
		return err
	}
	// Write the image
	_, err = (bapng.imageWriter).Write(buffer.Bytes())
	if err != nil {
		return err
	}
	dup := bapng.config
	dup.NumberOfFrames = dup.NumberOfFrames + 1
	bapng.config = dup
	return nil
}

func (bapng *BAPNG) Close() (error) {
	if bapng.imageWriter != nil {
		(bapng.imageWriter).(*os.File).Close()
	}
	if bapng.imageStream != nil {
		(bapng.imageStream).(*os.File).Close()
	}
	return nil
}

func (bapng *BAPNG) ReadNextFrame() (image.Image, *BAPNGFrame, error) {
	image, config, err := bapng.ReadNextFrameRAW()
	if err != nil {
		return nil, config, err
	}
	img, err := UniversalDecoder(image, config.Codec)
	if err != nil {
		return nil, config, err
	}
	return img, config, nil
}

func (bapng *BAPNG) ReadNextFrameRAW() ([]byte,*BAPNGFrame, error) {
	// Read the image length
	if bapng.config.framesread >= uint64(bapng.config.NumberOfFrames) {
		return nil,nil, io.EOF
	}

	imageLength, err := readAllUntilByte(bapng.imageStream, 0x0A)
	if err != nil {
		return nil,nil, err
	}
	stringed := string(imageLength)
	number, _ := strconv.ParseUint(stringed, 10, 64)
	// Read the image
	imageData := make([]byte, number)
	_, err = (bapng.imageStream).Read(imageData)
	reader := bytes.NewReader(imageData)
	// Read the frame config
	frameConfig, err := readAllUntilByte(reader, 0x0A)
	decoded , err := base64.StdEncoding.DecodeString(string(frameConfig))
	parsed := BAPNGFrame{}
	err = json.Unmarshal(decoded, &parsed)
	if err != nil {
		return nil,nil, err
	}
	bapng.config.framesread++
	return imageData[len(frameConfig)+1:],&parsed, nil
}


func (bapng *BAPNG) GetConfig() (BAPNGConfig) {
	return bapng.config
}

func (bapng *BAPNG) GetFramesRead() (uint64) {
	return bapng.config.framesread
}

func (bapng *BAPNG) GetFramesLeft() (uint64) {
	return uint64(bapng.config.NumberOfFrames) - bapng.config.framesread
}

func (bapng *BAPNG) GetFPS() (int) {
	return bapng.config.DesiredFPS
}

func (bapng *BAPNG) GetWidth() (int) {
	return bapng.config.Width
}

func (bapng *BAPNG) GetHeight() (int) {
	return bapng.config.Height
}

func (bapng *BAPNG) GetNumberOfFrames() (uint64) {
	return bapng.config.NumberOfFrames
}

func (bapng *BAPNG) ReadAllFramesAsRAW() ([][]byte, []*BAPNGFrame, error) {
	frames := make([][]byte, bapng.config.NumberOfFrames)
	frameconfigs := make([]*BAPNGFrame, bapng.config.NumberOfFrames)
	for i := uint64(0); i < bapng.config.NumberOfFrames; i++ {
		frame, config, err := bapng.ReadNextFrameRAW()
		if err != nil {
			return frames, frameconfigs, err
		}
		frameconfigs[i] = config
		frames[i] = frame
	}
	return frames,frameconfigs, nil
}

func (bapng *BAPNG) ReadAllFrames() ([]image.Image, []*BAPNGFrame, error) {
	frames := make([]image.Image, bapng.config.NumberOfFrames)
	frameconfigs := make([]*BAPNGFrame, bapng.config.NumberOfFrames)
	for i := uint64(0); i < bapng.config.NumberOfFrames; i++ {
		frame, config, err := bapng.ReadNextFrame()
		if err != nil {
			return frames, frameconfigs, err
		}
		frames[i] = frame
		frameconfigs[i] = config
	}
	return frames, frameconfigs, nil
}



func (bapng *BAPNG) WriteAllFrames(frames []image.Image) (error) {
	for i := 0; i < len(frames); i++ {
		err := bapng.WriteNextFrameAsPNG(frames[i])
		if err != nil {
			return err
		}
	}
	return nil
}