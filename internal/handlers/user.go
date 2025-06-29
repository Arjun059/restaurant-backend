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

type UserHandler struct {
	DB *gorm.DB
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User Created")
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}

	var user models.User

	if err := h.DB.Preload("Restaurant").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("Not Found User With Id  %d", id), http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)

	fmt.Println("User Get after json response")
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User Updated")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User Deleted")
}

func (h *UserHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	// change code [ no sign up flow in the current auth flow ]

	var user *models.Restaurant
	var body models.Restaurant

	// Decode the request body into 'body'
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Body Not Parsed", http.StatusBadRequest)
		return
	}

	// Check if a user with the same email already exists
	if err := h.DB.Where("email = ?", body.Email).First(&user).Error; err == nil {
		http.Error(w, "User Already Exists", http.StatusUnprocessableEntity)
		return
	}
	
	// hash password
	hashPassword, _ := utils.HashPassword(body.Password)
	body.Password = hashPassword
	
	// Create a new user since no user was h found
	if err := h.DB.Create(&body).Error; err != nil {
		log.Printf("Decoded body: %+v\n", body)
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Encode the newly created user in the response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

func (h *UserHandler) SigninUser(w http.ResponseWriter, r *http.Request) {

	var body models.Restaurant
	json.NewDecoder(r.Body).Decode(&body)

	var user models.Restaurant

	log.Printf("Decoded body: %+v\n", body)

	if e := h.DB.Where("email = ?", body.Email).First(&user).Error; e != nil {
		log.Printf("Decoded body: %+v\n", e)
	
		utils.WriteErrorResponse(w, "User not found", http.StatusNotFound)
		return 

	}

	if user.Email != body.Email || !utils.CheckPasswordHash(body.Password, user.Password) {
		utils.WriteErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	tokenString, err := utils.CreateToken(user.Email, user.ID)
	if err != nil {
		utils.WriteErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	type Response struct {
		Token   string `json:"token"`
		Restaurant models.Restaurant `json:"restaurant"`
	}
	var response Response
	response.Token = tokenString
	response.Restaurant = user

	utils.WriteSuccessResponse(w, "Login success" , http.StatusOK, response)
}

