package utils

import (
	qrcode "github.com/skip2/go-qrcode"
	"bytes"
	"log"
	"fmt"
)

var qrCodeFolder = "qr-codes"

func GenerateQrAndUploadToCloud(fileName string) (string, error) {
	// var qrCodeFileName = "abc.png"
	// err := qrcode.WriteFile("https://example.com", qrcode.Medium, 256, qrCodeFileName)
	// if err != nil {
	// 	return "", err
	// }

	// Step 1: Generate QR code as PNG byte slice
	png, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)
	if err != nil {
		log.Fatal("Error generating QR code:", err)
	}

	// Step 2: Wrap it in bytes.Reader
	pngReader := bytes.NewReader(png)

		// Step 3: Upload
	url, err := UploadFileToCloud(pngReader, "my_qr_code", qrCodeFolder)
	if err != nil {
		log.Fatal("Upload failed:", err)
		return "", err
	}
	
	fmt.Println("QR code uploaded to:", url)
	return url, nil
}

