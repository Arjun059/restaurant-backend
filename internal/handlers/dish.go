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

	// Remove categories from the form values so decoder doesn't choke
	if r.MultipartForm != nil {
		delete(r.MultipartForm.Value, "categories")
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

	if len(filesMeta) > 0 {
		if imagesJson, err := json.Marshal(filesMeta); err == nil {
			body.Images = datatypes.JSON(imagesJson)
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
	var products []models.Dish
	// .Order("created_at desc")

	var rVars = mux.Vars(r)
	restaurantID, err := uuid.Parse(rVars["restaurantID"])	
	fmt.Println(restaurantID, "restaurantID")

	if err != nil {
			http.Error(w, "Restaurant ID is required", http.StatusBadRequest)
	}

	if err := dh.DB.Where("restaurant_id = ?", restaurantID).Order("created_at desc").Find(&products).Error; err != nil {
		http.Error(w, "Internal Server ERror", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.WriteSuccessResponse(w, "success", http.StatusOK, products)
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
