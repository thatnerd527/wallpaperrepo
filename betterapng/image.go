package betterapng

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"
)

type BAPNGConfig struct {
	Width  int
	Height int
	DesiredFPS int
	NumberOfFrames int
	framesread uint64
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
	config.NumberOfFrames = int(decoded["NumberOfFrames"].(float64))
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

func (bapng *BAPNG) WriteNextFrame(img image.Image) (error) {
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, img)
	if err != nil {
		return err
	}
	return bapng.WriteNextFrameRAW(buffer.Bytes())
}

func (bapng *BAPNG) WriteNextFrameRAW(image []byte) (error) {
	// Encode the image
	buffer := new(bytes.Buffer)
	_, err := buffer.Write(image)
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

func (bapng *BAPNG) ReadNextFrame() (image.Image, error) {
	image, err := bapng.ReadNextFrameAsPNG()
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(image))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (bapng *BAPNG) ReadNextFrameAsPNG() ([]byte, error) {
	// Read the image length
	if bapng.config.framesread >= uint64(bapng.config.NumberOfFrames) {
		return nil, io.EOF
	}

	imageLength, err := readAllUntilByte(bapng.imageStream, 0x0A)
	if err != nil {
		return nil, err
	}
	stringed := string(imageLength)
	number, _ := strconv.ParseUint(stringed, 10, 64)
	// Read the image
	imageData := make([]byte, number)
	_, err = (bapng.imageStream).Read(imageData)
	if err != nil {
		return nil, err
	}
	bapng.config.framesread++
	return imageData, nil
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

func (bapng *BAPNG) GetNumberOfFrames() (int) {
	return bapng.config.NumberOfFrames
}

func (bapng *BAPNG) ReadAllFramesAsPNG() ([][]byte, error) {
	frames := make([][]byte, bapng.config.NumberOfFrames)
	for i := 0; i < bapng.config.NumberOfFrames; i++ {
		frame, err := bapng.ReadNextFrameAsPNG()
		if err != nil {
			return frames, err
		}
		frames[i] = frame
	}
	return frames, nil
}

func (bapng *BAPNG) ReadAllFrames() ([]image.Image, error) {
	frames := make([]image.Image, bapng.config.NumberOfFrames)
	for i := 0; i < bapng.config.NumberOfFrames; i++ {
		frame, err := bapng.ReadNextFrame()
		if err != nil {
			return frames, err
		}
		frames[i] = frame
	}
	return frames, nil
}



func (bapng *BAPNG) WriteAllFrames(frames []image.Image) (error) {
	for i := 0; i < len(frames); i++ {
		err := bapng.WriteNextFrame(frames[i])
		if err != nil {
			return err
		}
	}
	return nil
}