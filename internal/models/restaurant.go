package models

import (
	"time"

	"gorm.io/gorm"
)

type Restaurant struct {
  ID       uint   `gorm:"primaryKey"`
  Name     string
  Address  string

  Email    string `json:"email"`
	Password string `json:"password"`
 
  Users     []User     `gorm:"foreignKey:RestaurantID"` // future version add user for restaurant
  Dishes    []Dish     `gorm:"foreignKey:RestaurantID"`

  CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
  UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
  DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
