package models

import (
	"time"
	"gorm.io/gorm"
)

type Dish struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string         `json:"name" validate:"required,min=3,max=100"`
	Description   string         `json:"description" validate:"required"`
	Category      string 			   `json:"category" validate:"required"` 
	PreparationTime string 		   `json:"preparationTime"`

	Price         float64        `json:"price" validate:"required"`
	DisplayPrice  float64        `json:"displayPrice"`

	Veg bool `json:"veg"`

	Rating  float64        `json:"rating"`
  BestSeller bool 			 `json:"bestSeller"`
	
	RestaurantID  int            `json:"restaurantId"`
	Restaurant    Restaurant     `json:"user" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

