package models

import (
	"time"

	"gorm.io/gorm"
)

type Restaurant struct {
  ID       uint   `json:"id" gorm:"primaryKey"`
  Name     string `json:"name"`
  Address  string  `json:"address"`

  QrCodeURL  string `json:"qrCodeURL"` 
  
  URLPath string `json:"urlPath"`

  Users     []User     `json:"users" gorm:"foreignKey:RestaurantID"` // future version add user for restaurant
  Dishes    []Dish     `json:"dishes" gorm:"foreignKey:RestaurantID"`

  CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
  UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
  DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
