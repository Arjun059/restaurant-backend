package utils

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
)

var secretKey = os.Getenv("JWT_SECRET")

func CreateToken( user_id uuid.UUID, user_email string,restaurant_id uuid.UUID, restaurant_url_path string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userEmail": user_email,
			"userId":    user_id,
			"restaurantId": restaurant_id,
			"restaurantURLPath": restaurant_url_path,
			"exp":        time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (map[string]interface{}, error ) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Log the token details
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Printf("Token is valid. Claims: %v", claims)
		return claims, nil
	} else {
		log.Printf("Invalid token or claims: %v", token)
		return nil, errors.New("invalid token or claim")
	}
}
