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
	
	if os.Getenv("APP_ENV") != "production" {
			err := godotenv.Load(path.Join(cwd, ".env"))
			if err != nil {
					log.Fatalf("Failed to load env: %v", err)
			}
	}

	db, err := utils.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Auto-migrate  schema's
	db.AutoMigrate(
		&models.Restaurant{},
		&models.Blog{},
		&models.Dish{},
	)

	restaurantHandler := &handlers.RestaurantHandler{DB: db}
	blogHandler := &handlers.BlogHandler{DB: db}
	dishHandler := &handlers.DishHandler{DB: db}

	r := mux.NewRouter()

	r.HandleFunc("/restaurant/get/{id}", restaurantHandler.GetUser).Methods("GET")
	r.HandleFunc("/restaurant/create", restaurantHandler.CreateUser).Methods("POST")
	r.HandleFunc("/restaurant/update/{id}", restaurantHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/restaurant/delete/{id}", restaurantHandler.DeleteUser).Methods("DELETE")

	r.HandleFunc("/restaurant/sign-up", restaurantHandler.SignupUser).Methods("POST")
	r.HandleFunc("/restaurant/sign-in", restaurantHandler.SigninUser).Methods("POST")

	r.HandleFunc("/blog/create", utils.WithAuth(blogHandler.CreateBlog)).Methods("POST")
	r.HandleFunc("/blog/get/{id}", utils.WithAuth(blogHandler.GetBlog)).Methods("GET")
	r.HandleFunc("/blog/update/{id}", utils.WithAuth(blogHandler.UpdateBlog)).Methods("PUT")
	r.HandleFunc("/blog/delete/{id}", utils.WithAuth(blogHandler.DeleteBlog)).Methods("DELETE")
	r.HandleFunc("/blog/list", utils.WithAuth(blogHandler.ListBlogs)).Methods("GET")

	// r.HandleFunc("/admin/dashboard/dish/add", utils.WithAuth(dishHandler.AddDish)).Methods("POST")
	r.HandleFunc("/admin/dashboard/dish/add", dishHandler.AddDish).Methods("POST")
	r.HandleFunc("/admin/dashboard/dish/update/{id}", utils.WithAuth(dishHandler.UpdateDish)).Methods("PUT")
	r.HandleFunc("/admin/dashboard/dish/get/{id}", utils.WithAuth(dishHandler.GetDish)).Methods("GET")
	r.HandleFunc("/admin/dashboard/dish/delete/{id}", utils.WithAuth(dishHandler.DeleteDish)).Methods("DELETE")
	r.HandleFunc("/admin/dashboard/dish/list", dishHandler.ListDishes).Methods("GET")
	r.HandleFunc("/admin/dashboard/dish/image/upload", dishHandler.ImageUploadHandler).Methods("POST")

	r.HandleFunc("/dishes", dishHandler.ListDishes).Methods("GET")

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


	// cron job for ping service
	utils.PreventSleepCron()
	
	allowedCorsObj := muxHandler.AllowedOrigins([]string{"*"})
	allowedMethods := muxHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders :=  muxHandler.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept"})

	wrappedHandler := muxHandler.CORS(allowedCorsObj, allowedMethods, allowedHeaders)(r)
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	fmt.Println("Server running at:", fmt.Sprintf("%s:%s", host, port))

	loggedRouter := muxHandler.LoggingHandler(log.Writer(), wrappedHandler)
	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), loggedRouter)

}
