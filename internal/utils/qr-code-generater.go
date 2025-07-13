package utils

import (
	qrcode "github.com/skip2/go-qrcode"
	"bytes"
	"log"
	"fmt"
	"github.com/google/uuid"
)

var qrCodeFolder = "qr-codes"

func GenerateQrAndUploadToCloud(restaurantURL string, restaurantID uint) (string, error) {
	// var qrCodeFileName = "abc.png"
	// err := qrcode.WriteFile("https://example.com", qrcode.Medium, 256, qrCodeFileName)
	// if err != nil {
	// 	return "", err
	// }

	// Step 1: Generate QR code as PNG byte slice
	png, err := qrcode.Encode(restaurantURL, qrcode.Medium, 256)
	if err != nil {
		log.Fatal("Error generating QR code:", err)
	}

	// Step 2: Wrap it in bytes.Reader
	pngReader := bytes.NewReader(png)

	restaurantIDString := fmt.Sprintf("%d", restaurantID)

	filename := fmt.Sprintf("%s_%s", restaurantIDString, uuid.New().String())

		// Step 3: Upload
	uploadedURL, err := UploadFileToCloud(pngReader, filename, qrCodeFolder)
	if err != nil {
		log.Fatal("Upload failed:", err)
		return "", err
	}
	
	fmt.Println("QR code uploaded to:", uploadedURL)
	return uploadedURL, nil
}

