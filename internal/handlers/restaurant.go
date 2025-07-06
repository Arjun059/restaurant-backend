package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"restaurant/internal/models"
	utils "restaurant/internal/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type RestaurantHandler struct {
	DB *gorm.DB
}

func (h *RestaurantHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}

	var restaurant models.Restaurant

	if err := h.DB.First(&restaurant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("Not Found User With Id  %d", id), http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&restaurant)

	fmt.Println("User Get after json response")
}

func (h *RestaurantHandler) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User Updated")
}

func (h *RestaurantHandler) CreateRestaurantAccount(w http.ResponseWriter, r *http.Request) {

	type RestaurantAccountRequest struct {
		User models.User
		Restaurant models.Restaurant
	}

	var reqBody RestaurantAccountRequest

	// Decode the request body into 'body'
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Body Not Parsed", http.StatusBadRequest)
		return
	}

	// Check if a user with the same email already exists
	if err := h.DB.Where("email = ?", reqBody.User.Email).First(&models.Restaurant{}).Error; err == nil {
		http.Error(w, "User Already Exists", http.StatusUnprocessableEntity)
		return
	}
	
	// hash password
	hashPassword, _ := utils.HashPassword(reqBody.User.Password)
	reqBody.User.Password = hashPassword
	

  insertErr :=	h.DB.Transaction(func(tx *gorm.DB) error {

		// create restaurant
		if err := tx.Create(&reqBody.Restaurant).Error; err != nil {
			log.Printf("Decoded body: %+v\n", reqBody)
			log.Printf("Error creating user: %v", err)
			return err
		}

		// add restaurant id to user
		reqBody.User.RestaurantID = reqBody.Restaurant.ID
		// Create a new user since no user was found
		if err := tx.Create(&reqBody.User).Error; err != nil {
			// log.Printf("Decoded body reqBody: %+v\n", )
			log.Printf("Error creating user: %v", err)
			return err
		}

		return nil
	})

	if insertErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	utils.GenerateQr("http://hello.com")

	// Encode the newly created user in the response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reqBody)

}