package main

import (
	"fmt"
	handlers "restaurant/internal/handlers"
	"restaurant/internal/models"
	utils "restaurant/internal/utils"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	muxHandler "github.com/gorilla/handlers"
)

func main() {
	cwd, _ := os.Getwd()

	env := os.Getenv("APP_ENV")
	fmt.Println("APP_ENV:", env)

	err := godotenv.Load(path.Join(cwd, "config/.env." + env))
	if err != nil {
			log.Fatalf("Failed to load env: %v", err)
	}
	
	db, err := utils.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Auto-migrate  schema's
	db.AutoMigrate(
		&models.Restaurant{},
		&models.Dish{},
		&models.User{},
	)

	restaurantHandler := &handlers.RestaurantHandler{DB: db}
	dishHandler := &handlers.DishHandler{DB: db}
	userHandler := &handlers.UserHandler{DB: db}

	r := mux.NewRouter()

	r.HandleFunc("/qr/generate", restaurantHandler.GenerateQr).Methods("GET")
	
	r.HandleFunc("/user/get/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/user/update/{id}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/delete/{id}", userHandler.DeleteUser).Methods("GET")

	r.HandleFunc("/user/sign-up", userHandler.InviteUser).Methods("POST")
	r.HandleFunc("/user/sign-in", userHandler.SigninUser).Methods("POST")
	r.HandleFunc("/user/update/{id}", userHandler.UpdateUser).Methods("PUT")

	r.HandleFunc("/restaurant/get/{id}", restaurantHandler.GetRestaurant).Methods("GET")
	r.HandleFunc("/restaurant/get/url/{url}", restaurantHandler.GetRestaurantByUrl).Methods("GET")

	r.HandleFunc("/restaurant/register", restaurantHandler.CreateRestaurantAccount).Methods("POST")
	r.HandleFunc("/admin/dashboard/dish/add", utils.WithAuth(dishHandler.AddDish)).Methods("POST")
	r.HandleFunc("/admin/dashboard/dish/update/{id}", utils.WithAuth(dishHandler.UpdateDish)).Methods("PUT")
	r.HandleFunc("/admin/dashboard/dish/get/{id}", utils.WithAuth(dishHandler.GetDish)).Methods("GET")
	r.HandleFunc("/admin/dashboard/dish/delete/{id}", utils.WithAuth(dishHandler.DeleteDish)).Methods("DELETE")
	r.HandleFunc("/admin/dashboard/dish/list", dishHandler.ListDishes).Methods("GET")
	r.HandleFunc("/admin/dashboard/dish/image/upload", dishHandler.ImageUploadHandler).Methods("POST")
	r.HandleFunc("/admin/dashboard/restaurant/update", restaurantHandler.UpdateRestaurant).Methods("PUT")


	r.HandleFunc("/dishes/{restaurantID}", dishHandler.ListDishes).Methods("GET")

	r.HandleFunc("/protected", utils.WithAuth(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello Protected Route")
	})).Methods("GET")

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("!Pong"))
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Gorilla Mux!"))
	})
	
	allowedOrigins := muxHandler.AllowedOrigins([]string{"*"})
	allowedMethods := muxHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders :=  muxHandler.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept", "token"})

	wrappedHandler := muxHandler.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	fmt.Println("Server running at:", fmt.Sprintf("%s:%s", host, port))

	loggedRouter := muxHandler.LoggingHandler(log.Writer(), wrappedHandler)
	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), loggedRouter)

}
