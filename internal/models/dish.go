package models

import (
	"time"
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type FileMeta struct{
		Name string `json:"name"`
		Folder string `json:"folder"`
		Url string `json:"url"`
}

type Dish struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string         `json:"name" validate:"required,min=3,max=100" schema:"name"` 
	Description   string         `json:"description" validate:"required" schema:"description"`
	
	Categories datatypes.JSON    `json:"categories"`

	PreparationTime string 		   `json:"preparationTime" schema:"preparationTime"`

	Price         float64        `json:"price" validate:"required" schema:"price"`
	DisplayPrice  float64        `json:"displayPrice" schema:"displayPrice"`

	Veg bool                     `json:"veg" schema:"veg"`

	Rating  float64              `json:"rating" schema:"rating"`
  BestSeller bool 		       	 `json:"bestSeller" schema:"bestSeller"`

	Images datatypes.JSON        `json:"images"`

	UserID  uint                  `json:"userId" schema:"userId"`
	User    User                 `json:"user" gorm:"foreignKey:UserID;references:ID"`
	
	RestaurantID  uint            `json:"restaurantId" schema:"restaurantId"`
	Restaurant    Restaurant     `json:"restaurant" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

