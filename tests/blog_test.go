package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"restaurant/internal/handlers"
	"restaurant/internal/models"
	"restaurant/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// only that function is include in test that start with Test[RestName] name

// individual testing
func GetBlogHandler() (*handlers.BlogHandler, error) {
	db, err := utils.InitDB()

	if err != nil {
		return nil, errors.New("Error Occur on db Init")
	}

	db.AutoMigrate(models.Blog{})

	blogHandler := &handlers.BlogHandler{DB: db}

	return blogHandler, nil
}

func TestCreateBlog(t *testing.T) {
	blogHandler, _ := GetBlogHandler()

	body, err := json.Marshal(models.Blog{
		Title: "hello world",
	})

	if err != nil {
		t.Fatalf("Error occur on parsing data")
	}

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	blogHandler.CreateBlog(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Print log message (only visible in -v mode)
	t.Log("Success: Create Blog Test")

	var responseBlog models.Blog
	json.Unmarshal(w.Body.Bytes(), &responseBlog)

	if responseBlog.Title != "hello world" {
		t.Errorf("Expected title hello world, got %s", responseBlog.Title)
	}

}

func TestGetBlog(t *testing.T) {
	blogHandler, _ := GetBlogHandler()

	req := httptest.NewRequest("GET", "/blog/get/18", nil)
	req.Header.Set("content-type", "application/json")

	// Manually set the URL variable for `id` if using gorilla/mux
	req = mux.SetURLVars(req, map[string]string{"id": "18"})

	w := httptest.NewRecorder()

	blogHandler.GetBlog(w, req)

	if w.Code != http.StatusOK {
		t.Error("Code is not same ")
	}

	var response models.Blog

	json.Unmarshal(w.Body.Bytes(), &response)

	t.Logf("reponse %+v", response)

}

// Subtests  testing
func TestBlogApi(t *testing.T) {
	db, err := utils.InitDB()
	if err != nil {
		t.Error("Error Occur on db init")
	}

	db.AutoMigrate(models.Blog{})

	blogHandler := &handlers.BlogHandler{DB: db}

	t.Run("Get Blog", func(it *testing.T) {
		// ready req
		req := httptest.NewRequest("GET", "/blog/get/18", nil)
		req.Header.Set("content-type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": "18"})

		// ready w
		w := httptest.NewRecorder()

		// hit req
		blogHandler.GetBlog(w, req)

		/// process further
		if w.Code != http.StatusOK {
			it.Error("Blog Not Get")
		}

		var response models.Blog
		json.Unmarshal(w.Body.Bytes(), &response)

	})

	t.Run("Create Blog", func(it *testing.T) {

		body, err := json.Marshal(models.Blog{
			Title: "hello world",
		})

		if err != nil {
			t.Fatal("error on converting json")
		}

		// ready req
		req := httptest.NewRequest("POST", "/blog/create", bytes.NewBuffer(body))
		req.Header.Set("content-type", "application/json")

		// ready w
		w := httptest.NewRecorder()

		// hit req
		blogHandler.CreateBlog(w, req)

		// process further
		if w.Code != http.StatusCreated {
			it.Errorf("Blog Not Created")
		}

		// var response models.Blog
		// json.Unmarshal(w.Body.Bytes(), &response)

	})

}
