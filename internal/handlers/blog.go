package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"restaurant/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type BlogHandler struct {
	DB *gorm.DB
}

func (bh *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body models.Blog

	// Decode the request body and handle any potential errors
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("Error:", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := bh.DB.Create(&body).Error; err != nil {
		fmt.Println("Database error:", err) // Log the actual error for debugging
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set content type and return success message
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("Blog Created"))

	// Yes! When GORM's Create(&body) method is executed,
	// it inserts the new blog entry into the database
	// and updates the body variable with any fields that the database automatically generates
	// (such as ID, timestamps, default values, etc.)

	// The request body (req.Body) is parsed into body:
	// The blog entry is created in the database:
	// GORM inserts the record.
	// After insertion, GORM retrieves the newly created row from the database.
	// body is updated with database-generated fields like:
	// ID (if it's an auto-incrementing primary key)
	// CreatedAt (if GORM is managing timestamps)
	// Any default values from the database schema.

	// Now, the response contains both the original data and any new fields set by the DB.

	json.NewEncoder(w).Encode(&body)

}

func (bh *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var rVars = mux.Vars(r)
	id, err := strconv.Atoi(rVars["id"])

	if err != nil {
		http.Error(w, "Id is invalid", http.StatusBadRequest)
		return
	}

	var body models.Blog
	// Decode the request body and handle any potential errors
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := bh.DB.Model(&models.Blog{}).Where("Id = ? ", id).Updates(&body).Error; err != nil {
		fmt.Println("this is update error ", err)
		http.Error(w, "Error Occur", http.StatusInternalServerError)
		return
	}

	var updatedBlog models.Blog

	if err := bh.DB.Preload("User").First(&updatedBlog, id).Error; err != nil {
		fmt.Println("this is update error ", err)
		http.Error(w, "Error Occur", http.StatusInternalServerError)
		return
	}

	// Set content type and return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&updatedBlog)

}

func (bh *BlogHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		fmt.Println("this is err get", err)
		http.Error(w, "Internal server ERror", http.StatusBadRequest)
		return
	}

	var blog models.Blog

	if err := bh.DB.Preload("User").First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Blog not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server Error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&blog)

}

func (bh *BlogHandler) ListBlogs(w http.ResponseWriter, r *http.Request) {
	var blogs []models.Blog

	if err := bh.DB.Find(&blogs).Error; err != nil {
		http.Error(w, "Internal server ERror", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&blogs)

}

func (bh *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("errror not get id", err)
		http.Error(w, "Id Invalid", http.StatusBadRequest)
		return
	}
	if err := bh.DB.Delete(&models.Blog{}, id).Error; err != nil {
		fmt.Println("errror not get id", err)
		http.Error(w, "Id Invalid", http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, "Blog Delete Successfully")
}
