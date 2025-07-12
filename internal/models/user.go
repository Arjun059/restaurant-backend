package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	
	Email    string `json:"email"`
	Password string `json:"password"`

	Role string `json:"role"` // New field to define employee role (e.g., "owner", "manager", "editor")

	RestaurantID uint       `json:"restaurantId"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"` // Enables soft delete
}


