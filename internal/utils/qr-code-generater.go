package utils

import (
	"bytes"
	"fmt"
	"image"
	imagepng "image/png"

	"github.com/disintegration/imaging"
	qrcode "github.com/skip2/go-qrcode"
)

var qrCodeFolder = "qr-codes"

func GenerateQrAndUploadToCloud(restaurantURL string, restaurantURLPath string) (string, error) {

	// 1️⃣ Generate QR PNG bytes
	qrBytes, err := qrcode.Encode(restaurantURL, qrcode.Medium, 512)
	if err != nil {
		return "", err
	}

	// 2️⃣ Decode QR bytes into image.Image
	qrImage, err := imaging.Decode(bytes.NewReader(qrBytes))
	if err != nil {
		return "", err
	}

	// 3️⃣ Open Poster Background
	posterPath := "public/assets/qr-background.png"
	posterImg, err := imaging.Open(posterPath)
	if err != nil {
		return "", err
	}

	// 4️⃣ Resize QR
	qrSize := 566
	qrResized := imaging.Resize(qrImage, qrSize, qrSize, imaging.Lanczos)

	// 5️⃣ Calculate center
	bounds := posterImg.Bounds()
	centerX := (bounds.Dx() - qrSize) / 2
	centerY := (bounds.Dy() - qrSize) / 2 - 42
	// centerY := 550 // adjust as needed

	// 6️⃣ Overlay QR onto poster
	finalImage := imaging.Overlay(
		posterImg,
		qrResized,
		image.Pt(centerX, centerY),
		1.0,
	)

	// 7️⃣ Encode final image to memory
	var buf bytes.Buffer
	err = imagepng.Encode(&buf, finalImage)
	if err != nil {
		return "", err
	}

	// 8️⃣ Upload FINAL poster (not raw QR)
	filename := restaurantURLPath + ".png"

	uploadedURL, err := UploadFileToCloud(
		bytes.NewReader(buf.Bytes()),
		filename,
		qrCodeFolder,
	)
	if err != nil {
		return "", err
	}

	fmt.Println("Poster uploaded to:", uploadedURL)
	return uploadedURL, nil
}