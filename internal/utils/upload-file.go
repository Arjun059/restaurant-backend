package utils;

import (
	"io"
	"fmt"
	"os"
	"context"
	"strings"
	"path/filepath"
)
import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)


func UploadFileToCloud(file io.Reader, fileName string, folder string) (string, error) {
	// Load Cloudinary credentials from environment variables
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return "", err
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		PublicID: fileName,
		Folder:   folder,
	})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

// DeleteFileFromCloud logs a deletion request. Implement real deletion
// using the Cloudinary SDK if desired.
// func DeleteFileFromCloud(folder string, storedName string) error {
// 	fmt.Printf("Delete requested for cloud file: %s (folder=%s)\n", storedName, folder)
// 	return nil
// }

func DeleteFileFromCloud(fileName, folder string) error {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return err
	}

	// remove extension if present
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	publicID := folder + "/" + fileName

	fmt.Println("publicID", publicID)

	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID:   publicID,
	})

	if err != nil {
		fmt.Printf("there is an error when delete image %+v", err)
	}

	return err
}
