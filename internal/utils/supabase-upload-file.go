package utils;

import (
	"bytes"
	"net/http"
	"io"
	"fmt"
	"mime/multipart"
	"os"
	"context"
)
import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)



func UploadFileToCloud_Backup(file multipart.File, fileName string, contentType string) error {
	// Read file into bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	var (
		SUPABASE_URL     = os.Getenv("SUPABASE_STORAGE_URL")
		SUPABASE_BUCKET  = os.Getenv("SUPABASE_BUCKET") // your bucket name
		SUPABASE_API_KEY = os.Getenv("SUPABASE_SERVICE_ROLE_KEY") // get this from Supabase → API settings
	)
	
	// for _, env := range os.Environ() {
	// 	fmt.Println(env, "env")
	// }
	// fmt.Printf("URL: %s, BUCKET: %s, API_KEY: %s", SUPABASE_URL, SUPABASE_BUCKET , SUPABASE_API_KEY )

	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", SUPABASE_URL, SUPABASE_BUCKET, fileName)
	// url := fmt.Sprintf("%s/storage/v1/object/%s/public/%s", SUPABASE_URL, SUPABASE_BUCKET, fileName)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}

	fmt.Printf("url : %s", url)

	req.Header.Set("Authorization", "Bearer "+SUPABASE_API_KEY)
	req.Header.Set("apikey", SUPABASE_API_KEY)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error on client req %v", err )
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("response body %s", body)
	return nil
	// return fmt.Errorf("Supabase upload failed: %s", body)
}


func UploadFileToCloud(file multipart.File, fileName string, folder string) (string, error) {
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