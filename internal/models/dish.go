package models

import (
	"time"

	"gorm.io/gorm"
)

type Dish struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Rating string    `json:"rating"`
	Price string    `json:"price"`
	DisplayPrice string    `json:"displayPrice"`
	
	RestaurantID      int       `json:"restaurantId"`
	Restaurant        Restaurant      `json:"user" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"` // Automatically managed by GORM
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"` // Automatically managed by GORM

	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"` // Enables soft delete
}
