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
	"os"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"github.com/google/uuid"
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
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(&restaurant)

	fmt.Println("User Get after json response")
}

func (h *RestaurantHandler) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User Updated")
}

func (h *RestaurantHandler) CreateRestaurantAccount(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Hit CreateRestaurantAccount")

	type RestaurantAccountRequest struct {
		User models.User `json:"user"`
		Restaurant models.Restaurant `json:"restaurant"`
	}

	var reqBody RestaurantAccountRequest

	// Decode the request body into 'body'
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		utils.WriteErrorResponse(w, "Body Not Parsed", http.StatusBadRequest)
		return
	}

	// fmt.Printf("%+v req body", reqBody )

	// Check if a user with the same email already exists
	err := h.DB.Where("email = ?", reqBody.User.Email).First(&models.User{}).Error; 
  if	err == nil {
		// if record not found get the err with gorm not found error
	   	utils.WriteErrorResponse(w, "User Already Exists", http.StatusUnprocessableEntity)
	   	return
	}
	
	// hash password
	hashPassword, _ := utils.HashPassword(reqBody.User.Password)
	reqBody.User.Password = hashPassword
	
  insertErr := h.DB.Transaction(func(tx *gorm.DB) error {

		// create restaurant
		reqBody.Restaurant.URLPath = uuid.New().String()
		if err := tx.Create(&reqBody.Restaurant).Error; err != nil {
			log.Printf("Decoded body: %+v\n", reqBody)
			log.Printf("Error creating user: %v", err)
			return err
		}

		// add restaurant id to user
		reqBody.User.RestaurantID = reqBody.Restaurant.ID
		reqBody.User.Role = "owner"
		// Create a new user since no user was found
		if err := tx.Create(&reqBody.User).Error; err != nil {
			// log.Printf("Decoded body reqBody: %+v\n", )
			log.Printf("Error creating user: %v", err)
			return err
		}

		return nil
	})

	if insertErr != nil {
		utils.WriteErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

  siteUrl := os.Getenv("SITE_URL")
	url := fmt.Sprintf("%s/%s", siteUrl,	reqBody.Restaurant.URLPath)
	err = h.DB.Model(&reqBody.Restaurant).Where("id = ?", reqBody.Restaurant.ID).Update("QrCodePath", url).Error
	if err != nil {
		utils.WriteErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Encode the newly created user in the response
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reqBody)

}