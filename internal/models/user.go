package models

import (
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string `json:"name"`
	
	Email    string `json:"email"`
	Password string `json:"password"`

	Role string `json:"role"` // New field to define employee role (e.g., "owner", "manager", "editor", "super-admin")

	RestaurantID  uuid.UUID `json:"restaurantId"`
	Restaurant   Restaurant `json:"restaurant" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"` // Enables soft delete
}


