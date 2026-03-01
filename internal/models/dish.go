package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type FileMeta struct{
		Name string `json:"name"`
		Folder string `json:"folder"`
		Url string `json:"url"`
}

type Dish struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name          string         `json:"name" validate:"required,min=3,max=100" schema:"name"` 
	Description   string         `json:"description" validate:"required" schema:"description"`
	
	Categories datatypes.JSON    `json:"categories"`

	PreparationTime string 		   `json:"preparationTime" schema:"preparationTime"`

	// Price is only required when no variants are defined
	Price         float64        `json:"price" validate:"omitempty,gt=0" schema:"price"`
	DisplayPrice  float64        `json:"displayPrice" schema:"displayPrice"`

  Serving     string            `json:"serving"` // e.g., "1 person", "2 people"

	Veg bool                     `json:"veg" schema:"veg"`

	Rating  float64              `json:"rating" schema:"rating"`
  BestSeller bool 		       	 `json:"bestSeller" schema:"bestSeller"`

	Images datatypes.JSON        `json:"images"`

	UserID    uuid.UUID          `json:"userId" schema:"userId"`
	User    User                 `json:"user" gorm:"foreignKey:UserID;references:ID"`
	
	RestaurantID   uuid.UUID     `json:"restaurantId" schema:"restaurantId"`
	Restaurant    Restaurant     `json:"restaurant" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	
	CreatedAt     time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`

  IsAvailable bool             `json:"isAvailable" gorm:"default:true"`

  Variants    []DishVariant    `json:"variants" gorm:"foreignKey:DishID;references:ID"`

}
