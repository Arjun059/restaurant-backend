package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"restaurant/internal/models"
	"io"
	"net/http"
	utils "restaurant/internal/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/datatypes"

	"os"
	"path/filepath"
	"github.com/gorilla/schema"

	"strings"
	"github.com/google/uuid"
	"mime/multipart"
)

type DishHandler struct {
	DB *gorm.DB;
}

func (dh *DishHandler) AddDish(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	authContext := utils.GetAuthContext(r)
	// fmt.Printf("authContext %+v", authContext)

	fmt.Printf("Content-Type of request header: %s\n", r.Header.Get("Content-Type"))

	// Parse multipart form with a 50MB limit
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		fmt.Println("Error parsing request:", err)
		http.Error(w, "Error parsing multipart form data", http.StatusBadRequest)
		return
	}

 filesMeta := []models.FileMeta{}

	// Handle uploaded files
files := r.MultipartForm.File["images"]
for  _, fileHeader := range files {
	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, "Failed to open file: "+fileHeader.Filename, http.StatusBadRequest)
		return
	}

	func(file multipart.File) {
		defer file.Close()

		ext := filepath.Ext(fileHeader.Filename)
		name := strings.TrimSuffix(fileHeader.Filename, ext)

		// which folder to upload this assets 
		// every restaurant have their own folder for image
		folder := authContext.RestaurantID.String()

		uploadedFileName := fmt.Sprintf("%s_%s", name, uuid.New().String())

		fmt.Println("Before file upload:", uploadedFileName)

		// Call cloud uploader correctly
		uploadedURL, err := utils.UploadFileToCloud(file, uploadedFileName, folder)
		if err != nil {
			fmt.Printf("Failed to upload to Supabase: %v\n", err)
			http.Error(w, "Failed to upload to Supabase: "+err.Error(), http.StatusInternalServerError)
			return
		}

		filesMeta = append(filesMeta, models.FileMeta{
			Name: fmt.Sprintf("%s%s", uploadedFileName, ext), // add file ext when save to db
			Folder: folder,
			Url: uploadedURL,
		})

		fmt.Println("Uploaded URL:", uploadedURL)
	}(file)
}

	// Handle categories field (JSON encode string array)
	categories := r.Form["categories"]
	categoryJSON, err := json.Marshal(categories)
	if err != nil {
		fmt.Println("Error marshaling categories:", err)
		http.Error(w, "Failed to process categories", http.StatusBadRequest)
		return
	} 
	// Remove categories from the form so it doesn't interfere with decoding
  delete(r.MultipartForm.Value, "categories")

	// Decode remaining fields into Dish model
	var body models.Dish
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(&body, r.MultipartForm.Value); err != nil {
		fmt.Printf("Error decoding form: %v\n", err)
		http.Error(w, "Error decoding form data", http.StatusBadRequest)
		return
	}

	// Validate the dish object
	if err := utils.ValidateDish(body); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		utils.WriteErrorResponse(w, fmt.Sprintf("Validation Failed: %v", err), http.StatusBadRequest)
		return
	}

	imagesJson, err := json.Marshal(filesMeta)
	if err != nil {
		print("error in parsing email to json")
	} else {
		body.Images = datatypes.JSON(imagesJson)
	}
	// Assign category JSON
	body.Categories = datatypes.JSON(categoryJSON)
	
	// TODO: RestaurantID for test only remove after test
	body.RestaurantID = authContext.RestaurantID;
	body.UserID = authContext.UserID;

	// Save to DB
	if err := dh.DB.Create(&body).Error; err != nil {
		fmt.Printf("Database error: %v\n", err)
		utils.WriteErrorResponse(w, "Error occurred while saving dish", http.StatusBadRequest)
		return
	}

	utils.WriteSuccessResponse(w, "Dish added successfully", http.StatusOK, body)
}

func (dh *DishHandler) UpdateDish(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var rVars = mux.Vars(r)
	id := rVars["id"]

	restInfo := utils.GetAuthContext(r)
	var body models.Dish

	// Load current dish from DB to compare existing images later
	var currentDish models.Dish
	if err := dh.DB.Where("id = ?", id).Where("restaurant_id = ?", restInfo.RestaurantID).First(&currentDish).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "dish not found", http.StatusBadRequest)
			return
		}
		fmt.Printf("error fetching current dish: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}


	// Parse multipart form (allow up to 50MB like in AddDish)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		fmt.Println("Error parsing multipart form:", err)
		http.Error(w, "Error parsing multipart form data", http.StatusBadRequest)
		return
	}

	// Handle categories (array) similar to AddDish
	categories := r.Form["categories"]
	if len(categories) > 0 {
		categoryJSON, err := json.Marshal(categories)
		if err == nil {
			body.Categories = datatypes.JSON(categoryJSON)
		}
	}

	// ImageInfo describes stored/uploaded image metadata
	type ImageInfo struct {
		Url    string `json:"url"`
		Folder string `json:"folder"`
		Name   string `json:"name"`
	}

	// Handle existingImages field (client may send JSON array/object or plain strings)
	existingImagesSlice := r.Form["existingImages"]
	existing := make([]ImageInfo, 0)
	if len(existingImagesSlice) > 0 {
		for _, v := range existingImagesSlice {
			// v might be a JSON array, a JSON object, or a plain string
			// Try JSON array first
			var arr []ImageInfo
			if err := json.Unmarshal([]byte(v), &arr); err == nil {
				existing = append(existing, arr...)
				continue
			}

			// Try single JSON object
			var obj ImageInfo
			if err := json.Unmarshal([]byte(v), &obj); err == nil {
				existing = append(existing, obj)
				continue
			}

			// Fallback: treat as plain URL string
			existing = append(existing, ImageInfo{Url: v})
		}
	}

	// Remove categories and uploaded images from the form values so decoder doesn't choke
	if r.MultipartForm != nil {
		delete(r.MultipartForm.Value, "categories")
		delete(r.MultipartForm.Value, "existingImages")
	}

	// Decode remaining form values into Dish struct
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(&body, r.MultipartForm.Value); err != nil {
		fmt.Printf("Error decoding form: %v\n", err)
		http.Error(w, "Error decoding form data", http.StatusBadRequest)
		return
	}

	// Handle uploaded images if present
	filesMeta := []models.FileMeta{}
	if r.MultipartForm != nil {
		files := r.MultipartForm.File["images"]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "Failed to open file: "+fileHeader.Filename, http.StatusBadRequest)
				return
			}

			func(file multipart.File, fileHeader *multipart.FileHeader) {
				defer file.Close()

				ext := filepath.Ext(fileHeader.Filename)
				name := strings.TrimSuffix(fileHeader.Filename, ext)
				folder := restInfo.RestaurantID.String()
				uploadedFileName := fmt.Sprintf("%s_%s", name, uuid.New().String())

				uploadedURL, err := utils.UploadFileToCloud(file, uploadedFileName, folder)
				if err != nil {
					fmt.Printf("Failed to upload to Supabase: %v\n", err)
					// don't abort entire request for an upload error; surface as server error
					http.Error(w, "Failed to upload to Supabase: "+err.Error(), http.StatusInternalServerError)
					return
				}

				filesMeta = append(filesMeta, models.FileMeta{
					Name:   fmt.Sprintf("%s%s", uploadedFileName, ext),
					Folder: folder,
					Url:    uploadedURL,
				})
			}(file, fileHeader)
		}
	}

	var incoming []ImageInfo
	if len(filesMeta) > 0 {
		if imagesJson, err := json.Marshal(filesMeta); err == nil {
		
			if err := json.Unmarshal(imagesJson, &incoming); err != nil {
				fmt.Printf("Error unmarshaling incoming images: %v\n", err)
				http.Error(w, "Error processing uploaded images", http.StatusBadRequest)
				return
			}
		}
	}

	merged := append(existing, incoming...)
	mergedJSON, err := json.Marshal(merged)

	if err != nil {
		fmt.Printf("Error marshaling merged images: %v\n", err)
		http.Error(w, "Error processing images", http.StatusInternalServerError)
		return
	}

	fmt.Println("Total images len", len(merged))
	body.Images = datatypes.JSON(mergedJSON)

		// Delete images that exist in DB but are not present in the
		// client-provided `existing` list. Build sets of retained URLs
		retained := map[string]bool{}
		for _, m := range existing {
				retained[m.Name] = true
		}


		var currentFiles []models.FileMeta
		if len(currentDish.Images) > 0 {
			if err := json.Unmarshal(currentDish.Images, &currentFiles); err != nil {
				fmt.Printf("Error unmarshaling current dish images: %v\n", err)
			}
		}

		fmt.Printf("this is current file %+v", currentFiles)

		for _, cf := range currentFiles {
			// decide whether this file should be deleted: it's deleted
			// only when it's not present in the client's existing list
			// (neither by URL nor by name)
			kept := false
			if cf.Url != "" && retained[cf.Name] {
				kept = true
			}
			if !kept {
				fmt.Printf("this file is going to delete %+v", cf)
				if err := utils.DeleteFileFromCloud(cf.Name, cf.Folder); err != nil {
					fmt.Printf("failed to delete cloud file %s/%s: %v\n",cf.Name, cf.Folder,  err)
				} else {
					fmt.Printf("deleted cloud file %s/%s\n", cf.Name, cf.Folder)
				}
			}
		}

	// Ensure restaurant scoping and user info are preserved (don't allow changing owner)
	body.RestaurantID = restInfo.RestaurantID
	body.UserID = restInfo.UserID

	// Perform update scoped to the restaurant and id
	if err := dh.DB.
		Model(&models.Dish{}).
		Where("id = ?", id).
		Where("restaurant_id = ?", restInfo.RestaurantID).
		Updates(&body).Error; err != nil {
		fmt.Println("update error:", err)
		http.Error(w, "Error Occur", http.StatusInternalServerError)
		return
	}

	var updatedProduct models.Dish
	// Use pointer and scope retrieval to restaurant
	if err := dh.DB.Where("id = ?", id).Where("restaurant_id = ?", restInfo.RestaurantID).First(&updatedProduct).Error; err != nil {
		http.Error(w, "Error Occur on update dish", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&updatedProduct)

}

func (dh *DishHandler) GetDish(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		fmt.Println("this is err get", err)
		http.Error(w, "Internal server ERror", http.StatusBadRequest)
		return
	}

	var dish models.Dish

	if err := dh.DB.First(&dish, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "dish not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server Error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&dish)

}

func (dh *DishHandler) ListDishes(w http.ResponseWriter, r *http.Request) {
	var dishes []models.Dish
	// .Order("created_at desc")
	
	fmt.Println("Hit request")

	var rVars = mux.Vars(r)
	restaurantID, err := uuid.Parse(rVars["restaurantID"])	
	fmt.Println(restaurantID, "restaurantID")

	if err != nil {
			http.Error(w, "Restaurant ID is required", http.StatusBadRequest)
			return
	}

	if err := dh.DB.Where("restaurant_id = ?", restaurantID).Order("created_at desc").Find(&dishes).Error; err != nil {
		http.Error(w, "Internal Server ERror", http.StatusBadRequest)
		return
	}

  // modify image url
	for i := range dishes {

		if len(dishes[i].Images) == 0 {
			continue
		}

		var images []map[string]interface{}

		// 1️⃣ Unmarshal JSON
		if err := json.Unmarshal(dishes[i].Images, &images); err != nil {
			continue
		}

		// 2️⃣ Modify URLs
		for j := range images {
			if url, ok := images[j]["url"].(string); ok {
				smallUrl := strings.Replace(
					url,
					"/image/upload/",
					"/image/upload/w_300,c_limit,q_auto,f_auto/",
					1,
				);
				mediumUrl := strings.Replace(
					url,
					"/image/upload/",
					"/image/upload/w_900,c_limit,q_auto,f_auto/",
					1,
				)
				images[j]["originalUrl"] = url
				images[j]["smallUrl"] = smallUrl
				images[j]["mediumUrl"] = mediumUrl
			}
		}

		// 3️⃣ Marshal back to JSON
		updatedJSON, err := json.Marshal(images)
		if err != nil {
			continue
		}

		dishes[i].Images = updatedJSON
	}

	w.Header().Set("Content-Type", "application/json")
	utils.WriteSuccessResponse(w, "success", http.StatusOK, dishes)
}

func (dh *DishHandler) DeleteDish(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// check user can only deleted his own resource
  restInfo := utils.GetAuthContext(r)

	if err := dh.DB.
    Where("restaurant_id = ?", restInfo.RestaurantID).
    Where("ID = ?", id).
    Delete(&models.Dish{}).Error; err != nil {
	  utils.WriteErrorResponse(w, "Error occur on dish delete", http.StatusBadRequest)
		return
	}

	type DishDeleteRes struct {
		success bool
	}
	utils.WriteSuccessResponse(w, "Dish deleted successfully", http.StatusOK, DishDeleteRes{true})
}

func (dh *DishHandler) ImageUploadHandler(w http.ResponseWriter, r *http.Request) {


	fmt.Println("Image uploader start")

	// Parse multipart form with a max memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	fmt.Println("after parse image error")	

	// Get file from form
	file, handler, err := r.FormFile("file")
	
	fmt.Printf(" this is file %+v", file)

	if err != nil {
		http.Error(w, "Failed to read image from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create uploads directory if not exists
	os.MkdirAll("uploads", os.ModePerm)

	// Create destination file
	dst, err := os.Create(filepath.Join("uploads", handler.Filename))
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to write file to disk", http.StatusInternalServerError)
		return
	}
	
	type Response struct {
		Url string 	`json:"url"`
	}

	fmt.Println("before writing success response handler.Filename", handler.Filename)
	w.WriteHeader(200)
  json.NewEncoder(w).Encode(Response{Url: handler.Filename})	

}
