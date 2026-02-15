package utils

import (
	"fmt"
	"net/http"
	"context"
	"github.com/google/uuid"
)

type contextKey string

const (
	userIDKey            contextKey = "userID"
	userEmailKey         contextKey = "userEmail"
	restaurantIDKey      contextKey = "restaurantID"
	restaurantURLPathKey contextKey = "restaurantURLPath"
	isRestaurantVerifiedKey contextKey = "isRestaurantVerified"
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

		
		userIDStr := fmt.Sprintf("%v", tokenClaims["userId"])
		restaurantIDStr := fmt.Sprintf("%v", tokenClaims["restaurantId"])

		if err != nil {
			WriteErrorResponse(w, "Auth Failed", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			WriteErrorResponse(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		restaurantID, err := uuid.Parse(restaurantIDStr)
		if err != nil {
			WriteErrorResponse(w, "Invalid restaurant ID", http.StatusUnauthorized)
			return
		}

		userEmail := fmt.Sprintf("%v", tokenClaims["userEmail"])
		restaurantURLPath := fmt.Sprintf("%v", tokenClaims["restaurantURLPath"])
		isRestaurantVerified := tokenClaims["isRestaurantVerified"].(bool)

		// Set into context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, userEmailKey, userEmail)
		ctx = context.WithValue(ctx, restaurantIDKey, restaurantID)
		ctx = context.WithValue(ctx, restaurantURLPathKey, restaurantURLPath)
		ctx = context.WithValue(ctx, isRestaurantVerifiedKey, isRestaurantVerified)

		next(w, r.WithContext(ctx))
	}
}

type AuthContext struct {
	UserID            uuid.UUID
	UserEmail         string
	RestaurantID      uuid.UUID
	RestaurantURLPath string
	IsRestaurantVerified  bool
}

func GetAuthContext(r *http.Request) AuthContext {
	return AuthContext{
		UserID:            getUUIDValue(r.Context(), userIDKey),
		UserEmail:         getStringValue(r.Context(), userEmailKey),
		RestaurantID:      getUUIDValue(r.Context(), restaurantIDKey),
		RestaurantURLPath: getStringValue(r.Context(), restaurantURLPathKey),
		IsRestaurantVerified: getBoolValue(r.Context(), isRestaurantVerifiedKey),
	}
}

func getStringValue(ctx context.Context, key contextKey) string {
	if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	return ""
}

func getUUIDValue(ctx context.Context, key contextKey) uuid.UUID {
	if val, ok := ctx.Value(key).(uuid.UUID); ok {
		return val
	}
	return uuid.Nil
}

func getBoolValue(ctx context.Context, key contextKey) bool {
	if val, ok := ctx.Value(key).(bool); ok {
		return val
	}
	return false
}
