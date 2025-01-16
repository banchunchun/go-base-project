package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
)

func ToJpeg(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/png":
		img, err := png.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, nil); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("unable to convert %#v to jpeg", contentType)
}

func ConvertPNG2JPG(pngFileName string, jpgFileName string) error {
	imageBytes, err := os.ReadFile(pngFileName)
	if err != nil {
		return err
	}

	jpegBytes, err := ToJpeg(imageBytes)

	if err != nil {
		return err
	}

	err = DumpBytes(jpgFileName, jpegBytes)
	return err
}

func GetImgBase64(path string) (baseImg string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	var base64Encoding string

	imgByte, _ := io.ReadAll(file)
	mimeType := http.DetectContentType(imgByte)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += base64.StdEncoding.EncodeToString(imgByte)
	return base64Encoding, nil
}
