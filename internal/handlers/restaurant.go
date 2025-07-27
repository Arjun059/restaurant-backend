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

func (h *RestaurantHandler) GetRestaurantByUrl(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]

	var restaurant models.Restaurant

	fmt.Println("urlpath", url)

	if err := h.DB.Where("url_path = ?", url).First(&restaurant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("Restaurant not found with url  %s", url), http.StatusNotFound)
			return
		}
		fmt.Printf("db error %+v", err)
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}

	utils.WriteSuccessResponse(
		w,
		"success",
		http.StatusOK,
		&restaurant,
	)
}

func (h *RestaurantHandler) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User Updated")
}

func (h *RestaurantHandler) CreateRestaurantAccount(w http.ResponseWriter, r *http.Request) {

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

	var userExist models.User
	err :=  h.DB.Where("email = ?", reqBody.User.Email).First(&userExist).Error

	if err == nil {
			// Email exists 
			utils.WriteErrorResponse(w, "User email already exists", http.StatusBadRequest)
			return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Some other DB error (not "not found") → return 500
			utils.WriteErrorResponse(w, "Internal server error!", http.StatusInternalServerError)
			return
	}

	// fmt.Printf("find user %+v", userExist)

	// hash password
	hashPassword, _ := utils.HashPassword(reqBody.User.Password)
	reqBody.User.Password = hashPassword

	reqBody.Restaurant.URLPath = uuid.New().String()
	
  insertErr := h.DB.Transaction(func(tx *gorm.DB) error {
	  // create restaurant
		if err := tx.Create(&reqBody.Restaurant).Error; err != nil {
			return err
		}

		// add restaurant id to user
		reqBody.User.RestaurantID = reqBody.Restaurant.ID
		reqBody.User.Role = "owner"

		// Create a new owner user for the restaurant
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
	qrURL := fmt.Sprintf("%s/%s/%s/%s", siteUrl, "c", reqBody.Restaurant.URLPath, "home")
  uploadedURL, _ :=	utils.GenerateQrAndUploadToCloud(qrURL, reqBody.Restaurant.URLPath)

	if err := h.DB.Model(&reqBody.Restaurant).Where("id = ?", reqBody.Restaurant.ID).Update("QrCodeURL", uploadedURL).Error; err != nil {
		fmt.Println("error occur on create qrcode")
	}

	// Encode the newly created user in the response
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reqBody)

}