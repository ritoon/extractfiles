package qrcode

import (
	"os"
	"testing"
)

func TestReadFileQrcode(t *testing.T) {
	data := map[string]string{
		"1000-0.jpg": "",
		"12.png":     "12",
		"2.png":      "2",
		"23-1.jpg":   "",
		"23.jpg":     "",
		"26-1.jpg":   "",
		"26.jpg":     "",
		"bureau.jpg": "",
		"bureau.png": "12",
		"henri_l_textured_background_resembling_vintage_parchment_engr_1b930817-ab64-47d8-8497-8f158604ffb3_3.png": "",
		"qr-timbre.png":            "",
		"test-special-qrcode.png":  "12",
		"test-special-qrcode2.png": "12",
		"timbre-qrcode-2.jpg":      "",
		"timbre-qrcode.jpg":        "",
	}
	files, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("no files found in testdata directory")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == "" || len(file.Name()) < 5 {
			continue
		}
		lastElement := file.Name()[len(file.Name())-4:]
		if lastElement != ".jpg" && lastElement != ".png" && lastElement != ".jpeg" &&
			lastElement != ".JPG" && lastElement != ".PNG" && lastElement != ".JPEG" {
			continue
		}
		fileRes, ok := data[file.Name()]
		if !ok {
			t.Fatalf("file %s not found in data map", file.Name())
		}

		file, err := os.Open("testdata/" + file.Name())
		if err != nil {
			t.Fatalf("failed to open file: %v", err)
		}
		defer file.Close()

		res, err := ReadFileQrcode(file)
		if err != nil && err != ErrNotAQRCode {
			t.Fatalf("failed to read file: %v", err)
		}
		if res != fileRes {
			t.Fatalf("expected %s, got %s", fileRes, res)
		}

	}
}
