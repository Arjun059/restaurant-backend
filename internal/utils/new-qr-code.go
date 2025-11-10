package utils

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"

	qrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)


func GenerateHalftoneQRAndUpload(restaurantURL, restaurantURLPath string) (string, error) {
	qrc, err := qrcode.NewWith(restaurantURL, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest))
	if err != nil {
		return "", fmt.Errorf("failed to init qr: %w", err)
	}

	tmpFile := "./tmp/images/tmp-halftone-qr.png"
	halftoneImgPath := "./public/assets/logo1.png"

	// fgColor := color.RGBA{R: 255, G: 200, B: 55, A: 255} // Orange/Yellow fallback

	// optional gradient values for future use (not used by WithFgColor)
	// gradientColors := []color.RGBA{
	// 	{255, 200, 55, 255},  // Orange/Yellow
	// 	{255, 42, 104, 255},  // Pink/Red
	// 	{131, 58, 180, 255},  // Purple
	// }

	// use a single foreground color (fallback for gradient)
	options := []standard.ImageOption{
		standard.WithHalftone(halftoneImgPath),
		// standard.WithQRWidth(21),
		// standard.WithFgColor(fgColor),
		// standard.WithBgColor(color.White),
	}

	w0, err := standard.New(tmpFile, options...)
	if err != nil {
		return "", fmt.Errorf("writer init failed: %w", err)
	}

	if err := qrc.Save(w0); err != nil {
		return "", fmt.Errorf("QR save failed: %w", err)
	}

	f, err := os.Open(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to open tmp file: %w", err)
	}
	defer f.Close()

	uploadedURL, err := UploadFileToCloud(f, restaurantURLPath, qrCodeFolder)
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	// leave file during testing,
	// enable below in prod if folder grows too large
	// _ = os.Remove(tmpFile)

	return uploadedURL, nil
}
