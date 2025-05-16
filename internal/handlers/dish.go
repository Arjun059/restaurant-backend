package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gassu/internal/models"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type DishHandler struct {
	DB *gorm.DB
}

func (ph *DishHandler) AddDish(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body models.Dish

	// Decode the request body and handle any potential errors
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := ph.DB.Create(&body).Error; err != nil {
		fmt.Println("Database error:", err) // Log the actual error for debugging
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set content type and return success message
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Blog Created")
}

func (ph *DishHandler) UpdateDish(w http.ResponseWriter, r *http.Request) {
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

	if err := ph.DB.Model(&models.Dish{}).Where("Id = ? ", id).Updates(&body).Error; err != nil {
		fmt.Println("this is update error ", err)
		http.Error(w, "Error Occur", http.StatusInternalServerError)
		return
	}

	var updatedProduct models.Dish
	if err := ph.DB.First(updatedProduct, id).Error; err != nil {
		http.Error(w, "Error Occur on update dish", http.StatusBadRequest)
		return
	}

	// Set content type and return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&updatedProduct)

}

func (ph *DishHandler) GetDish(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		fmt.Println("this is err get", err)
		http.Error(w, "Internal server ERror", http.StatusBadRequest)
		return
	}

	var dish models.Dish

	if err := ph.DB.First(&dish, id).Error; err != nil {
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

func (ph *DishHandler) ListDishes(w http.ResponseWriter, r *http.Request) {
	var products []models.Dish

	if err := ph.DB.Find(&products).Error; err != nil {
		http.Error(w, "Internal Server ERror", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&products)

}

func (ph *DishHandler) DeleteDish(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("errror not get id", err)
		http.Error(w, "Id Invalid", http.StatusBadRequest)
		return
	}
	if err := ph.DB.Delete(&models.Blog{}, id).Error; err != nil {
		fmt.Println("error not get id", err)
		http.Error(w, "Id Invalid", http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, "Blog Delete Successfully")
}
