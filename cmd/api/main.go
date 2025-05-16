package main

import (
	"fmt"
	handlers "gassu/internal/handlers"
	"gassu/internal/models"
	"gassu/internal/utils"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	cwd, _ := os.Getwd()

	err := godotenv.Load(path.Join(cwd, "config", "local.env"))
	if err != nil {
		log.Fatalf("Failed to load evn: %v", err)
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

	r.HandleFunc("/user/get/{id}", restaurantHandler.GetUser).Methods("GET")
	r.HandleFunc("/user/create", restaurantHandler.CreateUser).Methods("POST")
	r.HandleFunc("/user/update/{id}", restaurantHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/delete/{id}", restaurantHandler.DeleteUser).Methods("DELETE")

	r.HandleFunc("/user/sign-up", restaurantHandler.SignupUser).Methods("POST")
	r.HandleFunc("/user/sign-in", restaurantHandler.SigninUser).Methods("POST")

	r.HandleFunc("/blog/create", utils.WithAuth(blogHandler.CreateBlog)).Methods("POST")
	r.HandleFunc("/blog/get/{id}", utils.WithAuth(blogHandler.GetBlog)).Methods("GET")
	r.HandleFunc("/blog/update/{id}", utils.WithAuth(blogHandler.UpdateBlog)).Methods("PUT")
	r.HandleFunc("/blog/delete/{id}", utils.WithAuth(blogHandler.DeleteBlog)).Methods("DELETE")
	r.HandleFunc("/blog/list", utils.WithAuth(blogHandler.ListBlogs)).Methods("GET")

	r.HandleFunc("/dish/add", utils.WithAuth(dishHandler.AddDish)).Methods("POST")
	r.HandleFunc("/dish/update/{id}", utils.WithAuth(dishHandler.UpdateDish)).Methods("PUT")
	r.HandleFunc("/dish/get/{id}", utils.WithAuth(dishHandler.GetDish)).Methods("GET")
	r.HandleFunc("/dish/delete/{id}", utils.WithAuth(dishHandler.DeleteDish)).Methods("DELETE")
	r.HandleFunc("/dish/list", utils.WithAuth(dishHandler.ListDishes)).Methods("GET")

	r.HandleFunc("/protected", utils.WithAuth(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello Protected Route")
	})).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Gorilla Mux!"))
	})

	fmt.Println("Server running at: http://localhost:8000")
	http.ListenAndServe("localhost:8000", r)
}
