package betterapng

import (
	"bytes"
	"io"
)

func readAllUntilByte(reader io.Reader, stopByte byte) ([]byte, error) {
	var buffer bytes.Buffer
	var bufferByte []byte = make([]byte, 1)
	for {
		_, err := reader.Read(bufferByte)
		if err != nil {
			return nil, err
		}
		if bufferByte[0] == stopByte {
			break
		}
		buffer.Write(bufferByte)
	}
	return buffer.Bytes(), nil
}