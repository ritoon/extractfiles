package qrcode

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"

	"roumet/logger"
)

var (
	// ErrNotAQRCode is returned when the image is not a QR code.
	ErrNotAQRCode = fmt.Errorf("qrcode: not a QR code")
)

func ReadFileQrcode(file *os.File) (string, error) {

	img, _, err := image.Decode(file)
	if err != nil {
		logger.Error("qrcode: decode error:", err)
		return "", err
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		logger.Error("qrcode: decode create binary:", err)
		return "", err
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		logger.Error("qrcode: read value", err)
		return "", ErrNotAQRCode
	}

	return fmt.Sprint(result), nil
}
