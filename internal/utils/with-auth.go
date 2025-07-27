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
	RestaurantURLPathKey  StringContextKey = ""
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
			WriteErrorResponse(w, "Auth Failed", http.StatusUnauthorized)
			return
		}

		userIDFloat, _ := tokenClaims["userId"].(float64)
		restaurantIDFloat, _ := tokenClaims["restaurantId"].(float64)

		userID := uint(userIDFloat)
		restaurantID := uint(restaurantIDFloat)
		userEmail := fmt.Sprintf("%v", tokenClaims["userEmail"])
		restaurantURLPath := fmt.Sprintf("%v", tokenClaims["restaurantURLPath"])

		// Set into context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserEmailKey, userEmail)
		ctx = context.WithValue(ctx, RestaurantIDKey, restaurantID)
		ctx = context.WithValue(ctx, RestaurantURLPathKey, restaurantURLPath)

		next(w, r.WithContext(ctx))
	}
}
type AuthContext struct {
	UserID uint
	UserEmail string
	RestaurantID uint
	RestaurantURLPath string
}

func GetAuthContext(r *http.Request) AuthContext {
	authContext := AuthContext{}
	authContext.UserID, _ = r.Context().Value(UserIDKey).(uint)
	authContext.UserEmail, _ = r.Context().Value(UserEmailKey).(string)
	authContext.RestaurantID, _ = r.Context().Value(RestaurantIDKey).(uint)
	authContext.RestaurantURLPath, _ = r.Context().Value(RestaurantURLPathKey).(string)
	return authContext
}

