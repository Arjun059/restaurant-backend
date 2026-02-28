package utils

// import (
// 	"bytes"
// 	"fmt"
// 	"image"
// 	imagepng "image/png"

// 	"github.com/disintegration/imaging"
// 	qrcode "github.com/skip2/go-qrcode"
// )

// var qrCodeFolder = "qr-codes"

// func GenerateQrAndUploadToCloud(restaurantURL string, restaurantURLPath string) (string, error) {

// 	// 1️⃣ Generate QR PNG bytes
// 	qrBytes, err := qrcode.Encode(restaurantURL, qrcode.Medium, 512)
// 	if err != nil {
// 		return "", err
// 	}

// 	// 2️⃣ Decode QR bytes into image.Image
// 	qrImage, err := imaging.Decode(bytes.NewReader(qrBytes))
// 	if err != nil {
// 		return "", err
// 	}

// 	// 3️⃣ Open Poster Background
// 	posterPath := "public/assets/qr-background.png"
// 	posterImg, err := imaging.Open(posterPath)
// 	if err != nil {
// 		return "", err
// 	}

// 	// 4️⃣ Resize QR
// 	qrSize := 566
// 	qrResized := imaging.Resize(qrImage, qrSize, qrSize, imaging.Lanczos)

// 	// 5️⃣ Calculate center
// 	bounds := posterImg.Bounds()
// 	centerX := (bounds.Dx() - qrSize) / 2
// 	centerY := (bounds.Dy() - qrSize) / 2 - 42
// 	// centerY := 550 // adjust as needed

// 	// 6️⃣ Overlay QR onto poster
// 	finalImage := imaging.Overlay(
// 		posterImg,
// 		qrResized,
// 		image.Pt(centerX, centerY),
// 		1.0,
// 	)

// 	// 7️⃣ Encode final image to memory
// 	var buf bytes.Buffer
// 	err = imagepng.Encode(&buf, finalImage)
// 	if err != nil {
// 		return "", err
// 	}

// 	// 8️⃣ Upload FINAL poster (not raw QR)
// 	filename := restaurantURLPath + ".png"

// 	uploadedURL, err := UploadFileToCloud(
// 		bytes.NewReader(buf.Bytes()),
// 		filename,
// 		qrCodeFolder,
// 	)
// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Println("Poster uploaded to:", uploadedURL)
// 	return uploadedURL, nil
// }


import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	imagepng "image/png"
	"io/ioutil"

	"github.com/disintegration/imaging"
	qrcode "github.com/skip2/go-qrcode"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var qrCodeFolder = "qr-codes"

func GenerateQrAndUploadToCloud(
	restaurantURL string,
	restaurantURLPath string,
	restaurantName string, // 👈 NEW PARAM
) (string, error) {

	// 1️⃣ Generate QR
	// qrBytes, err := qrcode.Encode(restaurantURL, qrcode.Medium, 512)
	// if err != nil {
	// 	return "", err
	// }


	qr, err := qrcode.New(restaurantURL, qrcode.Medium)
	if err != nil {
		return "", err
	}

	// 🔥 Control padding here
	qr.DisableBorder = true

	qrBytes, err := qr.PNG(512)
	if err != nil {
		return "", err
	}

	qrImage, err := imaging.Decode(bytes.NewReader(qrBytes))
	if err != nil {
		return "", err
	}

	// 2️⃣ Open Poster
	posterPath := "public/assets/qr-background.png"
	posterImg, err := imaging.Open(posterPath)
	if err != nil {
		return "", err
	}

	// 3️⃣ Resize QR
	qrSize := 480
	qrResized := imaging.Resize(qrImage, qrSize, qrSize, imaging.Lanczos)

	// 4️⃣ Center QR
	bounds := posterImg.Bounds()
	centerX := (bounds.Dx() - qrSize) / 2
	centerY := (bounds.Dy() - qrSize) / 2 - 20

	finalImage := imaging.Overlay(
		posterImg,
		qrResized,
		image.Pt(centerX, centerY),
		1.0,
	)

	// ===============================
	// 🆕 5️⃣ DRAW RESTAURANT NAME
	// ===============================

	fontBytes, err := ioutil.ReadFile("public/fonts/PlayfairDisplay-Bold.ttf")
	if err != nil {
		return "", err
	}

	ft, err := opentype.Parse(fontBytes)
	if err != nil {
		return "", err
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    72,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return "", err
	}

	// Create RGBA to draw text
	rgba := imaging.Clone(finalImage)

	col := color.RGBA{34, 34, 34, 255}

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: face,
	}

	// Center text horizontally
	textWidth := d.MeasureString(restaurantName).Round()
	x := (bounds.Dx() - textWidth) / 2
	y := 300 // adjust vertical position as needed

	d.Dot = fixed.P(x, y)
	d.DrawString(restaurantName)

	finalImage = rgba

	// ===============================

	var buf bytes.Buffer
	err = imagepng.Encode(&buf, finalImage)
	if err != nil {
		return "", err
	}

	filename := restaurantURLPath

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