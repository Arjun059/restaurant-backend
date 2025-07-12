package utils

import (
	"fmt"
	"net/http"
	"context"
)


type StringContextKey string
type IntContextKey uint

const (
	UserIDKey        IntContextKey = 0
	UserEmailKey     StringContextKey = ""
	RestaurantIDKey  IntContextKey = 0
)

// Inline middleware example (logging)
func WithAuth(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}
		tokenString = tokenString[len("Bearer "):]

		tokenClaims, err := VerifyToken(tokenString)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid token")
			return
		}

		userIDFloat, _ := tokenClaims["userId"].(float64)
		restaurantIDFloat, _ := tokenClaims["restaurantId"].(float64)

		userID := uint(userIDFloat)
		restaurantID := uint(restaurantIDFloat)
		userEmail := fmt.Sprintf("%v", tokenClaims["userEmail"])

		// Set into context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserEmailKey, userEmail)
		ctx = context.WithValue(ctx, RestaurantIDKey, restaurantID)

		next(w, r.WithContext(ctx))
	}
}

func GetAuthContext(r *http.Request) (userID uint, userEmail string, restaurantID uint) {
	userID = r.Context().Value(UserIDKey).(uint)
	userEmail = r.Context().Value(UserEmailKey).(string)
	restaurantID = r.Context().Value(RestaurantIDKey).(uint)
	return
}

