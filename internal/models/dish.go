package models

import (
	"time"
	"gorm.io/gorm"
)

type Dish struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string         `json:"name" validate:"required,min=3,max=100" schema:"name"` 
	Description   string         `json:"description" validate:"required" schema:"description"`
	Category      string 			   `json:"category" validate:"required" schema:"category"` 
	PreparationTime string 		   `json:"preparationTime" schema:"preparationTime"`

	Price         float64        `json:"price" validate:"required" schema:"price"`
	DisplayPrice  float64        `json:"displayPrice" schema:"displayPrice"`

	Veg bool `json:"veg" schema:"veg"`

	Rating  float64        `json:"rating" schema:"rating"`
  BestSeller bool 			 `json:"bestSeller" schema:"bestSeller"`
	
	RestaurantID  int            `json:"restaurantId" schema:"restaurantId"`
	Restaurant    Restaurant     `json:"user" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

