package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"restaurant/internal/models"
	"io"
	"net/http"
	"strconv"
	utils "restaurant/internal/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"os"
	"path/filepath"
	"github.com/gorilla/schema"

	"strings"
	"github.com/google/uuid"
)

var decoder = schema.NewDecoder()

type DishHandler struct {
	DB *gorm.DB;
}


func (dh *DishHandler) AddDish(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// var body models.Dish

	// // Decode the request body and handle any potential errors
	// if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
	// 	utils.WriteErrorResponse(w, fmt.Sprintf("Invalid JSON payload: %v", err), http.StatusBadRequest)
	// 	return
	// }

	// fmt.Println("Add Dish")
	// fmt.Printf("%v", body)
  // if err := utils.ValidateDish(body); err != nil {
	// 	fmt.Printf("%v error in validation ", err)
	// 	utils.WriteErrorResponse(w, fmt.Sprintf("Validation Failed: %v", err) , http.StatusBadRequest)
	// 	return
	// }

	// if err := dh.DB.Create(&body).Error; err != nil {
	// 	fmt.Println("Database error:", err) // Log the actual error for debugging
	// 	utils.WriteErrorResponse(w, "Error occur on Add Dish", http.StatusBadRequest)
	// 	return
	// }

	// fmt.Println("hit before success reponse")

	// // Set content type and return success message
	// w.Header().Set("Content-Type", "text/plain")
	// utils.WriteSuccessResponse(w, "Dish Added Successfully", 201, nil)

	fmt.Printf("content type of request header : %s ", r.Header.Get("content-type"))

  err := r.ParseMultipartForm(100 << 20) // 100MB max memory
	if err != nil {
		fmt.Println("error on parsing request")
		fmt.Printf("%v", err, )
		http.Error(w, "Error parsing multipart form data ", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open file: "+fileHeader.Filename, http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Create uploads directory if not exists
		os.MkdirAll("uploads", os.ModePerm)

		// Create destination file
    ext := filepath.Ext(fileHeader.Filename)
    name := strings.TrimSuffix(fileHeader.Filename , ext)
		uploadedFileName := fmt.Sprintf("%s_%s%s", name, uuid.New().String(), ext)
		dst, err := os.Create(filepath.Join("uploads", uploadedFileName))

		if err != nil {
			http.Error(w, "Could not save file: "+fileHeader.Filename, http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Could not copy file: "+fileHeader.Filename, http.StatusInternalServerError)
			return
		}
	}

	var body models.Dish
	if err := decoder.Decode(&body, r.MultipartForm.Value); err != nil {
		fmt.Printf("error occur on decode %v", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fmt.Printf("hit add dish ----------------- ")
	fmt.Printf("Decoded struct: %+v", body)
	
	utils.WriteSuccessResponse(w, "passed!", http.StatusOK, nil)

}

func (dh *DishHandler) UpdateDish(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var rVars = mux.Vars(r)
	id, err := strconv.Atoi(rVars["id"])

	if err != nil {
		http.Error(w, "Id is invalid", http.StatusBadRequest)
		return
	}

	var body models.Dish
	// Decode the request body and handle any potential errors
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := dh.DB.Model(&models.Dish{}).Where("Id = ? ", id).Updates(&body).Error; err != nil {
		fmt.Println("this is update error ", err)
		http.Error(w, "Error Occur", http.StatusInternalServerError)
		return
	}

	var updatedProduct models.Dish
	if err := dh.DB.First(updatedProduct, id).Error; err != nil {
		http.Error(w, "Error Occur on update dish", http.StatusBadRequest)
		return
	}

	// Set content type and return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&updatedProduct)

}

func (dh *DishHandler) GetDish(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

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

	if err := dh.DB.Find(&products).Error; err != nil {
		http.Error(w, "Internal Server ERror", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.WriteSuccessResponse(w, "success", http.StatusOK, products)
}

func (dh *DishHandler) DeleteDish(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("error not get id", err)
		http.Error(w, "Id Invalid", http.StatusBadRequest)
		return
	}
	if err := dh.DB.Delete(&models.Blog{}, id).Error; err != nil {
		fmt.Println("error not get id", err)
		http.Error(w, "Id Invalid", http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, "Blog Delete Successfully")
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
